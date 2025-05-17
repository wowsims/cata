import { ArmorType, Class, ItemSlot, Profession, PseudoStat, Race, RangedWeaponType, Spec, Stat, WeaponType } from '../proto/common';
import { ResourceType, SecondaryResourceType } from '../proto/spell';
import { DungeonDifficulty, RaidFilterOption, RepFaction, RepLevel, SourceFilterOption, StatCapType } from '../proto/ui';

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
	[Race.RaceAlliancePandaren, 'Pandaren (A)'],
	[Race.RaceHordePandaren, 'Pandaren (H)'],
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
	[Class.ClassMonk, 'Monk'],
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
	[Profession.Archeology, 'Archeology'],
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

export function getStatName(stat: Stat): string {
	if (stat == Stat.StatRangedAttackPower) {
		return 'Ranged AP';
	} else {
		return Stat[stat]
			.split(/(?<![A-Z])(?=[A-Z])/)
			.slice(1)
			.join(' ');
	}
}

export function getClassPseudoStatName(pseudoStat: PseudoStat, playerClass: Class): string {
	const genericName = PseudoStat[pseudoStat]
		.split(/(?<![A-Z])(?=[A-Z])/)
		.slice(2)
		.join(' ')
		.replace('Dps', 'DPS');

	if (playerClass == Class.ClassHunter) {
		return genericName.replace('Physical', 'Ranged');
	} else {
		return genericName.replace('Physical', 'Melee');
	}
}

// TODO: Make sure BE exports the spell schools properly
export enum SpellSchool {
	None = 0,
	Physical = 1 << 1,
	Arcane = 1 << 2,
	Fire = 1 << 3,
	Frost = 1 << 4,
	Holy = 1 << 5,
	Nature = 1 << 6,
	Shadow = 1 << 7,
}

export const spellSchoolNames: Map<number, string> = new Map([
	[SpellSchool.Physical, 'Physical'],
	[SpellSchool.Arcane, 'Arcane'],
	[SpellSchool.Fire, 'Fire'],
	[SpellSchool.Frost, 'Frost'],
	[SpellSchool.Holy, 'Holy'],
	[SpellSchool.Nature, 'Nature'],
	[SpellSchool.Shadow, 'Shadow'],
	[SpellSchool.Nature + SpellSchool.Arcane, 'Astral'],
	[SpellSchool.Shadow + SpellSchool.Fire, 'Shadowflame'],
	[SpellSchool.Fire + SpellSchool.Arcane, 'Spellfire'],
	[SpellSchool.Arcane + SpellSchool.Frost, 'Spellfrost'],
	[SpellSchool.Frost + SpellSchool.Fire, 'Frostfire'],
	[SpellSchool.Shadow + SpellSchool.Frost, 'Shadowfrost'],
]);

export const shortSecondaryStatNames: Map<Stat, string> = new Map([
	[Stat.StatSpirit, 'Spirit'],
	[Stat.StatHitRating, 'Hit'],
	[Stat.StatCritRating, 'Crit'],
	[Stat.StatHasteRating, 'Haste'],
	[Stat.StatExpertiseRating, 'Expertise'],
	[Stat.StatMasteryRating, 'Mastery'],
	[Stat.StatDodgeRating, 'Dodge'],
	[Stat.StatParryRating, 'Parry'],
]);

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
]);

export const resourceNames: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, 'None'],
	[ResourceType.ResourceTypeHealth, 'Health'],
	[ResourceType.ResourceTypeMana, 'Mana'],
	[ResourceType.ResourceTypeEnergy, 'Energy'],
	[ResourceType.ResourceTypeRage, 'Rage'],
	[ResourceType.ResourceTypeChi, 'Chi'],
	[ResourceType.ResourceTypeComboPoints, 'Combo Points'],
	[ResourceType.ResourceTypeFocus, 'Focus'],
	[ResourceType.ResourceTypeRunicPower, 'Runic Power'],
	[ResourceType.ResourceTypeBloodRune, 'Blood Rune'],
	[ResourceType.ResourceTypeFrostRune, 'Frost Rune'],
	[ResourceType.ResourceTypeUnholyRune, 'Unholy Rune'],
	[ResourceType.ResourceTypeDeathRune, 'Death Rune'],
	[ResourceType.ResourceTypeSolarEnergy, 'Solar Energy'],
	[ResourceType.ResourceTypeLunarEnergy, 'Lunar Energy'],
	[ResourceType.ResourceTypeGenericResource, 'Generic Resource'],
]);

export const resourceColors: Map<ResourceType, string> = new Map([
	[ResourceType.ResourceTypeNone, '#ffffff'],
	[ResourceType.ResourceTypeHealth, '#22ba00'],
	[ResourceType.ResourceTypeMana, '#2e93fa'],
	[ResourceType.ResourceTypeEnergy, '#ffd700'],
	[ResourceType.ResourceTypeRage, '#ff0000'],
	[ResourceType.ResourceTypeChi, '#00ff98'],
	[ResourceType.ResourceTypeComboPoints, '#ffa07a'],
	[ResourceType.ResourceTypeFocus, '#cd853f'],
	[ResourceType.ResourceTypeRunicPower, '#5b99ee'],
	[ResourceType.ResourceTypeBloodRune, '#ff0000'],
	[ResourceType.ResourceTypeFrostRune, '#0000ff'],
	[ResourceType.ResourceTypeUnholyRune, '#00ff00'],
	[ResourceType.ResourceTypeDeathRune, '#8b008b'],
	[ResourceType.ResourceTypeSolarEnergy, '#d2952b'],
	[ResourceType.ResourceTypeLunarEnergy, '#2c4f8f'],
	[ResourceType.ResourceTypeGenericResource, '#ffffff'],
]);

export function stringToResourceType(str: string): [ResourceType, SecondaryResourceType | undefined] {
	for (const [key, val] of resourceNames) {
		if (val.toLowerCase() == str.toLowerCase()) {
			return [key, undefined];
		}
	}

	for (const val of Object.keys(SecondaryResourceType).filter(key=> isNaN(Number(key)))) {
		if (val.toLowerCase() == str.toLowerCase()) {
			return [ResourceType.ResourceTypeGenericResource, (<any>SecondaryResourceType)[val]];
		}
	}

	return [ResourceType.ResourceTypeNone, undefined];
}

export const sourceNames: Map<SourceFilterOption, string> = new Map([
	[SourceFilterOption.SourceUnknown, 'Unknown'],
	[SourceFilterOption.SourceCrafting, 'Crafting'],
	[SourceFilterOption.SourceQuest, 'Quest'],
	[SourceFilterOption.SourceReputation, 'Reputation'],
	[SourceFilterOption.SourcePvp, 'PVP'],
	[SourceFilterOption.SourceDungeon, 'Dungeon'],
	[SourceFilterOption.SourceDungeonH, 'Dungeon (H)'],
	[SourceFilterOption.SourceRaid, 'Raid'],
	[SourceFilterOption.SourceRaidH, 'Raid (H)'],
	[SourceFilterOption.SourceRaidRF, 'Raid (RF)'],
]);
export const raidNames: Map<RaidFilterOption, string> = new Map([
	[RaidFilterOption.RaidUnknown, 'Unknown'],
	[RaidFilterOption.RaidIcecrownCitadel, 'Icecrown Citadel'],
	[RaidFilterOption.RaidRubySanctum, 'Ruby Sanctum'],
	[RaidFilterOption.RaidBlackwingDescent, 'Blackwing Descent'],
	[RaidFilterOption.RaidTheBastionOfTwilight, 'The Bastion of Twilight'],
	[RaidFilterOption.RaidBaradinHold, 'Baradin Hold'],
	[RaidFilterOption.RaidThroneOfTheFourWinds, 'Throne of the Four Winds'],
	[RaidFilterOption.RaidFirelands, 'Firelands'],
	[RaidFilterOption.RaidDragonSoul, 'Dragon Soul'],
]);

export const difficultyNames: Map<DungeonDifficulty, string> = new Map([
	[DungeonDifficulty.DifficultyUnknown, 'Unknown'],
	[DungeonDifficulty.DifficultyNormal, 'N'],
	[DungeonDifficulty.DifficultyHeroic, 'H'],
	[DungeonDifficulty.DifficultyTitanRuneAlpha, 'TRA'],
	[DungeonDifficulty.DifficultyTitanRuneBeta, 'TRB'],
	[DungeonDifficulty.DifficultyRaid10, '10N'],
	[DungeonDifficulty.DifficultyRaid10H, '10H'],
	[DungeonDifficulty.DifficultyRaid25RF, 'RF'],
	[DungeonDifficulty.DifficultyRaid25, 'RN'],
	[DungeonDifficulty.DifficultyRaid25H, 'RH'],
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
	[RepFaction.RepFactionTheEarthenRing]: 'The Earthen Ring',
	[RepFaction.RepFactionGuardiansOfHyjal]: 'Guardians of Hyjal',
	[RepFaction.RepFactionTherazane]: 'Therazane',
	[RepFaction.RepFactionDragonmawClan]: 'Dragonmaw Clan',
	[RepFaction.RepFactionRamkahen]: 'Ramkahen',
	[RepFaction.RepFactionWildhammerClan]: 'Wildhammer Clan',
	[RepFaction.RepFactionBaradinsWardens]: "Baradin's Wardens",
	[RepFaction.RepFactionHellscreamsReach]: "Hellscream's Reach",
	[RepFaction.RepFactionAvengersOfHyjal]: 'Avengers of Hyjal',
};

export const REP_FACTION_QUARTERMASTERS: Record<RepFaction, number> = {
	[RepFaction.RepFactionUnknown]: 0,
	[RepFaction.RepFactionTheEarthenRing]: 50324,
	[RepFaction.RepFactionGuardiansOfHyjal]: 50314,
	[RepFaction.RepFactionTherazane]: 45408,
	[RepFaction.RepFactionDragonmawClan]: 49387,
	[RepFaction.RepFactionRamkahen]: 48617,
	[RepFaction.RepFactionWildhammerClan]: 49386,
	[RepFaction.RepFactionBaradinsWardens]: 47328,
	[RepFaction.RepFactionHellscreamsReach]: 48531,
	[RepFaction.RepFactionAvengersOfHyjal]: 54401,
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
	[Spec.SpecGuardianDruid, 'Savage Defender'],
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
	[Spec.SpecBrewmasterMonk, 'Elusive Brawler'],
	[Spec.SpecMistweaverMonk, 'Gift of the Serpent'],
	[Spec.SpecWindwalkerMonk, 'Bottled Fury'],
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
	[Spec.SpecGuardianDruid, 77494],
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
	[Spec.SpecBrewmasterMonk, 117906],
	[Spec.SpecMistweaverMonk, 117907],
	[Spec.SpecWindwalkerMonk, 115636],
]);
export const statCapTypeNames = new Map<StatCapType, string>([
	[StatCapType.TypeHardCap, 'Hard cap'],
	[StatCapType.TypeSoftCap, 'Soft cap'],
	[StatCapType.TypeThreshold, 'Threshold'],
]);
