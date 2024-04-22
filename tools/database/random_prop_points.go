package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/tailscale/hujson"
	"github.com/wowsims/cata/sim/core/proto"
)

type RandomPropAllocation struct {
	Epic0     int32 `json:"Epic_0"`
	Epic1     int32 `json:"Epic_1"`
	Epic2     int32 `json:"Epic_2"`
	Epic3     int32 `json:"Epic_3"`
	Epic4     int32 `json:"Epic_4"`
	Superior0 int32 `json:"Superior_0"`
	Superior1 int32 `json:"Superior_1"`
	Superior2 int32 `json:"Superior_2"`
	Superior3 int32 `json:"Superior_3"`
	Superior4 int32 `json:"Superior_4"`
	Good0     int32 `json:"Good_0"`
	Good1     int32 `json:"Good_1"`
	Good2     int32 `json:"Good_2"`
	Good3     int32 `json:"Good_3"`
	Good4     int32 `json:"Good_4"`
}

type RandomPropAllocationMap map[proto.ItemQuality][5]int32

type RandomPropAllocationsByIlvl map[int32]RandomPropAllocationMap

func ParseRandPropPointsTable(contents string) RandomPropAllocationsByIlvl {
	var rawAllocationMap map[string]RandomPropAllocation
	standardized, err := hujson.Standardize([]byte(contents)) // Removes invalid JSON, such as trailing commas
	if err != nil {
		log.Fatalf("Failed to standardize json %s\n\n%s\n\n%s", err, contents[0:30], contents[len(contents)-30:])
	}

	err = json.Unmarshal(standardized, &rawAllocationMap)
	if err != nil {
		log.Fatalf("failed to parse RandPropPoints table to json %s\n\n%s", err, contents[0:30])
	}
	fmt.Printf("\n--\nRandom property allocations loaded: %d\n--\n", len(rawAllocationMap))

	processedAllocationMap := make(RandomPropAllocationsByIlvl)

	for ilvlStr, allocations := range rawAllocationMap {
		ilvl, err := strconv.ParseInt(ilvlStr, 10, 32)
		if err != nil {
			panic(err)
		}
		processedAllocationMap[int32(ilvl)] = RandomPropAllocationMap{
			proto.ItemQuality_ItemQualityEpic:     [5]int32{allocations.Epic0, allocations.Epic1, allocations.Epic2, allocations.Epic3, allocations.Epic4},
			proto.ItemQuality_ItemQualityRare:     [5]int32{allocations.Superior0, allocations.Superior1, allocations.Superior2, allocations.Superior3, allocations.Superior4},
			proto.ItemQuality_ItemQualityUncommon: [5]int32{allocations.Good0, allocations.Good1, allocations.Good2, allocations.Good3, allocations.Good4},
		}
	}

	return processedAllocationMap
}

func (allocationMap RandomPropAllocationsByIlvl) CalcItemAllocation(item *proto.UIItem) int32 {
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
			switch item.HandType {
			case proto.HandType_HandTypeTwoHand:
				idx = 0
			default:
				idx = 3
			}
		}
	}

	return allocationMap[item.Ilvl][item.Quality][idx]
}
