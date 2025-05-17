package dbc

import (
	"math"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

var MAX_UPGRADE_LEVELS = []int{1, 2}

const UPGRADE_SYSTEM_ACTIVE = true

type Item struct {
	Id                     int
	Name                   string
	InventoryType          int
	ItemDelay              int
	OverallQuality         ItemQuality
	DmgVariance            float64
	ItemLevel              int
	ItemClass              int
	ItemSubClass           int
	StatAlloc              []float64
	BonusStat              []int
	SocketEnchantmentId    int
	Flags0                 ItemStaticFlags0
	Flags1                 ItemStaticFlags1
	Flags2                 ItemStaticFlags2
	Flags3                 ItemStaticFlags3
	FDID                   int
	ItemSetName            string
	ItemSetId              int
	ClassMask              int
	RaceMask               int
	QualityModifier        float64
	RandomSuffixOptions    []int32
	StatPercentageOfSocket []float64
	BonusAmountCalculated  []float64
	Sockets                []int
	SocketModifier         []float64 // Todo: Figure out if this is socket modifier in disguise or something else - I call it that for now.
}

func (item *Item) ToUIItem() *proto.UIItem {
	return item.ToScaledUIItem(item.ItemLevel)
}

func (item *Item) ToScaledUIItem(itemLevel int) *proto.UIItem {
	scalingProperties := make(map[int32]*proto.ScalingItemProperties)
	var weaponType, handType, rangedType = item.GetWeaponTypes()
	uiItem := &proto.UIItem{
		Type:                MapInventoryTypeToItemType[item.InventoryType],
		Quality:             item.OverallQuality.ToProto(),
		SetName:             item.ItemSetName,
		SetId:               int32(item.ItemSetId),
		Name:                item.Name,
		ClassAllowlist:      GetClassesFromClassMask(item.ClassMask),
		Id:                  int32(item.Id),
		RandomSuffixOptions: item.RandomSuffixOptions,
		RangedWeaponType:    rangedType,
		HandType:            handType,
		WeaponType:          weaponType,
		WeaponSpeed:         float64(item.ItemDelay) / 1000,
		GemSockets:          item.GetGemSlots(),
		SocketBonus:         item.GetGemBonus().ToProtoArray(),
	}

	item.ParseItemFlags(uiItem)

	if item.ItemClass == ITEM_CLASS_ARMOR {
		uiItem.ArmorType = MapArmorSubclassToArmorType[item.ItemSubClass]
	}

	if item.ItemSubClass == ITEM_SUBCLASS_ARMOR_SHIELD && uiItem.HandType == proto.HandType_HandTypeUnknown {
		uiItem.HandType = proto.HandType_HandTypeOffHand
		uiItem.WeaponType = proto.WeaponType_WeaponTypeShield
	}

	// Append base itemlevel stats
	scalingProperties[int32(proto.ItemLevelState_Base)] = &proto.ScalingItemProperties{
		WeaponDamageMin: item.WeaponDmgMin(item.ItemLevel),
		WeaponDamageMax: item.WeaponDmgMax(item.ItemLevel),
		Stats:           item.GetStats(item.ItemLevel).ToProtoMap(),
		RandPropPoints:  item.GetRandPropPoints(item.ItemLevel),
		Ilvl:            int32(item.ItemLevel),
	}

	// Amount of upgrade steps is defined in MAX_UPGRADE_LEVELS
	// In P2 of MoP it is expected to be 2 steps
	//
	if item.ItemLevel > 458 && UPGRADE_SYSTEM_ACTIVE {
		for _, upgradeLevel := range MAX_UPGRADE_LEVELS {
			upgradedIlvl := item.ItemLevel + item.UpgradeItemLevelBy(upgradeLevel)
			upgradeStep := proto.ItemLevelState(upgradeLevel)
			scalingProperties[int32(upgradeStep)] = &proto.ScalingItemProperties{
				WeaponDamageMin: item.WeaponDmgMin(upgradedIlvl),
				WeaponDamageMax: item.WeaponDmgMax(upgradedIlvl),
				Stats:           item.GetStats(upgradedIlvl).ToProtoMap(),
				RandPropPoints:  item.GetRandPropPoints(upgradedIlvl),
				Ilvl:            int32(upgradedIlvl),
			}
		}
	}

	if item.ItemLevel > core.MaxChallengeModeIlvl {
		scalingProperties[int32(proto.ItemLevelState_ChallengeMode)] = &proto.ScalingItemProperties{
			WeaponDamageMin: item.WeaponDmgMin(core.MaxChallengeModeIlvl),
			WeaponDamageMax: item.WeaponDmgMax(core.MaxChallengeModeIlvl),
			Stats:           item.GetStats(core.MaxChallengeModeIlvl).ToProtoMap(),
			RandPropPoints:  item.GetRandPropPoints(core.MaxChallengeModeIlvl),
			Ilvl:            core.MaxChallengeModeIlvl,
		}
	}
	uiItem.ScalingOptions = scalingProperties
	return uiItem
}

func (item *Item) GetMaxIlvl() int {
	if item.ItemLevel > 458 {
		return item.ItemLevel + item.UpgradeItemLevelBy(MAX_UPGRADE_LEVELS[len(MAX_UPGRADE_LEVELS)-1])
	}
	return item.ItemLevel
}

func (item *Item) ParseItemFlags(uiItem *proto.UIItem) {
	if item.Flags1.Has(HORDE_SPECIFIC) {
		uiItem.FactionRestriction = proto.UIItem_FACTION_RESTRICTION_HORDE_ONLY
	}

	if item.Flags1.Has(ALLIANCE_SPECIFIC) {
		uiItem.FactionRestriction = proto.UIItem_FACTION_RESTRICTION_ALLIANCE_ONLY
	}

	if item.Flags0.Has(UNIQUE_EQUIPPABLE) {
		uiItem.Unique = true
	}

	if item.Flags0.Has(HEROIC_TOOLTIP) {
		uiItem.Heroic = true
	}
}

func (item *Item) GetStats(itemLevel int) *stats.Stats {
	stats := &stats.Stats{}
	for i, alloc := range item.BonusStat {
		stat, success := MapBonusStatIndexToStat(alloc)
		if !success {
			// Skip this stat then
			continue
		}
		stats[stat] = item.GetScaledStat(i, itemLevel)
		if stat == proto.Stat_StatAttackPower {
			stats[proto.Stat_StatRangedAttackPower] = item.GetScaledStat(i, itemLevel) // Apply RAP as well. Might not be true for 1.12 idk
		}
	}

	armor := item.GetArmorValue(itemLevel)
	if armor > 0 {
		stats[proto.Stat_StatArmor] = float64(armor)

		if item.QualityModifier > 0 {
			stats[proto.Stat_StatBonusArmor] = item.QualityModifier
		}
	}
	return stats
}
func (item *Item) GetRandPropPoints(itemLevel int) int32 {
	suffixType := item.GetRandomSuffixType()
	randomProperty := GetDBC().RandomPropertiesByIlvl[itemLevel]
	if suffixType < 0 {
		return 0
	}
	return randomProperty[item.OverallQuality.ToProto()][suffixType]
}
func (item *Item) GetScaledStat(index int, itemLevel int) float64 {
	//Todo check if overflow array

	if itemLevel == item.ItemLevel {
		// Maybe just return it?
		return item.BonusAmountCalculated[index]
	}

	slotType := item.GetRandomSuffixType()
	itemBudget := 0.0

	if slotType != -1 && item.OverallQuality > 0 {

		randomProperty := GetDBC().RandomPropertiesByIlvl[itemLevel]
		itemBudget = float64(randomProperty[item.OverallQuality.ToProto()][slotType])

		if item.StatAlloc[index] > 0 && itemBudget > 0 {
			rawValue := math.Round(item.StatAlloc[index] * itemBudget * 0.0001)

			//Todo: Figure out if this does anything in MoP
			//Not used right now in Cata
			//socket_penalty := math.RoundNearby item.StatPercentageOfSocket[index] * SocketCost(itemLevel)
			return rawValue - item.SocketModifier[index] // Todo: Could this be a calculated socket penalty?
		} else {
			return math.Floor(item.BonusAmountCalculated[index] * item.ApproximateScaleCoeff(item.ItemLevel, itemLevel))
		}
	}
	return 0
}

func (item *Item) GetGemSlots() []proto.GemColor {
	sockets := []proto.GemColor{}
	for _, socketType := range item.Sockets {
		if socketType == 0 {
			continue
		}
		var gemType = MapSocketTypeToGemColor[socketType]
		sockets = append(sockets, gemType)
	}
	return sockets
}

func (item *Item) GetGemBonus() stats.Stats {
	stats := stats.Stats{}
	if item.SocketEnchantmentId == 0 {
		return stats
	}
	bonus := GetDBC().ItemStatEffects[item.SocketEnchantmentId]

	for i, effectStat := range bonus.EffectArg {
		if effectStat == 0 {
			continue
		}
		stat, success := MapBonusStatIndexToStat(effectStat)
		if !success {
			return stats
		}
		value := bonus.EffectPointsMin[i]
		stats[stat] = float64(value)
		//Todo: check if this is always true
		if stat == proto.Stat_StatAttackPower {
			stats[proto.Stat_StatRangedAttackPower] = float64(value)
		}
	}
	return stats
}

func (item *Item) WeaponDmgMin(itemLevel int) float64 {
	if itemLevel == 0 {
		itemLevel = item.ItemLevel
	}
	total := item.WeaponDps(itemLevel)*(float64(item.ItemDelay)/1000.0)*(1-item.DmgVariance/2) +
		(item.QualityModifier * (float64(item.ItemDelay) / 1000.0))
	if total < 0 {
		total = 1
	}
	return math.Floor(total)
}

func (item *Item) WeaponDmgMax(itemLevel int) float64 {
	if itemLevel == 0 {
		itemLevel = item.ItemLevel
	}
	total := item.WeaponDps(itemLevel)*(float64(item.ItemDelay)/1000.0)*(1+item.DmgVariance/2) +
		(item.QualityModifier * (float64(item.ItemDelay) / 1000.0))
	if total < 0 {
		total = 1
	}
	return math.Floor(total + 0.5)
}

func (item *Item) ApproximateScaleCoeff(currIlvl int, newIlvl int) float64 {
	if currIlvl == 0 || newIlvl == 0 {
		return 1.0
	}

	diff := float64(currIlvl - newIlvl)
	return 1.0 / math.Pow(1.15, diff/15.0)
}

func (item *Item) GetArmorValue(itemLevel int) int {
	if item.Id == 0 || item.OverallQuality > 5 {
		return 0
	}
	ilvl := 0
	if itemLevel > 0 {
		ilvl = itemLevel
	} else {
		ilvl = item.ItemLevel
	}

	if item.ItemClass == ITEM_CLASS_ARMOR && item.ItemSubClass == ITEM_SUBCLASS_ARMOR_SHIELD {
		return int(math.Floor(GetDBC().ItemArmorShield[ilvl].Quality[item.OverallQuality] + 0.5))
	}

	if item.ItemSubClass == ITEM_SUBCLASS_ARMOR_MISC || item.ItemSubClass > ITEM_SUBCLASS_ARMOR_PLATE {
		return 0
	}
	total_armor := 0.0
	quality := 0.0
	//	3688.5300292969 * 1.37000000477 * 0.15999999642
	armorModifier := GetDBC().ArmorLocation[item.InventoryType] //	0.15999999642 	3688.5300292969 	1.37000000477
	if item.InventoryType == INVTYPE_ROBE {
		armorModifier = GetDBC().ArmorLocation[INVTYPE_CHEST]
	}
	switch item.InventoryType {
	case INVTYPE_HEAD, INVTYPE_SHOULDERS, INVTYPE_CHEST, INVTYPE_WAIST, INVTYPE_LEGS, INVTYPE_FEET, INVTYPE_WRISTS, INVTYPE_HANDS, INVTYPE_CLOAK, INVTYPE_ROBE:
		switch item.ItemSubClass {
		case ITEM_SUBCLASS_ARMOR_CLOTH:
			total_armor = GetDBC().ItemArmorTotal[ilvl].Cloth
		case ITEM_SUBCLASS_ARMOR_LEATHER:
			total_armor = GetDBC().ItemArmorTotal[ilvl].Leather
		case ITEM_SUBCLASS_ARMOR_MAIL:
			total_armor = GetDBC().ItemArmorTotal[ilvl].Mail
		case ITEM_SUBCLASS_ARMOR_PLATE:
			total_armor = GetDBC().ItemArmorTotal[ilvl].Plate
		}
		quality = GetDBC().ItemArmorQuality[ilvl].Quality[item.OverallQuality]
	default:
		return 0
	}
	return int(math.Floor(total_armor*quality*armorModifier.Modifier[item.ItemSubClass-1] + 0.5))
}

func (item *Item) GetWeaponTypes() (proto.WeaponType, proto.HandType, proto.RangedWeaponType) {
	weaponType := proto.WeaponType_WeaponTypeUnknown
	rangedWeaponType := proto.RangedWeaponType_RangedWeaponTypeUnknown
	handType := proto.HandType_HandTypeUnknown

	switch item.ItemClass {
	case ITEM_CLASS_ARMOR:
		switch item.ItemSubClass {
		case ITEM_SUBCLASS_ARMOR_MISC:
			if item.InventoryType == INVTYPE_HOLDABLE {

				handType = proto.HandType_HandTypeOffHand
				weaponType = proto.WeaponType_WeaponTypeOffHand
			}
		}
	case ITEM_CLASS_WEAPON:
		switch item.ItemSubClass {
		case ITEM_SUBCLASS_WEAPON_BOW:
			rangedWeaponType = proto.RangedWeaponType_RangedWeaponTypeBow
		case ITEM_SUBCLASS_WEAPON_CROSSBOW:
			rangedWeaponType = proto.RangedWeaponType_RangedWeaponTypeCrossbow
		case ITEM_SUBCLASS_WEAPON_GUN:
			rangedWeaponType = proto.RangedWeaponType_RangedWeaponTypeGun
		case ITEM_SUBCLASS_WEAPON_WAND:
			rangedWeaponType = proto.RangedWeaponType_RangedWeaponTypeWand
		case ITEM_SUBCLASS_WEAPON_THROWN:
			rangedWeaponType = proto.RangedWeaponType_RangedWeaponTypeThrown
		case ITEM_SUBCLASS_WEAPON_AXE2, ITEM_SUBCLASS_WEAPON_FISHING_POLE, ITEM_SUBCLASS_WEAPON_SPEAR, ITEM_SUBCLASS_WEAPON_STAFF, ITEM_SUBCLASS_WEAPON_POLEARM, ITEM_SUBCLASS_WEAPON_MACE2,
			ITEM_SUBCLASS_WEAPON_EXOTIC2, ITEM_SUBCLASS_WEAPON_SWORD2:
			handType = proto.HandType_HandTypeTwoHand
			weaponType = MapWeaponSubClassToWeaponType[item.ItemSubClass]
		case ITEM_SUBCLASS_WEAPON_EXOTIC, ITEM_SUBCLASS_WEAPON_MACE, ITEM_SUBCLASS_WEAPON_WARGLAIVE, ITEM_SUBCLASS_WEAPON_FIST, ITEM_SUBCLASS_WEAPON_SWORD,
			ITEM_SUBCLASS_WEAPON_DAGGER, ITEM_SUBCLASS_WEAPON_AXE:
			if item.InventoryType == INVTYPE_WEAPON {
				handType = proto.HandType_HandTypeOneHand
			} else if item.InventoryType == INVTYPE_WEAPONOFFHAND {
				handType = proto.HandType_HandTypeOffHand
			} else {
				handType = proto.HandType_HandTypeMainHand
			}
			weaponType = MapWeaponSubClassToWeaponType[item.ItemSubClass]
		}
	}
	return weaponType, handType, rangedWeaponType
}

func (item *Item) GetRandomSuffixType() int {
	switch item.ItemClass {
	case ITEM_CLASS_WEAPON:
		switch item.ItemSubClass {
		case ITEM_SUBCLASS_WEAPON_AXE2,
			ITEM_SUBCLASS_WEAPON_MACE2,
			ITEM_SUBCLASS_WEAPON_POLEARM,
			ITEM_SUBCLASS_WEAPON_SWORD2,
			ITEM_SUBCLASS_WEAPON_STAFF,
			ITEM_SUBCLASS_WEAPON_GUN,
			ITEM_SUBCLASS_WEAPON_BOW,
			ITEM_SUBCLASS_WEAPON_CROSSBOW:
			return 0

		case ITEM_SUBCLASS_WEAPON_THROWN:
			return 4

		default:
			return 3
		}

	case ITEM_CLASS_ARMOR:
		switch item.InventoryType {
		case INVTYPE_HEAD,
			INVTYPE_CHEST,
			INVTYPE_LEGS,
			INVTYPE_ROBE:
			return 0

		case INVTYPE_SHOULDERS,
			INVTYPE_WAIST,
			INVTYPE_FEET,
			INVTYPE_HANDS,
			INVTYPE_TRINKET:
			return 1

		case INVTYPE_NECK,
			INVTYPE_WEAPONOFFHAND,
			INVTYPE_HOLDABLE,
			INVTYPE_FINGER,
			INVTYPE_CLOAK,
			INVTYPE_WRISTS,
			INVTYPE_SHIELD:
			return 2

		default:
			return -1
		}

	default:
		return -1
	}
}

func (item *Item) UpgradeItemLevelBy(upgradeLevel int) int {
	if item.OverallQuality == 3 {
		return upgradeLevel * 8
	}
	if item.OverallQuality == 4 {
		return upgradeLevel * 4
	}
	if item.OverallQuality == 5 {
		return upgradeLevel * 4
	}
	return 0
}
