package database

import (
	"database/sql"
	"fmt"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type rawItemData struct {
	id                int
	name              string
	invType           int
	itemDelay         int
	overallQuality    int
	dmgVariance       float64
	dbMinDamage       string
	dbMaxDamage       string
	itemLevel         int
	itemClassName     string
	itemSubClassName  string
	rppEpic           string
	rppSuperior       string
	rppGood           string
	statPercent       string
	bonusStat         string
	clothArmorValue   float64
	leatherArmorValue float64
	mailArmorValue    float64
	plateArmorValue   float64
	armorLocID        int
}

func processItemRow(helper *DBHelper, rows *sql.Rows) (*proto.UIItem, error) {
	var raw rawItemData

	if err := rows.Scan(
		&raw.id, &raw.name, &raw.invType, &raw.itemDelay, &raw.overallQuality, &raw.dmgVariance,
		&raw.dbMinDamage, &raw.dbMaxDamage, &raw.itemLevel, &raw.itemClassName, &raw.itemSubClassName,
		&raw.rppEpic, &raw.rppSuperior, &raw.rppGood, &raw.statPercent, &raw.bonusStat,
		&raw.clothArmorValue, &raw.leatherArmorValue, &raw.mailArmorValue, &raw.plateArmorValue,
		&raw.armorLocID,
	); err != nil {
		return nil, err
	}

	item := &proto.UIItem{
		Type:    inventoryTypeMapToItemType[raw.invType],
		Quality: qualityToItemQualityMap[raw.overallQuality],
	}
	item.Name = raw.name
	item.Id = int32(raw.id)
	item.Ilvl = int32(raw.itemLevel)

	if isWeapon(raw.invType) {
		if err := processWeaponDamage(helper, raw, item); err != nil {
			fmt.Printf("processWeaponDamage error for item %d: %v\n", raw.id, err)
		}
		weaponType, handType, rangedType := determineWeaponTypes(raw.itemSubClassName, raw.invType)
		if raw.invType != 15 { // not a ranged weapon
			item.WeaponType = weaponType
			item.HandType = handType
		} else {
			item.RangedWeaponType = rangedType
		}
		item.WeaponSpeed = float64(raw.itemDelay) / 1000.0
	}

	if raw.itemClassName == "Armor" {
		item.ArmorType = subClassToArmorType[raw.itemSubClassName]
		applyArmorValue(item, raw)
	}

	if err := processStats(raw, item); err != nil {
		fmt.Printf("processStats error for item %d: %v\n", raw.id, err)
	}

	// processGemSlots(raw, item)

	return item, nil
}

func determineWeaponTypes(subClassName string, invType int) (proto.WeaponType, proto.HandType, proto.RangedWeaponType) {
	weaponType := proto.WeaponType_WeaponTypeUnknown
	handType := proto.HandType_HandTypeUnknown
	rangedType := proto.RangedWeaponType_RangedWeaponTypeUnknown

	if invType != 15 { // non-ranged weapon
		if w, ok := subClassNameToWeaponAndHandType[subClassName]; ok {
			weaponType = w.Weapon
			handType = w.Hand
		}
	} else { // ranged weapon
		if rt, ok := subClassNameRoRangedWeaponType[subClassName]; ok {
			rangedType = rt
		}
	}
	return weaponType, handType, rangedType
}

func applyArmorValue(item *proto.UIItem, raw rawItemData) {
	var armorValue float64
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
	st := stats.Stats{}
	st[proto.Stat_StatArmor] = armorValue
	item.Stats = st.ToProtoArray()
}

func isWeapon(invType int) bool {
	switch invType {
	case 13, 17, 21, 22, 15, 25, 26, 27:
		return true
	default:
		return false
	}
}

func CalcItemAllocation(item *proto.UIItem) int {
	idx := -1
	switch item.Type {
	case proto.ItemType_ItemTypeHead, proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeLegs:
		idx = 0
	case proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeWaist, proto.ItemType_ItemTypeFeet, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeTrinket:
		idx = 1
	case proto.ItemType_ItemTypeNeck, proto.ItemType_ItemTypeWrist, proto.ItemType_ItemTypeFinger, proto.ItemType_ItemTypeBack:
		idx = 2
	case proto.ItemType_ItemTypeRanged:
		switch item.RangedWeaponType {
		case proto.RangedWeaponType_RangedWeaponTypeBow, proto.RangedWeaponType_RangedWeaponTypeCrossbow, proto.RangedWeaponType_RangedWeaponTypeGun, proto.RangedWeaponType_RangedWeaponTypeThrown, proto.RangedWeaponType_RangedWeaponTypeWand:
			idx = 4
		}
	case proto.ItemType_ItemTypeWeapon:
		switch item.WeaponType {
		case proto.WeaponType_WeaponTypeOffHand, proto.WeaponType_WeaponTypeShield:
			idx = 2
		default:
			if item.HandType == proto.HandType_HandTypeTwoHand {
				idx = 0
			} else {
				idx = 3
			}
		}
	}
	return idx
}
