package database

import (
	"github.com/wowsims/cata/sim/core/proto"
)

var qualityToItemQualityMap = map[int]proto.ItemQuality{
	0: proto.ItemQuality_ItemQualityJunk,
	1: proto.ItemQuality_ItemQualityCommon,
	2: proto.ItemQuality_ItemQualityUncommon,
	3: proto.ItemQuality_ItemQualityRare,
	4: proto.ItemQuality_ItemQualityEpic,
	5: proto.ItemQuality_ItemQualityLegendary,
}

var inventoryTypeMap = map[int]string{
	0:  "None",
	1:  "Head",
	2:  "Neck",
	3:  "Shoulders",
	4:  "Shirt",
	5:  "Vest",
	6:  "Waist",
	7:  "Legs",
	8:  "Feet",
	9:  "Wrist",
	10: "Hands",
	11: "Ring",
	12: "Trinket",
	13: "One Hand",
	14: "Shield",
	15: "Bow",
	16: "Back",
	17: "Two Hand",
	18: "Bag",
	19: "Tabard",
	20: "Robe",
	21: "Main Hand",
	22: "Off Hand",
	23: "Held",
	24: "Ammo",
	25: "Ranged",
	26: "Ranged",
	27: "Ranged",
	28: "Ranged",
}

var inventoryTypeMapToItemType = map[int]proto.ItemType{
	0:  proto.ItemType_ItemTypeUnknown,
	1:  proto.ItemType_ItemTypeHead,
	2:  proto.ItemType_ItemTypeNeck,
	3:  proto.ItemType_ItemTypeShoulder,
	4:  proto.ItemType_ItemTypeUnknown,
	5:  proto.ItemType_ItemTypeChest,
	6:  proto.ItemType_ItemTypeWaist,
	7:  proto.ItemType_ItemTypeLegs,
	8:  proto.ItemType_ItemTypeFeet,
	9:  proto.ItemType_ItemTypeWrist,
	10: proto.ItemType_ItemTypeHands,
	11: proto.ItemType_ItemTypeFinger,
	12: proto.ItemType_ItemTypeTrinket,
	13: proto.ItemType_ItemTypeWeapon,
	14: proto.ItemType_ItemTypeWeapon,
	15: proto.ItemType_ItemTypeWeapon,
	16: proto.ItemType_ItemTypeBack,
	17: proto.ItemType_ItemTypeWeapon,
	18: proto.ItemType_ItemTypeUnknown,
	19: proto.ItemType_ItemTypeUnknown,
	20: proto.ItemType_ItemTypeChest,
	21: proto.ItemType_ItemTypeWeapon,
	22: proto.ItemType_ItemTypeWeapon,
	23: proto.ItemType_ItemTypeUnknown,
	24: proto.ItemType_ItemTypeUnknown,
	25: proto.ItemType_ItemTypeWeapon,
	26: proto.ItemType_ItemTypeWeapon,
	27: proto.ItemType_ItemTypeWeapon,
	28: proto.ItemType_ItemTypeWeapon,
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
	case 47: // SpellPenetration
		return proto.Stat_StatSpellPenetration, true
	case 35: // ResilienceRating
		return proto.Stat_StatResilienceRating, true
	case 56: // ArcaneResistance
		return proto.Stat_StatArcaneResistance, true
	case 51: // FireResistance
		return proto.Stat_StatFireResistance, true
	case 52: // FrostResistance
		return proto.Stat_StatFrostResistance, true
	case 55: // NatureResistance
		return proto.Stat_StatNatureResistance, true
	case 54: // ShadowResistance
		return proto.Stat_StatShadowResistance, true
	case 50: // ExtraArmor maps to BonusArmor (green armor)
		return proto.Stat_StatBonusArmor, true
	case 43: // ManaRegeneration
		return proto.Stat_StatMP5, true
	default:
		return 0, false
	}
}

var subClassToArmorType = map[string]proto.ArmorType{
	"Cloth":   proto.ArmorType_ArmorTypeCloth,
	"Leather": proto.ArmorType_ArmorTypeLeather,
	"Mail":    proto.ArmorType_ArmorTypeMail,
	"Plate":   proto.ArmorType_ArmorTypePlate,
	"":        proto.ArmorType_ArmorTypeUnknown,
}

type WeaponAndHand struct {
	Weapon proto.WeaponType
	Hand   proto.HandType
}

var subClassNameToWeaponAndHandType = map[string]WeaponAndHand{
	"One-Handed Axes":    {Weapon: proto.WeaponType_WeaponTypeAxe, Hand: proto.HandType_HandTypeOneHand},
	"Two-Handed Axes":    {Weapon: proto.WeaponType_WeaponTypeAxe, Hand: proto.HandType_HandTypeTwoHand},
	"Bows":               {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeUnknown},
	"Guns":               {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeUnknown},
	"One-Handed Maces":   {Weapon: proto.WeaponType_WeaponTypeMace, Hand: proto.HandType_HandTypeOneHand},
	"Two-Handed Maces":   {Weapon: proto.WeaponType_WeaponTypeMace, Hand: proto.HandType_HandTypeTwoHand},
	"Polearms":           {Weapon: proto.WeaponType_WeaponTypePolearm, Hand: proto.HandType_HandTypeTwoHand},
	"One-Handed Swords":  {Weapon: proto.WeaponType_WeaponTypeSword, Hand: proto.HandType_HandTypeOneHand},
	"Two-Handed Swords":  {Weapon: proto.WeaponType_WeaponTypeSword, Hand: proto.HandType_HandTypeTwoHand},
	"Staves":             {Weapon: proto.WeaponType_WeaponTypeStaff, Hand: proto.HandType_HandTypeTwoHand},
	"One-Handed Exotics": {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeOneHand},
	"Two-Handed Exotics": {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeTwoHand},
	"Fist Weapons":       {Weapon: proto.WeaponType_WeaponTypeFist, Hand: proto.HandType_HandTypeOneHand},
	"Daggers":            {Weapon: proto.WeaponType_WeaponTypeDagger, Hand: proto.HandType_HandTypeOneHand},
	"Thrown":             {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeUnknown},
	"Spears":             {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeTwoHand},
	"Crossbows":          {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeUnknown},
	"Wands":              {Weapon: proto.WeaponType_WeaponTypeUnknown, Hand: proto.HandType_HandTypeUnknown},
	"Fishing Poles":      {Weapon: proto.WeaponType_WeaponTypePolearm, Hand: proto.HandType_HandTypeTwoHand},
}

var subClassNameRoRangedWeaponType = map[string]proto.RangedWeaponType{
	"Crossbows": proto.RangedWeaponType_RangedWeaponTypeCrossbow,
	"Bows":      proto.RangedWeaponType_RangedWeaponTypeBow,
	"Guns":      proto.RangedWeaponType_RangedWeaponTypeGun,
	"Wands":     proto.RangedWeaponType_RangedWeaponTypeWand,
	"Relic":     proto.RangedWeaponType_RangedWeaponTypeRelic,
}
