package dbc

import (
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// ItemEffect represents an item effect in the game.
type ItemEffect struct {
	ID                   int // Effect ID
	LegacySlotIndex      int // Legacy slot index
	TriggerType          int // Trigger type
	Charges              int // Number of charges
	CoolDownMSec         int // Cooldown in milliseconds
	CategoryCoolDownMSec int // Category cooldown in milliseconds
	SpellCategoryID      int // Spell category ID
	SpellID              int // Spell ID
	ChrSpecializationID  int // Character specialization ID
	ParentItemID         int // Parent item ID
}

// ToMap returns a generic representation of the effect.
func (e *ItemEffect) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":                   e.ID,
		"LegacySlotIndex":      e.LegacySlotIndex,
		"TriggerType":          e.TriggerType,
		"Charges":              e.Charges,
		"CoolDownMSec":         e.CoolDownMSec,
		"CategoryCoolDownMSec": e.CategoryCoolDownMSec,
		"SpellCategoryID":      e.SpellCategoryID,
		"SpellID":              e.SpellID,
		"ChrSpecializationID":  e.ChrSpecializationID,
		"ParentItemID":         e.ParentItemID,
	}
}

func GetItemEffect(effectId int) ItemEffect {
	return dbcInstance.ItemEffects[effectId]
}

func makeBaseProto(e *ItemEffect, statsSpellID int) *proto.ItemEffect {
	sp := dbcInstance.Spells[e.SpellID]
	base := &proto.ItemEffect{
		BuffId:           int32(e.SpellID),
		BuffName:         sp.NameLang,
		EffectDurationMs: int32(sp.Duration),
		ScalingOptions:   make(map[int32]*proto.ScalingItemEffectProperties),
	}
	// override duration if stats spell defines its own
	if dur := dbcInstance.Spells[statsSpellID].Duration; dur > 0 {
		base.EffectDurationMs = int32(dur)
	}
	return base
}

func assignTrigger(e *ItemEffect, statsSpellID int, pe *proto.ItemEffect) {
	spTop := dbcInstance.Spells[e.SpellID]
	statsSP := dbcInstance.Spells[statsSpellID]
	switch resolveTriggerType(e.TriggerType, e.SpellID) {
	case ITEM_SPELLTRIGGER_ON_USE:
		pe.Effect = &proto.ItemEffect_OnUse{OnUse: &proto.OnUseEffect{
			CooldownMs:         int32(e.CoolDownMSec),
			CategoryId:         int32(e.SpellCategoryID),
			CategoryCooldownMs: int32(e.CategoryCoolDownMSec),
		}}
	case ITEM_SPELLTRIGGER_CHANCE_ON_HIT:
		proc := &proto.ProcEffect{
			ProcChance: float64(spTop.ProcChance) / 100,
			IcdMs:      spTop.ProcCategoryRecovery,
			Ppm:        spTop.Rppm,
			RppmScale:  uint32(realPpmScale(&spTop)),
		}
		// If proc chance is above 100 something weird is happening so we set ppm to 1 since we cant accurately proc it 100% of the time
		if spTop.ProcChance == 0 || spTop.ProcChance > 100 {
			proc.Ppm = 1
			proc.ProcChance = 0
		}
		pe.BuffId = statsSP.ID
		pe.BuffName = statsSP.NameLang
		pe.Effect = &proto.ItemEffect_Proc{Proc: proc}
	}
}

func (e *ItemEffect) ToProto(itemLevel int, levelState proto.ItemLevelState) (*proto.ItemEffect, bool) {
	statsSpellID := resolveStatsSpell(e.SpellID)

	pe := makeBaseProto(e, statsSpellID)
	assignTrigger(e, statsSpellID, pe)

	// build scaling properties and skip if empty
	props := buildScalingProps(statsSpellID, itemLevel)

	if len(props.Stats) == 0 {
		return nil, false
	}

	pe.ScalingOptions[int32(levelState)] = props

	return pe, true
}

func resolveStatsSpell(spellID int) int {
	for _, se := range dbcInstance.SpellEffects[spellID] {
		switch se.EffectAura {
		case A_MOD_STAT, A_MOD_RATING, A_MOD_RANGED_ATTACK_POWER, A_MOD_ATTACK_POWER, A_MOD_DAMAGE_DONE, A_MOD_TARGET_RESISTANCE, A_MOD_RESISTANCE, A_MOD_INCREASE_ENERGY,
			A_MOD_INCREASE_HEALTH_2, A_PERIODIC_TRIGGER_SPELL:
			return spellID
		}
	}
	// If we cant resolve the spell in the first loop, we follow proc triggers downwards
	for _, se := range dbcInstance.SpellEffects[spellID] {
		if se.EffectAura == A_PROC_TRIGGER_SPELL {
			return resolveStatsSpell(se.EffectTriggerSpell)
		}
	}
	return spellID
}

func resolveTriggerType(topType, spellID int) int {
	if topType == ITEM_SPELLTRIGGER_ON_USE || topType == ITEM_SPELLTRIGGER_CHANCE_ON_HIT {
		return topType
	}
	for _, se := range dbcInstance.SpellEffects[spellID] {
		if se.EffectAura == A_PROC_TRIGGER_SPELL {
			return resolveTriggerType(ITEM_SPELLTRIGGER_CHANCE_ON_HIT, se.EffectTriggerSpell)
		}
	}
	return topType
}

func buildScalingProps(spellID, itemLevel int) *proto.ScalingItemEffectProperties {
	total := collectStats(spellID, itemLevel)
	return &proto.ScalingItemEffectProperties{Stats: total.ToProtoMap()}
}

func collectStats(spellID, itemLevel int) stats.Stats {
	var total stats.Stats

	var emptyStats = stats.Stats{}
	visited := make(map[int]bool)

	var recurse func(int)
	recurse = func(id int) {
		if visited[id] {
			return
		}
		visited[id] = true

		sp := dbcInstance.Spells[id]
		for _, se := range dbcInstance.SpellEffects[id] {
			s := se.ParseStatEffect(sp.HasAttributeAt(11, 0x4), itemLevel)
			if s != nil && *s != emptyStats {
				total.AddInplace(s)
			} else if se.EffectAura == A_PROC_TRIGGER_SPELL {
				recurse(se.EffectTriggerSpell)
			}
		}
	}

	recurse(spellID)
	return total
}

func ParseItemEffects(itemID, itemLevel int, levelState proto.ItemLevelState) []*proto.ItemEffect {
	raw := dbcInstance.ItemEffectsByParentID[itemID]
	out := make([]*proto.ItemEffect, 0, len(raw))
	for _, ie := range raw {
		if pe, ok := ie.ToProto(itemLevel, levelState); ok {
			out = append(out, pe)
		}
	}
	return out
}

// Parses a UIItem and loops through Scaling Options for that item.
func MergeItemEffectsForAllStates(parsed *proto.UIItem) *proto.ItemEffect {
	// pick a base effect that has stats if there is more than one effect on the item
	var baseEff *ItemEffect
	for i := range dbcInstance.ItemEffectsByParentID[int(parsed.Id)] {

		e := &dbcInstance.ItemEffectsByParentID[int(parsed.Id)][i]
		props := buildScalingProps(resolveStatsSpell(e.SpellID), int(parsed.ScalingOptions[int32(proto.ItemLevelState_Base)].Ilvl))
		if len(props.Stats) > 0 {
			baseEff = e
			break
		}
	}
	if baseEff == nil {
		return nil
	}
	statsSpellID := resolveStatsSpell(baseEff.SpellID)
	pe := makeBaseProto(baseEff, statsSpellID)
	assignTrigger(baseEff, statsSpellID, pe)

	// add scaling for each saved state
	for state, opt := range parsed.ScalingOptions {
		ilvl := int(opt.Ilvl)
		props := buildScalingProps(statsSpellID, ilvl)
		pe.ScalingOptions[state] = props

		if proc := pe.GetProc(); proc != nil {
			spell := dbcInstance.Spells[statsSpellID]
			_, specMods := realPpmModifier(&spell, ilvl)
			proc.SpecModifiers = specMods
		}
	}

	return pe
}

func realPpmScale(spell *Spell) proto.RppmScalingType {
	scale := proto.RppmScalingType_RppmScalingTypeNone
	for _, mod := range spell.RppmModifiers {
		switch mod.ModifierType {
		case RPPMModifierHaste:
			scale |= proto.RppmScalingType_RppmScalingTypeHaste
		case RPPMModifierCrit:
			scale |= proto.RppmScalingType_RppmScalingTypeCrit
		}
	}
	return scale
}

func realPpmModifier(spell *Spell, itemLevel int) (float64, map[int32]float64) {
	specModifier := make(map[int32]float64)
	ilvlModifier := 1.0
	for _, mod := range spell.RppmModifiers {
		switch mod.ModifierType {
		case RPPMModifierSpec:
			spec := SpecFromID(mod.Param)
			specModifier[int32(spec)] = 1.0 * (1.0 + mod.Coeff)

		case RPPMModifierIlevel:
			basePoints := dbcInstance.RandomPropertiesByIlvl[int(mod.Param)][proto.ItemQuality_ItemQualityRare][0]
			ilvlPoints := dbcInstance.RandomPropertiesByIlvl[itemLevel][proto.ItemQuality_ItemQualityRare][0]
			if basePoints != ilvlPoints {
				ilvlModifier *= 1.0 + ((float64(ilvlPoints)/float64(basePoints))-1.0)*mod.Coeff
			}
		}
	}
	return ilvlModifier, specModifier
}
