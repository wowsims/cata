import { CHARACTER_LEVEL } from '../../../constants/mechanics';
import { IndividualSimUI } from '../../../individual_sim_ui';
import { ItemSlot, Spec } from '../../../proto/common';
import { raceNames } from '../../../proto_utils/names';
import { IndividualWowheadGearPlannerImporter } from '../importers';
import { IndividualExporter } from './individual_exporter';

const c = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';

function writeBits(value: number): number[] {
	let e = value;
	let t = 0;
	const bits: number[] = [];

	for (let a = 1; a <= 5; a++) {
		const n = 5 * a;
		if (e < 1 << n) {
			const nArray = [];
			while (nArray.length < a) {
				const t = e & 63;
				e >>= 6;
				nArray.unshift(t);
			}
			nArray[0] = nArray[0] | t;
			bits.push(...nArray);
			return bits;
		}
		e -= 1 << n;
		t = (64 | t) >> 1;
	}
	throw new Error('Value too large to encode.');
}

function writeTalents(talentStr: string): number[] {
	const bits: number[] = [];
	const trees = talentStr.split('-');

	for (let a = 0; a < 3; a++) {
		const tree = trees[a] || '';
		bits.push(...writeBits(tree.length));

		let l = 0;
		while (l < tree.length) {
			let chunk = 0;
			let s = 0;
			while (s < 7 && l < tree.length) {
				const digit = parseInt(tree[l], 10);
				chunk = (chunk << 3) | digit;
				l++;
				s++;
			}
			bits.push(...writeBits(chunk));
		}
	}

	return bits;
}

// Function to write glyphs (reverse of parseGlyphs)
function writeGlyphs(glyphIds: number[]): string {
	const base32 = '0123456789abcdefghjkmnpqrstvwxyz'; // Base32 character set
	let glyphStr = '0'; // insert random 0

	for (let i = 0; i < glyphIds.length; i++) {
		const spellId = glyphIds[i];
		if (spellId) {
			const glyphSlotChar = base32[i];
			const c1 = (spellId >> 15) & 31;
			const c2 = (spellId >> 10) & 31;
			const c3 = (spellId >> 5) & 31;
			const c4 = spellId & 31;

			if (c1 < 0 || c2 < 0 || c3 < 0 || c4 < 0) {
				continue; // Invalid spell ID
			}

			glyphStr += glyphSlotChar + base32[c1] + base32[c2] + base32[c3] + base32[c4];
		}
	}

	return glyphStr;
}

// Function to write the hash (reverse of readHash)
function writeHash(data: any): string {
	let hash = '';

	// Initialize bits array
	const bits: number[] = [];

	// Starting character (B for gear planner)
	const idx = 1; // Assuming idx is 1 for gear planner

	// Include idx in the hash (as first character)
	hash += c[idx];

	// Gender (assuming genderId is 0 or 1)
	bits.push(0);

	// Level
	bits.push(...writeBits(data.level ?? 0));

	// Talents
	const talentBits = writeTalents(data.talents.join('-'));
	bits.push(...talentBits);

	// Glyphs
	const glyphStr = writeGlyphs(data.glyphs ?? []);
	const glyphBytes = glyphStr.split('').map(ch => c.indexOf(ch));
	bits.push(...writeBits(glyphBytes.length));
	bits.push(...glyphBytes);

	// Items
	const items = data.items ?? [];
	bits.push(...writeBits(items.length));

	for (const item of items) {
		let e = 0;
		const itemBits: number[] = [];

		// Encode flags into e
		if (item.randomEnchantId) e |= 1 << 6;
		if (item.reforge) e |= 1 << 5;
		const gemCount = Math.min((item.gemItemIds ?? []).length, 7);
		e |= gemCount << 2;
		const enchantCount = Math.min((item.enchantIds ?? []).length, 3);
		e |= enchantCount;

		// Item slot and ID
		itemBits.push(...writeBits(item.slotId ?? 0));
		itemBits.push(...writeBits(item.itemId ?? 0));

		// Random Enchant ID
		if (item.randomEnchantId) {
			let enchant = item.randomEnchantId;
			const negative = enchant < 0 ? 1 : 0;
			if (negative) enchant *= -1;
			enchant = (enchant << 1) | negative;
			itemBits.push(...writeBits(enchant));
		}

		// Reforge
		if (item.reforge) {
			itemBits.push(...writeBits(item.reforge));
		}

		// Gems
		const gems = item.gemItemIds ?? [];
		for (let i = 0; i < gemCount; i++) {
			itemBits.push(...writeBits(gems[i]));
		}

		// Enchants
		const enchants = item.enchantIds ?? [];
		for (let i = 0; i < enchantCount; i++) {
			itemBits.push(...writeBits(enchants[i]));
		}

		// e is the item flags; add it at the start of itemBits
		bits.push(...writeBits(e));
		bits.push(...itemBits);
	}

	// Encode bits into characters
	let hashData = '';
	for (const bit of bits) {
		hashData += c.charAt(bit);
	}

	// Append the hash data to the URL
	if (hashData) {
		hash += hashData;
	}

	return hash;
}

export interface WowheadGearPlannerData {
	class?: string;
	race?: string;
	genderId?: number;
	level: number;
	talents: string[];
	glyphs: number[];
	items: WowheadItemData[];
}

export interface WowheadItemData {
	slotId: number;
	itemId: number;
	randomEnchantId?: number;
	reforge?: number;
	gemItemIds?: number[];
	enchantIds?: number[];
}

export function createWowheadGearPlannerLink(data: WowheadGearPlannerData): string {
	const baseUrl = '';
	const hash = writeHash(data);
	return baseUrl + hash;
}

export class IndividualWowheadGearPlannerExporter<SpecType extends Spec> extends IndividualExporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Wowhead Export', allowDownload: true });
	}

	getData(): string {
		const player = this.simUI.player;

		const classStr = player.getPlayerClass().friendlyName.replaceAll(/\s/g, '-').toLowerCase();
		const raceStr = raceNames.get(player.getRace())!.replaceAll(/\s/g, '-').toLowerCase();
		const url = `https://www.wowhead.com/mop-classic/gear-planner/${classStr}/${raceStr}/`;

		const addGlyph = (glyphItemId: number): number => {
			const spellId = this.simUI.sim.db.glyphItemToSpellId(glyphItemId);
			if (!spellId) {
				return 0;
			}
			return spellId;
		};

		const glyphs = player.getGlyphs();

		const data = {
			level: CHARACTER_LEVEL,
			talents: player.getTalentsString().split('-'),
			glyphs: [
				addGlyph(glyphs.major1),
				addGlyph(glyphs.major2),
				addGlyph(glyphs.major3),
				addGlyph(glyphs.minor1),
				addGlyph(glyphs.minor2),
				addGlyph(glyphs.minor3),
			],
			items: [],
		} as WowheadGearPlannerData;

		const gear = player.getGear();

		gear.getItemSlots()
			.sort((slot1, slot2) => IndividualWowheadGearPlannerImporter.slotIDs[slot1] - IndividualWowheadGearPlannerImporter.slotIDs[slot2])
			.forEach(itemSlot => {
				const item = gear.getEquippedItem(itemSlot);
				if (!item) {
					return;
				}

				const slotId = IndividualWowheadGearPlannerImporter.slotIDs[itemSlot];
				const itemData = {
					slotId: slotId,
					itemId: item.id,
				} as WowheadItemData;
				if (item._randomSuffix?.id) {
					itemData.randomEnchantId = item._randomSuffix.id;
				}
				if (item._enchant) {
					itemData.enchantIds = [item._enchant.spellId];
				}

				if (ItemSlot.ItemSlotHands == itemSlot) {
					//Todo: IF Hands we want to append any tinkers if existing
				}

				if (item._gems) {
					itemData.gemItemIds = item._gems.map(gem => {
						return gem?.id ?? 0;
					});
				}
				if (item._reforge) {
					itemData.reforge = item._reforge.id;
				}
				data.items.push(itemData);
			});

		const hash = createWowheadGearPlannerLink(data);

		return url + hash;
	}
}
