import tippy from 'tippy.js';

import { SimResult, SimResultFilter,UnitMetrics } from '../../proto_utils/sim_result.js';
import { maxIndex, sum } from '../../utils.js';
import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';
import { ResultsFilter } from './results_filter.js';
import { SourceChart } from './source_chart.js';

export class PlayerDamageTakenMetricsTable extends MetricsTable<UnitMetrics> {
	private readonly resultsFilter: ResultsFilter;

	// Cached values from most recent result.
	private resultData: SimResultData | undefined;
	private raidDtps: number;
	private maxDtps: number;

	constructor(config: ResultComponentConfig, resultsFilter: ResultsFilter) {
		config.rootCssClass = 'player-damage-taken-metrics-root';
		super(config, [
			MetricsTable.playerNameCellConfig(),
			{
				name: 'Amount',
				tooltip: 'Player Damage Taken / Raid Damage Taken',
				headerCellClass: 'amount-header-cell',
				fillCell: (player: UnitMetrics, cellElem: HTMLElement, rowElem: HTMLElement) => {
					cellElem.classList.add('amount-cell');

					let chart: HTMLElement | null = null;
					const makeChart = () => {
						const chartContainer = document.createElement('div');
						rowElem.appendChild(chartContainer);
						if (this.resultData) {
							const targets = this.resultData.result.getTargets(this.resultData.filter);
							const playerFilter = {
								player: player.unitIndex,
							}
							const targetActions = targets.map(target => target.getMeleeActions().concat(target.getSpellActions()).map(action => action.forTarget(playerFilter))).flat();
							const sourceChart = new SourceChart(chartContainer, targetActions);
						}
						return chartContainer;
					};

					tippy(rowElem, {
						content: 'Loading...',
						placement: 'bottom',
						ignoreAttributes: true,
						onShow(instance: any) {
							if (!chart) {
								chart = makeChart();
								instance.setContent(chart);
							}
						},
					});

					const playerDtps = this.getPlayerDtps(player)
					cellElem.innerHTML = `
						<div class="player-damage-percent">
							<span>${(playerDtps / this.raidDtps * 100).toFixed(2)}%</span>
						</div>
						<div class="player-damage-bar-container">
							<div class="player-damage-bar bg-${player.classColor}" style="width:${playerDtps / this.maxDtps * 100}%"></div>
						</div>
						<div class="player-damage-total">
							<span>${(playerDtps * this.resultData!.result.duration / 1000).toFixed(1)}k</span>
						</div>
					`;
				},
			},
			{
				name: 'DTPS',
				tooltip: 'Damage Taken / Encounter Duration',
				columnClass: 'dps-cell',
				sort: ColumnSortType.Descending,
				getValue: (player: UnitMetrics) => this.getPlayerDtps(player),
				getDisplayString: (player: UnitMetrics) => this.getPlayerDtps(player).toFixed(1),
			},
		]);
		this.resultsFilter = resultsFilter;
		this.raidDtps = 0;
		this.maxDtps = 0;
	}

	private getPlayerDtps(player: UnitMetrics): number {
		const targets = this.resultData!.result.getTargets(this.resultData!.filter);
		const targetActions = targets.map(target => target.getPlayerAndPetActions().map(action => action.forTarget({ player: player.unitIndex }))).flat();
		const playerDtps = sum(targetActions.map(action => action.dps))
		return playerDtps
	}

	customizeRowElem(player: UnitMetrics, rowElem: HTMLElement) {
		rowElem.classList.add('player-damage-row');
		rowElem.addEventListener('click', event => {
			this.resultsFilter.setPlayer(this.getLastSimResult().eventID, player.index);
		});
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<UnitMetrics>> {
		this.resultData = resultData;
		const players = resultData.result.getPlayers(resultData.filter);

		const targets = resultData.result.getTargets(resultData.filter);
		const targetActions = targets.map(target => target.getPlayerAndPetActions().map(action => action.forTarget(resultData.filter))).flat();

		this.raidDtps = sum(targetActions.map(action => action.dps));
		const maxDpsIndex = maxIndex(players.map(player => {
			const targetActions = targets.map(target => target.getPlayerAndPetActions().map(action => action.forTarget({ player: player.unitIndex }))).flat();
			return sum(targetActions.map(action => action.dps))
		}))!;

		const maxDtpsTargetActions = targets.map(target => target.getPlayerAndPetActions().map(action => action.forTarget({ player: players[maxDpsIndex].unitIndex }))).flat();
		this.maxDtps = sum(maxDtpsTargetActions.map(action => action.dps));

		return players.map(player => [player]);
	}
}
