package database

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/wowsims/cata/sim/core/proto"
)

func processWeaponDamage(helper *DBHelper, raw RawItemData, item *proto.UIItem) error {
	invTypeStr, ok := inventoryTypeMap[raw.invType]
	if !ok {
		invTypeStr = "Unknown"
	}

	tableSuffix := invTypeToTableNameSuffix(raw.invType)
	tableName := "ItemDamage" + tableSuffix
	colName := fmt.Sprintf("Quality_%d", raw.overallQuality)

	qualityQuery := fmt.Sprintf("SELECT %s FROM %s WHERE ItemLevel = ?", colName, tableName)
	var qualityValue float64
	if err := helper.db.QueryRow(qualityQuery, raw.itemLevel).Scan(&qualityValue); err != nil {
		log.Printf("failed to query quality from %s: %v", tableName, err)
		return err
	}

	multiplier := float64(raw.itemDelay) / 1000.0
	baseDamage := qualityValue * multiplier
	calcMinDamage := baseDamage * (1 - raw.dmgVariance/2)
	calcMaxDamage := baseDamage * (1 + raw.dmgVariance/2)

	// For now we log the damage calculations
	fmt.Printf("processWeaponDamage - Item ID: %d\n", raw.id)
	fmt.Printf("Name: %s, InvType: %s, ItemLevel: %d\n", raw.name, invTypeStr, raw.itemLevel)
	fmt.Printf("Delay: %d (Multiplier: %.2f), Quality: %d (Value: %.2f)\n", raw.itemDelay, multiplier, raw.overallQuality, qualityValue)
	fmt.Printf("Damage Variance: %.2f\n", raw.dmgVariance)
	fmt.Printf("Calculated Base Damage: %.2f, Min: %.2f, Max: %.2f\n", baseDamage, calcMinDamage, calcMaxDamage)
	fmt.Printf("DB Min: %s, DB Max: %s\n", raw.dbMinDamage, raw.dbMaxDamage)
	fmt.Println("-----------")
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
	default:
		return clean
	}
}
