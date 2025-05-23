package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/tools/database/dbc"
	"github.com/wowsims/mop/tools/tooltip"
)

// Loading tables
// Below is the definition and loading of tables
//

// Raw Item Data

func ScanRawItemData(rows *sql.Rows) (dbc.Item, error) {
	var raw dbc.Item
	var randomSuffixOptions sql.NullString
	var statPercentageOfSocket string
	var bonusAmountCalculated string
	var bonusStatString string
	var statValue string
	var socketTypes string
	var statPercentEditor string
	err := rows.Scan(&raw.Id, &raw.Name, &raw.InventoryType, &raw.ItemDelay, &raw.OverallQuality, &raw.DmgVariance,
		&raw.ItemLevel,
		&statValue, &bonusStatString,
		&statPercentEditor, &socketTypes, &raw.SocketEnchantmentId, &raw.Flags0, &raw.FDID, &raw.ItemSetName, &raw.ItemSetId, &raw.Flags1, &raw.ClassMask, &raw.RaceMask, &raw.QualityModifier, &randomSuffixOptions, &statPercentageOfSocket, &bonusAmountCalculated, &raw.ItemClass, &raw.ItemSubClass)
	if err != nil {
		panic(err)
	}
	var parseErr error
	raw.RandomSuffixOptions, parseErr = ParseRandomSuffixOptions(randomSuffixOptions)
	if parseErr != nil {
		return raw, fmt.Errorf("failed to parse RandomSuffixOptions: %w", parseErr)
	}

	raw.StatPercentageOfSocket, err = parseFloatArrayField(statPercentageOfSocket, 10)
	if err != nil {
		return raw, fmt.Errorf("failed to parse StatPercentageOfSocket: %w %s", err, statPercentageOfSocket)
	}
	raw.StatAlloc, err = parseFloatArrayField(statValue, 10)
	if err != nil {
		return raw, fmt.Errorf("failed to parse StatAlloc: %w", err)
	}
	raw.BonusAmountCalculated, err = parseFloatArrayField(bonusAmountCalculated, 10)
	if err != nil {
		return raw, fmt.Errorf("failed to parse BonusAmountCalculated: %w", err)
	}
	raw.BonusStat, err = parseIntArrayField(bonusStatString, 10)
	if err != nil {
		return raw, fmt.Errorf("failed to parse BonusStat: %w", err)
	}
	raw.Sockets, err = parseIntArrayField(socketTypes, 3)
	if err != nil {
		return raw, fmt.Errorf("failed to parse Sockets: %w", err)
	}
	raw.SocketModifier, err = parseFloatArrayField(statPercentEditor, 10)
	if err != nil {
		return raw, fmt.Errorf("failed to parse SocketModifier: %w", err)
	}
	return raw, err
}

func LoadAndWriteRawItems(dbHelper *DBHelper, filter string, inputsDir string) ([]dbc.Item, error) {
	baseQuery := `
		SELECT
			i.ID,
			s.Display_lang AS Name,
			i.InventoryType,
			s.ItemDelay,
			s.OverallQualityID,
			s.DmgVariance,
			s.ItemLevel,
			s.Field_1_15_3_55112_014 as StatValue,
			s.StatModifier_bonusStat as bonusStat,
			s.StatPercentEditor as StatPercentEditor,
			s.SocketType as SocketTypes,
			s.Field_1_15_7_59706_036 as SocketEnchantmentId,
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
			) AS RandomSuffixOptions,
			 s.StatPercentageOfSocket,
			 s.StatModifier_bonusAmount,
			 i.ClassID,
			 i.SubClassID
		FROM Item i
		JOIN ItemSparse s ON i.ID = s.ID
		JOIN ItemClass ic ON i.ClassID = ic.ClassID
		JOIN ItemSubClass isc ON i.ClassID = isc.ClassID AND i.SubClassID = isc.SubClassID
		JOIN RandPropPoints rpp ON s.ItemLevel = rpp.ID
		LEFT JOIN ItemArmorShield ias ON s.ItemLevel = ias.ItemLevel
		LEFT JOIN ItemSet itemset ON s.ItemSet = itemset.ID
		LEFT JOIN ItemArmorQuality iaq ON s.ItemLevel = iaq.ID
		JOIN ItemArmorTotal at ON s.ItemLevel = at.ItemLevel
		`

	if strings.TrimSpace(filter) != "" {
		baseQuery += " WHERE " + filter
	}

	items, err := LoadRows(dbHelper.db, baseQuery, ScanRawItemData)
	fmt.Println("Loaded Items:", len(items))
	if err != nil {
		fmt.Println("Error loading items:", err.Error())
		return nil, err
	}
	json, _ := json.Marshal(items)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/items.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}

	return items, nil
}

//ItemStatEffects
// Used for straight up item stat effects from SpellItemEnchantment (socket bonuses for now, single stat)
//

func ScanItemStatEffects(rows *sql.Rows) (dbc.ItemStatEffect, error) {
	var raw dbc.ItemStatEffect
	var ePointsMin, epointsMax, eArgs string
	err := rows.Scan(&raw.ID, &ePointsMin, &epointsMax, &eArgs)
	if err != nil {
		panic("Error scanning item stat effects")
	}
	raw.EffectPointsMin, err = parseIntArrayField(ePointsMin, 3)
	if err != nil {
		return raw, fmt.Errorf("failed to parse EffectPointsMin: %w", err)
	}
	raw.EffectPointsMax, err = parseIntArrayField(epointsMax, 3)
	if err != nil {
		return raw, fmt.Errorf("failed to parse EffectPointsMax: %w", err)
	}
	raw.EffectArg, err = parseIntArrayField(eArgs, 3)
	if err != nil {
		return raw, fmt.Errorf("failed to parse EffectArg: %w", err)
	}
	return raw, err
}

func LoadAndWriteItemStatEffects(dbHelper *DBHelper, inputsDir string) ([]dbc.ItemStatEffect, error) {
	query := `SELECT ID, EffectPointsMin, EffectPointsMax, EffectArg FROM SpellItemEnchantment WHERE Effect_0 = 5`
	items, err := LoadRows(dbHelper.db, query, ScanItemStatEffects)
	if err != nil {
		return nil, fmt.Errorf("error in query load items")
	}
	json, _ := json.Marshal(items)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/item_stat_effects.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return items, nil
}

func ScanItemDamageTable(rows *sql.Rows) (dbc.ItemDamageTable, error) {
	var raw dbc.ItemDamageTable
	var qualityString string
	err := rows.Scan(&raw.ItemLevel, &qualityString)
	if err != nil {
		return raw, fmt.Errorf("scanning item damage table: %w", err)
	}

	raw.Quality, err = parseFloatArrayField(qualityString, 7)
	if err != nil {
		return raw, fmt.Errorf("parsing quality string '%s': %w", qualityString, err)
	}

	return raw, nil
}

var ItemDamageByTableAndItemLevel = make(map[string]map[int]dbc.ItemDamageTable)
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

func LoadAndWriteItemDamageTables(dbHelper *DBHelper, inputsDir string) (map[string]map[int]dbc.ItemDamageTable, error) {
	for _, tableName := range itemDamageTableNames {
		query := fmt.Sprintf("SELECT ItemLevel, Quality FROM %s", tableName)
		items, err := LoadRows(dbHelper.db, query, ScanItemDamageTable)
		if err != nil {
			return nil, fmt.Errorf("error loading items for table %s: %w", tableName, err)
		}

		// Cache the slice of ItemDamageTable into a map keyed by ItemLevel.
		ItemDamageByTableAndItemLevel[tableName] = CacheBy(items, func(table dbc.ItemDamageTable) int {
			return table.ItemLevel
		})
	}
	json, _ := json.Marshal(ItemDamageByTableAndItemLevel)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/item_damage_tables.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return ItemDamageByTableAndItemLevel, nil
}

func LoadAndWriteItemArmorQuality(dbHelper *DBHelper, inputsDir string) (map[int]dbc.ItemArmorQuality, error) {
	query := "SELECT ID, Qualitymod FROM ItemArmorQuality"
	result, err := LoadRows(dbHelper.db, query, ScanItemArmorQualityTable)

	cache := CacheBy(result, func(table dbc.ItemArmorQuality) int {
		return table.ItemLevel
	})
	json, _ := json.Marshal(cache)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/item_armor_quality.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return cache, err
}

func ScanItemArmorQualityTable(rows *sql.Rows) (dbc.ItemArmorQuality, error) {
	var raw dbc.ItemArmorQuality
	var qualityString string
	err := rows.Scan(&raw.ItemLevel, &qualityString)
	if err != nil {
		return raw, fmt.Errorf("scanning item armor quality: %w", err)
	}

	raw.Quality, err = parseFloatArrayField(qualityString, 7)
	if err != nil {
		return raw, fmt.Errorf("parsing quality string '%s': %w", qualityString, err)
	}

	return raw, nil
}

func LoadAndWriteItemArmorShield(dbHelper *DBHelper, inputsDir string) (map[int]dbc.ItemArmorShield, error) {
	query := "SELECT ItemLevel, Quality FROM ItemArmorShield"
	result, err := LoadRows(dbHelper.db, query, ScanItemArmorShieldTable)

	cache := CacheBy(result, func(table dbc.ItemArmorShield) int {
		return table.ItemLevel
	})
	json, _ := json.Marshal(cache)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/item_armor_shield.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return cache, err
}

func ScanItemArmorShieldTable(rows *sql.Rows) (dbc.ItemArmorShield, error) {
	var raw dbc.ItemArmorShield
	var qualityString string
	err := rows.Scan(&raw.ItemLevel, &qualityString)
	if err != nil {
		return raw, fmt.Errorf("scanning item armor shield: %w", err)
	}

	raw.Quality, err = parseFloatArrayField(qualityString, 7)
	if err != nil {
		return raw, fmt.Errorf("parsing quality string '%s': %w", qualityString, err)
	}

	return raw, nil
}
func LoadAndWriteItemArmorTotal(dbHelper *DBHelper, inputsDir string) (map[int]dbc.ItemArmorTotal, error) {
	query := "SELECT ItemLevel, Cloth, Leather, Mail, Plate FROM ItemArmorTotal"
	result, err := LoadRows(dbHelper.db, query, ScanItemArmorTotalTable)
	cached := CacheBy(result, func(table dbc.ItemArmorTotal) int {
		return table.ItemLevel
	})
	json, _ := json.Marshal(cached)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/item_armor_total.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}

	return cached, err
}

func ScanItemArmorTotalTable(rows *sql.Rows) (dbc.ItemArmorTotal, error) {
	var raw dbc.ItemArmorTotal
	var qualityString string
	err := rows.Scan(&raw.ItemLevel, &raw.Cloth, &raw.Leather, &raw.Mail, &raw.Plate)
	if err != nil {
		fmt.Println(err.Error(), 3, qualityString)
		return raw, fmt.Errorf("error loading ScanItemArmorTotalTable")
	}
	return raw, err
}

func LoadAndWriteArmorLocation(dbHelper *DBHelper, inputsDir string) (map[int]dbc.ArmorLocation, error) {
	query := "SELECT ID, Clothmodifier, Leathermodifier, Chainmodifier, Platemodifier, Modifier FROM ArmorLocation"
	result, err := LoadRows(dbHelper.db, query, ScanArmorLocation)
	cache := CacheBy(result, func(table dbc.ArmorLocation) int {
		return table.Id
	})
	json, _ := json.Marshal(cache)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/armor_location.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return cache, err
}

func ScanArmorLocation(rows *sql.Rows) (dbc.ArmorLocation, error) {
	var raw dbc.ArmorLocation
	raw.Modifier = [5]float64{}
	err := rows.Scan(&raw.Id, &raw.Modifier[0], &raw.Modifier[1], &raw.Modifier[2], &raw.Modifier[3], &raw.Modifier[4])
	return raw, err
}

// ItemDamage tables
func ScanGemTable(rows *sql.Rows) (dbc.Gem, error) {
	var raw dbc.Gem
	var statListString string
	var statBonusString string
	var effectString string
	err := rows.Scan(&raw.ItemId, &raw.Name, &raw.FDID, &raw.GemType, &statListString, &statBonusString, &raw.MinItemLevel, &raw.Quality, &effectString, &raw.IsJc, &raw.Flags0)
	if err != nil {
		return raw, fmt.Errorf("scanning gem data: %w", err)
	}

	raw.EffectPoints, err = parseIntArrayField(statListString, 3)
	if err != nil {
		return raw, fmt.Errorf("parsing effect points for gem %d (%s): %w", raw.ItemId, statListString, err)
	}

	raw.EffectArgs, err = parseIntArrayField(statBonusString, 3)
	if err != nil {
		return raw, fmt.Errorf("parsing effect args for gem %d (%s): %w", raw.ItemId, statBonusString, err)
	}

	raw.Effects, err = parseIntArrayField(effectString, 3)
	if err != nil {
		return raw, fmt.Errorf("parsing effects for gem %d (%s): %w", raw.ItemId, effectString, err)
	}

	return raw, nil
}

func LoadAndWriteRawGems(dbHelper *DBHelper, inputsDir string) ([]dbc.Gem, error) {
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
		JOIN GemProperties gp ON s.Field_1_15_7_59706_035  = gp.ID
		JOIN SpellItemEnchantment sie ON gp.Enchant_ID = sie.ID
		WHERE i.ClassID = 3`
	items, err := LoadRows(dbHelper.db, query, ScanGemTable)
	if err != nil {
		return nil, fmt.Errorf("error loading items for GemTables: %w", err)
	}
	json, _ := json.Marshal(items)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/gems.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return items, nil
}

func ScanEnchantsTable(rows *sql.Rows) (dbc.Enchant, error) {
	var raw dbc.Enchant
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
		&raw.InventoryType,
		&raw.SubClassMask,
		&raw.ClassMask, &raw.FDID, &raw.Quality, &raw.RequiredProfession, &raw.EffectName)
	if err != nil {
		return raw, fmt.Errorf("scanning enchant data for effect ID %d: %w", raw.EffectId, err)
	}

	raw.Effects, err = parseIntArrayField(effectsString, 3)
	if err != nil {
		return raw, fmt.Errorf("parsing effects for enchant %d (%s): %w", raw.EffectId, effectsString, err)
	}

	raw.EffectPoints, err = parseIntArrayField(effectPointsString, 3)
	if err != nil {
		return raw, fmt.Errorf("parsing effect points for enchant %d (%s): %w", raw.EffectId, effectPointsString, err)
	}

	raw.EffectArgs, err = parseIntArrayField(effectArgsString, 3)
	if err != nil {
		return raw, fmt.Errorf("parsing effect args for enchant %d (%s): %w", raw.EffectId, effectArgsString, err)
	}

	return raw, nil
}

func LoadAndWriteRawEnchants(dbHelper *DBHelper, inputsDir string) ([]dbc.Enchant, error) {
	query := `SELECT DISTINCT
		sie.ID as effectId,
		CASE
		    WHEN sn.Name_lang LIKE '%+%' THEN COALESCE(isp.Display_lang, sn.Name_lang)
		    ELSE sn.Name_lang
		END AS name,
		se.SpellID as spellId,
		COALESCE(ie.ParentItemID, 0) as ItemId,
		sie.Field_1_15_3_55112_014 as professionId,
		sie.Effect as Effect,
		sie.EffectPointsMin as EffectPoints,
		sie.EffectArg as EffectArgs,
		CASE
			WHEN sei.EquippedItemClass = 4 THEN false
			ELSE true
		END AS isWeaponEnchant,
		COALESCE(sei.EquippedItemInvTypes, 0) as InvTypes,
		COALESCE(sei.EquippedItemSubclass, 0),
		COALESCE(sla.ClassMask, 0),
		COALESCE(it.IconFileDataID, 0),
		COALESCE(isp.OverallQualityID, 1),
		COALESCE(sie.Field_1_15_3_55112_014, 0) as RequiredProfession,
		COALESCE(sie.Name_lang, "")
		FROM SpellEffect se
		JOIN Spell s ON se.SpellID = s.ID
		JOIN SpellName sn ON se.SpellID = sn.ID
		JOIN SpellItemEnchantment sie ON se.EffectMiscValue_0 = sie.ID
		LEFT JOIN ItemEffect ie ON se.SpellID = ie.SpellID
		LEFT JOIN SpellEquippedItems sei ON se.SpellId = sei.SpellID
		LEFT JOIN SkillLineAbility sla ON se.SpellID = sla.Spell
		LEFT JOIN Item it ON ie.ParentItemId = it.ID
		LEFT JOIN ItemSparse isp ON ie.ParentItemId = isp.ID
		WHERE se.Effect = 53 GROUP BY sn.Name_lang, sie.ID`
	items, err := LoadRows(dbHelper.db, query, ScanEnchantsTable)
	if err != nil {
		return nil, fmt.Errorf("error loading items for GemTables: %w", err)
	}
	json, _ := json.Marshal(items)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/enchants.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return items, nil
}

//RandPropPoints

type RandPropAllocationRow struct {
	Ilvl       int32
	Allocation dbc.RandomPropAllocation
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

func LoadAndWriteRandomPropAllocations(dbHelper *DBHelper, inputsDir string) (map[int32]RandPropAllocationRow, error) {
	query := `SELECT ID, DamageReplaceStat, Epic_0, Epic_1, Epic_2, Epic_3, Epic_4, Superior_0, Superior_1, Superior_2, Superior_3, Superior_4, Good_0, Good_1, Good_2 ,Good_3, Good_4 FROM RandPropPoints`
	rowsData, err := LoadRows(dbHelper.db, query, ScanRandPropAllocationRow)
	if err != nil {
		return nil, fmt.Errorf("error loading random property allocations: %w", err)
	}

	processed := make(map[int32]RandPropAllocationRow)
	for _, r := range rowsData {
		processed[r.Ilvl] = r
	}

	randProps := make(dbc.RandomPropAllocationsByIlvl)
	for _, r := range processed {
		randProps[int(r.Ilvl)] = dbc.RandomPropAllocationMap{
			proto.ItemQuality_ItemQualityEpic:     [5]int32{r.Allocation.Epic0, r.Allocation.Epic1, r.Allocation.Epic2, r.Allocation.Epic3, r.Allocation.Epic4},
			proto.ItemQuality_ItemQualityRare:     [5]int32{r.Allocation.Superior0, r.Allocation.Superior1, r.Allocation.Superior2, r.Allocation.Superior3, r.Allocation.Superior4},
			proto.ItemQuality_ItemQualityUncommon: [5]int32{r.Allocation.Good0, r.Allocation.Good1, r.Allocation.Good2, r.Allocation.Good3, r.Allocation.Good4},
		}
	}
	json, _ := json.Marshal(randProps)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/rand_prop_points.json", inputsDir), json); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
	return processed, nil
}

func ScanSpellEffect(rows *sql.Rows) (dbc.SpellEffect, error) {
	var raw dbc.SpellEffect
	raw.EffectMinRange = []float64{0, 0}
	raw.EffectMaxRange = []float64{0, 0}

	// Temporary strings to hold the concatenated JSON for grouped fields.
	var miscValuesStr, spellClassMasksStr, implicitTargetsStr string
	err := rows.Scan(
		&raw.ID,
		&raw.DifficultyID,
		&raw.EffectIndex,
		&raw.EffectType,
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
		&spellClassMasksStr,
		&implicitTargetsStr,
		&raw.SpellID,
		&raw.ScalingType,
		&raw.EffectMinRange[0],
		&raw.EffectMaxRange[0],
		&raw.EffectMinRange[1],
		&raw.EffectMaxRange[1],
	)
	if err != nil {
		return raw, err
	}

	raw.EffectMiscValues, err = parseIntArrayField(miscValuesStr, 2)
	if err != nil {
		return raw, fmt.Errorf("error parsing EffectMiscValues: %w", err)
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

func LoadAndWriteRawSpellEffects(dbHelper *DBHelper, inputsDir string) (map[int]map[int]dbc.SpellEffect, error) {
	query := `
	SELECT
		se.ID,
		se.DifficultyID,
		se.EffectIndex,
		se.Effect,
		se.EffectAmplitude,
		se.EffectAttributes,
		se.EffectAura,
		se.EffectAuraPeriod,
		se.EffectBasePoints,
		se.EffectBonusCoefficient,
		se.EffectChainAmplitude,
		se.EffectChainTargets,
		se.EffectDieSides,
		se.EffectItemType,
		se.EffectMechanic,
		se.EffectPointsPerResource,
		se.EffectPos_facing,
		se.EffectRealPointsPerLevel,
		se.EffectTriggerSpell,
		se.BonusCoefficientFromAP,
		se.PvpMultiplier,
		se.Coefficient,
		se.Variance,
		se.ResourceCoefficient,
		se.GroupSizeBasePointsCoefficient,
		se.EffectMiscValue,
		se.EffectSpellClassMask,
		se.ImplicitTarget,
		se.SpellID,
		COALESCE(ss.Class, 0),
		COALESCE(sr1.RadiusMin, 0),
		COALESCE(sr1.RadiusMax, 0),
		COALESCE(sr2.RadiusMin, 0),
		COALESCE(sr2.RadiusMax, 0)
	FROM SpellEffect se
	LEFT JOIN SpellScaling ss ON se.SpellID = ss.SpellID
	LEFT JOIN SpellRadius sr1 ON sr1.ID = se.EffectRadiusIndex_0
	LEFT JOIN SpellRadius sr2 ON sr2.ID = se.EffectRadiusIndex_1
	`
	items, err := LoadRows(dbHelper.db, query, ScanSpellEffect)
	if err != nil {
		return nil, fmt.Errorf("error loading SpellEffects: %w", err)
	}
	groupedBySpellID := make(map[int][]dbc.SpellEffect)
	for _, effect := range items {
		groupedBySpellID[effect.SpellID] = append(groupedBySpellID[effect.SpellID], effect)
	}

	RawSpellEffectBySpellIdAndIndex := make(map[int]map[int]dbc.SpellEffect)
	for spellID, effects := range groupedBySpellID {
		RawSpellEffectBySpellIdAndIndex[spellID] = CacheBy(effects, func(e dbc.SpellEffect) int {
			return e.EffectIndex
		})
	}
	json, _ := json.Marshal(RawSpellEffectBySpellIdAndIndex)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/spell_effects.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return RawSpellEffectBySpellIdAndIndex, nil
}

// RawRandomSuffix represents the combined result of the ItemRandomSuffix row
// with its joined SpellItemEnchantment columns.

// ScanRawRandomSuffix scans one row from the query result into a RawRandomSuffix struct.
func ScanRawRandomSuffix(rows *sql.Rows) (dbc.RandomSuffix, error) {
	var raw dbc.RandomSuffix
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

var RawRandomSuffixes []dbc.RandomSuffix
var RawRandomSuffixesById map[int]dbc.RandomSuffix

func LoadAndWriteRawRandomSuffixes(dbHelper *DBHelper, inputsDir string) ([]dbc.RandomSuffix, error) {
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
		return nil, fmt.Errorf("error loading RawRandomSuffixes: %w", err)
	}

	RawRandomSuffixes = items
	RawRandomSuffixesById = CacheBy(items, func(suffix dbc.RandomSuffix) int {
		return suffix.ID
	})
	json, _ := json.Marshal(items)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/random_suffix.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return items, nil
}

func ScanConsumable(rows *sql.Rows) (dbc.Consumable, error) {
	var consumable dbc.Consumable
	var itemEffectsStr sql.NullString

	err := rows.Scan(
		&consumable.Id,
		&consumable.Name,
		&consumable.ItemLevel,
		&consumable.RequiredLevel,
		&consumable.ClassId,
		&consumable.SubClassId,
		&consumable.IconFileDataID,
		&consumable.SpellCategoryID,
		&consumable.SpellCategoryFlags,
		&itemEffectsStr,
		&consumable.ElixirType,
		&consumable.Duration,
	)
	if err != nil {
		return consumable, fmt.Errorf("scanning consumable data: %w", err)
	}

	if !itemEffectsStr.Valid || itemEffectsStr.String == "" || itemEffectsStr.String == "null" {
		consumable.ItemEffects = []int{}
	} else {
		parts := strings.Split(itemEffectsStr.String, ",")
		effects := make([]int, 0, len(parts))
		for _, part := range parts {
			token := strings.TrimSpace(part)
			num, err := strconv.Atoi(token)
			if err != nil {
				fmt.Printf("Warning: parsing item effects for consumable %d (%s): %v\n", consumable.Id, token, err)
				effects = []int{}
				break
			}
			effects = append(effects, num)
		}
		consumable.ItemEffects = effects
	}

	return consumable, nil
}

func LoadAndWriteConsumables(dbHelper *DBHelper, inputsDir string) ([]dbc.Consumable, error) {
	query := `
		SELECT
				i.ID,
				s.Display_lang AS Name,
				s.ItemLevel,
				s.RequiredLevel,
				i.ClassID,
				i.SubClassID,
				i.IconFileDataID,
				COALESCE(ie.SpellCategoryID, 0) AS SpellCategoryID,
				COALESCE(sc.Flags, 0) AS SpellCategoryFlags,
				(
					SELECT group_concat(ie2.ID, ',')
					FROM ItemEffect ie2
					WHERE ie2.ParentItemID = i.ID
				) AS ItemEffects,
				CASE
					WHEN sp.Description_lang LIKE '%Guardian Elixir%' THEN 1
					WHEN sp.Description_lang LIKE '%Battle Elixir%' THEN 2
					ELSE 0
				END AS ElixirType,
				COALESCE(sd.Duration, 0) as Duration
			FROM Item i
			JOIN ItemSparse s ON i.ID = s.ID
			LEFT JOIN ItemEffect ie ON i.ID = ie.ParentItemID
			LEFT JOIN SpellCategory sc ON ie.SpellCategoryID = sc.ID
			LEFT JOIN Spell sp ON ie.SpellID = sp.ID
			LEFT JOIN SpellMisc sm ON ie.SpellId = sm.SpellID
			LEFT JOIN SpellDuration sd ON sm.DurationIndex = sd.ID
			WHERE ((i.ClassID = 0 AND i.SubclassID IS NOT 0 AND i.SubclassID IS NOT 8 AND i.SubclassID IS NOT 6) OR (i.ClassID = 7 AND i.SubclassID = 2)) AND ItemEffects is not null AND (s.RequiredLevel >= 70 OR i.ID = 22788 OR i.ID = 13442)
			AND s.Display_lang != ''
			AND s.Display_lang NOT LIKE '%Test%'
			AND s.Display_lang NOT LIKE 'QA%'
			GROUP BY i.ID
	`

	consumables, err := LoadRows(dbHelper.db, query, ScanConsumable)
	if err != nil {
		return nil, fmt.Errorf("error loading consumables: %w", err)
	}

	fmt.Println("Loaded Consumables:", len(consumables))
	json, _ := json.Marshal(consumables)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/consumables.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return consumables, nil
}

func ScanItemEffect(rows *sql.Rows) (dbc.ItemEffect, error) {
	var effect dbc.ItemEffect

	err := rows.Scan(
		&effect.ID,
		&effect.LegacySlotIndex,
		&effect.TriggerType,
		&effect.Charges,
		&effect.CoolDownMSec,
		&effect.CategoryCoolDownMSec,
		&effect.SpellCategoryID,
		&effect.SpellID,
		&effect.ChrSpecializationID,
		&effect.ParentItemID,
	)
	if err != nil {
		return effect, fmt.Errorf("scanning item effect data: %w", err)
	}

	return effect, nil
}

func LoadAndWriteItemEffects(dbHelper *DBHelper, inputsDir string) ([]dbc.ItemEffect, error) {
	query := `
	SELECT
		ID,
		LegacySlotIndex,
		TriggerType,
		Charges,
		CoolDownMSec,
		CategoryCoolDownMSec,
		SpellCategoryID,
		SpellID,
		ChrSpecializationID,
		ParentItemID
	FROM ItemEffect
	`

	effects, err := LoadRows(dbHelper.db, query, ScanItemEffect)
	if err != nil {
		return nil, fmt.Errorf("error loading item effects: %w", err)
	}

	fmt.Println("Loaded ItemEffects:", len(effects))
	json, _ := json.Marshal(effects)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/item_effects.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return effects, nil
}

type RawGlyph struct {
	ItemId      int32
	Name        string
	SpellId     int32
	Description string
	GlyphType   int32
	ClassMask   int32
	FDID        int32
}

func ScanGlyphs(rows *sql.Rows) (RawGlyph, error) {
	var glyph RawGlyph

	err := rows.Scan(
		&glyph.ItemId,
		&glyph.SpellId,
		&glyph.Name,
		&glyph.Description,
		&glyph.GlyphType,
		&glyph.ClassMask,
		&glyph.FDID,
	)
	if err != nil {
		return glyph, fmt.Errorf("scanning glyph data: %w", err)
	}

	return glyph, nil
}

func LoadGlyphs(dbHelper *DBHelper) ([]RawGlyph, error) {
	query := `
SELECT DISTINCT i.ID, gp.SpellID, is2.Display_lang, glyphSpell.Description_lang, gp.Field_3_4_0_43659_001, i.SubclassID, sm.SpellIconFileDataID
FROM Item i
LEFT JOIN ItemSparse is2 ON i.ID = is2.ID
LEFT JOIN ItemEffect ie ON ie.ParentItemID  = i.ID
JOIN SpellEffect se ON se.SpellID = ie.SpellID AND se.Effect=74
LEFT JOIN GlyphProperties gp ON gp.ID = se.EffectMiscValue_0
LEFT JOIN Spell glyphSpell ON glyphSpell.ID = gp.SpellID
LEFT JOIN SpellEffect gse ON gse.SpellID = glyphSpell.ID
LEFT JOIN SpellMisc sm ON sm.SpellID = glyphSpell.ID
WHERE i.ClassID = 16
GROUP BY i.ID

	`

	effects, err := LoadRows(dbHelper.db, query, ScanGlyphs)
	if err != nil {
		return nil, fmt.Errorf("error loading glyphs : %w", err)
	}

	fmt.Println("Loaded glyphs:", len(effects))
	return effects, nil
}

type RawTalent struct {
	TierID      int
	TalentName  string
	ColumnIndex int
	ClassMask   int
	SpellID     int
}

func ScanTalent(rows *sql.Rows) (RawTalent, error) {
	var talent RawTalent

	err := rows.Scan(
		&talent.TierID,
		&talent.TalentName,
		&talent.ColumnIndex,
		&talent.ClassMask,
		&talent.SpellID,
	)
	if err != nil {
		return talent, fmt.Errorf("scanning talent data: %w", err)
	}

	return talent, nil
}

func LoadTalents(dbHelper *DBHelper) ([]RawTalent, error) {
	query := `
SELECT
  t.TierID,
  sn.Name_lang,
  t.ColumnIndex,
  t.ClassID,
  t.SpellID
FROM Talent t
JOIN SpellName sn ON sn.ID = t.SpellID
WHERE sn.Name_lang IS NOT "Dummy 5.0 Talent"
ORDER BY t.ClassID;
`

	talents, err := LoadRows(dbHelper.db, query, ScanTalent)
	if err != nil {
		return nil, fmt.Errorf("error loading talents: %w", err)
	}

	fmt.Println("Loaded talents:", len(talents))
	return talents, nil
}

type SpellIcon struct {
	SpellID int
	FDID    int
	HasBuff bool
	Name    string
}

func ScanSpellIcon(rows *sql.Rows) (SpellIcon, error) {
	var icon SpellIcon

	err := rows.Scan(
		&icon.SpellID,
		&icon.FDID,
		&icon.HasBuff,
		&icon.Name,
	)
	if err != nil {
		return icon, fmt.Errorf("scanning talent data: %w", err)
	}

	return icon, nil
}

func LoadSpellIcons(dbHelper *DBHelper) (map[int]SpellIcon, error) {
	query := `
		SELECT
  sm.SpellID,
  sm.SpellIconFileDataID,
  (
    (sm.Attributes_4 & 0x00001000) <> 0
    OR EXISTS (
      SELECT 1
      FROM SpellEffect se
      WHERE se.SpellID = sm.SpellID
        AND se.Effect = 6
    ) OR (ss.AuraDescription_lang != '' and ss.AuraDescription_lang is not null)
  ) AS HasBuff,
  sn.Name_lang
FROM SpellMisc sm
LEFT JOIN Spell ss ON ss.ID = sm.SpellID
LEFT JOIN SpellName sn ON sn.ID = sm.SpellID;
`

	talents, err := LoadRows(dbHelper.db, query, ScanSpellIcon)
	if err != nil {
		return nil, fmt.Errorf("error loading spellicons: %w", err)
	}
	iconsByID := make(map[int]SpellIcon, len(talents))
	for _, icon := range talents {
		iconsByID[icon.SpellID] = icon
	}
	fmt.Println("Loaded spellicons:", len(talents))
	return iconsByID, nil
}

var iconsMap, _ = LoadArtTexturePaths("./tools/DB2ToSqlite/listfile.csv")

func ScanSpells(rows *sql.Rows) (dbc.Spell, error) {
	var spell dbc.Spell

	var stringAttr string
	var stringClassMask string
	var stringProcType string              //2
	var stringAuraIFlags string            //2
	var stringChannelInterruptFlags string // 2
	var stringShapeShift string            //2
	var iconId int                         //

	err := rows.Scan(
		&spell.NameLang,
		&spell.ID,
		&spell.SchoolMask,
		&spell.Speed,
		&spell.LaunchDelay,
		&spell.MinDuration,
		&spell.MaxScalingLevel,
		&spell.MinScalingLevel,
		&spell.ScalesFromItemLevel,
		&spell.SpellLevel,
		&spell.BaseLevel,
		&spell.MaxLevel,
		&spell.MaxPassiveAuraLevel,
		&spell.Cooldown,
		&spell.GCD,
		&spell.MinRange,
		&spell.MaxRange,
		&stringAttr,
		&spell.CategoryFlags,
		&spell.MaxCharges,
		&spell.ChargeRecoveryTime,
		&spell.CategoryTypeMask,
		&spell.Category,
		&spell.Duration,
		&spell.ProcChance,
		&spell.ProcCharges,
		&stringProcType,
		&spell.ProcCategoryRecovery,
		&spell.SpellProcsPerMinute,
		&spell.EquippedItemClass,
		&spell.EquippedItemInvTypes,
		&spell.EquippedItemSubclass,
		&spell.CastTimeMin,
		&stringClassMask,
		&spell.SpellClassSet,
		&stringAuraIFlags,
		&stringChannelInterruptFlags,
		&stringShapeShift,
		&spell.Description,
		&spell.Variables,
		&spell.MaxCumulativeStacks,
		&spell.MaxTargets,
		&iconId,
	)
	if err != nil {
		return spell, fmt.Errorf("scanning spell data: %w", err)
	}

	spell.Attributes, err = parseIntArrayField(stringAttr, 16)
	if err != nil {
		return spell, fmt.Errorf("parsing attributes args for spell %d (%s): %w", spell.ID, stringAttr, err)
	}
	spell.SpellClassMask, err = parseIntArrayField(stringClassMask, 4)
	if err != nil {
		return spell, fmt.Errorf("parsing classmask args for spell %d (%s): %w", spell.ID, stringClassMask, err)
	}

	spell.ProcTypeMask, err = parseIntArrayField(stringProcType, 2)
	if err != nil {
		return spell, fmt.Errorf("parsing ProcTypeMask args for spell %d (%s): %w", spell.ID, stringProcType, err)
	}
	spell.AuraInterruptFlags, err = parseIntArrayField(stringAuraIFlags, 2)
	if err != nil {
		return spell, fmt.Errorf("parsing stringAuraIFlags args for spell %d (%s): %w", spell.ID, stringAuraIFlags, err)
	}
	spell.ChannelInterruptFlags, err = parseIntArrayField(stringChannelInterruptFlags, 2)
	if err != nil {
		return spell, fmt.Errorf("parsing stringChannelInterruptFlags args for spell %d (%s): %w", spell.ID, stringChannelInterruptFlags, err)
	}
	spell.ShapeshiftMask, err = parseIntArrayField(stringShapeShift, 2)
	if err != nil {
		return spell, fmt.Errorf("parsing stringShapeShift args for spell %d (%s): %w", spell.ID, stringShapeShift, err)
	}

	spell.IconPath = iconsMap[iconId]
	return spell, nil
}

func LoadAndWriteSpells(dbHelper *DBHelper, inputsDir string) ([]dbc.Spell, error) {
	query := `
	SELECT DISTINCT
		sn.Name_lang,
		sn.ID,
		sm.SchoolMask,
		sm.Speed,
		sm.LaunchDelay,
		sm.MinDuration,
		COALESCE(ss.MaxScalingLevel, 0),
		COALESCE(ss.MinScalingLevel, 0),
		COALESCE(ss.ScalesFromItemLevel, 0),
		COALESCE(sl.SpellLevel, 0),
		COALESCE(sl.BaseLevel, 0),
		COALESCE(sl.MaxLevel, 0),
		COALESCE(sl.MaxPassiveAuraLevel, 0),
		COALESCE(sc.RecoveryTime, 0),
		COALESCE(sc.StartRecoveryTime, 0),
		COALESCE(sr.RangeMin_0, 0.0),
		COALESCE(sr.RangeMax_0, 0.0),
		COALESCE(sm."Attributes", ""),
		COALESCE(ssc.Flags, 0),
		COALESCE(ssc.MaxCharges, 0),
		COALESCE(ssc.ChargeRecoveryTime, 0),
		COALESCE(ssc.TypeMask, 0),
		COALESCE(scs.Category,0),
		COALESCE(sd.Duration,0),
		COALESCE(sao.ProcChance,0),
		COALESCE(sao.ProcCharges,0),
		COALESCE(sao.ProcTypeMask, ""),
		COALESCE(sao.ProcCategoryRecovery, 0),
		COALESCE(spm.BaseProcRate, 0),
		COALESCE(sei.EquippedItemClass, 0),
		COALESCE(sei.EquippedItemInvTypes, 0),
		COALESCE(sei.EquippedItemSubclass,0),
		COALESCE(ss.CastTimeMin, 0),
		COALESCE(sco.SpellClassMask, ""),
		COALESCE(sco.SpellClassSet, 0),
		COALESCE(si.AuraInterruptFlags, ""),
		COALESCE(si.ChannelInterruptFlags, ""),
		COALESCE(ssp.ShapeshiftMask, ""),
		COALESCE(s.Description_lang, ""),
		COALESCE(sdv.Variables, ""),
		COALESCE(sao.CumulativeAura, 0),
		COALESCE(str.MaxTargets, 0),
		COALESCE(sm.SpellIconFileDataID, 0)
		FROM Spell s
		LEFT JOIN SpellName sn ON s.ID = sn.ID
		LEFT JOIN SpellEffect se ON s.ID = se.SpellID
		LEFT JOIN SpellMisc sm ON s.ID = sm.SpellID
		LEFT JOIN SpellLevels sl ON s.ID = sl.SpellID
		LEFT JOIN SpellCooldowns sc ON s.ID = sc.SpellID
		LEFT JOIN SpellScaling ss ON s.ID = ss.SpellID
		LEFT JOIN SpellLabel slb ON s.ID = slb.SpellID
		LEFT JOIN SpellCategories scs ON s.ID = scs.SpellID
		LEFT JOIN SpellCategory ssc ON ssc.ID = scs.Category
		LEFT JOIN SpellDuration sd ON sm.DurationIndex = sd.ID
		LEFT JOIN SpellPower sp ON sp.SpellID = s.ID
		LEFT JOIN SpellInterrupts si ON si.SpellID = s.ID
		LEFT JOIN SpellEquippedItems sei ON sei.SpellID = s.ID
		LEFT JOIN SpellAuraOptions sao ON sao.SpellID = s.ID
		LEFT JOIN SpellClassOptions sco ON s.ID = sco.SpellID
		LEFT JOIN SpellShapeshift ssp ON ssp.SpellID = s.ID
		LEFT JOIN SpellXDescriptionVariables sxd ON s.ID = sxd.SpellID
		LEFT JOIN SpellDescriptionVariables sdv ON sdv.ID = sxd.SpellDescriptionVariablesID
		LEFT JOIN SpellTargetRestrictions str ON s.ID = str.SpellID
		LEFT JOIN SpellRange sr ON sr.ID = sm.RangeIndex
		LEFT JOIN SpellProcsPerMinute spm ON spm.ID = sao.SpellProcsPerMinuteID
		WHERE sco.SpellClassSet is not null
		GROUP BY s.ID
`

	spells, err := LoadRows(dbHelper.db, query, ScanSpells)
	if err != nil {
		return nil, fmt.Errorf("error loading spells: %w", err)
	}

	fmt.Println("Loaded spells:", len(spells))
	json, _ := json.Marshal(spells)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/spells.json", inputsDir), json); err != nil {
		panic(fmt.Sprintf("Error loading DBC data %v", err))
	}
	return spells, nil
}

func LoadAndWriteEnchantDescriptions(outputPath string, db *WowDatabase, instance *dbc.DBC) error {
	descriptions := make(map[int32]string)

	dataProvider := tooltip.DBCTooltipDataProvider{DBC: instance}
	for _, enchant := range db.Enchants {
		dbcEnch := instance.Enchants[int(enchant.EffectId)]
		tooltip, err := tooltip.ParseTooltip(dbcEnch.EffectName, dataProvider, int64(enchant.EffectId))
		if err != nil {
			fmt.Printf("Could not parse enchant (%d), '%s'\n", enchant.EffectId, dbcEnch.EffectName)
		} else {
			descriptions[enchant.EffectId] = tooltip.String()
		}
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(descriptions); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func ScanDropRow(rows *sql.Rows) (itemID int, ds *proto.DropSource, instanceName string, err error) {
	var (
		mask       int
		dropSource proto.DropSource
		jiName     string
	)
	err = rows.Scan(
		&itemID,
		&mask,
		&dropSource.NpcId,
		&dropSource.ZoneId,
		&dropSource.OtherName,
		&jiName,
	)
	if err != nil {
		return 0, nil, "", fmt.Errorf("scanning drop row: %w", err)
	}
	return itemID, &dropSource, jiName, nil
}

func LoadAndWriteDropSources(dbHelper *DBHelper, inputsDir string) (
	sourcesByItem map[int][]*proto.DropSource,
	namesByZone map[int]string,
	err error,
) {
	const query = `
		SELECT DISTINCT
		jei.ItemID,
		jei.DifficultyMask,
		je.ID                               AS NpcId,
		COALESCE(COALESCE(
			NULLIF(ji.AreaID, 0),
			at.ID
		), 0)                                AS ZoneId,
		je.Name_lang                     AS OtherName,
		ji.Name_lang
		FROM JournalEncounterItem AS jei
		INNER JOIN JournalEncounter AS je
		ON je.ID = jei.JournalEncounterID
		INNER JOIN JournalInstance AS ji
		ON ji.ID = je.JournalInstanceID
		LEFT JOIN AreaTable AS at
		ON (
			at.ZoneName       = ji.Name_lang
			OR at.AreaName_lang  = ji.Name_lang
		)
		GROUP BY jei.ItemID
    `

	rows, err := dbHelper.db.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("querying drop sources: %w", err)
	}
	defer rows.Close()

	sourcesByItem = make(map[int][]*proto.DropSource)
	namesByZone = make(map[int]string)

	for rows.Next() {
		itemID, ds, jiName, scanErr := ScanDropRow(rows)
		if scanErr != nil {
			return nil, nil, scanErr
		}
		sourcesByItem[itemID] = append(sourcesByItem[itemID], ds)
		namesByZone[int(ds.ZoneId)] = jiName
	}
	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterating drop rows: %w", err)
	}
	json, _ := json.Marshal(sourcesByItem)
	if err := dbc.WriteGzipFile(fmt.Sprintf("%s/dbc/dropSources.json", inputsDir), json); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
	return sourcesByItem, namesByZone, nil
}
