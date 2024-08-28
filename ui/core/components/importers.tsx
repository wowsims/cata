import { JsonObject } from '@protobuf-ts/runtime';
import { default as pako } from 'pako';
import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../individual_sim_ui';
import { Class, EquipmentSpec, Glyphs, ItemSlot, ItemSpec, Profession, Race, Spec } from '../proto/common';
import { IndividualSimSettings } from '../proto/ui';
import { Database } from '../proto_utils/database';
import { classNames, nameToClass, nameToProfession, nameToRace } from '../proto_utils/names';
import { SimSettingCategories } from '../sim';
import { SimUI } from '../sim_ui';
import { classGlyphsConfig, talentSpellIdsToTalentString } from '../talents/factory';
import { GlyphConfig } from '../talents/glyphs_picker';
import { TypedEvent } from '../typed_event';
import { buf2hex, getEnumValues } from '../utils';
import { BaseModal } from './base_modal';
import Toast from './toast';

export abstract class Importer extends BaseModal {
	protected readonly textElem: HTMLTextAreaElement;
	protected readonly descriptionElem: HTMLElement;
	protected readonly importButton: HTMLButtonElement;
	private readonly includeFile: boolean;

	constructor(parent: HTMLElement, simUI: SimUI, title: string, includeFile: boolean) {
		super(parent, 'importer', { title: title, footer: true, disposeOnClose: false });
		this.includeFile = includeFile;
		const uploadInputId = 'upload-input-' + title.toLowerCase().replaceAll(' ', '-');

		const descriptionElemRef = ref<HTMLDivElement>();
		const textElemRef = ref<HTMLTextAreaElement>();
		const importButtonRef = ref<HTMLButtonElement>();
		const uploadInputRef = ref<HTMLInputElement>();

		this.body.replaceChildren(
			<>
				<div ref={descriptionElemRef} className="import-description"></div>
				<textarea ref={textElemRef} className="importer-textarea form-control" attributes={{ spellcheck: false }}></textarea>
			</>,
		);

		this.footer!.appendChild(
			<>
				{this.includeFile && (
					<>
						<label htmlFor={uploadInputId} className="importer-button btn btn-primary upload-button me-2">
							<i className="fas fa-file-arrow-up me-1"></i>
							Upload File
						</label>
						<input ref={uploadInputRef} type="file" id={uploadInputId} className="importer-upload-input d-none" hidden />
					</>
				)}
				<button ref={importButtonRef} className="importer-button btn btn-primary import-button">
					<i className="fa fa-download me-1"></i>
					Import
				</button>
			</>,
		);

		this.descriptionElem = descriptionElemRef.value!;
		this.textElem = textElemRef.value!;

		if (this.includeFile && uploadInputRef.value) {
			uploadInputRef.value.addEventListener('change', async event => {
				this.textElem.textContent = await (event as any).target.files[0].text();
			});
		}

		this.importButton = importButtonRef.value!;
		this.importButton.addEventListener('click', async _event => {
			try {
				await this.onImport(this.textElem.value || '');
			} catch (error: any) {
				new Toast({ variant: 'error', body: `Import error: ${error?.message || error}` });
			}
		});
	}

	abstract onImport(data: string): Promise<void>;

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

interface UrlParseData {
	settings: IndividualSimSettings;
	categories: Array<SimSettingCategories>;
}

// For now this just holds static helpers to match the exporter, so it doesn't extend Importer.
export class IndividualLinkImporter {
	// Exclude UISettings by default, since most users don't intend to export those.
	static readonly DEFAULT_CATEGORIES = getEnumValues(SimSettingCategories).filter(c => c != SimSettingCategories.UISettings) as Array<SimSettingCategories>;

	static readonly CATEGORY_PARAM = 'i';
	static readonly CATEGORY_KEYS: Map<SimSettingCategories, string> = (() => {
		const map = new Map();
		// Use single-letter abbreviations since these will be included in sim links.
		map.set(SimSettingCategories.Gear, 'g');
		map.set(SimSettingCategories.Talents, 't');
		map.set(SimSettingCategories.Rotation, 'r');
		map.set(SimSettingCategories.Consumes, 'c');
		map.set(SimSettingCategories.Miscellaneous, 'm');
		map.set(SimSettingCategories.External, 'x');
		map.set(SimSettingCategories.Encounter, 'e');
		map.set(SimSettingCategories.UISettings, 'u');
		return map;
	})();

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

		let exportCategories = IndividualLinkImporter.DEFAULT_CATEGORIES;
		const urlParams = new URLSearchParams(window.location.search);
		if (urlParams.has(IndividualLinkImporter.CATEGORY_PARAM)) {
			const categoryChars = urlParams.get(IndividualLinkImporter.CATEGORY_PARAM)!.split('');
			exportCategories = categoryChars
				.map(char => [...IndividualLinkImporter.CATEGORY_KEYS.entries()].find(e => e[1] == char))
				.filter(e => e)
				.map(e => e![0]);
		}

		return {
			settings: settings,
			categories: exportCategories,
		};
	}
}

export class IndividualJsonImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'JSON Import', true);
		this.simUI = simUI;

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

export class Individual60UImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, '60 Upgrades Cataclysm Import', true);
		this.simUI = simUI;

		this.descriptionElem.appendChild(
			<>
				<p>
					Import settings from{' '}
					<a href="https://sixtyupgrades.com/cata" target="_blank">
						60 Upgrades
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
			throw new Error('Please use a valid 60U export.');
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
						<p>60U currently exports the wrong Random Suffixes. We have removed the random suffix on the following item(s):</p>
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

export class IndividualWowheadGearPlannerImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'Wowhead Import', true);
		this.simUI = simUI;

		this.descriptionElem.appendChild(
			<>
				<p>
					Import settings from{' '}
					<a href="https://www.wowhead.com/cata/gear-planner" target="_blank">
						Wowhead Gear Planner
					</a>
					.
				</p>
				<p>This feature imports gear, race, and (optionally) talents. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.</p>
				<p>To import, paste the gear planner link below and click, 'Import'.</p>
			</>,
		);
	}

	async onImport(url: string) {
		const match = url.match(/www\.wowhead\.com\/cata\/gear-planner\/([a-z\-]+)\/([a-z\-]+)\/([a-zA-Z0-9_\-]+)/);
		if (!match) {
			throw new Error(`Invalid WCL URL ${url}, must look like "https://www.wowhead.com/cata/gear-planner/CLASS/RACE/XXXX"`);
		}

		// Parse all the settings.
		const charClass = nameToClass(match[1].replaceAll('-', ''));
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class: ' + match[1]);
		}

		const race = nameToRace(match[2].replaceAll('-', ''));
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race: ' + match[2]);
		}

		const base64Data = match[3].replaceAll('_', '/').replaceAll('-', '+');
		//console.log('Base64: ' + base64Data);
		const data = Uint8Array.from(atob(base64Data), c => c.charCodeAt(0));
		//console.log('Hex: ' + buf2hex(data));

		// Binary schema
		// Byte 00: ??
		// Byte 01: ?? Seems related to aesthetics (e.g. body type)
		// Byte 02: 8-bit Player Level
		// Byte 03: 8-bit length of talents bytes
		// Next N Bytes: Talents in hex string format

		// Talent hex string looks like '230005232100330150323102505321f03f023203001f'
		// Just like regular wowhead talents string except 'f' instead of '-'.
		const numTalentBytes = data[3];
		const talentBytes = data.subarray(4, 4 + numTalentBytes);
		const talentsHexStr = buf2hex(talentBytes);
		//console.log('Talents hex: ' + talentsHexStr);
		const talentsStr = talentsHexStr.split('f').slice(0, 3).join('-');
		//console.log('Talents: ' + talentsStr);

		let cur = 4 + numTalentBytes;
		const numGlyphBytes = data[cur];
		cur++;
		const glyphBytes = data.subarray(cur, cur + numGlyphBytes);
		const gearBytes = data.subarray(cur + numGlyphBytes);
		//console.log(`Glyphs have ${numGlyphBytes} bytes: ${buf2hex(glyphBytes)}`);
		//console.log(`Remaining ${gearBytes.length} bytes: ${buf2hex(gearBytes)}`);

		// First byte in glyphs section seems to always be 0x30
		cur = 1;
		let hasGlyphs = false;
		const d = '0123456789abcdefghjkmnpqrstvwxyz';
		const glyphStr = String.fromCharCode(...glyphBytes);
		const glyphIds = [0, 0, 0, 0, 0, 0];
		while (cur < glyphBytes.length) {
			// First byte for each glyph is 0x3z, where z is the glyph position.
			// 0, 1, 2 are major glyphs, 3, 4, 5 are minor glyphs.
			const glyphPosition = d.indexOf(glyphStr[cur]);
			cur++;

			// For some reason, wowhead uses the spell IDs for the glyphs and
			// applies a ridiculous hashing scheme.
			const spellId =
				0 +
				(d.indexOf(glyphStr[cur + 0]) << 15) +
				(d.indexOf(glyphStr[cur + 1]) << 10) +
				(d.indexOf(glyphStr[cur + 2]) << 5) +
				(d.indexOf(glyphStr[cur + 3]) << 0);
			const itemId = this.simUI.sim.db.glyphSpellToItemId(spellId);
			//console.log(`Glyph position: ${glyphPosition}, spellID: ${spellId}`);

			hasGlyphs = true;
			glyphIds[glyphPosition] = itemId;
			cur += 4;
		}
		const glyphs = Glyphs.create({
			major1: glyphIds[0],
			major2: glyphIds[1],
			major3: glyphIds[2],
			minor1: glyphIds[3],
			minor2: glyphIds[4],
			minor3: glyphIds[5],
		});

		// Binary schema for each item:
		// 8-bit slotNumber, high bit = is enchanted
		// 8-bit upper 3 bits for gem count
		// 16-bit item id
		// if enchant bit is set:
		//   8-bit ??, possibly enchant position for multiple enchants?
		//   16-bit enchant id
		// for each gem:
		//   8-bit upper 3 bits for gem position
		//   16-bit gem item id
		const equipmentSpec = EquipmentSpec.create();
		cur = 0;
		while (cur < gearBytes.length) {
			const itemSpec = ItemSpec.create();
			const slotId = gearBytes[cur] & 0b00111111;
			const isEnchanted = Boolean(gearBytes[cur] & 0b10000000);
			const randomEnchant = Boolean(gearBytes[cur] & 0b01000000);
			cur++;

			const numGems = (gearBytes[cur] & 0b11100000) >> 5;
			const highid = gearBytes[cur] & 0b00011111;
			cur++;

			itemSpec.id = (highid << 16) + (gearBytes[cur] << 8) + gearBytes[cur + 1];
			cur += 2;
			//console.log(`Slot ID: ${slotId}, isEnchanted: ${isEnchanted}, numGems: ${numGems}, itemID: ${itemSpec.id}`);

			if (isEnchanted) {
				// Note: this is the enchant SPELL id, not the effect ID.
				const enchantSpellId = (gearBytes[cur] << 16) + (gearBytes[cur + 1] << 8) + gearBytes[cur + 2];
				itemSpec.enchant = this.simUI.sim.db.enchantSpellIdToEffectId(enchantSpellId);
				cur += 3;
				//console.log(`Enchant ID: ${itemSpec.enchant}. Spellid: ${enchantSpellId}`);
			}

			for (let gemIdx = 0; gemIdx < numGems; gemIdx++) {
				const gemPosition = (gearBytes[cur] & 0b11100000) >> 5;
				const highgemid = gearBytes[cur] & 0b00011111;
				cur++;

				const gemId = (highgemid << 16) + (gearBytes[cur] << 8) + gearBytes[cur + 1];
				cur += 2;
				//console.log(`Gem position: ${gemPosition}, gemID: ${gemId}`);

				if (!itemSpec.gems) {
					itemSpec.gems = [];
				}
				while (itemSpec.gems.length < gemPosition) {
					itemSpec.gems.push(0);
				}
				itemSpec.gems[gemPosition] = gemId;
			}

			// Ignore tabard / shirt slots
			const itemSlotEntry = Object.entries(IndividualWowheadGearPlannerImporter.slotIDs).find(e => e[1] == slotId);
			if (itemSlotEntry != null) {
				equipmentSpec.items.push(itemSpec);
			}
		}
		const gear = this.simUI.sim.db.lookupEquipmentSpec(equipmentSpec);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, hasGlyphs ? glyphs : null, []);
	}

	static slotIDs: Record<ItemSlot, number> = {
		[ItemSlot.ItemSlotHead]: 1,
		[ItemSlot.ItemSlotNeck]: 2,
		[ItemSlot.ItemSlotShoulder]: 3,
		[ItemSlot.ItemSlotBack]: 15,
		[ItemSlot.ItemSlotChest]: 5,
		[ItemSlot.ItemSlotWrist]: 9,
		[ItemSlot.ItemSlotHands]: 10,
		[ItemSlot.ItemSlotWaist]: 6,
		[ItemSlot.ItemSlotLegs]: 7,
		[ItemSlot.ItemSlotFeet]: 8,
		[ItemSlot.ItemSlotFinger1]: 11,
		[ItemSlot.ItemSlotFinger2]: 12,
		[ItemSlot.ItemSlotTrinket1]: 13,
		[ItemSlot.ItemSlotTrinket2]: 14,
		[ItemSlot.ItemSlotMainHand]: 16,
		[ItemSlot.ItemSlotOffHand]: 17,
		[ItemSlot.ItemSlotRanged]: 18,
	};
}

export class IndividualAddonImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, 'Addon Import', true);
		this.simUI = simUI;

		const warningRef = ref<HTMLDivElement>();
		this.descriptionElem.appendChild(
			<>
				<p>
					Import settings from the{' '}
					<a href="https://www.curseforge.com/wow/addons/wowsimsexporter" target="_blank">
						WoWSims Importer In-Game Addon
					</a>
					.
				</p>
				<p>
					This feature imports gear, race, talents, glyphs, and professions. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.
				</p>
				<p>To import, paste the output from the addon below and click, 'Import'.</p>
				<div ref={warningRef} />
			</>,
		);

		if (warningRef.value)
			new Toast({
				title: 'Reforging issues',
				body: (
					<>
						There are known issues with Reforging when using the WSE addon.
						<br />
						Always make sure to double check your reforges after importing.
					</>
				),
				additionalClasses: ['toast-import-warning'],
				container: warningRef.value,
				variant: 'warning',
				canClose: false,
				autoShow: true,
				autohide: false,
			});
	}

	async onImport(data: string) {
		let importJson: any | null;
		try {
			importJson = JSON.parse(data);
		} catch {
			throw new Error('Please use a valid Addon export.');
		}

		// Parse all the settings.
		const charClass = nameToClass((importJson['class'] as string) || '');
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class!');
		}

		const race = nameToRace((importJson['race'] as string) || '');
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race!');
		}

		const professions = (importJson['professions'] as Array<{ name: string; level: number }>).map(profData => nameToProfession(profData.name));
		professions.forEach((prof, i) => {
			if (prof == Profession.ProfessionUnknown) {
				throw new Error(`Could not parse profession '${importJson['professions'][i]}'`);
			}
		});

		const talentsStr = (importJson['talents'] as string) || '';
		const glyphsConfig = classGlyphsConfig[charClass];

		const db = await Database.get();
		const primeGlyphIDs = (importJson['glyphs']['prime'] as Array<string | JsonObject>).map(g => glyphToID(g, db, glyphsConfig.primeGlyphs));
		const majorGlyphIDs = (importJson['glyphs']['major'] as Array<string | JsonObject>).map(g => glyphToID(g, db, glyphsConfig.majorGlyphs));
		const minorGlyphIDs = (importJson['glyphs']['minor'] as Array<string | JsonObject>).map(g => glyphToID(g, db, glyphsConfig.minorGlyphs));

		const glyphs = Glyphs.create({
			prime1: primeGlyphIDs[0] || 0,
			prime2: primeGlyphIDs[1] || 0,
			prime3: primeGlyphIDs[2] || 0,
			major1: majorGlyphIDs[0] || 0,
			major2: majorGlyphIDs[1] || 0,
			major3: majorGlyphIDs[2] || 0,
			minor1: minorGlyphIDs[0] || 0,
			minor2: minorGlyphIDs[1] || 0,
			minor3: minorGlyphIDs[2] || 0,
		});

		const gearJson = importJson['gear'];
		gearJson.items = (gearJson.items as Array<any>).filter(item => item != null);
		delete gearJson.version;

		(gearJson.items as Array<any>).forEach(item => {
			if (item.gems) {
				item.gems = (item.gems as Array<any>).map(gem => gem || 0);
			}
		});
		const equipmentSpec = EquipmentSpec.fromJson(gearJson);

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, talentsStr, glyphs, professions);
	}
}

function glyphNameToID(glyphName: string, glyphsConfig: Record<number, GlyphConfig>): number {
	if (!glyphName) {
		return 0;
	}

	for (const glyphIDStr in glyphsConfig) {
		if (glyphsConfig[glyphIDStr].name == glyphName) {
			return parseInt(glyphIDStr);
		}
	}
	throw new Error(`Unknown glyph name '${glyphName}'`);
}

function glyphToID(glyph: string | JsonObject, db: Database, glyphsConfig: Record<number, GlyphConfig>): number {
	if (typeof glyph === 'string') {
		// Legacy version: AddOn exports Glyphs by name (string) only. Names must be in English.
		return glyphNameToID(glyph, glyphsConfig);
	}

	// Cata version exports glyph information in a table that includes the name and the glyph spell ID.
	return db.glyphSpellToItemId(glyph.spellID as number);
}
