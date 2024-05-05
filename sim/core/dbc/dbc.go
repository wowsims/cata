package dbc

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
)

//go:embed generated/data/spell_effect.json
var spellEffectDataJson string

//go:embed generated/data/spell_data.json
var spellDataJson string

//go:embed generated/data/spell_power_data.json
var spellPowerDataJson string

//go:embed generated/GameTables/SpellScaling.txt
var spellScalingFile string

// DBC represents a database connection or access layer for game data.
type DBC struct {
	spellIndex       map[uint]*SpellData
	spellEffectIndex map[uint]*SpellEffectData
	SpellScalings    map[int]*SpellScaling
	classFamilyIndex map[uint][]*SpellData
	categoryMapping  SpellMapping
}

func NewDBC() *DBC {
	dbc := &DBC{
		SpellScalings:    make(map[int]*SpellScaling),
		spellEffectIndex: make(map[uint]*SpellEffectData),
		spellIndex:       make(map[uint]*SpellData),
		classFamilyIndex: make(map[uint][]*SpellData),
		categoryMapping: SpellMapping{
			Spells:  make(map[uint][]*SpellData),
			Effects: make(map[uint][]*SpellEffectData),
		},
	}
	dbc.initSpells() // Initialize spells immediately upon creation
	return dbc
}

type SpellScaling struct {
	Level  int
	Values map[proto.Class]float64
}
type SpellMapping struct {
	Spells  map[uint][]*SpellData
	Effects map[uint][]*SpellEffectData
}

func (m *SpellMapping) AddSpell(category uint, spell *SpellData) {
	m.Spells[category] = append(m.Spells[category], spell)
}

func (m *SpellMapping) AddEffect(category uint, effect *SpellEffectData) {
	m.Effects[category] = append(m.Effects[category], effect)
}

func (m *SpellMapping) SpellsByCategory(category uint) []*SpellData {
	return m.Spells[category]
}

func (m *SpellMapping) EffectsByCategory(category uint) []*SpellEffectData {
	return m.Effects[category]
}
func (dbc *DBC) LoadSpellScaling() error {
	scanner := bufio.NewScanner(strings.NewReader(spellScalingFile))
	dbc.SpellScalings = make(map[int]*SpellScaling)
	scanner.Scan() // Skip first line

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 14 {
			continue // consider handling or logging this situation
		}

		level, err := strconv.Atoi(parts[0])
		if err != nil {
			continue // consider handling or logging this situation
		}

		scaling := &SpellScaling{
			Level: level,
			Values: map[proto.Class]float64{
				proto.Class_ClassWarrior:     parseScalingValue(parts[1]),
				proto.Class_ClassPaladin:     parseScalingValue(parts[2]),
				proto.Class_ClassHunter:      parseScalingValue(parts[3]),
				proto.Class_ClassRogue:       parseScalingValue(parts[4]),
				proto.Class_ClassPriest:      parseScalingValue(parts[5]),
				proto.Class_ClassDeathKnight: parseScalingValue(parts[6]),
				proto.Class_ClassShaman:      parseScalingValue(parts[7]),
				proto.Class_ClassMage:        parseScalingValue(parts[8]),
				proto.Class_ClassWarlock:     parseScalingValue(parts[9]),
				proto.Class_ClassDruid:       parseScalingValue(parts[11]),
				proto.Class_ClassUnknown:     parseScalingValue(parts[12]),
			},
		}
		dbc.SpellScalings[level] = scaling
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
func (dbc *DBC) SpellScaling(class proto.Class, level int) float64 {
	if scaling, ok := dbc.SpellScalings[level]; ok {
		if value, ok := scaling.Values[class]; ok {
			return value
		}
	}
	return 0.0 // return a default or error value if not found
}
func parseScalingValue(value string) float64 {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0 // consider how to handle or log this error properly
	}
	return v
}
func (dbc *DBC) initSpells() {
	// Load spells
	err := dbc.loadSpellData()
	if err != nil {
		log.Fatalf("Failed to load spells: %v", err)
	}
	log.Println("Spells loaded successfully.")

	// Load spell effects
	err = dbc.loadSpellEffectData()
	if err != nil {
		log.Fatalf("Failed to load spell effects: %v", err)
	}
	log.Println("Spell effects loaded successfully.")

	err = dbc.LoadSpellScaling()
	if err != nil {
		log.Fatalf("Failed to load spell effects: %v", err)
	}
	log.Println("Spell scaling loaded successfully.")
}

func (dbc *DBC) loadSpellEffectData() error {
	var rawEffects [][]interface{}
	if err := json.Unmarshal([]byte(spellEffectDataJson), &rawEffects); err != nil {
		return err
	}

	dbc.spellEffectIndex = make(map[uint]*SpellEffectData)
	for _, raw := range rawEffects {
		effect := new(SpellEffectData)
		effect.ID = uint(raw[0].(float64))
		effect.SpellID = uint(raw[1].(float64))
		effect.Index = uint(raw[2].(float64))
		effect.Type = uint(raw[3].(float64))
		effect.Subtype = uint(raw[4].(float64))
		effect.ScalingType = int(raw[5].(float64))
		effect.MCoeff = raw[6].(float64)
		effect.MDelta = raw[7].(float64)
		effect.MUnk = raw[8].(float64)
		effect.SPCoeff = raw[9].(float64)
		effect.APCoeff = raw[10].(float64)
		effect.Amplitude = raw[11].(float64)
		effect.Radius = raw[12].(float64)
		effect.RadiusMax = raw[13].(float64)
		effect.BaseValue = raw[14].(float64)
		effect.MiscValue = int(raw[15].(float64))
		effect.MiscValue2 = int(raw[16].(float64))
		classFlagsRaw := raw[17].([]interface{})
		for i, val := range classFlagsRaw {
			effect.ClassFlags[i] = uint(val.(float64))
		}
		effect.TriggerSpellID = uint(raw[18].(float64))
		effect.MChain = raw[19].(float64)
		effect.PPComboPoints = raw[20].(float64)
		effect.RealPPL = raw[21].(float64)
		effect.Mechanic = uint(raw[22].(float64))
		effect.ChainTarget = int(raw[23].(float64))
		effect.Targeting1 = uint(raw[24].(float64))
		effect.Targeting2 = uint(raw[25].(float64))
		effect.MValue = raw[26].(float64)
		effect.PVPCoeff = raw[27].(float64)

		if spell, ok := dbc.spellIndex[effect.SpellID]; ok {
			spell.Effects = append(spell.Effects, effect)
			spell.EffectsCount += 1
		}

		if contains(effect.Subtype) {
			dbc.categoryMapping.AddEffect(uint(effect.MiscValue), effect)
		}

		// Cache the effect
		dbc.spellEffectIndex[effect.ID] = effect
	}

	return nil
}

func (dbc *DBC) loadSpellData() error {
	var rawSpells [][]interface{}
	if err := json.Unmarshal([]byte(spellDataJson), &rawSpells); err != nil {
		return err
	}

	dbc.spellIndex = make(map[uint]*SpellData)
	for _, raw := range rawSpells {
		spell := new(SpellData)
		spell.Name = raw[0].(string)
		spell.ID = uint(raw[1].(float64))
		spell.School = uint(raw[2].(float64))
		spell.PrjSpeed = raw[3].(float64)
		spell.PrjDelay = raw[4].(float64)
		spell.PrjMinDuration = raw[5].(float64)
		raceMaskStr := raw[6].(string) // Example: "-0x0000000000000001" or "0x0000000000000001"
		if len(raceMaskStr) >= 2 && (raceMaskStr[:2] == "0x" || (len(raceMaskStr) >= 3 && raceMaskStr[:3] == "-0x")) {
			// Remove the '0x' or '-0x' prefix properly before parsing
			baseIndex := 2
			if raceMaskStr[0] == '-' {
				baseIndex = 3 // Start after "-0x"
			}

			// Parse the hex string as an int64
			val, err := strconv.ParseInt(raceMaskStr[baseIndex:], 16, 64)
			if err != nil {
				log.Fatalf("Failed to parse RaceMask: %v", err)
			}

			// Convert int64 to int and assign it to spell.RaceMask
			spell.RaceMask = int(val)
		} else {
			log.Fatalf("Invalid raceMaskStr format: %v", raceMaskStr)
		}
		classMaskStr := raw[7].(string) // Example: "-0x00000001" or "0x00000080"
		if len(classMaskStr) >= 2 && (classMaskStr[:2] == "0x" || (len(classMaskStr) >= 3 && classMaskStr[:3] == "-0x")) {
			// Remove the '0x' or '-0x' prefix properly before parsing
			baseIndex := 2
			if classMaskStr[0] == '-' {
				baseIndex = 3 // Start after "-0x"
			}

			// Parse the hex string as an int64
			val, err := strconv.ParseInt(classMaskStr[baseIndex:], 16, 32) // Use 32-bit size for parsing
			if err != nil {
				log.Fatalf("Failed to parse ClassMask: %v", err)
			}

			// Convert int64 to int and assign it to spell.ClassMask
			spell.ClassMask = int(val)
		} else {
			log.Fatalf("Invalid ClassMaskStr format: %v", classMaskStr)
		}

		spell.MaxScalingLevel = int(raw[8].(float64))
		spell.SpellLevel = int(raw[9].(float64))
		spell.MaxLevel = int(raw[10].(float64))
		spell.ReqMaxLevel = int(raw[11].(float64))
		spell.MinRange = raw[12].(float64)
		spell.MaxRange = raw[13].(float64)
		spell.Cooldown = time.Duration(raw[14].(float64)) * time.Millisecond
		spell.GCD = time.Duration(raw[15].(float64)) * time.Millisecond
		spell.CategoryCooldown = time.Duration(raw[16].(float64)) * time.Millisecond
		spell.Charges = uint(raw[17].(float64))
		spell.ChargeCooldown = time.Duration(raw[18].(float64)) * time.Millisecond
		spell.Category = uint(raw[19].(float64))
		spell.DmgClass = uint(raw[20].(float64))
		spell.MaxTargets = int(raw[21].(float64))
		spell.Duration = time.Duration(raw[22].(float64)) * time.Millisecond
		spell.MaxStack = uint(raw[23].(float64))
		spell.ProcChance = uint(raw[24].(float64))
		spell.ProcCharges = int(raw[25].(float64))
		spell.ProcFlags = uint64(raw[26].(float64))
		spell.InternalCooldown = time.Duration(raw[27].(float64)) * time.Millisecond
		spell.RPPM = raw[28].(float64)
		spell.EquippedClass = uint(raw[29].(float64))
		spell.EquippedInvtypeMask = uint(raw[30].(float64))
		spell.EquippedSubclassMask = uint(raw[31].(float64))
		spell.CastTime = time.Duration(raw[32].(float64)) * time.Millisecond
		// Handle attributes and class flags, assuming they are formatted as JSON arrays of integers
		attributeRaw := raw[33].([]interface{})
		for i, val := range attributeRaw {
			spell.Attributes[i] = uint(val.(float64))
		}
		classFlagsRaw := raw[34].([]interface{})
		for i, val := range classFlagsRaw {
			spell.ClassFlags[i] = uint(val.(float64))
		}
		spell.ClassFlagsFamily = uint(raw[35].(float64))
		spell.StanceMask = uint(raw[36].(float64))
		spell.Mechanic = uint(raw[37].(float64))
		spell.PowerID = uint(raw[38].(float64))
		spell.EssenceID = uint(raw[39].(float64))
		spell.Effects = make([]*SpellEffectData, 0)
		spell.Power = make([]*SpellPowerData, 0)
		spell.PowerCount = 0
		spell.EffectsCount = 0

		// Cache the spell
		dbc.spellIndex[spell.ID] = spell
		if spell.ClassFlagsFamily != 0 {
			if _, exists := dbc.classFamilyIndex[spell.ClassFlagsFamily]; !exists {
				dbc.classFamilyIndex[spell.ClassFlagsFamily] = []*SpellData{}
			}
			dbc.classFamilyIndex[spell.ClassFlagsFamily] = append(dbc.classFamilyIndex[spell.ClassFlagsFamily], spell)
		}

		// Index by category, if applicable
		if spell.Category != 0 {
			dbc.categoryMapping.AddSpell(spell.Category, spell)
		}
	}

	var rawJson [][]interface{} // To accommodate the large array of arrays
	if err := json.Unmarshal([]byte(spellPowerDataJson), &rawJson); err != nil {
		log.Fatal(err)
	}

	// Assuming the first item in rawJson is the large array of SpellPowerData arrays
	for _, item := range rawJson[0] {
		entry := item.([]interface{}) // Each item is an array representing one SpellPowerData
		spellPowerData := SpellPowerData{
			ID:             uint(entry[0].(float64)),
			SpellID:        uint(entry[1].(float64)),
			AuraID:         uint(entry[2].(float64)),
			PowerType:      int(entry[3].(float64)),
			Cost:           int(entry[4].(float64)),
			CostMax:        int(entry[5].(float64)),
			CostPerTick:    int(entry[6].(float64)),
			PctCost:        entry[7].(float64),
			PctCostMax:     entry[8].(float64),
			PctCostPerTick: entry[9].(float64),
		}
		if dbc.spellIndex[spellPowerData.SpellID] != nil {
			dbc.spellIndex[spellPowerData.SpellID].Power = append(dbc.spellIndex[spellPowerData.SpellID].Power, &spellPowerData)
			dbc.spellIndex[spellPowerData.SpellID].PowerCount += 1
		}
	}
	return nil
}

// FetchSpellEffect retrieves a spell effect based on an ID.
func (dbc *DBC) FetchSpellEffect(effectID uint) *SpellEffectData {
	return dbc.spellEffectIndex[effectID]
}

// FetchSpell retrieves a spell based on an ID. It returns nil and an error if no spell is found.
func (dbc *DBC) FetchSpell(spellID uint) (*SpellData, error) {
	spell, found := dbc.spellIndex[spellID]
	if !found {
		return nil, fmt.Errorf("no spell found with ID %d", spellID)
	}
	return spell, nil
}

// EffectAverage calculates the average value of an effect at a given level.
func (dbc *DBC) EffectAverage(e *SpellEffectData, level int) float64 {
	if e == nil || level <= 0 || level > MAX_SCALING_LEVEL {
		return 0
	}

	scale := e.ScalingClass()

	// Todo: DF stuff?
	// if scale == PLAYER_NONE && e.Spell.MaxScalingLevel > 0 {
	// 	scale = PLAYER_SPECIAL_SCALE8
	// }

	if e.MCoeff != 0 && scale != proto.Class_ClassUnknown {
		if e.Spell.MaxScalingLevel > 0 {
			level = min(level, e.Spell.MaxScalingLevel)
		}
		return e.MCoeff * dbc.SpellScaling(scale, level)
	} else if e.RealPPL != 0 {
		if e.Spell.MaxLevel > 0 {
			return e.BaseValue + (float64(min(level, e.Spell.MaxLevel))-float64(e.Spell.SpellLevel))*e.RealPPL
		}
		return e.BaseValue + (float64(level)-float64(e.Spell.SpellLevel))*e.RealPPL
	}
	return e.BaseValue
}

// EffectDelta calculates the delta value for an effect.
func (dbc *DBC) EffectDelta(e *SpellEffectData, level int) float64 {
	if e == nil || level <= 0 || level > MAX_SCALING_LEVEL {
		return 0
	}

	if e.MDelta != 0 && e.ScalingClass() != proto.Class_ClassUnknown {
		return e.MDelta
	}

	return 0
}
func (dbc *DBC) EffectBonusById(effectId uint, level int) float64 {
	spellEffectData := dbc.spellEffectIndex[effectId]
	return dbc.EffectBonus(spellEffectData, level)
}

// EffectBonus calculates additional bonus effects.
func (dbc *DBC) EffectBonus(e *SpellEffectData, level int) float64 {
	if e == nil || level <= 0 || level > MAX_SCALING_LEVEL {
		return 0
	}

	if e.MiscValue != 0 && e.ScalingClass() != proto.Class_ClassUnknown {
		scalingLevel := min(level, e.Spell.MaxScalingLevel)
		mScale := dbc.SpellScaling(e.ScalingClass(), scalingLevel)
		return float64(e.MiscValue) * mScale
	}

	return float64(e.MiscValue)
}

// EffectMin calculates the minimum value for an effect.
func (dbc *DBC) EffectMin(e *SpellEffectData, level int) float64 {
	if e == nil || level <= 0 || level > MAX_SCALING_LEVEL {
		return 0
	}

	avg := dbc.EffectAverage(e, level)
	delta := dbc.EffectDelta(e, level)
	result := avg - (avg * delta / 2)

	if e.Type == E_WEAPON_PERCENT_DAMAGE {
		result *= 0.01
	}

	return result
}

// EffectMax calculates the maximum value for an effect.
func (dbc *DBC) EffectMax(e *SpellEffectData, level int) float64 {
	if e == nil || level <= 0 || level > MAX_SCALING_LEVEL {
		return 0
	}

	avg := dbc.EffectAverage(e, level)
	delta := dbc.EffectDelta(e, level)
	result := avg + (avg * delta / 2)

	if e.Type == E_WEAPON_PERCENT_DAMAGE {
		result *= 0.01
	}

	return result
}

// Returns effects affecting spells by checking class family flags
func (dbc *DBC) EffectAffectsSpells(family uint, effect *SpellEffectData) []*SpellData {
	var affectedSpells []*SpellData
	if family == 0 {
		return affectedSpells
	}

	index, exists := dbc.classFamilyIndex[family]
	if !exists {
		return affectedSpells // Early return if no spells are indexed under this family
	}

	// Using a map to check if a spell has already been added
	spellsMap := make(map[*SpellData]bool)

	for _, s := range index {
		for j := uint(0); j < NUM_CLASS_FAMILY_FLAGS; j++ {
			if effect.ClassFlags[j]&s.ClassFlags[j] != 0 { // Corrected bitwise AND operation and condition check
				if _, found := spellsMap[s]; !found {
					spellsMap[s] = true
					affectedSpells = append(affectedSpells, s)
				}
			}
		}
	}

	return affectedSpells
}

// AffectedBy fetches SpellEffectData for a given category.
func (dbc *DBC) EffectsByCategory(category uint) []*SpellEffectData {
	if effects, ok := dbc.categoryMapping.Effects[category]; ok {
		return effects
	}
	return nil
}

// Retrieve effects based on categories that affect a spell
func (dbc *DBC) EffectCategoriesAffectingSpell(spell *SpellData) []*SpellEffectData {
	var effects []*SpellEffectData
	categoryEffects := dbc.EffectsByCategory(spell.Category) // Fetch effects for the spell's category

	effectMap := make(map[uint]bool) // Using a map to track already included effect IDs for de-duplication

	for _, effect := range categoryEffects {
		if !effectMap[effect.ID] { // Check if the effect has already been added
			effects = append(effects, effect)
			effectMap[effect.ID] = true // Mark this effect as added
		}
	}

	return effects
}
func (dbc *DBC) EffectsAffectingSpell(spell *SpellData) []*SpellEffectData {
	var affectingEffects []*SpellEffectData
	if spell.ClassFlagsFamily == 0 {
		return affectingEffects
	}

	index, ok := dbc.classFamilyIndex[spell.ClassFlagsFamily]
	if !ok {
		return affectingEffects // Return if there is no entry for the class family
	}

	effectMap := make(map[*SpellEffectData]bool) // Use a map to track unique effects

	for _, s := range index {
		if s.ID == spell.ID {
			continue // Skip the spell itself
		}
		for _, effect := range s.Effects {
			// Assume ClassFlags is an array of uint32 and contains one for each class family flag
			for j := uint(0); j < NUM_CLASS_FAMILY_FLAGS; j++ {
				// Correct bitwise operation to check flags
				if effect.ClassFlags[j]&spell.ClassFlags[j] != 0 {
					if !effectMap[effect] { // Check if already added
						affectingEffects = append(affectingEffects, effect)
						effectMap[effect] = true
					}
				}
			}
		}
	}
	return affectingEffects
}
