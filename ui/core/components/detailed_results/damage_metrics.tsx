import { SpellType } from '../../proto/api';
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
				headerCellClass: 'text-center',
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

					const hitValues = metric.damageDone.hit;
					const critValues = metric.damageDone.crit;

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						headerValues={[, 'Amount']}
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
					/>;
				},
			},
			{
				name: 'Casts',
				tooltip: 'Casts',
				getValue: (metric: ActionMetrics, _isChildRow) => {
					if (metric.isProc) return 0;
					return metric.casts;
				},
				getDisplayString: (metric: ActionMetrics, _isChildRow) => {
					if (metric.isProc) return '-';
					return formatToNumber(metric.casts);
				},
			},
			{
				name: 'Avg Cast',
				tooltip: 'Damage / Casts',
				getValue: (metric: ActionMetrics) => metric.avgCast,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToCompactNumber(metric.avgCast)}</>);

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						tooltipConfig={{
							onShow: () => {
								const hideThreatMetrics = !!document.querySelector('.hide-threat-metrics');
								if (hideThreatMetrics) return false;
							},
						}}
						headerValues={[, 'Amount']}
						spellSchool={metric.spellSchool}
						total={metric.avgCastThreat}
						totalPercentage={100}
						hasFooter={false}
						values={[
							{
								name: 'Threat',
								value: metric.avgCastThreat,
								percentage: 100,
							},
						]}
					/>;
				},
			},
			{
				name: 'Hits',
				tooltip: 'Hits',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					if (!metric.landedHits) return '-';

					cellElem.appendChild(<>{formatToNumber(metric.landedHits)}</>);

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						spellSchool={metric.spellSchool}
						total={metric.landedHits}
						totalPercentage={100}
						values={[
							...(metric.spellType === SpellType.SpellTypeAll
								? [
										{
											name: 'Hit',
											value: metric.hits,
											percentage: metric.hitPercent,
										},
										{
											name: `Critical Hit`,
											value: metric.crits,
											percentage: metric.critPercent,
										},
								  ]
								: []),
							...(metric.spellType === SpellType.SpellTypeCast
								? [
										{
											name: 'Hit',
											value: metric.hits,
											percentage: metric.hitPercent,
										},
										{
											name: `Critical Hit`,
											value: metric.crits,
											percentage: metric.critPercent,
										},
								  ]
								: []),
							...(metric.spellType === SpellType.SpellTypePeriodic
								? [
										{
											name: 'Tick',
											value: metric.hits,
											percentage: metric.hitPercent,
										},
										{
											name: `Critical Tick`,
											value: metric.crits,
											percentage: metric.critPercent,
										},
								  ]
								: []),
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
					/>;
				},
			},
			{
				name: 'Avg Hit',
				tooltip: 'Damage / Hits',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToCompactNumber(metric.avgHit)}</>);

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						tooltipConfig={{
							onShow: () => {
								const hideThreatMetrics = !!document.querySelector('.hide-threat-metrics');
								if (hideThreatMetrics) return false;
							},
						}}
						headerValues={[, 'Amount']}
						spellSchool={metric.spellSchool}
						total={metric.avgHitThreat}
						totalPercentage={100}
						hasFooter={false}
						values={[
							{
								name: 'Threat',
								value: metric.avgHitThreat,
								percentage: 100,
							},
						]}
					/>;
				},
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
					cellElem.appendChild(<>{metric.totalMissesPercent ? formatToPercent(metric.totalMissesPercent) : '-'}</>);
					if (!metric.totalMissesPercent) return;

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
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
					/>;
				},
			},
			{
				name: 'DPS',
				tooltip: 'Damage / Encounter Duration',
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.dps,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.dps, { minimumFractionDigits: 2 })}</>);

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						tooltipConfig={{
							onShow: () => {
								const hideThreatMetrics = !!document.querySelector('.hide-threat-metrics');
								if (hideThreatMetrics) return false;
							},
						}}
						headerValues={[, 'Amount']}
						spellSchool={metric.spellSchool}
						total={metric.tps}
						totalPercentage={100}
						hasFooter={false}
						values={[
							{
								name: 'Threat',
								value: metric.tps,
								percentage: 100,
							},
						]}
					/>;
				},
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
		const isCastSpellType = metrics.some(m => m.spellType === SpellType.SpellTypeCast);
		const isDotSpellType = metrics.some(m => m.spellType === SpellType.SpellTypePeriodic);

		return ActionMetrics.merge(metrics, {
			removeTag: true,
			actionIdOverride: metrics[0]?.unit?.petActionId || undefined,
			spellTypeOverride: isCastSpellType && isDotSpellType ? SpellType.SpellTypeAll : undefined,
		});
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}