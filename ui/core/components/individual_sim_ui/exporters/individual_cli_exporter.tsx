import { IndividualSimUI } from '../../../individual_sim_ui';
import { RaidSimRequest } from '../../../proto/api';
import { Spec } from '../../../proto/common';
import { IndividualExporter } from './individual_exporter';

export class IndividualCLIExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'CLI Export', allowDownload: true });
	}

	getData(): string {
		const raidSimJson: any = RaidSimRequest.toJson(this.simUI.sim.makeRaidSimRequest(false));
		delete raidSimJson.raid?.parties[0]?.players[0]?.database;
		return JSON.stringify(raidSimJson, null, 2);
	}
}
