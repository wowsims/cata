import { Stats } from '../../core/proto_utils/stats';
import { CHARACTER_LEVEL } from '../constants/mechanics.js';
import {
	ConsumableType,
	EquipmentSpec,
	GemColor,
	ItemLevelState,
	ItemRandomSuffix,
	ItemSlot,
	ItemSpec,
	ItemSwap,
	PresetEncounter,
	PresetTarget,
	ReforgeStat,
	Stat,
} from '../proto/common.js';
import { Consumable, ItemEffectRandPropPoints, SimDatabase } from '../proto/db';
import { SpellEffect } from '../proto/spell';
import { GlyphID, IconData, UIDatabase, UIEnchant as Enchant, UIGem as Gem, UIItem as Item, UINPC as Npc, UIZone as Zone } from '../proto/ui.js';
import { distinct } from '../utils.js';
import { WOWHEAD_EXPANSION_ENV } from '../wowhead';
import { EquippedItem } from './equipped_item.js';
import { Gear, ItemSwapGear } from './gear.js';
import { gemEligibleForSocket, gemMatchesSocket } from './gems.js';
import { getEligibleEnchantSlots, getEligibleItemSlots } from './utils.js';

const dbUrlJson = '/mop/assets/database/db.json';
const dbUrlBin = '/mop/assets/database/db.bin';
const leftoversUrlJson = '/mop/assets/database/leftover_db.json';
const leftoversUrlBin = '/mop/assets/database/leftover_db.bin';
// When changing this value, don't forget to change the html <link> for preloading!
const READ_JSON = true;

export class Database {
	private static loadPromise: Promise<Database> | null = null;
	private static instance: Database | null = null;

	static async get(options: { signal?: AbortSignal } = {}): Promise<Database> {
		if (!Database.loadPromise) {
			Database.loadPromise = (async () => {
				let dbData: UIDatabase;
				if (READ_JSON) {
					const resp = await fetch(dbUrlJson, { signal: options?.signal });
					const json = await resp.json();
					dbData = UIDatabase.fromJson(json);
				} else {
					const buf = await fetch(dbUrlBin, { signal: options?.signal }).then(r => r.arrayBuffer());
					const bytes = new Uint8Array(buf);
					dbData = UIDatabase.fromBinary(bytes);
				}
				const db = new Database(dbData);
				Database.instance = db;
				return db;
			})();
		}
		return Database.loadPromise;
	}

	static getSync(): Database {
		if (!Database.instance) {
			throw new Error('Database not yet loaded; call `await Database.get()` before using getSync()');
		}
		return Database.instance;
	}

	static async getLeftovers(): Promise<UIDatabase> {
		if (READ_JSON) {
			return fetch(leftoversUrlJson)
				.then(response => response.json())
				.then(json => UIDatabase.fromJson(json));
		} else {
			return fetch(leftoversUrlBin)
				.then(response => response.arrayBuffer())
				.then(buffer => UIDatabase.fromBinary(new Uint8Array(buffer)));
		}
	}

	// Checks if any items in the equipment are missing from the current DB. If so, loads the leftover DB.
	static async loadLeftoversIfNecessary(equipment: EquipmentSpec): Promise<Database> {
		const db = await Database.get();
		if (db.loadedLeftovers) {
			return db;
		}

		const shouldLoadLeftovers = equipment.items.some(item => item.id != 0 && !db.items.has(item.id));
		if (shouldLoadLeftovers) {
			const leftoverDb = await Database.getLeftovers();
			db.loadProto(leftoverDb);
			db.loadedLeftovers = true;
		}
		return db;
	}

	private readonly items = new Map<number, Item>();
	private readonly randomSuffixes = new Map<number, ItemRandomSuffix>();
	private readonly reforgeStats = new Map<number, ReforgeStat>();
	private readonly itemEffectRandPropPoints = new Map<number, ItemEffectRandPropPoints>();
	private readonly enchantsBySlot: Partial<Record<ItemSlot, Enchant[]>> = {};
	private readonly gems = new Map<number, Gem>();
	private readonly npcs = new Map<number, Npc>();
	private readonly zones = new Map<number, Zone>();
	private readonly presetEncounters = new Map<string, PresetEncounter>();
	private readonly presetTargets = new Map<string, PresetTarget>();
	private readonly itemIcons: Record<number, Promise<IconData>> = {};
	private readonly spellIcons: Record<number, Promise<IconData>> = {};
	private readonly glyphIds: Array<GlyphID> = [];
	private readonly consumables = new Map<number, Consumable>();
	private readonly spellEffects = new Map<number, SpellEffect>();

	private loadedLeftovers = false;

	private constructor(db: UIDatabase) {
		this.loadProto(db);
	}

	// Add all data from the db proto into this database.
	private loadProto(db: UIDatabase) {
		db.items.forEach(item => {
			const itemCopy = { ...item };
			// Pre populate the item with stats from the highest base state the item can have.
			// We use this in EP calculations
			const maxScaling = item.scalingOptions[Math.max(...Object.keys(item.scalingOptions).map(Number))];
			itemCopy.weaponDamageMax = maxScaling.weaponDamageMax;
			itemCopy.weaponDamageMin = maxScaling.weaponDamageMin;
			itemCopy.randPropPoints = maxScaling.randPropPoints;
			itemCopy.ilvl = maxScaling.ilvl;
			itemCopy.stats = Stats.fromMap(maxScaling.stats).asProtoArray();

			this.items.set(itemCopy.id, itemCopy);
		});
		db.randomSuffixes.forEach(randomSuffix => this.randomSuffixes.set(randomSuffix.id, randomSuffix));
		db.reforgeStats.forEach(reforgeStat => this.reforgeStats.set(reforgeStat.id, reforgeStat));
		db.itemEffectRandPropPoints.forEach(ieRpp => this.itemEffectRandPropPoints.set(ieRpp.ilvl, ieRpp));
		db.enchants.forEach(enchant => {
			const slots = getEligibleEnchantSlots(enchant);
			slots.forEach(slot => {
				if (!this.enchantsBySlot[slot]) {
					this.enchantsBySlot[slot] = [];
				}
				this.enchantsBySlot[slot]!.push(enchant);
			});
		});
		db.gems.forEach(gem => this.gems.set(gem.id, gem));

		db.npcs.forEach(npc => this.npcs.set(npc.id, npc));
		db.zones.forEach(zone => this.zones.set(zone.id, zone));
		db.encounters.forEach(encounter => this.presetEncounters.set(encounter.path, encounter));
		db.encounters
			.map(e => e.targets)
			.flat()
			.forEach(target => this.presetTargets.set(target.path, target));

		db.items.forEach(
			item =>
				(this.itemIcons[item.id] = Promise.resolve(
					IconData.create({
						id: item.id,
						name: item.name,
						icon: item.icon,
					}),
				)),
		);
		db.gems.forEach(
			gem =>
				(this.itemIcons[gem.id] = Promise.resolve(
					IconData.create({
						id: gem.id,
						name: gem.name,
						icon: gem.icon,
					}),
				)),
		);
		db.itemIcons.forEach(data => (this.itemIcons[data.id] = Promise.resolve(data)));
		db.spellIcons.forEach(data => (this.spellIcons[data.id] = Promise.resolve(data)));
		db.glyphIds.forEach(id => this.glyphIds.push(id));
		db.consumables.forEach(consumable => this.consumables.set(consumable.id, consumable));
	}

	getAllItems(): Array<Item> {
		return Array.from(this.items.values());
	}

	getAllItemIds(): Array<number> {
		return Array.from(this.items.keys());
	}

	getItems(slot: ItemSlot): Array<Item> {
		return this.getAllItems().filter(item => getEligibleItemSlots(item).includes(slot));
	}

	getItemById(id: number): Item | undefined {
		return this.items.get(id);
	}

	getItemIdsForSet(setId: number): Array<number> {
		return this.getAllItemIds().filter(itemId => this.getItemById(itemId)!.setId === setId);
	}
	getSpellEffect(effectId: number): SpellEffect | undefined {
		return this.spellEffects.get(effectId);
	}
	getConsumables(): Array<Consumable> {
		return Array.from(this.consumables.values());
	}
	getConsumable(itemId: number): Consumable | undefined {
		return this.consumables.get(itemId);
	}
	getConsumablesByType(type: ConsumableType): Array<Consumable> {
		return this.getConsumables().filter(consume => consume.type == type);
	}
	getConsumablesByTypeAndStats(type: ConsumableType, stats: Array<Stat>): Array<Consumable> {
		return this.getConsumablesByType(type).filter(consume => consume.buffsMainStat || stats.some(index => consume.stats[index] > 0));
	}
	getRandomSuffixById(id: number): ItemRandomSuffix | undefined {
		return this.randomSuffixes.get(id);
	}

	getReforgeById(id: number): ReforgeStat | undefined {
		return this.reforgeStats.get(id);
	}

	getItemEffectRandPropPoints(ilvl: number) {
		return this.itemEffectRandPropPoints.get(ilvl);
	}

	getAvailableReforges(item: Item): ReforgeStat[] {
		return Array.from(this.reforgeStats.values()).filter(reforgeStat => item.stats[reforgeStat.fromStat] > 0 && item.stats[reforgeStat.toStat] == 0);
	}

	getEnchants(slot: ItemSlot): Array<Enchant> {
		return this.enchantsBySlot[slot] || [];
	}

	getGems(socketColor?: GemColor): Array<Gem> {
		if (!socketColor) return Array.from(this.gems.values());

		const ret = [];
		for (const g of this.gems.values()) {
			if (gemEligibleForSocket(g, socketColor)) ret.push(g);
		}
		return ret;
	}

	getNpc(npcId: number): Npc | null {
		return this.npcs.get(npcId) || null;
	}
	getZone(zoneId: number): Zone | null {
		return this.zones.get(zoneId) || null;
	}

	getMatchingGems(socketColor: GemColor): Array<Gem> {
		const ret = [];
		for (const g of this.gems.values()) {
			if (gemMatchesSocket(g, socketColor)) ret.push(g);
		}
		return ret;
	}

	lookupGem(itemID: number): Gem | null {
		return this.gems.get(itemID) || null;
	}

	lookupItemSpec(itemSpec: ItemSpec): EquippedItem | null {
		const item = this.items.get(itemSpec.id);
		if (!item) return null;

		let enchant: Enchant | null = null;
		if (itemSpec.enchant) {
			const slots = getEligibleItemSlots(item);
			for (let i = 0; i < slots.length; i++) {
				enchant =
					(this.enchantsBySlot[slots[i]] || []).find(enchant => [enchant.effectId, enchant.itemId, enchant.spellId].includes(itemSpec.enchant)) ||
					null;
				if (enchant) {
					break;
				}
			}
		}
		let tinker: Enchant | null = null;
		if (itemSpec.tinker) {
			const slots = getEligibleItemSlots(item);
			for (let i = 0; i < slots.length; i++) {
				tinker =
					(this.enchantsBySlot[slots[i]] || []).find(enchant => [enchant.effectId, enchant.itemId, enchant.spellId].includes(itemSpec.tinker)) ||
					null;
				if (tinker) {
					break;
				}
			}
		}

		const gems = itemSpec.gems.map(gemId => this.lookupGem(gemId));

		let randomSuffix: ItemRandomSuffix | null = null;
		if (itemSpec.randomSuffix && !!this.getRandomSuffixById(itemSpec.randomSuffix)) {
			randomSuffix = this.getRandomSuffixById(itemSpec.randomSuffix)!;
		}

		let reforge: ReforgeStat | null = null;
		if (itemSpec.reforging) {
			reforge = this.getReforgeById(itemSpec.reforging) || null;
		}

		return new EquippedItem({
			item,
			enchant,
			tinker,
			gems,
			randomSuffix,
			reforge,
			upgrade: itemSpec.upgradeStep ?? ItemLevelState.Base,
			challengeMode: itemSpec.challengeMode ?? false,
		});
	}

	lookupEquipmentSpec(equipSpec: EquipmentSpec): Gear {
		// EquipmentSpec is supposed to be indexed by slot, but here we assume
		// it isn't just in case.
		const gearMap: Partial<Record<ItemSlot, EquippedItem | null>> = {};
		equipSpec.items.forEach(itemSpec => {
			const item = this.lookupItemSpec(itemSpec);
			if (!item) return;

			const itemSlots = getEligibleItemSlots(item.item);

			const assignedSlot = itemSlots.find(slot => !gearMap[slot]);
			if (assignedSlot == null) throw new Error('No slots left to equip ' + Item.toJsonString(item.item));

			gearMap[assignedSlot] = item;
		});

		return new Gear(gearMap);
	}

	lookupItemSwap(itemSwap: ItemSwap): ItemSwapGear {
		const gearMap = itemSwap.items.reduce<Partial<Record<ItemSlot, EquippedItem | null>>>((gearMap, itemSpec, slot) => {
			const item = this.lookupItemSpec(itemSpec);
			if (item) {
				const eligibleItemSlots = getEligibleItemSlots(item.item);
				const isSwapSlotMatch = eligibleItemSlots.some(eligibleItemSlot => eligibleItemSlot === slot);
				const assignedSlot = isSwapSlotMatch ? (slot as ItemSlot) : eligibleItemSlots[0];
				if (typeof assignedSlot === 'number') gearMap[assignedSlot] = item;
			}
			return gearMap;
		}, {});
		return new ItemSwapGear(gearMap);
	}

	enchantSpellIdToEffectId(enchantSpellId: number): number {
		const enchant = Object.values(this.enchantsBySlot)
			.flat()
			.find(enchant => enchant.spellId == enchantSpellId);
		return enchant ? enchant.effectId : 0;
	}

	glyphItemToSpellId(itemId: number): number {
		return this.glyphIds.find(gid => gid.itemId == itemId)?.spellId || 0;
	}
	glyphSpellToItemId(spellId: number): number {
		return this.glyphIds.find(gid => gid.spellId == spellId)?.itemId || 0;
	}

	getPresetEncounter(path: string): PresetEncounter | null {
		return this.presetEncounters.get(path) || null;
	}
	getPresetTarget(path: string): PresetTarget | null {
		return this.presetTargets.get(path) || null;
	}
	getAllPresetEncounters(): Array<PresetEncounter> {
		return Array.from(this.presetEncounters.values());
	}
	getAllPresetTargets(): Array<PresetTarget> {
		return Array.from(this.presetTargets.values());
	}

	static async getItemIconData(itemId: number, options: { signal?: AbortSignal } = {}): Promise<IconData> {
		const db = await Database.get({ signal: options?.signal });
		const data = await db.spellIcons[itemId];

		if (!data?.icon) {
			db.itemIcons[itemId] = Database.getWowheadItemTooltipData(itemId, { signal: options?.signal });
		}
		return await db.itemIcons[itemId];
	}

	static async getSpellIconData(spellId: number, options: { signal?: AbortSignal } = {}): Promise<IconData> {
		const db = await Database.get({ signal: options?.signal });
		const data = await db.spellIcons[spellId];

		if (!data?.icon) {
			db.spellIcons[spellId] = Database.getWowheadSpellTooltipData(spellId, { signal: options?.signal });
		}
		return db.spellIcons[spellId];
	}

	private static async getWowheadItemTooltipData(id: number, options: { signal?: AbortSignal } = {}): Promise<IconData> {
		return Database.getWowheadTooltipData(id, 'item', { signal: options?.signal });
	}
	private static async getWowheadSpellTooltipData(id: number, options: { signal?: AbortSignal } = {}): Promise<IconData> {
		return Database.getWowheadTooltipData(id, 'spell', { signal: options?.signal });
	}
	private static async getWowheadTooltipData(id: number, tooltipPostfix: string, options: { signal?: AbortSignal } = {}): Promise<IconData> {
		const url = `https://nether.wowhead.com/mop-classic/tooltip/${tooltipPostfix}/${id}?lvl=${CHARACTER_LEVEL}&dataEnv=${WOWHEAD_EXPANSION_ENV}`;
		try {
			const response = await fetch(url, { signal: options?.signal });
			const json = await response.json();
			return IconData.create({
				id: id,
				name: json['name'],
				icon: json['icon'],
				hasBuff: json['buff'] !== '',
			});
		} catch (e) {
			if (e instanceof DOMException && e.name === 'AbortError') {
				return IconData.create();
			}
			console.error('Error while fetching url: ' + url + '\n\n' + e);
			return IconData.create();
		}
	}

	public static mergeSimDatabases(db1: SimDatabase, db2: SimDatabase): SimDatabase {
		return SimDatabase.create({
			items: distinct(db1.items.concat(db2.items), (a, b) => a.id == b.id),
			randomSuffixes: distinct(db1.randomSuffixes.concat(db2.randomSuffixes), (a, b) => a.id == b.id),
			reforgeStats: distinct(db1.reforgeStats.concat(db2.reforgeStats), (a, b) => a.id == b.id),
			itemEffectRandPropPoints: distinct(db1.itemEffectRandPropPoints.concat(db2.itemEffectRandPropPoints), (a, b) => a.ilvl == b.ilvl),
			enchants: distinct(db1.enchants.concat(db2.enchants), (a, b) => a.effectId == b.effectId),
			gems: distinct(db1.gems.concat(db2.gems), (a, b) => a.id == b.id),
			spellEffects: distinct(db1.spellEffects.concat(db2.spellEffects), (a, b) => a.id == b.id),
			consumables: distinct(db1.consumables.concat(db2.consumables), (a, b) => a.id == b.id),
		});
	}
}
