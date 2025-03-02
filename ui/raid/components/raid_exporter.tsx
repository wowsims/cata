import { Exporter, ExporterOptions } from '../../core/components/exporter';
import { RaidSimUI } from '../raid_sim_ui';

export abstract class RaidExporter extends Exporter {
	protected readonly simUI: RaidSimUI;

	constructor(parent: HTMLElement, simUI: RaidSimUI, options: ExporterOptions) {
		super(parent, options);

		this.simUI = simUI;
	}
}
