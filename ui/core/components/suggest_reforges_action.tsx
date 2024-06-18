import clsx from 'clsx';
import tippy, { hideAll } from 'tippy.js';
import { ref } from 'tsx-vanilla';
import { Constraint, greaterEq, lessEq, Model, Options, Solution, solve } from 'yalps';

import * as Mechanics from '../constants/mechanics.js';
import { IndividualSimUI } from '../individual_sim_ui';
import { Player } from '../player';
import { ItemSlot, Spec, Stat } from '../proto/common';
import { StatCapConfig, StatCapType } from '../proto/ui';
import { Gear } from '../proto_utils/gear';
import { getClassStatName } from '../proto_utils/names';
import { statPercentageOrPointsToNumber, Stats, statToPercentageOrPoints } from '../proto_utils/stats';
import { SpecTalents } from '../proto_utils/utils';
import { Sim } from '../sim';
import { ActionGroupItem } from '../sim_ui';
import { TypedEvent } from '../typed_event';
import { isDevMode, sleep } from '../utils';
import { BooleanPicker } from './pickers/boolean_picker';
import { NumberPicker, NumberPickerConfig } from './pickers/number_picker';
import { renderSavedEPWeights } from './saved_data_managers/ep_weights';
import Toast from './toast';

type YalpsCoefficients = Map<string, number>;
type YalpsVariables = Map<string, YalpsCoefficients>;
type YalpsConstraints = Map<string, Constraint>;

const EXCLUDED_STATS = [
	Stat.StatStamina,
	Stat.StatHealth,
	Stat.StatStrength,
	Stat.StatAgility,
	Stat.StatAttackPower,
	Stat.StatRangedAttackPower,
	Stat.StatIntellect,
	Stat.StatSpellPower,
	Stat.StatSpellPenetration,
	Stat.StatSpirit,
	Stat.StatMana,
	Stat.StatMP5,
	Stat.StatBlock,
	Stat.StatBonusArmor,
	Stat.StatArcaneResistance,
	Stat.StatNatureResistance,
	Stat.StatFireResistance,
	Stat.StatFrostResistance,
	Stat.StatShadowResistance,
];

const STAT_TOOLTIP: { [key in Stat]?: () => Element | string } = {
	[Stat.StatMastery]: () => (
		<>
			Rating: Excludes your base mastery
			<br />
			%: Includes base mastery
		</>
	),
};

export type ReforgeOptimizerOptions = {
	// Allows you to modify the stats before they are returned for the calculations
	// For example: Adding class specific Glyphs/Talents that are not added by the backend
	updateGearStatsModifier?: (baseStats: Stats) => Stats;
};

export class ReforgeOptimizer {
	protected readonly simUI: IndividualSimUI<any>;
	protected readonly player: Player<any>;
	protected readonly isHybridCaster: boolean;
	protected readonly sim: Sim;
	protected readonly defaults: IndividualSimUI<any>['individualConfig']['defaults'];
	protected _statCaps: Stats;
	protected updateGearStatsModifier: ReforgeOptimizerOptions['updateGearStatsModifier'];
	protected softCapsConfig: StatCapConfig[];

	constructor(simUI: IndividualSimUI<any>, options?: ReforgeOptimizerOptions) {
		this.simUI = simUI;
		this.player = simUI.player;
		this.isHybridCaster = [Spec.SpecBalanceDruid, Spec.SpecShadowPriest, Spec.SpecElementalShaman].includes(this.player.getSpec());
		this.sim = simUI.sim;
		this.defaults = simUI.individualConfig.defaults;
		this.updateGearStatsModifier = options?.updateGearStatsModifier;
		this.softCapsConfig = this.defaults.softCapBreakpoints || [];
		this._statCaps = this.statCaps;

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
					new Toast({
						variant: 'success',
						body: 'Reforge optimization complete!',
					});
				} catch {
					new Toast({
						variant: 'error',
						body: 'Reforge optimization failed. Please try again, or report the issue if it persists.',
					});
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

		const [_startReforgeOptimizationButton, contextMenuButton] = simUI.addActionGroup([startReforgeOptimizationEntry, contextMenuEntry], {
			cssClass: 'suggest-reforges-settings-group d-flex',
		});

		if (!!this.softCapsConfig?.length)
			tippy(_startReforgeOptimizationButton, {
				content: this.buildReforgeButtonTooltip(),
				placement: 'bottom',
				maxWidth: 310,
			});

		tippy(contextMenuButton, {
			placement: 'bottom',
			content: 'Change Reforge Optimizer settings',
		});

		this.buildContextMenu(contextMenuButton);
	}

	get statCaps() {
		return this.sim.getUseCustomEPValues() ? this.player.getStatCaps() : this.defaults.statCaps || new Stats();
	}
	setStatCap(stat: Stat, value: number) {
		this._statCaps = this._statCaps.withStat(stat, value);
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
		let weights = this.sim.getUseCustomEPValues() ? this.player.getEpWeights() : this.defaults.epWeights;

		// Replace Spirit EP for hybrid casters with a small value in order to break ties between Spirit and Hit Reforges
		if (this.isHybridCaster) {
			weights = weights.withStat(Stat.StatSpirit, 0.01);
		}

		return weights;
	}

	buildReforgeButtonTooltip() {
		return (
			<>
				The following soft caps / breakpoints have been implemented for this spec:
				<ul className="mt-1 mb-0">{this.softCapsConfig?.map(({ stat }) => <li>{getClassStatName(stat, this.player.getClass())}</li>)}</ul>
			</>
		);
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
					label: 'Enable custom EP Weights',
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
						label: 'Enable soft cap breakpoints',
						inline: true,
						changedEvent: () => this.sim.useSoftCapBreakpointsChangeEmitter,
						getValue: () => this.sim.getUseSoftCapBreakpoints(),
						setValue: (eventID, _player, newValue) => {
							this.sim.setUseSoftCapBreakpoints(eventID, newValue);
						},
					});
				}

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
						{this.buildEPWeightsToggle({ useCustomEPValuesInput: useCustomEPValuesInput })}
					</>,
				);
			},
			onHidden: () => {
				instance.setContent(<></>);
			},
		});
	}

	buildCapsList({ useCustomEPValuesInput, description }: { useCustomEPValuesInput: BooleanPicker<Player<any>>; description: HTMLElement }) {
		const numberPickerSharedConfig: Pick<NumberPickerConfig<Player<any>>, 'float' | 'showZeroes' | 'positive' | 'extraCssClasses' | 'changedEvent'> = {
			float: true,
			showZeroes: false,
			positive: true,
			extraCssClasses: ['mb-0'],
			changedEvent: _ => this.player.statCapsChangeEmitter,
		};

		const tableRef = ref<HTMLTableElement>();
		const statCapTooltipRef = ref<HTMLButtonElement>();
		const defaultStatCapsButtonRef = ref<HTMLButtonElement>();

		const stats = new Stats(this.simUI.individualConfig.displayStats);

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
						<th className="text-end">Rating</th>
						<th className="text-end">%</th>
					</tr>
				</thead>
				<tbody>
					{this.simUI.individualConfig.displayStats.map(stat => {
						if (EXCLUDED_STATS.includes(stat)) return;

						const listElementRef = ref<HTMLTableRowElement>();
						const statName = getClassStatName(stat, this.player.getClass());
						const ratingPicker = new NumberPicker(null, this.player, {
							...numberPickerSharedConfig,
							id: `character-bonus-stat-${stat}-rating`,
							getValue: () => {
								let statValue = this.statCaps.getStat(stat);
								if (stat === Stat.StatMastery) statValue = this.toVisualBaseMasteryRating(statValue);

								return statValue;
							},
							setValue: (_eventID, _player, newValue) => {
								let statValue = newValue;
								if (stat === Stat.StatMastery) statValue += this.baseMastery;

								this.setStatCap(stat, statValue);
							},
						});
						const percentagePicker = new NumberPicker(null, this.player, {
							...numberPickerSharedConfig,
							id: `character-bonus-stat-${stat}-percentage`,
							getValue: () => {
								const statValue = this.statCaps.getStat(stat);
								let statInPercentageOrPoints = statToPercentageOrPoints(stat, statValue, stats);
								if (stat === Stat.StatMastery)
									statInPercentageOrPoints = this.toVisualTotalMasteryPercentage(statInPercentageOrPoints, statValue);

								return statInPercentageOrPoints;
							},
							setValue: (_eventID, _player, newValue) => {
								let statInNumber = statPercentageOrPointsToNumber(stat, newValue, stats);
								if (stat === Stat.StatMastery) statInNumber /= this.player.getMasteryPerPointModifier();

								this.setStatCap(stat, statInNumber);
							},
						});

						const tooltipText = STAT_TOOLTIP[stat];
						const statTooltipRef = ref<HTMLButtonElement>();

						const row = (
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
								<td>{ratingPicker.rootElem}</td>
								<td>{percentagePicker.rootElem}</td>
							</tr>
						);
						console.log(tooltipText);
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
						}}>Edit weights</button>
				)}
			</>
		);
	}

	async optimizeReforges() {
		if (isDevMode()) console.log('Starting Reforge optimization...');

		// First, clear all existing Reforges
		if (isDevMode()) console.log('Clearing existing Reforges...');
		const baseGear = this.player.getGear().withoutReforges(this.player.canDualWield2H());
		const baseStats = await this.updateGear(baseGear);

		// Compute effective stat caps for just the Reforge contribution
		const reforgeCaps = baseStats.computeStatCapsDelta(this.statCaps);
		if (isDevMode()) {
			console.log('Stat caps for Reforge contribution:');
			console.log(reforgeCaps);
		}

		// Do the same for any soft cap breakpoints that were configured
		const reforgeSoftCaps = this.computeReforgeSoftCaps(baseStats);

		// Set up YALPS model
		const variables = this.buildYalpsVariables(baseGear);
		const constraints = this.buildYalpsConstraints(baseGear);

		// Solve in multiple passes to enforce caps
		await this.solveModel(baseGear, this.preCapEPs, reforgeCaps, reforgeSoftCaps, variables, constraints);
	}

	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		let baseStats = Stats.fromProto(this.player.getCurrentStats().finalStats);
		baseStats = baseStats.addStat(Stat.StatMastery, this.player.getBaseMastery() * Mechanics.MASTERY_RATING_PER_MASTERY_POINT);
		if (this.updateGearStatsModifier) baseStats = this.updateGearStatsModifier(baseStats);
		return baseStats;
	}

	computeReforgeSoftCaps(baseStats: Stats): StatCapConfig[] {
		const reforgeSoftCaps: StatCapConfig[] = [];

		if (this.sim.getUseSoftCapBreakpoints() && this.softCapsConfig) {
			this.softCapsConfig
				.slice()
				.filter(config => this.statCaps.getStat(config.stat) == 0)
				.forEach(config => {
					let weights = config.postCapEPs.slice();
					const relativeBreakpoints = [];

					for (const breakpoint of config.breakpoints) {
						relativeBreakpoints.push(breakpoint - baseStats.getStat(config.stat));
					}

					// For stats that are configured as thresholds rather than
					// soft caps, reverse the order of evaluation of the
					// breakpoints so that the largest relevant threshold is
					// always targeted. Likewise, make sure that post-cap EPs
					// are always 0 for these stats.
					if (config.capType == StatCapType.TypeThreshold) {
						relativeBreakpoints.reverse();
						weights = Array(relativeBreakpoints.length).fill(0);
					}

					reforgeSoftCaps.push({
						stat: config.stat,
						breakpoints: relativeBreakpoints,
						capType: config.capType,
						postCapEPs: weights,
					});
				});
		}

		return reforgeSoftCaps;
	}

	buildYalpsVariables(gear: Gear): YalpsVariables {
		const variables = new Map<string, YalpsCoefficients>();

		for (const slot of gear.getItemSlots()) {
			const item = gear.getEquippedItem(slot);

			if (!item) {
				continue;
			}

			for (const reforgeData of this.player.getAvailableReforgings(item)) {
				const variableKey = `${slot}_${reforgeData.id}`;
				const coefficients = new Map<string, number>();
				coefficients.set(ItemSlot[slot], 1);

				for (const fromStat of reforgeData.fromStat) {
					this.applyReforgeStat(coefficients, fromStat, reforgeData.fromAmount);
				}

				for (const toStat of reforgeData.toStat) {
					this.applyReforgeStat(coefficients, toStat, reforgeData.toAmount);
				}

				variables.set(variableKey, coefficients);
			}
		}

		return variables;
	}

	applyReforgeStat(coefficients: YalpsCoefficients, stat: Stat, amount: number) {
		// Apply Spirit to Spell Hit conversion for hybrid casters before setting optimization coefficients
		let appliedStat = stat;
		let appliedAmount = amount;

		if (stat == Stat.StatSpirit && this.isHybridCaster) {
			appliedStat = Stat.StatSpellHit;

			switch (this.player.getSpec()) {
				case Spec.SpecBalanceDruid:
					appliedAmount *= 0.5 * (this.player.getTalents() as SpecTalents<Spec.SpecBalanceDruid>).balanceOfPower;
					break;
				case Spec.SpecShadowPriest:
					appliedAmount *= 0.5 * (this.player.getTalents() as SpecTalents<Spec.SpecShadowPriest>).twistedFaith;
					break;
				case Spec.SpecElementalShaman:
					appliedAmount *= [0, 0.33, 0.66, 1][(this.player.getTalents() as SpecTalents<Spec.SpecElementalShaman>).elementalPrecision];
					break;
			}

			// Also set the Spirit coefficient directly in order to break ties between Hit and Spirit Reforges
			coefficients.set(Stat[stat], amount);
		}

		const currentValue = coefficients.get(Stat[appliedStat]) || 0;
		coefficients.set(Stat[appliedStat], currentValue + appliedAmount);
	}

	buildYalpsConstraints(gear: Gear): YalpsConstraints {
		const constraints = new Map<string, Constraint>();

		for (const slot of gear.getItemSlots()) {
			constraints.set(ItemSlot[slot], lessEq(1));
		}

		return constraints;
	}

	async solveModel(
		gear: Gear,
		weights: Stats,
		reforgeCaps: Stats,
		reforgeSoftCaps: StatCapConfig[],
		variables: YalpsVariables,
		constraints: YalpsConstraints,
	) {
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
			timeout: 60000,
			maxIterations: Infinity,
			tolerance: 0.01,
		};
		const solution = solve(model, options);
		if (isDevMode()) {
			console.log('LP solution for this iteration:');
			console.log(solution);
		}
		// Apply the current solution
		await this.applyLPSolution(gear, solution);

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
		} else {
			if (isDevMode()) console.log('One or more stat caps were exceeded, starting constrained iteration...');
			await sleep(100);
			await this.solveModel(gear, updatedWeights, reforgeCaps, reforgeSoftCaps, updatedVariables, updatedConstraints);
		}
	}

	updateReforgeScores(variables: YalpsVariables, weights: Stats): YalpsVariables {
		const updatedVariables = new Map<string, YalpsCoefficients>();

		for (const [variableKey, coefficients] of variables.entries()) {
			let score = 0;
			const updatedCoefficients = new Map<string, number>();

			for (const [coefficientKey, value] of coefficients.entries()) {
				updatedCoefficients.set(coefficientKey, value);

				// Determine whether the key corresponds to a stat change. If
				// so, apply current EP for that stat. It is assumed that the
				// supplied weights have already been updated to post-cap
				// values for any stats that were constrained to be capped in
				// a previous iteration.
				if (coefficientKey.includes('Stat')) {
					const statKey = (Stat as any)[coefficientKey] as Stat;
					score += weights.getStat(statKey) * value;
				}
			}

			updatedCoefficients.set('score', score);
			updatedVariables.set(variableKey, updatedCoefficients);
		}

		return updatedVariables;
	}

	async applyLPSolution(gear: Gear, solution: Solution) {
		let updatedGear = gear.withoutReforges(this.player.canDualWield2H());

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
	}

	checkCaps(
		solution: Solution,
		reforgeCaps: Stats,
		reforgeSoftCaps: StatCapConfig[],
		variables: YalpsVariables,
		constraints: YalpsConstraints,
		currentWeights: Stats,
	): [boolean, YalpsConstraints, Stats] {
		// First add up the total stat changes from the solution
		let reforgeStatContribution = new Stats();

		for (const [variableKey, _coefficient] of solution.variables) {
			for (const [coefficientKey, value] of variables.get(variableKey)!.entries()) {
				if (coefficientKey.includes('Stat')) {
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

		for (const [statKey, value] of reforgeStatContribution.asArray().entries()) {
			const cap = reforgeCaps.getStat(statKey);
			const statName = Stat[statKey];

			if (cap !== 0 && value > cap && !constraints.has(statName)) {
				updatedConstraints.set(statName, greaterEq(cap));
				anyCapsExceeded = true;
				if (isDevMode()) console.log('Cap exceeded for: %s', statName);

				// Set EP to 0 for hard capped stats
				updatedWeights = updatedWeights.withStat(statKey, 0);
			}
		}

		// If hard caps are all taken care of, then deal with any remaining soft cap breakpoints
		while (!anyCapsExceeded && reforgeSoftCaps.length > 0) {
			const nextSoftCap = reforgeSoftCaps[0];
			const statKey = nextSoftCap.stat;
			const statName = Stat[statKey];
			const currentValue = reforgeStatContribution.getStat(statKey);

			let idx = 0;
			for (const breakpoint of nextSoftCap.breakpoints) {
				if (currentValue > breakpoint) {
					updatedConstraints.set(statName, greaterEq(breakpoint));
					updatedWeights = updatedWeights.withStat(statKey, nextSoftCap.postCapEPs[idx]);
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

	private toVisualBaseMasteryRating(value: number) {
		// If the value is less than or equal to the base mastery, then set it to 0,
		// because we assume you want to reset this stat cap.
		if (value <= this.baseMastery) {
			value = 0;
		} else {
			// Visually we show the mastery rating without the base mastery included
			value -= this.baseMastery;
		}
		return value;
	}
}
