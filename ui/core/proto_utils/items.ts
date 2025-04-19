import { ArmorType, HandType, ItemQuality, ItemType, Stat, WeaponType } from '../proto/common';
import { ItemQualityValue, QualityValues } from '../proto/db';
import { UIItem } from '../proto/ui';
import { approximateScaleCoeff, randPropPoints, rangedTypes, thrownTypes, valueForQuality, wandTypes } from '../utils';
import { Database } from './database';

export const weaponDps = (item: UIItem, itemLevel: number): number => {
	const ilvl = itemLevel > 0 ? itemLevel : 0;
	const caster = (item.stats[Stat.StatSpellPower] ?? 0) > 0;
	const db = Database.getSync();
	let table: Array<ItemQualityValue> = [];
	switch (item.type) {
		case ItemType.ItemTypeWeapon:
			if (item.handType === HandType.HandTypeTwoHand) {
				table = caster ? db.weaponDamageValues.caster2H : db.weaponDamageValues.melee2H;
			} else {
				table = caster ? db.weaponDamageValues.caster1H : db.weaponDamageValues.melee1H;
			}
			break;

		case ItemType.ItemTypeRanged:
			const rw = item.rangedWeaponType;
			if (rw && rangedTypes.has(rw)) {
				table = db.weaponDamageValues.ranged;
			} else if (rw && thrownTypes.has(rw)) {
				table = db.weaponDamageValues.thrown;
			} else if (rw && wandTypes.has(rw)) {
				table = db.weaponDamageValues.wand;
			} else {
				return 0;
			}
			break;

		default:
			return 0;
	}

	const idx = Math.floor(ilvl);
	if (idx < 0 || idx >= table.length || !table[idx]) {
		return 0;
	}
	return valueForQuality(item.quality, table[idx].quality ?? QualityValues.create());
};

export const damageMin = (item: UIItem, itemLevel?: number): number => {
	const ilvl = itemLevel ?? item.ilvl;
	const baseDps = weaponDps(item, ilvl);
	const swing = item.weaponSpeed;
	let total = baseDps * swing * (1 - item.dmgVariance / 2) + item.qualityModifier * swing;
	if (total < 0) total = 1;
	return Math.floor(total);
};

export const damageMax = (item: UIItem, itemLevel?: number): number => {
	const ilvl = itemLevel ?? item.ilvl;
	const baseDps = weaponDps(item, ilvl);
	const swing = item.weaponSpeed;
	let total = baseDps * swing * (1 + item.dmgVariance / 2) + item.qualityModifier * (swing / 1000);
	if (total < 0) total = 1;
	return Math.floor(total + 0.5);
};

export const getArmorValue = (item: UIItem, itemLevel: number): number => {
	const db = Database.getSync();
	if (item.quality > ItemQuality.ItemQualityLegendary) {
		return 0;
	}

	const ilvl = itemLevel > 0 ? itemLevel : item.ilvl;

	// Shields have their own table
	if (item.weaponType === WeaponType.WeaponTypeShield) {
		const q = db.armorValues.shieldArmorValues[ilvl].quality;
		return Math.floor(valueForQuality(item.quality, q ?? QualityValues.create()) + 0.5);
	}

	if (item.armorType === ArmorType.ArmorTypeUnknown) {
		return 0;
	}

	const armorTotal = db.itemArmorTotal.get(ilvl);
	let baseArmor: number;
	switch (item.armorType) {
		case ArmorType.ArmorTypeCloth:
			baseArmor = armorTotal?.cloth ?? 0;
			break;
		case ArmorType.ArmorTypeLeather:
			baseArmor = armorTotal?.leather ?? 0;
			break;
		case ArmorType.ArmorTypeMail:
			baseArmor = armorTotal?.mail ?? 0;
			break;
		case ArmorType.ArmorTypePlate:
			baseArmor = armorTotal?.plate ?? 0;
			break;
		default:
			return 0;
	}

	const q = db.armorValues.armorValues[ilvl].quality;
	const qualityFactor = valueForQuality(item.quality, q ?? QualityValues.create());
	return Math.floor(baseArmor * qualityFactor * (this as any).ArmorModifier + 0.5);
};

export const getScaledStat = (item: UIItem, stat: Stat, itemLevel: number): number => {
	if (item.type === ItemType.ItemTypeUnknown) {
		return 0;
	}

	const budget = randPropPoints(itemLevel, item);
	const alloc = item.statAllocation[stat] ?? 0;
	if (alloc > 0 && budget > 0) {
		const raw = Math.round(alloc * budget * 0.0001);
		return raw - (item.socketModifier[stat] ?? 0);
	}

	const base = item.stats[stat] ?? 0;
	return Math.floor(base * approximateScaleCoeff(item.ilvl, itemLevel));
};

export const getStats = (item: UIItem, itemLevel: number): number[] => {
	const result: number[] = [];
	for (const statKey in item.statAllocation) {
		const s = Number(statKey) as Stat;
		result[s] = getScaledStat(item, s, itemLevel);
		if (s === Stat.StatAttackPower) {
			result[Stat.StatRangedAttackPower] = getScaledStat(item, s, itemLevel);
		}
	}
	return result;
};
