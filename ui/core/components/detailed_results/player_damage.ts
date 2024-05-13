import tippy from 'tippy.js';

import { SimResult, SimResultFilter,UnitMetrics } from '../../proto_utils/sim_result.js';
import { maxIndex, sum } from '../../utils.js';
import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';
import { ResultsFilter } from './results_filter.js';
import { SourceChart } from './source_chart.js';

export class PlayerDamageMetricsTable extends MetricsTable<UnitMetrics> {
	private readonly resultsFilter: ResultsFilter;

	// Cached values from most recent result.
	private raidDps: number;
	private maxDps: number;

	constructor(config: ResultComponentConfig, resultsFilter: ResultsFilter) {
		config.rootCssClass = 'player-damage-metrics-root';
		super(config, [
			MetricsTable.playerNameCellConfig(),
			{
				name: 'Amount',
				tooltip: 'Player Damage / Raid Damage',
				headerCellClass: 'amount-header-cell',
				fillCell: (player: UnitMetrics, cellElem: HTMLElement, rowElem: HTMLElement) => {
					cellElem.classList.add('amount-cell');

					let chart: HTMLElement | null = null;
					const makeChart = () => {
						const chartContainer = document.createElement('div');
						rowElem.appendChild(chartContainer);
						const sourceChart = new SourceChart(chartContainer, player.actions);
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

					const playerDps = this.getPlayerDps(player)
					cellElem.innerHTML = `
						<div class="player-damage-percent">
							<span>${(playerDps / this.raidDps * 100).toFixed(2)}%</span>
						</div>
						<div class="player-damage-bar-container">
							<div class="player-damage-bar bg-${player.classColor}" style="width:${playerDps / this.maxDps * 100}%"></div>
						</div>
						<div class="player-damage-total">
							<span>${(playerDps * this.getLastSimResult().result.duration / 1000).toFixed(1)}k</span>
						</div>
					`;
				},
			},
			{
				name: 'DPS',
				tooltip: 'Damage / Encounter Duration',
				columnClass: 'dps-cell',
				sort: ColumnSortType.Descending,
				getValue: (player: UnitMetrics) => this.getPlayerDps(player),
				getDisplayString: (player: UnitMetrics) => this.getPlayerDps(player).toFixed(1),
			},
		]);
		this.resultsFilter = resultsFilter;
		this.raidDps = 0;
		this.maxDps = 0;
	}

	private getPlayerDps(player:UnitMetrics): number {
		const playerActions = player.getPlayerAndPetActions().map(action => action.forTarget(this.resultsFilter.getFilter())).flat();
		const playerDps = sum(playerActions.map(action => action.dps))
		return playerDps
	}

	customizeRowElem(player: UnitMetrics, rowElem: HTMLElement) {
		rowElem.classList.add('player-damage-row');
		rowElem.addEventListener('click', event => {
			this.resultsFilter.setPlayer(this.getLastSimResult().eventID, player.index);
		});
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<UnitMetrics>> {
		const players = resultData.result.getPlayers(resultData.filter);

		const targetActions = players.map(player => player.getPlayerAndPetActions().map(action => action.forTarget(resultData.filter))).flat();

		this.raidDps = sum(targetActions.map(action => action.dps));
		const maxDpsIndex = maxIndex(players.map(player => {
			const targetActions = player.getPlayerAndPetActions().map(action => action.forTarget(resultData.filter)).flat();
			return sum(targetActions.map(action => action.dps));
		}))!;

		const maxDpsTargetActions = players[maxDpsIndex].getPlayerAndPetActions().map(action => action.forTarget(resultData.filter)).flat();
		this.maxDps = sum(maxDpsTargetActions.map(action => action.dps));

		return players.map(player => [player]);
	}
}
