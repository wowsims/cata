import { TOOLTIP_METRIC_LABELS } from '../../constants/tooltips';
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
				headerCellClass: 'text-center metrics-table-cell--primary-metric',
				columnClass: 'metrics-table-cell--primary-metric',
				getValue: (metric: ActionMetrics) => metric.avgDamage,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
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
					const critHitValues = metric.damageDone.critHit;
					const tickValues = metric.damageDone.tick;
					const critTickValues = metric.damageDone.critTick;
					const glanceValues = metric.damageDone.glance;
					const glanceBlockValues = metric.damageDone.glanceBlock;
					const blockValues = metric.damageDone.block;
					const critBlockValues = metric.damageDone.critBlock;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							headerValues={[, 'Amount']}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.damage,
									totalPercentage: 100,
									data: [
										{
											name: 'Hit',
											...hitValues,
										},
										{
											name: `Critical Hit`,
											...critHitValues,
										},
										{
											name: 'Tick',
											...tickValues,
										},
										{
											name: `Critical Tick`,
											...critTickValues,
										},
										{
											name: 'Glancing Blow',
											...glanceValues,
										},
										{
											name: 'Blocked Glancing Blow',
											...glanceBlockValues,
										},
										{
											name: 'Blocked Hit',
											...blockValues,
										},
										{
											name: 'Blocked Critical Hit',
											...critBlockValues,
										},
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'Casts',
				getValue: (metric: ActionMetrics) => metric.casts,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.casts, { fallbackString: '-' })}</>);

					if ((!metric.landedHits && !metric.totalMisses) || metric.isPassiveAction) return;
					const relativeHitPercent = ((metric.landedHits || metric.casts) / ((metric.landedHits || metric.casts) + metric.totalMisses)) * 100;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.casts,
									totalPercentage: 100,
									data: [
										{
											name: 'Hits',
											value: metric.landedHits || metric.casts - metric.totalMisses,
											percentage: relativeHitPercent,
										},
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
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'Avg Cast',
				tooltip: TOOLTIP_METRIC_LABELS['Damage Avg Cast'],
				getValue: (metric: ActionMetrics) => {
					if (metric.isPassiveAction) return 0;
					return metric.avgCastHit || metric.avgCastTick;
				},
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<>
							{metric.isPassiveAction ? (
								'-'
							) : (
								<>
									{formatToCompactNumber(metric.avgCastHit || metric.avgCastTick, { fallbackString: '-' })}
									{metric.avgCastHit && metric.avgCastTick ? (
										<> ({formatToCompactNumber(metric.avgCastTick, { fallbackString: '-' })})</>
									) : undefined}
								</>
							)}
						</>,
					);

					if (!metric.avgCastHit && !metric.avgCastTick) return;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							tooltipConfig={{
								onShow: () => {
									const hideThreatMetrics = !!document.querySelector('.hide-threat-metrics');
									if (hideThreatMetrics) return false;
								},
							}}
							headerValues={[, 'Amount']}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.avgCastThreat,
									totalPercentage: 100,
									data: [
										{
											name: 'Threat',
											value: metric.avgCastThreat,
											percentage: 100,
										},
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'Hits',
				getValue: (metric: ActionMetrics) => metric.landedHits || metric.landedTicks,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<>
							{formatToNumber(metric.landedHits || metric.landedTicks, { fallbackString: '-' })}
							{metric.landedHits && metric.landedTicks ? <> ({formatToNumber(metric.landedTicks, { fallbackString: '-' })})</> : undefined}
						</>,
					);
					if (!metric.landedHits && !metric.landedTicks) return;

					const relativeHitPercent = (metric.hits / metric.landedHits) * 100;
					const relativeCritPercent = (metric.crits / metric.landedHits) * 100;
					const relativeTickPercent = (metric.ticks / metric.landedTicks) * 100;
					const relativeCritTickPercent = (metric.critTicks / metric.landedTicks) * 100;
					const relativeGlancePercent = (metric.glances / metric.landedHits) * 100;
					const relativeGlanceBlockPercent = (metric.glanceBlocks / metric.landedHits) * 100;
					const relativeBlockPercent = (metric.blocks / metric.landedHits) * 100;
					const relativeCritBlockPercent = (metric.critBlocks / metric.landedHits) * 100;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.landedHits,
									totalPercentage: 100,
									data: [
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
										{
											name: 'Glancing Blow',
											value: metric.glances,
											percentage: relativeGlancePercent,
										},
										{
											name: 'Blocked Glancing Blow',
											value: metric.glanceBlocks,
											percentage: relativeGlanceBlockPercent,
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
									],
								},
								{
									spellSchool: metric.spellSchool,
									total: metric.landedTicks,
									totalPercentage: 100,
									data: [
										{
											name: 'Tick',
											value: metric.ticks,
											percentage: relativeTickPercent,
										},
										{
											name: `Critical Tick`,
											value: metric.critTicks,
											percentage: relativeCritTickPercent,
										},
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'Avg Hit',
				getValue: (metric: ActionMetrics) => metric.avgHit || metric.avgTick,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(
						<>
							{formatToCompactNumber(metric.avgHit || metric.avgTick, { fallbackString: '-' })}
							{metric.avgHit && metric.avgTick ? <> ({formatToCompactNumber(metric.avgTick, { fallbackString: '-' })})</> : undefined}
						</>,
					);

					if (!metric.avgHitThreat) return;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							tooltipConfig={{
								onShow: () => {
									const hideThreatMetrics = !!document.querySelector('.hide-threat-metrics');
									if (hideThreatMetrics) return false;
								},
							}}
							headerValues={[, 'Amount']}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.avgHitThreat,
									totalPercentage: 100,
									data: [
										{
											name: 'Threat',
											value: metric.avgHitThreat,
											percentage: 100,
										},
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'Crit %',
				getValue: (metric: ActionMetrics) => (metric.critPercent + metric.critBlockPercent) || metric.critTickPercent,
				getDisplayString: (metric: ActionMetrics) =>
					`${formatToPercent(metric.critPercent + metric.critBlockPercent || metric.critTickPercent, { fallbackString: '-' })}${
						metric.critPercent + metric.critBlockPercent && metric.critTickPercent ? ` (${formatToPercent(metric.critTickPercent, { fallbackString: '-' })})` : ''
					}`,
			},
			{
				name: 'Miss %',
				getValue: (metric: ActionMetrics) => metric.totalMissesPercent,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToPercent(metric.totalMissesPercent, { fallbackString: '-' })}</>);
					if (!metric.totalMissesPercent) return;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.totalMisses,
									totalPercentage: metric.totalMissesPercent,
									data: [
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
									],
								},
							]}
						/>,
					);
				},
			},
			{
				name: 'DPET',
				getValue: (metric: ActionMetrics) => metric.damageThroughput,
				getDisplayString: (metric: ActionMetrics) => formatToCompactNumber(metric.damageThroughput, { fallbackString: '-' }),
			},
			{
				name: 'DPS',
				headerCellClass: 'text-body',
				columnClass: 'text-success',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.dps,
				fillCell: (metric: ActionMetrics, cellElem: HTMLElement) => {
					cellElem.appendChild(<>{formatToNumber(metric.dps, { minimumFractionDigits: 2, fallbackString: '-' })}</>);
					if (!metric.dps) return;

					cellElem.appendChild(
						<MetricsCombinedTooltipTable
							tooltipElement={cellElem}
							tooltipConfig={{
								onShow: () => {
									const hideThreatMetrics = !!document.querySelector('.hide-threat-metrics');
									if (hideThreatMetrics) return false;
								},
							}}
							headerValues={[, 'Amount']}
							groups={[
								{
									spellSchool: metric.spellSchool,
									total: metric.tps,
									totalPercentage: 100,
									data: [
										{
											name: 'Threat',
											value: metric.tps,
											percentage: 100,
										},
									],
								},
							]}
						/>,
					);
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
		return ActionMetrics.merge(metrics, {
			removeTag: true,
			actionIdOverride: metrics[0]?.unit?.petActionId || undefined,
		});
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
