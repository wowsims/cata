package dbc

import "github.com/wowsims/cata/sim/core/proto"

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
