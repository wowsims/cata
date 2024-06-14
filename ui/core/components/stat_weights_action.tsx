import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../individual_sim_ui.jsx';
import { Player } from '../player.js';
import { ProgressMetrics, StatWeightsResult, StatWeightValues } from '../proto/api.js';
import { GemColor, ItemSlot, Profession, PseudoStat, Stat, UnitStats } from '../proto/common.js';
import { UIGem as Gem } from '../proto/ui.js';
import { Gear } from '../proto_utils/gear.js';
import * as Gems from '../proto_utils/gems.js';
import { getClassStatName } from '../proto_utils/names.js';
import { Stats, UnitStat } from '../proto_utils/stats.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { cloneChildren, combinationsWithDups, permutations, stDevToConf90, sum } from '../utils.js';
import { BaseModal } from './base_modal.jsx';
import { BooleanPicker } from './pickers/boolean_picker.js';
import { NumberPicker } from './pickers/number_picker.js';
import { ResultsViewer } from './results_viewer.jsx';

export function addStatWeightsAction(simUI: IndividualSimUI<any>, epStats: Stat[], epPseudoStats: PseudoStat[] | undefined, epReferenceStat: Stat) {
	simUI.addAction('Stat Weights', 'ep-weights-action', () => {
		// TODO: Make this so we can initialize the menu once outside of this function
		new EpWeightsMenu(simUI, epStats, epPseudoStats || [], epReferenceStat).open();
	});
}

// Create the config for modal in separate function, as constructor cannot
// contain any logic before `super' call. Use modal-xl to accommodate the extra
// TMI & p(death) EP in the UI.
function getModalConfig(simUI: IndividualSimUI<any>) {
	const baseConfig = { footer: true, scrollContents: true };
	if (simUI.sim.getShowThreatMetrics() && simUI.sim.getShowExperimental()) {
		return { size: 'xl' as const, ...baseConfig };
	}
	return baseConfig;
}

function scaledEpValue(stat: UnitStat, epRatios: number[], result: StatWeightsResult | null): number {
	if (!result) return 0;

	return (
		(result.dps?.epValues ? epRatios[0] * stat.getProtoValue(result.dps.epValues) : 0) +
		(result.hps?.epValues ? epRatios[1] * stat.getProtoValue(result.hps.epValues) : 0) +
		(result.tps?.epValues ? epRatios[2] * stat.getProtoValue(result.tps.epValues) : 0) +
		(result.dtps?.epValues ? epRatios[3] * stat.getProtoValue(result.dtps.epValues) : 0) +
		(result.tmi?.epValues ? epRatios[4] * stat.getProtoValue(result.tmi.epValues) : 0) +
		(result.pDeath?.epValues ? epRatios[5] * stat.getProtoValue(result.pDeath.epValues) : 0)
	);
}

class EpWeightsMenu extends BaseModal {
	private readonly simUI: IndividualSimUI<any>;
	private readonly container: HTMLElement;
	private readonly table: HTMLElement;
	private readonly tableBody: HTMLElement;
	private readonly resultsViewer: ResultsViewer;

	private statsType: string;
	private epStats: Stat[];
	private epPseudoStats: PseudoStat[];
	private epReferenceStat: Stat;
	private showAllStats = false;

	constructor(simUI: IndividualSimUI<any>, epStats: Stat[], epPseudoStats: PseudoStat[], epReferenceStat: Stat) {
		super(simUI.rootElem, 'ep-weights-menu', { ...getModalConfig(simUI), disposeOnClose: false });
		this.header?.insertAdjacentElement('afterbegin', <h5 className="modal-title">Calculate Stat Weights</h5>);

		this.simUI = simUI;
		this.statsType = 'ep';
		this.epStats = epStats;
		this.epPseudoStats = epPseudoStats;
		this.epReferenceStat = epReferenceStat;

		const statsTable = this.buildStatsTable();
		const containerRef = ref<HTMLDivElement>();
		const tableRef = ref<HTMLTableElement>();
		const tableBodyRef = ref<HTMLTableSectionElement>();
		const damageMetricsSelectRef = ref<HTMLSelectElement>();
		const healingMetricsSelectRef = ref<HTMLSelectElement>();
		const threatMetricsSelectRef = ref<HTMLSelectElement>();
		const typeSelectRef = ref<HTMLSelectElement>();
		const computeEpRef = ref<HTMLButtonElement>();
		const optimizeGemsButtonRef = ref<HTMLButtonElement>();
		const calcWeightsButtonRef = ref<HTMLButtonElement>();
		const allStatsContainerRef = ref<HTMLDivElement>();

		const getNameFromStat = (stat: Stat | undefined) => (stat !== undefined ? getClassStatName(stat, this.simUI.player.getClass()) : '??');
		const getStatFromName = (value: string) => Object.values(this.epStats).find(stat => getNameFromStat(stat) === value);
		const epRefSelectOptions = (
			<>
				{this.epStats.map(stat => (
					<option>{getNameFromStat(stat)}</option>
				))}
			</>
		);

		this.body.appendChild(
			<>
				<div className="ep-weights-options row">
					<div className="col col-sm-3">
						<select ref={typeSelectRef} className="ep-type-select form-select">
							<option value="ep">EP</option>
							<option value="weight">Weights</option>
						</select>
					</div>
					<div ref={allStatsContainerRef} className="show-all-stats-container col col-sm-3"></div>
				</div>
				<div className="ep-reference-options row experimental">
					<div className="col col-sm-4 damage-metrics">
						<span>DPS/TPS reference:</span>
						<select ref={damageMetricsSelectRef} className="ref-stat-select form-select damage-metrics">
							{epRefSelectOptions.cloneNode(true)}
						</select>
					</div>
					<div className="col col-sm-4 healing-metrics">
						<span>Healing reference:</span>
						<select ref={healingMetricsSelectRef} className="ref-stat-select form-select healing-metrics">
							{epRefSelectOptions.cloneNode(true)}
						</select>
					</div>
					<div className="col col-sm-4 threat-metrics">
						<span>Mitigation reference:</span>
						<select ref={threatMetricsSelectRef} className="ref-stat-select form-select threat-metrics">
							{epRefSelectOptions.cloneNode(true)}
						</select>
					</div>
					<p>The above stat selectors control which reference stat is used for EP normalisation for the different EP columns.</p>
				</div>
				<p>
					The 'Current EP' column displays the values currently used by the item pickers to sort items.
					<br />
					Use the <button className="fa fa-copy" /> icon above the EPs to use newly calculated EPs.
				</p>
				<div ref={containerRef} className="results-ep-table-container modal-scroll-table">
					<table ref={tableRef} className="results-ep-table">
						<thead>
							<tr>
								<th>Stat</th>
								{statsTable.map(({ metric, type, label, metricRef }) => {
									const isAction = type === 'action';
									return (
										<th className={clsx(metric && `${metric}-metrics`, isAction ? 'text-center' : `type-${type}`)}>
											<span>{label}</span>
											<button ref={metricRef} className="col-action">
												<i className={clsx('fas', isAction ? 'fa-arrows-rotate' : 'fa-copy')} />
											</button>
										</th>
									);
								})}
							</tr>
							<tr className="ep-ratios">
								<td>EP Ratio</td>
								{statsTable
									.filter(({ type }) => type !== 'action')
									.map(({ metric, type, experimental, ratioRef }) => (
										<td
											ref={ratioRef}
											className={clsx('type-ratio', `${metric}-metrics`, `type-${type}`, experimental && 'experimental')}
										/>
									))}
								<td className="text-center align-middle">
									<button ref={computeEpRef} className="btn btn-primary compute-ep">
										<i className="fas fa-calculator" />
										<span className="not-tiny">Update </span>EP
									</button>
								</td>
							</tr>
						</thead>
						<tbody ref={tableBodyRef}></tbody>
					</table>
				</div>
			</>,
		);

		this.footer!.appendChild(
			<>
				<button ref={optimizeGemsButtonRef} className="btn btn-primary optimize-gems experimental me-2">
					Optimize Gems
				</button>
				<button ref={calcWeightsButtonRef} className="btn btn-primary calc-weights">
					<i className="fas fa-calculator me-1" />
					Calculate
				</button>
			</>,
		);

		this.container = containerRef.value!;
		this.table = tableRef.value!;
		this.tableBody = tableBodyRef.value!;

		const pendingDiv = (<div className="results-pending-overlay" />) as HTMLDivElement;
		this.resultsViewer = new ResultsViewer(pendingDiv);

		const updateType = () => {
			if (this.statsType === 'ep') {
				this.table.classList.remove('stats-type-weight');
				this.table.classList.add('stats-type-ep');
			} else {
				this.table.classList.add('stats-type-weight');
				this.table.classList.remove('stats-type-ep');
			}
		};

		const selectElem = typeSelectRef.value!;
		selectElem.addEventListener('input', () => {
			this.statsType = selectElem.value;
			updateType();
		});
		selectElem.value = this.statsType;
		updateType();

		const updateEpRefStat = () => {
			this.simUI.player.epRefStatChangeEmitter.emit(TypedEvent.nextEventID());
			this.simUI.prevEpSimResult = this.calculateEp(this.getPrevSimResult());
			this.updateTable();
		};

		const damageMetricsSelect = damageMetricsSelectRef.value;
		if (damageMetricsSelect) {
			damageMetricsSelect.addEventListener('input', () => {
				this.simUI.dpsRefStat = getStatFromName(damageMetricsSelect.value);
				updateEpRefStat();
			});
			damageMetricsSelect.value = getNameFromStat(this.getDpsEpRefStat());
		}

		const healingMetricsSelect = healingMetricsSelectRef.value;
		if (healingMetricsSelect) {
			healingMetricsSelect.addEventListener('input', () => {
				this.simUI.healRefStat = getStatFromName(healingMetricsSelect.value);
				updateEpRefStat();
			});
			healingMetricsSelect.value = getNameFromStat(this.getHealEpRefStat());
		}
		const threatMetricsSelect = threatMetricsSelectRef.value;
		if (threatMetricsSelect) {
			threatMetricsSelect.addEventListener('input', () => {
				this.simUI.tankRefStat = getStatFromName(threatMetricsSelect.value);
				updateEpRefStat();
			});
			threatMetricsSelect.value = getNameFromStat(this.getTankEpRefStat());
		}

		const optimizeGemsButton = optimizeGemsButtonRef.value;
		if (optimizeGemsButton) {
			tippy(optimizeGemsButton, {
				content: (
					<>
						<p>
							<span className="warning link-warning">
								<i className="fa fa-exclamation-triangle"></i> WARNING
							</span>{' '}
							This feature is experimental, and will not always produce the most optimal gems especially when interacting with soft/hard stat
							caps.
						</p>
						<p>
							Optimizes equipped gems to maximize EP, based on the values in <b>Current EP</b>.
						</p>
						<p className="mb-0">
							Does not change the meta gem, but ensures that its condition is met. Uses JC gems if Jewelcrafting is a selected profession.
						</p>
					</>
				),
			});
			optimizeGemsButton.addEventListener('click', async () => {
				const defaultState = cloneChildren(optimizeGemsButton);
				optimizeGemsButton.classList.add('disabled');
				optimizeGemsButton.style.width = `${optimizeGemsButton.getBoundingClientRect().width.toFixed(3)}px`;
				optimizeGemsButton.replaceChildren(
					<>
						<i className=" fa fa-spinner fa-spin me-1" />
						Running
					</>,
				);
				await this.optimizeGems(TypedEvent.nextEventID());
				optimizeGemsButton.replaceChildren(...defaultState);
				optimizeGemsButton.classList.remove('disabled');
			});
		}

		const calcButton = calcWeightsButtonRef.value;
		calcButton?.addEventListener('click', async () => {
			this.simUI.rootElem.classList.add('blurred');
			this.simUI.rootElem.insertAdjacentElement('afterend', pendingDiv);

			this.container.scrollTo({ top: 0 });
			this.container.classList.add('pending');
			this.resultsViewer.setPending();
			const iterations = this.simUI.sim.getIterations();
			const result = await this.simUI.player.computeStatWeights(
				TypedEvent.nextEventID(),
				this.epStats,
				this.epPseudoStats,
				this.epReferenceStat,
				progress => {
					this.setSimProgress(progress);
				},
			);
			this.simUI.rootElem.classList.remove('blurred');
			pendingDiv.remove();
			this.container.classList.remove('pending');
			this.resultsViewer.hideAll();
			if (!result) return;
			this.simUI.prevEpIterations = iterations;
			this.simUI.prevEpSimResult = this.calculateEp(result);
			this.updateTable();
		});

		const makeUpdateWeights = (
			button: HTMLButtonElement,
			labelTooltip: string,
			tooltip: string,
			weightsFunc: () => UnitStats | undefined,
			epRefStat?: () => Stat,
		) => {
			const label = button.previousElementSibling as HTMLElement;
			const title = () => {
				if (!epRefStat) return labelTooltip;

				const refStatName = getNameFromStat(epRefStat());
				return `${labelTooltip} Normalized by ${refStatName}.`;
			};

			tippy(label, {
				content: title,
			});
			tippy(button, {
				content: tooltip,
			});

			button.addEventListener('click', () => {
				this.simUI.player.setEpWeights(TypedEvent.nextEventID(), Stats.fromProto(weightsFunc()));
				this.updateTable();
			});
		};
		statsTable.forEach(({ metricRef, labelTooltip, actionTooltip, getWeights, getEpRefStat }) =>
			makeUpdateWeights(metricRef!.value!, labelTooltip, actionTooltip, getWeights, getEpRefStat),
		);

		new BooleanPicker(allStatsContainerRef.value!, this, {
			id: 'ep-show-all-stats',
			label: 'Show All Stats',
			inline: true,
			changedEvent: () => new TypedEvent(),
			getValue: () => this.showAllStats,
			setValue: (_eventID: EventID, _menu: EpWeightsMenu, newValue: boolean) => {
				this.showAllStats = newValue;
				this.updateTable();
			},
		});

		this.updateTable();

		const makeEpRatioCell = (cell: HTMLElement, idx: number) => {
			new NumberPicker(cell, this.simUI.player, {
				id: `ep-ratio-${idx}`,
				float: true,
				changedEvent: player => player.epRatiosChangeEmitter,
				getValue: () => this.simUI.player.getEpRatios()[idx],
				setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
					const epRatios = player.getEpRatios();
					epRatios[idx] = newValue;
					player.setEpRatios(eventID, epRatios);
				},
			});
		};
		const epRatioCells = statsTable.filter(({ type, ratioRef }) => type === 'ep' && !!ratioRef?.value).map(({ ratioRef }) => ratioRef!.value!);
		epRatioCells.forEach(makeEpRatioCell);
		this.simUI.player.epRatiosChangeEmitter.on(_eventID => this.updateTable());

		const weightRatioCells = statsTable.filter(({ type, ratioRef }) => type === 'weight' && !!ratioRef?.value).map(({ ratioRef }) => ratioRef!.value!);
		weightRatioCells.forEach(makeEpRatioCell);

		const updateButton = computeEpRef.value!;
		tippy(updateButton, {
			content: 'Compute Weighted EP',
		});

		updateButton.addEventListener('click', () => {
			const results = this.getPrevSimResult();
			const epRatios = this.simUI.player.getEpRatios();
			if (this.statsType === 'ep') {
				const scaledDpsEp = Stats.fromProto(results.dps!.epValues).scale(epRatios[0]);
				const scaledHpsEp = Stats.fromProto(results.hps!.epValues).scale(epRatios[1]);
				const scaledTpsEp = Stats.fromProto(results.tps!.epValues).scale(epRatios[2]);
				const scaledDtpsEp = Stats.fromProto(results.dtps!.epValues).scale(epRatios[3]);
				const scaledTmiEp = Stats.fromProto(results.tmi!.epValues).scale(epRatios[4]);
				const scaledPDeathEp = Stats.fromProto(results.pDeath!.epValues).scale(epRatios[5]);
				const newEp = scaledDpsEp.add(scaledHpsEp).add(scaledTpsEp).add(scaledDtpsEp).add(scaledTmiEp).add(scaledPDeathEp);
				this.simUI.player.setEpWeights(TypedEvent.nextEventID(), newEp);
			} else {
				const scaledDpsWeights = Stats.fromProto(results.dps!.weights).scale(epRatios[0]);
				const scaledHpsWeights = Stats.fromProto(results.hps!.weights).scale(epRatios[1]);
				const scaledTpsWeights = Stats.fromProto(results.tps!.weights).scale(epRatios[2]);
				const scaledDtpsWeights = Stats.fromProto(results.dtps!.weights).scale(epRatios[3]);
				const scaledTmiWeights = Stats.fromProto(results.tmi!.weights).scale(epRatios[4]);
				const scaledPDeathWeights = Stats.fromProto(results.pDeath!.weights).scale(epRatios[5]);
				const newWeights = scaledDpsWeights
					.add(scaledHpsWeights)
					.add(scaledTpsWeights)
					.add(scaledDtpsWeights)
					.add(scaledTmiWeights)
					.add(scaledPDeathWeights);
				this.simUI.player.setEpWeights(TypedEvent.nextEventID(), newWeights);
			}
			this.updateTable();
		});
	}

	private setSimProgress(progress: ProgressMetrics) {
		this.resultsViewer.setContent(
			<div className="results-sim">
				<div>
					{progress.completedSims} / {progress.totalSims}
					<br />
					simulations complete
				</div>
				<div>
					{progress.completedIterations} / {progress.totalIterations}
					<br />
					iterations complete
				</div>
			</div>,
		);
	}

	private updateTable() {
		const tempTable = <></>;
		EpWeightsMenu.epUnitStats.forEach(stat => {
			// Don't show extra stats when 'Show all stats' is not selected
			if (
				(!this.showAllStats && stat.isStat() && !this.epStats.includes(stat.getStat())) ||
				(stat.isPseudoStat() && !this.epPseudoStats.includes(stat.getPseudoStat()))
			) {
				return;
			}
			const row = this.makeTableRow(stat);
			tempTable.appendChild(row);
		});
		this.tableBody.replaceChildren(tempTable);
	}

	private makeTableRow(stat: UnitStat): HTMLElement {
		const result = this.simUI.prevEpSimResult;
		const epRatios = this.simUI.player.getEpRatios();
		const rowTotalEp = scaledEpValue(stat, epRatios, result);
		const currentEpRef = ref<HTMLTableCellElement>();
		const row = (
			<tr>
				<td>{stat.getName(this.simUI.player.getClass())}</td>
				{this.makeTableRowCells(stat, result?.dps, 'damage-metrics', rowTotalEp, epRatios[0])}
				{this.makeTableRowCells(stat, result?.hps, 'healing-metrics', rowTotalEp, epRatios[1])}
				{this.makeTableRowCells(stat, result?.tps, 'threat-metrics', rowTotalEp, epRatios[2])}
				{this.makeTableRowCells(stat, result?.dtps, 'threat-metrics', rowTotalEp, epRatios[3])}
				{this.makeTableRowCells(stat, result?.tmi, 'threat-metrics experimental', rowTotalEp, epRatios[4])}
				{this.makeTableRowCells(stat, result?.pDeath, 'threat-metrics experimental', rowTotalEp, epRatios[5])}
				<td ref={currentEpRef} className="current-ep"></td>
			</tr>
		) as HTMLElement;

		const currentEpCell = currentEpRef.value!;
		new NumberPicker(currentEpCell, this.simUI.player, {
			id: `ep-weight-stat-${stat}`,
			float: true,
			changedEvent: (player: Player<any>) => player.epWeightsChangeEmitter,
			getValue: () => this.simUI.player.getEpWeights().getUnitStat(stat),
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const epWeights = player.getEpWeights().withUnitStat(stat, newValue);
				player.setEpWeights(eventID, epWeights);
			},
		});

		return row;
	}

	private makeTableRowCells(stat: UnitStat, statWeights: StatWeightValues | undefined, className: string, epTotal: number, epRatio: number) {
		let weightCell: Element | null = null;
		let epCell: Element | null = null;
		const isZeroEpRatio = epRatio === 0;
		const weightRef = ref<HTMLTableCellElement>();
		const epRef = ref<HTMLTableCellElement>();

		if (statWeights) {
			const weightAvg = stat.getProtoValue(statWeights.weights!);
			const weightStdev = stat.getProtoValue(statWeights.weightsStdev!);
			weightCell = this.makeTableCellContents(weightAvg, weightStdev);

			const epAvg = stat.getProtoValue(statWeights.epValues!);
			const epStdev = stat.getProtoValue(statWeights.epValuesStdev!);
			epCell = this.makeTableCellContents(epAvg, epStdev);
		} else {
			weightCell = <span className="results-avg notapplicable">N/A</span>;
			epCell = weightCell.cloneNode(true) as Element;
		}

		const row = (
			<>
				<td ref={weightRef} className={clsx('stdev-cell', 'type-weight', statWeights && isZeroEpRatio && 'unused-ep', className)}>
					{weightCell}
				</td>
				<td ref={epRef} className={clsx('stdev-cell', 'type-ep', statWeights && isZeroEpRatio && 'unused-ep', className)}>
					{epCell}
				</td>
			</>
		);

		if (!statWeights || isZeroEpRatio) return row;

		const epCurrent = this.simUI.player.getEpWeights().getUnitStat(stat);
		const epDelta = epTotal - epCurrent;

		const epAvgElem = epRef.value!.querySelector('.type-ep .results-avg')!;
		if (epDelta.toFixed(2) === '0.00') epAvgElem; // no-op
		else if (epDelta > 0) epAvgElem.classList.add('positive');
		else if (epDelta < 0) epAvgElem.classList.add('negative');

		return row;
	}

	private makeTableCellContents(value: number, stdev: number) {
		const iterations = this.simUI.prevEpIterations || 1;
		return (
			<>
				<span className="results-avg">{value.toFixed(2)}</span>
				<span className="results-stdev">
					(<i className="fas fa-plus-minus fa-xs"></i>
					{stDevToConf90(stdev, iterations).toFixed(2)})
				</span>
			</>
		) as HTMLElement;
	}

	private calculateEp(weights: StatWeightsResult) {
		const result = StatWeightsResult.clone(weights);

		if (this.simUI.dpsRefStat !== undefined) {
			EpWeightsMenu.normaliseEpValue(this.simUI.dpsRefStat, result.dps!);
			EpWeightsMenu.normaliseEpValue(this.simUI.dpsRefStat, result.tps!);
		}
		if (this.simUI.healRefStat !== undefined) {
			EpWeightsMenu.normaliseEpValue(this.simUI.healRefStat, result.hps!);
		}
		if (this.simUI.tankRefStat !== undefined) {
			EpWeightsMenu.normaliseEpValue(this.simUI.tankRefStat, result.dtps!);
			EpWeightsMenu.normaliseEpValue(this.simUI.tankRefStat, result.tmi!);
			EpWeightsMenu.normaliseEpValue(this.simUI.tankRefStat, result.pDeath!);
		}
		return result;
	}

	private static normaliseEpValue(refStat: Stat, values: StatWeightValues) {
		const refUnitStat = UnitStat.fromStat(refStat);
		const refWeight = refUnitStat.getProtoValue(values.weights!);
		const refStdev = refUnitStat.getProtoValue(values.weightsStdev!);
		EpWeightsMenu.epUnitStats.forEach(stat => {
			const value = stat.getProtoValue(values.weights!);
			stat.setProtoValue(values.epValues!, refWeight === 0 ? 0 : value / refWeight);

			const valueStdev = stat.getProtoValue(values.weightsStdev!);
			stat.setProtoValue(values.epValuesStdev!, refStdev === 0 ? 0 : valueStdev / refStdev);
		});
	}

	private getDpsEpRefStat(): Stat {
		return this.simUI.dpsRefStat !== undefined ? this.simUI.dpsRefStat : this.epReferenceStat;
	}

	private getHealEpRefStat(): Stat {
		return this.simUI.healRefStat !== undefined ? this.simUI.healRefStat : this.epReferenceStat;
	}

	private getTankEpRefStat(): Stat {
		return this.simUI.tankRefStat !== undefined ? this.simUI.tankRefStat : Stat.StatArmor;
	}

	private getPrevSimResult(): StatWeightsResult {
		return (
			this.simUI.prevEpSimResult ||
			StatWeightsResult.create({
				dps: {
					weights: new Stats().toProto(),
					weightsStdev: new Stats().toProto(),
					epValues: new Stats().toProto(),
					epValuesStdev: new Stats().toProto(),
				},
				hps: {
					weights: new Stats().toProto(),
					weightsStdev: new Stats().toProto(),
					epValues: new Stats().toProto(),
					epValuesStdev: new Stats().toProto(),
				},
				tps: {
					weights: new Stats().toProto(),
					weightsStdev: new Stats().toProto(),
					epValues: new Stats().toProto(),
					epValuesStdev: new Stats().toProto(),
				},
				dtps: {
					weights: new Stats().toProto(),
					weightsStdev: new Stats().toProto(),
					epValues: new Stats().toProto(),
					epValuesStdev: new Stats().toProto(),
				},
				tmi: {
					weights: new Stats().toProto(),
					weightsStdev: new Stats().toProto(),
					epValues: new Stats().toProto(),
					epValuesStdev: new Stats().toProto(),
				},
				pDeath: {
					weights: new Stats().toProto(),
					weightsStdev: new Stats().toProto(),
					epValues: new Stats().toProto(),
					epValuesStdev: new Stats().toProto(),
				},
			})
		);
	}

	private async optimizeGems(eventID: EventID) {
		// Replace 0 weights with a very tiny value, so we always prefer to take free stats even if the user gave a 0 weight.
		let epWeights = this.simUI.player.getEpWeights();
		epWeights = new Stats(epWeights.asArray().map(w => (w === 0 ? 1e-8 : w)));

		const gear = this.simUI.player.getGear();
		const allGems = this.simUI.sim.db.getGems();
		const phase = this.simUI.sim.getPhase();
		const isBlacksmithing = this.simUI.player.isBlacksmithing();
		const isJewelcrafting = this.simUI.player.hasProfession(Profession.Jewelcrafting);

		const optimizedGear = EpWeightsMenu.optimizeGemsForWeights(epWeights, gear, allGems, phase, isBlacksmithing, isJewelcrafting);
		this.simUI.player.setGear(eventID, optimizedGear);
	}

	private static optimizeGemsForWeights(
		epWeights: Stats,
		gear: Gear,
		allGems: Gem[],
		phase: number,
		isBlacksmithing: boolean,
		isJewelcrafting: boolean,
	): Gear {
		const unrestrictedGems = allGems.filter(gem => Gems.isUnrestrictedGem(gem, phase));

		const {
			bestGemForColor: bestGemForColor,
			bestGemForColorEP: _bestGemForColorEP,
			bestGemForSocket: bestGemForSocket,
			bestGemForSocketEP: bestGemForSocketEP,
			bestGem: bestGem,
			bestGemEP: bestGemEP,
		} = EpWeightsMenu.findBestGems(unrestrictedGems, epWeights);

		const items = gear.asMap();
		const socketBonusEPs = Object.values(items).map(item => (!!item ? new Stats(item.item.socketBonus).computeEP(epWeights) : 0));

		// Start by optimally filling all items, ignoring meta condition.
		Object.entries(items).forEach(([itemSlot, equippedItem], i) => {
			if (!equippedItem) return;

			//const item = equippedItem.item;
			const socketColors = equippedItem.curSocketColors(isBlacksmithing);

			// Compare whether its better to match sockets + get socket bonus, or just use best gems.
			const bestGemEPNotMatchingSockets = sum(socketColors.map(socketColor => (socketColor === GemColor.GemColorMeta ? 0 : bestGemEP)));
			const bestGemEPMatchingSockets =
				socketBonusEPs[i] + sum(socketColors.map(socketColor => (socketColor === GemColor.GemColorMeta ? 0 : bestGemForSocketEP[socketColor])));

			if (bestGemEPNotMatchingSockets > bestGemEPMatchingSockets) {
				socketColors.forEach((socketColor, i) => {
					if (socketColor !== GemColor.GemColorMeta) {
						equippedItem = equippedItem!.withGem(bestGem, i);
					}
				});
			} else {
				socketColors.forEach((socketColor, i) => {
					if (socketColor !== GemColor.GemColorMeta) {
						equippedItem = equippedItem!.withGem(bestGemForSocket[socketColor], i);
					}
				});
			}

			items[Number(itemSlot) as ItemSlot] = equippedItem;
		});
		gear = new Gear(items);

		const allSockets: { itemSlot: ItemSlot; socketIdx: number }[] = Object.keys(items)
			.map(itemSlotStr => {
				const itemSlot: ItemSlot = parseInt(itemSlotStr);
				const item = items[itemSlot];
				if (!item) {
					return [];
				}

				const numSockets = item.numSockets(isBlacksmithing);
				return [...Array(numSockets).keys()]
					.filter(socketIdx => item.item.gemSockets[socketIdx] !== GemColor.GemColorMeta)
					.map(socketIdx => ({
						itemSlot,
						socketIdx,
					}));
			})
			.flat();
		const threeSocketCombos = permutations(allSockets, 3);
		const calculateGearGemsEP = (gear: Gear): number => gear.statsFromGems(isBlacksmithing).computeEP(epWeights);

		// Now make adjustments to satisfy meta condition.
		// Use a wrapper function so we can return for readability.
		gear = ((gear: Gear): Gear => {
			const metaGem = gear.getMetaGem();
			if (!metaGem) {
				return gear;
			}

			const condition = Gems.getMetaGemCondition(metaGem.id);
			// Only TBC gems use compare color conditions, so just ignore them.
			if (!condition || condition.isCompareColorCondition()) {
				return gear;
			}

			// If there are very few non-meta gem slots, just skip because it's annoying to deal with.
			if (gear.getAllGems(isBlacksmithing).length - 1 < 3) {
				return gear;
			}

			// In wrath, all meta gems use min colors condition (numRed >= r && numYellow >= y && numBlue >= b)
			// All conditions require 3 gems, e.g. 3 of a single color, 2 of one color and 1 of another, or 1 of each.
			// So the maximum number of gems that ever need to change is 3.

			const colorCombos = EpWeightsMenu.getColorCombosToSatisfyCondition(condition);

			let bestGear = gear;
			let bestGearEP = calculateGearGemsEP(gear);

			// Use brute-force to try every possibility.
			colorCombos.forEach(colorCombo => {
				threeSocketCombos.forEach(socketCombo => {
					const curItems = gear.asMap();
					for (let i = 0; i < colorCombo.length; i++) {
						const gemColor = colorCombo[i];
						const { itemSlot, socketIdx } = socketCombo[i];
						curItems[itemSlot] = curItems[itemSlot]!.withGem(bestGemForColor[gemColor], socketIdx);
					}
					const curGear = new Gear(curItems);
					if (curGear.hasActiveMetaGem(isBlacksmithing)) {
						const curGearEP = calculateGearGemsEP(curGear);
						if (curGearEP > bestGearEP) {
							bestGear = curGear;
							bestGearEP = curGearEP;
						}
					}
				});
			});

			return bestGear;
		})(gear);

		// Now insert 3 JC gems, if Jewelcrafting is selected.
		// Use a wrapper function so we can return for readability.
		gear = ((gear: Gear): Gear => {
			if (!isJewelcrafting) {
				return gear;
			}

			const jcGems = allGems.filter(gem => gem.requiredProfession === Profession.Jewelcrafting);

			const {
				bestGemForColor: bestJcGemForColor,
				bestGemForColorEP: _bestJcGemForColorEP,
				bestGemForSocket: _bestJcGemForSocket,
				bestGemForSocketEP: _bestJcGemForSocketEP,
				bestGem: _bestJcGem,
				bestGemEP: _bestJcGemEP,
			} = EpWeightsMenu.findBestGems(jcGems, epWeights);

			let bestGear = gear;
			let bestGearEP = calculateGearGemsEP(gear);

			threeSocketCombos.forEach(socketCombo => {
				const curItems = gear.asMap();
				for (let i = 0; i < socketCombo.length; i++) {
					const { itemSlot, socketIdx } = socketCombo[i];
					const ei = curItems[itemSlot]!;
					const gemColor = ei.gems[socketIdx]!.color;
					curItems[itemSlot] = ei.withGem(bestJcGemForColor[gemColor], socketIdx);
				}

				const curGear = new Gear(curItems);
				if (curGear.hasActiveMetaGem(isBlacksmithing)) {
					const curGearEP = calculateGearGemsEP(curGear);
					if (curGearEP > bestGearEP) {
						bestGear = curGear;
						bestGearEP = curGearEP;
					}
				}
			});

			return bestGear;
		})(gear);

		return gear;
	}

	// Returns every possible way we could satisfy the gem condition.
	private static getColorCombosToSatisfyCondition(condition: Gems.MetaGemCondition): GemColor[][] {
		if (condition.isOneOfEach()) {
			return [Gems.PRIMARY_COLORS, [GemColor.GemColorPrismatic]].concat(
				Gems.SECONDARY_COLORS.map((secondaryColor, i) => {
					const remainingColor = Gems.PRIMARY_COLORS[i];
					return Gems.socketToMatchingColors.get(remainingColor)!.map(matchingColor => [matchingColor, secondaryColor]);
				}).flat(),
			);
		} else if (condition.isTwoAndOne()) {
			const oneColor = Gems.PRIMARY_COLORS[[condition.minRed, condition.minYellow, condition.minBlue].indexOf(1)];
			const twoColor = Gems.PRIMARY_COLORS[[condition.minRed, condition.minYellow, condition.minBlue].indexOf(2)];
			const secondaryColor = Gems.SECONDARY_COLORS.find(
				color => Gems.gemColorMatchesSocket(color, oneColor) && Gems.gemColorMatchesSocket(color, twoColor),
			)!;

			return [
				// All the ways to get 1 point in both colors. These are partial combos,
				// which still need 1 more gem in the 2-color.
				[GemColor.GemColorPrismatic],
				[secondaryColor],
				[oneColor, twoColor],
			]
				.map(partialCombo => {
					return Gems.socketToMatchingColors.get(twoColor)!.map(matchingColor => partialCombo.concat([matchingColor]));
				})
				.flat();
		} else if (condition.isThreeOfAColor()) {
			const threeColor = Gems.PRIMARY_COLORS[[condition.minRed, condition.minYellow, condition.minBlue].indexOf(3)];
			const matchingColors = Gems.socketToMatchingColors.get(threeColor)!;
			return combinationsWithDups(matchingColors, 3);
		} else {
			return [];
		}
	}

	private static findBestGems(gemList: Gem[], epWeights: Stats): BestGemsResult {
		// Best gem when we need a gem of a specific color.
		const bestGemForColor: Gem[] = new Array(Gems.GEM_COLORS.length).fill(undefined);
		const bestGemForColorEP: number[] = new Array(Gems.GEM_COLORS.length).fill(0);
		// Best gem when we need to match a socket to activate a bonus.
		const bestGemForSocket: Gem[] = bestGemForColor.slice();
		const bestGemForSocketEP: number[] = bestGemForColorEP.slice();
		// The single best gem, when color doesn't matter.
		let bestGem = gemList[0];
		let bestGemEP = 0;
		gemList.forEach(gem => {
			const gemEP = new Stats(gem.stats).computeEP(epWeights);
			if (gemEP > bestGemForColorEP[gem.color]) {
				bestGemForColorEP[gem.color] = gemEP;
				bestGemForColor[gem.color] = gem;

				if (gem.color !== GemColor.GemColorMeta && gemEP > bestGemEP) {
					bestGemEP = gemEP;
					bestGem = gem;
				}
			}

			Gems.GEM_COLORS.forEach(socketColor => {
				if (Gems.gemMatchesSocket(gem, socketColor) && gemEP > bestGemForSocketEP[socketColor]) {
					bestGemForSocketEP[socketColor] = gemEP;
					bestGemForSocket[socketColor] = gem;
				}
			});
		});

		return {
			bestGemForColor: bestGemForColor,
			bestGemForColorEP: bestGemForColorEP,
			bestGemForSocket: bestGemForSocket,
			bestGemForSocketEP: bestGemForSocketEP,
			bestGem: bestGem,
			bestGemEP: bestGemEP,
		};
	}

	private static epUnitStats: UnitStat[] = UnitStat.getAll().filter(stat => {
		if (stat.isStat()) {
			return true;
		} else {
			return [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps, PseudoStat.PseudoStatRangedDps].includes(stat.getPseudoStat());
		}
	});

	private buildStatsTable(): StatsTableEntry[] {
		const copyToCurrentEpText = 'Copy to Current EP';
		const createRefs = () => ({
			metricRef: ref<HTMLButtonElement>(),
			ratioRef: ref<HTMLTableCellElement>(),
		});
		return [
			{
				metric: 'damage',
				type: 'weight',
				label: 'DPS Weight',
				labelTooltip: 'Per-point increase in DPS (Damage Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().dps!.weights,
				...createRefs(),
			},
			{
				metric: 'damage',
				type: 'ep',
				label: 'DPS EP',
				labelTooltip: 'EP (Equivalency Points) for DPS (Damage Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().dps!.epValues,
				getEpRefStat: () => this.getDpsEpRefStat(),
				...createRefs(),
			},
			{
				metric: 'healing',
				type: 'weight',
				label: 'HPS Weight',
				labelTooltip: 'Per-point increase in HPS (Healing Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().hps!.weights,
				...createRefs(),
			},
			{
				metric: 'healing',
				type: 'ep',
				label: 'HPS EP',
				labelTooltip: 'EP (Equivalency Points) for HPS (Healing Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().hps!.epValues,
				getEpRefStat: () => this.getHealEpRefStat(),
				...createRefs(),
			},
			{
				metric: 'threat',
				type: 'weight',
				label: 'TPS Weight',
				labelTooltip: 'Per-point increase in TPS (Threat Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().tps!.weights,
				...createRefs(),
			},
			{
				metric: 'threat',
				type: 'ep',
				label: 'TPS EP',
				labelTooltip: 'EP (Equivalency Points) for TPS (Threat Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().tps!.epValues,
				getEpRefStat: () => this.getDpsEpRefStat(),
				...createRefs(),
			},
			{
				metric: 'threat',
				type: 'weight',
				label: 'DTPS Weight',
				labelTooltip: 'Per-point increase in DTPS (Damage Taken Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().dtps!.weights,
				...createRefs(),
			},
			{
				metric: 'threat',
				type: 'ep',
				label: 'DTPS EP',
				labelTooltip: 'EP (Equivalency Points) for DTPS (Damage Taken Per Second) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().dtps!.epValues,
				getEpRefStat: () => this.getTankEpRefStat(),
				...createRefs(),
			},
			{
				experimental: true,
				metric: 'threat',
				type: 'weight',
				label: 'TMI Weight',
				labelTooltip: 'Per-point decrease in TMI (Theck-Meloree Index) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().tmi!.weights,
				...createRefs(),
			},
			{
				experimental: true,
				metric: 'threat',
				type: 'ep',
				label: 'TMI EP',
				labelTooltip: 'EP (Equivalency Points) for TMI (Theck-Meloree Index) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().tmi!.epValues,
				getEpRefStat: () => this.getTankEpRefStat(),
				...createRefs(),
			},
			{
				experimental: true,
				metric: 'threat',
				type: 'weight',
				label: 'Death Weight',
				labelTooltip: 'Per-point decrease in p(death) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().pDeath!.weights,
				...createRefs(),
			},
			{
				experimental: true,
				metric: 'threat',
				type: 'ep',
				label: 'Death EP',
				labelTooltip: 'EP (Equivalency Points) for p(death) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().pDeath!.epValues,
				getEpRefStat: () => this.getTankEpRefStat(),
				...createRefs(),
			},
			{
				type: 'action',
				label: 'Current EP',
				labelTooltip: 'Current EP Weights. Used to sort the gear selector menus.',
				actionTooltip: 'Restore Default EP',
				getWeights: () => this.simUI.individualConfig.defaults.epWeights.toProto(),
				...createRefs(),
			},
		];
	}
}

type StatsTableEntry = {
	experimental?: boolean;
	metric?: 'damage' | 'healing' | 'threat';
	type: 'ep' | 'weight' | 'action';
	label: string;
	labelTooltip: string;
	actionTooltip: string;
	getWeights: () => UnitStats | undefined;
	getEpRefStat?: () => Stat;
	metricRef: ReturnType<typeof ref<HTMLButtonElement>>;
	ratioRef: ReturnType<typeof ref<HTMLTableCellElement>>;
};

type BestGemsResult = {
	bestGemForColor: Gem[];
	bestGemForColorEP: number[];
	bestGemForSocket: Gem[];
	bestGemForSocketEP: number[];
	bestGem: Gem;
	bestGemEP: number;
};
