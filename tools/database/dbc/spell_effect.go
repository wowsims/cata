package dbc

import (
	"math"
	"slices"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const MAX_SCALING_LEVEL = 100
const BASE_LEVEL = 85

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
	EffectMiscValues      []int // from EffectMiscValue_0, EffectMiscValue_1
	EffectRadiusIndices   []int // from EffectRadiusIndex_0, EffectRadiusIndex_1
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
	return math.Max(float64(s.EffectRadiusIndices[0]), float64(s.EffectRadiusIndices[1]))
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
	case 11:
		return proto.Class_ClassDruid
	default:
		return proto.Class_ClassUnknown
	}
}
func (s *SpellEffect) Delta(pLevel int, level int) float64 {
	if level > 85 {
		level = 85
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
	// Ensure the operation is performed within uint32 context
	return uint32(data.EffectSpellClassMasks[index/32]) & (1 << (index % 32))
}

func (effect *SpellEffect) ParseStatEffect() *stats.Stats {
	stats := &stats.Stats{}
	if effect.EffectAura == A_MOD_STAT && effect.EffectType == E_APPLY_AURA {
		if effect.EffectMiscValues[0] > 0 {
			stat, _ := MapMainStatToStat(effect.EffectMiscValues[0])

			stats[stat] = float64(effect.EffectBasePoints)
			return stats
		}
	}

	if effect.EffectAura == A_MOD_DAMAGE_DONE && effect.EffectType == E_APPLY_AURA {
		stats[proto.Stat_StatSpellPower] = float64(effect.EffectBasePoints)
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
