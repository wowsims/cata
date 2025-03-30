package database

import (
	"math"
	"strings"

	"github.com/wowsims/cata/sim/core/proto"
)

func processWeaponDamage(helper *DBHelper, raw RawItemData, item *proto.UIItem) error {
	tableSuffix := invTypeToTableNameSuffix(raw.invType)
	tableName := "ItemDamage" + tableSuffix
	if raw.overallQuality == 7 { // Skip heirlooms for now lazy
		return nil
	}

	if raw.flags1&0x200 != 0 {
		tableName += "Caster"
	}
	qualityValue := ItemDamageByTableAndItemLevel[tableName][raw.itemLevel].Quality[raw.overallQuality]
	multiplier := float64(raw.itemDelay) / 1000.0
	baseDamage := qualityValue * multiplier
	qualityAdjustment := raw.qualityModifier * multiplier
	calcMinDamage := baseDamage*(1-raw.dmgVariance/2) + qualityAdjustment
	calcMaxDamage := baseDamage*(1+raw.dmgVariance/2) + qualityAdjustment

	// For now we log the damage calculations
	// fmt.Printf("processWeaponDamage - Item ID: %d\n", raw.id)
	// fmt.Printf("Name: %s, Ilvl: %s, QualityModifier: %.2f\n", raw.name, raw.itemLevel, raw.qualityModifier)
	// fmt.Printf("Delay: %d (Multiplier: %.2f), Quality: %d (Value: %.2f)\n", raw.itemDelay, multiplier, raw.overallQuality, qualityValue)
	// fmt.Printf("Damage Variance: %.2f\n", raw.dmgVariance)
	// fmt.Printf("Calculated Base Damage: %.2f, Min: %.2f, Max: %.2f\n", baseDamage, calcMinDamage, calcMaxDamage)
	// fmt.Printf("DB Min: %f, DB Max: %f\n", raw.dbMinDamage, raw.dbMaxDamage)
	// fmt.Printf("Multi = %f * %f\n", qualityValue, multiplier)
	// fmt.Printf("BaseDam = %f + %f\n", raw.qualityModifier, multiplier)
	// fmt.Printf("MinDam = %f * 1 (1-%f/2) + %f\n", baseDamage, raw.dmgVariance, qualityAdjustment)
	// fmt.Printf("MaxDam = %f * 1 (1+%f/2) + %f\n", baseDamage, raw.dmgVariance, qualityAdjustment)
	// fmt.Println("-----------")
	item.WeaponDamageMax = math.Round(calcMaxDamage)
	item.WeaponDamageMin = math.Round(calcMinDamage)

	return nil
}

func invTypeToTableNameSuffix(invType int) string {
	invTypeStr, ok := inventoryTypeMap[invType]
	if !ok {
		invTypeStr = "Unknown"
	}

	clean := strings.ReplaceAll(invTypeStr, " ", "")
	switch clean {
	case "Bow", "Crossbow", "Gun":
		return "Ranged"
	case "MainHand", "OffHand":
		return "OneHand"
	default:
		return clean
	}
}
