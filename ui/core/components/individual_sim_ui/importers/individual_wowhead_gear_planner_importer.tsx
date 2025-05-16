import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../../individual_sim_ui';
import { Class, EquipmentSpec, Glyphs, ItemSlot, ItemSpec, Race, Spec } from '../../../proto/common';
import { nameToClass, nameToRace } from '../../../proto_utils/names';
import Toast from '../../toast';
import { IndividualImporter } from './individual_importer';

const c = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';

// Taken from Wowhead
function readBits(e: number[]): number {
	if (!e.length) return 0;
	let t = 0,
		a = 1,
		n = e[0];
	while ((32 & n) > 0) {
		a++;
		n <<= 1;
	}
	const l = 63 >> a;
	let s = e.shift()! & l;
	a--;
	for (let n = 1; n <= a; n++) {
		t += 1 << (5 * n);
		s = (s << 6) | (e.shift() || 0);
	}
	return s + t;
}

// Taken from Wowhead
function parseTalents(e: number[]): string {
	let t = '';
	for (let a = 0; a < 3; a++) {
		let len = readBits(e);
		while (len > 0) {
			let n = '',
				l = readBits(e);
			while (len > 0 && n.length < 7) {
				const digit = 7 & l;
				l >>= 3;
				n = `${digit}${n}`;
				len--;
			}
			t += n;
		}
		t += '-';
	}
	return t.replace(/-+$/, '');
}

// Taken from Wowhead
function readHash(e: string): any {
	const t: any = {};
	const a = e.match(/^#?~(-?\d+)$/);
	if (a) {
		return t;
	}
	const n = /^([a-z-]+)\/([a-z-]+)(?:\/([a-zA-Z0-9_-]+))?$/.exec(e);
	if (!n) return t;
	{
		t.class = n[1];
	}
	{
		t.race = n[2];
	}
	let o = n[3];
	if (!o) return t;
	const idx = c.indexOf(o.substring(0, 1));
	o = o.substring(1);
	if (!o.length) return t;
	const u: number[] = [];
	for (let e = 0; e < o.length; e++) u.push(c.indexOf(o.substring(e, e + 1)));
	if (idx > 1) return t;
	{
		const e = readBits(u) - 1;
		if (e >= 0) t.genderId = e;
	}
	{
		const e = readBits(u);
		if (e) t.level = e;
	}
	{
		const e = parseTalents(u),
			a = readBits(u),
			n = u
				.splice(0, a)
				.map(e => c[e])
				.join(''),
			l = e + (a ? `_${n}` : '');
		if ('' !== l) t.talentHash = l;
		if (t.talentHash) {
			const talents = parseTalentString(t.talentHash);
			t.talents = talents.talents;
			t.glyphs = talents.glyphs;
		}
	}
	{
		let itemCount = readBits(u);
		t.items = [];
		while (itemCount--) {
			const a: any = {};
			let e: number;
			if (idx < 1) {
				e = u.shift()!;
				const t = (e >> 5) & 1;
				e &= 31;
				if (t) e |= 64;
			} else e = readBits(u);
			a.slotId = readBits(u);
			a.itemId = readBits(u);
			if (0 != ((e >> 6) & 1)) {
				let enchant = readBits(u);
				const t = 1 & enchant;
				enchant >>= 1;
				if (t) enchant *= -1;
				a.randomEnchantId = enchant;
			}
			if (0 != ((e >> 5) & 1)) a.reforge = readBits(u);
			{
				let gemCount = (e >> 2) & 7;
				while (gemCount--) {
					if (!a.gemItemIds) a.gemItemIds = [];
					a.gemItemIds.push(readBits(u));
				}
			}
			{
				let enchantCount = e & 3;
				while (enchantCount--) {
					if (!a.enchantIds) a.enchantIds = [];
					a.enchantIds.push(readBits(u));
				}
			}
			t.items.push(a);
		}
	}
	return t;
}

// Function to parse glyphs from the glyph string
function parseGlyphs(glyphStr: string): number[] {
	const glyphIds = Array(9).fill(0); // Nine potential glyph slots
	const base32 = '0123456789abcdefghjkmnpqrstvwxyz'; // Base32 character set
	let cur = 1; // we skip the first index for whatever reason

	while (cur < glyphStr.length) {
		// Get glyph slot index
		const glyphSlotChar = glyphStr[cur];
		const glyphSlotIndex = base32.indexOf(glyphSlotChar);
		cur++;

		if (glyphSlotIndex < 0 || glyphSlotIndex >= glyphIds.length) {
			continue; // Skip invalid glyph slots
		}

		if (cur + 4 > glyphStr.length) {
			break; // Not enough characters for a glyph ID
		}

		// Decode the spellId using base32 encoding (each character represents 5 bits)
		const c1 = base32.indexOf(glyphStr[cur]);
		const c2 = base32.indexOf(glyphStr[cur + 1]);
		const c3 = base32.indexOf(glyphStr[cur + 2]);
		const c4 = base32.indexOf(glyphStr[cur + 3]);
		cur += 4;

		if (c1 < 0 || c2 < 0 || c3 < 0 || c4 < 0) {
			continue; // Invalid character in spell ID
		}

		const spellId = (c1 << 15) | (c2 << 10) | (c3 << 5) | c4;

		glyphIds[glyphSlotIndex] = spellId;
	}

	return glyphIds;
}

function parseTalentString(talentString: string): { talents: string; glyphs: number[] } {
	const [talentPart, glyphPart] = talentString.split('_');

	// Parse the talents
	// Talent string is something like '001-2301-33223203120220120321'
	// Each part separated by '-' corresponds to a talent tree
	const talents = talentPart;

	// Parse the glyphs
	let glyphs: number[] = [];
	if (glyphPart) {
		glyphs = parseGlyphs(glyphPart);
	}

	return { talents, glyphs };
}

function parseWowheadGearLink(link: string): any {
	// Extract the part after 'mop-classic/gear-planner/'
	const match = link.match(/mop-classic\/gear-planner\/(.+)/);
	if (!match) {
		throw new Error(`Invalid WCL URL ${link}, must look like "https://www.wowhead.com/mop-classic/gear-planner/CLASS/RACE/XXXX"`);
	}
	const e = match[1];
	return readHash(e);
}

export class IndividualWowheadGearPlannerImporter<SpecType extends Spec> extends IndividualImporter<SpecType> {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Wowhead Import', allowFileUpload: true });

		const warningRef = ref<HTMLDivElement>();
		this.descriptionElem.appendChild(
			<>
				<p>
					Import settings from{' '}
					<a href="https://www.wowhead.com/mop-classic/gear-planner" target="_blank">
						Wowhead Gear Planner
					</a>
					.
				</p>
				<p>This feature imports gear, race, and (optionally) talents. It does NOT import buffs, debuffs, consumes, rotation, or custom stats.</p>
				<p>To import, paste the gear planner link below and click, 'Import'.</p>
				<div ref={warningRef} />
			</>,
		);

		if (warningRef.value)
			new Toast({
				title: 'Tinker issues',
				body: (
					<>
						There are known issues importing tinkers from Wowhead.
						<br />
						Always make sure to double check your tinkers after importing.
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

	async onImport(url: string) {
		const match = url.match(/www\.wowhead\.com\/mop-classic\/gear-planner\/([a-z\-]+)\/([a-z\-]+)\/([a-zA-Z0-9_\-]+)/);
		if (!match) {
			throw new Error(`Invalid WCL URL ${url}, must look like "https://www.wowhead.com/mop-classic/gear-planner/CLASS/RACE/XXXX"`);
		}
		console.log(url);

		const parsed = parseWowheadGearLink(url);
		console.log(parsed);
		const glyphIds = parsed.glyphs;

		const charClass = nameToClass(parsed.class.replaceAll('-', ''));
		if (charClass == Class.ClassUnknown) {
			throw new Error('Could not parse Class: ' + parsed.class);
		}

		const race = nameToRace(parsed.race.replaceAll('-', ''));
		if (race == Race.RaceUnknown) {
			throw new Error('Could not parse Race: ' + parsed.race);
		}

		const equipmentSpec = EquipmentSpec.create();

		parsed.items.forEach((item: any) => {
			const itemSpec = ItemSpec.create();
			const slotId = item.slotId;
			const isEnchanted = item.enchantIds?.length > 0;
			itemSpec.id = item.itemId;
			if (isEnchanted) {
				itemSpec.enchant = this.simUI.sim.db.enchantSpellIdToEffectId(item.enchantIds[0]);
			}
			if (item.gemItemIds) {
				itemSpec.gems = item.gemItemIds;
			}
			if (item.randomEnchantId) {
				itemSpec.randomSuffix = item.randomEnchantId;
			}
			if (item.reforge) {
				itemSpec.reforging = item.reforge;
			}
			const itemSlotEntry = Object.entries(IndividualWowheadGearPlannerImporter.slotIDs).find(e => e[1] == slotId);
			if (itemSlotEntry != null) {
				equipmentSpec.items.push(itemSpec);
			}
		});

		const glyphs = Glyphs.create({
			major1: this.simUI.sim.db.glyphSpellToItemId(glyphIds[3]),
			major2: this.simUI.sim.db.glyphSpellToItemId(glyphIds[4]),
			major3: this.simUI.sim.db.glyphSpellToItemId(glyphIds[5]),
			minor1: this.simUI.sim.db.glyphSpellToItemId(glyphIds[6]),
			minor2: this.simUI.sim.db.glyphSpellToItemId(glyphIds[7]),
			minor3: this.simUI.sim.db.glyphSpellToItemId(glyphIds[8]),
		});

		this.finishIndividualImport(this.simUI, charClass, race, equipmentSpec, parsed.talents ?? '', glyphs, []);
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
	};
}
