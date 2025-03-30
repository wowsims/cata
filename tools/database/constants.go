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
	26: "Wand",
	27: "Ranged",
	28: "Ranged",
}
var SocketTypeToGemColorMap = map[int]proto.GemColor{
	0: proto.GemColor_GemColorUnknown,
	1: proto.GemColor_GemColorMeta,
	2: proto.GemColor_GemColorRed,
	3: proto.GemColor_GemColorYellow,
	4: proto.GemColor_GemColorBlue,
	5: proto.GemColor_GemColorUnknown,
	6: proto.GemColor_GemColorCogwheel,
	7: proto.GemColor_GemColorPrismatic,
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
	15: proto.ItemType_ItemTypeRanged,
	16: proto.ItemType_ItemTypeBack,
	17: proto.ItemType_ItemTypeWeapon,
	18: proto.ItemType_ItemTypeUnknown,
	19: proto.ItemType_ItemTypeUnknown,
	20: proto.ItemType_ItemTypeChest,
	21: proto.ItemType_ItemTypeWeapon,
	22: proto.ItemType_ItemTypeWeapon,
	23: proto.ItemType_ItemTypeWeapon,
	24: proto.ItemType_ItemTypeUnknown,
	25: proto.ItemType_ItemTypeWeapon,
	26: proto.ItemType_ItemTypeRanged,
	27: proto.ItemType_ItemTypeWeapon,
	28: proto.ItemType_ItemTypeWeapon,
}

type InventoryType int64

// Define each inventory slot as a bit flag.
const (
	InvTypeNonEquip       InventoryType = 0
	InvTypeHead           InventoryType = 1 << 1  // 2
	InvTypeNeck           InventoryType = 1 << 2  // 4
	InvTypeShoulder       InventoryType = 1 << 3  // 8
	InvTypeShirt          InventoryType = 1 << 4  // 16
	InvTypeChest          InventoryType = 1 << 5  // 32
	InvTypeWaist          InventoryType = 1 << 6  // 64
	InvTypeLegs           InventoryType = 1 << 7  // 128
	InvTypeFeet           InventoryType = 1 << 8  // 256
	InvTypeWrist          InventoryType = 1 << 9  // 512
	InvTypeHand           InventoryType = 1 << 10 // 1024
	InvTypeFinger         InventoryType = 1 << 11 // 2048
	InvTypeTrinket        InventoryType = 1 << 12 // 4096
	InvTypeWeapon         InventoryType = 1 << 13 // 8192
	InvTypeShield         InventoryType = 1 << 14 // 16384
	InvTypeRanged         InventoryType = 1 << 15 // 32768
	InvTypeCloak          InventoryType = 1 << 16 // 65536
	InvType2HWeapon       InventoryType = 1 << 17 // 131072
	InvTypeBag            InventoryType = 1 << 18 // 262144
	InvTypeTabard         InventoryType = 1 << 19 // 524288
	InvTypeRobe           InventoryType = 1 << 20 // 1048576
	InvTypeMainHand       InventoryType = 1 << 21 // 2097152
	InvTypeOffHand        InventoryType = 1 << 22 // 4194304
	InvTypeHoldable       InventoryType = 1 << 23 // 8388608
	InvTypeAmmo           InventoryType = 1 << 24 // 16777216
	InvTypeThrown         InventoryType = 1 << 25 // 33554432
	InvTypeRangedRight    InventoryType = 1 << 26 // 67108864
	InvTypeQuiver         InventoryType = 1 << 27 // 134217728
	InvTypeRelic          InventoryType = 1 << 28 // 268435456
	InvTypeProfessionTool InventoryType = 1 << 29 // 536870912
	InvTypeProfessionGear InventoryType = 1 << 30 // 1073741824
	// … additional entries could follow if needed.
)

type ItemSubClass int

// Define each item subclass as a bit flag (only those with a name).
const (
	OneHandedAxes    ItemSubClass = 1 << 0  // 1    from "One-Handed Axes" (SubClassID 0)
	TwoHandedAxes    ItemSubClass = 1 << 1  // 2    from "Two-Handed Axes" (SubClassID 1)
	Bows             ItemSubClass = 1 << 2  // 4    from "Bows" (SubClassID 2)
	Guns             ItemSubClass = 1 << 3  // 8    from "Guns" (SubClassID 3)
	OneHandedMaces   ItemSubClass = 1 << 4  // 16   from "One-Handed Maces" (SubClassID 4)
	TwoHandedMaces   ItemSubClass = 1 << 5  // 32   from "Two-Handed Maces" (SubClassID 5)
	Polearms         ItemSubClass = 1 << 6  // 64   from "Polearms" (SubClassID 6)
	OneHandedSwords  ItemSubClass = 1 << 7  // 128  from "One-Handed Swords" (SubClassID 7)
	TwoHandedSwords  ItemSubClass = 1 << 8  // 256  from "Two-Handed Swords" (SubClassID 8)
	Staves           ItemSubClass = 1 << 10 // 1024 from "Staves" (SubClassID 10)
	OneHandedExotics ItemSubClass = 1 << 11 // 2048 from "One-Handed Exotics" (SubClassID 11)
	TwoHandedExotics ItemSubClass = 1 << 12 // 4096 from "Two-Handed Exotics" (SubClassID 12)
	FistWeapons      ItemSubClass = 1 << 13 // 8192 from "Fist Weapons" (SubClassID 13)
	Daggers          ItemSubClass = 1 << 15 // 32768 from "Daggers" (SubClassID 15)
)

var SubclassNames = map[ItemSubClass]string{
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

type EnchantType struct {
	ItemType   proto.ItemType
	WeaponType proto.WeaponType
}

// SELECT
//
//	 se.ID as effectId,
//	 sn.Name_lang as name,
//	 se.SpellID as spellId,
//	 ie.ParentItemID as ItemId,
//	 sie.Field_1_15_3_55112_014 as professionId,
//	 sie.Effect as Effect,
//	 sie.EffectPointsMax as EffectPoints,
//	 sie.EffectArg as EffectArgs,
//	 sie.ID,
//	 CASE
//		WHEN sei.EquippedItemClass = 4 THEN 1
//		ELSE 0
//
// END AS isWeaponEnchant,
//
//	sei.EquippedItemInvTypes as InvTypes,
//	sei.EquippedItemSubclass
//
// FROM SpellEffect se
// JOIN Spell s ON se.SpellID = s.ID
// JOIN SpellName sn ON se.SpellID = sn.ID
// JOIN SpellItemEnchantment sie ON se.EffectMiscValue_0 = sie.ID
// LEFT JOIN ItemEffect ie ON se.SpellID = ie.SpellID
// LEFT JOIN SpellEquippedItems sei ON se.SpellId = sei.SpellID
// WHERE se.Effect = 53
var inventoryNames = map[InventoryType]EnchantType{
	InvTypeHead:     {ItemType: proto.ItemType_ItemTypeHead, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeNeck:     {ItemType: proto.ItemType_ItemTypeNeck, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeShoulder: {ItemType: proto.ItemType_ItemTypeShoulder, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeChest:    {ItemType: proto.ItemType_ItemTypeChest, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeWaist:    {ItemType: proto.ItemType_ItemTypeWaist, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeLegs:     {ItemType: proto.ItemType_ItemTypeLegs, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeFeet:     {ItemType: proto.ItemType_ItemTypeFeet, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeWrist:    {ItemType: proto.ItemType_ItemTypeWrist, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeHand:     {ItemType: proto.ItemType_ItemTypeHands, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeFinger:   {ItemType: proto.ItemType_ItemTypeFinger, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeTrinket:  {ItemType: proto.ItemType_ItemTypeTrinket, WeaponType: proto.WeaponType_WeaponTypeUnknown},

	InvTypeWeapon:   {ItemType: proto.ItemType_ItemTypeWeapon, WeaponType: proto.WeaponType_WeaponTypeUnknown}, // One-Hand
	InvTypeShield:   {ItemType: proto.ItemType_ItemTypeWeapon, WeaponType: proto.WeaponType_WeaponTypeShield},  // Off Hand
	InvTypeRanged:   {ItemType: proto.ItemType_ItemTypeRanged, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvTypeCloak:    {ItemType: proto.ItemType_ItemTypeBack, WeaponType: proto.WeaponType_WeaponTypeUnknown},
	InvType2HWeapon: {ItemType: proto.ItemType_ItemTypeWeapon, WeaponType: proto.WeaponType_WeaponTypeUnknown},
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

var ProfessionIdToProfession = map[int]proto.Profession{
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
	"Shields":            {Weapon: proto.WeaponType_WeaponTypeShield, Hand: proto.HandType_HandTypeOffHand},
}

var subClassNameRoRangedWeaponType = map[string]proto.RangedWeaponType{
	"Crossbows": proto.RangedWeaponType_RangedWeaponTypeCrossbow,
	"Bows":      proto.RangedWeaponType_RangedWeaponTypeBow,
	"Guns":      proto.RangedWeaponType_RangedWeaponTypeGun,
	"Wands":     proto.RangedWeaponType_RangedWeaponTypeWand,
	"Relic":     proto.RangedWeaponType_RangedWeaponTypeRelic,
}

type ItemFlags int

const (
	Unknown1       ItemFlags = 1 << iota // 1
	ConjuredItem                         // 2
	OpenableItem                         // 4
	HeroicItem                           // 8
	DeprecatedItem                       // 16
	Totem                                // 32
	Spelltriggerer                       // 64
	Unknown3                             // 128
	Wand                                 // 256
	Wrap                                 // 512
	Producer                             // 1024
	MultiLoot                            // 2048
	BGItem                               // 4096
	Charter                              // 8192
	ReadableItem                         // 16384
	PvPItem                              // 32768
	Duration                             // 65536
	Unknown4                             // 131072
	Prospectable                         // 262144
	UniqueEquipped                       // 524288
	Unknown5                             // 1048576
	Unknown6                             // 2097152
	ThrowingWeapon                       // 4194304
	Unknown7                             // 8388608
)

type Db2Class struct {
	protoClass proto.Class
	ID         int
}

var classes = []Db2Class{
	{proto.Class_ClassWarrior, 1},
	{proto.Class_ClassPaladin, 2},
	{proto.Class_ClassHunter, 3},
	{proto.Class_ClassRogue, 4},
	{proto.Class_ClassPriest, 5},
	{proto.Class_ClassDeathKnight, 6},
	{proto.Class_ClassShaman, 7},
	{proto.Class_ClassMage, 8},
	{proto.Class_ClassWarlock, 9},
	{proto.Class_ClassDruid, 11},
}

type GemType int

const (
	Meta   GemType = 0x1
	Red    GemType = 0x2
	Yellow GemType = 0x4
	Blue   GemType = 0x8
	// Combined colors:
	Orange    GemType = Red | Yellow        // 0x6
	Purple    GemType = Red | Blue          // 0xa
	Green     GemType = Yellow | Blue       // 0xc
	Prismatic GemType = Red | Yellow | Blue // 0xe
	Cogwheel  GemType = 0x20
)
