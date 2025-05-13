package dbc

import (
	"github.com/wowsims/cata/sim/core"
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

var emptyStats = stats.Stats{}

func (e *ItemEffect) ToProto(itemLevel int, levelState proto.ItemLevelState) *proto.ItemEffect {
	pe := newProtoShell(e)

	statsSpellID, ilvlRppmMod, _ := applyTrigger(e, pe, itemLevel)

	pe.ScalingOptions[int32(levelState)] = buildScalingProps(statsSpellID, itemLevel)
	if ilvlRppmMod != 0.0 && ilvlRppmMod != 1.0 {
		pe.ScalingOptions[int32(levelState)].RppmIlvlModifier = ilvlRppmMod
	}
	return pe
}

func newProtoShell(e *ItemEffect) *proto.ItemEffect {
	sp := dbcInstance.Spells[e.SpellID]
	return &proto.ItemEffect{
		BuffId:         int32(e.SpellID),
		BuffName:       sp.NameLang,
		Type:           proto.ItemEffectType_EffectTypeNone,
		EffectDuration: int32(sp.Duration) / 1000,
		ScalingOptions: make(map[int32]*proto.ScalingItemEffectProperties),
	}
}

func applyTrigger(e *ItemEffect, pe *proto.ItemEffect, itemLevel int) (int, float64, bool) {
	trig, statsSpellID := resolveTrigger(e.TriggerType, e.SpellID)
	sp := dbcInstance.Spells[statsSpellID]
	if sp.Duration > 0 {
		pe.EffectDuration = sp.Duration / 1000
	}
	switch trig {
	case ITEM_SPELLTRIGGER_ON_USE:
		pe.Type = proto.ItemEffectType_EffectTypeOnUse
		pe.Effect = &proto.ItemEffect_OnUse{
			OnUse: &proto.OnUseEffect{
				Cooldown:         int32(e.CoolDownMSec / 1000),
				CategoryId:       int32(e.SpellCategoryID),
				CategoryCooldown: int32(e.CategoryCoolDownMSec / 1000),
			},
		}
		return statsSpellID, 0, true
	case ITEM_SPELLTRIGGER_CHANCE_ON_HIT:
		// For procchance and ICD we always use the original spell id
		spTop := dbcInstance.Spells[e.SpellID]
		effect := &proto.ProcEffect{
			ProcChance: float64(spTop.ProcChance) / 100,
			Icd:        int32(spTop.ProcCategoryRecovery / 1000),
			Ppm:        spTop.Rppm,
			RppmScale:  int32(realPpmScale(spTop)),
		}
		// On procs we want the lower name though
		pe.BuffName = sp.NameLang
		pe.BuffId = sp.ID
		ilvlMod, specMods := realPpmModifier(spTop, itemLevel)
		effect.SpecModifiers = specMods

		pe.Type = proto.ItemEffectType_EffectTypeProc
		pe.Effect = &proto.ItemEffect_Proc{
			Proc: effect,
		}
		return statsSpellID, ilvlMod, true
	default:
		// leave as NONE
	}

	return statsSpellID, 0, false
}

func realPpmScale(spell Spell) int {
	scale := 0
	for _, mod := range spell.RppmModifiers {
		switch mod.ModifierType {
		case RPPMModifierHaste:
			scale |= core.RPPM_HASTE
		case RPPMModifierCrit:
			scale |= core.RPPM_CRIT
		}
	}
	return scale
}

func realPpmModifier(spell Spell, itemLevel int) (float64, map[int32]float64) {
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

func resolveTrigger(topType, spellID int) (triggerType, statsSpellID int) {
	if topType == ITEM_SPELLTRIGGER_ON_USE || topType == ITEM_SPELLTRIGGER_CHANCE_ON_HIT {
		return topType, spellID
	}
	for _, se := range dbcInstance.SpellEffects[spellID] {
		if se.EffectAura == A_PROC_TRIGGER_SPELL {
			// stats come from the triggered spell
			return resolveTrigger(ITEM_SPELLTRIGGER_CHANCE_ON_HIT, se.EffectTriggerSpell)
		}
	}
	return topType, spellID
}

func buildScalingProps(spellID, itemLevel int) *proto.ScalingItemEffectProperties {
	total := collectStats(spellID, itemLevel)
	return &proto.ScalingItemEffectProperties{Stats: total.ToProtoMap()}
}

func collectStats(spellID, itemLevel int) stats.Stats {
	var total stats.Stats
	visited := make(map[int]bool)

	var recurse func(int)
	recurse = func(id int) {
		if visited[id] {
			return
		}
		visited[id] = true

		sp := dbcInstance.Spells[id]
		for _, se := range dbcInstance.SpellEffects[id] {
			if s := se.ParseStatEffect(sp.HasAttributeAt(11, 0x4), itemLevel); s != &emptyStats {
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
		out = append(out, ie.ToProto(itemLevel, levelState))
	}
	return out
}

func MergeItemEffectsForAllStatesNew(parsed *proto.UIItem) *proto.ItemEffect {
	itemID := int(parsed.Id)
	raws := dbcInstance.ItemEffectsByParentID[itemID]

	merged := &proto.ItemEffect{}
	statsSpellIDs := make([]int, len(raws))
	baseMods := make([]float64, len(raws))

	for idx, ie := range raws {
		pe := newProtoShell(&ie)
		spellID, ilvlMod, success := applyTrigger(&ie, pe, 0)
		if !success {
			continue
		}
		merged = pe
		statsSpellIDs[idx] = spellID
		baseMods[idx] = ilvlMod
	}

	for key, props := range parsed.ScalingOptions {
		ilvl := int(props.Ilvl)

		for idx := range raws {
			spellID := statsSpellIDs[idx]

			ilvlMod, specMods := realPpmModifier(dbcInstance.Spells[spellID], ilvl)

			scaling := buildScalingProps(spellID, ilvl)
			if ilvlMod != 0 && ilvlMod != 1 {
				scaling.RppmIlvlModifier = ilvlMod
			}
			if proc := merged.GetProc(); proc != nil {
				proc.SpecModifiers = specMods
			}
			merged.ScalingOptions[key] = scaling
		}
	}

	return merged
}
