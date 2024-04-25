import { ResourceType } from '../proto/api.js';
import { ArmorType, Class, ItemSlot, Profession, PseudoStat, Race, RangedWeaponType, Spec, Stat, WeaponType } from '../proto/common.js';
import { DungeonDifficulty, RaidFilterOption, RepFaction, RepLevel, SourceFilterOption } from '../proto/ui.js';

export const armorTypeNames: Map<ArmorType, string> = new Map([
	[ArmorType.ArmorTypeUnknown, 'Unknown'],
	[ArmorType.ArmorTypeCloth, 'Cloth'],
	[ArmorType.ArmorTypeLeather, 'Leather'],
	[ArmorType.ArmorTypeMail, 'Mail'],
	[ArmorType.ArmorTypePlate, 'Plate'],
]);

export const weaponTypeNames: Map<WeaponType, string> = new Map([
	[WeaponType.WeaponTypeUnknown, 'Unknown'],
	[WeaponType.WeaponTypeAxe, 'Axe'],
	[WeaponType.WeaponTypeDagger, 'Dagger'],
	[WeaponType.WeaponTypeFist, 'Fist'],
	[WeaponType.WeaponTypeMace, 'Mace'],
	[WeaponType.WeaponTypeOffHand, 'Misc'],
	[WeaponType.WeaponTypePolearm, 'Polearm'],
	[WeaponType.WeaponTypeShield, 'Shield'],
	[WeaponType.WeaponTypeStaff, 'Staff'],
	[WeaponType.WeaponTypeSword, 'Sword'],
]);

export const rangedWeaponTypeNames: Map<RangedWeaponType, string> = new Map([
	[RangedWeaponType.RangedWeaponTypeUnknown, 'Unknown'],
	[RangedWeaponType.RangedWeaponTypeBow, 'Bow'],
	[RangedWeaponType.RangedWeaponTypeCrossbow, 'Crossbow'],
	[RangedWeaponType.RangedWeaponTypeGun, 'Gun'],
	[RangedWeaponType.RangedWeaponTypeRelic, 'Relic'],
	[RangedWeaponType.RangedWeaponTypeThrown, 'Thrown'],
	[RangedWeaponType.RangedWeaponTypeWand, 'Wand'],
]);

export const raceNames: Map<Race, string> = new Map([
	[Race.RaceUnknown, 'None'],
	[Race.RaceBloodElf, 'Blood Elf'],
	[Race.RaceDraenei, 'Draenei'],
	[Race.RaceDwarf, 'Dwarf'],
	[Race.RaceGnome, 'Gnome'],
	[Race.RaceGoblin, 'Goblin'],
	[Race.RaceHuman, 'Human'],
	[Race.RaceNightElf, 'Night Elf'],
	[Race.RaceOrc, 'Orc'],
	[Race.RaceTauren, 'Tauren'],
	[Race.RaceTroll, 'Troll'],
	[Race.RaceUndead, 'Undead'],
	[Race.RaceWorgen, 'Worgen'],
]);

export function nameToRace(name: string): Race {
	const normalized = name.toLowerCase().replaceAll(' ', '');
	for (const [key, value] of raceNames) {
		if (value.toLowerCase().replaceAll(' ', '') == normalized) {
			return key;
		}
	}
	return Race.RaceUnknown;
}

export const classNames: Map<Class, string> = new Map([
	[Class.ClassUnknown, 'None'],
	[Class.ClassDruid, 'Druid'],
	[Class.ClassHunter, 'Hunter'],
	[Class.ClassMage, 'Mage'],
	[Class.ClassPaladin, 'Paladin'],
	[Class.ClassPriest, 'Priest'],
	[Class.ClassRogue, 'Rogue'],
	[Class.ClassShaman, 'Shaman'],
	[Class.ClassWarlock, 'Warlock'],
	[Class.ClassWarrior, 'Warrior'],
	[Class.ClassDeathKnight, 'Death Knight'],
]);

export function nameToClass(name: string): Class {
	const lower = name.toLowerCase();
	for (const [key, value] of classNames) {
		if (value.toLowerCase().replace(/\s+/g, '') == lower) {
			return key;
		}
	}
	return Class.ClassUnknown;
}

export const professionNames: Map<Profession, string> = new Map([
	[Profession.ProfessionUnknown, 'None'],
	[Profession.Alchemy, 'Alchemy'],
	[Profession.Blacksmithing, 'Blacksmithing'],
	[Profession.Enchanting, 'Enchanting'],
	[Profession.Engineering, 'Engineering'],
	[Profession.Herbalism, 'Herbalism'],
	[Profession.Inscription, 'Inscription'],
	[Profession.Jewelcrafting, 'Jewelcrafting'],
	[Profession.Leatherworking, 'Leatherworking'],
	[Profession.Mining, 'Mining'],
	[Profession.Skinning, 'Skinning'],
	[Profession.Tailoring, 'Tailoring'],
]);

export function nameToProfession(name: string): Profession {
	const lower = name.toLowerCase();
	for (const [key, value] of professionNames) {
		if (value.toLowerCase() == lower) {
			return key;
		}
	}
	return Profession.ProfessionUnknown;
}

export const statOrder: Array<Stat> = [
	Stat.StatHealth,
	Stat.StatMana,
	Stat.StatArmor,
	Stat.StatBonusArmor,
	Stat.StatStamina,
	Stat.StatStrength,
	Stat.StatAgility,
	Stat.StatIntellect,
	Stat.StatSpirit,
	Stat.StatSpellPower,
	Stat.StatSpellHit,
	Stat.StatSpellCrit,
	Stat.StatSpellHaste,
	Stat.StatSpellPenetration,
	Stat.StatMP5,
	Stat.StatAttackPower,
	Stat.StatRangedAttackPower,
	Stat.StatMeleeHit,
	Stat.StatMeleeCrit,
	Stat.StatMeleeHaste,
	Stat.StatArmorPenetration,
	Stat.StatExpertise,
	Stat.StatMastery,
	Stat.StatDefense,
	Stat.StatBlock,
	Stat.StatBlockValue,
	Stat.StatDodge,
	Stat.StatParry,
	Stat.StatResilience,
	Stat.StatArcaneResistance,
	Stat.StatFireResistance,
	Stat.StatFrostResistance,
	Stat.StatNatureResistance,
	Stat.StatShadowResistance,
];

export const statNames: Map<Stat, string> = new Map([
	[Stat.StatStrength, 'Strength'],
	[Stat.StatAgility, 'Agility'],
	[Stat.StatStamina, 'Stamina'],
	[Stat.StatIntellect, 'Intellect'],
	[Stat.StatSpirit, 'Spirit'],
	[Stat.StatSpellPower, 'Spell Power'],
	[Stat.StatMP5, 'MP5'],
	[Stat.StatSpellHit, 'Spell Hit'],
	[Stat.StatSpellCrit, 'Spell Crit'],
	[Stat.StatSpellHaste, 'Spell Haste'],
	[Stat.StatSpellPenetration, 'Spell Pen'],
	[Stat.StatAttackPower, 'Attack Power'],
	[Stat.StatMeleeHit, 'Melee Hit'],
	[Stat.StatMeleeCrit, 'Melee Crit'],
	[Stat.StatMeleeHaste, 'Melee Haste'],
	[Stat.StatArmorPenetration, 'Armor Pen'],
	[Stat.StatMastery, 'Mastery'],
	[Stat.StatExpertise, 'Expertise'],
	[Stat.StatMana, 'Mana'],
	[Stat.StatArmor, 'Armor'],
	[Stat.StatRangedAttackPower, 'Ranged AP'],
	[Stat.StatDefense, 'Defense'],
	[Stat.StatBlock, 'Block'],
	[Stat.StatBlockValue, 'Block Value'],
	[Stat.StatDodge, 'Dodge'],
	[Stat.StatParry, 'Parry'],
	[Stat.StatResilience, 'Resilience'],
	[Stat.StatHealth, 'Health'],
	[Stat.StatArcaneResistance, 'Arcane Resistance'],
	[Stat.StatFireResistance, 'Fire Resistance'],
	[Stat.StatFrostResistance, 'Frost Resistance'],
	[Stat.StatNatureResistance, 'Nature Resistance'],
	[Stat.StatShadowResistance, 'Shadow Resistance'],
	[Stat.StatBonusArmor, 'Bonus Armor'],
	[Stat.StatMastery, 'Mastery'],
]);

export const shortSecondaryStatNames: Map<Stat, string> = new Map([
	[Stat.StatSpirit, 'Spirit'],
	[Stat.StatSpellHit, 'Hit'],
	[Stat.StatSpellCrit, 'Crit'],
	[Stat.StatSpellHaste, 'Haste'],
	[Stat.StatMeleeHit, 'Hit'],
	[Stat.StatMeleeCrit, 'Crit'],
	[Stat.StatMeleeHaste, 'Haste'],
	[Stat.StatExpertise, 'Expertise'],
	[Stat.StatMastery, 'Mastery'],
	[Stat.StatDefense, 'Defense'],
	[Stat.StatBlock, 'Block'],
	[Stat.StatBlockValue, 'Block Value'],
	[Stat.StatDodge, 'Dodge'],
	[Stat.StatParry, 'Parry'],
]);

export const pseudoStatOrder: Array<PseudoStat> = [
	PseudoStat.PseudoStatMainHandDps,
	PseudoStat.PseudoStatOffHandDps,
	PseudoStat.PseudoStatRangedDps,
	PseudoStat.PseudoStatBlockValueMultiplier,
];
export const pseudoStatNames: Map<PseudoStat, string> = new Map([
	[PseudoStat.PseudoStatMainHandDps, 'Main Hand DPS'],
	[PseudoStat.PseudoStatOffHandDps, 'Off Hand DPS'],
	[PseudoStat.PseudoStatRangedDps, 'Ranged DPS'],
	[PseudoStat.PseudoStatBlockValueMultiplier, 'Block Value Multiplier'],
	[PseudoStat.PseudoStatDodge, 'Dodge Chance'],
	[PseudoStat.PseudoStatParry, 'Parry Chance'],
]);

export function getClassStatName(stat: Stat, playerClass: Class): string {
	const statName = statNames.get(stat);
	if (!statName) return 'UnknownStat';
	if (playerClass == Class.ClassHunter) {
		return statName.replace('Melee', 'Ranged');
	} else {
		return statName;
	}
}

export const slotNames: Map<ItemSlot, string> = new Map([
	[ItemSlot.ItemSlotHead, 'Head'],
	[ItemSlot.ItemSlotNeck, 'Neck'],
	[ItemSlot.ItemSlotShoulder, 'Shoulders'],
	[ItemSlot.ItemSlotBack, 'Back'],
	[ItemSlot.ItemSlotChest, 'Chest'],
	[ItemSlot.ItemSlotWrist, 'Wrist'],
	[ItemSlot.ItemSlotHands, 'Hands'],
	[ItemSlot.ItemSlotWaist, 'Waist'],
	[ItemSlot.ItemSlotLegs, 'Legs'],
	[ItemSlot.ItemSlotFeet, 'Feet'],
	[ItemSlot.ItemSlotFinger1, 'Finger 1'],
	[ItemSlot.ItemSlotFinger2, 'Finger 2'],
	[ItemSlot.ItemSlotTrinket1, 'Trinket 1'],
	[ItemSlot.ItemSlotTrinket2, 'Trinket 2'],
	[ItemSlot.ItemSlotMainHand, 'Main Hand'],
	[ItemSlot.ItemSlotOffHand, 'Off Hand'],
	[ItemSlot.ItemSlotRanged, 'Ranged'],
]);

export const resourceNames: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, 'None'],
	[ResourceType.ResourceTypeHealth, 'Health'],
	[ResourceType.ResourceTypeMana, 'Mana'],
	[ResourceType.ResourceTypeEnergy, 'Energy'],
	[ResourceType.ResourceTypeRage, 'Rage'],
	[ResourceType.ResourceTypeComboPoints, 'Combo Points'],
	[ResourceType.ResourceTypeFocus, 'Focus'],
	[ResourceType.ResourceTypeRunicPower, 'Runic Power'],
	[ResourceType.ResourceTypeBloodRune, 'Blood Rune'],
	[ResourceType.ResourceTypeFrostRune, 'Frost Rune'],
	[ResourceType.ResourceTypeUnholyRune, 'Unholy Rune'],
	[ResourceType.ResourceTypeDeathRune, 'Death Rune'],
	[ResourceType.ResourceTypeSolarEnergy, 'Solar Energy'],
	[ResourceType.ResourceTypeLunarEnergy, 'Lunar Energy'],
]);

export const resourceColors: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, '#ffffff'],
	[ResourceType.ResourceTypeHealth, '#22ba00'],
	[ResourceType.ResourceTypeMana, '#2e93fa'],
	[ResourceType.ResourceTypeEnergy, '#ffd700'],
	[ResourceType.ResourceTypeRage, '#ff0000'],
	[ResourceType.ResourceTypeComboPoints, '#ffa07a'],
	[ResourceType.ResourceTypeFocus, '#cd853f'],
	[ResourceType.ResourceTypeRunicPower, '#5b99ee'],
	[ResourceType.ResourceTypeBloodRune, '#ff0000'],
	[ResourceType.ResourceTypeFrostRune, '#0000ff'],
	[ResourceType.ResourceTypeUnholyRune, '#00ff00'],
	[ResourceType.ResourceTypeDeathRune, '#8b008b'],
	[ResourceType.ResourceTypeSolarEnergy, '#d2952b'],
	[ResourceType.ResourceTypeLunarEnergy, '#2c4f8f'],
]);

export function stringToResourceType(str: string): ResourceType {
	for (const [key, val] of resourceNames) {
		if (val.toLowerCase() == str.toLowerCase()) {
			return key;
		}
	}
	return ResourceType.ResourceTypeNone;
}

export const sourceNames: Map<SourceFilterOption, string> = new Map([
	[SourceFilterOption.SourceUnknown, 'Unknown'],
	[SourceFilterOption.SourceCrafting, 'Crafting'],
	[SourceFilterOption.SourceQuest, 'Quest'],
	[SourceFilterOption.SourceDungeon, 'Dungeon'],
	[SourceFilterOption.SourceDungeonH, 'Dungeon (H)'],
	[SourceFilterOption.SourceDungeonTRA, 'Dungeon (TRA)'],
	[SourceFilterOption.SourceDungeonTRB, 'Dungeon (TRB)'],
	[SourceFilterOption.SourceRaid10, 'Raid (10N)'],
	[SourceFilterOption.SourceRaid10H, 'Raid (10H)'],
	[SourceFilterOption.SourceRaid25, 'Raid (25N)'],
	[SourceFilterOption.SourceRaid25H, 'Raid (25H)'],
]);
export const raidNames: Map<RaidFilterOption, string> = new Map([
	[RaidFilterOption.RaidUnknown, 'Unknown'],
	[RaidFilterOption.RaidVanilla, 'Vanilla'],
	[RaidFilterOption.RaidTbc, 'TBC'],
	[RaidFilterOption.RaidNaxxramas, 'Naxxramas'],
	[RaidFilterOption.RaidEyeOfEternity, 'Eye of Eternity'],
	[RaidFilterOption.RaidObsidianSanctum, 'Obsidian Sanctum'],
	[RaidFilterOption.RaidVaultOfArchavon, 'Vault of Archavon'],
	[RaidFilterOption.RaidUlduar, 'Ulduar'],
	[RaidFilterOption.RaidTrialOfTheCrusader, 'Trial of the Crusader'],
	[RaidFilterOption.RaidOnyxiasLair, "Onyxia's Lair"],
	[RaidFilterOption.RaidIcecrownCitadel, 'Icecrown Citadel'],
	[RaidFilterOption.RaidRubySanctum, 'Ruby Sanctum'],
]);

export const difficultyNames: Map<DungeonDifficulty, string> = new Map([
	[DungeonDifficulty.DifficultyUnknown, 'Unknown'],
	[DungeonDifficulty.DifficultyNormal, 'N'],
	[DungeonDifficulty.DifficultyHeroic, 'H'],
	[DungeonDifficulty.DifficultyTitanRuneAlpha, 'TRA'],
	[DungeonDifficulty.DifficultyTitanRuneBeta, 'TRB'],
	[DungeonDifficulty.DifficultyRaid10, '10N'],
	[DungeonDifficulty.DifficultyRaid10H, '10H'],
	[DungeonDifficulty.DifficultyRaid25, '25N'],
	[DungeonDifficulty.DifficultyRaid25H, '25H'],
]);

export const REP_LEVEL_NAMES: Record<RepLevel, string> = {
	[RepLevel.RepLevelUnknown]: 'Unknown',
	[RepLevel.RepLevelHated]: 'Hated',
	[RepLevel.RepLevelHostile]: 'Hostile',
	[RepLevel.RepLevelUnfriendly]: 'Unfriendly',
	[RepLevel.RepLevelNeutral]: 'Neutral',
	[RepLevel.RepLevelFriendly]: 'Friendly',
	[RepLevel.RepLevelHonored]: 'Honored',
	[RepLevel.RepLevelRevered]: 'Revered',
	[RepLevel.RepLevelExalted]: 'Exalted',
};

export const REP_FACTION_NAMES: Record<RepFaction, string> = {
	[RepFaction.RepFactionUnknown]: 'Unknown',
};

export const masterySpellNames: Map<Spec, string> = new Map([
	[Spec.SpecAssassinationRogue, 'Potent Poisons'],
	[Spec.SpecCombatRogue, 'Main Gauche'],
	[Spec.SpecSubtletyRogue, 'Executioner'],
	[Spec.SpecBloodDeathKnight, 'Blood Shield'],
	[Spec.SpecFrostDeathKnight, 'Frozen Heart'],
	[Spec.SpecUnholyDeathKnight, 'Dreadblade'],
	[Spec.SpecBalanceDruid, 'Total Eclipse'],
	[Spec.SpecFeralDruid, 'Razor Claws'],
	[Spec.SpecRestorationDruid, 'Harmony'],
	[Spec.SpecHolyPaladin, 'Illuminated Healing'],
	[Spec.SpecProtectionPaladin, 'Divine Bulwark'],
	[Spec.SpecRetributionPaladin, 'Hand of Light'],
	[Spec.SpecElementalShaman, 'Elemental Overload'],
	[Spec.SpecEnhancementShaman, 'Enhanced Elements'],
	[Spec.SpecRestorationShaman, 'Deep Healing'],
	[Spec.SpecBeastMasteryHunter, 'Master of Beasts'],
	[Spec.SpecMarksmanshipHunter, 'Wild Quiver'],
	[Spec.SpecSurvivalHunter, 'Essence of the Viper'],
	[Spec.SpecArmsWarrior, 'Strikes of Opportunity'],
	[Spec.SpecFuryWarrior, 'Unshackled Fury'],
	[Spec.SpecProtectionWarrior, 'Critical Block'],
	[Spec.SpecArcaneMage, 'Mana Adept'],
	[Spec.SpecFireMage, 'Flashburn'],
	[Spec.SpecFrostMage, 'Frostburn'],
	[Spec.SpecDisciplinePriest, 'Shield Discipline'],
	[Spec.SpecHolyPriest, 'Echo of Light'],
	[Spec.SpecShadowPriest, 'Shadow Orb Power'],
	[Spec.SpecAfflictionWarlock, 'Potent Afflictions'],
	[Spec.SpecDemonologyWarlock, 'Master Demonologist'],
	[Spec.SpecDestructionWarlock, 'Fiery Apocalypse'],
]);

export const masterySpellIDs: Map<Spec, number> = new Map([
	[Spec.SpecAssassinationRogue, 76803],
	[Spec.SpecCombatRogue, 76806],
	[Spec.SpecSubtletyRogue, 76808],
	[Spec.SpecBloodDeathKnight, 77513],
	[Spec.SpecFrostDeathKnight, 77514],
	[Spec.SpecUnholyDeathKnight, 77515],
	[Spec.SpecBalanceDruid, 77492],
	[Spec.SpecFeralDruid, 77493],
	[Spec.SpecRestorationDruid, 77495],
	[Spec.SpecHolyPaladin, 76669],
	[Spec.SpecProtectionPaladin, 76671],
	[Spec.SpecRetributionPaladin, 76672],
	[Spec.SpecElementalShaman, 77222],
	[Spec.SpecEnhancementShaman, 77223],
	[Spec.SpecRestorationShaman, 77226],
	[Spec.SpecBeastMasteryHunter, 76657],
	[Spec.SpecMarksmanshipHunter, 76659],
	[Spec.SpecSurvivalHunter, 76658],
	[Spec.SpecArmsWarrior, 76838],
	[Spec.SpecFuryWarrior, 76856],
	[Spec.SpecProtectionWarrior, 76857],
	[Spec.SpecArcaneMage, 76547],
	[Spec.SpecFireMage, 76595],
	[Spec.SpecFrostMage, 76613],
	[Spec.SpecDisciplinePriest, 77484],
	[Spec.SpecHolyPriest, 77485],
	[Spec.SpecShadowPriest, 77486],
	[Spec.SpecAfflictionWarlock, 77215],
	[Spec.SpecDemonologyWarlock, 77219],
	[Spec.SpecDestructionWarlock, 77220],
]);
