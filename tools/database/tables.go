package database

import (
	"database/sql"
	"fmt"
)

// Tables
var RawItems []RawItemData
var RandPropPoints []RandPropPointsStruct
var ItemStatEffects []ItemStatEffect
var ItemStatEffectById map[int]ItemStatEffect

// Loading tables
// Below is the definition and loading of tables
//

// Raw Item Data
//
//

type RawItemData struct {
	id                  int
	name                string
	invType             int
	itemDelay           int
	overallQuality      int
	dmgVariance         float64
	dbMinDamage         string
	dbMaxDamage         string
	itemLevel           int
	itemClassName       string
	itemSubClassName    string
	rppEpic             string
	rppSuperior         string
	rppGood             string
	statValue           string
	bonusStat           string
	clothArmorValue     float64
	leatherArmorValue   float64
	mailArmorValue      float64
	plateArmorValue     float64
	armorLocID          int
	shieldArmorValues   string
	statPercentEditor   string
	socketTypes         string
	socketEnchantmentId int
	flags0              ItemFlags
	flags1              ItemFlags
	FDID                int
	itemSetName         string
	itemSetId           int
}

func ScanRawItemData(rows *sql.Rows) (RawItemData, error) {
	var raw RawItemData
	err := rows.Scan(&raw.id, &raw.name, &raw.invType, &raw.itemDelay, &raw.overallQuality, &raw.dmgVariance,
		&raw.dbMinDamage, &raw.dbMaxDamage, &raw.itemLevel, &raw.itemClassName, &raw.itemSubClassName,
		&raw.rppEpic, &raw.rppSuperior, &raw.rppGood, &raw.statValue, &raw.bonusStat,
		&raw.clothArmorValue, &raw.leatherArmorValue, &raw.mailArmorValue, &raw.plateArmorValue,
		&raw.armorLocID, &raw.shieldArmorValues, &raw.statPercentEditor, &raw.socketTypes, &raw.socketEnchantmentId, &raw.flags0, &raw.FDID, &raw.itemSetName, &raw.itemSetId, &raw.flags1)
	return raw, err
}

func LoadRawItems(dbHelper *DBHelper) {
	query := `
		SELECT 
			i.ID, 
			s.Display_lang AS Name,
			i.InventoryType,
			s.ItemDelay,
			s.OverallQualityID,
			s.DmgVariance,
			s.MinDamage, 
			s.MaxDamage,
			s.ItemLevel,
			ic.ClassName_lang AS ItemClassName,
			isc.VerboseName_lang AS ItemSubClassName,
			rpp.Epic as RPPEpic,
			rpp.Superior as RPPSuperior,
			rpp.Good as RPPGood,
			s.Field_1_15_3_55112_014 as StatValue,
			s.StatModifier_bonusStat as bonusStat,
			COALESCE(at.Cloth, 0) * COALESCE(al.Clothmodifier, 1) AS clothArmorValue,
			COALESCE(at.Leather, 0) * COALESCE(al.LeatherModifier, 1) AS leatherArmorValue,
			COALESCE(at.Mail, 0) * COALESCE(al.Chainmodifier, 1) AS mailArmorValue,
			COALESCE(at.Plate, 0) * COALESCE(al.Platemodifier, 1) AS plateArmorValue,
			CASE 
				WHEN s.InventoryType = 20 THEN 5 
				ELSE s.InventoryType 
			END AS ArmorLocationID,
			ias.Quality as shieldArmorValues,
			s.StatPercentEditor as StatPercentEditor,
			s.SocketType as SocketTypes,
			s.Socket_match_enchantment_ID as SocketEnchantmentId,
			s.Flags_0 as Flags_0,
			i.IconFileDataId as FDID,
			COALESCE(itemset.Name_lang, '') as ItemSetName,
			COALESCE(itemset.ID, 0) as ItemSetID,
			s.Flags_1 as Flags_1
		FROM Item i
		JOIN ItemSparse s ON i.ID = s.ID
		JOIN ItemClass ic ON i.ClassID = ic.ClassID
		JOIN ItemSubClass isc ON i.ClassID = isc.ClassID AND i.SubClassID = isc.SubClassID
		JOIN RandPropPoints rpp ON s.ItemLevel = rpp.ID
		LEFT JOIN ArmorLocation al ON al.ID = ArmorLocationId
		LEFT JOIN ItemArmorShield ias ON s.ItemLevel = ias.ItemLevel
		LEFT JOIN ItemSet itemset ON s.ItemSet = itemset.ID
		JOIN ItemArmorTotal at ON s.ItemLevel = at.ItemLevel;
	`
	//WHERE s.ID = 78737
	items, err := LoadRows(dbHelper.db, query, ScanRawItemData)
	fmt.Println("Loaded Items", len(items))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Errorf("Error in query load items")
	}
	RawItems = items
}

// RandPropPoints
//
//
//

type RandPropPointsStruct struct {
	ItemLevel int
	Epic      []int
	Superior  []int
	Good      []int
}

func ScanRandPropPoints(rows *sql.Rows) (RandPropPointsStruct, error) {
	var raw RandPropPointsStruct
	var epicString, superiorString, goodString string

	err := rows.Scan(&raw.ItemLevel, &epicString, &superiorString, &goodString)
	raw.Epic, err = parseIntArrayField(epicString, 5)
	if err != nil {
		fmt.Errorf("Error loading items")
	}
	raw.Superior, err = parseIntArrayField(superiorString, 5)
	if err != nil {
		fmt.Errorf("Error loading items")
	}
	raw.Good, err = parseIntArrayField(goodString, 5)
	if err != nil {
		fmt.Errorf("Error loading items")
	}
	return raw, err
}

func LoadRandPropPoints(dbHelper *DBHelper) {
	query := `SELECT ID as ItemLevel, Epic, Superiors, Good FROM RandPropPoints;`
	items, err := LoadRows(dbHelper.db, query, ScanRandPropPoints)
	if err != nil {
		fmt.Errorf("Error in query load items")
	}

	RandPropPoints = items
}

//ItemStatEffects
// Used for straight up item stat effects from SpellItemEnchantment (socket bonuses for now, single stat)
//

type ItemStatEffect struct {
	ID              int
	EffectPointsMin []int
	EffectPointsMax []int
	EffectArg       []int
}

func ScanItemStatEffects(rows *sql.Rows) (ItemStatEffect, error) {
	var raw ItemStatEffect
	var ePointsMin, epointsMax, eArgs string
	err := rows.Scan(&raw.ID, &ePointsMin, &epointsMax, &eArgs)
	raw.EffectPointsMin, err = parseIntArrayField(ePointsMin, 3)
	if err != nil {
		fmt.Errorf("Error loading ItemStatEffects ePointsMin")
	}
	raw.EffectPointsMax, err = parseIntArrayField(epointsMax, 3)
	if err != nil {
		fmt.Errorf("Error loading ItemStatEffects epointsMax")
	}
	raw.EffectArg, err = parseIntArrayField(eArgs, 3)
	if err != nil {
		fmt.Errorf("Error loading ItemStatEffects eArgs")
	}
	return raw, err
}

func LoadItemStatEffects(dbHelper *DBHelper) {
	query := `SELECT ID, EffectPointsMin, EffectPointsMax, EffectArg FROM SpellItemEnchantment WHERE Effect_0 = 5`
	items, err := LoadRows(dbHelper.db, query, ScanItemStatEffects)
	if err != nil {
		fmt.Errorf("Error in query load items")
	}

	ItemStatEffects = items
	ItemStatEffectById = CacheBy(ItemStatEffects, func(effect ItemStatEffect) int {
		return effect.ID
	})
}

// ItemDamage tables
type ItemDamageTable struct {
	ItemLevel int
	Quality   []float64
}

func ScanItemDamageTable(rows *sql.Rows) (ItemDamageTable, error) {
	var raw ItemDamageTable
	var qualityString string
	err := rows.Scan(&raw.ItemLevel, &qualityString)
	raw.Quality, err = parseFloatArrayField(qualityString, 7)
	if err != nil {
		fmt.Errorf("Error loading ItemDamageTable qualityString")
	}
	return raw, err
}

var ItemDamageByTableAndItemLevel = make(map[string]map[int]ItemDamageTable)
var itemDamageTableNames = []string{
	"ItemDamageAmmo",
	"ItemDamageOneHand",
	"ItemDamageOneHandCaster",
	"ItemDamageRanged",
	"ItemDamageTwoHand",
	"ItemDamageTwoHandCaster",
	"ItemDamageThrown",
	"ItemDamageWand",
}

func LoadItemDamageTables(dbHelper *DBHelper) error {
	for _, tableName := range itemDamageTableNames {
		query := fmt.Sprintf("SELECT ItemLevel, Quality FROM %s", tableName)
		items, err := LoadRows(dbHelper.db, query, ScanItemDamageTable)
		if err != nil {
			return fmt.Errorf("error loading items for table %s: %w", tableName, err)
		}

		// Cache the slice of ItemDamageTable into a map keyed by ItemLevel.
		ItemDamageByTableAndItemLevel[tableName] = CacheBy(items, func(table ItemDamageTable) int {
			return table.ItemLevel
		})
	}
	return nil
}
