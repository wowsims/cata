import { RaidSimSettings } from '../../../core/proto/ui';
import { TypedEvent } from '../../../core/typed_event';
import { RaidSimUI } from '../../raid_sim_ui';
import { RaidImporter } from '../raid_importer';

export class RaidJsonImporter extends RaidImporter {
	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, simUI, { title: 'JSON Import', allowFileUpload: true });

		this.descriptionElem.appendChild(
			<>
				<p>Import settings from a JSON text file, which can be created using the JSON Export feature of this site.</p>
				<p>To import, paste the JSON text below and click, 'Import'.</p>
			</>,
		);
	}

	async onImport(data: string) {
		const settings = RaidSimSettings.fromJsonString(data, { ignoreUnknownFields: true });
		this.simUI.fromProto(TypedEvent.nextEventID(), settings);
		this.close();
	}
}
