import { IndividualSimUI } from '../../../individual_sim_ui';
import { Class, EquipmentSpec, ItemSpec, Race, Spec } from '../../../proto/common';
import { nameToClass, nameToRace } from '../../../proto_utils/names';
import { talentSpellIdsToTalentString } from '../../../talents/factory';
import Toast from '../../toast';
import { IndividualImporter } from './individual_importer';

export class Individual60UImporter<SpecType extends Spec> extends IndividualImporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Sixty Upgrades Cataclysm Import', allowFileUpload: true });

		this.descriptionElem.appendChild(
			<>
				<p>
					Import settings from{' '}
					<a href="https://sixtyupgrades.com/cata" target="_blank">
						Sixty Upgrades
					</a>
					.
				</p>
				<p>This feature imports gear, race, and (optionally) talents. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.</p>
				<p>To import, paste the output from the site's export option below and click, 'Import'.</p>
			</>,
		);
	}

	async onImport(data: string) {
		let importJson: any | null;
		try {
			importJson = JSON.parse(data);
		} catch {
			throw new Error('Please use a valid Sixty Upgrades export.');
		}

		// Parse all the settings.
		const charClass = nameToClass((importJson?.character?.gameClass as string) || '');
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class!');
		}

		const race = nameToRace((importJson?.character?.race as string) || '');
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race!');
		}

		let talentsStr = '';
		if (importJson?.talents?.length > 0) {
			const talentIds = (importJson.talents as Array<any>).map(talentJson => talentJson.spellId);
			talentsStr = talentSpellIdsToTalentString(charClass, talentIds);
		}

		let hasRemovedRandomSuffix = false;
		const modifiedItemNames: string[] = [];
		const equipmentSpec = EquipmentSpec.create();
		(importJson.items as Array<any>).forEach(itemJson => {
			const itemSpec = ItemSpec.create();
			itemSpec.id = itemJson.id;
			if (itemJson.enchant?.id) {
				itemSpec.enchant = itemJson.enchant.id;
			}
			if (itemJson.gems) {
				itemSpec.gems = (itemJson.gems as Array<any>).filter(gemJson => gemJson?.id).map(gemJson => gemJson.id);
			}

			// As long as 60U exports the wrong suffixes we should
			// inform the user that they need to manually add them.
			// Due to this we also remove the reforge on the item.
			if (itemJson.suffixId) {
				hasRemovedRandomSuffix = true;
				if (itemJson.reforge?.id) {
					itemJson.reforge.id = null;
				}
				modifiedItemNames.push(itemJson.name);
			}
			if (itemJson.reforge?.id) {
				itemSpec.reforging = itemJson.reforge.id;
			}
			equipmentSpec.items.push(itemSpec);
		});

		this.simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, null, []);

		if (hasRemovedRandomSuffix && modifiedItemNames.length) {
			new Toast({
				variant: 'warning',
				body: (
					<>
						<p>Sixty Upgrades currently exports the wrong Random Suffixes. We have removed the random suffix on the following item(s):</p>
						<ul>
							{modifiedItemNames.map(itemName => (
								<li>
									<strong>{itemName}</strong>
								</li>
							))}
						</ul>
					</>
				),
				delay: 8000,
			});
		}
	}
}
