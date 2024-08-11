import clsx from 'clsx';

import { spellSchoolNames } from '../../proto_utils/names';
import { ActionMetrics } from '../../proto_utils/sim_result.js';
import { formatToCompactNumber, formatToNumber, formatToPercent } from '../../utils.js';
import { ColumnSortType, MetricsTable } from './metrics_table/metrics_table.jsx';
import { MetricsTotalBar } from './metrics_table/metrics_total_bar';
import { ResultComponentConfig, SimResultData } from './result_component.js';

export class HealingMetricsTable extends MetricsTable<ActionMetrics> {
	maxHealingAmount: number | null = null;
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'healing-metrics-root';
		config.resultsEmitter.on((_, resultData) => {
			const lastResult = resultData
				? this.getGroupedMetrics(resultData)
						.filter(g => g.length)
						.map(groups => this.mergeMetrics(groups))
				: undefined;
			this.maxHealingAmount = Math.max(...(lastResult || []).map(a => a.healing));
		});
		super(config, [
			MetricsTable.nameCellConfig((metric: ActionMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
					metricType: metric.constructor?.name,
				};
			}),
			{
				name: 'Healing done',
				tooltip: 'Total Healing done',
				headerCellClass: 'text-start',
				getValue: (metric: ActionMetrics) => metric.avgHealing,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.classList.add('metric-total');
					cellElem.appendChild(
						<MetricsTotalBar
							spellSchool={metric.spellSchool}
							percentage={metric.totalHealingPercent}
							max={this.maxHealingAmount}
							total={metric.avgHealing}
							value={metric.healing}
							overlayValue={metric.shielding}
						/>,
					);
				},
			},
			{
				name: 'Casts',
				tooltip: 'Casts',
				getValue: (metric: ActionMetrics) => metric.casts,
				getDisplayString: (metric: ActionMetrics) => formatToNumber(metric.casts, { minimumFractionDigits: 1 }),
			},
			{
				name: 'CPM',
				tooltip: 'Casts / (Encounter Duration / 60 Seconds)',
				getValue: (metric: ActionMetrics) => metric.castsPerMinute,
				getDisplayString: (metric: ActionMetrics) => formatToNumber(metric.castsPerMinute, { minimumFractionDigits: 1 }),
			},
			{
				name: 'Cast Time',
				tooltip: 'Average cast time in seconds',
				getValue: (metric: ActionMetrics) => metric.avgCastTimeMs,
				getDisplayString: (metric: ActionMetrics) => (metric.avgCastTimeMs / 1000).toFixed(2),
			},
			{
				name: 'HPM',
				tooltip: 'Healing / Mana',
				getValue: (metric: ActionMetrics) => metric.hpm,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.hpm),
			},
			{
				name: 'Avg Cast',
				tooltip: 'Healing / Casts',
				getValue: (metric: ActionMetrics) => metric.avgCastHealing,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgCastHealing),
			},

			{
				name: 'Avg Cast',
				tooltip: 'Threat / Casts',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgCastThreat,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgCastThreat),
			},
			{
				name: 'Crit %',
				tooltip: 'Crits / Hits',
				getValue: (metric: ActionMetrics) => metric.critPercent,
				getDisplayString: (metric: ActionMetrics) => formatToPercent(metric.critPercent / 1000),
			},
			{
				name: 'HPET',
				tooltip: 'Healing / Avg Cast Time',
				getValue: (metric: ActionMetrics) => metric.healingThroughput,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.healingThroughput),
			},
			{
				name: 'TPS',
				tooltip: 'Threat / Encounter Duration',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.tps,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.tps),
			},
			{
				name: 'HPS',
				tooltip: 'Healing / Encounter Duration',
				sort: ColumnSortType.Descending,
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				getValue: (metric: ActionMetrics) => metric.hps,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.hps),
			},
		]);
	}

	customizeRowElem(action: ActionMetrics, rowElem: HTMLElement) {
		if (action.hitAttempts == 0 && action.hps == 0) {
			rowElem.classList.add('threat-metrics');
		}
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<ActionMetrics>> {
		const players = resultData.result.getRaidIndexedPlayers(resultData.filter);
		if (players.length != 1) {
			return [];
		}
		const player = players[0];

		//const actions = player.getSpellActions().map(action => action.forTarget(resultData.filter));
		// TODO: Do we want to show 0 hps metrics in here? Make it conditional for healing sims
		// in case they want to show the threat for non healing spells
		const actions = player.getHealingActions().filter(action => action.hps > 0);
		const actionGroups = ActionMetrics.groupById(actions);

		return actionGroups;
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		return ActionMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
