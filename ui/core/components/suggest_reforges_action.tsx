import clsx from 'clsx';
import tippy, { hideAll } from 'tippy.js';
import { ref } from 'tsx-vanilla';
import { Constraint, greaterEq, lessEq, Model, Options, Solution, solve } from 'yalps';

import * as Mechanics from '../constants/mechanics.js';
import { IndividualSimUI } from '../individual_sim_ui';
import { Player } from '../player';
import { Class, ItemSlot, PseudoStat, ReforgeStat, Spec, Stat } from '../proto/common';
import { IndividualSimSettings, StatCapType } from '../proto/ui';
import { ReforgeData } from '../proto_utils/equipped_item';
import { Gear } from '../proto_utils/gear';
import { shortSecondaryStatNames, slotNames, statCapTypeNames } from '../proto_utils/names';
import { pseudoStatIsCapped, StatCap, Stats, UnitStat, UnitStatPresets } from '../proto_utils/stats';
import { SpecTalents } from '../proto_utils/utils';
import { Sim } from '../sim';
import { ActionGroupItem } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event';
import { isDevMode, sleep } from '../utils';
import { CopyButton } from './copy_button';
import { BooleanPicker } from './pickers/boolean_picker';
import { EnumPicker } from './pickers/enum_picker';
import { NumberPicker, NumberPickerConfig } from './pickers/number_picker';
import { renderSavedEPWeights } from './saved_data_managers/ep_weights';
import Toast from './toast';

type YalpsCoefficients = Map<string, number>;
type YalpsVariables = Map<string, YalpsCoefficients>;
type YalpsConstraints = Map<string, Constraint>;

const INCLUDED_STATS = [
	Stat.StatHitRating,
	Stat.StatCritRating,
	Stat.StatHasteRating,
	Stat.StatExpertiseRating,
	Stat.StatMasteryRating,
	Stat.StatDodgeRating,
	Stat.StatParryRating,
];

type StatTooltipContent = { [key in Stat]?: () => Element | string };

const STAT_TOOLTIPS: StatTooltipContent = {
	[Stat.StatMasteryRating]: () => (
		<>
			Total <strong>percentage</strong>
		</>
	),
	[Stat.StatHasteRating]: () => (
		<>
			Final percentage value <strong>including</strong> all buffs/gear.
		</>
	),
};

export type ReforgeOptimizerOptions = {
	experimental?: true;
	statTooltips?: StatTooltipContent;
	statSelectionPresets?: UnitStatPresets[];
	// Allows you to enable breakpoint limits for Treshold type caps
	enableBreakpointLimits?: boolean;
	// Allows you to modify the stats before they are returned for the calculations
	// For example: Adding class specific Glyphs/Talents that are not added by the backend
	updateGearStatsModifier?: (baseStats: Stats) => Stats;
	// Allows you to get alternate default EPs
	// For example for Fury where you have SMF and TG EPs
	getEPDefaults?: (player: Player<any>) => Stats;
	// Allows you to modify default softCaps
	// For example you wish to add breakpoints for Berserking / Bloodlust if enabled
	updateSoftCaps?: (softCaps: StatCap[]) => StatCap[];
	// Allows you to specifiy additional information for the soft cap tooltips
	additionalSoftCapTooltipInformation?: StatTooltipContent;
};

// Used to force a particular proc from trinkets like Matrix Restabilizer and Apparatus of Khaz'goroth.
class RelativeStatCap {
	static relevantStats: Stat[] = [Stat.StatCritRating, Stat.StatHasteRating, Stat.StatMasteryRating];
	readonly forcedHighestStat: UnitStat;
	readonly constrainedStats: UnitStat[];
	readonly constraintKeys: string[];

	// Not comprehensive, add any other relevant offsets here as needed.
	static procTrinketOffsets: Map<Stat, Map<number, number>> = new Map([
		[
			Stat.StatCritRating,
			new Map([
				[69167, 460], // Vessel of Acceleration (H)
				[68995, 410], // Vessel of Acceleration (N)
			]),
		],
		[
			Stat.StatHasteRating,
			new Map([
				[69112, 1730], // The Hungerer (H)
				[68927, 1532], // The Hungerer (N)
			]),
		],
		[Stat.StatMasteryRating, new Map([])],
	]);

	static canEnable(player: Player<any>): boolean {
		const variableStatTrinkets: number[] = [69150, 68994, 69113, 68972];
		return player.getGear().hasTrinketFromOptions(variableStatTrinkets);
	}

	constructor(forcedHighestStat: Stat, playerClass: Class) {
		if (!RelativeStatCap.relevantStats.includes(forcedHighestStat)) {
			throw new Error('Forced highest stat must be either Crit, Haste, or Mastery!');
		}

		this.forcedHighestStat = UnitStat.fromStat(forcedHighestStat);
		this.constrainedStats = RelativeStatCap.relevantStats.filter(stat => stat !== forcedHighestStat).map(stat => UnitStat.fromStat(stat));
		this.constraintKeys = this.constrainedStats.map(
			unitStat => this.forcedHighestStat.getShortName(playerClass) + 'Minus' + unitStat.getShortName(playerClass),
		);
	}

	updateCoefficients(coefficients: YalpsCoefficients, stat: Stat, amount: number) {
		if (!RelativeStatCap.relevantStats.includes(stat)) {
			return;
		}

		for (const [idx, constrainedStat] of this.constrainedStats.entries()) {
			const coefficientKey = this.constraintKeys[idx];
			const currentValue = coefficients.get(coefficientKey) || 0;

			if (this.forcedHighestStat.equalsStat(stat)) {
				coefficients.set(coefficientKey, currentValue + amount);
			} else if (constrainedStat.equalsStat(stat)) {
				coefficients.set(coefficientKey, currentValue - amount);
			}
		}
	}

	updateConstraints(constraints: YalpsConstraints, gear: Gear, baseStats: Stats) {
		for (const [idx, constrainedStat] of this.constrainedStats.entries()) {
			const weightedStatsArray = new Stats().withUnitStat(this.forcedHighestStat, 1).withUnitStat(constrainedStat, -1);
			let minReforgeContribution = 1 - baseStats.computeEP(weightedStatsArray);
			const procOffsetMap = RelativeStatCap.procTrinketOffsets.get(constrainedStat.getStat())!;

			for (const trinket of gear.getTrinkets()) {
				if (!trinket) {
					continue;
				}

				const trinketId = trinket.item.id;

				if (procOffsetMap.has(trinketId)) {
					minReforgeContribution += procOffsetMap.get(trinketId)!;
					break;
				}
			}

			constraints.set(this.constraintKeys[idx], greaterEq(minReforgeContribution));
		}
	}
}

export class ReforgeOptimizer {
	protected readonly simUI: IndividualSimUI<any>;
	protected readonly player: Player<any>;
	protected readonly playerClass: Class;
	protected readonly isExperimental: ReforgeOptimizerOptions['experimental'];
	protected readonly isHybridCaster: boolean;
	protected readonly sim: Sim;
	protected readonly defaults: IndividualSimUI<any>['individualConfig']['defaults'];
	protected getEPDefaults: ReforgeOptimizerOptions['getEPDefaults'];
	protected _statCaps: Stats;
	protected updateGearStatsModifier: ReforgeOptimizerOptions['updateGearStatsModifier'];
	protected _softCapsConfig: StatCap[];
	protected updateSoftCaps: ReforgeOptimizerOptions['updateSoftCaps'];
	protected enableBreakpointLimits: ReforgeOptimizerOptions['enableBreakpointLimits'];
	protected statTooltips: StatTooltipContent = {};
	protected additionalSoftCapTooltipInformation: StatTooltipContent = {};
	protected statSelectionPresets: ReforgeOptimizerOptions['statSelectionPresets'];
	readonly freezeItemSlotsChangeEmitter = new TypedEvent<void>();
	protected freezeItemSlots = false;
	protected frozenItemSlots = new Map<ItemSlot, boolean>();
	protected previousGear: Gear | null = null;
	protected previousReforges = new Map<ItemSlot, ReforgeData>();
	protected currentReforges = new Map<ItemSlot, ReforgeData>();
	protected relativeStatCap: RelativeStatCap | null = null;

	constructor(simUI: IndividualSimUI<any>, options?: ReforgeOptimizerOptions) {
		this.simUI = simUI;
		this.player = simUI.player;
		this.playerClass = this.player.getClass();
		this.isExperimental = options?.experimental;
		this.isHybridCaster = [Spec.SpecBalanceDruid, Spec.SpecShadowPriest, Spec.SpecElementalShaman, Spec.SpecMistweaverMonk].includes(this.player.getSpec());
		this.sim = simUI.sim;
		this.defaults = simUI.individualConfig.defaults;
		this.getEPDefaults = options?.getEPDefaults;
		this.updateSoftCaps = options?.updateSoftCaps;
		this.updateGearStatsModifier = options?.updateGearStatsModifier;
		this._softCapsConfig = this.defaults.softCapBreakpoints || [];
		this.statTooltips = { ...STAT_TOOLTIPS, ...options?.statTooltips };
		this.additionalSoftCapTooltipInformation = { ...options?.additionalSoftCapTooltipInformation };
		this.statSelectionPresets = options?.statSelectionPresets;
		this._statCaps = this.statCaps;
		this.enableBreakpointLimits = !!options?.enableBreakpointLimits;

		const startReforgeOptimizationEntry: ActionGroupItem = {
			label: 'Suggest Reforges',
			cssClass: 'suggest-reforges-action-button flex-grow-1',
			onClick: async ({ currentTarget }) => {
				const button = currentTarget as HTMLButtonElement;
				if (button) {
					button.classList.add('loading');
					button.disabled = true;
				}
				try {
					performance.mark('reforge-optimization-start');
					await this.optimizeReforges();
					this.onReforgeDone();
				} catch (error) {
					this.onReforgeError(error);
				} finally {
					performance.mark('reforge-optimization-end');
					if (isDevMode())
						console.log(
							'Reforge optimization took:',
							`${performance
								.measure('reforge-optimization-measure', 'reforge-optimization-start', 'reforge-optimization-end')
								.duration.toFixed(2)}ms`,
						);
					if (button) {
						button.classList.remove('loading');
						button.disabled = false;
					}
				}
			},
		};

		const contextMenuEntry: ActionGroupItem = {
			cssClass: 'suggest-reforges-button-settings',
			children: (
				<>
					<i className="fas fa-cog" />
				</>
			),
		};

		const {
			group,
			children: [startReforgeOptimizationButton, contextMenuButton],
		} = simUI.addActionGroup([startReforgeOptimizationEntry, contextMenuEntry], {
			cssClass: clsx('suggest-reforges-settings-group d-flex', this.isExperimental && !this.player.sim.getShowExperimental() && 'hide'),
		});

		this.bindToggleExperimental(group);

		if (!!this.softCapsConfig?.length)
			tippy(startReforgeOptimizationButton, {
				theme: 'suggest-reforges-softcaps',
				placement: 'bottom',
				maxWidth: 310,
				interactive: true,
				onShow: instance => instance.setContent(this.buildReforgeButtonTooltip()),
			});

		tippy(contextMenuButton, {
			placement: 'bottom',
			content: 'Change Reforge Optimizer settings',
		});

		this.buildContextMenu(contextMenuButton);
	}

	private bindToggleExperimental(element: Element) {
		const toggle = () => element.classList[this.isExperimental && !this.player.sim.getShowExperimental() ? 'add' : 'remove']('hide');
		toggle();
		this.player.sim.showExperimentalChangeEmitter.on(() => {
			toggle();
		});
	}

	get softCapsConfig() {
		return this.updateSoftCaps?.(StatCap.cloneSoftCaps(this._softCapsConfig)) || this._softCapsConfig;
	}

	get softCapsConfigWithLimits() {
		if (!this.enableBreakpointLimits) return this.softCapsConfig;

		const softCaps = StatCap.cloneSoftCaps(this.softCapsConfig);
		for (const [unitStat, limit] of this.player.getBreakpointLimits().asUnitStatArray()) {
			if (!limit) continue;
			const config = softCaps.find(config => config.unitStat.equals(unitStat));
			if (config) config.breakpoints = config.breakpoints.filter(breakpoint => breakpoint <= limit);
		}
		return softCaps;
	}

	get statCaps() {
		return this.sim.getUseCustomEPValues() ? this.player.getStatCaps() : this.defaults.statCaps || new Stats();
	}
	setStatCap(unitStat: UnitStat, value: number) {
		this._statCaps = this._statCaps.withUnitStat(unitStat, value);
		if (this.sim.getUseCustomEPValues()) {
			this.player.setStatCaps(TypedEvent.nextEventID(), this._statCaps);
		}
		return this.statCaps;
	}
	setDefaultStatCaps() {
		this._statCaps = this.defaults.statCaps || new Stats();
		this.player.setStatCaps(TypedEvent.nextEventID(), this._statCaps);
		return this.statCaps;
	}

	get preCapEPs(): Stats {
		let weights = this.sim.getUseCustomEPValues() ? this.player.getEpWeights() : this.getEPDefaults?.(this.player) || this.defaults.epWeights;

		// Replace Spirit EP for hybrid casters with a small value in order to break ties between Spirit and Hit Reforges
		if (this.isHybridCaster) {
			weights = weights.withStat(Stat.StatSpirit, 0.01);
		}

		return weights;
	}

	// Checks that school-specific weights for Rating stats are set whenever there is a school-specific stat cap configured, and ensures that the
	// EPs for such stats are not double counted.
	static checkWeights(weights: Stats, reforgeCaps: Stats, reforgeSoftCaps: StatCap[]): Stats {
		let validatedWeights = weights;

		// Loop through Hit/Crit/Haste pure Rating stats.
		for (const parentStat of [Stat.StatHitRating, Stat.StatCritRating, Stat.StatHasteRating]) {
			const children = UnitStat.getChildren(parentStat);
			const specificSchoolWeights = children.map(childStat => weights.getPseudoStat(childStat));

			// If any of the children have non-zero EP, then set pure Rating EP
			// to 0 and continue.
			if (specificSchoolWeights.some(weight => weight !== 0)) {
				validatedWeights = validatedWeights.withStat(parentStat, 0);
				continue;
			}

			// If all children have 0 EP, then loop through children and check whether a cap has been configured for that child.
			for (const childStat of children) {
				if (pseudoStatIsCapped(childStat, reforgeCaps, reforgeSoftCaps)) {
					// The first time a cap is detected, set EP for that child to re-scaled parent Rating EP, set parent Rating EP
					// to 0, and break.
					const rescaledWeight = UnitStat.fromPseudoStat(childStat).convertPercentToRating(weights.getStat(parentStat));
					validatedWeights = validatedWeights.withPseudoStat(childStat, rescaledWeight!);
					validatedWeights = validatedWeights.withStat(parentStat, 0);
					break;
				}
			}
		}

		return validatedWeights;
	}

	buildReforgeButtonTooltip() {
		return (
			<>
				<p>The following breakpoints have been implemented for this spec:</p>
				<table className="w-100">
					<tbody>
						{this.softCapsConfigWithLimits?.map(({ unitStat, breakpoints, capType, postCapEPs }, index) => (
							<>
								<tr>
									<th className="text-nowrap" colSpan={2}>
										{unitStat.getShortName(this.playerClass)}
									</th>
									<td className="text-end">{statCapTypeNames.get(capType)}</td>
								</tr>
								{this.additionalSoftCapTooltipInformation[unitStat.getRootStat()] && (
									<>
										<tr>
											<td colSpan={3}>{this.additionalSoftCapTooltipInformation[unitStat.getRootStat()]?.()}</td>
										</tr>
										<tr>
											<td colSpan={3} className="pb-2"></td>
										</tr>
									</>
								)}
								<tr>
									<th className="text-end">
										<em>%</em>
									</th>
									<th colSpan={2} className="text-nowrap text-end">
										<em>Post cap EP</em>
									</th>
								</tr>
								{breakpoints.map((breakpoint, breakpointIndex) => (
									<tr>
										<td className="text-end">{this.breakpointValueToDisplayPercentage(breakpoint, unitStat)}</td>
										<td colSpan={2} className="text-end">
											{unitStat
												.convertEpToRatingScale(capType === StatCapType.TypeThreshold ? postCapEPs[0] : postCapEPs[breakpointIndex])
												.toFixed(2)}
										</td>
									</tr>
								))}
								{index !== this.softCapsConfigWithLimits.length - 1 && (
									<>
										<tr>
											<td colSpan={3} className="border-bottom pb-2"></td>
										</tr>
										<tr>
											<td colSpan={3} className="pb-2"></td>
										</tr>
									</>
								)}
							</>
						))}
					</tbody>
				</table>
			</>
		);
	}

	setFreezeItemSlots(eventID: EventID, newValue: boolean) {
		if (this.freezeItemSlots !== newValue) {
			this.freezeItemSlots = newValue;
			this.frozenItemSlots.clear();
			this.freezeItemSlotsChangeEmitter.emit(eventID);
		}
	}

	buildContextMenu(button: HTMLButtonElement) {
		const instance = tippy(button, {
			interactive: true,
			trigger: 'click',
			theme: 'reforge-optimiser-popover',
			placement: 'right-start',
			onShow: instance => {
				const useCustomEPValuesInput = new BooleanPicker(null, this.player, {
					id: 'reforge-optimizer-enable-custom-ep-weights',
					label: 'Use custom EP Weights',
					inline: true,
					changedEvent: () => this.sim.useCustomEPValuesChangeEmitter,
					getValue: () => this.sim.getUseCustomEPValues(),
					setValue: (eventID, _player, newValue) => {
						this.sim.setUseCustomEPValues(eventID, newValue);
					},
				});
				let useSoftCapBreakpointsInput: BooleanPicker<Player<any>> | null = null;
				if (!!this.softCapsConfig?.length) {
					useSoftCapBreakpointsInput = new BooleanPicker(null, this.player, {
						id: 'reforge-optimizer-enable-soft-cap-breakpoints',
						label: 'Use soft cap breakpoints',
						inline: true,
						changedEvent: () => this.sim.useSoftCapBreakpointsChangeEmitter,
						getValue: () => this.sim.getUseSoftCapBreakpoints(),
						setValue: (eventID, _player, newValue) => {
							this.sim.setUseSoftCapBreakpoints(eventID, newValue);
						},
					});
				}

				const forcedProcInput = new EnumPicker(null, this.player, {
					id: 'reforge-optimizer-force-stat-proc',
					label: 'Force Matrix/Apparatus proc',
					values: [
						{ name: 'Any', value: -1 },
						...[...RelativeStatCap.relevantStats].map(stat => {
							return {
								name: UnitStat.fromStat(stat).getShortName(this.playerClass),
								value: stat,
							};
						}),
					],
					changedEvent: () => this.player.gearChangeEmitter,
					getValue: () => {
						if (!this.relativeStatCap) {
							return -1;
						} else {
							return this.relativeStatCap!.forcedHighestStat.getStat();
						}
					},
					setValue: (_eventID, _player, newValue) => {
						if (newValue == -1) {
							this.relativeStatCap = null;
						} else {
							this.relativeStatCap = new RelativeStatCap(newValue, this.playerClass);
						}
					},
					showWhen: () => {
						const canEnable = RelativeStatCap.canEnable(this.player);

						if (!canEnable) {
							this.relativeStatCap = null;
						}

						return canEnable;
					},
				});

				const freezeItemSlotsInput = new BooleanPicker(null, this.player, {
					id: 'reforge-optimizer-freeze-item-slots',
					label: 'Freeze item slots',
					labelTooltip:
						'Flag one or more item slots to be "frozen", which will prevent the optimizer from changing the Reforge in that slot from its current setting. This can be useful for hybrid classes who use the same gear piece for multiple raid roles.',
					inline: true,
					changedEvent: () => this.freezeItemSlotsChangeEmitter,
					getValue: () => this.freezeItemSlots,
					setValue: (eventID, _player, newValue) => {
						this.setFreezeItemSlots(eventID, newValue);
					},
				});

				const descriptionRef = ref<HTMLParagraphElement>();
				instance.setContent(
					<>
						{useCustomEPValuesInput.rootElem}
						<div ref={descriptionRef} className={clsx('mb-0', this.sim.getUseCustomEPValues() && 'hide')}>
							<p>This will enable modification of the default EP weights and setting custom stat caps.</p>
							<p>Ep weights can be modified in the Stat Weights editor.</p>
							<p>If you want to hard cap a stat make sure to put the EP for that stat higher.</p>
						</div>
						{this.buildCapsList({
							useCustomEPValuesInput: useCustomEPValuesInput,
							description: descriptionRef.value!,
						})}
						{useSoftCapBreakpointsInput?.rootElem}
						{forcedProcInput.rootElem}
						{this.buildSoftCapBreakpointsLimiter({ useSoftCapBreakpointsInput })}
						{freezeItemSlotsInput.rootElem}
						{this.buildFrozenSlotsInputs()}
						{this.buildEPWeightsToggle({ useCustomEPValuesInput: useCustomEPValuesInput })}
					</>,
				);
			},
			onHidden: () => {
				instance.setContent(<></>);
			},
		});
	}

	buildFrozenSlotsInputs() {
		const allSlots = this.player.getGear().getItemSlots();
		const numRows = Math.floor(allSlots.length / 2) + 1;
		const slotsByRow: ItemSlot[][] = [];

		for (let rowIdx = 0; rowIdx < numRows; rowIdx++) {
			slotsByRow.push(allSlots.slice(rowIdx * 2, (rowIdx + 1) * 2));
		}

		const tableRef = ref<HTMLTableElement>();
		const content = (
			<table ref={tableRef}>
				{slotsByRow.map(slots => {
					const rowRef = ref<HTMLTableRowElement>();
					const row = (
						<tr ref={rowRef}>
							{slots.map(slot => {
								const picker = new BooleanPicker(null, this.player, {
									id: 'reforge-optimizer-freeze-' + ItemSlot[slot],
									label: slotNames.get(slot),
									inline: true,
									changedEvent: () => this.freezeItemSlotsChangeEmitter,
									getValue: () => this.frozenItemSlots.get(slot) || false,
									setValue: (_eventID, _player, newValue) => {
										this.frozenItemSlots.set(slot, newValue);
									},
									showWhen: () => this.freezeItemSlots,
								});
								const column = <td>{picker.rootElem}</td>;
								return column;
							})}
						</tr>
					);
					return row;
				})}
			</table>
		);

		return content;
	}

	buildCapsList({ useCustomEPValuesInput, description }: { useCustomEPValuesInput: BooleanPicker<Player<any>>; description: HTMLElement }) {
		const sharedInputConfig: Pick<NumberPickerConfig<Player<any>>, 'changedEvent'> = {
			changedEvent: _ => TypedEvent.onAny([this.sim.useSoftCapBreakpointsChangeEmitter, this.player.statCapsChangeEmitter]),
		};

		const tableRef = ref<HTMLTableElement>();
		const statCapTooltipRef = ref<HTMLButtonElement>();
		const defaultStatCapsButtonRef = ref<HTMLButtonElement>();

		const content = (
			<table ref={tableRef} className={clsx('reforge-optimizer-stat-cap-table mb-2', !this.sim.getUseCustomEPValues() && 'hide')}>
				<thead>
					<tr>
						<th colSpan={3} className="pb-3">
							<div className="d-flex">
								<h6 className="content-block-title mb-0 me-1">Edit stat caps</h6>
								<button ref={statCapTooltipRef} className="d-inline">
									<i className="fa-regular fa-circle-question" />
								</button>
								<button ref={defaultStatCapsButtonRef} className="d-inline ms-auto" onclick={() => this.setDefaultStatCaps()}>
									<i className="fas fa-arrow-rotate-left" />
								</button>
							</div>
						</th>
					</tr>
					<tr>
						<th>Stat</th>
						<th colSpan={2} className="text-end">
							%
						</th>
					</tr>
				</thead>
				<tbody>
					{this.simUI.individualConfig.displayStats.map(unitStat => {
						if (!unitStat.hasRootStat()) return;
						const rootStat = unitStat.getRootStat();
						if (!INCLUDED_STATS.includes(rootStat)) return;

						const listElementRef = ref<HTMLTableRowElement>();
						const statName = unitStat.getShortName(this.player.getClass());

						const sharedStatInputConfig: Pick<NumberPickerConfig<Player<any>>, 'getValue' | 'setValue'> = {
							getValue: () => {
								return this.toVisualUnitStatPercentage(this.statCaps.getUnitStat(unitStat), unitStat);
							},
							setValue: (_eventID, _player, newValue) => {
								this.setStatCap(unitStat, this.toDefaultUnitStatValue(newValue, unitStat));
							},
						};

						const percentagePicker = new NumberPicker(null, this.player, {
							id: `reforge-optimizer-${statName}-percentage`,
							float: true,
							maxDecimalDigits: 5,
							showZeroes: false,
							positive: true,
							extraCssClasses: ['mb-0'],
							enableWhen: () => this.isAllowedToOverrideStatCaps || !this.softCapsConfig.some(config => config.unitStat.equals(unitStat)),
							...sharedInputConfig,
							...sharedStatInputConfig,
						});
						const statPresets = this.statSelectionPresets?.find(entry => entry.unitStat.equals(unitStat))?.presets;
						const presets = !!statPresets
							? new EnumPicker(null, this.player, {
									id: `reforge-optimizer-${statName}-presets`,
									extraCssClasses: ['mb-0'],
									label: '',
									values: [
										{ name: 'Select preset', value: 0 },
										...[...statPresets.keys()].map(key => {
											const percentValue = statPresets.get(key)!;

											return {
												name: `${key} - ${percentValue.toFixed(2)}%`,
												value: percentValue,
											};
										}),
									].sort((a, b) => a.value - b.value),
									enableWhen: () => this.isAllowedToOverrideStatCaps || !this.softCapsConfig.some(config => config.unitStat.equals(unitStat)),
									...sharedInputConfig,
									...sharedStatInputConfig,
							  })
							: null;

						const tooltipText = this.statTooltips[rootStat];
						const statTooltipRef = ref<HTMLButtonElement>();

						const row = (
							<>
								<tr ref={listElementRef} className="reforge-optimizer-stat-cap-item">
									<td>
										<div className="reforge-optimizer-stat-cap-item-label">
											{statName}{' '}
											{tooltipText && (
												<button ref={statTooltipRef} className="d-inline">
													<i className="fa-regular fa-circle-question" />
												</button>
											)}
										</div>
									</td>
									<td colSpan={2}>{percentagePicker.rootElem}</td>
								</tr>
								{presets && (
									<tr>
										<td></td>
										<td colSpan={2}>{presets.rootElem}</td>
									</tr>
								)}
							</>
						);

						const tooltip = tooltipText
							? tippy(statTooltipRef.value!, {
									content: tooltipText,
							  })
							: null;

						useCustomEPValuesInput.addOnDisposeCallback(() => tooltip?.destroy());

						return row;
					})}
				</tbody>
			</table>
		);

		if (statCapTooltipRef.value) {
			const tooltip = tippy(statCapTooltipRef.value, {
				content:
					'Stat caps are the maximum amount of a stat that can be gained from Reforging. If a stat exceeds its cap, the optimizer will attempt to reduce it to the cap value.',
			});
			useCustomEPValuesInput.addOnDisposeCallback(() => tooltip.destroy());
		}
		if (defaultStatCapsButtonRef.value) {
			const tooltip = tippy(defaultStatCapsButtonRef.value, {
				content: 'Reset to stat cap defaults',
			});
			useCustomEPValuesInput.addOnDisposeCallback(() => tooltip.destroy());
		}

		const event = this.sim.useCustomEPValuesChangeEmitter.on(() => {
			const isUsingCustomEPValues = this.sim.getUseCustomEPValues();
			tableRef.value?.classList[isUsingCustomEPValues ? 'remove' : 'add']('hide');
			description?.classList[!isUsingCustomEPValues ? 'remove' : 'add']('hide');
		});

		useCustomEPValuesInput.addOnDisposeCallback(() => {
			content.remove();
			event.dispose();
		});

		return content;
	}

	buildEPWeightsToggle({ useCustomEPValuesInput }: { useCustomEPValuesInput: BooleanPicker<Player<any>> }) {
		const extraCssClasses = ['mt-3'];
		if (!this.sim.getUseCustomEPValues()) extraCssClasses.push('hide');
		const savedEpWeights = renderSavedEPWeights(null, this.simUI, { extraCssClasses, loadOnly: true });
		const event = this.sim.useCustomEPValuesChangeEmitter.on(() => {
			const isUsingCustomEPValues = this.sim.getUseCustomEPValues();
			savedEpWeights.rootElem?.classList[isUsingCustomEPValues ? 'remove' : 'add']('hide');
		});

		useCustomEPValuesInput.addOnDisposeCallback(() => {
			savedEpWeights.dispose();
			savedEpWeights.rootElem.remove();
			event.dispose();
		});

		return (
			<>
				{savedEpWeights.rootElem}
				{this.simUI.epWeightsModal && (
					<button
						className="btn btn-outline-primary"
						onclick={() => {
							this.simUI.epWeightsModal?.open();
							hideAll();
						}}>
						Edit weights
					</button>
				)}
			</>
		);
	}

	buildSoftCapBreakpointsLimiter({ useSoftCapBreakpointsInput }: { useSoftCapBreakpointsInput: BooleanPicker<Player<any>> | null }) {
		if (!this.enableBreakpointLimits || !useSoftCapBreakpointsInput) return null;

		const tableRef = ref<HTMLTableElement>();
		const breakpointsLimitTooltipRef = ref<HTMLButtonElement>();

		const content = (
			<table ref={tableRef} className={clsx('reforge-optimizer-stat-cap-table mb-2', !this.sim.getUseSoftCapBreakpoints() && 'hide')}>
				<thead>
					<tr>
						<th colSpan={3} className="pb-3">
							<div className="d-flex">
								<h6 className="content-block-title mb-0 me-1">Breakpoint limit</h6>
								<button ref={breakpointsLimitTooltipRef} className="d-inline">
									<i className="fa-regular fa-circle-question" />
								</button>
							</div>
						</th>
					</tr>
				</thead>
				<tbody>
					{this.softCapsConfig
						.filter(config => config.capType === StatCapType.TypeThreshold && config.breakpoints.length > 1)
						.map(({ breakpoints, unitStat }) => {
							if (!unitStat.hasRootStat()) return;
							const rootStat = unitStat.getRootStat();
							if (!INCLUDED_STATS.includes(rootStat)) return;

							const listElementRef = ref<HTMLTableRowElement>();
							const statName = unitStat.getShortName(this.player.getClass());
							const picker = !!breakpoints
								? new EnumPicker(null, this.player, {
										id: `reforge-optimizer-${statName}-presets`,
										extraCssClasses: ['mb-0'],
										label: '',
										values: [
											{ name: 'No limit set', value: 0 },
											...breakpoints.map(breakpoint => ({
												name: `${this.breakpointValueToDisplayPercentage(breakpoint, unitStat)}%`,
												value: breakpoint,
											})),
										].sort((a, b) => a.value - b.value),
										changedEvent: _ => TypedEvent.onAny([this.sim.useSoftCapBreakpointsChangeEmitter]),
										getValue: () => {
											return this.player.getBreakpointLimits().getUnitStat(unitStat) || 0;
										},
										setValue: (eventID, _player, newValue) => {
											this.player.setBreakpointLimits(eventID, this.player.getBreakpointLimits().withUnitStat(unitStat, newValue));
										},
								  })
								: null;

							if (!picker?.rootElem) return null;

							const row = (
								<>
									<tr ref={listElementRef} className="reforge-optimizer-stat-cap-item">
										<td>
											<div className="reforge-optimizer-stat-cap-item-label">{statName}</div>
										</td>
										<td colSpan={2}>{picker.rootElem}</td>
									</tr>
								</>
							);

							return row;
						})}
				</tbody>
			</table>
		);

		if (breakpointsLimitTooltipRef.value) {
			const tooltip = tippy(breakpointsLimitTooltipRef.value, {
				content: 'Allows you to set a custom breakpoint limit.',
			});
			useSoftCapBreakpointsInput.addOnDisposeCallback(() => tooltip.destroy());
		}

		const event = this.sim.useSoftCapBreakpointsChangeEmitter.on(() => {
			const isUsingBreakpoints = this.sim.getUseSoftCapBreakpoints();
			tableRef.value?.classList[isUsingBreakpoints ? 'remove' : 'add']('hide');
		});

		useSoftCapBreakpointsInput.addOnDisposeCallback(() => {
			content.remove();
			event?.dispose();
		});

		return content;
	}

	get isAllowedToOverrideStatCaps() {
		return !(this.sim.getUseSoftCapBreakpoints() && this.softCapsConfig);
	}

	get processedStatCaps() {
		let statCaps = this.statCaps;
		if (!this.isAllowedToOverrideStatCaps)
			this.softCapsConfigWithLimits.forEach(({ unitStat }) => {
				statCaps = statCaps.withUnitStat(unitStat, 0);
			});

		return statCaps;
	}

	async optimizeReforges() {
		if (isDevMode()) console.log('Starting Reforge optimization...');

		// First, clear all existing Reforges
		if (isDevMode()) {
			console.log('Clearing existing Reforges...');
			console.log('The following slots will not be cleared:');
			console.log(Array.from(this.frozenItemSlots.keys()).filter(key => this.frozenItemSlots.get(key)));
		}
		this.previousGear = this.player.getGear();
		this.previousReforges = this.previousGear.getAllReforges();
		const baseGear = this.previousGear.withoutReforges(this.player.canDualWield2H(), this.frozenItemSlots);
		const baseStats = await this.updateGear(baseGear);

		// Compute effective stat caps for just the Reforge contribution
		const reforgeCaps = baseStats.computeStatCapsDelta(this.processedStatCaps);
		if (isDevMode()) {
			console.log('Stat caps for Reforge contribution:');
			console.log(reforgeCaps);
		}

		// Do the same for any soft cap breakpoints that were configured
		const reforgeSoftCaps = this.computeReforgeSoftCaps(baseStats);

		// Perform any required processing on the pre-cap EPs to make them internally consistent with the
		// configured hard caps and soft caps.
		const validatedWeights = ReforgeOptimizer.checkWeights(this.preCapEPs, reforgeCaps, reforgeSoftCaps);

		// Set up YALPS model
		const variables = this.buildYalpsVariables(baseGear, validatedWeights);
		const constraints = this.buildYalpsConstraints(baseGear, baseStats);

		// Solve in multiple passes to enforce caps
		await this.solveModel(baseGear, validatedWeights, reforgeCaps, reforgeSoftCaps, variables, constraints, 75000);
		this.currentReforges = this.player.getGear().getAllReforges();
	}

	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		let baseStats = Stats.fromProto(this.player.getCurrentStats().finalStats);
		baseStats = baseStats.addStat(Stat.StatMasteryRating, this.player.getBaseMastery() * Mechanics.MASTERY_RATING_PER_MASTERY_POINT);
		if (this.updateGearStatsModifier) baseStats = this.updateGearStatsModifier(baseStats);
		return baseStats;
	}

	computeReforgeSoftCaps(baseStats: Stats): StatCap[] {
		const reforgeSoftCaps: StatCap[] = [];

		if (!this.isAllowedToOverrideStatCaps) {
			this.softCapsConfigWithLimits.slice().forEach(config => {
				let weights = config.postCapEPs.slice();
				const relativeBreakpoints = [];

				for (const breakpoint of config.breakpoints) {
					relativeBreakpoints.push(baseStats.computeGapToCap(config.unitStat, breakpoint));
				}

				// For stats that are configured as thresholds rather than soft caps,
				// reverse the order of evaluation of the breakpoints so that the
				// largest relevant threshold is always targeted. Likewise, use a
				// single value for the post-cap EP for these stats, which should be
				// interpreted (and computed) as the residual stat value just after
				// passing a threshold discontinuity.
				if (config.capType == StatCapType.TypeThreshold) {
					relativeBreakpoints.reverse();
					weights = Array(relativeBreakpoints.length).fill(weights[0]);
				}

				reforgeSoftCaps.push(new StatCap(config.unitStat, relativeBreakpoints, config.capType, weights));
			});
		}

		return reforgeSoftCaps;
	}

	buildYalpsVariables(gear: Gear, preCapEPs: Stats): YalpsVariables {
		const variables = new Map<string, YalpsCoefficients>();
		const epStats = this.simUI.individualConfig.epStats;

		for (const slot of gear.getItemSlots()) {
			const item = gear.getEquippedItem(slot);

			if (!item || this.frozenItemSlots.get(slot)) {
				continue;
			}

			for (const reforgeData of this.player.getAvailableReforgings(item.withDynamicStats())) {
				if (!epStats.includes(reforgeData.toStat)) {
					continue;
				}

				const variableKey = `${slot}_${reforgeData.id}`;
				const coefficients = new Map<string, number>();
				coefficients.set(ItemSlot[slot], 1);
				this.applyReforgeStat(coefficients, reforgeData.fromStat, reforgeData.fromAmount, preCapEPs);
				this.applyReforgeStat(coefficients, reforgeData.toStat, reforgeData.toAmount, preCapEPs);
				variables.set(variableKey, coefficients);
			}
		}

		return variables;
	}

	// Apply stat dependencies before setting optimization coefficients
	applyReforgeStat(coefficients: YalpsCoefficients, stat: Stat, amount: number, preCapEPs: Stats) {
		// Handle Spirit to Spell Hit conversion for hybrid casters separately from standard dependencies
		if ((stat == Stat.StatSpirit && this.isHybridCaster) || stat == Stat.StatExpertiseRating) {
			this.setPseudoStatCoefficient(coefficients, PseudoStat.PseudoStatSpellHitPercent, amount / Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT);
		}

		// If a highest Stat constraint is to be enforced, then update the
		// associated coefficient if applicable.
		this.relativeStatCap?.updateCoefficients(coefficients, stat, amount);

		// If the pre-cap EP for the root stat is non-zero, then apply
		// the root stat directly and don't look for any children.
		if (preCapEPs.getStat(stat) != 0) {
			coefficients.set(Stat[stat], amount);
			return;
		}

		// Loop over all dependent PseudoStats
		for (const childStat of UnitStat.getChildren(stat)) {
			// Only add a dependency if the child has an EP value associated with it
			if (preCapEPs.getPseudoStat(childStat) != 0) {
				this.setPseudoStatCoefficient(coefficients, childStat, UnitStat.fromPseudoStat(childStat).convertRatingToPercent(amount)!);
			}
		}
	}

	setPseudoStatCoefficient(coefficients: YalpsCoefficients, pseudoStat: PseudoStat, amount: number) {
		const currentValue = coefficients.get(PseudoStat[pseudoStat]) || 0;
		coefficients.set(PseudoStat[pseudoStat], currentValue + amount);
	}

	buildYalpsConstraints(gear: Gear, baseStats: Stats): YalpsConstraints {
		const constraints = new Map<string, Constraint>();

		for (const slot of gear.getItemSlots()) {
			constraints.set(ItemSlot[slot], lessEq(1));
		}

		if (this.relativeStatCap) {
			const statsWithoutBaseMastery = baseStats.addStat(
				Stat.StatMasteryRating,
				-this.player.getBaseMastery() * Mechanics.MASTERY_RATING_PER_MASTERY_POINT,
			);
			this.relativeStatCap.updateConstraints(constraints, gear, statsWithoutBaseMastery);
		}

		return constraints;
	}

	async solveModel(
		gear: Gear,
		weights: Stats,
		reforgeCaps: Stats,
		reforgeSoftCaps: StatCap[],
		variables: YalpsVariables,
		constraints: YalpsConstraints,
		maxIterations: number,
	): Promise<number> {
		// Calculate EP scores for each Reforge option
		if (isDevMode()) {
			console.log('Stat weights for this iteration:');
			console.log(weights);
		}
		const updatedVariables = this.updateReforgeScores(variables, weights);
		if (isDevMode()) {
			console.log('Optimization variables and constraints for this iteration:');
			console.log(updatedVariables);
			console.log(constraints);
		}

		// Set up and solve YALPS model
		const model: Model = {
			direction: 'maximize',
			objective: 'score',
			constraints: constraints,
			variables: updatedVariables,
			binaries: true,
		};
		const options: Options = {
			timeout: Infinity,
			maxIterations: maxIterations,
			tolerance: 0.01,
		};
		const solution = solve(model, options);

		if (isDevMode()) {
			console.log('LP solution for this iteration:');
			console.log(solution);
		}

		if (isNaN(solution.result)) {
			if (maxIterations > 500000) {
				throw solution;
			} else {
				if (isDevMode()) console.log('No feasible solution was found, doubling max iterations...');
				return await this.solveModel(gear, weights, reforgeCaps, reforgeSoftCaps, variables, constraints, maxIterations * 2);
			}
		}

		// Apply the current solution
		const updatedGear = await this.applyLPSolution(gear, solution);

		// Check if any unconstrained stats exceeded their specified cap.
		// If so, add these stats to the constraint list and re-run the solver.
		// If no unconstrained caps were exceeded, then we're done.
		const [anyCapsExceeded, updatedConstraints, updatedWeights] = this.checkCaps(
			solution,
			reforgeCaps,
			reforgeSoftCaps,
			updatedVariables,
			constraints,
			weights,
		);

		if (!anyCapsExceeded) {
			if (isDevMode()) console.log('Reforge optimization has finished!');
			return solution.result;
		} else {
			if (isDevMode()) console.log('One or more stat caps were exceeded, starting constrained iteration...');
			await sleep(100);
			return await this.solveModel(updatedGear, updatedWeights, reforgeCaps, reforgeSoftCaps, updatedVariables, updatedConstraints, maxIterations);
		}
	}

	updateReforgeScores(variables: YalpsVariables, weights: Stats): YalpsVariables {
		const updatedVariables = new Map<string, YalpsCoefficients>();

		for (const [variableKey, coefficients] of variables.entries()) {
			let score = 0;
			const updatedCoefficients = new Map<string, number>();

			for (const [coefficientKey, value] of coefficients.entries()) {
				updatedCoefficients.set(coefficientKey, value);

				// Determine whether the key corresponds to a stat change. If so, apply
				// current EP for that stat. It is assumed that the supplied weights have
				// already been updated to post-cap values for any stats that were
				// constrained to be capped in a previous iteration.
				if (coefficientKey.includes('PseudoStat')) {
					const statKey = (PseudoStat as any)[coefficientKey] as PseudoStat;
					score += weights.getPseudoStat(statKey) * value;
				} else if (coefficientKey.includes('Stat')) {
					const statKey = (Stat as any)[coefficientKey] as Stat;
					score += weights.getStat(statKey) * value;
				}
			}

			updatedCoefficients.set('score', score);
			updatedVariables.set(variableKey, updatedCoefficients);
		}

		return updatedVariables;
	}

	async applyLPSolution(gear: Gear, solution: Solution): Promise<Gear> {
		let updatedGear = gear.withoutReforges(this.player.canDualWield2H(), this.frozenItemSlots);

		for (const [variableKey, _coefficient] of solution.variables) {
			const splitKey = variableKey.split('_');
			const slot = parseInt(splitKey[0]) as ItemSlot;
			const reforgeId = parseInt(splitKey[1]);
			const equippedItem = gear.getEquippedItem(slot);

			if (equippedItem) {
				updatedGear = updatedGear.withEquippedItem(
					slot,
					equippedItem.withReforge(this.sim.db.getReforgeById(reforgeId)!),
					this.player.canDualWield2H(),
				);
			}
		}

		await this.updateGear(updatedGear);
		return updatedGear;
	}

	checkCaps(
		solution: Solution,
		reforgeCaps: Stats,
		reforgeSoftCaps: StatCap[],
		variables: YalpsVariables,
		constraints: YalpsConstraints,
		currentWeights: Stats,
	): [boolean, YalpsConstraints, Stats] {
		// First add up the total stat changes from the solution
		let reforgeStatContribution = new Stats();

		for (const [variableKey, _coefficient] of solution.variables) {
			for (const [coefficientKey, value] of variables.get(variableKey)!.entries()) {
				if (coefficientKey.includes('PseudoStat')) {
					const statKey = (PseudoStat as any)[coefficientKey] as PseudoStat;
					reforgeStatContribution = reforgeStatContribution.addPseudoStat(statKey, value);
				} else if (coefficientKey.includes('Stat')) {
					const statKey = (Stat as any)[coefficientKey] as Stat;
					reforgeStatContribution = reforgeStatContribution.addStat(statKey, value);
				}
			}
		}

		if (isDevMode()) {
			console.log('Total stat contribution from Reforging:');
			console.log(reforgeStatContribution);
		}

		// Then check whether any unconstrained stats exceed their cap
		let anyCapsExceeded = false;
		const updatedConstraints = new Map<string, Constraint>(constraints);
		let updatedWeights = currentWeights;

		for (const [unitStat, value] of reforgeStatContribution.asUnitStatArray()) {
			const cap = reforgeCaps.getUnitStat(unitStat);
			const statName = unitStat.getKey();

			if (cap !== 0 && value > cap && !constraints.has(statName)) {
				updatedConstraints.set(statName, greaterEq(cap));
				anyCapsExceeded = true;
				if (isDevMode()) console.log('Cap exceeded for: %s', statName);

				// Set EP to 0 for hard capped stats
				updatedWeights = updatedWeights.withUnitStat(unitStat, 0);
			}
		}

		// If hard caps are all taken care of, then deal with any remaining soft cap breakpoints
		while (!anyCapsExceeded && reforgeSoftCaps.length > 0) {
			const nextSoftCap = reforgeSoftCaps[0];
			const unitStat = nextSoftCap.unitStat;
			const statName = unitStat.getKey();
			const currentValue = reforgeStatContribution.getUnitStat(unitStat);

			let idx = 0;
			for (const breakpoint of nextSoftCap.breakpoints) {
				if (currentValue > breakpoint) {
					updatedConstraints.set(statName, greaterEq(breakpoint));
					updatedWeights = updatedWeights.withUnitStat(unitStat, nextSoftCap.postCapEPs[idx]);
					anyCapsExceeded = true;
					if (isDevMode()) console.log('Breakpoint exceeded for: %s', statName);
					break;
				}

				idx++;
			}

			// For true soft cap stats (evaluated in ascending order), remove any breakpoint that was
			// exceeded from the configuration. If no breakpoints were exceeded or there are none
			// remaining, then remove the entry completely from reforgeSoftCaps. In contrast, for threshold
			// stats (evaluated in descending order), always remove the entry completely after the first
			// pass.
			if (nextSoftCap.capType == StatCapType.TypeSoftCap) {
				nextSoftCap.breakpoints = nextSoftCap.breakpoints.slice(idx + 1);
				nextSoftCap.postCapEPs = nextSoftCap.postCapEPs.slice(idx + 1);
			}

			if (nextSoftCap.capType == StatCapType.TypeThreshold || nextSoftCap.breakpoints.length == 0) {
				reforgeSoftCaps.shift();
			}
		}

		return [anyCapsExceeded, updatedConstraints, updatedWeights];
	}

	private get baseMastery() {
		return this.player.getBaseMastery() * Mechanics.MASTERY_RATING_PER_MASTERY_POINT;
	}

	private toVisualTotalMasteryPercentage(statPoints: number, statValue: number) {
		// If the value is less than or equal to the base mastery, then set it to 0,
		// because we assume you want to reset this stat cap.
		if (statValue - this.baseMastery <= 0) {
			statPoints = 0;
		} else {
			// When displaying the mastery percentage we want to include the base mastery
			statPoints *= this.player.getMasteryPerPointModifier();
		}
		return statPoints;
	}

	private toVisualUnitStatPercentage(statValue: number, unitStat: UnitStat) {
		const rawStatValue = statValue;
		let percentOrPointsValue = unitStat.convertDefaultUnitsToPercent(rawStatValue)!;
		if (unitStat.equalsStat(Stat.StatMasteryRating)) percentOrPointsValue = this.toVisualTotalMasteryPercentage(percentOrPointsValue, rawStatValue);

		return percentOrPointsValue;
	}

	private toDefaultUnitStatValue(value: number, unitStat: UnitStat) {
		let statValue = unitStat.convertPercentToDefaultUnits(value)!;
		if (unitStat.equalsStat(Stat.StatMasteryRating)) statValue /= this.player.getMasteryPerPointModifier();
		return statValue;
	}

	private breakpointValueToDisplayPercentage(value: number, unitStat: UnitStat) {
		return unitStat.equalsStat(Stat.StatMasteryRating)
			? ((value / Mechanics.MASTERY_RATING_PER_MASTERY_POINT) * this.player.getMasteryPerPointModifier()).toFixed(2)
			: unitStat.convertDefaultUnitsToPercent(value)!.toFixed(2);
	}

	onReforgeDone() {
		const itemSlots = this.player.getGear().getItemSlots();
		const changedSlots = new Map<ItemSlot, ReforgeData | undefined>();
		for (const slot of itemSlots) {
			const prev = this.previousReforges.get(slot);
			const current = this.currentReforges.get(slot);
			if (!ReforgeStat.equals(prev?.reforge, current?.reforge)) changedSlots.set(slot, current);
		}
		const hasReforgeChanges = changedSlots.size;

		const copyButtonContainerRef = ref<HTMLDivElement>();
		const changedReforgeMessage = (
			<>
				<p className="mb-0">The following items were reforged:</p>
				<ul>
					{[...changedSlots].map(([slot, reforge]) => {
						if (reforge) {
							const slotName = slotNames.get(slot);
							const { fromStat, toStat } = reforge;
							const fromText = shortSecondaryStatNames.get(fromStat);
							const toText = shortSecondaryStatNames.get(toStat);
							return (
								<li>
									{slotName}: {fromText} → {toText}
								</li>
							);
						} else {
							return <li>{slotNames.get(slot)}: Removed reforge</li>;
						}
					})}
				</ul>
				<div ref={copyButtonContainerRef} />
			</>
		);

		if (hasReforgeChanges) {
			const settingsExport = IndividualSimSettings.toJson(this.simUI.toProto());
			if (settingsExport)
				new CopyButton(copyButtonContainerRef.value!, {
					extraCssClasses: ['btn-outline-primary'],
					getContent: () => JSON.stringify(settingsExport),
					text: 'Copy to Reforge Lite',
				});
		}

		new Toast({
			variant: 'success',
			body: hasReforgeChanges ? changedReforgeMessage : <>No reforge changes were made!</>,
			delay: hasReforgeChanges ? 5000 : 3000,
		});
	}

	onReforgeError(error: any) {
		if (isDevMode()) console.log(error);

		if (this.previousGear) this.updateGear(this.previousGear);
		new Toast({
			variant: 'error',
			body: 'Reforge optimization failed. Please try again, or report the issue if it persists.',
		});
	}
}
