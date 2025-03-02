import { IndividualSimUI } from '../../../individual_sim_ui';
import { Spec } from '../../../proto/common';
import { IndividualSimSettings } from '../../../proto/ui';
import { Database } from '../../../proto_utils/database';
import { TypedEvent } from '../../../typed_event';
import { IndividualImporter } from './individual_importer';

export class IndividualJsonImporter<SpecType extends Spec> extends IndividualImporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'JSON Import', allowFileUpload: true });

		this.descriptionElem.appendChild(
			<>
				<p>Import settings from a JSON file, which can be created using the JSON Export feature.</p>
				<p>To import, upload the file or paste the text below, then click, 'Import'.</p>
			</>,
		);
	}

	async onImport(data: string) {
		let proto: ReturnType<typeof IndividualSimSettings.fromJsonString> | null = null;
		try {
			proto = IndividualSimSettings.fromJsonString(data, { ignoreUnknownFields: true });
		} catch {
			throw new Error('Please use a valid JSON object.');
		}
		if (proto.player?.equipment) {
			await Database.loadLeftoversIfNecessary(proto.player.equipment);
		}
		if (this.simUI.isWithinRaidSim) {
			if (proto.player) {
				this.simUI.player.fromProto(TypedEvent.nextEventID(), proto.player);
			}
		} else {
			this.simUI.fromProto(TypedEvent.nextEventID(), proto);
		}
		this.close();
	}
}
