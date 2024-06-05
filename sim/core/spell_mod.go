package core

import (
	"strconv"
	"time"
)

/*
SpellMod implementation.
*/

type SpellModConfig struct {
	ClassMask  int64
	Kind       SpellModType
	School     SpellSchool
	ProcMask   ProcMask
	IntValue   int64
	TimeValue  time.Duration
	FloatValue float64
}

type SpellMod struct {
	ClassMask      int64
	Kind           SpellModType
	School         SpellSchool
	ProcMask       ProcMask
	floatValue     float64
	intValue       int64
	timeValue      time.Duration
	Apply          SpellModApply
	Remove         SpellModRemove
	IsActive       bool
	AffectedSpells []*Spell
}

type SpellModApply func(mod *SpellMod, spell *Spell)
type SpellModRemove func(mod *SpellMod, spell *Spell)
type SpellModFunctions struct {
	Apply  SpellModApply
	Remove SpellModRemove
}

func buildMod(unit *Unit, config SpellModConfig) *SpellMod {
	functions := spellModMap[config.Kind]
	if functions == nil {
		panic("SpellMod " + strconv.Itoa(int(config.Kind)) + " not implmented")
	}

	mod := &SpellMod{
		ClassMask:  config.ClassMask,
		Kind:       config.Kind,
		School:     config.School,
		ProcMask:   config.ProcMask,
		floatValue: config.FloatValue,
		intValue:   config.IntValue,
		timeValue:  config.TimeValue,
		Apply:      functions.Apply,
		Remove:     functions.Remove,
		IsActive:   false,
	}

	unit.OnSpellRegistered(func(spell *Spell) {
		if shouldApply(spell, mod) {
			mod.AffectedSpells = append(mod.AffectedSpells, spell)

			if mod.IsActive {
				mod.Apply(mod, spell)
			}
		}
	})

	return mod
}

func (unit *Unit) AddStaticMod(config SpellModConfig) {
	mod := buildMod(unit, config)
	mod.Activate()
}

func (unit *Unit) AddDynamicMod(config SpellModConfig) *SpellMod {
	return buildMod(unit, config)
}

func shouldApply(spell *Spell, mod *SpellMod) bool {
	if spell.Flags.Matches(SpellFlagNoSpellMods) {
		return false
	}

	if mod.ClassMask > 0 && mod.ClassMask&spell.ClassSpellMask == 0 {
		return false
	}

	if mod.School > 0 && !mod.School.Matches(spell.SpellSchool) {
		return false
	}

	if mod.ProcMask > 0 && !mod.ProcMask.Matches(spell.ProcMask) {
		return false
	}

	return true
}

func (mod *SpellMod) UpdateIntValue(value int64) {
	if mod.IsActive {
		mod.Deactivate()
		mod.intValue = value
		mod.Activate()
	} else {
		mod.intValue = value
	}
}

func (mod *SpellMod) UpdateTimeValue(value time.Duration) {
	if mod.IsActive {
		mod.Deactivate()
		mod.timeValue = value
		mod.Activate()
	} else {
		mod.timeValue = value
	}
}

func (mod *SpellMod) UpdateFloatValue(value float64) {
	if mod.IsActive {
		mod.Deactivate()
		mod.floatValue = value
		mod.Activate()
	} else {
		mod.floatValue = value
	}
}

func (mod *SpellMod) GetIntValue() int64 {
	return mod.intValue
}

func (mod *SpellMod) GetFloatValue() float64 {
	return mod.floatValue
}

func (mod *SpellMod) GetTimeValue() time.Duration {
	return mod.timeValue
}

func (mod *SpellMod) Activate() {
	if mod.IsActive {
		return
	}

	for _, spell := range mod.AffectedSpells {
		mod.Apply(mod, spell)
	}

	mod.IsActive = true
}

func (mod *SpellMod) Deactivate() {
	if !mod.IsActive {
		return
	}

	for _, spell := range mod.AffectedSpells {
		mod.Remove(mod, spell)
	}

	mod.IsActive = false
}

// Mod implmentations
type SpellModType uint32

const (
	// Will multiply the spell.DamageDoneMultiplier. +5% = 0.05
	// Uses FloatValue
	SpellMod_DamageDone_Pct SpellModType = 1 << iota

	// Will add the value spell.DamageDoneAddMultiplier
	// Uses FloatValue
	SpellMod_DamageDone_Flat

	// Will reduce spell.DefaultCast.Cost by % amount. -5% = -0.05
	// Uses FloatValue
	SpellMod_PowerCost_Pct

	// Increases or decreases spell.DefaultCast.Cost by flat amount
	// Uses FloatValue
	SpellMod_PowerCost_Flat

	// Increases or decreases RuneCost.RunicPowerCost by flat amount
	// Uses FloatValue
	SpellMod_RunicPowerCost_Flat

	// Will add time.Duration to spell.CD.Duration
	// Uses TimeValue
	SpellMod_Cooldown_Flat

	// Increases or decreases spell.CD.Multiplier by flat amount
	// Uses FloatValue
	SpellMod_Cooldown_Multiplier

	// Will increase the CritMultiplier. +100% = 1.0
	// Uses FloatValue
	SpellMod_CritMultiplier_Pct

	// Will add / substract % amount from the cast time multiplier.
	// Ueses: FloatValue
	SpellMod_CastTime_Pct

	// Will add / substract time from the cast time.
	// Ueses: TimeValue
	SpellMod_CastTime_Flat

	// Add/subtract bonus crit rating
	// Uses: FloatValue
	SpellMod_BonusCrit_Rating

	// Add/subtract bonus hit rating
	// Uses: FloatValue
	SpellMod_BonusHit_Rating

	// Add/subtract to the dots max ticks
	// Uses: IntValue
	SpellMod_DotNumberOfTicks_Flat

	// Add/subtract to the casts gcd
	// Uses: TimeValue
	SpellMod_GlobalCooldown_Flat

	// Add/substrct to the base tick frequency
	// Uses: TimeValue
	SpellMod_DotTickLength_Flat

	// Add/subtract bonus coefficient
	// Uses: FloatValue
	SpellMod_BonusCoeffecient_Flat

	// Enables casting while moving
	SpellMod_AllowCastWhileMoving

	// Add/subtract bonus spell power
	// Uses: FloatValue
	SpellMod_BonusSpellPower_Flat
)

var spellModMap = map[SpellModType]*SpellModFunctions{
	SpellMod_DamageDone_Pct: {
		Apply:  applyDamageDonePercent,
		Remove: removeDamageDonePercent,
	},

	SpellMod_DamageDone_Flat: {
		Apply:  applyDamageDoneAdd,
		Remove: removeDamageDonAdd,
	},

	SpellMod_PowerCost_Pct: {
		Apply:  applyPowerCostPercent,
		Remove: removePowerCostPercent,
	},

	SpellMod_PowerCost_Flat: {
		Apply:  applyPowerCostFlat,
		Remove: removePowerCostFlat,
	},

	SpellMod_RunicPowerCost_Flat: {
		Apply:  applyRunicPowerCostFlat,
		Remove: removeRunicPowerCostFlat,
	},

	SpellMod_Cooldown_Flat: {
		Apply:  applyCooldownFlat,
		Remove: removeCooldownFlat,
	},

	SpellMod_Cooldown_Multiplier: {
		Apply:  applyCooldownMultiplier,
		Remove: removeCooldownMultiplier,
	},

	SpellMod_CritMultiplier_Pct: {
		Apply:  applyCritMultiplier,
		Remove: removeCritMultiplier,
	},

	SpellMod_CastTime_Pct: {
		Apply:  applyCastTimePercent,
		Remove: removeCastTimePercent,
	},

	SpellMod_CastTime_Flat: {
		Apply:  applyCastTimeFlat,
		Remove: removeCastTimeFlat,
	},

	SpellMod_BonusCrit_Rating: {
		Apply:  applyBonusCritRating,
		Remove: removeBonusCritRating,
	},

	SpellMod_BonusHit_Rating: {
		Apply:  applyBonusHitRating,
		Remove: removeBonusHitRating,
	},

	SpellMod_DotNumberOfTicks_Flat: {
		Apply:  applyDotNumberOfTicks,
		Remove: removeDotNumberOfTicks,
	},

	SpellMod_GlobalCooldown_Flat: {
		Apply:  applyGlobalCooldownFlat,
		Remove: removeGlobalCooldownFlat,
	},
	SpellMod_DotTickLength_Flat: {
		Apply:  applyDotTickLengthFlat,
		Remove: removeDotTickLengthFlat,
	},

	SpellMod_BonusCoeffecient_Flat: {
		Apply:  applyBonusCoefficientFlat,
		Remove: removeBonusCoefficientFlat,
	},

	SpellMod_AllowCastWhileMoving: {
		Apply:  applyAllowCastWhileMoving,
		Remove: removeAllowCastWhileMoving,
	},

	SpellMod_BonusSpellPower_Flat: {
		Apply:  applyBonusSpellPowerFlat,
		Remove: removeBonusSpellPowerFlat,
	},
}

func applyDamageDonePercent(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplier *= 1 + mod.floatValue
}

func removeDamageDonePercent(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplier /= 1 + mod.floatValue
}

func applyDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplierAdditive += mod.floatValue
}

func removeDamageDonAdd(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplierAdditive -= mod.floatValue
}

func applyPowerCostPercent(mod *SpellMod, spell *Spell) {
	spell.CostMultiplier += mod.floatValue
}

func removePowerCostPercent(mod *SpellMod, spell *Spell) {
	spell.CostMultiplier -= mod.floatValue
}

func applyPowerCostFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.Cost += mod.floatValue
}

func removePowerCostFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.Cost -= mod.floatValue
}

func applyRunicPowerCostFlat(mod *SpellMod, spell *Spell) {
	cost := spell.RuneCostImpl()
	cost.RunicPowerCost += mod.floatValue
	spell.Cost = newRuneCost(spell, cost.GetConfig())
}

func removeRunicPowerCostFlat(mod *SpellMod, spell *Spell) {
	cost := spell.RuneCostImpl()
	cost.RunicPowerCost -= mod.floatValue
	spell.Cost = newRuneCost(spell, cost.GetConfig())
}

func applyCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.CD.Duration += mod.timeValue
}

func removeCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.CD.Duration -= mod.timeValue
}

func applyCooldownMultiplier(mod *SpellMod, spell *Spell) {
	spell.CdMultiplier += mod.floatValue
}

func removeCooldownMultiplier(mod *SpellMod, spell *Spell) {
	spell.CdMultiplier -= mod.floatValue
}

func applyCritMultiplier(mod *SpellMod, spell *Spell) {
	spell.CritMultiplier = 1 + (spell.CritMultiplier-1)*(mod.floatValue+1)
}

func removeCritMultiplier(mod *SpellMod, spell *Spell) {
	spell.CritMultiplier = 1 + (spell.CritMultiplier-1)/(mod.floatValue+1)
}

func applyCastTimePercent(mod *SpellMod, spell *Spell) {
	spell.CastTimeMultiplier += mod.floatValue
}

func removeCastTimePercent(mod *SpellMod, spell *Spell) {
	spell.CastTimeMultiplier -= mod.floatValue
}

func applyCastTimeFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.CastTime += mod.timeValue
}

func removeCastTimeFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.CastTime -= mod.timeValue
}

func applyBonusCritRating(mod *SpellMod, spell *Spell) {
	spell.BonusCritRating += mod.floatValue
}

func removeBonusCritRating(mod *SpellMod, spell *Spell) {
	spell.BonusCritRating -= mod.floatValue
}

func applyBonusHitRating(mod *SpellMod, spell *Spell) {
	spell.BonusHitRating += mod.floatValue
}

func removeBonusHitRating(mod *SpellMod, spell *Spell) {
	spell.BonusHitRating -= mod.floatValue
}

func applyDotNumberOfTicks(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.AddTicks(int32(mod.intValue))
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.AddTicks(int32(mod.intValue))
	}
}

func removeDotNumberOfTicks(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.AddTicks(-int32(mod.intValue))
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.AddTicks(-int32(mod.intValue))
	}
}

func applyGlobalCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.GCD += mod.timeValue
}

func removeGlobalCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.GCD -= mod.timeValue
}

func applyDotTickLengthFlat(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.TickLength += mod.timeValue
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.TickLength += mod.timeValue
	}
}

func removeDotTickLengthFlat(mod *SpellMod, spell *Spell) {
	if spell.dots != nil {
		for _, dot := range spell.dots {
			if dot != nil {
				dot.TickLength -= mod.timeValue
			}
		}
	}
	if spell.aoeDot != nil {
		spell.aoeDot.TickLength -= mod.timeValue
	}
}

func applyBonusCoefficientFlat(mod *SpellMod, spell *Spell) {
	spell.BonusCoefficient += mod.floatValue
}

func removeBonusCoefficientFlat(mod *SpellMod, spell *Spell) {
	spell.BonusCoefficient -= mod.floatValue
}

func applyAllowCastWhileMoving(mod *SpellMod, spell *Spell) {
	spell.Flags |= SpellFlagCanCastWhileMoving
}

func removeAllowCastWhileMoving(mod *SpellMod, spell *Spell) {
	spell.Flags ^= SpellFlagCanCastWhileMoving
}

func applyBonusSpellPowerFlat(mod *SpellMod, spell *Spell) {
	spell.BonusSpellPower += mod.floatValue
}

func removeBonusSpellPowerFlat(mod *SpellMod, spell *Spell) {
	spell.BonusSpellPower -= mod.floatValue
}
