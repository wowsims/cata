package database

import (
	"fmt"
	"math"

	"github.com/wowsims/cata/sim/core/proto"
)

func processGemStats(raw RawGem, gem *proto.UIGem) error {
	for i, effect := range raw.Effect {
		if effect == 5 {
			stat, err := MapBonusStatIndexToStat(raw.StatList[i])
			if err != true {
				return fmt.Errorf("Error mapping bonus stat to stat")
			}
			amount := raw.StatBonus[i]
			gem.Stats[stat] = float64(amount)
		}
	}
	return nil
}

func processStats(raw RawItemData, item *proto.UIItem) error {
	epic, err := parseIntArrayField(raw.rppEpic, 5)
	if err != nil {
		fmt.Printf("Error parsing rppEpic: %v\n", err)
		return err
	}
	superior, err := parseIntArrayField(raw.rppSuperior, 5)
	if err != nil {
		fmt.Printf("Error parsing rppSuperior: %v\n", err)
		return err
	}
	good, err := parseIntArrayField(raw.rppGood, 5)
	if err != nil {
		fmt.Printf("Error parsing rppGood: %v\n", err)
		return err
	}
	percent, err := parseIntArrayField(raw.statValue, 10)
	if err != nil {
		fmt.Printf("Error parsing percent: %v\n", err)
		return err
	}
	bonusStats, err := parseIntArrayField(raw.bonusStat, 10)
	if err != nil {
		fmt.Printf("Error parsing bonusStat: %v\n", err)
		return err
	}
	statMods, err := parseIntArrayField(raw.statPercentEditor, 10)
	if err != nil {
		fmt.Printf("Error parsing statMods: %v\n", err)
		return err
	}
	alloc := CalcItemAllocation(item)
	if alloc != -1 {
		rpp := 0
		if raw.overallQuality >= 4 {
			rpp = epic[alloc]
		} else if raw.overallQuality < 4 {
			rpp = superior[alloc]
		} else if raw.overallQuality <= 2 {
			rpp = good[alloc]
		}

		for i, statIndex := range bonusStats {
			if statIndex != -1 {
				if stat, ok := MapBonusStatIndexToStat(statIndex); ok {
					calculated := percent[i] * rpp

					var statMod = statMods[i]
					value := math.Round(float64(calculated)/10000) - float64(statMod)
					// Remap Armor stat to BonusArmor if needed idk
					if stat == proto.Stat_StatArmor {
						stat = proto.Stat_StatBonusArmor
					}
					item.Stats[stat] = value
				}
			}
		}
		if raw.qualityModifier > 0 {
			item.Stats[proto.Stat_StatBonusArmor] = raw.qualityModifier
		}
	}

	return ParseStats(raw.id, raw.name, raw.invType, raw.itemLevel)
}

func ParseStats(id int, name string, invType int, itemLevel int) error {
	//fmt.Printf("ParseStats - Item ID: %d: Stats parsing not implemented yet.\n", id)
	return nil
}
