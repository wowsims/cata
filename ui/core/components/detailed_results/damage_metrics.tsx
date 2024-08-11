import clsx from 'clsx';
import tippy from 'tippy.js';

import { spellSchoolNames } from '../../proto_utils/names';
import { ActionMetrics } from '../../proto_utils/sim_result';
import { bucket, formatToCompactNumber, formatToNumber, formatToPercent } from '../../utils';
import { MetricsCombinedTooltipTable } from './metrics_table/metrics_combined_tooltip_table';
import { ColumnSortType, MetricsTable } from './metrics_table/metrics_table';
import { MetricsTotalBar } from './metrics_table/metrics_total_bar';
import { ResultComponentConfig, SimResultData } from './result_component';

export class DamageMetricsTable extends MetricsTable<ActionMetrics> {
	maxDamageAmount: number | null = null;
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'damage-metrics-root';
		config.resultsEmitter.on((_, resultData) => {
			const lastResult = resultData
				? this.getGroupedMetrics(resultData)
						.filter(g => g.length)
						.map(groups => this.mergeMetrics(groups))
				: undefined;
			this.maxDamageAmount = Math.max(...(lastResult || []).map(a => a.damage));
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
				name: 'Damage done',
				tooltip: 'Total Damage done',
				headerCellClass: 'text-start',
				getValue: (metric: ActionMetrics) => metric.avgDamage,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.classList.add('metric-total');
					cellElem.appendChild(
						<MetricsTotalBar
							spellSchool={metric.spellSchool}
							percentage={metric.totalDamagePercent}
							max={this.maxDamageAmount}
							total={metric.avgDamage}
							value={metric.damage}
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
				name: 'Avg Cast',
				tooltip: 'Damage / Casts',
				getValue: (metric: ActionMetrics) => metric.avgCast,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgCast),
			},
			{
				name: 'Avg Cast',
				tooltip: 'Threat / Casts',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgCastThreat,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgCastThreat),
			},
			{
				name: 'Hits',
				tooltip: 'Hits',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				// getDisplayString: (metric: ActionMetrics) => formatToNumber(metric.landedHits, { minimumFractionDigits: 1 }),
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					if (!metric.landedHits) return '-';

					cellElem.appendChild(<>{formatToNumber(metric.landedHits)}</>);

					tippy(cellElem, {
						placement: 'right',
						theme: 'metrics-table',
						content: (
							<>
								<MetricsCombinedTooltipTable
									spellSchool={metric.spellSchool}
									total={metric.landedHits}
									totalPercentage={100}
									values={[
										{
											name: metric.isPeriodic ? 'Tick' : 'Hit',
											value: metric.hits,
											percentage: metric.hitPercent,
										},
										{
											name: `Critical ${metric.isPeriodic ? 'Tick' : 'Hit'}`,
											value: metric.crits,
											percentage: metric.critPercent,
										},
										{
											name: 'Glancing Blow',
											value: metric.glances,
											percentage: metric.glancePercent,
										},
										{
											name: 'Blocked Blow',
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
				name: 'Avg Hit',
				tooltip: 'Damage / Hits',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgHit),
			},
			{
				name: 'Avg Hit',
				tooltip: 'Threat / Hits',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgHitThreat,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgHitThreat),
			},
			{
				name: 'Crit %',
				tooltip: 'Crits / Hits',
				getValue: (metric: ActionMetrics) => metric.critPercent,
				getDisplayString: (metric: ActionMetrics) => formatToPercent(metric.critPercent),
			},
			{
				name: 'Miss %',
				tooltip: 'Misses / Casts',
				getValue: (metric: ActionMetrics) => metric.totalMissesPercent,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					if (!metric.totalMissesPercent) return '-';

					cellElem.appendChild(<>{formatToPercent(metric.totalMissesPercent)}</>);

					tippy(cellElem, {
						placement: 'right',
						theme: 'metrics-table',
						content: (
							<>
								<MetricsCombinedTooltipTable
									spellSchool={metric.spellSchool}
									total={metric.totalMisses}
									totalPercentage={metric.totalMissesPercent}
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
									]}
								/>
							</>
						),
					});
				},
			},
			{
				name: 'TPS',
				tooltip: 'Threat / Encounter Duration',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.tps,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.tps),
			},
			{
				name: 'DPS',
				tooltip: 'Damage / Encounter Duration',
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.dps,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.dps),
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

		const actions = player.getDamageActions().map(action => action.forTarget(resultData.filter));
		const actionGroups = ActionMetrics.groupById(actions);

		const petsByName = bucket(player.pets, pet => pet.name);
		const petGroups = Object.values(petsByName).map(pets =>
			ActionMetrics.joinById(
				pets.flatMap(pet => pet.getDamageActions().map(action => action.forTarget(resultData.filter))),
				true,
			),
		);

		return actionGroups.concat(petGroups);
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		return ActionMetrics.merge(metrics, true, metrics[0]?.unit?.petActionId || undefined);
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
