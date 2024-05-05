package dbc

import (
	"math"

	"github.com/wowsims/cata/sim/core/proto"
)

// Constants
const MAX_SCALING_LEVEL = 100

type SpellEffectData struct {
	ID             uint
	SpellID        uint
	Index          uint
	Type           uint
	Subtype        uint
	ScalingType    int
	MCoeff         float64
	MDelta         float64
	MUnk           float64
	SPCoeff        float64
	APCoeff        float64
	Amplitude      float64
	Radius         float64
	RadiusMax      float64
	BaseValue      float64
	MiscValue      int
	MiscValue2     int
	ClassFlags     [NUM_CLASS_FAMILY_FLAGS]uint
	TriggerSpellID uint
	MChain         float64
	PPComboPoints  float64
	RealPPL        float64
	Mechanic       uint
	ChainTarget    int
	Targeting1     uint
	Targeting2     uint
	MValue         float64
	PVPCoeff       float64
	Spell          *SpellData
	TriggerSpell   *SpellData
}

// Get our proto class from scalingclass property
func (s *SpellEffectData) ScalingClass() proto.Class {
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

func (s *SpellEffectData) GetRadiusMax() float64 {
	return math.Max(s.RadiusMax, s.Radius)
}

func (effect *SpellEffectData) IsDirectDamageEffect() bool {
	types := []EffectType{
		E_HEAL, E_SCHOOL_DAMAGE, E_HEALTH_LEECH,
		E_NORMALIZED_WEAPON_DMG, E_WEAPON_DAMAGE, E_WEAPON_PERCENT_DAMAGE,
	}
	for _, t := range types {
		if effect.Type == uint(t) {
			return true
		}
	}
	return false
}

// Checks if the effect is a periodic damage type
func (effect *SpellEffectData) IsPeriodicDamageEffect() bool {
	subtypes := []EffectSubType{
		A_PERIODIC_DAMAGE, A_PERIODIC_LEECH, A_PERIODIC_HEAL, A_PERIODIC_HEAL_PCT,
	}
	if effect.Type == E_APPLY_AURA {
		for _, st := range subtypes {
			if effect.Subtype == uint(st) {
				return true
			}
		}
	}
	return false
}

// Method to calculate delta for a player
func (s *SpellEffectData) Delta(dbc *DBC, pLevel int, level int) float64 {
	if level > 85 {
		level = 85
	}

	var mScale float64
	spell := s.GetSpell(dbc)
	if s.MDelta != 0 && s.ScalingClass() != 0 {
		scalingLevel := level
		if scalingLevel == 0 {
			scalingLevel = pLevel
		}
		if spell.MaxScalingLevel > 0 {
			scalingLevel = min(scalingLevel, spell.MaxScalingLevel)
		}
		mScale = dbc.SpellScaling(s.ScalingClass(), scalingLevel)
	}

	return s.scaledDelta(mScale)
}

// Bonus method for a player
func (s *SpellEffectData) Bonus(dbc *DBC, pLevel int, level int) float64 {
	if level == 0 {
		level = pLevel
	}
	return dbc.EffectBonusById(s.GetSpell(dbc).ID, level)
}

// Minimum value calculation for player
func (s *SpellEffectData) Min(dbc *DBC, pLevel int, level int) float64 {
	return s.scaledMin(s.Average(dbc, pLevel, level), s.Delta(dbc, pLevel, level))
}

// Maximum value calculation for player
func (s *SpellEffectData) Max(dbc *DBC, pLevel int, level int) float64 {
	return s.scaledMax(s.Average(dbc, pLevel, level), s.Delta(dbc, pLevel, level))
}

func (s *SpellEffectData) GetSpell(dbc *DBC) *SpellData {
	return dbc.spellIndex[s.SpellID]
}

// Average calculation for a player
func (s *SpellEffectData) Average(dbc *DBC, pLevel int, level int) float64 {
	if level == 0 {
		level = pLevel
	}

	scale := s.ScalingClass()
	//Todo: DF stuff_
	// if scale == proto.Class_ClassUnknown && s.Spell.MaxScalingLevel > 0 {
	// 	scale = PLAYER_SPECIAL_SCALE8
	// }
	spell := s.GetSpell(dbc)

	if s.MCoeff != 0 && scale != proto.Class_ClassUnknown {
		if spell.MaxScalingLevel > 0 {
			level = min(level, spell.MaxScalingLevel)
		}
		scaler := dbc.SpellScaling(scale, level)
		value := s.MCoeff * scaler
		//todo: df stuff?
		// if scale == PLAYER_SPECIAL_SCALE7 {
		// 	value = itemDatabase.ApplyCombatRatingMultiplier(p, CR_MULTIPLIER_ARMOR, 1, value)
		// }
		return value
	} else if s.RealPPL != 0 {
		if spell.MaxLevel > 0 {
			return s.BaseValue + float64(min(level, spell.MaxLevel)-spell.SpellLevel)*s.RealPPL
		}
		return s.BaseValue + float64(level-spell.SpellLevel)*s.RealPPL
	}
	return s.BaseValue
}

// Helper function to check class flag in an effect or spell
func (data *SpellEffectData) classFlag(index uint) uint32 {
	// Ensure the operation is performed within uint32 context
	return uint32(data.ClassFlags[index/32]) & (1 << (index % 32))
}

func (s *SpellEffectData) scaledDelta(budget float64) float64 {
	if s.MDelta != 0 && budget > 0 {
		return s.MCoeff * s.MDelta * budget
	}
	return 0
}

// Scaled minimum calculation
func (s *SpellEffectData) scaledMin(avg, delta float64) float64 {
	result := avg - delta/2
	if s.Type == E_WEAPON_PERCENT_DAMAGE {
		result *= 0.01
	}
	return result
}

// Scaled maximum calculation
func (s *SpellEffectData) scaledMax(avg, delta float64) float64 {
	result := avg + delta/2
	if s.Type == E_WEAPON_PERCENT_DAMAGE {
		result *= 0.01
	}
	return result
}
