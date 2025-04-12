package dbc

import (
	"math"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type SpellEffect struct {
	ID                             int
	DifficultyID                   int
	EffectIndex                    int
	EffectType                     SpellEffectType
	EffectAmplitude                float64
	EffectAttributes               int
	EffectAura                     EffectAuraType
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

func (effect *SpellEffect) ParseStatEffect() *stats.Stats {
	stats := &stats.Stats{}
	if effect.EffectAura == A_MOD_STAT && effect.EffectType == E_APPLY_AURA {
		stat, success := MapMainStatToStat(effect.EffectMiscValues[0])
		if success == false {
			// error
		}
		stats[stat] = float64(effect.EffectBasePoints)
		return stats
	}

	if effect.EffectMiscValues[0] == -1 &&
		effect.EffectType == E_APPLY_AURA &&
		effect.EffectAura == A_MOD_STAT {
		// Apply bonus to all stats
		stats[proto.Stat_StatAgility] = float64(effect.EffectBasePoints)
		stats[proto.Stat_StatIntellect] = float64(effect.EffectBasePoints)
		stats[proto.Stat_StatSpirit] = float64(effect.EffectBasePoints)
		stats[proto.Stat_StatStamina] = float64(effect.EffectBasePoints)
		stats[proto.Stat_StatStrength] = float64(effect.EffectBasePoints)
		return stats
	}
	school := SpellSchool(effect.EffectMiscValues[0])
	if effect.EffectAura == A_MOD_TARGET_RESISTANCE {
		if school == SPELL_PENETRATION {
			stats[proto.Stat_StatSpellPenetration] += math.Abs(float64(effect.EffectBasePoints))
			return stats
		}
	}

	if effect.EffectAura == A_OBS_MOD_HEALTH {
		return stats
	}

	if effect.EffectAura == A_MOD_RESISTANCE {
		if school.Has(FIRE) {
			stats[proto.Stat_StatFireResistance] += float64(effect.EffectBasePoints)
		}
		if school.Has(ARCANE) {
			stats[proto.Stat_StatArcaneResistance] += float64(effect.EffectBasePoints)
		}
		if school.Has(NATURE) {
			stats[proto.Stat_StatNatureResistance] += float64(effect.EffectBasePoints)
		}
		if school.Has(FROST) {
			stats[proto.Stat_StatFrostResistance] += float64(effect.EffectBasePoints)
		}
		if school.Has(SHADOW) {
			stats[proto.Stat_StatShadowResistance] += float64(effect.EffectBasePoints)
		}
		if school.Has(PHYSICAL) {
			stats[proto.Stat_StatArmor] += float64(effect.EffectBasePoints)
		}
	}
	if effect.EffectAura == A_MOD_RATING {
		matching := getMatchingRatingMods(effect.EffectMiscValues[0])
		for _, rating := range matching {
			statMod := RatingModToStat[rating]
			if statMod != -1 {
				stats[statMod] = float64(effect.EffectBasePoints)
			}
		}
	}

	if effect.EffectAura == A_MOD_INCREASE_ENERGY {
		stats[proto.Stat_StatMana] = float64(effect.EffectBasePoints)
	}

	if effect.EffectAura == A_PERIODIC_TRIGGER_SPELL && effect.EffectAuraPeriod == 10000 { // Make sure if its a food
		subEffects := dbcInstance.SpellEffects[effect.EffectTriggerSpell]
		for _, subEffect := range subEffects {
			stats.AddInplace(subEffect.ParseStatEffect())
		}
	}

	return stats
}
