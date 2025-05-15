package database

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/tools"
)

func ReadAtlasLootData(dbHelper *DBHelper) *WowDatabase {
	db := NewWowDatabase()

	// Read these in reverse order, because some items are listed in multiple expansions
	// and we want to overwrite with the earliest value.
	readAtlasLootSourceData(db, proto.Expansion_ExpansionWotlk, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_Data/source-wrath.lua")
	readAtlasLootSourceData(db, proto.Expansion_ExpansionTbc, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_Data/source-tbc.lua")
	readAtlasLootSourceData(db, proto.Expansion_ExpansionVanilla, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_Data/source.lua")

	readAtlasLootDungeonData(db, proto.Expansion_ExpansionVanilla, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_DungeonsAndRaids/data.lua")
	readAtlasLootDungeonData(db, proto.Expansion_ExpansionTbc, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_DungeonsAndRaids/data-tbc.lua")
	readAtlasLootDungeonData(db, proto.Expansion_ExpansionWotlk, "https://raw.githubusercontent.com/Hoizame/AtlasLootClassic/master/AtlasLootClassic_DungeonsAndRaids/data-wrath.lua")

	// Cata addon
	readAtlasLootSourceData(db, proto.Expansion_ExpansionCata, "https://raw.githubusercontent.com/snowflame0/AtlasLootClassic_Cata/main/AtlasLootClassic_Data/source-cata.lua")
	readAtlasLootDungeonData(db, proto.Expansion_ExpansionCata, "https://raw.githubusercontent.com/snowflame0/AtlasLootClassic_Cata/main/AtlasLootClassic_DungeonsAndRaids/data-cata.lua")
	readAtlasLootFactionData(db, "https://raw.githubusercontent.com/snowflame0/AtlasLootClassic_Cata/main/AtlasLootClassic_Factions/data-cata.lua")

	readZoneData(db, dbHelper)

	return db
}

func readAtlasLootSourceData(db *WowDatabase, expansion proto.Expansion, srcUrl string) {
	srcTxt, err := tools.ReadWeb(srcUrl)
	if err != nil {
		log.Fatalf("Error reading atlasloot file %s", err)
	}

	itemPattern := regexp.MustCompile(`^\[([0-9]+)\] = {(.*)},$`)
	typePattern := regexp.MustCompile(`\[3\] = (\d+),.*\[4\] = (\d+)`)
	lines := strings.Split(srcTxt, "\n")
	for _, line := range lines {
		match := itemPattern.FindStringSubmatch(line)
		if match != nil {
			idStr := match[1]
			id, _ := strconv.Atoi(idStr)
			item := &proto.UIItem{Id: int32(id), Expansion: expansion}
			if _, ok := db.Items[item.Id]; ok {
				continue
			}

			paramsStr := match[2]
			typeMatch := typePattern.FindStringSubmatch(paramsStr)
			if typeMatch != nil {
				itemType, _ := strconv.Atoi(typeMatch[1])
				spellID, _ := strconv.Atoi(typeMatch[2])
				if prof, ok := AtlasLootProfessionIDs[itemType]; ok {
					item.Sources = append(item.Sources, &proto.UIItemSource{
						Source: &proto.UIItemSource_Crafted{
							Crafted: &proto.CraftedSource{
								Profession: prof,
								SpellId:    int32(spellID),
							},
						},
					})
				}
			}

			db.MergeItem(item)
		}
	}
}

func readAtlasLootDungeonData(db *WowDatabase, expansion proto.Expansion, srcUrl string) {
	srcTxt, err := tools.ReadWeb(srcUrl)
	if err != nil {
		log.Fatalf("Error reading atlasloot file %s", err)
	}

	// Substitute the raw string for END_TIME_ECHO_LOOT wherever it is referenced in the lua source.
	srcTxt = strings.ReplaceAll(srcTxt, "END_TIME_ECHO_LOOT,", EndTimeEchoLootString)

	// Convert newline to '@@@' so we can do regexes on the whole file as 1 line.
	srcTxt = strings.ReplaceAll(srcTxt, "\n", "@@@")

	dungeonPattern := regexp.MustCompile(`data\["([^"]+)"] = {.*?\sMapID = (\d+),.*?ContentType = ([^"]+),.*?items = {(.*?)@@@}@@@`)
	npcNameAndIDPattern := regexp.MustCompile(`^[^@]*?AL\["(.*?)"\]\)?,(.*?(@@@\s*npcID = {?(\d+),))?`)
	diffItemsPattern := regexp.MustCompile(`\[([A-Z0-9]+_DIFF)\] = (({.*?@@@\s*},?@@@)|(.*?@@@\s*\),?@@@))`)
	itemsPattern := regexp.MustCompile(`@@@\s+{(.*?)},`)
	itemParamPattern := regexp.MustCompile(`AL\["(.*?)"\]`)
	for _, dungeonMatch := range dungeonPattern.FindAllStringSubmatch(srcTxt, -1) {
		fmt.Printf("Zone: %s\n", dungeonMatch[1])
		zoneID, _ := strconv.Atoi(dungeonMatch[2])
		db.MergeZone(&proto.UIZone{
			Id:        int32(zoneID),
			Expansion: expansion,
		})
		contentType := dungeonMatch[3]

		npcSplits := strings.Split(dungeonMatch[4], "name = ")[1:]
		for _, npcSplit := range npcSplits {
			npcSplit = strings.ReplaceAll(npcSplit, "AtlasLoot:GetRetByFaction(", "")
			npcMatch := npcNameAndIDPattern.FindStringSubmatch(npcSplit)
			if npcMatch == nil {
				panic("No npc match: " + npcSplit)
			}
			npcName := npcMatch[1]
			npcID := 0
			if len(npcMatch) > 3 {
				npcID, _ = strconv.Atoi(npcMatch[4])
			}
			if npcName == "Onyxia" { // AtlasLoot uses 15956 for some reason, which is the ID for Anub'Rekan.
				npcID = 10184
			} else if npcName == "Yogg-Saron" { // AtlasLoot uses 33271 for some reason, which is the ID for General Vezax.
				npcID = 33288
			}
			fmt.Printf("NPC: %s/%d\n", npcName, npcID)
			if npcID != 0 {
				db.MergeNpc(&proto.UINPC{
					Id:     int32(npcID),
					ZoneId: int32(zoneID),
					Name:   npcName,
				})
			}

			for _, difficultyMatch := range diffItemsPattern.FindAllStringSubmatch(npcSplit, -1) {
				diffString := difficultyMatch[1]
				if expansion == proto.Expansion_ExpansionCata && contentType == "RAID_CONTENT" {
					diffString = AtlasLootDungeonToRaidDifficulty[diffString]
				}
				difficulty, ok := AtlasLootDifficulties[diffString]
				if !ok {
					log.Fatalf("Invalid difficulty for NPC %s: %s", npcName, diffString)
				}

				curCategory := ""
				curLocation := 0

				for _, itemMatch := range itemsPattern.FindAllStringSubmatch(difficultyMatch[0], -1) {
					itemParams := core.MapSlice(strings.Split(itemMatch[1], ","), strings.TrimSpace)
					location, _ := strconv.Atoi(itemParams[0]) // Location within AtlasLoot's menu.

					idStr := itemParams[1]
					if idStr[0] == 'n' || idStr[0] == '"' { // nil or "xxx"
						if len(itemParams) > 3 {
							if paramMatch := itemParamPattern.FindStringSubmatch(itemParams[3]); paramMatch != nil {
								curCategory = paramMatch[1]
								curLocation = location
							}
						}
						if len(itemParams) > 4 {
							if paramMatch := itemParamPattern.FindStringSubmatch(itemParams[4]); paramMatch != nil {
								curCategory = paramMatch[1]
								curLocation = location
							}
						}
					} else { // item ID
						itemID, _ := strconv.Atoi(idStr)
						fmt.Printf("Item: %d\n", itemID)
						dropSource := &proto.DropSource{
							Difficulty: difficulty,
							ZoneId:     int32(zoneID),
						}
						if npcID == 0 {
							dropSource.OtherName = npcName
						} else {
							dropSource.NpcId = int32(npcID)
						}

						if curCategory != "" && location == curLocation+1 {
							curLocation = location
							dropSource.Category = curCategory
						}

						item := &proto.UIItem{Id: int32(itemID), Sources: []*proto.UIItemSource{{
							Source: &proto.UIItemSource_Drop{
								Drop: dropSource,
							},
						}}}
						db.MergeItem(item)
					}
				}
			}
		}
	}
}

func readAtlasLootFactionData(db *WowDatabase, srcUrl string) {
	srcTxt, err := tools.ReadWeb(srcUrl)
	if err != nil {
		log.Fatalf("Error reading atlasloot file %s", err)
	}

	// Convert newline to '@@@' so we can do regexes on the whole file as 1 line.
	srcTxt = strings.ReplaceAll(srcTxt, "\n", "@@@")

	factionPattern := regexp.MustCompile(`data\["([^"]+)"] = {.*?\sFactionID = (\d+),.*?ContentType = ([^,]+),.*?items = {(.*?)@@@}`)
	levelPattern := regexp.MustCompile(`^[^@]*?ALIL\["(.*?)"\]\)?,(.*?(@@@\s*npcID = {?(\d+),))?`)
	diffItemsPattern := regexp.MustCompile(`\[([A-Z0-9]+_DIFF)\] = (({.*?@@@\s*},?@@@)|(.*?@@@\s*\),?@@@))`)
	itemsPattern := regexp.MustCompile(`@@@\s+{(.*?)},`)
	for _, factionMatch := range factionPattern.FindAllStringSubmatch(srcTxt, -1) {
		fmt.Printf("Faction: %s\n", factionMatch[1])
		if factionMatch[1] == "DUMMY" {
			continue
		}

		factionID, _ := strconv.Atoi(factionMatch[2])
		contentType := factionMatch[3]
		faction, ok := AtlasLootFactions[contentType]
		if !ok {
			log.Fatalf("Invalid faction for Content Type %s", contentType)
		}

		npcSplits := strings.Split(factionMatch[4], "name = ")[1:]
		for _, levelSplit := range npcSplits {
			levelMatch := levelPattern.FindStringSubmatch(levelSplit)
			if levelMatch == nil {
				panic("No level match: " + levelSplit)
			}
			levelName := levelMatch[1]
			fmt.Printf("Level: %s\n", levelName)

			repLevel, ok := AtlasLootRepLevel[levelName]
			if !ok {
				log.Fatalf("Invalid Rep Level for %s", levelName)
			}

			for _, difficultyMatch := range diffItemsPattern.FindAllStringSubmatch(levelSplit, -1) {
				for _, itemMatch := range itemsPattern.FindAllStringSubmatch(difficultyMatch[0], -1) {
					itemParams := core.MapSlice(strings.Split(itemMatch[1], ","), strings.TrimSpace)

					idStr := itemParams[1]
					if idStr[0] == 'n' || idStr[0] == '"' { // nil or "xxx"
					} else { // item ID
						itemID, _ := strconv.Atoi(idStr)
						fmt.Printf("Item: %d\n", itemID)
						factionSource := &proto.RepSource{
							FactionId:    faction,
							RepFactionId: proto.RepFaction(factionID),
							RepLevel:     repLevel,
						}

						item := &proto.UIItem{Id: int32(itemID), FactionRestriction: AtlasLootFactionRestrictions[faction], Sources: []*proto.UIItemSource{{
							Source: &proto.UIItemSource_Rep{
								Rep: factionSource,
							},
						}}}
						db.MergeItem(item)
					}
				}
			}
		}
	}
}

// func readAtlasLootCraftingData(db *WowDatabase, srcUrl string) {
// 	srcTxt, err := tools.ReadWeb(srcUrl)
// 	if err != nil {
// 		log.Fatalf("Error reading atlasloot file %s", err)
// 	}

// 	// Convert newline to '@@@' so we can do regexes on the whole file as 1 line.
// 	srcTxt = strings.ReplaceAll(srcTxt, "\n", "@@@")

// 	craftingPattern := regexp.MustCompile(`data\["([^"]+)"] = {.*?\name = ALIL["(\d+)"],.*?ContentType = ([^,]+),.*?items = {(.*?)@@@}`)
// 	diffItemsPattern := regexp.MustCompile(`\[([A-Z0-9]+_DIFF)\] = (({.*?@@@\s*},?@@@)|(.*?@@@\s*\),?@@@))`)
// 	itemsPattern := regexp.MustCompile(`@@@\s+{(.*?)},`)
// 	for _, craftingMatch := range craftingPattern.FindAllStringSubmatch(srcTxt, -1) {
// 		fmt.Printf("Profession: %s\n", craftingMatch[1])

// 		profession, ok := AtlasLootProfessionNames[craftingMatch[1]]
// 		if !ok {
// 			log.Fatalf("Invalid Profession for %s", craftingMatch[1])
// 		}

// 		npcSplits := strings.Split(craftingMatch[4], "name = ")[1:]
// 		for _, levelSplit := range npcSplits {
// 			for _, difficultyMatch := range diffItemsPattern.FindAllStringSubmatch(levelSplit, -1) {
// 				for _, itemMatch := range itemsPattern.FindAllStringSubmatch(difficultyMatch[0], -1) {
// 					itemParams := core.MapSlice(strings.Split(itemMatch[1], ","), strings.TrimSpace)

// 					idStr := itemParams[1]
// 					if idStr[0] == 'n' || idStr[0] == '"' { // nil or "xxx"
// 					} else { // craft ID
// 						craftID, _ := strconv.Atoi(idStr)
// 						fmt.Printf("CraftID: %d\n", craftID)
// 						craftedSource := &proto.CraftedSource{
// 							Profession: profession,
// 							SpellId: int32(craftID),
// 						}

// 						itemID := int32(0)

// 						// Dont add new items from craft ids
// 						if _, ok := db.Items[itemID]; ok {
// 							item := &proto.UIItem{Id: int32(itemID), Sources: []*proto.UIItemSource{{
// 								Source: &proto.UIItemSource_Crafted{
// 									Crafted: craftedSource,
// 								},
// 							}}}
// 							db.MergeItem(item)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func readZoneData(db *WowDatabase, dbHelper *DBHelper) {
	zoneIDs := make([]int32, 0, len(db.Zones))
	for zoneID := range db.Zones {
		zoneIDs = append(zoneIDs, zoneID)
	}

	zoneNames, error := loadZones(dbHelper)
	if error != nil {
		panic(error)
	}

	for _, zoneID := range zoneIDs {
		db.Zones[zoneID].Name = zoneNames[zoneID]
	}
}

func loadZones(dbHelper *DBHelper) (map[int32]string, error) {
	const query = `SELECT ID, AreaName_lang FROM AreaTable`

	rows, err := dbHelper.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying drop sources: %w", err)
	}
	defer rows.Close()

	var zoneId int32
	var zoneName string
	namesByZone := make(map[int32]string)
	for rows.Next() {
		e := rows.Scan(&zoneId, &zoneName)
		if e != nil {
			return nil, e
		}

		namesByZone[zoneId] = zoneName
	}

	return namesByZone, nil
}

var AtlasLootProfessionIDs = map[int]proto.Profession{
	//4: proto.Profession_FirstAid,
	5: proto.Profession_Blacksmithing,
	6: proto.Profession_Leatherworking,
	7: proto.Profession_Alchemy,
	//9: proto.Profession_Cooking,
	10: proto.Profession_Mining,
	11: proto.Profession_Tailoring,
	12: proto.Profession_Engineering,
	13: proto.Profession_Enchanting,
	17: proto.Profession_Jewelcrafting,
	18: proto.Profession_Inscription,
}
var AtlasLootFactions = map[string]proto.Faction{
	"FACTIONS_CONTENT":       proto.Faction_Unknown,
	"FACTIONS_ALLI_CONTENT":  proto.Faction_Alliance,
	"FACTIONS_HORDE_CONTENT": proto.Faction_Horde,
}
var AtlasLootFactionRestrictions = map[proto.Faction]proto.UIItem_FactionRestriction{
	proto.Faction_Unknown:  proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED,
	proto.Faction_Alliance: proto.UIItem_FACTION_RESTRICTION_ALLIANCE_ONLY,
	proto.Faction_Horde:    proto.UIItem_FACTION_RESTRICTION_HORDE_ONLY,
}
var AtlasLootRepLevel = map[string]proto.RepLevel{
	"Exalted":  proto.RepLevel_RepLevelExalted,
	"Revered":  proto.RepLevel_RepLevelRevered,
	"Honored":  proto.RepLevel_RepLevelHonored,
	"Friendly": proto.RepLevel_RepLevelFriendly,
}
var AtlasLootDifficulties = map[string]proto.DungeonDifficulty{
	"NORMAL_DIFF":   proto.DungeonDifficulty_DifficultyNormal,
	"HEROIC_DIFF":   proto.DungeonDifficulty_DifficultyHeroic,
	"ALPHA_DIFF":    proto.DungeonDifficulty_DifficultyTitanRuneAlpha,
	"BETA_DIFF":     proto.DungeonDifficulty_DifficultyTitanRuneBeta,
	"RAID10_DIFF":   proto.DungeonDifficulty_DifficultyRaid10,
	"RAID10H_DIFF":  proto.DungeonDifficulty_DifficultyRaid10H,
	"RAID25RF_DIFF": proto.DungeonDifficulty_DifficultyRaid25RF,
	"RAID25_DIFF":   proto.DungeonDifficulty_DifficultyRaid25,
	"RAID25H_DIFF":  proto.DungeonDifficulty_DifficultyRaid25H,
}
var AtlasLootDungeonToRaidDifficulty = map[string]string{
	"RF_DIFF":     "RAID25RF_DIFF",
	"NORMAL_DIFF": "RAID25_DIFF",
	"HEROIC_DIFF": "RAID25H_DIFF",
}

const EndTimeEchoLootString = `{
    { 1, "INV_Box_01", nil, AL["Echo of Baine"], nil },	--Echo of Baine
    { 2, 72815 },	-- Bloodhoof Legguards
    { 3, 72814 },	-- Axe of the Tauren Chieftains
    { 4, "INV_Box_01", nil, AL["Echo of Jaina"], nil },	--Echo of Jaina
    { 5, 72808 },	-- Jaina's Staff
    { 6, 72809 },	-- Ward of Incantations
    { 7, "INV_Box_01", nil, AL["Echo of Sylvanas"], nil },	--Echo of Sylvanas
    { 8, 72811 },	-- Cloak of the Banshee Queen
    { 9, 72810 },	-- Windrunner's Bow
    { 10, "ac6130" },
    { 11, "INV_Box_01", nil, AL["Echo of Tyrande"], nil },	--Echo of Tyrande
    { 12, 72813 },	-- Whisperwind Robes
    { 13, 72812 },	-- Crescent Moon
    { 14, "ac5995" },
    { 16, "INV_Box_01", nil, AL["Shared"], nil },	--Shared
    { 17, 72802 },	-- Time Traveler's Leggings
    { 18, 72805 },	-- Gloves of the Hollow
    { 19, 72798 },	-- Cord of Lost Hope
    { 20, 72806 },	-- Echoing Headguard
    { 21, 72799 },	-- Dead End Boots
    { 22, 72801 },	-- Breastplate of Sorrow
    { 23, 72800 },	-- Gauntlets of Temporal Interference
    { 24, 72803 },	-- Girdle of Lost Heroes
    { 25, 72807 },	-- Waistguard of Lost Time
    { 26, 72804 },	-- Dragonshrine Scepter
},`
