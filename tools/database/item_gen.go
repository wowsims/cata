package database

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/wowsims/cata/sim/core/proto"
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
	LoadRawItems(helper)

	total := len(RawItems)
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
		//First we process equippable items (UIITems)

		//Then we handle Gems
		//Then we handle food and pots
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
	return items, nil
}
