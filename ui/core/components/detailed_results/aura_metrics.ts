import { OtherAction } from '../../proto/common';
import { AuraMetrics, SimResult, SimResultFilter } from '../../proto_utils/sim_result';
import { ColumnSortType, MetricsTable } from './metrics_table/metrics_table';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component';

export class AuraMetricsTable extends MetricsTable<AuraMetrics> {
	private readonly useDebuffs: boolean;

	constructor(config: ResultComponentConfig, useDebuffs: boolean) {
		if (useDebuffs) {
			config.rootCssClass = 'debuff-metrics-root';
		} else {
			config.rootCssClass = 'buff-metrics-root';
		}
		super(config, [
			MetricsTable.nameCellConfig((metric: AuraMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
					metricType: metric?.constructor?.name,
				};
			}),
			{
				name: 'Procs',
				tooltip: 'Procs',
				getValue: (metric: AuraMetrics) => metric.averageProcs,
				getDisplayString: (metric: AuraMetrics) => metric.averageProcs.toFixed(2),
			},
			{
				name: 'PPM',
				tooltip: 'Procs Per Minute',
				getValue: (metric: AuraMetrics) => metric.ppm,
				getDisplayString: (metric: AuraMetrics) => metric.ppm.toFixed(2),
			},
			{
				name: 'Uptime',
				tooltip: 'Uptime / Encounter Duration',
				sort: ColumnSortType.Descending,
				getValue: (metric: AuraMetrics) => metric.uptimePercent,
				getDisplayString: (metric: AuraMetrics) => metric.uptimePercent.toFixed(2) + '%',
			},
		]);
		this.useDebuffs = useDebuffs;
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<AuraMetrics>> {
		if (this.useDebuffs) {
			return AuraMetrics.groupById(resultData.result.getDebuffMetrics(resultData.filter));
		} else {
			const players = resultData.result.getRaidIndexedPlayers(resultData.filter);
			if (players.length != 1) {
				return [];
			}
			const player = players[0];

			const auras = this.filterMetrics(player.auras);
			const actionGroups = AuraMetrics.groupById(auras);
			const petGroups = player.pets.map(pet => this.filterMetrics(pet.auras));
			return actionGroups.concat(petGroups);
		}
	}

	mergeMetrics(metrics: Array<AuraMetrics>): AuraMetrics {
		return AuraMetrics.merge(metrics, {
			removeTag: true,
			actionIdOverride: metrics[0].unit?.petActionId || undefined,
		});
	}

	shouldCollapse(metric: AuraMetrics): boolean {
		return !metric.unit?.isPet;
	}

	filterMetrics(metrics: Array<AuraMetrics>): Array<AuraMetrics> {
		return metrics.filter(aura => aura.actionId.otherId !== OtherAction.OtherActionMove);
	}
}
