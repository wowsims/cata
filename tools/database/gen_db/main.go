package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/wowsims/mop/sim"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	_ "github.com/wowsims/mop/sim/encounters" // Needed for preset encounters.
	"github.com/wowsims/mop/tools"
	"github.com/wowsims/mop/tools/database"
	"github.com/wowsims/mop/tools/database/dbc"
)

// To do a full re-scrape, delete the previous output file first.
// go run ./tools/database/gen_db -outDir=assets -gen=atlasloot
// go run ./tools/database/gen_db -outDir=assets -gen=db

var outDir = flag.String("outDir", "assets", "Path to output directory for writing generated .go files.")
var genAsset = flag.String("gen", "", "Asset to generate. Valid values are 'db', 'atlasloot', 'wowhead-items', 'wowhead-spells', 'wowhead-itemdb', 'mop-items', and 'wago-db2-items'")
var dbPath = flag.String("dbPath", "./tools/database/wowsims.db", "Location of wowsims.db file from the DB2ToSqliteTool")

func main() {
	flag.Parse()

	database.DatabasePath = *dbPath

	if *outDir == "" {
		panic("outDir flag is required!")
	}

	dbDir := fmt.Sprintf("%s/database", *outDir)
	inputsDir := fmt.Sprintf("%s/db_inputs", *outDir)

	if *genAsset == "atlasloot" {
		helper, err := database.NewDBHelper()
		if err != nil {
			log.Fatalf("failed to initialize database: %v", err)
		}
		defer helper.Close()

		db := database.ReadAtlasLootData(helper)
		db.WriteJson(fmt.Sprintf("%s/atlasloot_db.json", inputsDir))
		return
	} else if *genAsset == "reforge-stats" {
		//Todo: fill this when we have information from wowhead @ Neteyes - Gehennas
		// For now, the version we have was taken from https://web.archive.org/web/20120201045249js_/http://www.wowhead.com/data=item-scaling
		return
	} else if *genAsset != "db" {
		panic("Invalid gen value")
	}
	helper, err := database.NewDBHelper()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer helper.Close()

	if err := database.RunOverrides(helper, "tools/database/overrides"); err != nil {
		log.Fatalf("failed to run overrides: %v", err)
	}

	_, err = database.LoadAndWriteRawRandomSuffixes(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}

	_, err = database.LoadAndWriteRawItems(helper, "s.OverallQualityId != 7 AND NOT (s.Bonding = 2 AND ind.Description_lang IS NOT NULL and ind.Description_lang NOT LIKE '%Season%') AND s.Field_1_15_7_59706_054 = 0 AND s.OverallQualityId != 0 AND (i.ClassID = 2 OR i.ClassID = 4) AND s.Display_lang != '' AND (s.ID != 34219 AND s.Display_lang NOT LIKE '%Test%' AND s.Display_lang NOT LIKE 'QA%')", inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}

	_, err = database.LoadAndWriteRandomPropAllocations(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}

	_, err = database.LoadAndWriteRawGems(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteRawEnchants(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteRawSpellEffects(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteItemStatEffects(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteItemDamageTables(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteItemArmorTotal(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteItemArmorQuality(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteItemArmorShield(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteArmorLocation(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteItemEffects(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	_, err = database.LoadAndWriteSpells(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	consumables, err := database.LoadAndWriteConsumables(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	dropSources, names, err := database.LoadAndWriteDropSources(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	craftingSources := database.LoadCraftedItems(helper)
	repSources := database.LoadRepItems(helper)
	//Todo: See if we cant get rid of these as well
	atlaslootDB := database.ReadDatabaseFromJson(tools.ReadFile(fmt.Sprintf("%s/atlasloot_db.json", inputsDir)))

	// Todo: https://web.archive.org/web/20120201045249js_/http://www.wowhead.com/data=item-scaling
	reforgeStats := database.ParseWowheadReforgeStats(tools.ReadFile(fmt.Sprintf("%s/wowhead_reforge_stats.json", inputsDir)))

	db := database.NewWowDatabase()
	db.Encounters = core.PresetEncounters
	db.ReforgeStats = reforgeStats.ToProto()

	iconsMap, _ := database.LoadArtTexturePaths("./tools/DB2ToSqlite/listfile.csv")
	var instance = dbc.GetDBC()
	instance.LoadSpellScaling()

	database.GenerateProtos(instance, db)

	processItems(instance, iconsMap, names, dropSources, craftingSources, repSources, db)

	for _, gem := range instance.Gems {
		parsed := gem.ToProto()
		if parsed.Icon == "" {
			parsed.Icon = strings.ToLower(database.GetIconName(iconsMap, gem.FDID))
		}
		db.MergeGem(parsed)
	}

	for _, enchant := range instance.Enchants {
		parsed := enchant.ToProto()
		if parsed.Icon == "" {
			parsed.Icon = strings.ToLower(database.GetIconName(iconsMap, enchant.FDID))
		}
		db.MergeEnchant(parsed)
	}

	for _, item := range atlaslootDB.Items {
		if _, ok := db.Items[item.Id]; ok {
			db.MergeItem(item)
		}
	}

	bestByStat := make(map[int]map[int]*dbc.Consumable)

	// Phase 1: find the best consumable per (subclass, stat-index)
	for i := range consumables {
		c := &consumables[i]
		subclass := int(c.SubClassId)

		// ensure the inner map exists
		if _, ok := bestByStat[subclass]; !ok {
			bestByStat[subclass] = make(map[int]*dbc.Consumable)
		}
		bucket := bestByStat[subclass]

		// pull the raw stats array once
		stats := c.ToProto().Stats
		for idx, val := range stats {
			if existing, seen := bucket[idx]; !seen || val > existing.ToProto().Stats[idx] {
				bucket[idx] = c
			}
		}
	}

	// Phase 2: merge each unique consumable exactly once
	seen := make(map[int]bool)
	for _, bucket := range bestByStat {
		for _, c := range bucket {
			if seen[c.Id] {
				continue
			}
			p := c.ToProto()
			p.Icon = strings.ToLower(
				database.GetIconName(iconsMap, c.IconFileDataID),
			)
			db.MergeConsumable(p)
			seen[c.Id] = true
		}
	}

	for _, consumable := range database.ConsumableOverrides {
		db.MergeConsumable(consumable)
	}

	db.MergeItems(database.ItemOverrides)
	db.MergeGems(database.GemOverrides)
	db.MergeEnchants(database.EnchantOverrides)
	ApplyGlobalFilters(db)
	leftovers := db.Clone()
	ApplyNonSimmableFilters(leftovers)
	leftovers.WriteBinaryAndJson(fmt.Sprintf("%s/leftover_db.bin", dbDir), fmt.Sprintf("%s/leftover_db.json", dbDir))
	ApplySimmableFilters(db)
	for _, enchant := range db.Enchants {
		if enchant.ItemId != 0 {
			db.AddItemIcon(enchant.ItemId, enchant.Icon, enchant.Name)
		}
		if enchant.SpellId != 0 {
			db.AddSpellIcon(enchant.SpellId, enchant.Icon, enchant.Name)
		}
	}

	for _, consume := range db.Consumables {
		if len(consume.EffectIds) > 0 {
			for _, se := range consume.EffectIds {
				effect := instance.SpellEffectsById[int(se)]
				db.MergeEffect(effect.ToProto())
			}
		}
	}

	for _, randomSuffix := range instance.RandomSuffix {
		if _, exists := db.RandomSuffixes[int32(randomSuffix.ID)]; !exists {
			db.RandomSuffixes[int32(randomSuffix.ID)] = randomSuffix.ToProto()
		}
	}

	icons, err := database.LoadSpellIcons(helper)
	if err != nil {
		panic("error loading icons")
	}

	addSpellIcons(db, database.SharedSpellsIcons, icons, iconsMap)

	for _, group := range GetAllTalentSpellIds(&inputsDir) {
		addSpellIcons(db, group, icons, iconsMap)
	}

	for _, group := range GetAllRotationSpellIds() {
		addSpellIcons(db, group, icons, iconsMap)
	}

	craftedSpellIds := []int32{}
	for _, item := range db.Items {
		for _, source := range item.Sources {
			if crafted := source.GetCrafted(); crafted != nil {
				craftedSpellIds = append(craftedSpellIds, crafted.SpellId)
			}
		}
		if item.Phase < 2 {
			item.Phase = InferPhase(item)
		}
	}
	addSpellIcons(db, craftedSpellIds, icons, iconsMap)

	database.LoadAndWriteEnchantDescriptions("assets/enchants/descriptions.json", db, instance)

	atlasDBProto := atlaslootDB.ToUIProto()
	db.MergeZones(atlasDBProto.Zones)
	db.MergeNpcs(atlasDBProto.Npcs)
	db.WriteBinaryAndJson(fmt.Sprintf("%s/db.bin", dbDir), fmt.Sprintf("%s/db.json", dbDir))
}

func InferPhase(item *proto.UIItem) int32 {
	ilvl := item.ScalingOptions[int32(proto.ItemLevelState_Base)].Ilvl
	name := item.Name
	quality := item.Quality

	if strings.Contains(name, "Necklace of the Terra-Cotta") {
		return 4
	}

	//- Any blue pvp ''Crafted'' item of ilvl 458 is 5.2
	//- Any blue pvp ''Crafted'' item of ilvl 476 is 5.4
	if strings.Contains(name, "Crafted") {
		switch ilvl {
		case 458:
			return 3
		case 476:
			return 5
		}
	}

	//- Any "Tyrannical" item is 5.2
	//- Any "Grievous" item is 5.4
	//- Any "Prideful" item is 5.4
	switch {
	case strings.Contains(name, "Grievous"),
		strings.Contains(name, "Prideful"):
		return 5
	case strings.Contains(name, "Tyrannical"):
		return 3
	}

	//- Any 476 epic item with random stats is 5.1
	//- Any 496 epic item with random stats is 5.4
	//- Any 516 epic items with random stats are 5.3
	//- Any 535 epic items with random stats are 5.4
	//- Any 489 random stat epic is 5.3
	if item.RandPropPoints > 0 {
		switch ilvl {
		case 476:
			return 2
		case 489:
			return 4
		case 496:
			return 5
		case 516:
			return 4
		case 535:
			return 5
		}
	}

	//iLvl 600 legendary vs. epic
	if ilvl == 600 {
		if quality == proto.ItemQuality_ItemQualityLegendary {
			return 5
		}
		if quality == proto.ItemQuality_ItemQualityEpic {
			return 4
		}
	}

	//- Any item above ilvl 542 is 5.4 (except the 600 ilvl Epic Cloaks from the legendary questline)
	if ilvl > 542 && quality < proto.ItemQuality_ItemQualityLegendary {
		return 5
	}

	//- Any 483 green item is a boosted level 90 item in 5.4
	if ilvl == 483 && quality == proto.ItemQuality_ItemQualityUncommon {
		return 5
	}

	//- All pve tier items of ilvl 502/522/535 are 5.2
	//- All pve tier items of ilvl 528/540/553/566 are 5.4
	if item.SetId > 0 {
		switch ilvl {
		case 528, 540, 553, 566:
			return 5
		case 502, 522, 535:
			return 3
		}
	}

	// Timeless Isle trinkets are all ilvl 496 and does not have a source listed.
	if item.Sources == nil {
		if item.Type == proto.ItemType_ItemTypeTrinket && ilvl == 496 {
			return 3
		}
	}

	//AtlasLoot‐style source checks
	for _, src := range item.Sources {
		//- All items with Reputation requirements of "Shado-Pan Assault" are 5.2
		if rep := src.GetRep(); rep != nil {
			if rep.RepFactionId == proto.RepFaction_RepFactionShadoPanAssault {
				return 3
			}
			if rep.RepFactionId == proto.RepFaction_RepFactionOperationShieldwall || rep.RepFactionId == proto.RepFaction_RepFactionDominanceOffensive {
				return 2
			}
		}
		if craft := src.GetCrafted(); craft != nil {
			switch ilvl {
			case 476, 496:
				return 1
			case 502:
				return 4
			case 522:
				return 3
			case 553:
				return 4
			}
		}
		if drop := src.GetDrop(); drop != nil {
			switch drop.ZoneId {
			case 6297, 6125, 6067:
				return 1
			case 6622:
				return 3
			case 6738:
				return 5
			}
			//- All "Oondasta (World Boss)" items are 5.2
			if drop.NpcId == 826 {
				return 3
			}
			//- All "Ordos (World Boss)" items are 5.4
			if drop.NpcId == 861 {
				return 5
			}
		}
	}

	// Any 489 random stat epic is 5.3
	if ilvl >= 489 && len(item.RandomSuffixOptions) > 0 {
		return 2
	}

	// high ilvl greens probably boosted
	if ilvl > 440 && quality < proto.ItemQuality_ItemQualityRare {
		return 5
	}

	if ilvl <= 463 {
		return 1
	}

	switch ilvl {
	case 476, 483, 489, 496:
		return 1
	case 502, 522, 535, 541:
		return 3
	case 553, 528, 566, 540:
		return 5
	}

	return 0
}

func processItems(instance *dbc.DBC, iconsMap map[int]string, names map[int]string, dropSources map[int][]*proto.DropSource, craftingSources map[int][]*proto.CraftedSource, repSources map[int][]*proto.RepSource, db *database.WowDatabase) {
	sourceMap := make(map[string][]*proto.UIItemSource, len(instance.Items))
	parsedItems := make([]*proto.UIItem, 0, len(instance.Items))
	for _, item := range instance.Items {
		if item.Flags2&0x10 != 0 && (item.StatAlloc[0] > 0 && item.StatAlloc[0] < 600) {
			continue
		}
		parsed := item.ToUIItem()
		if parsed.Icon == "" {
			parsed.Icon = strings.ToLower(database.GetIconName(iconsMap, item.FDID))
		}

		drops := dropSources[int(item.Id)]
		if drops != nil {
			sources := make([]*proto.UIItemSource, 0, len(drops))
			for _, drop := range drops {
				sources = append(sources, &proto.UIItemSource{
					Source: &proto.UIItemSource_Drop{Drop: drop},
				})
				db.MergeZone(&proto.UIZone{Id: drop.ZoneId, Name: names[int(drop.ZoneId)]})
				db.MergeNpc(&proto.UINPC{Id: drop.NpcId, Name: drop.OtherName, ZoneId: drop.ZoneId})
			}
			parsed.Sources = sources
			sourceMap[parsed.Name] = sources
		}

		crafted := craftingSources[int(item.Id)]
		if crafted != nil {
			sources := make([]*proto.UIItemSource, 0, len(crafted))
			for _, craft := range crafted {
				sources = append(sources, &proto.UIItemSource{
					Source: &proto.UIItemSource_Crafted{Crafted: craft},
				})
			}
			parsed.Sources = sources
			sourceMap[parsed.Name] = sources
		}

		rep := repSources[int(item.Id)]
		if rep != nil {
			sources := make([]*proto.UIItemSource, 0, len(rep))
			for _, repItem := range rep {
				sources = append(sources, &proto.UIItemSource{
					Source: &proto.UIItemSource_Rep{Rep: repItem},
				})
			}
			parsed.Sources = sources
			sourceMap[parsed.Name] = sources
		}

		parsedItems = append(parsedItems, parsed)
	}

	for _, parsed := range parsedItems {
		if len(parsed.Sources) == 0 {
			if fallbacks, ok := sourceMap[parsed.Name]; ok {
				parsed.Sources = fallbacks
			}
		}
	}
	db.MergeItems(parsedItems)
}

// Filters out entities which shouldn't be included anywhere.
func ApplyGlobalFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := database.ItemAllowList[item.Id]; ok {
			return true
		}
		if _, ok := database.ItemDenyList[item.Id]; ok {
			return false
		}
		if len(item.ScalingOptions) <= 0 {
			return false
		}
		if item.ScalingOptions[0].Ilvl > 600 || item.ScalingOptions[0].Ilvl < 100 {
			return false
		}
		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(item.Name) {
				return false
			}
		}
		return true
	})

	// There is an 'unavailable' version of every naxx set, e.g. https://www.wowhead.com/mop-classic/item=43728/bonescythe-gauntlets
	heroesItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return strings.HasPrefix(item.Name, "Heroes' ")
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		nameToMatch := "Heroes' " + item.Name
		for _, heroItem := range heroesItems {
			if heroItem.Name == nameToMatch {
				return false
			}
		}
		return true
	})

	// There is an 'unavailable' version of many t8 set pieces, e.g. https://www.wowhead.com/mop-classic/item=46235/darkruned-gauntlets
	valorousItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return strings.HasPrefix(item.Name, "Valorous ")
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		nameToMatch := "Valorous " + item.Name
		for _, item := range valorousItems {
			if item.Name == nameToMatch {
				return false
			}
		}
		return true
	})

	// There is an 'unavailable' version of many t9 set pieces, e.g. https://www.wowhead.com/mop-classic/item=48842/thralls-hauberk
	triumphItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return strings.HasSuffix(item.Name, "of Triumph")
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		nameToMatch := item.Name + " of Triumph"
		for _, item := range triumphItems {
			if item.Name == nameToMatch {
				return false
			}
		}
		return true
	})

	db.Gems = core.FilterMap(db.Gems, func(_ int32, gem *proto.UIGem) bool {
		if _, ok := database.GemDenyList[gem.Id]; ok {
			return false
		}

		if gem.Quality == proto.ItemQuality_ItemQualityLegendary || gem.Id == 95348 {
			gem.Phase = 3
		} else {
			gem.Phase = 1
		}

		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(gem.Name) {
				return false
			}
		}
		return true
	})

	db.Gems = core.FilterMap(db.Gems, func(_ int32, gem *proto.UIGem) bool {
		if strings.HasSuffix(gem.Name, "Stormjewel") {
			gem.Unique = false
		}
		return true
	})

	db.ItemIcons = core.FilterMap(db.ItemIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
	db.SpellIcons = core.FilterMap(db.SpellIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != "" && icon.Id != 0
	})

	db.Enchants = core.FilterMap(db.Enchants, func(_ database.EnchantDBKey, enchant *proto.UIEnchant) bool {
		// MoP no longer has head enchants, so filter them.
		if enchant.Type == proto.ItemType_ItemTypeHead {
			return false
		}
		if _, ok := database.EnchantDenyListSpells[enchant.SpellId]; ok {
			return false
		}
		if _, ok := database.EnchantDenyListItems[enchant.ItemId]; ok {
			return false
		}
		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(enchant.Name) {
				return false
			}
		}
		return !strings.HasPrefix(enchant.Name, "QA") && !strings.HasPrefix(enchant.Name, "Test") && !strings.HasPrefix(enchant.Name, "TEST")
	})

	db.Consumables = core.FilterMap(db.Consumables, func(_ int32, consumable *proto.Consumable) bool {
		if slices.Contains(database.ConsumableAllowList, consumable.Id) {
			return true
		}
		if slices.Contains(database.ConsumableDenyList, consumable.Id) {
			return false
		}
		if allZero(consumable.Stats) && consumable.Type != proto.ConsumableType_ConsumableTypePotion {
			return false
		}

		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(consumable.Name) {
				return false
			}
		}

		if consumable.Type == proto.ConsumableType_ConsumableTypeUnknown || consumable.Type == proto.ConsumableType_ConsumableTypeScroll {
			return false
		}

		return !strings.HasPrefix(consumable.Name, "QA") && !strings.HasPrefix(consumable.Name, "Test") && !strings.HasPrefix(consumable.Name, "TEST")
	})
}

func allZero(stats []float64) bool {
	for _, val := range stats {
		if val != 0 {
			return false
		}
	}
	return true
}

// Filters out entities which shouldn't be included in the sim.
func ApplySimmableFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, simmableItemFilter)
	db.Gems = core.FilterMap(db.Gems, simmableGemFilter)
}

func ApplyNonSimmableFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(id int32, item *proto.UIItem) bool {
		return !simmableItemFilter(id, item)
	})
	db.Gems = core.FilterMap(db.Gems, func(id int32, gem *proto.UIGem) bool {
		return !simmableGemFilter(id, gem)
	})
}
func simmableItemFilter(_ int32, item *proto.UIItem) bool {
	if _, ok := database.ItemAllowList[item.Id]; ok {
		return true
	}

	if item.Quality < proto.ItemQuality_ItemQualityUncommon {
		return false
	} else if item.Quality == proto.ItemQuality_ItemQualityArtifact {
		return false
	} else if item.Quality >= proto.ItemQuality_ItemQualityHeirloom {
		return false
	} else if item.Quality <= proto.ItemQuality_ItemQualityEpic {
		if item.ScalingOptions[0].Ilvl < 372 {
			return false
		}
	} else {
		// Epic and legendary items might come from classic, so use a lower ilvl threshold.
		if item.ScalingOptions[0].Ilvl <= 359 {
			return false
		}
	}
	if item.ScalingOptions[0].Ilvl == 0 {
		fmt.Printf("Missing ilvl: %s\n", item.Name)
	}

	return true
}
func simmableGemFilter(_ int32, gem *proto.UIGem) bool {
	if _, ok := database.GemAllowList[gem.Id]; ok {
		return true
	}

	// Arbitrary to filter out old gems
	if gem.Id < 46000 {
		return false
	}

	return gem.Quality >= proto.ItemQuality_ItemQualityUncommon
}

type TalentConfig struct {
	FieldName string `json:"fieldName"`
	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	SpellId int32 `json:"spellId"`
}

type TalentTreeConfig struct {
	BackgroundUrl string         `json:"backgroundUrl"`
	Talents       []TalentConfig `json:"talents"`
}

func getSpellIdsFromTalentJson(infile *string) []int32 {
	data, err := os.ReadFile(*infile)

	if err != nil {
		log.Fatalf("failed to load talent json file: %s", err)
	}

	var buf bytes.Buffer

	err = json.Compact(&buf, []byte(data))
	if err != nil {
		log.Fatalf("failed to compact json: %s", err)
	}

	var talentTree TalentTreeConfig

	err = json.Unmarshal(buf.Bytes(), &talentTree)
	if err != nil {
		log.Fatalf("failed to parse talent to json %s", err)
	}
	spellIds := make([]int32, 0)

	for _, talent := range talentTree.Talents {
		spellIds = append(spellIds, talent.SpellId)
	}

	return spellIds
}

func GetAllTalentSpellIds(inputsDir *string) map[string][]int32 {
	talentsDir := fmt.Sprintf("%s/../../ui/core/talents/trees", *inputsDir)
	specFiles := []string{
		"death_knight.json",
		"druid.json",
		"hunter.json",
		"mage.json",
		"paladin.json",
		"priest.json",
		"rogue.json",
		"shaman.json",
		"warlock.json",
		"warrior.json",
		"monk.json",
	}

	ret_db := make(map[string][]int32, 0)

	for _, specFile := range specFiles {
		specPath := fmt.Sprintf("%s/%s", talentsDir, specFile)
		ret_db[specFile[:len(specFile)-5]] = getSpellIdsFromTalentJson(&specPath)
	}

	return ret_db

}

type GlyphID struct {
	ItemID  int32 `json:"itemId"`
	SpellID int32 `json:"spellId"`
}

func CreateTempAgent(r *proto.Raid) core.Agent {
	encounter := core.MakeSingleTargetEncounter(0.0)
	env, _, _ := core.NewEnvironment(r, encounter, false)
	return env.Raid.Parties[0].Players[0]
}

type RotContainer struct {
	Name string
	Raid *proto.Raid
}

func GetAllRotationSpellIds() map[string][]int32 {
	sim.RegisterAll()

	rotMapping := []RotContainer{
		// Death Knight
		{Name: "bloodDeathKnight", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDeathKnight,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_BloodDeathKnight{BloodDeathKnight: &proto.BloodDeathKnight{Options: &proto.BloodDeathKnight_Options{ClassOptions: &proto.DeathKnightOptions{}}, Rotation: &proto.BloodDeathKnight_Rotation{}}}), nil, nil, nil)},
		{Name: "frostDeathKnight", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDeathKnight,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_FrostDeathKnight{FrostDeathKnight: &proto.FrostDeathKnight{Options: &proto.FrostDeathKnight_Options{ClassOptions: &proto.DeathKnightOptions{}}}}), nil, nil, nil)},
		{Name: "unholyDeathKnight", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDeathKnight,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_UnholyDeathKnight{UnholyDeathKnight: &proto.UnholyDeathKnight{Options: &proto.UnholyDeathKnight_Options{ClassOptions: &proto.DeathKnightOptions{}}}}), nil, nil, nil)},

		// Druid
		{Name: "balanceDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_BalanceDruid{BalanceDruid: &proto.BalanceDruid{Options: &proto.BalanceDruid_Options{ClassOptions: &proto.DruidOptions{}}}}), nil, nil, nil)},
		{Name: "feralDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_FeralDruid{FeralDruid: &proto.FeralDruid{Options: &proto.FeralDruid_Options{ClassOptions: &proto.DruidOptions{}}, Rotation: &proto.FeralDruid_Rotation{}}}), nil, nil, nil)},
		// {Name: "guardianDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
		// 	Class:         proto.Class_ClassDruid,
		// 	Equipment:     &proto.EquipmentSpec{},
		// 	TalentsString: "000000",
		// }, &proto.Player_FeralTankDruid{FeralTankDruid: &proto.FeralTankDruid{Options: &proto.FeralTankDruid_Options{ClassOptions: &proto.DruidOptions{}}}}), nil, nil, nil)},
		{Name: "restorationDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_RestorationDruid{RestorationDruid: &proto.RestorationDruid{Options: &proto.RestorationDruid_Options{ClassOptions: &proto.DruidOptions{}}}}), nil, nil, nil)},

		// Hunter
		{Name: "beastMasteryHunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassHunter,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_BeastMasteryHunter{BeastMasteryHunter: &proto.BeastMasteryHunter{Options: &proto.BeastMasteryHunter_Options{ClassOptions: &proto.HunterOptions{}}}}), nil, nil, nil)},
		{Name: "marksmanshipHunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassHunter,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_MarksmanshipHunter{MarksmanshipHunter: &proto.MarksmanshipHunter{Options: &proto.MarksmanshipHunter_Options{ClassOptions: &proto.HunterOptions{}}}}), nil, nil, nil)},
		{Name: "survivalHunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassHunter,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_SurvivalHunter{SurvivalHunter: &proto.SurvivalHunter{Options: &proto.SurvivalHunter_Options{ClassOptions: &proto.HunterOptions{}}}}), nil, nil, nil)},

		// Mage
		{Name: "arcaneMage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_ArcaneMage{ArcaneMage: &proto.ArcaneMage{Options: &proto.ArcaneMage_Options{ClassOptions: &proto.MageOptions{}}}}), nil, nil, nil)},
		{Name: "fireMage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_FireMage{FireMage: &proto.FireMage{Options: &proto.FireMage_Options{ClassOptions: &proto.MageOptions{}}}}), nil, nil, nil)},
		{Name: "frostMage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_FrostMage{FrostMage: &proto.FrostMage{Options: &proto.FrostMage_Options{ClassOptions: &proto.MageOptions{}}}}), nil, nil, nil)},

		// Paladin
		{Name: "holyPaladin", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_HolyPaladin{HolyPaladin: &proto.HolyPaladin{Options: &proto.HolyPaladin_Options{ClassOptions: &proto.PaladinOptions{}}}}), nil, nil, nil)},
		{Name: "protPaladin", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_ProtectionPaladin{ProtectionPaladin: &proto.ProtectionPaladin{Options: &proto.ProtectionPaladin_Options{ClassOptions: &proto.PaladinOptions{}}}}), nil, nil, nil)},
		{Name: "retPaladin", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Race:          proto.Race_RaceBloodElf,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_RetributionPaladin{RetributionPaladin: &proto.RetributionPaladin{Options: &proto.RetributionPaladin_Options{ClassOptions: &proto.PaladinOptions{}}}}), nil, nil, nil)},

		// Priest
		{Name: "disciplinePriest", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_DisciplinePriest{DisciplinePriest: &proto.DisciplinePriest{Options: &proto.DisciplinePriest_Options{ClassOptions: &proto.PriestOptions{}}}}), nil, nil, nil)},
		{Name: "holyPriest", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_HolyPriest{HolyPriest: &proto.HolyPriest{Options: &proto.HolyPriest_Options{ClassOptions: &proto.PriestOptions{}}}}), nil, nil, nil)},
		{Name: "shadowPriest", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_ShadowPriest{ShadowPriest: &proto.ShadowPriest{Options: &proto.ShadowPriest_Options{ClassOptions: &proto.PriestOptions{}}}}), nil, nil, nil)},

		// Rogue
		{Name: "assassinationRogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassRogue,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_AssassinationRogue{AssassinationRogue: &proto.AssassinationRogue{Options: &proto.AssassinationRogue_Options{ClassOptions: &proto.RogueOptions{}}}}), nil, nil, nil)},
		{Name: "combatRogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassRogue,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_CombatRogue{CombatRogue: &proto.CombatRogue{Options: &proto.CombatRogue_Options{ClassOptions: &proto.RogueOptions{}}}}), nil, nil, nil)},
		{Name: "subtletyRogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassRogue,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_SubtletyRogue{SubtletyRogue: &proto.SubtletyRogue{Options: &proto.SubtletyRogue_Options{ClassOptions: &proto.RogueOptions{}}}}), nil, nil, nil)},

		// Shaman
		{Name: "elementalShaman", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}}), nil, nil, nil)},
		{Name: "enhancementShaman", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}}}}), nil, nil, nil)},
		{Name: "restorationShaman", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_RestorationShaman{RestorationShaman: &proto.RestorationShaman{Options: &proto.RestorationShaman_Options{ClassOptions: &proto.ShamanOptions{}}}}), nil, nil, nil)},

		// Warlock
		{Name: "afflictionWarlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_AfflictionWarlock{AfflictionWarlock: &proto.AfflictionWarlock{Options: &proto.AfflictionWarlock_Options{ClassOptions: &proto.WarlockOptions{}}}}), nil, nil, nil)},
		{Name: "demonologyWarlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_DemonologyWarlock{DemonologyWarlock: &proto.DemonologyWarlock{Options: &proto.DemonologyWarlock_Options{ClassOptions: &proto.WarlockOptions{}}}}), nil, nil, nil)},
		{Name: "destructionWarlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_DestructionWarlock{DestructionWarlock: &proto.DestructionWarlock{Options: &proto.DestructionWarlock_Options{ClassOptions: &proto.WarlockOptions{}}}}), nil, nil, nil)},

		// Warrior
		{Name: "armsWarrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_ArmsWarrior{ArmsWarrior: &proto.ArmsWarrior{Options: &proto.ArmsWarrior_Options{ClassOptions: &proto.WarriorOptions{}}}}), nil, nil, nil)},
		{Name: "furyWarrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_FuryWarrior{FuryWarrior: &proto.FuryWarrior{Options: &proto.FuryWarrior_Options{ClassOptions: &proto.WarriorOptions{}}}}), nil, nil, nil)},
		{Name: "protectionWarrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{Options: &proto.ProtectionWarrior_Options{ClassOptions: &proto.WarriorOptions{}}}}), nil, nil, nil)},

		// Monk
		{Name: "brewmasterMonk", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMonk,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_BrewmasterMonk{BrewmasterMonk: &proto.BrewmasterMonk{Options: &proto.BrewmasterMonk_Options{ClassOptions: &proto.MonkOptions{}, Stance: proto.MonkStance_SturdyOx}}}), nil, nil, nil)},
		{Name: "mistweaverMonk", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMonk,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_MistweaverMonk{MistweaverMonk: &proto.MistweaverMonk{Options: &proto.MistweaverMonk_Options{ClassOptions: &proto.MonkOptions{}, Stance: proto.MonkStance_WiseSerpent}}}), nil, nil, nil)},
		{Name: "windwalkerMonk", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMonk,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "000000",
		}, &proto.Player_WindwalkerMonk{WindwalkerMonk: &proto.WindwalkerMonk{Options: &proto.WindwalkerMonk_Options{ClassOptions: &proto.MonkOptions{}}}}), nil, nil, nil)},
	}

	ret_db := make(map[string][]int32, 0)

	for _, r := range rotMapping {
		f := CreateTempAgent(r.Raid).GetCharacter()

		spells := make([]int32, 0, len(f.Spellbook))

		for _, s := range f.Spellbook {
			if s.SpellID != 0 {
				spells = append(spells, s.SpellID)
			}
		}

		for _, s := range f.GetAuras() {
			if s.ActionID.SpellID != 0 {
				spells = append(spells, s.ActionID.SpellID)
			}
		}

		ret_db[r.Name] = spells
	}
	return ret_db
}

func addSpellIcons(db *database.WowDatabase, spellIds []int32, icons map[int]database.SpellIcon, iconsMap map[int]string) {
	for _, spellId := range spellIds {
		iconEntry := icons[int(spellId)]
		if iconEntry.Name == "" {
			continue
		}
		db.SpellIcons[spellId] = &proto.IconData{
			Id:      int32(iconEntry.SpellID),
			Name:    iconEntry.Name,
			Icon:    strings.ToLower(database.GetIconName(iconsMap, iconEntry.FDID)),
			HasBuff: iconEntry.HasBuff,
		}
	}
}
