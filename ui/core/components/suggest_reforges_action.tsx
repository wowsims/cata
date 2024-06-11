import clsx from 'clsx';
import tippy, { Instance as TippyInstance } from 'tippy.js';
import { ref } from 'tsx-vanilla';
import { Constraint, greaterEq, lessEq, Model, Solution, solve } from 'yalps';

import * as Mechanics from '../constants/mechanics.js';
import { IndividualSimUI } from '../individual_sim_ui';
import { Player } from '../player';
import { ItemSlot, PseudoStat, Stat } from '../proto/common';
import { Gear } from '../proto_utils/gear';
import { getClassStatName } from '../proto_utils/names';
import { Stats } from '../proto_utils/stats';
import { Sim } from '../sim';
import { ActionGroupItem } from '../sim_ui';
import { TypedEvent } from '../typed_event';
import { isDevMode, noop, sleep } from '../utils';
import { BooleanPicker } from './boolean_picker';
import { NumberPicker } from './number_picker';
import Toast from './toast';

interface StatWeightsDefaults {
	statCaps: Stats;
	preCapEPs: Stats;
}

type YalpsCoefficients = Map<string, number>;
type YalpsVariables = Map<string, YalpsCoefficients>;
type YalpsConstraints = Map<string, Constraint>;

export class ReforgeOptimizer {
	protected readonly simUI: IndividualSimUI<any>;
	protected readonly player: Player<any>;
	protected readonly sim: Sim;
	protected readonly defaults: StatWeightsDefaults;
	protected _statCaps: Stats;

	constructor(simUI: IndividualSimUI<any>, defaults: StatWeightsDefaults) {
		this.simUI = simUI;
		this.player = simUI.player;
		this.sim = simUI.sim;
		this.defaults = defaults;
		this._statCaps = this.defaults.statCaps;

		const startReforgeOptimizationEntry: ActionGroupItem = {
			label: 'Suggest Reforges',
			cssClass: 'flex-grow-1 suggest-reforges-action',
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
			cssClass: '',
			children: (
				<>
					<i className="fas fa-cog" />
				</>
			),
		};

		const [_startReforgeOptimizationButton, contextMenuButton] = simUI.addActionGroup([startReforgeOptimizationEntry, contextMenuEntry], {
			cssClass: 'd-flex',
		});

		this.buildContextMenu(contextMenuButton);
	}

	get statCaps(): Stats {
		return this.sim.getUseCustomEPValues() ? this._statCaps : this.defaults.statCaps;
	}
	setStatCap(stat: Stat, value: number): Stats {
		this._statCaps = this._statCaps.withStat(stat, value);
		return this.statCaps;
	}

	get preCapEPs(): Stats {
		return this.sim.getUseCustomEPValues() ? this.player.getEpWeights() : this.defaults.preCapEPs;
	}

	buildContextMenu(button: HTMLButtonElement) {
		tippy(button, {
			content: 'Change Reforge Optimizer settings',
		});

		tippy(button!, {
			interactive: true,
			trigger: 'click',
			theme: 'reforge-optimiser-popover',
			placement: 'right-start',
			onShow: instance => {
				const content = (<></>) as HTMLElement;

				const useCustomEPValuesInput = new BooleanPicker(content, this.player, {
					id: 'reforge-optimizer-enable-custom-ep-weights',
					label: 'Enable custom EP Weights',
					inline: true,
					changedEvent: player => player.epWeightsChangeEmitter,
					getValue: () => this.sim.getUseCustomEPValues(),
					setValue: (eventID, _player, newValue) => {
						console.log('Setting useCustomEPValues to', newValue);
						this.sim.setUseCustomEPValues(eventID, newValue);
					},
				});

				const descriptionRef = ref<HTMLParagraphElement>();
				content.appendChild(
					<p ref={descriptionRef} className={clsx('fst-italic mb-0', this.sim.getUseCustomEPValues() && 'hide')}>
						This will enable modication to the default EP weights and set custom stat caps. Ep weights can be modified in the Stat Weights editor.
					</p>,
				);

				content.appendChild(this.buildCapsList({ instance, input: useCustomEPValuesInput, description: descriptionRef.value! }));

				instance.setContent(content);
			},
		});
	}

	buildCapsList({ instance, input, description }: { instance: TippyInstance; input: BooleanPicker<Player<any>>; description: HTMLElement }) {
		const tableRef = ref<HTMLUListElement>();
		const statCapTooltipRef = ref<HTMLButtonElement>();

		const event = this.sim.useCustomEPValuesChangeEmitter.on(() => {
			const isUsingCustomEPValues = this.sim.getUseCustomEPValues();
			tableRef.value?.classList[isUsingCustomEPValues ? 'remove' : 'add']('hide');
			description?.classList[!isUsingCustomEPValues ? 'remove' : 'add']('hide');
		});

		input.addOnDisposeCallback(() => event.dispose());

		const statCaps = this.statCaps;
		const stats = new Stats(this.simUI.individualConfig.displayStats);

		const content = (
			<ul ref={tableRef} className={clsx('reforge-optimizer-stat-cap-list list-reset d-grid gap-2', !this.sim.getUseCustomEPValues() && 'hide')}>
				<li>
					<label className="me-1">Edit stat caps</label>
					<button ref={statCapTooltipRef} className="d-inline">
						<i className="fa-regular fa-circle-question" />
					</button>
				</li>
				{this.simUI.individualConfig.displayStats.map(stat => {
					if (stat === Stat.StatHealth) return;

					const listElementRef = ref<HTMLLIElement>();
					const statName = getClassStatName(stat, this.player.getClass());

					const picker = new NumberPicker(null, this.player, {
						id: `character-bonus-stat-${stat}`,
						inline: true,
						label: `${statName} ${stat === Stat.StatMastery ? 'Points' : '%'}`,
						extraCssClasses: ['mb-0'],
						changedEvent: _ => this.sim.useCustomEPValuesChangeEmitter,
						getValue: () => {
							const currentStat = statCaps.getStat(stat);
							let statInPercentage = 0;
							switch (stat) {
								case Stat.StatMeleeHit:
									statInPercentage = currentStat / Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE;
									break;
								case Stat.StatSpellHit:
									statInPercentage = currentStat / Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE;
									break;
								case Stat.StatMeleeCrit:
								case Stat.StatSpellCrit:
									statInPercentage = currentStat / Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE;
									break;
								case Stat.StatMeleeHaste:
									statInPercentage = currentStat / Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
									break;
								case Stat.StatSpellHaste:
									statInPercentage = currentStat / Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
									break;
								case Stat.StatExpertise:
									statInPercentage = currentStat / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION / 4;
									break;
								case Stat.StatBlock:
									statInPercentage = currentStat / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE + 5.0;
									break;
								case Stat.StatDodge:
									statInPercentage = stats.getPseudoStat(PseudoStat.PseudoStatDodge) / 100;
									break;
								case Stat.StatParry:
									statInPercentage = stats.getPseudoStat(PseudoStat.PseudoStatParry) / 100;
									break;
								case Stat.StatResilience:
									statInPercentage = currentStat / Mechanics.RESILIENCE_RATING_PER_CRIT_REDUCTION_CHANCE;
									break;
								case Stat.StatMastery:
									statInPercentage = currentStat / Mechanics.MASTERY_RATING_PER_MASTERY_POINT;
									break;
							}
							return statInPercentage;
						},
						setValue: (_eventID, _player, newValue) => {
							let statInPoints = 0;
							switch (stat) {
								case Stat.StatMeleeHit:
									statInPoints = newValue * Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE;
									break;
								case Stat.StatSpellHit:
									statInPoints = newValue * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE;
									break;
								case Stat.StatMeleeCrit:
								case Stat.StatSpellCrit:
									statInPoints = newValue * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE;
									break;
								case Stat.StatMeleeHaste:
									statInPoints = newValue * Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
									break;
								case Stat.StatSpellHaste:
									statInPoints = newValue * Mechanics.HASTE_RATING_PER_HASTE_PERCENT;
									break;
								case Stat.StatExpertise:
									statInPoints = newValue * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION * 4;
									break;
								case Stat.StatBlock:
									statInPoints = newValue * Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE - 5.0;
									break;
								case Stat.StatDodge:
									statInPoints = stats.getPseudoStat(PseudoStat.PseudoStatDodge) * 100;
									break;
								case Stat.StatParry:
									statInPoints = stats.getPseudoStat(PseudoStat.PseudoStatParry) * 100;
									break;
								case Stat.StatResilience:
									statInPoints = newValue * Mechanics.RESILIENCE_RATING_PER_CRIT_REDUCTION_CHANCE;
									break;
								case Stat.StatMastery:
									statInPoints = newValue * Mechanics.MASTERY_RATING_PER_MASTERY_POINT;
									break;
							}
							this.setStatCap(stat, statInPoints);
						},
					});

					return (
						<li ref={listElementRef} className="reforge-optimizer-stat-cap-item">
							{picker.rootElem}
						</li>
					);
				})}
			</ul>
		);

		if (statCapTooltipRef.value) {
			tippy(statCapTooltipRef.value, {
				content:
					'Stat caps are the maximum amount of a stat that can be gained from Reforging. If a stat exceeds its cap, the optimizer will attempt to reduce it to the cap value.',
			});
		}

		return content;
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
		// Set up YALPS model
		const variables = this.buildYalpsVariables(baseGear);
		const constraints = this.buildYalpsConstraints(baseGear);

		// Solve in multiple passes to enforce caps
		await this.solveModel(baseGear, reforgeCaps, variables, constraints);
	}

	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		return Stats.fromProto(this.player.getCurrentStats().finalStats);
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
					coefficients.set(Stat[fromStat], reforgeData.fromAmount);
				}

				for (const toStat of reforgeData.toStat) {
					coefficients.set(Stat[toStat], reforgeData.toAmount);
				}

				variables.set(variableKey, coefficients);
			}
		}

		return variables;
	}

	buildYalpsConstraints(gear: Gear): YalpsConstraints {
		const constraints = new Map<string, Constraint>();

		for (const slot of gear.getItemSlots()) {
			constraints.set(ItemSlot[slot], lessEq(1));
		}

		return constraints;
	}

	async solveModel(gear: Gear, reforgeCaps: Stats, variables: YalpsVariables, constraints: YalpsConstraints) {
		// Calculate EP scores for each Reforge option
		const updatedVariables = this.updateReforgeScores(variables, constraints);
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
		const solution = solve(model);
		if (isDevMode()) {
			console.log('LP solution for this iteration:');
			console.log(solution);
		}
		// Apply the current solution
		await this.applyLPSolution(gear, solution);

		// Check if any unconstrained stats exceeded their specified cap.
		// If so, add these stats to the constraint list and re-run the solver.
		// If no unconstrained caps were exceeded, then we're done.
		const [anyCapsExceeded, updatedConstraints] = this.checkCaps(solution, reforgeCaps, updatedVariables, constraints);

		if (!anyCapsExceeded) {
			if (isDevMode()) console.log('Reforge optimization has finished!');
		} else {
			if (isDevMode()) console.log('One or more stat caps were exceeded, starting constrained iteration...');
			await sleep(100);
			await this.solveModel(gear, reforgeCaps, updatedVariables, updatedConstraints);
		}
	}

	updateReforgeScores(variables: YalpsVariables, constraints: YalpsConstraints): YalpsVariables {
		const updatedVariables = new Map<string, YalpsCoefficients>();

		for (const [variableKey, coefficients] of variables.entries()) {
			let score = 0;
			const updatedCoefficients = new Map<string, number>();

			for (const [coefficientKey, value] of coefficients.entries()) {
				updatedCoefficients.set(coefficientKey, value);

				// Determine whether the key corresponds to a stat change.
				// If so, check whether the stat has already been constrained to be capped in a previous iteration.
				// Apply stored EP only for unconstrained stats.
				if (coefficientKey.includes('Stat') && !constraints.has(coefficientKey)) {
					const statKey = (Stat as any)[coefficientKey] as Stat;
					score += this.preCapEPs.getStat(statKey) * value;
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

	checkCaps(solution: Solution, reforgeCaps: Stats, variables: YalpsVariables, constraints: YalpsConstraints): [boolean, YalpsConstraints] {
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

		for (const [statKey, value] of reforgeStatContribution.asArray().entries()) {
			const cap = reforgeCaps.getStat(statKey);
			const statName = Stat[statKey];

			if (cap !== 0 && value > cap && !constraints.has(statName)) {
				updatedConstraints.set(statName, greaterEq(cap));
				anyCapsExceeded = true;
				if (isDevMode()) console.log('Cap exceeded for: %s', statName);
			}
		}

		return [anyCapsExceeded, updatedConstraints];
	}
}
