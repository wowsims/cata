package database

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/tools"
)

func printProgressBar(current, total int) {
	percent := float64(current) / float64(total)
	barLength := 50
	filledLength := int(percent * float64(barLength))
	bar := ""
	for i := 0; i < filledLength; i++ {
		bar += "="
	}
	for i := filledLength; i < barLength; i++ {
		bar += " "
	}
	fmt.Printf("\r[%s] %d%% (%d of %d)", bar, int(percent*100), current, total)
}

func QueryItems() ([]*proto.UIItem, error) {
	iconsMap, err := LoadArtTexturePaths("./assets/db_inputs/ArtTextureID.lua")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize icons map: %v", err)
	}
	helper, err := NewDBHelper("./tools/database/wowsims.db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database helper: %v", err)
	}
	defer helper.Close()
	var items []*proto.UIItem

	LoadItemDamageTables(helper)
	LoadRawItems(helper, "s.OverallQualityId != 7 AND s.ScalingStatDistributionID == 0 AND s.ItemLevel > 240")

	total := len(RawItems)
	fmt.Println("Parsing items\n\n")
	for i, rawItem := range RawItems {
		switch rawItem.itemClassName {
		case "Weapon", "Armor":
			item, err := RawItemToUIItem(helper, rawItem)
			if err != nil {
				log.Printf("Error processing item row: %v", err)
				continue
			}
			items = append(items, item)
			item.Icon = GetIconName(iconsMap, rawItem.FDID)
			break
		default:
			//fmt.Println("Currently not processing", rawItem.itemClassName)
			break
		}

		printProgressBar(i+1, total)
	}
	var buffer bytes.Buffer

	tools.WriteProtoArrayToBuffer(items, &buffer, "items")

	jsonString := buffer.String()

	filePath := "./assets/db_inputs/testItems.json"
	os.WriteFile(filePath, []byte(jsonString), 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON file: %s", err.Error())
	}
	log.Printf("JSON successfully written to %s", filePath)

	//Start loading gems
	LoadRawGems(helper)
	totalGems := len(RawGems)
	fmt.Println("Parsing gems\n\n")
	for i, rawGem := range RawGems {
		gem := &proto.UIGem{
			Id:      int32(rawGem.ItemId),
			Icon:    GetIconName(iconsMap, rawGem.FDID),
			Quality: qualityToItemQualityMap[rawGem.Quality],
			Stats:   stats.Stats{}.ToProtoArray(),
			Color:   ConvertGemTypeToProto(rawGem.GemType),
			Unique:  rawGem.Flags.Has(UniqueEquipped),
		}
		if rawGem.IsJc {
			gem.RequiredProfession = proto.Profession_Jewelcrafting
		}
		processGemStats(rawGem, gem)

		printProgressBar(i+1, totalGems)
	}

	//Start loading enchants
	LoadRawEnchants(helper)
	totalEnchants := len(RawEnchants)
	fmt.Println("Parsing enchants\n\n")
	for i, rawEnchant := range RawEnchants {
		enchant := &proto.UIEnchant{
			Name:           rawEnchant.Name,
			ItemId:         int32(rawEnchant.ItemId),
			SpellId:        int32(rawEnchant.SpellId),
			ClassAllowlist: GetClassesFromClassMask(rawEnchant.ClassMask),
			ExtraTypes:     []proto.ItemType{},
		}
		if rawEnchant.IsWeaponEnchant {
			//Handle this diff
		} else {
			for flag, name := range inventoryNames {
				if InventoryType(rawEnchant.InvTypesMask)&flag != 0 {
					if enchant.Type != proto.ItemType_ItemTypeUnknown {
						enchant.ExtraTypes = append(enchant.ExtraTypes, name.ItemType)
					} else {
						enchant.Type = name.ItemType
					}
				}
			}
		}
		printProgressBar(i+1, totalEnchants)
	}
	//Start loading enchants

	return items, nil
}
