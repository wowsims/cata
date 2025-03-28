package database

import (
	"database/sql"
	"fmt"
	"strings"
)

// Tables
var RawItems []RawItemData
var RandPropPoints []RandPropPointsStruct
var ItemStatEffects []ItemStatEffect
var ItemStatEffectById map[int]ItemStatEffect
var RawGems []RawGem
var RawEnchants []RawEnchant

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
	classMask           int
	raceMask            int
	qualityModifier     float64 //In seemingly all cases this is bonus armor
}

func ScanRawItemData(rows *sql.Rows) (RawItemData, error) {
	var raw RawItemData
	err := rows.Scan(&raw.id, &raw.name, &raw.invType, &raw.itemDelay, &raw.overallQuality, &raw.dmgVariance,
		&raw.dbMinDamage, &raw.dbMaxDamage, &raw.itemLevel, &raw.itemClassName, &raw.itemSubClassName,
		&raw.rppEpic, &raw.rppSuperior, &raw.rppGood, &raw.statValue, &raw.bonusStat,
		&raw.clothArmorValue, &raw.leatherArmorValue, &raw.mailArmorValue, &raw.plateArmorValue,
		&raw.armorLocID, &raw.shieldArmorValues, &raw.statPercentEditor, &raw.socketTypes, &raw.socketEnchantmentId, &raw.flags0, &raw.FDID, &raw.itemSetName, &raw.itemSetId, &raw.flags1, &raw.classMask, &raw.raceMask, &raw.qualityModifier)
	return raw, err
}
func LoadRawItems(dbHelper *DBHelper, filter string) {
	baseQuery := `
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
			s.Flags_1 as Flags_1,
			s.AllowableClass as ClassMask,
			s.AllowableRace as RaceMask,
			s.QualityModifier
		FROM Item i
		JOIN ItemSparse s ON i.ID = s.ID
		JOIN ItemClass ic ON i.ClassID = ic.ClassID
		JOIN ItemSubClass isc ON i.ClassID = isc.ClassID AND i.SubClassID = isc.SubClassID
		JOIN RandPropPoints rpp ON s.ItemLevel = rpp.ID
		LEFT JOIN ArmorLocation al ON al.ID = ArmorLocationId
		LEFT JOIN ItemArmorShield ias ON s.ItemLevel = ias.ItemLevel
		LEFT JOIN ItemSet itemset ON s.ItemSet = itemset.ID
		JOIN ItemArmorTotal at ON s.ItemLevel = at.ItemLevel`

	// Filter string can be provided provided, we just append it to the query. For multiple conditions, the filter string should include real SQL ("s.ID = 78737" or "s.ItemLevel > 50 AND s.OverallQualityID = 4").
	if strings.TrimSpace(filter) != "" {
		baseQuery += " WHERE " + filter
	}

	// For debugging, you might want to see the complete query:
	fmt.Println("Executing query:", baseQuery)

	// LoadRows is assumed to be a function that executes the query and maps the results using ScanRawItemData.
	items, err := LoadRows(dbHelper.db, baseQuery, ScanRawItemData)
	fmt.Println("Loaded Items:", len(items))
	if err != nil {
		fmt.Println("Error loading items:", err.Error())
		return
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

// ItemDamage tables
type RawGem struct {
	ItemId       int
	Name         string
	FDID         int
	GemType      int
	Effect       []int
	StatList     []int
	StatBonus    []int
	MinItemLevel int
	Quality      int
	IsJc         bool
	Flags        ItemFlags
}

func ScanGemTable(rows *sql.Rows) (RawGem, error) {
	var raw RawGem
	var statListString string
	var statBonusString string
	var effectString string
	err := rows.Scan(&raw.ItemId, &raw.Name, &raw.FDID, &statListString, statBonusString, &raw.MinItemLevel, &raw.Quality, effectString, &raw.IsJc, &raw.Flags)
	raw.StatList, err = parseIntArrayField(statListString, 3)
	if err != nil {
		fmt.Errorf("Error loading GemTable statListString")
	}
	raw.StatBonus, err = parseIntArrayField(statBonusString, 3)
	if err != nil {
		fmt.Errorf("Error loading GemTable statBonusString")
	}
	raw.Effect, err = parseIntArrayField(effectString, 3)
	if err != nil {
		fmt.Errorf("Error loading GemTable effectString")
	}
	return raw, err
}

func LoadRawGems(dbHelper *DBHelper) error {
	query := `SELECT
		s.ID,
		s.Display_lang as Name,
		i.IconFileDataID as FDID,
		gp.'Type' as GemType,
		sie.EffectPointsMax as StatList,
		sie.EffectArg as StatBonus,
		gp.Min_item_level MinItemLevel,
		s.OverallQualityId Quality,
		s.Effect,
		CASE 
			WHEN s.RequiredSkill = 755 THEN 1
			ELSE 0
		END AS IsJc,
		s.Flags_0
		FROM ItemSparse s 
		JOIN Item i ON s.ID = i.ID
		JOIN GemProperties gp ON s.Gem_properties = gp.ID
		JOIN SpellItemEnchantment sie ON gp.Enchant_ID = sie.ID
		WHERE i.ClassID = 3`
	items, err := LoadRows(dbHelper.db, query, ScanGemTable)
	if err != nil {
		return fmt.Errorf("error loading items for GemTables: %w", err)
	}

	RawGems = items

	return nil
}

type RawEnchant struct {
	EffectId        int
	Name            string
	SpellId         int
	ItemId          int
	ProfessionId    int
	Effects         []int
	EffectPoints    []int
	EffectArgs      []int
	IsWeaponEnchant bool
	InvTypesMask    int
	SubClassMask    int
	ClassMask       int
}

func ScanEnchantsTable(rows *sql.Rows) (RawEnchant, error) {
	var raw RawEnchant
	var effectsString string
	var effectPointsString string
	var effectArgsString string
	err := rows.Scan(&raw.EffectId, &raw.Name, &raw.SpellId, &raw.ItemId, &raw.ProfessionId, &effectsString, &effectPointsString, &effectArgsString, &raw.IsWeaponEnchant, &raw.InvTypesMask, &raw.SubClassMask, &raw.ClassMask)
	raw.Effects, err = parseIntArrayField(effectsString, 3)
	if err != nil {
		fmt.Errorf("Error loading GemTable effectsString")
	}
	raw.EffectPoints, err = parseIntArrayField(effectPointsString, 3)
	if err != nil {
		fmt.Errorf("Error loading GemTable effectPointsString")
	}
	raw.EffectArgs, err = parseIntArrayField(effectArgsString, 3)
	if err != nil {
		fmt.Errorf("Error loading GemTable effectArgsString")
	}
	return raw, err
}

func LoadRawEnchants(dbHelper *DBHelper) error {
	query := `SELECT
		se.ID as effectId,
		sn.Name_lang as name,
		se.SpellID as spellId,
		ie.ParentItemID as ItemId,
		sie.Field_1_15_3_55112_014 as professionId,
		sie.Effect as Effect,
		sie.EffectPointsMax as EffectPoints,
		sie.EffectArg as EffectArgs,
		CASE 
			WHEN sei.EquippedItemClass = 4 THEN 1
			ELSE 0
		END AS isWeaponEnchant,
		sei.EquippedItemInvTypes as InvTypes,
		sei.EquippedItemSubclass,
		sla.ClassMask
		FROM SpellEffect se 
		JOIN Spell s ON se.SpellID = s.ID
		JOIN SpellName sn ON se.SpellID = sn.ID
		JOIN SpellItemEnchantment sie ON se.EffectMiscValue_0 = sie.ID
		LEFT JOIN ItemEffect ie ON se.SpellID = ie.SpellID
		LEFT JOIN SpellEquippedItems sei ON se.SpellId = sei.SpellID
		LEFT JOIN SkillLineAbility sla ON se.SpellID = sla.Spell
		WHERE se.Effect = 53`
	items, err := LoadRows(dbHelper.db, query, ScanEnchantsTable)
	if err != nil {
		return fmt.Errorf("error loading items for GemTables: %w", err)
	}

	RawEnchants = items

	return nil
}
