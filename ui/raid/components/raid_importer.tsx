import { Importer, ImporterOptions } from '../../core/components/importer';
import { RaidSimUI } from '../raid_sim_ui';

export abstract class RaidImporter extends Importer {
	protected readonly simUI: RaidSimUI;

	constructor(parent: HTMLElement, simUI: RaidSimUI, options: ImporterOptions) {
		super(parent, options);

		this.simUI = simUI;
	}
}
