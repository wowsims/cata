package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// QueryUIItems loads raw items from the database, converts them to UIItems,
// and returns a slice of UIItems.
func QueryUIItems() ([]*proto.UIItem, error) {
	iconsMap, err := LoadArtTexturePaths("./assets/db_inputs/ArtTextureID.lua")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize icons map: %v", err)
	}
	helper, err := NewDBHelper("./tools/database/wowsims.db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database helper: %v", err)
	}
	defer helper.Close()

	LoadItemStatEffects(helper)
	// Load necessary raw item data.
	LoadItemDamageTables(helper)
	LoadRawItems(helper, "s.OverallQualityId != 7 AND s.ScalingStatDistributionID == 0 AND (ItemClassName = 'Armor' OR ItemClassName = 'Weapon') AND s.Display_lang != '' AND (s.ID != 34219 AND s.Display_lang NOT LIKE '%Test%' AND s.Display_lang NOT LIKE 'QA%')")

	var items []*proto.UIItem
	total := len(RawItems)
	fmt.Println("Parsing items\n")
	for i, rawItem := range RawItems {
		switch rawItem.itemClassName {
		case "Weapon", "Armor":
			item, err := RawItemToUIItem(helper, rawItem)
			if err != nil {
				log.Printf("Error processing item row: %v", err)
				continue
			}
			item.Icon = strings.ToLower(GetIconName(iconsMap, rawItem.FDID))
			items = append(items, item)
		default:
			// Skip items we are not processing.
		}
		if i%500 == 0 {

			printProgressBar(i+1, total)
		}
	}
	fmt.Println() // Newline after progress bar.
	return items, nil
}

// QueryUIGems loads raw gems from the database, converts them to UIGems,
// and returns a slice of UIGems.
func QueryUIGems() ([]*proto.UIGem, error) {
	iconsMap, err := LoadArtTexturePaths("./assets/db_inputs/ArtTextureID.lua")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize icons map: %v", err)
	}
	helper, err := NewDBHelper("./tools/database/wowsims.db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database helper: %v", err)
	}
	defer helper.Close()

	LoadRawGems(helper)
	LoadRawSpellEffects(helper)
	totalGems := len(RawGems)
	var gems []*proto.UIGem
	fmt.Println("Parsing gems\n")
	for i, rawGem := range RawGems {
		gem := &proto.UIGem{
			Id:      int32(rawGem.ItemId),
			Name:    rawGem.Name,
			Icon:    strings.ToLower(GetIconName(iconsMap, rawGem.FDID)),
			Quality: qualityToItemQualityMap[rawGem.Quality],
			Stats:   stats.Stats{}.ToProtoArray(),
			Color:   ConvertGemTypeToProto(rawGem.GemType),
			Unique:  rawGem.Flags.Has(UniqueEquipped),
		}

		if rawGem.IsJc {
			gem.RequiredProfession = proto.Profession_Jewelcrafting
		}
		processGemStats(rawGem, gem)
		gems = append(gems, gem)
		printProgressBar(i+1, totalGems)
	}
	fmt.Println() // Newline after progress bar.
	return gems, nil
}

// QueryUIEnchants loads raw enchants, converts them to UIEnchants,
// and returns two slices: one for enchants and one for complex enchant effect IDs (thats not as simple as adding stats).
func QueryUIEnchants() ([]*proto.UIEnchant, []int, error) {
	iconsMap, err := LoadArtTexturePaths("./assets/db_inputs/ArtTextureID.lua")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize icons map: %v", err)
	}
	helper, err := NewDBHelper("./tools/database/wowsims.db")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize database helper: %v", err)
	}
	defer helper.Close()

	LoadRawEnchants(helper)
	totalEnchants := len(RawEnchants)
	var enchants []*proto.UIEnchant
	var complexEnchants []int
	fmt.Println("Parsing enchants\n")
	for i, rawEnchant := range RawEnchants {
		enchant := &proto.UIEnchant{
			Name:               rawEnchant.Name,
			ItemId:             int32(rawEnchant.ItemId),
			SpellId:            int32(rawEnchant.SpellId),
			EffectId:           int32(rawEnchant.EffectId),
			ClassAllowlist:     GetClassesFromClassMask(rawEnchant.ClassMask),
			ExtraTypes:         []proto.ItemType{},
			Stats:              stats.Stats{}.ToProtoArray(),
			Quality:            qualityToItemQualityMap[rawEnchant.Quality],
			RequiredProfession: GetProfession(rawEnchant.RequiredProfession),
		}

		if rawEnchant.FDID == 0 {
			enchant.Icon = "trade_engraving"
		} else {
			enchant.Icon = strings.ToLower(GetIconName(iconsMap, rawEnchant.FDID))
		}

		if rawEnchant.IsWeaponEnchant {
			// Process weapon enchants.
			enchant.Type = proto.ItemType_ItemTypeWeapon
			if rawEnchant.SubClassMask == 1024 {
				// Staff only.
				enchant.EnchantType = proto.EnchantType_EnchantTypeStaff
			}
			if rawEnchant.SubClassMask == 262156 {
				enchant.Type = proto.ItemType_ItemTypeRanged
			}
			if rawEnchant.SubClassMask == 136546 {
				// Two-handed weapon.
				enchant.EnchantType = proto.EnchantType_EnchantTypeTwoHand
			}
		} else {
			// Process non-weapon enchants.
			if rawEnchant.SubClassMask == 65 {
				enchant.EnchantType = proto.EnchantType_EnchantTypeOffHand
				enchant.Type = proto.ItemType_ItemTypeWeapon
			}
			if rawEnchant.SubClassMask == 96 || rawEnchant.SubClassMask == 64 {
				enchant.EnchantType = proto.EnchantType_EnchantTypeShield
				enchant.Type = proto.ItemType_ItemTypeWeapon
			}
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
		processEnchantStats(rawEnchant, enchant)
		if enchant.Type == proto.ItemType_ItemTypeUnknown {
			fmt.Println(rawEnchant)
		}
		enchants = append(enchants, enchant)
		if enchantHasComplexEffect(rawEnchant) {
			complexEnchants = append(complexEnchants, rawEnchant.EffectId)
		}
		printProgressBar(i+1, totalEnchants)
	}
	fmt.Println() // Newline after progress bar.
	return enchants, complexEnchants, nil
}
