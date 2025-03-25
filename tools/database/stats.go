package database

import (
	"fmt"
	"math"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func processStats(raw rawItemData, item *proto.UIItem) error {
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
	percent, err := parseIntArrayField(raw.statPercent, 10)
	if err != nil {
		fmt.Printf("Error parsing statPercent: %v\n", err)
		return err
	}
	bonusStats, err := parseIntArrayField(raw.bonusStat, 10)
	if err != nil {
		fmt.Printf("Error parsing bonusStat: %v\n", err)
		return err
	}

	alloc := CalcItemAllocation(item)
	rpp := 0
	if raw.overallQuality >= 4 {
		rpp = epic[alloc]
	} else if raw.overallQuality < 4 {
		rpp = superior[alloc]
	} else if raw.overallQuality <= 2 {
		rpp = good[alloc]
	}
	fmt.Println(item.Type, alloc)
	st := stats.Stats{}

	for i, statIndex := range bonusStats {
		if statIndex != -1 {
			if stat, ok := MapBonusStatIndexToStat(statIndex); ok {
				calculated := percent[i] * rpp
				value := math.Round(float64(calculated) / 10000)
				// Remap Armor stat to BonusArmor if needed idk
				if stat == proto.Stat_StatArmor {
					stat = proto.Stat_StatBonusArmor
				}
				st[stat] = value
				fmt.Println("STAT:", stat.String(), value, rpp, percent[i], alloc, raw.overallQuality)
			}
		}
	}
	item.Stats = st.ToProtoArray()

	return ParseStats(raw.id, raw.name, raw.invType, raw.itemLevel)
}

func ParseStats(id int, name string, invType int, itemLevel int) error {
	fmt.Printf("ParseStats - Item ID: %d: Stats parsing not implemented yet.\n", id)
	return nil
}
