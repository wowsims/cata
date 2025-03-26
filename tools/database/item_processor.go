package database

import (
	"fmt"
	"math"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func RawItemToUIItem(helper *DBHelper, raw RawItemData) (*proto.UIItem, error) {
	item := &proto.UIItem{
		Type:    inventoryTypeMapToItemType[raw.invType],
		Quality: qualityToItemQualityMap[raw.overallQuality],
		Stats:   stats.Stats{}.ToProtoArray(),
		SetName: raw.itemSetName,
		SetId:   int32(raw.itemSetId),
		Name:    raw.name,

		Id:   int32(raw.id),
		Ilvl: int32(raw.itemLevel),
	}

	if raw.flags0.Has(UniqueEquipped) {
		item.Unique = true
	}
	if raw.flags0.Has(HeroicItem) {
		item.Heroic = true
	}
	if item.Type == proto.ItemType_ItemTypeWeapon || item.Type == proto.ItemType_ItemTypeRanged {
		weaponType, handType, rangedType := determineWeaponTypes(raw.itemSubClassName, raw.invType)
		if weaponType != proto.WeaponType_WeaponTypeShield {
			if err := processWeaponDamage(helper, raw, item); err != nil {
				fmt.Printf("processWeaponDamage error for item %d: %v\n", raw.id, err)
			}
			item.WeaponType = weaponType
			item.HandType = handType
			item.RangedWeaponType = rangedType

			item.WeaponSpeed = float64(raw.itemDelay) / 1000.0
		}
	}

	if raw.itemClassName == "Armor" {
		item.ArmorType = subClassToArmorType[raw.itemSubClassName]
		applyArmorValue(item, raw)
	}

	if err := processStats(raw, item); err != nil {
		fmt.Printf("processStats error for item %d: %v\n", raw.id, err)
	}

	socketTypes, err := parseIntArrayField(raw.socketTypes, 3)
	if err != nil {
		fmt.Printf("Error parsing socketTypes: %v\n", err)
	}
	for _, socketType := range socketTypes {
		if socketType == 0 {
			continue
		}
		var gemType = SocketTypeToGemColorMap[socketType]
		item.GemSockets = append(item.GemSockets, gemType)
	}

	//Orm-style
	LoadItemStatEffects(helper)

	var gemBonus = ItemStatEffectById[raw.socketEnchantmentId]
	//since its a socket bonus we know it should be straight forward to use min value?

	stats := &stats.Stats{}
	for i, effectStat := range gemBonus.EffectArg {
		if effectStat == 0 {
			continue
		}
		stat, err := MapBonusStatIndexToStat(effectStat)

		if err == false {
			fmt.Println("Error parsing statValue: \n", err, effectStat)
			fmt.Printf("Error parsing statValue: %v\n", err, effectStat)
		}
		value := gemBonus.EffectPointsMin[i]

		stats[stat] = float64(value)

	}

	item.SocketBonus = stats.ToProtoArray()

	return item, nil
}

func determineWeaponTypes(subClassName string, invType int) (proto.WeaponType, proto.HandType, proto.RangedWeaponType) {
	weaponType := proto.WeaponType_WeaponTypeUnknown
	handType := proto.HandType_HandTypeUnknown
	rangedType := proto.RangedWeaponType_RangedWeaponTypeUnknown
	switch invType {
	case 15, 26:
		if rt, ok := subClassNameRoRangedWeaponType[subClassName]; ok {
			rangedType = rt
		}
		break
	default:
		if w, ok := subClassNameToWeaponAndHandType[subClassName]; ok {
			weaponType = w.Weapon
			handType = w.Hand
		}
		break
	}
	return weaponType, handType, rangedType
}

func applyArmorValue(item *proto.UIItem, raw RawItemData) {
	var armorValue float64
	if raw.itemSubClassName == "Shields" {
		shieldArmor, err := parseIntArrayField(raw.shieldArmorValues, 7)
		if err != nil {
			fmt.Printf("Error parsing shieldArmor: %v\n", err)
		}
		armorValue = float64(shieldArmor[raw.overallQuality])
	}
	switch item.ArmorType {
	case proto.ArmorType_ArmorTypeCloth:
		armorValue = raw.clothArmorValue
	case proto.ArmorType_ArmorTypeLeather:
		armorValue = raw.leatherArmorValue
	case proto.ArmorType_ArmorTypeMail:
		armorValue = raw.mailArmorValue
	case proto.ArmorType_ArmorTypePlate:
		armorValue = raw.plateArmorValue
	}

	item.Stats[proto.Stat_StatArmor] = math.Round(armorValue)
}

func CalcItemAllocation(item *proto.UIItem) int {
	idx := -1
	switch item.Type {
	case proto.ItemType_ItemTypeHead, proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeLegs:
		idx = 0
		break
	case proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeWaist, proto.ItemType_ItemTypeFeet, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeTrinket:
		idx = 1
		break
	case proto.ItemType_ItemTypeNeck, proto.ItemType_ItemTypeWrist, proto.ItemType_ItemTypeFinger, proto.ItemType_ItemTypeBack:
		idx = 2
		break
	case proto.ItemType_ItemTypeRanged:
		switch item.RangedWeaponType {
		case proto.RangedWeaponType_RangedWeaponTypeBow, proto.RangedWeaponType_RangedWeaponTypeCrossbow, proto.RangedWeaponType_RangedWeaponTypeGun, proto.RangedWeaponType_RangedWeaponTypeThrown, proto.RangedWeaponType_RangedWeaponTypeWand:
			idx = 4
			break
		}
	case (proto.ItemType_ItemTypeWeapon):
		switch item.WeaponType {
		case proto.WeaponType_WeaponTypeOffHand, proto.WeaponType_WeaponTypeShield:
			idx = 2
			break
		default:
			if item.Type == proto.ItemType_ItemTypeRanged {
				break
			}
			if item.HandType == proto.HandType_HandTypeTwoHand {
				idx = 0
			} else {
				idx = 3
			}
		}
		break
	}
	return idx
}
