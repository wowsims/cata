import clsx from 'clsx';
import tippy from 'tippy.js';

import { SpellType } from '../../proto/api';
import { spellSchoolNames } from '../../proto_utils/names';
import { ActionMetrics } from '../../proto_utils/sim_result';
import { formatToCompactNumber, formatToPercent } from '../../utils';
import { MetricsCombinedTooltipTable } from './metrics_table/metrics_combined_tooltip_table';
import { ColumnSortType, MetricsTable } from './metrics_table/metrics_table';
import { MetricsTotalBar } from './metrics_table/metrics_total_bar';
import { ResultComponentConfig, SimResultData } from './result_component';

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
					cellElem.appendChild(
						<MetricsTotalBar
							spellSchool={metric.spellSchool}
							percentage={metric.totalDamagePercent}
							max={this.maxDtpsAmount}
							total={metric.avgDamage}
							value={metric.damage}
						/>,
					);

					const hitValues = metric.damageDone.hit;
					const critValues = metric.damageDone.crit;

					tippy(cellElem, {
						maxWidth: 'none',
						placement: 'auto',
						theme: 'metrics-table',
						content: (
							<>
								<MetricsCombinedTooltipTable
									spellSchool={metric.spellSchool}
									total={metric.damage}
									totalPercentage={100}
									hasFooter={false}
									values={[
										...(metric.spellType === SpellType.SpellTypeAll
											? [
													{
														name: 'Hit',
														...hitValues,
													},
													{
														name: `Critical Hit`,
														...critValues,
													},
											  ]
											: []),
										...(metric.spellType === SpellType.SpellTypeCast
											? [
													{
														name: 'Hit',
														...hitValues,
													},
													{
														name: `Critical Hit`,
														...critValues,
													},
											  ]
											: []),
										...(metric.spellType === SpellType.SpellTypePeriodic
											? [
													{
														name: 'Tick',
														...hitValues,
													},
													{
														name: `Critical Tick`,
														...critValues,
													},
											  ]
											: []),
										// {
										// 	name: 'Glancing Blow',
										// 	value: metric.glances,
										// 	percentage: metric.glancePercent,
										// },
										// {
										// 	name: 'Blocked Blow',
										// 	value: metric.blocks,
										// 	percentage: metric.blockPercent,
										// },
									]}
								/>
							</>
						),
					});
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
				getValue: (metric: ActionMetrics) => metric.totalMissesPercent,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{metric.totalMissesPercent ? formatToPercent(metric.totalMissesPercent) : '-'}</>);
					if (!metric.totalMissesPercent) return;

					tippy(cellElem, {
						placement: 'right',
						theme: 'metrics-table',
						content: (
							<>
								<MetricsCombinedTooltipTable
									spellSchool={metric.spellSchool}
									total={metric.totalMisses + metric.blocks}
									totalPercentage={metric.totalMissesPercent + metric.blockPercent}
									values={[
										{
											name: 'Miss',
											value: metric.misses,
											percentage: metric.missPercent,
										},
										{
											name: 'Parry',
											value: metric.parries,
											percentage: metric.parryPercent,
										},
										{
											name: 'Dodge',
											value: metric.dodges,
											percentage: metric.dodgePercent,
										},
										{
											name: 'Block',
											value: metric.blocks,
											percentage: metric.blockPercent,
										},
									]}
								/>
							</>
						),
					});
				},
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
		return ActionMetrics.merge(metrics, {
			removeTag: true,
			actionIdOverride: metrics[0].unit?.petActionId || undefined,
		});
	}
}
