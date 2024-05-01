package dbc

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
)

type Character interface {
	GetLevel() int
	// other methods needed by dbc
}

// Enums and constants
const (
	NUM_SPELL_FLAGS        = 15
	NUM_CLASS_FAMILY_FLAGS = 4
)

// Struct Definitions
type SpellLabelData struct {
	ID      uint
	SpellID uint
	Label   int16
}

type SpellPowerData struct {
	ID             uint
	SpellID        uint
	AuraID         uint
	PowerType      int
	Cost           int
	CostMax        int
	CostPerTick    int
	PctCost        float64
	PctCostMax     float64
	PctCostPerTick float64
}

func (s *SpellEffectData) GetRadiusMax() float64 {
	return math.Max(s.RadiusMax, s.Radius)
}

type SpellData struct {
	Name                 string
	ID                   uint
	School               uint
	PrjSpeed             float64
	PrjDelay             float64
	PrjMinDuration       float64
	RaceMask             int
	ClassMask            int
	MaxScalingLevel      int
	SpellLevel           int
	MaxLevel             int
	ReqMaxLevel          int
	MinRange             float64
	MaxRange             float64
	Cooldown             time.Duration
	GCD                  time.Duration
	CategoryCooldown     time.Duration
	Charges              uint
	ChargeCooldown       time.Duration
	Category             uint
	DmgClass             uint
	MaxTargets           int
	Duration             time.Duration
	MaxStack             uint
	ProcChance           uint
	ProcCharges          int
	ProcFlags            uint64
	InternalCooldown     time.Duration
	RPPM                 float64
	EquippedClass        uint
	EquippedInvtypeMask  uint
	EquippedSubclassMask uint
	CastTime             time.Duration
	Attributes           [NUM_SPELL_FLAGS]uint
	ClassFlags           [NUM_CLASS_FAMILY_FLAGS]uint
	ClassFlagsFamily     uint
	StanceMask           uint
	Mechanic             uint
	PowerID              uint
	EssenceID            uint
	Effects              []*SpellEffectData
	Power                []*SpellPowerData
	Driver               []*SpellData
	Labels               []*SpellLabelData
	EffectsCount         uint8
	PowerCount           uint8
	DriverCount          uint8
	LabelsCount          uint8
}

func (s *SpellData) Ok() bool {
	return s.ID != 0 // assuming an ID of 0 means uninitialized or invalid data
}

// Flags checks for specific attributes
func (s *SpellData) Flags(attr uint) bool {
	bit := attr % 32
	index := attr / 32
	if index >= uint(len(s.Attributes)) {
		return false
	}
	return (s.Attributes[index] & (1 << bit)) != 0
}
func (s *SpellData) HasDirectDamageEffect() bool {
	for _, effect := range s.Effects {
		if effect.IsDirectDamageEffect() {
			return true
		}
	}
	return false
}

// Determines if any effect in the spell is a periodic damage effect
func (s *SpellData) HasPeriodicDamageEffect() bool {
	for _, effect := range s.Effects {
		if effect.IsPeriodicDamageEffect() {
			return true
		}
	}
	return false
}
func (s *SpellData) CooldownMillis() time.Duration {
	return s.Cooldown
}

func (sd *SpellData) EffectN(idx int) (*SpellEffectData, error) {
	if idx <= 0 {
		return nil, fmt.Errorf("effect index must not be zero or less")
	}

	if sd == nil {
		return nil, fmt.Errorf("spell data is nil or not found")
	}

	if idx > int(sd.EffectsCount) {
		return nil, fmt.Errorf("effect index out of bound")
	}

	return sd.Effects[idx-1], nil
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

// Minimum value calculation for item
// func (s *SpellEffectData) MinItem(item *core.Item) float64 {
// 	return s.scaledMin(s.AverageItem(item), s.DeltaItem(item))
// }

// Maximum value calculation for player
func (s *SpellEffectData) Max(dbc *DBC, pLevel int, level int) float64 {
	return s.scaledMax(s.Average(dbc, pLevel, level), s.Delta(dbc, pLevel, level))
}

// Maximum value calculation for item
//
//	func (s *SpellEffectData) MaxItem(item *Item) float64 {
//		return s.scaledMax(s.AverageItem(item), s.DeltaItem(item))
//	}
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

// Average calculation for an item
// func (s *SpellEffectData) AverageItem(item *core.Item) float64 {
// 	if item == nil {
// 		return 0
// 	}

// 	budget := item.Budget()
// 	if s.ScalingClass() == PLAYER_SPECIAL_SCALE7 {
// 		budget = itemDatabase.ApplyCombatRatingMultiplier(*item, budget)
// 	} else if s.ScalingClass() == PLAYER_SPECIAL_SCALE8 {
// 		props := item.Player.RandomProperty(item.ItemLevel())
// 		budget = props.DamageReplaceStat
// 	} else if (s.ScalingClass() == PLAYER_NONE || s.ScalingClass() == PLAYER_SPECIAL_SCALE9) && s.Spell.Flags(SX_SCALE_ILEVEL) {
// 		props := item.Player.RandomProperty(item.ItemLevel())
// 		budget = props.DamageSecondary
// 	}

// 	return s.MCoeff * budget
// }

// Scaled delta calculation
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

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *SpellData) AffectedByAll(effect *SpellEffectData) bool {
	return s.AffectedByEffect(effect) || s.AffectedByCategory(effect)
}

func (s *SpellData) AffectedByCategory(effect *SpellEffectData) bool {
	return s.AffectedByCategoryValue(effect.MiscValue)
}

func (s *SpellData) AffectedByCategoryValue(category int) bool {
	return category > 0 && s.Category == uint(category)
}

func (s *SpellData) AffectedBy(effect *SpellData) bool {
	if s.ClassFlagsFamily != effect.ClassFlagsFamily {
		return false
	}

	for flagIdx := 0; flagIdx < NUM_CLASS_FAMILY_FLAGS; flagIdx++ {
		if s.ClassFlags[flagIdx] == 0 {
			continue
		}

		for _, e := range effect.Effects {
			if e.ClassFlags[flagIdx]&s.ClassFlags[flagIdx] != 0 {
				return true
			}
		}
	}
	return false
}

func (s *SpellData) AffectedByEffect(effect *SpellEffectData) bool {
	if s.ClassFlagsFamily != effect.Spell.ClassFlagsFamily {
		return false
	}

	for flagIdx := 0; flagIdx < NUM_CLASS_FAMILY_FLAGS; flagIdx++ {
		if s.ClassFlags[flagIdx] == 0 {
			continue
		}

		if effect.ClassFlags[flagIdx]&s.ClassFlags[flagIdx] != 0 {
			return true
		}
	}
	return false
}
