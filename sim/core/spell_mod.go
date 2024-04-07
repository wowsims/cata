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
	IntValue   int64
	TimeValue  time.Duration
	FloatValue float64
}

type SpellMod struct {
	ClassMask      int64
	Kind           SpellModType
	School         SpellSchool
	floatValue     float64
	intValue       int64
	Stacks         int32
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
		floatValue: config.FloatValue,
		intValue:   config.IntValue,
		timeValue:  config.TimeValue,
		Apply:      functions.Apply,
		Remove:     functions.Remove,
		IsActive:   false,
		Stacks:     1,
	}

	for _, spell := range unit.Spellbook {
		if shouldApply(spell, mod) {
			mod.AffectedSpells = append(mod.AffectedSpells, spell)
		}
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
	unit.ActiveSpellMods = append(unit.ActiveSpellMods, mod)
	mod.Activate()
}

// Never use dynamic mods for Auras that have ExpireNever and activate on reset
// Those mods will be overwritten potentilly during sim reset
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

func (mod *SpellMod) AddStack() {
	if mod.IsActive {
		mod.Deactivate()
		mod.Stacks += 1
		mod.Activate()
	} else {
		mod.Stacks += 1
		mod.Activate()
	}
}

func (mod *SpellMod) RemoveStack() {
	if mod.IsActive {
		mod.Deactivate()
		mod.Stacks = max(mod.Stacks-1, 0)
		if mod.Stacks > 0 {
			mod.Activate()
		}
	} else {
		mod.Stacks = max(mod.Stacks-1, 0)
	}
}

// Mod implmentations
type SpellModType uint32

const (
	SpellMod_DamageDonePercent SpellModType = 1 << iota
	SpellMod_DamageDoneAdd
	SpellMod_PowerCostPercent
	SpellMod_CooldownFlat
	SpellMod_CritMultiplier
	SpellMod_CastTimePercent
)

var spellModMap = map[SpellModType]*SpellModFunctions{
	SpellMod_DamageDonePercent: {
		Apply:  applyDamageDonePercent,
		Remove: removeDamageDonePercent,
	},

	SpellMod_DamageDoneAdd: {
		Apply:  applyDamageDoneAdd,
		Remove: removeDamageDonAdd,
	},

	SpellMod_PowerCostPercent: {
		Apply:  applyPowerCostPercent,
		Remove: removePowerCostPercent,
	},

	SpellMod_CooldownFlat: {
		Apply:  applyCooldownFlat,
		Remove: removeCooldownFlat,
	},

	SpellMod_CritMultiplier: {
		Apply:  applyCritMultiplier,
		Remove: removeCritMultiplier,
	},

	SpellMod_CastTimePercent: {
		Apply:  applyCastTimePercent,
		Remove: removeCastTimePercent,
	},
}

func applyDamageDonePercent(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplier *= 1 + mod.floatValue*float64(mod.Stacks)
}

func removeDamageDonePercent(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplier /= 1 + mod.floatValue*float64(mod.Stacks)
}

func applyDamageDoneAdd(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplierAdditive += mod.floatValue * float64(mod.Stacks)
}

func removeDamageDonAdd(mod *SpellMod, spell *Spell) {
	spell.DamageMultiplierAdditive -= mod.floatValue * float64(mod.Stacks)
}

func applyPowerCostPercent(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.Cost *= 1 + mod.floatValue*float64(mod.Stacks)
}

func removePowerCostPercent(mod *SpellMod, spell *Spell) {
	spell.DefaultCast.Cost /= 1 + mod.floatValue*float64(mod.Stacks)
}

func applyCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.CD.Duration += mod.timeValue * time.Duration(mod.Stacks)
}

func removeCooldownFlat(mod *SpellMod, spell *Spell) {
	spell.CD.Duration -= mod.timeValue * time.Duration(mod.Stacks)
}

func applyCritMultiplier(mod *SpellMod, spell *Spell) {
	spell.CritMultiplier = 1 + (spell.CritMultiplier-1)*(mod.floatValue*float64(mod.Stacks)+1)
}

func removeCritMultiplier(mod *SpellMod, spell *Spell) {
	spell.CritMultiplier = 1 + (spell.CritMultiplier-1)/(mod.floatValue*float64(mod.Stacks)+1)
}

func applyCastTimePercent(mod *SpellMod, spell *Spell) {
	spell.CastTimeMultiplier += mod.floatValue * float64(mod.Stacks)
}

func removeCastTimePercent(mod *SpellMod, spell *Spell) {
	spell.CastTimeMultiplier -= mod.floatValue * float64(mod.Stacks)
}
