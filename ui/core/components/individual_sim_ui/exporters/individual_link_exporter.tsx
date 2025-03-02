import { default as pako } from 'pako';

import { SIM_CATEGORY_KEYS, SimSettingCategories } from '../../../constants/sim_settings';
import { IndividualSimUI } from '../../../individual_sim_ui';
import { Spec } from '../../../proto/common';
import { IndividualSimSettings } from '../../../proto/ui';
import { arrayEquals, getEnumValues } from '../../../utils';
import { IndividualImporter } from '../importers/individual_importer';
import { IndividualExporter } from './individual_exporter';

export class IndividualLinkExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Shareable Link', selectCategories: true });
	}

	getData(): string {
		return IndividualLinkExporter.createLink(
			this.simUI,
			(getEnumValues(SimSettingCategories) as Array<SimSettingCategories>).filter(c => this.exportCategories[c]),
		);
	}

	static createLink(simUI: IndividualSimUI<any>, exportCategories?: Array<SimSettingCategories>): string {
		if (!exportCategories) {
			exportCategories = IndividualImporter.DEFAULT_CATEGORIES;
		}

		const proto = simUI.toProto(exportCategories);

		const protoBytes = IndividualSimSettings.toBinary(proto);
		// @ts-ignore Pako did some weird stuff between versions and the @types package doesn't correctly support this syntax for version 2.0.4 but it's completely valid
		// The syntax was removed in 2.1.0 and there were several complaints but the project seems to be largely abandoned now
		const deflated = pako.deflate(protoBytes, { to: 'string' });
		const encoded = btoa(String.fromCharCode(...deflated));

		const linkUrl = new URL(window.location.href);
		linkUrl.hash = encoded;
		if (arrayEquals(exportCategories, IndividualImporter.DEFAULT_CATEGORIES)) {
			linkUrl.searchParams.delete(IndividualImporter.CATEGORY_PARAM);
		} else {
			const categoryCharString = exportCategories.map(c => SIM_CATEGORY_KEYS.get(c)).join('');
			linkUrl.searchParams.set(IndividualImporter.CATEGORY_PARAM, categoryCharString);
		}
		return linkUrl.toString();
	}
}
