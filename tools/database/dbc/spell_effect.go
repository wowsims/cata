package dbc

import (
	"math"
	"slices"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const MAX_SCALING_LEVEL = 100
const BASE_LEVEL = 90

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
	EffectMiscValues      []int     // from EffectMiscValue_0, EffectMiscValue_1
	EffectMinRange        []float64 // from EffectRadiusIndex_0, EffectRadiusIndex_1
	EffectMaxRange        []float64
	EffectSpellClassMasks []int // from EffectSpellClassMask_0, EffectSpellClassMask_1, EffectSpellClassMask_2, EffectSpellClassMask_3
	ImplicitTargets       []int // from ImplicitTarget_0, ImplicitTarget_1
	SpellID               int
	ScalingType           int
}

func (se *SpellEffect) ToProto() *proto.SpellEffect {
	spellEffect := &proto.SpellEffect{
		Id:            int32(se.ID),
		SpellId:       int32(se.SpellID),
		Index:         int32(se.EffectIndex),
		Type:          proto.EffectType(se.EffectType),
		EffectSpread:  se.Delta(BASE_LEVEL, BASE_LEVEL), // Todo: something weird here only true for pots?
		MinEffectSize: se.Min(BASE_LEVEL, BASE_LEVEL),
	}
	if spellEffect.EffectSpread == 0 {
		spellEffect.EffectSpread = float64(se.EffectDieSides)
	}
	switch se.EffectType {
	case E_ENERGIZE:
		spellEffect.MiscValue0 = &proto.SpellEffect_ResourceType{ResourceType: MapPowerTypeEnumToResourceType[int32(se.EffectMiscValues[0])]}
	case E_HEAL:
		spellEffect.MiscValue0 = &proto.SpellEffect_ResourceType{ResourceType: proto.ResourceType_ResourceTypeHealth}
	}

	return spellEffect
}

func (s *SpellEffect) GetRadiusMax() float64 {
	return math.Max(s.EffectMaxRange[0], s.EffectMaxRange[1])
}

func (s *SpellEffect) GetRadiusMin() float64 {
	return math.Min(s.EffectMinRange[0], s.EffectMinRange[1])
}

func (s *SpellEffect) ScalingClass() proto.Class {
	switch s.ScalingType {
	case 1:
		return proto.Class_ClassWarrior
	case 2:
		return proto.Class_ClassPaladin
	case 3:
		return proto.Class_ClassHunter
	case 4:
		return proto.Class_ClassRogue
	case 5:
		return proto.Class_ClassPriest
	case 6:
		return proto.Class_ClassDeathKnight
	case 7:
		return proto.Class_ClassShaman
	case 8:
		return proto.Class_ClassMage
	case 9:
		return proto.Class_ClassWarlock
	case 10:
		return proto.Class_ClassMonk
	case 11:
		return proto.Class_ClassDruid
	default:
		return proto.Class_ClassUnknown
	}
}
func (s *SpellEffect) Delta(pLevel int, level int) float64 {
	if level > 90 {
		level = 90
	}

	var mScale float64
	spell := dbcInstance.Spells[s.SpellID]
	if s.Variance != 0 && s.ScalingClass() != 0 {
		scalingLevel := level
		if scalingLevel == 0 {
			scalingLevel = pLevel
		}
		if spell.MaxScalingLevel > 0 {
			scalingLevel = min(scalingLevel, spell.MaxScalingLevel)
		}
		mScale = dbcInstance.SpellScaling(s.ScalingClass(), scalingLevel)
	}

	return s.scaledDelta(mScale)
}

// func (s *SpellEffect) Bonus(dbc *DBC, pLevel int, level int) float64 {
// 	if level == 0 {
// 		level = pLevel
// 	}
// 	return dbc.EffectBonusById(s.GetSpell(dbc).ID, level)
// }

func (s *SpellEffect) Average(pLevel int, level int) float64 {
	if level == 0 {
		level = pLevel
	}

	scale := s.ScalingClass()
	spell := dbcInstance.Spells[s.ID]

	if s.Coefficient != 0 && scale != proto.Class_ClassUnknown {
		if spell.MaxScalingLevel > 0 {
			level = min(level, spell.MaxScalingLevel)
		}
		scaler := dbcInstance.SpellScaling(scale, level)
		value := s.Coefficient * scaler
		return value
	} else if s.EffectRealPointsPerLevel != 0 {
		if spell.MaxLevel > 0 {
			return float64(s.EffectBasePoints) + float64(min(level, spell.MaxLevel)-spell.SpellLevel)*s.EffectRealPointsPerLevel
		}
		return float64(s.EffectBasePoints) + float64(level-spell.SpellLevel)*s.EffectRealPointsPerLevel
	}
	return float64(s.EffectBasePoints)
}

// Minimum value calculation for player
func (s *SpellEffect) Min(pLevel int, level int) float64 {
	return s.scaledMin(s.Average(pLevel, level), s.Delta(pLevel, level))
}

// Maximum value calculation for player
func (s *SpellEffect) Max(pLevel int, level int) float64 {
	return s.scaledMax(s.Average(pLevel, level), s.Delta(pLevel, level))
}

func (s *SpellEffect) scaledDelta(budget float64) float64 {
	if s.Variance != 0 && budget > 0 {
		return s.Coefficient * float64(s.Variance) * budget
	}
	return 0
}

// Scaled minimum calculation
func (s *SpellEffect) scaledMin(avg, delta float64) float64 {
	result := avg - delta/2
	if s.EffectType == E_WEAPON_PERCENT_DAMAGE {
		result *= 0.01
	}
	return result
}

// Scaled maximum calculation
func (s *SpellEffect) scaledMax(avg, delta float64) float64 {
	result := avg + delta/2
	if s.EffectType == E_WEAPON_PERCENT_DAMAGE {
		result *= 0.01
	}
	return result
}

func (effect *SpellEffect) IsDirectDamageEffect() bool {
	types := []SpellEffectType{
		E_HEAL, E_SCHOOL_DAMAGE, E_HEALTH_LEECH,
		E_NORMALIZED_WEAPON_DMG, E_WEAPON_DAMAGE, E_WEAPON_PERCENT_DAMAGE,
	}
	return slices.Contains(types, effect.EffectType)
}

func (effect *SpellEffect) IsPeriodicDamageEffect() bool {
	subtypes := []EffectAuraType{
		A_PERIODIC_DAMAGE, A_PERIODIC_LEECH, A_PERIODIC_HEAL,
	}
	if effect.EffectType == E_APPLY_AURA {
		return slices.Contains(subtypes, effect.EffectAura)
	}
	return false
}

func (data *SpellEffect) ClassFlag(index uint) uint32 {
	return uint32(data.EffectSpellClassMasks[index/32]) & (1 << (index % 32))
}

func (effect *SpellEffect) ParseStatEffect() *stats.Stats {
	stats := &stats.Stats{}
	scale := effect.ScalingClass()
	spell := dbcInstance.Spells[effect.SpellID]
	stat, _ := MapMainStatToStat(effect.EffectMiscValues[0])

	switch {
	case effect.EffectAura == A_MOD_STAT && effect.EffectType == E_APPLY_AURA:
		stats[stat] = math.Round(float64(effect.Coefficient * dbcInstance.SpellScaling(scale, spell.MaxScalingLevel)))
		if effect.Coefficient <= 0 {
			// if Coefficient is not set, we fall back to EffectBasePoints
			stats[stat] = float64(effect.EffectBasePoints)
		}

	case effect.EffectAura == A_MOD_DAMAGE_DONE && effect.EffectType == E_APPLY_AURA:
		// Apply spell power, A_MOD_HEALING_DONE is also a possibility for healing power
		stats[proto.Stat_StatSpellPower] = float64(effect.EffectBasePoints)

	case effect.EffectMiscValues[0] == -1 && effect.EffectAura == A_MOD_STAT && effect.EffectType == E_APPLY_AURA:
		// -1 represents ALL STATS if present in MiscValue 0
		for _, s := range []proto.Stat{
			proto.Stat_StatAgility, proto.Stat_StatIntellect, proto.Stat_StatSpirit,
			proto.Stat_StatStamina, proto.Stat_StatStrength,
		} {
			stats[s] = float64(effect.EffectBasePoints)
		}

	case effect.EffectAura == A_MOD_RESISTANCE:
		school := SpellSchool(effect.EffectMiscValues[0])
		for schoolType, stat := range SpellSchoolToStat {
			if school.Has(schoolType) {
				stats[stat] += float64(effect.EffectBasePoints)
			}
		}

	case effect.EffectAura == A_MOD_RATING:
		for _, rating := range getMatchingRatingMods(effect.EffectMiscValues[0]) {
			if statMod := RatingModToStat[rating]; statMod != -1 {
				stats[statMod] = float64(effect.EffectBasePoints)
			}
		}

	case effect.EffectAura == A_MOD_INCREASE_ENERGY:
		stats[proto.Stat_StatMana] = float64(effect.EffectBasePoints)

	case effect.EffectAura == A_PERIODIC_TRIGGER_SPELL && effect.EffectAuraPeriod == 10000:
		for _, sub := range dbcInstance.SpellEffects[effect.EffectTriggerSpell] {
			stats.AddInplace(sub.ParseStatEffect())
		}
	}

	return stats
}
