import clsx from 'clsx';

import { spellSchoolNames } from '../../proto_utils/names';
import { ActionMetrics } from '../../proto_utils/sim_result.js';
import { formatToCompactNumber, formatToPercent } from '../../utils';
import { ColumnSortType, MetricsTable } from './metrics_table.jsx';
import { ResultComponentConfig, SimResultData } from './result_component.js';

export class DtpsMetricsTable extends MetricsTable<ActionMetrics> {
	maxDtpsAmount: number | null = null;
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'dtps-metrics-root';
		config.resultsEmitter.on((_, resultData) => {
			const lastResult = resultData
				? this.getGroupedMetrics(resultData)
						.filter(g => g.length)
						.map(groups => this.mergeMetrics(groups))
				: undefined;
			this.maxDtpsAmount = Math.max(...(lastResult || []).map(a => a.damage));
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
				name: 'Damage taken',
				tooltip: 'Total Damage taken',
				headerCellClass: 'text-start',
				getValue: (metric: ActionMetrics) => metric.avgDamage,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.classList.add('metric-total');
					const spellSchoolString = typeof metric.spellSchool === 'number' ? spellSchoolNames.get(metric.spellSchool) : undefined;
					cellElem.appendChild(
						<div className="d-flex gap-1">
							<div className="metrics-total-percentage">{formatToPercent(metric.totalDamagePercent || 0)}</div>
							<div className="metrics-total-bar">
								<div
									className={clsx('metrics-total-bar-fill', spellSchoolString && `spell-school-${spellSchoolString.toLowerCase()}`)}
									style={{ '--percentage': formatToPercent((metric.damage / (this.maxDtpsAmount ?? 1)) * 100) }}></div>
							</div>
							<div className="metrics-total-damage">{formatToCompactNumber(metric.avgDamage)}</div>
						</div>,
					);
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
				name: 'Hits',
				tooltip: 'Hits + Crits + Glances + Blocks',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				getDisplayString: (metric: ActionMetrics) => metric.landedHits.toFixed(1),
			},
			{
				name: 'Avg Hit',
				tooltip: 'Damage / (Hits + Crits + Glances + Blocks)',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				getDisplayString: (metric: ActionMetrics) => metric.avgHit.toFixed(1),
			},
			{
				name: 'Miss %',
				tooltip: 'Misses / Swings',
				getValue: (metric: ActionMetrics) => metric.missPercent,
				getDisplayString: (metric: ActionMetrics) => metric.missPercent.toFixed(2) + '%',
			},
			{
				name: 'Dodge %',
				tooltip: 'Dodges / Swings',
				getValue: (metric: ActionMetrics) => metric.dodgePercent,
				getDisplayString: (metric: ActionMetrics) => metric.dodgePercent.toFixed(2) + '%',
			},
			{
				name: 'Parry %',
				tooltip: 'Parries / Swings',
				getValue: (metric: ActionMetrics) => metric.parryPercent,
				getDisplayString: (metric: ActionMetrics) => metric.parryPercent.toFixed(2) + '%',
			},
			{
				name: 'Block %',
				tooltip: 'Blocks / Swings',
				getValue: (metric: ActionMetrics) => metric.blockPercent,
				getDisplayString: (metric: ActionMetrics) => metric.blockPercent.toFixed(2) + '%',
			},
			{
				name: 'Crit %',
				tooltip: 'Crits / Swings',
				getValue: (metric: ActionMetrics) => metric.critPercent,
				getDisplayString: (metric: ActionMetrics) => metric.critPercent.toFixed(2) + '%',
			},
			{
				name: 'DTPS',
				tooltip: 'Damage Taken / Encounter Duration',
				sort: ColumnSortType.Descending,
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				getValue: (metric: ActionMetrics) => metric.dps,
				getDisplayString: (metric: ActionMetrics) => metric.dps.toFixed(1),
			},
		]);
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<ActionMetrics>> {
		const players = resultData.result.getRaidIndexedPlayers(resultData.filter);
		if (players.length != 1) {
			return [];
		}
		const player = players[0];

		const targets = resultData.result.getTargets(resultData.filter);
		const targetActions = targets.map(target => target.getDamageActions().map(action => action.forTarget({ player: player.unitIndex }))).flat();
		const actionGroups = ActionMetrics.groupById(targetActions);

		return actionGroups;
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		// TODO: Use NPC ID here instead of pet ID.
		return ActionMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}
}
