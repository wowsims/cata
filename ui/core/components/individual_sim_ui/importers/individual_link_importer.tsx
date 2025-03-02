import pako from 'pako';

import { SIM_CATEGORY_KEYS, SimSettingCategories } from '../../../constants/sim_settings';
import { IndividualSimSettings } from '../../../proto/ui';
import { IndividualImporter } from './individual_importer';

interface UrlParseData {
	settings: IndividualSimSettings;
	categories: Array<SimSettingCategories>;
}

// For now this just holds static helpers to match the exporter, so it doesn't extend Importer.
export class IndividualLinkImporter {
	static tryParseUrlLocation(location: Location | URL): UrlParseData | null {
		let hash = location.hash;
		if (hash.length <= 1) {
			return null;
		}

		// Remove leading '#'
		hash = hash.substring(1);
		const binary = atob(hash);
		const bytes = new Uint8Array(binary.length);
		for (let i = 0; i < bytes.length; i++) {
			bytes[i] = binary.charCodeAt(i);
		}

		const settingsBytes = pako.inflate(bytes);
		const settings = IndividualSimSettings.fromBinary(settingsBytes);

		let exportCategories = IndividualImporter.DEFAULT_CATEGORIES;
		const urlParams = new URLSearchParams(window.location.search);
		if (urlParams.has(IndividualImporter.CATEGORY_PARAM)) {
			const categoryChars = urlParams.get(IndividualImporter.CATEGORY_PARAM)!.split('');
			exportCategories = categoryChars
				.map(char => [...SIM_CATEGORY_KEYS.entries()].find(e => e[1] == char))
				.filter(e => e)
				.map(e => e![0]);
		}

		return {
			settings: settings,
			categories: exportCategories,
		};
	}
}
