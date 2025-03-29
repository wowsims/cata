package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/wowsims/cata/sim/core/proto"
)

// Tables
var RawItems []RawItemData
var RandPropPoints []RandPropPointsStruct
var ItemStatEffects []ItemStatEffect
var ItemStatEffectById map[int]ItemStatEffect
var RawGems []RawGem
var RawEnchants []RawEnchant
var RawSpellEffectBySpellIdAndIndex map[int]map[int]RawSpellEffect

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
	RandomSuffixOptions []int32
}

func ScanRawItemData(rows *sql.Rows) (RawItemData, error) {
	var raw RawItemData
	var randomSuffixOptions sql.NullString
	err := rows.Scan(&raw.id, &raw.name, &raw.invType, &raw.itemDelay, &raw.overallQuality, &raw.dmgVariance,
		&raw.dbMinDamage, &raw.dbMaxDamage, &raw.itemLevel, &raw.itemClassName, &raw.itemSubClassName,
		&raw.rppEpic, &raw.rppSuperior, &raw.rppGood, &raw.statValue, &raw.bonusStat,
		&raw.clothArmorValue, &raw.leatherArmorValue, &raw.mailArmorValue, &raw.plateArmorValue,
		&raw.armorLocID, &raw.shieldArmorValues, &raw.statPercentEditor, &raw.socketTypes, &raw.socketEnchantmentId, &raw.flags0, &raw.FDID, &raw.itemSetName, &raw.itemSetId, &raw.flags1, &raw.classMask, &raw.raceMask, &raw.qualityModifier, &randomSuffixOptions)
	raw.RandomSuffixOptions = ParseRandomSuffixOptions(randomSuffixOptions)
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
			s.QualityModifier,
			(
			SELECT group_concat(-ench, ',')
				FROM item_enchantment_template
				WHERE entry = s.ItemRandomSuffixGroupID
			) AS RandomSuffixOptions
		FROM Item i
		JOIN ItemSparse s ON i.ID = s.ID
		JOIN ItemClass ic ON i.ClassID = ic.ClassID
		JOIN ItemSubClass isc ON i.ClassID = isc.ClassID AND i.SubClassID = isc.SubClassID
		JOIN RandPropPoints rpp ON s.ItemLevel = rpp.ID
		LEFT JOIN ArmorLocation al ON al.ID = ArmorLocationId
		LEFT JOIN ItemArmorShield ias ON s.ItemLevel = ias.ItemLevel
		LEFT JOIN ItemSet itemset ON s.ItemSet = itemset.ID
		JOIN ItemArmorTotal at ON s.ItemLevel = at.ItemLevel
		`

	if strings.TrimSpace(filter) != "" {
		baseQuery += " WHERE " + filter
	}

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
	query := `SELECT ID as ItemLevel, Epic, Superior, Good FROM RandPropPoints;`
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
	err := rows.Scan(&raw.ItemId, &raw.Name, &raw.FDID, &raw.GemType, &statListString, &statBonusString, &raw.MinItemLevel, &raw.Quality, &effectString, &raw.IsJc, &raw.Flags)
	raw.StatList, err = parseIntArrayField(statListString, 3)
	if err != nil {
		fmt.Println(err.Error(), 3, statListString)
		fmt.Errorf("Error loading GemTable statListString")
	}
	raw.StatBonus, err = parseIntArrayField(statBonusString, 3)
	if err != nil {
		fmt.Println(err.Error(), 1, statBonusString, raw.ItemId)
		fmt.Errorf("Error loading GemTable statBonusString")
	}
	raw.Effect, err = parseIntArrayField(effectString, 3)
	if err != nil {
		fmt.Println(err.Error(), 2, effectString, raw.ItemId)
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
		sie.Effect,
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
		fmt.Println(err.Error())
		return fmt.Errorf("error loading items for GemTables: %w", err)
	}

	RawGems = items

	return nil
}

type RawEnchant struct {
	EffectId           int
	Name               string
	SpellId            int
	ItemId             int
	ProfessionId       int
	Effects            []int
	EffectPoints       []int
	EffectArgs         []int
	IsWeaponEnchant    bool
	InvTypesMask       int
	SubClassMask       int
	ClassMask          int
	FDID               int
	Quality            int
	RequiredProfession int
}

func ScanEnchantsTable(rows *sql.Rows) (RawEnchant, error) {
	var raw RawEnchant
	var effectsString string
	var effectPointsString string
	var effectArgsString string
	err := rows.Scan(
		&raw.EffectId,
		&raw.Name,
		&raw.SpellId,
		&raw.ItemId,
		&raw.ProfessionId,
		&effectsString,
		&effectPointsString,
		&effectArgsString,
		&raw.IsWeaponEnchant,
		&raw.InvTypesMask,
		&raw.SubClassMask,
		&raw.ClassMask, &raw.FDID, &raw.Quality, &raw.RequiredProfession)
	if raw.InvTypesMask > 0 {
		//fmt.Println(raw.InvTypesMask)
	}
	raw.Effects, err = parseIntArrayField(effectsString, 3)
	if err != nil {
		fmt.Println(err.Error(), 3, effectsString)
		fmt.Errorf("Error loading ScanEnchantsTable effectsString")
	}
	raw.EffectPoints, err = parseIntArrayField(effectPointsString, 3)
	if err != nil {
		fmt.Println(err.Error(), 1, effectPointsString)
		fmt.Errorf("Error loading ScanEnchantsTable effectPointsString")
	}
	raw.EffectArgs, err = parseIntArrayField(effectArgsString, 3)
	if err != nil {
		fmt.Println(err.Error(), 2, effectArgsString)
		fmt.Errorf("Error loading ScanEnchantsTable effectArgsString")
	}
	return raw, err
}

func LoadRawEnchants(dbHelper *DBHelper) error {
	query := `SELECT DISTINCT
		sie.ID as effectId,
		sn.Name_lang as name,
		se.SpellID as spellId,
		COALESCE(ie.ParentItemID, 0) as ItemId,
		sie.Field_1_15_3_55112_014 as professionId,
		sie.Effect as Effect,
		sie.EffectPointsMax as EffectPoints,
		sie.EffectArg as EffectArgs,
		CASE 
			WHEN sei.EquippedItemClass = 4 THEN false
			ELSE true
		END AS isWeaponEnchant,
		sei.EquippedItemInvTypes as InvTypes,
		sei.EquippedItemSubclass,
		COALESCE(sla.ClassMask, 0),
		COALESCE(it.IconFileDataID, 0),
		COALESCE(isp.OverallQualityID, 1),
		COALESCE(sie.Field_1_15_3_55112_014, 0) as RequiredProfession
		FROM SpellEffect se 
		JOIN Spell s ON se.SpellID = s.ID
		JOIN SpellName sn ON se.SpellID = sn.ID
		JOIN SpellItemEnchantment sie ON se.EffectMiscValue_0 = sie.ID
		LEFT JOIN ItemEffect ie ON se.SpellID = ie.SpellID
		LEFT JOIN SpellEquippedItems sei ON se.SpellId = sei.SpellID
		LEFT JOIN SkillLineAbility sla ON se.SpellID = sla.Spell
		LEFT JOIN Item it ON ie.ParentItemId = it.ID
		LEFT JOIN ItemSparse isp ON ie.ParentItemId = isp.ID
		WHERE se.Effect = 53 GROUP BY sie.ID`
	items, err := LoadRows(dbHelper.db, query, ScanEnchantsTable)
	if err != nil {
		return fmt.Errorf("error loading items for GemTables: %w", err)
	}

	RawEnchants = items

	return nil
}

//RandPropPoints

type RandPropAllocationRow struct {
	Ilvl       int32
	Allocation RandomPropAllocation
}

func ScanRandPropAllocationRow(rows *sql.Rows) (RandPropAllocationRow, error) {
	var row RandPropAllocationRow
	var damageReplaceStat int
	err := rows.Scan(
		&row.Ilvl,
		&damageReplaceStat,
		&row.Allocation.Epic0,
		&row.Allocation.Epic1,
		&row.Allocation.Epic2,
		&row.Allocation.Epic3,
		&row.Allocation.Epic4,
		&row.Allocation.Superior0,
		&row.Allocation.Superior1,
		&row.Allocation.Superior2,
		&row.Allocation.Superior3,
		&row.Allocation.Superior4,
		&row.Allocation.Good0,
		&row.Allocation.Good1,
		&row.Allocation.Good2,
		&row.Allocation.Good3,
		&row.Allocation.Good4,
	)
	return row, err
}

func LoadRandomPropAllocations(dbHelper *DBHelper) (RandomPropAllocationsByIlvl, error) {
	query := `SELECT ID, DamageReplaceStat, Epic_0, Epic_1, Epic_2, Epic_3, Epic_4, Superior_0, Superior_1, Superior_2, Superior_3, Superior_4, Good_0, Good_1, Good_2 ,Good_3, Good_4 FROM RandPropPoints`
	rowsData, err := LoadRows[RandPropAllocationRow](dbHelper.db, query, ScanRandPropAllocationRow)
	if err != nil {
		return nil, fmt.Errorf("error loading random property allocations: %w", err)
	}

	processed := make(RandomPropAllocationsByIlvl)
	for _, r := range rowsData {
		processed[r.Ilvl] = RandomPropAllocationMap{
			proto.ItemQuality_ItemQualityEpic:     [5]int32{r.Allocation.Epic0, r.Allocation.Epic1, r.Allocation.Epic2, r.Allocation.Epic3, r.Allocation.Epic4},
			proto.ItemQuality_ItemQualityRare:     [5]int32{r.Allocation.Superior0, r.Allocation.Superior1, r.Allocation.Superior2, r.Allocation.Superior3, r.Allocation.Superior4},
			proto.ItemQuality_ItemQualityUncommon: [5]int32{r.Allocation.Good0, r.Allocation.Good1, r.Allocation.Good2, r.Allocation.Good3, r.Allocation.Good4},
		}
	}

	return processed, nil
}

type RawSpellEffect struct {
	ID                             int
	DifficultyID                   int
	EffectIndex                    int
	Effect                         int
	EffectAmplitude                float64
	EffectAttributes               int
	EffectAura                     int
	EffectAuraPeriod               int
	EffectBasePoints               int
	EffectBonusCoefficient         float64
	EffectChainAmplitude           float64
	EffectChainTargets             int
	EffectDieSides                 int
	EffectItemType                 int
	EffectMechanic                 int
	EffectPointsPerResource        float64
	EffectPosFacing                float64
	EffectRealPointsPerLevel       float64
	EffectTriggerSpell             int
	BonusCoefficientFromAP         float64
	PvpMultiplier                  float64
	Coefficient                    float64
	Variance                       float64
	ResourceCoefficient            float64
	GroupSizeBasePointsCoefficient float64
	// Grouped properties parsed from JSON strings:
	EffectMiscValues      []int // from EffectMiscValue, EffectMiscValue_0, EffectMiscValue_1
	EffectRadiusIndices   []int // from EffectRadiusIndex, EffectRadiusIndex_0, EffectRadiusIndex_1
	EffectSpellClassMasks []int // from EffectSpellClassMask, EffectSpellClassMask_0, EffectSpellClassMask_1, EffectSpellClassMask_2, EffectSpellClassMask_3
	ImplicitTargets       []int // from ImplicitTarget, ImplicitTarget_0, ImplicitTarget_1
	SpellID               int
}

func ScanSpellEffect(rows *sql.Rows) (RawSpellEffect, error) {
	var raw RawSpellEffect
	// Temporary strings to hold the concatenated JSON for grouped fields.
	var miscValuesStr, radiusIndicesStr, spellClassMasksStr, implicitTargetsStr string

	err := rows.Scan(
		&raw.ID,
		&raw.DifficultyID,
		&raw.EffectIndex,
		&raw.Effect,
		&raw.EffectAmplitude,
		&raw.EffectAttributes,
		&raw.EffectAura,
		&raw.EffectAuraPeriod,
		&raw.EffectBasePoints,
		&raw.EffectBonusCoefficient,
		&raw.EffectChainAmplitude,
		&raw.EffectChainTargets,
		&raw.EffectDieSides,
		&raw.EffectItemType,
		&raw.EffectMechanic,
		&raw.EffectPointsPerResource,
		&raw.EffectPosFacing,
		&raw.EffectRealPointsPerLevel,
		&raw.EffectTriggerSpell,
		&raw.BonusCoefficientFromAP,
		&raw.PvpMultiplier,
		&raw.Coefficient,
		&raw.Variance,
		&raw.ResourceCoefficient,
		&raw.GroupSizeBasePointsCoefficient,
		&miscValuesStr,
		&radiusIndicesStr,
		&spellClassMasksStr,
		&implicitTargetsStr,
		&raw.SpellID,
	)
	if err != nil {
		return raw, err
	}

	raw.EffectMiscValues, err = parseIntArrayField(miscValuesStr, 2)
	if err != nil {
		return raw, fmt.Errorf("error parsing EffectMiscValues: %w", err)
	}
	raw.EffectRadiusIndices, err = parseIntArrayField(radiusIndicesStr, 2)
	if err != nil {
		return raw, fmt.Errorf("error parsing EffectRadiusIndices: %w", err)
	}
	raw.EffectSpellClassMasks, err = parseIntArrayField(spellClassMasksStr, 4)
	if err != nil {
		return raw, fmt.Errorf("error parsing EffectSpellClassMasks: %w", err)
	}
	raw.ImplicitTargets, err = parseIntArrayField(implicitTargetsStr, 2)
	if err != nil {
		return raw, fmt.Errorf("error parsing ImplicitTargets: %w", err)
	}

	return raw, nil
}

func LoadRawSpellEffects(dbHelper *DBHelper) error {
	query := `
	SELECT
		ID,
		DifficultyID,
		EffectIndex,
		Effect,
		EffectAmplitude,
		EffectAttributes,
		EffectAura,
		EffectAuraPeriod,
		EffectBasePoints,
		EffectBonusCoefficient,
		EffectChainAmplitude,
		EffectChainTargets,
		EffectDieSides,
		EffectItemType,
		EffectMechanic,
		EffectPointsPerResource,
		EffectPos_facing,
		EffectRealPointsPerLevel,
		EffectTriggerSpell,
		BonusCoefficientFromAP,
		PvpMultiplier,
		Coefficient,
		Variance,
		ResourceCoefficient,
		GroupSizeBasePointsCoefficient,
		EffectMiscValue,
		EffectRadiusIndex,
		EffectSpellClassMask,
		ImplicitTarget,
		SpellID
	FROM SpellEffect
	`
	items, err := LoadRows(dbHelper.db, query, ScanSpellEffect)
	if err != nil {
		return fmt.Errorf("error loading SpellEffects: %w", err)
	}
	groupedBySpellID := make(map[int][]RawSpellEffect)
	for _, effect := range items {
		groupedBySpellID[effect.SpellID] = append(groupedBySpellID[effect.SpellID], effect)
	}

	RawSpellEffectBySpellIdAndIndex = make(map[int]map[int]RawSpellEffect)
	for spellID, effects := range groupedBySpellID {
		RawSpellEffectBySpellIdAndIndex[spellID] = CacheBy(effects, func(e RawSpellEffect) int {
			return e.EffectIndex
		})
	}
	return nil
}

// RawRandomSuffix represents the combined result of the ItemRandomSuffix row
// with its joined SpellItemEnchantment columns.
type RawRandomSuffix struct {
	ID            int
	Name          string
	AllocationPct []int // AllocationPct_0-4
	EffectArgs    []int // EffectArg_0-4
	Effects       []int // Effect_0-4
}

// ScanRawRandomSuffix scans one row from the query result into a RawRandomSuffix struct.
func ScanRawRandomSuffix(rows *sql.Rows) (RawRandomSuffix, error) {
	var raw RawRandomSuffix
	var (
		id     int
		name   string
		alloc0 int
		alloc1 int
		alloc2 int
		alloc3 int
		alloc4 int
		eArg0  int
		eArg1  int
		eArg2  int
		eArg3  int
		eArg4  int
		eff0   int
		eff1   int
		eff2   int
		eff3   int
		eff4   int
	)

	// The order here must match the SELECT list order.
	err := rows.Scan(
		&id,
		&name,
		&alloc0,
		&alloc1,
		&alloc2,
		&alloc3,
		&alloc4,
		&eArg0,
		&eArg1,
		&eArg2,
		&eArg3,
		&eArg4,
		&eff0,
		&eff1,
		&eff2,
		&eff3,
		&eff4,
	)
	if err != nil {
		return raw, err
	}

	raw.ID = id
	raw.Name = name
	raw.AllocationPct = []int{alloc0, alloc1, alloc2, alloc3, alloc4}
	raw.EffectArgs = []int{eArg0, eArg1, eArg2, eArg3, eArg4}
	raw.Effects = []int{eff0, eff1, eff2, eff3, eff4}

	return raw, nil
}

var RawRandomSuffixes []RawRandomSuffix
var RawRandomSuffixesById map[int]RawRandomSuffix

func LoadRawRandomSuffixes(dbHelper *DBHelper) error {
	query := `
	SELECT 
		COALESCE(-irs.ID, 0) AS ID,
		COALESCE(irs.Name_lang, '') AS Name_lang,
		COALESCE(irs.AllocationPct_0, 0) AS AllocationPct_0,
		COALESCE(irs.AllocationPct_1, 0) AS AllocationPct_1,
		COALESCE(irs.AllocationPct_2, 0) AS AllocationPct_2,
		COALESCE(irs.AllocationPct_3, 0) AS AllocationPct_3,
		COALESCE(irs.AllocationPct_4, 0) AS AllocationPct_4,
		COALESCE(sie0.EffectArg_0, 0) AS EffectArg_0,
		COALESCE(sie1.EffectArg_0, 0) AS EffectArg_1,
		COALESCE(sie2.EffectArg_0, 0) AS EffectArg_2,
		COALESCE(sie3.EffectArg_0, 0) AS EffectArg_3,
		COALESCE(sie4.EffectArg_0, 0) AS EffectArg_4,
		COALESCE(sie0.Effect_0, 0) AS Effect_0,
		COALESCE(sie1.Effect_0, 0) AS Effect_1,
		COALESCE(sie2.Effect_0, 0) AS Effect_2,
		COALESCE(sie3.Effect_0, 0) AS Effect_3,
		COALESCE(sie4.Effect_0, 0) AS Effect_4
	FROM ItemRandomSuffix irs
	LEFT JOIN SpellItemEnchantment sie0 ON irs.Enchantment_0 = sie0.ID
	LEFT JOIN SpellItemEnchantment sie1 ON irs.Enchantment_1 = sie1.ID
	LEFT JOIN SpellItemEnchantment sie2 ON irs.Enchantment_2 = sie2.ID
	LEFT JOIN SpellItemEnchantment sie3 ON irs.Enchantment_3 = sie3.ID
	LEFT JOIN SpellItemEnchantment sie4 ON irs.Enchantment_4 = sie4.ID;
`
	// Use your generic LoadRows function to scan each row into a RawRandomSuffix.
	items, err := LoadRows(dbHelper.db, query, ScanRawRandomSuffix)
	if err != nil {
		return fmt.Errorf("error loading RawRandomSuffixes: %w", err)
	}

	RawRandomSuffixes = items
	RawRandomSuffixesById = CacheBy(items, func(suffix RawRandomSuffix) int {
		return suffix.ID
	})
	return nil
}
