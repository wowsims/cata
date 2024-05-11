package dbc

import (
	"fmt"
	"time"
)

// Struct Definitions
type SpellLabelData struct {
	ID      uint
	SpellID uint
	Label   int16
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

func (s *SpellData) IsValid() bool {
	return s.ID != 0 // assuming an ID of 0 means uninitialized or invalid data
}

// Flags checks for specific attributes
func (s *SpellData) HasAttributeFlag(attr uint) bool {
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
	if sd == nil {
		return nil, fmt.Errorf("spell data is nil")
	}
	if idx <= 0 {
		return nil, fmt.Errorf("effect index must not be zero or less")
	}
	if idx > int(sd.EffectsCount) {
		return nil, fmt.Errorf("effect index out of bounds")
	}
	return sd.Effects[idx-1], nil
}

// Helper function to extract a specific flag from SpellData's ClassFlags
func (data *SpellData) classFlag(index uint) (uint32, error) {
	if index/32 >= uint(len(data.ClassFlags)) {
		return 0, fmt.Errorf("index out of range")
	}
	return uint32(data.ClassFlags[index/32]) & (1 << (index % 32)), nil
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
