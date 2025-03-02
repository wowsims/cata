import { SimSettingCategories } from '../../../constants/sim_settings';
import { IndividualSimUI } from '../../../individual_sim_ui';
import { Class, EquipmentSpec, Glyphs, Profession, Race, Spec } from '../../../proto/common';
import { Database } from '../../../proto_utils/database';
import { classNames } from '../../../proto_utils/names';
import { TypedEvent } from '../../../typed_event';
import { getEnumValues } from '../../../utils';
import { Importer, ImporterOptions } from '../../importer';
import Toast from '../../toast';

// For now this just holds static helpers to match the exporter, so it doesn't extend Importer.
export abstract class IndividualImporter<SpecType extends Spec> extends Importer {
	// Exclude UISettings by default, since most users don't intend to export those.
	static readonly DEFAULT_CATEGORIES = getEnumValues(SimSettingCategories).filter(c => c != SimSettingCategories.UISettings) as Array<SimSettingCategories>;
	static readonly CATEGORY_PARAM = 'i';

	protected readonly simUI: IndividualSimUI<any>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>, options: ImporterOptions) {
		super(parent, options);
		this.simUI = simUI;
	}

	protected async finishIndividualImport<SpecType extends Spec>(
		simUI: IndividualSimUI<SpecType>,
		charClass: Class,
		race: Race,
		equipmentSpec: EquipmentSpec,
		talentsStr: string,
		glyphs: Glyphs | null,
		professions: Array<Profession>,
	): Promise<void> {
		if (charClass != simUI.player.getClass()) {
			throw new Error(`Wrong Class! Expected ${simUI.player.getPlayerClass().friendlyName} but found ${classNames.get(charClass)}!`);
		}

		await Database.loadLeftoversIfNecessary(equipmentSpec);

		const gear = simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		const expectedEnchantIds = equipmentSpec.items.map(item => item.enchant);
		const foundEnchantIds = gear.asSpec().items.map(item => item.enchant);
		const missingEnchants = expectedEnchantIds.filter(expectedId => expectedId != 0 && !foundEnchantIds.includes(expectedId));

		const expectedItemIds = equipmentSpec.items.map(item => item.id);
		const foundItemIds = gear.asSpec().items.map(item => item.id);
		const missingItems = expectedItemIds.filter(expectedId => !foundItemIds.includes(expectedId));

		// Now update settings using the parsed values.
		const eventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			simUI.player.setRace(eventID, race);
			simUI.player.setGear(eventID, gear);
			if (talentsStr && talentsStr != '--') {
				simUI.player.setTalentsString(eventID, talentsStr);
			}
			if (glyphs) {
				simUI.player.setGlyphs(eventID, glyphs);
			}
			if (professions.length > 0) {
				simUI.player.setProfessions(eventID, professions);
			}
		});

		this.close();

		if (missingItems.length == 0 && missingEnchants.length == 0) {
			new Toast({ variant: 'success', body: `Import successful!` });
		} else {
			new Toast({
				variant: 'info',
				body:
					'Import successful, but the following IDs were not found in the sim database:' +
					(missingItems.length == 0 ? '' : '\n\nItems: ' + missingItems.join(', ')) +
					(missingEnchants.length == 0 ? '' : '\n\nEnchants: ' + missingEnchants.join(', ')),
			});
		}
	}
}
