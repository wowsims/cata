import { TOOLTIP_METRIC_LABELS } from '../../constants/tooltips';
import { SpellType } from '../../proto/api';
import { ActionMetrics } from '../../proto_utils/sim_result';
import { formatToCompactNumber, formatToNumber, formatToPercent } from '../../utils';
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
				name: 'Damage Taken',
				headerCellClass: 'text-center',
				getValue: (metric: ActionMetrics) => metric.avgDamage,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
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
					const glanceValues = metric.damageDone.glance;
					const blockValues = metric.damageDone.block;
					const critBlockValues = metric.damageDone.critBlock;

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
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
							{
								name: 'Glancing Blow',
								...glanceValues,
							},
							{
								name: 'Blocked Hit',
								...blockValues,
							},
							{
								name: 'Blocked Critical Hit',
								...critBlockValues,
							},
						]}
					/>;
				},
			},
			{
				name: 'Casts',
				getValue: (metric: ActionMetrics) => metric.casts,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.casts, { fallbackString: '-' })}</>);
					if (!metric.casts) return;

					const relativeHitPercent = (metric.landedHits / (metric.landedHits + metric.totalMisses)) * 100;

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						spellSchool={metric.spellSchool}
						total={metric.casts}
						totalPercentage={100}
						hasFooter={false}
						values={[
							{
								name: 'Hits',
								value: metric.landedHits,
								percentage: relativeHitPercent,
							},
							{
								name: `Misses`,
								value: metric.totalMisses,
								percentage: metric.totalMissesPercent,
							},
						]}
					/>;
				},
			},
			{
				name: 'Avg Cast',
				tooltip: TOOLTIP_METRIC_LABELS['Damage Avg Cast'],
				getValue: (metric: ActionMetrics) => metric.avgCast,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgCast, { fallbackString: '-' }),
			},
			{
				name: 'Hits',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.landedHits, { fallbackString: '-' })}</>);
					if (!metric.landedHits) return;

					const relativeHitPercent = (metric.hits / metric.landedHits) * 100;
					const relativeCritPercent = (metric.crits / metric.landedHits) * 100;
					const relativeGlancePercent = (metric.glances / metric.landedHits) * 100;
					const relativeBlockPercent = (metric.blocks / metric.landedHits) * 100;
					const relativeCritBlockPercent = (metric.critBlocks / metric.landedHits) * 100;

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						spellSchool={metric.spellSchool}
						total={metric.landedHits}
						totalPercentage={100}
						hasFooter={false}
						values={[
							...(metric.spellType === SpellType.SpellTypeAll
								? [
										{
											name: 'Hit',
											value: metric.hits,
											percentage: relativeHitPercent,
										},
										{
											name: `Critical Hit`,
											value: metric.crits,
											percentage: relativeCritPercent,
										},
								  ]
								: []),
							...(metric.spellType === SpellType.SpellTypeCast
								? [
										{
											name: 'Hit',
											value: metric.hits,
											percentage: relativeHitPercent,
										},
										{
											name: `Critical Hit`,
											value: metric.crits,
											percentage: relativeCritPercent,
										},
								  ]
								: []),
							...(metric.spellType === SpellType.SpellTypePeriodic
								? [
										{
											name: 'Tick',
											value: metric.hits,
											percentage: relativeHitPercent,
										},
										{
											name: `Critical Tick`,
											value: metric.crits,
											percentage: relativeCritPercent,
										},
								  ]
								: []),
							{
								name: 'Glancing Blow',
								value: metric.glances,
								percentage: relativeGlancePercent,
							},
							{
								name: 'Blocked Hit',
								value: metric.blocks,
								percentage: relativeBlockPercent,
							},
							{
								name: 'Blocked Critical Hit',
								value: metric.critBlocks,
								percentage: relativeCritBlockPercent,
							},
						]}
					/>;
				},
			},
			{
				name: 'Avg Hit',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.avgHit, { fallbackString: '-' }),
			},
			{
				name: 'Miss %',
				tooltip: TOOLTIP_METRIC_LABELS['Hit Miss %'],
				getValue: (metric: ActionMetrics) => metric.totalMissesPercent,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToPercent(metric.totalMissesPercent, { fallbackString: '-' })}</>);
					if (!metric.totalMissesPercent) return;

					<MetricsCombinedTooltipTable
						tooltipElement={cellElem}
						spellSchool={metric.spellSchool}
						total={metric.totalMisses}
						totalPercentage={metric.totalMissesPercent}
						hasFooter={false}
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
				name: 'Crit %',
				getValue: (metric: ActionMetrics) => metric.critPercent,
				getDisplayString: (metric: ActionMetrics) => formatToPercent(metric.critPercent, { fallbackString: '-' }),
			},
			{
				name: 'DTPS',
				sort: ColumnSortType.Descending,
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				getValue: (metric: ActionMetrics) => metric.dps,
				getDisplayString: (metric: ActionMetrics) => formatToNumber(metric.dps, { minimumFractionDigits: 2, fallbackString: '-' }),
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
