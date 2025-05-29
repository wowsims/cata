package dbc

import "github.com/wowsims/mop/sim/core/proto"

func MapResistanceToStat(index int) (proto.Stat, bool) {
	switch index {
	case 0:
		return proto.Stat_StatBonusArmor, true
	}
	return proto.Stat_StatBonusArmor, false
}

var MapArmorSubclassToArmorType = map[int]proto.ArmorType{
	ITEM_SUBCLASS_ARMOR_CLOTH:   proto.ArmorType_ArmorTypeCloth,
	ITEM_SUBCLASS_ARMOR_LEATHER: proto.ArmorType_ArmorTypeLeather,
	ITEM_SUBCLASS_ARMOR_MAIL:    proto.ArmorType_ArmorTypeMail,
	ITEM_SUBCLASS_ARMOR_PLATE:   proto.ArmorType_ArmorTypePlate,
	0:                           proto.ArmorType_ArmorTypeUnknown,
}

func MapMainStatToStat(index int) (proto.Stat, bool) {
	switch index {
	case 0:
		return proto.Stat_StatStrength, true
	case 1:
		return proto.Stat_StatAgility, true
	case 2:
		return proto.Stat_StatStamina, true
	case 3:
		return proto.Stat_StatIntellect, true
	case 4:
		return proto.Stat_StatSpirit, true
	}
	return 0, false
}
func MapBonusStatIndexToStat(index int) (proto.Stat, bool) {
	switch index {
	case 0: // Mana
		return proto.Stat_StatMana, true
	case 1: // Health
		return proto.Stat_StatHealth, true
	case 7: // Stamina
		return proto.Stat_StatStamina, true
	case 3: // Agility
		return proto.Stat_StatAgility, true
	case 4: // Strength
		return proto.Stat_StatStrength, true
	case 5: // Intellect
		return proto.Stat_StatIntellect, true
	case 6: // Spirit
		return proto.Stat_StatSpirit, true

	// Secondary ratings (reforge-able)
	case 16, 17, 18, 31: // MeleeHitRating, RangedHitRating, SpellHitRating, or generic HitRating
		return proto.Stat_StatHitRating, true
	case 19, 20, 21, 32: // MeleeCritRating, RangedCritRating, SpellCritRating, or generic CritRating
		return proto.Stat_StatCritRating, true
	case 36: // HasteRating (non-obsolete)
		return proto.Stat_StatHasteRating, true
	case 37: // ExpertiseRating
		return proto.Stat_StatExpertiseRating, true
	case 13: // DodgeRating
		return proto.Stat_StatDodgeRating, true
	case 14: // ParryRating
		return proto.Stat_StatParryRating, true
	case 49: // Mastery
		return proto.Stat_StatMasteryRating, true
	case 38: // AttackPower
		return proto.Stat_StatAttackPower, true
	case 39: // RangedAttackPower
		return proto.Stat_StatRangedAttackPower, true
	case 41, 42, 45: // SpellHealing, SpellDamage, or SpellPower
		return proto.Stat_StatSpellPower, true
	case 57: // PvPPowerRating
		return proto.Stat_StatPvpPowerRating, true
	case 35: // ResilienceRating
		return proto.Stat_StatPvpResilienceRating, true
	case 50: // ExtraArmor maps to BonusArmor (green armor)
		return proto.Stat_StatBonusArmor, true
	case 43: // ManaRegeneration
		return proto.Stat_StatMP5, true
	default:
		return 0, false
	}
}

var MapProfessionIdToProfession = map[int]proto.Profession{
	0:   proto.Profession_ProfessionUnknown,
	164: proto.Profession_Blacksmithing,
	165: proto.Profession_Leatherworking,
	171: proto.Profession_Alchemy,
	182: proto.Profession_Herbalism,
	186: proto.Profession_Mining,
	197: proto.Profession_Tailoring,
	202: proto.Profession_Engineering,
	333: proto.Profession_Enchanting,
	393: proto.Profession_Skinning,
	755: proto.Profession_Jewelcrafting,
	773: proto.Profession_Inscription,
	794: proto.Profession_Archeology,
}

var MapItemSubclassNames = map[ItemSubClass]string{
	OneHandedAxes:    "One-Handed Axes",
	TwoHandedAxes:    "Two-Handed Axes",
	Bows:             "Bows",
	Guns:             "Guns",
	OneHandedMaces:   "One-Handed Maces",
	TwoHandedMaces:   "Two-Handed Maces",
	Polearms:         "Polearms",
	OneHandedSwords:  "One-Handed Swords",
	TwoHandedSwords:  "Two-Handed Swords",
	Staves:           "Staves",
	OneHandedExotics: "One-Handed Exotics",
	TwoHandedExotics: "Two-Handed Exotics",
	FistWeapons:      "Fist Weapons",
	Daggers:          "Daggers",
}

var MapSocketTypeToGemColor = map[int]proto.GemColor{
	0: proto.GemColor_GemColorUnknown,
	1: proto.GemColor_GemColorMeta,
	2: proto.GemColor_GemColorRed,
	3: proto.GemColor_GemColorYellow,
	4: proto.GemColor_GemColorBlue,
	5: proto.GemColor_GemColorShaTouched,
	6: proto.GemColor_GemColorCogwheel,
	7: proto.GemColor_GemColorPrismatic,
}

var MapInventoryTypeToItemType = map[int]proto.ItemType{
	0:                      proto.ItemType_ItemTypeUnknown,
	INVTYPE_HEAD:           proto.ItemType_ItemTypeHead,
	INVTYPE_NECK:           proto.ItemType_ItemTypeNeck,
	INVTYPE_SHOULDERS:      proto.ItemType_ItemTypeShoulder,
	INVTYPE_CHEST:          proto.ItemType_ItemTypeChest,
	INVTYPE_WAIST:          proto.ItemType_ItemTypeWaist,
	INVTYPE_LEGS:           proto.ItemType_ItemTypeLegs,
	INVTYPE_FEET:           proto.ItemType_ItemTypeFeet,
	INVTYPE_WRISTS:         proto.ItemType_ItemTypeWrist,
	INVTYPE_HANDS:          proto.ItemType_ItemTypeHands,
	INVTYPE_FINGER:         proto.ItemType_ItemTypeFinger,
	INVTYPE_TRINKET:        proto.ItemType_ItemTypeTrinket,
	INVTYPE_WEAPON:         proto.ItemType_ItemTypeWeapon,
	INVTYPE_SHIELD:         proto.ItemType_ItemTypeWeapon,
	INVTYPE_RANGED:         proto.ItemType_ItemTypeRanged,
	INVTYPE_CLOAK:          proto.ItemType_ItemTypeBack,
	INVTYPE_2HWEAPON:       proto.ItemType_ItemTypeWeapon,
	INVTYPE_BAG:            proto.ItemType_ItemTypeUnknown,
	INVTYPE_TABARD:         proto.ItemType_ItemTypeUnknown,
	INVTYPE_ROBE:           proto.ItemType_ItemTypeChest,
	INVTYPE_WEAPONMAINHAND: proto.ItemType_ItemTypeWeapon,
	INVTYPE_WEAPONOFFHAND:  proto.ItemType_ItemTypeWeapon,
	INVTYPE_HOLDABLE:       proto.ItemType_ItemTypeWeapon,
	INVTYPE_AMMO:           proto.ItemType_ItemTypeUnknown,
	INVTYPE_THROWN:         proto.ItemType_ItemTypeRanged,
	INVTYPE_RANGEDRIGHT:    proto.ItemType_ItemTypeRanged,
	INVTYPE_QUIVER:         proto.ItemType_ItemTypeRanged,
	INVTYPE_RELIC:          proto.ItemType_ItemTypeRanged,
}
var MapInventoryTypeFlagToItemType = map[InventoryTypeFlag]proto.ItemType{
	0:                proto.ItemType_ItemTypeUnknown,
	HEAD:             proto.ItemType_ItemTypeHead,
	NECK:             proto.ItemType_ItemTypeNeck,
	SHOULDER:         proto.ItemType_ItemTypeShoulder,
	CHEST:            proto.ItemType_ItemTypeChest,
	WAIST:            proto.ItemType_ItemTypeWaist,
	LEGS:             proto.ItemType_ItemTypeLegs,
	FEET:             proto.ItemType_ItemTypeFeet,
	WRIST:            proto.ItemType_ItemTypeWrist,
	HAND:             proto.ItemType_ItemTypeHands,
	FINGER:           proto.ItemType_ItemTypeFinger,
	TRINKET:          proto.ItemType_ItemTypeTrinket,
	MAIN_HAND:        proto.ItemType_ItemTypeWeapon,
	OFF_HAND:         proto.ItemType_ItemTypeWeapon,
	RANGED:           proto.ItemType_ItemTypeRanged,
	CLOAK:            proto.ItemType_ItemTypeBack,
	TWO_H_WEAPON:     proto.ItemType_ItemTypeWeapon,
	BAG:              proto.ItemType_ItemTypeUnknown,
	TABARD:           proto.ItemType_ItemTypeUnknown,
	ROBE:             proto.ItemType_ItemTypeChest,
	WEAPON_MAIN_HAND: proto.ItemType_ItemTypeWeapon,
	WEAPON_OFF_HAND:  proto.ItemType_ItemTypeWeapon,
	HOLDABLE:         proto.ItemType_ItemTypeWeapon,
	AMMO:             proto.ItemType_ItemTypeUnknown,
	THROWN:           proto.ItemType_ItemTypeRanged,
	RANGED_RIGHT:     proto.ItemType_ItemTypeRanged,
	QUIVER:           proto.ItemType_ItemTypeRanged,
	RELIC:            proto.ItemType_ItemTypeRanged,
}

var MapWeaponSubClassToWeaponType = map[int]proto.WeaponType{
	ITEM_SUBCLASS_WEAPON_AXE:          proto.WeaponType_WeaponTypeAxe,
	ITEM_SUBCLASS_WEAPON_AXE2:         proto.WeaponType_WeaponTypeAxe,
	ITEM_SUBCLASS_WEAPON_BOW:          proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_GUN:          proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_MACE:         proto.WeaponType_WeaponTypeMace,
	ITEM_SUBCLASS_WEAPON_MACE2:        proto.WeaponType_WeaponTypeMace,
	ITEM_SUBCLASS_WEAPON_POLEARM:      proto.WeaponType_WeaponTypePolearm,
	ITEM_SUBCLASS_WEAPON_SWORD:        proto.WeaponType_WeaponTypeSword,
	ITEM_SUBCLASS_WEAPON_SWORD2:       proto.WeaponType_WeaponTypeSword,
	ITEM_SUBCLASS_WEAPON_WARGLAIVE:    proto.WeaponType_WeaponTypePolearm, // assuming polearm idk
	ITEM_SUBCLASS_WEAPON_STAFF:        proto.WeaponType_WeaponTypeStaff,
	ITEM_SUBCLASS_WEAPON_EXOTIC:       proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_EXOTIC2:      proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_FIST:         proto.WeaponType_WeaponTypeFist,
	ITEM_SUBCLASS_WEAPON_MISC:         proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_DAGGER:       proto.WeaponType_WeaponTypeDagger,
	ITEM_SUBCLASS_WEAPON_THROWN:       proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_SPEAR:        proto.WeaponType_WeaponTypePolearm,
	ITEM_SUBCLASS_WEAPON_CROSSBOW:     proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_WAND:         proto.WeaponType_WeaponTypeUnknown,
	ITEM_SUBCLASS_WEAPON_FISHING_POLE: proto.WeaponType_WeaponTypeUnknown,
}

type EnchantMetaType struct {
	ItemType   proto.ItemType
	WeaponType proto.WeaponType
}

var SpellSchoolToStat = map[SpellSchool]proto.Stat{
	FIRE:     -1,
	ARCANE:   -1,
	NATURE:   -1,
	FROST:    -1,
	SHADOW:   -1,
	PHYSICAL: proto.Stat_StatArmor,
}
var MapInventoryTypeToEnchantMetaType = map[InventoryTypeFlag]EnchantMetaType{
	HEAD:     {ItemType: proto.ItemType_ItemTypeHead, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	NECK:     {ItemType: proto.ItemType_ItemTypeNeck, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	SHOULDER: {ItemType: proto.ItemType_ItemTypeShoulder, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	CHEST:    {ItemType: proto.ItemType_ItemTypeChest, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	WAIST:    {ItemType: proto.ItemType_ItemTypeWaist, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	LEGS:     {ItemType: proto.ItemType_ItemTypeLegs, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	FEET:     {ItemType: proto.ItemType_ItemTypeFeet, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	WRIST:    {ItemType: proto.ItemType_ItemTypeWrist, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	HAND:     {ItemType: proto.ItemType_ItemTypeHands, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	FINGER:   {ItemType: proto.ItemType_ItemTypeFinger, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	TRINKET:  {ItemType: proto.ItemType_ItemTypeTrinket, WeaponType: proto.WeaponType_WeaponTypeUnknown},

	WEAPON_MAIN_HAND: {ItemType: proto.ItemType_ItemTypeWeapon, WeaponType: proto.WeaponType_WeaponTypeUnknown}, // One-Hand
	WEAPON_OFF_HAND:  {ItemType: proto.ItemType_ItemTypeWeapon, WeaponType: proto.WeaponType_WeaponTypeShield},  // Off Hand
	RANGED:           {ItemType: proto.ItemType_ItemTypeRanged, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	CLOAK:            {ItemType: proto.ItemType_ItemTypeBack, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	TWO_H_WEAPON:     {ItemType: proto.ItemType_ItemTypeWeapon, WeaponType: proto.WeaponType_WeaponTypeUnknown},
}
var consumableClassToProto = map[ConsumableClass]proto.ConsumableType{
	EXPLOSIVES_AND_DEVICES: proto.ConsumableType_ConsumableTypeExplosive,
	POTION:                 proto.ConsumableType_ConsumableTypePotion,
	FLASK:                  proto.ConsumableType_ConsumableTypeFlask,
	SCROLL:                 proto.ConsumableType_ConsumableTypeScroll,
	FOOD:                   proto.ConsumableType_ConsumableTypeFood,
	BANDAGE:                proto.ConsumableType_ConsumableTypeUnknown,
	OTHER:                  proto.ConsumableType_ConsumableTypeUnknown,
}

var MapPowerTypeEnumToResourceType = map[int32]proto.ResourceType{
	0:  proto.ResourceType_ResourceTypeMana,
	1:  proto.ResourceType_ResourceTypeRage,
	2:  proto.ResourceType_ResourceTypeFocus,
	3:  proto.ResourceType_ResourceTypeEnergy,
	4:  proto.ResourceType_ResourceTypeComboPoints,
	5:  proto.ResourceType_ResourceTypeDeathRune | proto.ResourceType_ResourceTypeBloodRune,
	6:  proto.ResourceType_ResourceTypeRunicPower,
	7:  proto.ResourceType_ResourceTypeNone, // Soulshards
	8:  proto.ResourceType_ResourceTypeLunarEnergy,
	9:  proto.ResourceType_ResourceTypeNone, // Holy Power
	12: proto.ResourceType_ResourceTypeChi,
	20: proto.ResourceType_ResourceTypeBloodRune,
	21: proto.ResourceType_ResourceTypeFrostRune,
	22: proto.ResourceType_ResourceTypeUnholyRune,
	29: proto.ResourceType_ResourceTypeDeathRune,
}

func ClassNameFromDBC(dbc DbcClass) string {
	switch dbc.ID {
	case 1:
		return "Warrior"
	case 2:
		return "Paladin"
	case 3:
		return "Hunter"
	case 4:
		return "Rogue"
	case 5:
		return "Priest"
	case 6:
		return "Death_Knight"
	case 7:
		return "Shaman"
	case 8:
		return "Mage"
	case 9:
		return "Warlock"
	case 10:
		return "Monk"
	case 11:
		return "Druid"
	default:
		return "Unknown"
	}
}
func getMatchingRatingMods(value int) []RatingModType {
	allMods := []RatingModType{
		RATING_MOD_DODGE,
		RATING_MOD_PARRY,
		RATING_MOD_HIT_MELEE,
		RATING_MOD_HIT_RANGED,
		RATING_MOD_HIT_SPELL,
		RATING_MOD_CRIT_MELEE,
		RATING_MOD_CRIT_RANGED,
		RATING_MOD_CRIT_SPELL,
		RATING_MOD_MULTISTRIKE,
		RATING_MOD_READINESS,
		RATING_MOD_SPEED,
		RATING_MOD_RESILIENCE,
		RATING_MOD_LEECH,
		RATING_MOD_HASTE_MELEE,
		RATING_MOD_HASTE_RANGED,
		RATING_MOD_HASTE_SPELL,
		RATING_MOD_AVOIDANCE,
		RATING_MOD_EXPERTISE,
		RATING_MOD_MASTERY,
		RATING_MOD_PVP_POWER,
		RATING_MOD_VERS_DAMAGE,
		RATING_MOD_VERS_HEAL,
		RATING_MOD_VERS_MITIG,
	}

	var result []RatingModType
	for _, mod := range allMods {
		if value&int(mod) != 0 {
			result = append(result, mod)
		}
	}
	return result
}

var RatingModToStat = map[RatingModType]proto.Stat{
	RATING_MOD_DODGE:        proto.Stat_StatDodgeRating,
	RATING_MOD_PARRY:        proto.Stat_StatParryRating,
	RATING_MOD_HIT_MELEE:    proto.Stat_StatHitRating,
	RATING_MOD_HIT_RANGED:   proto.Stat_StatHitRating,
	RATING_MOD_HIT_SPELL:    proto.Stat_StatHitRating,
	RATING_MOD_CRIT_MELEE:   proto.Stat_StatCritRating,
	RATING_MOD_CRIT_RANGED:  proto.Stat_StatCritRating,
	RATING_MOD_CRIT_SPELL:   proto.Stat_StatCritRating,
	RATING_MOD_MULTISTRIKE:  -1,
	RATING_MOD_READINESS:    -1,
	RATING_MOD_SPEED:        -1,
	RATING_MOD_RESILIENCE:   proto.Stat_StatPvpResilienceRating,
	RATING_MOD_LEECH:        -1,
	RATING_MOD_HASTE_MELEE:  proto.Stat_StatHasteRating,
	RATING_MOD_HASTE_RANGED: proto.Stat_StatHasteRating,
	RATING_MOD_HASTE_SPELL:  proto.Stat_StatHasteRating,
	RATING_MOD_AVOIDANCE:    -1,
	RATING_MOD_EXPERTISE:    proto.Stat_StatExpertiseRating,
	RATING_MOD_MASTERY:      proto.Stat_StatMasteryRating,
	RATING_MOD_PVP_POWER:    proto.Stat_StatPvpPowerRating,

	RATING_MOD_VERS_DAMAGE: -1,
	RATING_MOD_VERS_HEAL:   -1,
	RATING_MOD_VERS_MITIG:  -1,
}

type DbcClass struct {
	ProtoClass proto.Class
	ID         int
}

var Classes = []DbcClass{
	{proto.Class_ClassWarrior, 1},
	{proto.Class_ClassPaladin, 2},
	{proto.Class_ClassHunter, 3},
	{proto.Class_ClassRogue, 4},
	{proto.Class_ClassPriest, 5},
	{proto.Class_ClassDeathKnight, 6},
	{proto.Class_ClassShaman, 7},
	{proto.Class_ClassMage, 8},
	{proto.Class_ClassWarlock, 9},
	{proto.Class_ClassMonk, 10},
	{proto.Class_ClassDruid, 11},
}
