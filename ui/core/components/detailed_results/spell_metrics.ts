import { ActionMetrics, SimResult, SimResultFilter } from '../../proto_utils/sim_result.js';
import { bucket } from '../../utils.js';
import { ColumnSortType, MetricsTable } from './metrics_table/metrics_table.jsx';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class SpellMetricsTable extends MetricsTable<ActionMetrics> {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'spell-metrics-root';
		super(config, [
			MetricsTable.nameCellConfig((metric: ActionMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
					metricType: metric.constructor?.name,
				};
			}),
			{
				name: 'Damage done',
				tooltip: 'Total Damage done',
				getValue: (metric: ActionMetrics) => metric.damage,
				getDisplayString: (metric: ActionMetrics) => {
					const metricDamageDone = new Intl.NumberFormat('en-US', { notation: 'compact', maximumFractionDigits: 2 }).format(metric.avgDamage);
					console.log({ metricDamageDone, totalDamagePercent: metric.totalDamagePercent });
					return metricDamageDone;
				},
			},
			{
				name: 'Casts',
				tooltip: 'Casts',
				getValue: (metric: ActionMetrics) => metric.casts,
				getDisplayString: (metric: ActionMetrics) => metric.casts.toFixed(1),
			},
			{
				name: 'Avg Cast',
				tooltip: 'Damage / Casts',
				getValue: (metric: ActionMetrics) => metric.avgCast,
				getDisplayString: (metric: ActionMetrics) => metric.avgCast.toFixed(1),
			},
			{
				name: 'Avg Cast',
				tooltip: 'Threat / Casts',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgCastThreat,
				getDisplayString: (metric: ActionMetrics) => metric.avgCastThreat.toFixed(1),
			},
			{
				name: 'Hits',
				tooltip: 'Hits',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				getDisplayString: (metric: ActionMetrics) => metric.landedHits.toFixed(1),
			},
			{
				name: 'Avg Hit',
				tooltip: 'Damage / Hits',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				getDisplayString: (metric: ActionMetrics) => metric.avgHit.toFixed(1),
			},
			{
				name: 'Avg Hit',
				tooltip: 'Threat / Hits',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgHitThreat,
				getDisplayString: (metric: ActionMetrics) => metric.avgHitThreat.toFixed(1),
			},
			{
				name: 'Crit %',
				tooltip: 'Crits / Hits',
				getValue: (metric: ActionMetrics) => metric.critPercent,
				getDisplayString: (metric: ActionMetrics) => metric.critPercent.toFixed(2) + '%',
			},
			{
				name: 'Miss %',
				tooltip: 'Misses / Casts',
				getValue: (metric: ActionMetrics) => metric.missPercent,
				getDisplayString: (metric: ActionMetrics) => metric.missPercent.toFixed(2) + '%',
			},
			{
				name: 'DPS',
				tooltip: 'Damage / Encounter Duration',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.dps,
				getDisplayString: (metric: ActionMetrics) => metric.dps.toFixed(1),
			},
			{
				name: 'TPS',
				tooltip: 'Threat / Encounter Duration',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.tps,
				getDisplayString: (metric: ActionMetrics) => metric.tps.toFixed(1),
			},
		]);
	}

	customizeRowElem(action: ActionMetrics, rowElem: HTMLElement) {
		if (action.hitAttempts == 0 && action.dps == 0) {
			rowElem.classList.add('threat-metrics');
		}
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<ActionMetrics>> {
		const players = resultData.result.getRaidIndexedPlayers(resultData.filter);
		if (players.length != 1) {
			return [];
		}
		const player = players[0];

		const actions = player.getSpellDamageActions().map(action => action.forTarget(resultData.filter));
		const actionGroups = ActionMetrics.groupById(actions);

		const petsByName = bucket(player.pets, pet => pet.name);
		const petGroups = Object.values(petsByName).map(pets =>
			ActionMetrics.joinById(
				pets.flatMap(pet => pet.getSpellDamageActions().map(action => action.forTarget(resultData.filter))),
				true,
			),
		);

		return actionGroups.concat(petGroups);
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		return ActionMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
