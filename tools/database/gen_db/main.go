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

	"github.com/wowsims/cata/sim"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	_ "github.com/wowsims/cata/sim/encounters" // Needed for preset encounters.
	"github.com/wowsims/cata/tools"
	"github.com/wowsims/cata/tools/database"
	"github.com/wowsims/cata/tools/database/dbc"
)

// To do a full re-scrape, delete the previous output file first.
// go run ./tools/database/gen_db -outDir=assets -gen=atlasloot
// go run ./tools/database/gen_db -outDir=assets -gen=db

var outDir = flag.String("outDir", "assets", "Path to output directory for writing generated .go files.")
var genAsset = flag.String("gen", "", "Asset to generate. Valid values are 'db', 'atlasloot', 'wowhead-items', 'wowhead-spells', 'wowhead-itemdb', 'cata-items', and 'wago-db2-items'")
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

	database.GenerateProtos()

	_, err = database.LoadAndWriteRawRandomSuffixes(helper, inputsDir)
	if err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}

	_, err = database.LoadAndWriteRawItems(helper, "s.OverallQualityId != 7 AND s.ScalingStatDistributionID = 0 AND s.OverallQualityId != 0 AND (i.ClassID = 2 OR i.ClassID = 4) AND s.Display_lang != '' AND (s.ID != 34219 AND s.Display_lang NOT LIKE '%Test%' AND s.Display_lang NOT LIKE 'QA%')", inputsDir)
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
	//Todo: See if we cant get rid of these as well
	atlaslootDB := database.ReadDatabaseFromJson(tools.ReadFile(fmt.Sprintf("%s/atlasloot_db.json", inputsDir)))

	// Todo: https://web.archive.org/web/20120201045249js_/http://www.wowhead.com/data=item-scaling
	reforgeStats := database.ParseWowheadReforgeStats(tools.ReadFile(fmt.Sprintf("%s/wowhead_reforge_stats.json", inputsDir)))

	db := database.NewWowDatabase()
	db.Encounters = core.PresetEncounters
	db.GlyphIDs = getGlyphIDsFromJson(fmt.Sprintf("%s/glyph_id_map.json", inputsDir))
	db.ReforgeStats = reforgeStats.ToProto()

	iconsMap, _ := database.LoadArtTexturePaths("./tools/DB2ToSqlite/listfile.csv")
	var instance = dbc.GetDBC()
	instance.LoadSpellScaling()
	for _, item := range instance.Items {
		parsed := item.ToUIItem()
		if parsed.Icon == "" {
			parsed.Icon = strings.ToLower(database.GetIconName(iconsMap, item.FDID))
		}
		parsed.ItemEffect = dbc.MergeItemEffectsForAllStatesNew(parsed)
		db.MergeItem(parsed)
	}

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
	for _, consumable := range consumables {
		protoConsumable := consumable.ToProto()
		protoConsumable.Icon = strings.ToLower(database.GetIconName(iconsMap, consumable.IconFileDataID))
		db.MergeConsumable(protoConsumable)
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

// Uses heuristics on ilvl + source to infer release phase of an item when missing.
func InferPhase(item *proto.UIItem) int32 {
	if item.Ilvl <= 352 {
		return 1
	}
	// Since Atlasloot populates before we run inferphase, we can use Atlasloot data to help us infer the phase
	// Such as here where I check the Zone Id to correctly place Firelands and DS Dungeons in the correct phase
	for _, source := range item.Sources {
		if drop := source.GetDrop(); drop != nil {
			switch drop.ZoneId {
			case 5723: // Firelands
				return 3
			case 5789, 5844, 5788: // Dragon Soul dungeons
				return 4
			}
		}
	}

	if item.Ilvl >= 397 {
		return 4 // Heroic Rag loot should already be tagged correctly by Wowhead.
	}

	switch item.Ilvl {
	case 353:
		return 2
	case 358, 371, 391:
		return 3
	case 359:
		if item.Quality == proto.ItemQuality_ItemQualityUncommon {
			return 4
		}

		return 1
	case 372, 379:
		return 1
	case 377, 390:
		return 4
	case 365:
		if strings.Contains(item.Name, "Vicious") {
			return 1
		}

		return 3
	case 378:
		for _, itemSource := range item.Sources {
			dropSource := itemSource.GetDrop()

			if (dropSource != nil) && slices.Contains([]int32{5788, 5789, 5844}, dropSource.ZoneId) {
				return 4
			}
		}

		return 3
	case 384:
		if strings.Contains(item.Name, "Ruthless") {
			return 3
		}

		return 4
	default:
		return 0
	}
}

// Filters out entities which shouldn't be included anywhere.
func ApplyGlobalFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := database.ItemDenyList[item.Id]; ok {
			return false
		}
		if len(item.ScalingOptions) <= 0 {
			return false
		}
		if item.ScalingOptions[0].Ilvl > 416 || item.ScalingOptions[0].Ilvl < 100 {
			return false
		}
		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(item.Name) {
				return false
			}
		}
		return true
	})

	// There is an 'unavailable' version of every naxx set, e.g. https://www.wowhead.com/cata/item=43728/bonescythe-gauntlets
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

	// There is an 'unavailable' version of many t8 set pieces, e.g. https://www.wowhead.com/cata/item=46235/darkruned-gauntlets
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

	// There is an 'unavailable' version of many t9 set pieces, e.g. https://www.wowhead.com/cata/item=48842/thralls-hauberk
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

	// Theres an invalid 251 t10 set for every class
	// The invalid set has a higher item id than the 'correct' ones
	t10invalidItems := core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		return item.SetName != "" && item.Ilvl == 251
	})
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		for _, t10item := range t10invalidItems {
			if t10item.Name == item.Name && item.Ilvl == t10item.Ilvl && item.Id > t10item.Id {
				return false
			}
		}
		return true
	})

	db.Gems = core.FilterMap(db.Gems, func(_ int32, gem *proto.UIGem) bool {
		if _, ok := database.GemDenyList[gem.Id]; ok {
			return false
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
		return icon.Name != "" && icon.Icon != ""
	})

	db.Enchants = core.FilterMap(db.Enchants, func(_ database.EnchantDBKey, enchant *proto.UIEnchant) bool {
		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(enchant.Name) {
				return false
			}
		}
		if strings.Contains(enchant.Name, "Template") {
			return false
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
		if item.ScalingOptions[0].Ilvl < 277 {
			return false
		}
	} else {
		// Epic and legendary items might come from classic, so use a lower ilvl threshold.
		if item.ScalingOptions[0].Ilvl <= 200 {
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
	SpellIds  []int32 `json:"spellIds"`
	MaxPoints int32   `json:"maxPoints"`
}

type TalentTreeConfig struct {
	Name          string         `json:"name"`
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

	var talents []TalentTreeConfig

	err = json.Unmarshal(buf.Bytes(), &talents)
	if err != nil {
		log.Fatalf("failed to parse talent to json %s", err)
	}

	spellIds := make([]int32, 0)

	for _, tree := range talents {
		for _, talent := range tree.Talents {
			spellIds = append(spellIds, talent.SpellIds...)

			// Infer omitted spell IDs.
			if len(talent.SpellIds) < int(talent.MaxPoints) {
				curSpellId := talent.SpellIds[len(talent.SpellIds)-1]
				for i := len(talent.SpellIds); i < int(talent.MaxPoints); i++ {
					curSpellId++
					spellIds = append(spellIds, curSpellId)
				}
			}
		}
	}

	return spellIds
}

func GetAllTalentSpellIds(inputsDir *string) map[string][]int32 {
	talentsDir := fmt.Sprintf("%s/../../ui/core/talents/trees", *inputsDir)
	specFiles := []string{
		"death_knight.json",
		"druid.json",
		"hunter.json",
		"hunter_cunning.json",
		"hunter_ferocity.json",
		"hunter_tenacity.json",
		"mage.json",
		"paladin.json",
		"priest.json",
		"rogue.json",
		"shaman.json",
		"warlock.json",
		"warrior.json",
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

func getGlyphIDsFromJson(infile string) []*proto.GlyphID {
	data, err := os.ReadFile(infile)
	if err != nil {
		log.Fatalf("failed to load glyph json file: %s", err)
	}

	var buf bytes.Buffer
	err = json.Compact(&buf, []byte(data))
	if err != nil {
		log.Fatalf("failed to compact json: %s", err)
	}

	var glyphIDs []GlyphID

	err = json.Unmarshal(buf.Bytes(), &glyphIDs)
	if err != nil {
		log.Fatalf("failed to parse glyph IDs to json %s", err)
	}

	return core.MapSlice(glyphIDs, func(gid GlyphID) *proto.GlyphID {
		return &proto.GlyphID{
			ItemId:  gid.ItemID,
			SpellId: gid.SpellID,
		}
	})
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
			TalentsString: "23323223130222311321",
		}, &proto.Player_BloodDeathKnight{BloodDeathKnight: &proto.BloodDeathKnight{Options: &proto.BloodDeathKnight_Options{ClassOptions: &proto.DeathKnightOptions{}}, Rotation: &proto.BloodDeathKnight_Rotation{}}}), nil, nil, nil)},
		{Name: "frostDeathKnight", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDeathKnight,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-32231222233112212331",
		}, &proto.Player_FrostDeathKnight{FrostDeathKnight: &proto.FrostDeathKnight{Options: &proto.FrostDeathKnight_Options{ClassOptions: &proto.DeathKnightOptions{}}}}), nil, nil, nil)},
		{Name: "unholyDeathKnight", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDeathKnight,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "--22322321232331121231",
		}, &proto.Player_UnholyDeathKnight{UnholyDeathKnight: &proto.UnholyDeathKnight{Options: &proto.UnholyDeathKnight_Options{ClassOptions: &proto.DeathKnightOptions{}}}}), nil, nil, nil)},

		// Druid
		{Name: "balanceDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "33233221123212111231-2",
		}, &proto.Player_BalanceDruid{BalanceDruid: &proto.BalanceDruid{Options: &proto.BalanceDruid_Options{ClassOptions: &proto.DruidOptions{}}}}), nil, nil, nil)},
		{Name: "feralDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-2322322312012221222311",
		}, &proto.Player_FeralDruid{FeralDruid: &proto.FeralDruid{Options: &proto.FeralDruid_Options{ClassOptions: &proto.DruidOptions{}}, Rotation: &proto.FeralDruid_Rotation{}}}), nil, nil, nil)},
		// {Name: "guardianDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
		// 	Class:         proto.Class_ClassDruid,
		// 	Equipment:     &proto.EquipmentSpec{},
		// 	TalentsString: "-503232132322010353120300313511-20350001",
		// }, &proto.Player_FeralTankDruid{FeralTankDruid: &proto.FeralTankDruid{Options: &proto.FeralTankDruid_Options{ClassOptions: &proto.DruidOptions{}}}}), nil, nil, nil)},
		{Name: "restorationDruid", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassDruid,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-2322322312012221222311",
		}, &proto.Player_RestorationDruid{RestorationDruid: &proto.RestorationDruid{Options: &proto.RestorationDruid_Options{ClassOptions: &proto.DruidOptions{}}}}), nil, nil, nil)},

		// Hunter
		{Name: "beastMasteryHunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassHunter,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "003-2322321232122231221",
		}, &proto.Player_BeastMasteryHunter{BeastMasteryHunter: &proto.BeastMasteryHunter{Options: &proto.BeastMasteryHunter_Options{ClassOptions: &proto.HunterOptions{}}}}), nil, nil, nil)},
		{Name: "marksmanshipHunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassHunter,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-0320321232120131221-332002",
		}, &proto.Player_MarksmanshipHunter{MarksmanshipHunter: &proto.MarksmanshipHunter{Options: &proto.MarksmanshipHunter_Options{ClassOptions: &proto.HunterOptions{}}}}), nil, nil, nil)},
		{Name: "survivalHunter", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassHunter,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "--33222223123222121321",
		}, &proto.Player_SurvivalHunter{SurvivalHunter: &proto.SurvivalHunter{Options: &proto.SurvivalHunter_Options{ClassOptions: &proto.HunterOptions{}}}}), nil, nil, nil)},

		// Mage
		{Name: "arcaneMage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "323322221232122212121",
		}, &proto.Player_ArcaneMage{ArcaneMage: &proto.ArcaneMage{Options: &proto.ArcaneMage_Options{ClassOptions: &proto.MageOptions{}}}}), nil, nil, nil)},
		{Name: "fireMage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "001-232332221121121213231",
		}, &proto.Player_FireMage{FireMage: &proto.FireMage{Options: &proto.FireMage_Options{ClassOptions: &proto.MageOptions{}}}}), nil, nil, nil)},
		{Name: "frostMage", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassMage,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "--2323223213331321221",
		}, &proto.Player_FrostMage{FrostMage: &proto.FrostMage{Options: &proto.FrostMage_Options{ClassOptions: &proto.MageOptions{}}}}), nil, nil, nil)},

		// Paladin
		{Name: "holyPaladin", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "23332221222131312321",
		}, &proto.Player_HolyPaladin{HolyPaladin: &proto.HolyPaladin{Options: &proto.HolyPaladin_Options{ClassOptions: &proto.PaladinOptions{}}}}), nil, nil, nil)},
		{Name: "protPaladin", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "003-32223223122121121231",
		}, &proto.Player_ProtectionPaladin{ProtectionPaladin: &proto.ProtectionPaladin{Options: &proto.ProtectionPaladin_Options{ClassOptions: &proto.PaladinOptions{}}}}), nil, nil, nil)},
		{Name: "retPaladin", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPaladin,
			Race:          proto.Race_RaceBloodElf,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-3-23223213211113212321",
		}, &proto.Player_RetributionPaladin{RetributionPaladin: &proto.RetributionPaladin{Options: &proto.RetributionPaladin_Options{ClassOptions: &proto.PaladinOptions{}}}}), nil, nil, nil)},

		// Priest
		{Name: "disciplinePriest", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "233213121213222312221",
		}, &proto.Player_DisciplinePriest{DisciplinePriest: &proto.DisciplinePriest{Options: &proto.DisciplinePriest_Options{ClassOptions: &proto.PriestOptions{}}}}), nil, nil, nil)},
		{Name: "holyPriest", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "013-233122221211221123211",
		}, &proto.Player_HolyPriest{HolyPriest: &proto.HolyPriest{Options: &proto.HolyPriest_Options{ClassOptions: &proto.PriestOptions{}}}}), nil, nil, nil)},
		{Name: "shadowPriest", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassPriest,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-002-322232212211222121231",
		}, &proto.Player_ShadowPriest{ShadowPriest: &proto.ShadowPriest{Options: &proto.ShadowPriest_Options{ClassOptions: &proto.PriestOptions{}}}}), nil, nil, nil)},

		// Rogue
		{Name: "assassinationRogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassRogue,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "2333232213322112321",
		}, &proto.Player_AssassinationRogue{AssassinationRogue: &proto.AssassinationRogue{Options: &proto.AssassinationRogue_Options{ClassOptions: &proto.RogueOptions{}}}}), nil, nil, nil)},
		{Name: "combatRogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassRogue,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-2332232312232212321",
		}, &proto.Player_CombatRogue{CombatRogue: &proto.CombatRogue{Options: &proto.CombatRogue_Options{ClassOptions: &proto.RogueOptions{}}}}), nil, nil, nil)},
		{Name: "subtletyRogue", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassRogue,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "--2331232321313312321",
		}, &proto.Player_SubtletyRogue{SubtletyRogue: &proto.SubtletyRogue{Options: &proto.SubtletyRogue_Options{ClassOptions: &proto.RogueOptions{}}}}), nil, nil, nil)},

		// Shaman
		{Name: "elementalShaman", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "3230323212231121321-22",
		}, &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{ClassOptions: &proto.ShamanOptions{}}}}), nil, nil, nil)},
		{Name: "enhancementShaman", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-2331322313223212321",
		}, &proto.Player_EnhancementShaman{EnhancementShaman: &proto.EnhancementShaman{Options: &proto.EnhancementShaman_Options{ClassOptions: &proto.ShamanOptions{}}}}), nil, nil, nil)},
		{Name: "restorationShaman", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassShaman,
			Race:          proto.Race_RaceTroll,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "--23322232132123121321",
		}, &proto.Player_RestorationShaman{RestorationShaman: &proto.RestorationShaman{Options: &proto.RestorationShaman_Options{ClassOptions: &proto.ShamanOptions{}}}}), nil, nil, nil)},

		// Warlock
		{Name: "afflictionWarlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "223222123213321321-31",
		}, &proto.Player_AfflictionWarlock{AfflictionWarlock: &proto.AfflictionWarlock{Options: &proto.AfflictionWarlock_Options{ClassOptions: &proto.WarlockOptions{}}}}), nil, nil, nil)},
		{Name: "demonologyWarlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "013-3322222312312212211",
		}, &proto.Player_DemonologyWarlock{DemonologyWarlock: &proto.DemonologyWarlock{Options: &proto.DemonologyWarlock_Options{ClassOptions: &proto.WarlockOptions{}}}}), nil, nil, nil)},
		{Name: "destructionWarlock", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarlock,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-022-3322222312221312211",
		}, &proto.Player_DestructionWarlock{DestructionWarlock: &proto.DestructionWarlock{Options: &proto.DestructionWarlock_Options{ClassOptions: &proto.WarlockOptions{}}}}), nil, nil, nil)},

		// Warrior
		{Name: "armsWarrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "32222323122212312211-2",
		}, &proto.Player_ArmsWarrior{ArmsWarrior: &proto.ArmsWarrior{Options: &proto.ArmsWarrior_Options{ClassOptions: &proto.WarriorOptions{}}}}), nil, nil, nil)},
		{Name: "furyWarrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "002-332222131321111223211",
		}, &proto.Player_FuryWarrior{FuryWarrior: &proto.FuryWarrior{Options: &proto.FuryWarrior_Options{ClassOptions: &proto.WarriorOptions{}}}}), nil, nil, nil)},
		{Name: "protectionWarrior", Raid: core.SinglePlayerRaidProto(core.WithSpec(&proto.Player{
			Class:         proto.Class_ClassWarrior,
			Equipment:     &proto.EquipmentSpec{},
			TalentsString: "-002-33233221121212212231",
		}, &proto.Player_ProtectionWarrior{ProtectionWarrior: &proto.ProtectionWarrior{Options: &proto.ProtectionWarrior_Options{ClassOptions: &proto.WarriorOptions{}}}}), nil, nil, nil)},
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
