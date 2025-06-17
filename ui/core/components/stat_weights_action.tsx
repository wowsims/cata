import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { CURRENT_API_VERSION } from '../constants/other';
import { IndividualSimUI } from '../individual_sim_ui.jsx';
import { Player } from '../player.js';
import { ProgressMetrics, StatWeightsResult, StatWeightValues } from '../proto/api';
import { PseudoStat, Stat, UnitStats } from '../proto/common.js';
import { SavedStatWeightSettings } from '../proto/ui';
import { getStatName } from '../proto_utils/names.js';
import { Stats, UnitStat } from '../proto_utils/stats.js';
import { RequestTypes } from '../sim_signal_manager';
import { SimUI } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event.js';
import { sanitizeId, stDevToConf90 } from '../utils.js';
import { BaseModal } from './base_modal.jsx';
import { BooleanPicker } from './pickers/boolean_picker.js';
import { NumberPicker } from './pickers/number_picker.js';
import { ResultsViewer } from './results_viewer.jsx';
import { renderSavedEPWeights } from './saved_data_managers/ep_weights';
export class StatWeightActionSettings {
	private readonly storageKey: string;
	readonly changeEmitter = new TypedEvent<void>();

	_excludedStats: Stat[] = [];
	_excludedPseudoStats: PseudoStat[] = [];

	constructor(simUI: SimUI) {
		this.storageKey = simUI.getStorageKey('__statweight_settings__');
		this.changeEmitter.on(() => {
			const json = SavedStatWeightSettings.toJsonString(this.toProto());
			window.localStorage.setItem(this.storageKey, json);
		});
	}

	set excludedStats(value: Stat[]) {
		this._excludedStats = value;
	}
	get excludedStats(): Stat[] {
		return this._excludedStats.slice();
	}

	set excludedPseudoStats(value: PseudoStat[]) {
		this._excludedPseudoStats = value;
	}
	get excludedPseudoStats(): PseudoStat[] {
		return this._excludedPseudoStats.slice();
	}

	static updateProtoVersion(_: SavedStatWeightSettings) {
		// No-op, as there are no proto version migrations currently
	}

	applyDefaults(eventID: EventID) {
		this.excludedStats = [];
		this.excludedPseudoStats = [];
		this.changeEmitter.emit(eventID);
	}

	load(eventID: EventID) {
		const storageValue = window.localStorage.getItem(this.storageKey);
		if (storageValue) {
			const settingsProto = SavedStatWeightSettings.fromJsonString(storageValue, { ignoreUnknownFields: true });
			StatWeightActionSettings.updateProtoVersion(settingsProto);

			const { excludedStats, excludedPseudoStats } = settingsProto;
			this.excludedStats = excludedStats || [];
			this.excludedPseudoStats = excludedPseudoStats || [];
			this.changeEmitter.emit(eventID);
		}
	}

	toProto(): SavedStatWeightSettings {
		return SavedStatWeightSettings.create({
			apiVersion: CURRENT_API_VERSION,
			excludedStats: this.excludedStats,
			excludedPseudoStats: this.excludedPseudoStats,
		});
	}

	/**
	 * Check if a stat should be excluded from weight calculation.
	 * @param stat
	 * @returns true if stat should be excluded.
	 */
	isStatExcludedFromCalc(stat: Stat): boolean {
		return !!this.excludedStats.includes(stat);
	}

	/**
	 * Check if a pseudostat should be excluded from weight calculation.
	 * @param pseudoStat
	 * @returns true if pseudostat should be excluded.
	 */
	isPseudoStatExcludedFromCalc(pseudoStat: PseudoStat): boolean {
		return !!this.excludedPseudoStats.includes(pseudoStat);
	}

	/**
	 * Check if a unitstat should be excluded from weight calculation.
	 * @param unitstat
	 * @returns true if unitstat should be excluded.
	 */
	isUnitStatExcludedFromCalc(unitstat: UnitStat): boolean {
		return unitstat.isStat() ? this.isStatExcludedFromCalc(unitstat.getStat()) : this.isPseudoStatExcludedFromCalc(unitstat.getPseudoStat());
	}

	/**
	 * Set whether a stat should be excluded from calculation.
	 * @param stat
	 * @param exclude
	 */
	setStatExcluded(eventID: EventID, stat: UnitStat, exclude: boolean) {
		const updateStatEntry = <T extends Stat | PseudoStat>(s: T, target: T[]) => {
			const currentIdx = target.indexOf(s);
			if (exclude) {
				if (currentIdx === -1) target.push(s);
			} else if (currentIdx !== -1) {
				target.splice(currentIdx, 1);
			}
			return target;
		};
		if (stat.isStat()) {
			this.excludedStats = updateStatEntry(stat.getStat(), this.excludedStats);
		} else {
			this.excludedPseudoStats = updateStatEntry(stat.getPseudoStat(), this.excludedPseudoStats);
		}
		this.changeEmitter.emit(eventID);
	}
}

export const addStatWeightsAction = (simUI: IndividualSimUI<any>, settings: StatWeightActionSettings) => {
	const epWeightsModal = new EpWeightsMenu(simUI, settings);
	simUI.addAction('Stat Weights', 'ep-weights-action', () => {
		epWeightsModal.open();
	});

	return epWeightsModal;
};

// Create the config for modal in separate function, as constructor cannot
// contain any logic before `super' call. Use modal-xl to accommodate the extra
// TMI & p(death) EP in the UI.
const getModalConfig = (simUI: IndividualSimUI<any>) => {
	const baseConfig = { footer: true, scrollContents: true };
	if (simUI.sim.getShowThreatMetrics()) return { size: 'xl' as const, ...baseConfig };
	return baseConfig;
};

export const scaledEpValue = (stat: UnitStat, epRatios: number[], result: StatWeightsResult | null): number => {
	if (!result) return 0;

	return (
		(result.dps?.epValues ? epRatios[0] * stat.getProtoValue(result.dps.epValues) : 0) +
		(result.hps?.epValues ? epRatios[1] * stat.getProtoValue(result.hps.epValues) : 0) +
		(result.tps?.epValues ? epRatios[2] * stat.getProtoValue(result.tps.epValues) : 0) +
		(result.dtps?.epValues ? epRatios[3] * stat.getProtoValue(result.dtps.epValues) : 0) +
		(result.tmi?.epValues ? epRatios[4] * stat.getProtoValue(result.tmi.epValues) : 0) +
		(result.pDeath?.epValues ? epRatios[5] * stat.getProtoValue(result.pDeath.epValues) : 0)
	);
};

export class EpWeightsMenu extends BaseModal {
	private readonly simUI: IndividualSimUI<any>;
	private readonly container: HTMLElement;
	private readonly sidebar: HTMLElement;
	private readonly table: HTMLElement;
	private readonly tableBody: HTMLElement;
	private readonly resultsViewer: ResultsViewer;
	private readonly settings: StatWeightActionSettings;

	private statsType: string;
	private epStats: Stat[];
	private epPseudoStats: PseudoStat[];
	private epReferenceStat: Stat;
	private showAllStats = false;

	constructor(simUI: IndividualSimUI<any>, settings: StatWeightActionSettings) {
		super(simUI.rootElem, 'ep-weights-menu', { ...getModalConfig(simUI), disposeOnClose: false });
		this.header?.insertAdjacentElement('afterbegin', <h5 className="modal-title">Calculate Stat Weights</h5>);

		this.simUI = simUI;
		this.statsType = 'ep';
		this.epStats = this.simUI.individualConfig.epStats;
		this.epPseudoStats = this.simUI.individualConfig.epPseudoStats || [];
		this.epReferenceStat = this.simUI.individualConfig.epReferenceStat;
		this.settings = settings;

		const statsTable = this.buildStatsTable();
		const containerRef = ref<HTMLDivElement>();
		const sidebarRef = ref<HTMLDivElement>();
		const tableRef = ref<HTMLTableElement>();
		const tableBodyRef = ref<HTMLTableSectionElement>();
		const damageMetricsSelectRef = ref<HTMLSelectElement>();
		const healingMetricsSelectRef = ref<HTMLSelectElement>();
		const threatMetricsSelectRef = ref<HTMLSelectElement>();
		const typeSelectRef = ref<HTMLSelectElement>();
		const computeEpRef = ref<HTMLButtonElement>();
		const calcWeightsButtonRef = ref<HTMLButtonElement>();
		const allStatsContainerRef = ref<HTMLDivElement>();

		const getNameFromStat = (stat: Stat | undefined) => (stat !== undefined ? getStatName(stat) : '??');
		const getStatFromName = (value: string) => Object.values(this.epStats).find(stat => getNameFromStat(stat) === value);
		const epRefSelectOptions = (
			<>
				{this.epStats.map(stat => (
					<option>{getNameFromStat(stat)}</option>
				))}
			</>
		);

		this.body.appendChild(
			<div className="d-flex flex-column flex-lg-row align-items-lg-start gap-3">
				<div className="ep-weights-content order-1 order-lg-0">
					<div className="ep-weights-options row">
						<div className="col col-sm-3">
							<select ref={typeSelectRef} className="ep-type-select form-select">
								<option value="ep">EP</option>
								<option value="weight">Weights</option>
							</select>
						</div>
						<div ref={allStatsContainerRef} className="show-all-stats-container col col-sm-3"></div>
					</div>
					<div className="ep-reference-options row">
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
									<th>Update</th>
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
									<td></td>
									{statsTable
										.filter(({ type }) => type !== 'action')
										.map(({ metric, type, ratioRef }) => (
											<td ref={ratioRef} className={clsx('type-ratio', `${metric}-metrics`, `type-${type}`)} />
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
				</div>
				<div ref={sidebarRef} className="ep-weights-sidebar sticky-lg-top order-0 order-lg-1" />
			</div>,
		);

		this.footer!.appendChild(
			<>
				<button ref={calcWeightsButtonRef} className="btn btn-primary calc-weights">
					<i className="fas fa-calculator me-1" />
					Calculate
				</button>
			</>,
		);

		this.container = containerRef.value!;
		this.sidebar = sidebarRef.value!;
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

		const calcButton = calcWeightsButtonRef.value;
		let isRunning = false;
		calcButton?.addEventListener('click', async () => {
			if (isRunning) return;
			isRunning = true;

			try {
				await this.simUI.sim.signalManager.abortType(RequestTypes.All);
			} catch (error) {
				console.error(error);
				return;
			}

			calcButton.disabled = true;
			this.simUI.rootElem.classList.add('blurred');
			this.simUI.rootElem.insertAdjacentElement('afterend', pendingDiv);

			this.container.scrollTo({ top: 0 });
			this.container.classList.add('pending');
			this.resultsViewer.setPending();
			const iterations = this.simUI.sim.getIterations();

			let waitAbort = false;
			this.resultsViewer.addAbortButton(async () => {
				if (waitAbort) return;
				try {
					waitAbort = true;
					await simUI.sim.signalManager.abortType(RequestTypes.StatWeights);
				} catch (error) {
					console.error('Error on stat weight abort!');
					console.error(error);
				} finally {
					waitAbort = false;
					if (!isRunning) calcButton.disabled = false;
				}
			});

			const epStatsToCalc = this.epStats.filter(s => !this.settings.isStatExcludedFromCalc(s));
			const epPseudoStatsToCalc = this.epPseudoStats.filter(ps => !this.settings.isPseudoStatExcludedFromCalc(ps));

			const result = await this.simUI.player.computeStatWeights(
				TypedEvent.nextEventID(),
				epStatsToCalc,
				epPseudoStatsToCalc,
				this.epReferenceStat,
				progress => {
					this.setSimProgress(progress);
				},
			);
			this.simUI.rootElem.classList.remove('blurred');
			pendingDiv.remove();
			this.container.classList.remove('pending');
			this.resultsViewer.hideAll();
			isRunning = false;
			if (!waitAbort) calcButton.disabled = false;

			if (!result) return;
			this.simUI.prevEpIterations = iterations;
			this.simUI.prevEpSimResult = this.calculateEp(result);
			this.updateTable();
		});

		this.addOnHideCallback(() => {
			this.simUI.sim.signalManager.abortType(RequestTypes.StatWeights).catch(console.error);
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
				this.setEpWeightsWithoutExcluded(Stats.fromProto(weightsFunc()));
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
				this.setEpWeightsWithoutExcluded(newEp);
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
				this.setEpWeightsWithoutExcluded(newWeights);
			}
			this.updateTable();
		});

		this.buildSavedEPWeightsPicker();
	}

	/**
	 * Set new ep weights while leaving excluded stats at their old value.
	 * @param newWeights
	 */
	private setEpWeightsWithoutExcluded(newWeights: Stats) {
		const { excludedStats, excludedPseudoStats } = this.settings;
		const oldWeights = this.simUI.player.getEpWeights();
		for (const stat of excludedStats) {
			newWeights = newWeights.withStat(stat, oldWeights.getStat(stat));
		}
		for (const pseudoStat of excludedPseudoStats) {
			newWeights = newWeights.withPseudoStat(pseudoStat, oldWeights.getPseudoStat(pseudoStat));
		}
		this.simUI.player.setEpWeights(TypedEvent.nextEventID(), newWeights);
	}

	/**
	 * Check if a specific stat is included in the EP stats for this spec.
	 * @param stat
	 * @returns
	 */
	private isEpStat(stat: UnitStat) {
		if (stat.isStat()) return this.epStats.includes(stat.getStat());
		return this.epPseudoStats.includes(stat.getPseudoStat());
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
		const result = !this.settings.isUnitStatExcludedFromCalc(stat) ? this.simUI.prevEpSimResult : null;
		const epRatios = this.simUI.player.getEpRatios();

		const rowTotalEp = scaledEpValue(stat, epRatios, result);
		const currentEpRef = ref<HTMLTableCellElement>();
		const includeToggleRef = ref<HTMLTableCellElement>();
		const row = (
			<tr>
				<td>{stat.getFullName(this.simUI.player.getClass())}</td>
				<td ref={includeToggleRef} className="swcalc-include-toggle"></td>
				{this.makeTableRowCells(stat, result?.dps, 'damage-metrics', rowTotalEp, epRatios[0])}
				{this.makeTableRowCells(stat, result?.hps, 'healing-metrics', rowTotalEp, epRatios[1])}
				{this.makeTableRowCells(stat, result?.tps, 'threat-metrics', rowTotalEp, epRatios[2])}
				{this.makeTableRowCells(stat, result?.dtps, 'threat-metrics', rowTotalEp, epRatios[3])}
				{this.makeTableRowCells(stat, result?.tmi, 'threat-metrics', rowTotalEp, epRatios[4])}
				{this.makeTableRowCells(stat, result?.pDeath, 'threat-metrics', rowTotalEp, epRatios[5])}
				<td ref={currentEpRef} className="current-ep"></td>
			</tr>
		) as HTMLElement;

		if (includeToggleRef.value && this.isEpStat(stat)) {
			new BooleanPicker(includeToggleRef.value, this, {
				id: 'sw-stat-toggle-' + stat.getFullName(this.simUI.player.getClass()),
				getValue: epWeightsModal => !epWeightsModal.settings.isUnitStatExcludedFromCalc(stat),
				setValue: (eventID, epWeightsModal, newValue) => epWeightsModal.settings.setStatExcluded(eventID, stat, !newValue),
				changedEvent: epWeightsModal => epWeightsModal.settings.changeEmitter,
				enableWhen: epWeightsModal => !stat.isStat() || epWeightsModal.epReferenceStat != stat.getStat(),
			});
		}

		const currentEpCell = currentEpRef.value!;

		new NumberPicker(currentEpCell, this.simUI.player, {
			id: `ep-weight-stat-${sanitizeId(stat.getShortName(this.simUI.player.playerClass.classID))}`,
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

	private static epUnitStats: UnitStat[] = UnitStat.getAll().filter(stat => {
		if (stat.isStat()) {
			return true;
		} else {
			return [
				PseudoStat.PseudoStatMainHandDps,
				PseudoStat.PseudoStatOffHandDps,
				PseudoStat.PseudoStatRangedDps,
				PseudoStat.PseudoStatPhysicalHitPercent,
				PseudoStat.PseudoStatSpellHitPercent,
				PseudoStat.PseudoStatPhysicalCritPercent,
				PseudoStat.PseudoStatSpellCritPercent,
			].includes(stat.getPseudoStat());
		}
	});

	private buildSavedEPWeightsPicker() {
		renderSavedEPWeights(this.sidebar, this.simUI);
	}

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
				metric: 'threat',
				type: 'weight',
				label: 'TMI Weight',
				labelTooltip: 'Per-point decrease in TMI (Theck-Meloree Index) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().tmi!.weights,
				...createRefs(),
			},
			{
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
				metric: 'threat',
				type: 'weight',
				label: 'Death Weight',
				labelTooltip: 'Per-point decrease in p(death) for each stat.',
				actionTooltip: copyToCurrentEpText,
				getWeights: () => this.getPrevSimResult().pDeath!.weights,
				...createRefs(),
			},
			{
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
