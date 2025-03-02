import { SimSettingCategories } from '../../../constants/sim_settings';
import { IndividualSimUI } from '../../../individual_sim_ui';
import { Spec } from '../../../proto/common';
import { IndividualSimSettings } from '../../../proto/ui';
import { getEnumValues, jsonStringifyWithFlattenedPaths } from '../../../utils';
import { IndividualImporter } from '../importers/individual_importer';
import { IndividualExporter } from './individual_exporter';

export class IndividualJsonExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'JSON Export', allowDownload: true, selectCategories: true });
	}

	getData(): string {
		return IndividualJsonExporter.getData(
			this.simUI,
			(getEnumValues(SimSettingCategories) as Array<SimSettingCategories>).filter(c => this.exportCategories[c]),
		);
	}

	static getData(simUI: IndividualSimUI<any>, exportCategories?: Array<SimSettingCategories>): string {
		if (!exportCategories) {
			exportCategories = IndividualImporter.DEFAULT_CATEGORIES;
		}

		return jsonStringifyWithFlattenedPaths(IndividualSimSettings.toJson(simUI.toProto(exportCategories)), 2, (value, path) => {
			if (['stats', 'pseudoStats'].includes(path[path.length - 1])) {
				return true;
			}

			if (['player', 'equipment', 'items'].every((v, i) => path[i] == v)) {
				return path.length > 3;
			}

			if (path[0] == 'player' && path[1] == 'rotation' && ['prepullActions', 'priorityList'].includes(path[2])) {
				return path.length > 3;
			}

			return false;
		});
	}
}
