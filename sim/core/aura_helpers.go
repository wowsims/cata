// Functions for creating common types of auras.
package core

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core/stats"
)

type AuraCallback uint16

func (c AuraCallback) Matches(other AuraCallback) bool {
	return (c & other) != 0
}

const (
	CallbackEmpty AuraCallback = 0

	CallbackOnSpellHitDealt AuraCallback = 1 << iota
	CallbackOnSpellHitTaken
	CallbackOnPeriodicDamageDealt
	CallbackOnHealDealt
	CallbackOnPeriodicHealDealt
	CallbackOnCastComplete
	CallbackOnApplyEffects
)

type ProcHandler func(sim *Simulation, spell *Spell, result *SpellResult)
type ProcExtraCondition func(sim *Simulation, spell *Spell, result *SpellResult) bool

type ProcTrigger struct {
	Name            string
	ActionID        ActionID
	Duration        time.Duration
	Callback        AuraCallback
	ProcMask        ProcMask
	ProcMaskExclude ProcMask
	SpellFlags      SpellFlag
	Outcome         HitOutcome
	Harmful         bool
	ProcChance      float64
	PPM             float64
	ICD             time.Duration
	Handler         ProcHandler
	ClassSpellMask  int64
	ExtraCondition  ProcExtraCondition
}

func ApplyProcTriggerCallback(unit *Unit, aura *Aura, config ProcTrigger) {
	var icd Cooldown
	if config.ICD != 0 {
		icd = Cooldown{
			Timer:    unit.NewTimer(),
			Duration: config.ICD,
		}
		aura.Icd = &icd
	}

	var ppmm PPMManager
	if config.PPM > 0 {
		ppmm = unit.AutoAttacks.NewPPMManager(config.PPM, config.ProcMask)
		aura.Ppmm = &ppmm
	}

	handler := config.Handler
	callback := func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
		if config.SpellFlags != SpellFlagNone && !spell.Flags.Matches(config.SpellFlags) {
			return
		}
		if config.ClassSpellMask > 0 && config.ClassSpellMask&spell.ClassSpellMask == 0 {
			return
		}
		if config.ProcMaskExclude != ProcMaskUnknown && spell.ProcMask.Matches(config.ProcMaskExclude) {
			return
		}
		if config.ProcMask != ProcMaskUnknown && !spell.ProcMask.Matches(config.ProcMask) {
			return
		}
		if config.Outcome != OutcomeEmpty && !result.Outcome.Matches(config.Outcome) {
			return
		}
		if config.Harmful && result.Damage == 0 {
			return
		}
		if icd.Duration != 0 && !icd.IsReady(sim) {
			return
		}
		if config.ExtraCondition != nil && !config.ExtraCondition(sim, spell, result) {
			return
		}
		if config.ProcChance != 1 && sim.RandomFloat(config.Name) > config.ProcChance {
			return
		} else if config.PPM != 0 && !ppmm.Proc(sim, spell.ProcMask, config.Name) {
			return
		}

		if icd.Duration != 0 {
			icd.Use(sim)
		}
		handler(sim, spell, result)
	}

	if config.ProcChance == 0 {
		config.ProcChance = 1
	}

	if config.Callback.Matches(CallbackOnSpellHitDealt) {
		aura.OnSpellHitDealt = callback
	}
	if config.Callback.Matches(CallbackOnSpellHitTaken) {
		aura.OnSpellHitTaken = callback
	}
	if config.Callback.Matches(CallbackOnPeriodicDamageDealt) {
		aura.OnPeriodicDamageDealt = callback
	}
	if config.Callback.Matches(CallbackOnHealDealt) {
		aura.OnHealDealt = callback
	}
	if config.Callback.Matches(CallbackOnPeriodicHealDealt) {
		aura.OnPeriodicHealDealt = callback
	}
	if config.Callback.Matches(CallbackOnCastComplete) {
		aura.OnCastComplete = func(aura *Aura, sim *Simulation, spell *Spell) {
			if config.SpellFlags != SpellFlagNone && !spell.Flags.Matches(config.SpellFlags) {
				return
			}
			if config.ClassSpellMask > 0 && config.ClassSpellMask&spell.ClassSpellMask == 0 {
				return
			}
			if config.ProcMask != ProcMaskUnknown && !spell.ProcMask.Matches(config.ProcMask) {
				return
			}
			if config.ProcMaskExclude != ProcMaskUnknown && spell.ProcMask.Matches(config.ProcMaskExclude) {
				return
			}
			if icd.Duration != 0 && !icd.IsReady(sim) {
				return
			}
			if config.ProcChance != 1 && sim.RandomFloat(config.Name) > config.ProcChance {
				return
			}

			if icd.Duration != 0 {
				icd.Use(sim)
			}
			handler(sim, spell, nil)
		}
	}
	if config.Callback.Matches(CallbackOnApplyEffects) {
		aura.OnApplyEffects = func(aura *Aura, sim *Simulation, target *Unit, spell *Spell) {
			if config.SpellFlags != SpellFlagNone && !spell.Flags.Matches(config.SpellFlags) {
				return
			}
			if config.ClassSpellMask > 0 && config.ClassSpellMask&spell.ClassSpellMask == 0 {
				return
			}
			if config.ProcMask != ProcMaskUnknown && !spell.ProcMask.Matches(config.ProcMask) {
				return
			}
			if config.ProcMaskExclude != ProcMaskUnknown && spell.ProcMask.Matches(config.ProcMaskExclude) {
				return
			}
			if icd.Duration != 0 && !icd.IsReady(sim) {
				return
			}
			if config.ProcChance != 1 && sim.RandomFloat(config.Name) > config.ProcChance {
				return
			}

			if icd.Duration != 0 {
				icd.Use(sim)
			}
			handler(sim, spell, &SpellResult{Target: target})
		}
	}
}

func MakeProcTriggerAura(unit *Unit, config ProcTrigger) *Aura {
	aura := Aura{
		Label:           config.Name,
		ActionIDForProc: config.ActionID,
		Duration:        config.Duration,
	}
	if config.Duration == 0 {
		aura.Duration = NeverExpires
		aura.OnReset = func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		}
	}

	ApplyProcTriggerCallback(unit, &aura, config)

	return unit.GetOrRegisterAura(aura)
}

type StackingStatAura struct {
	Aura          Aura
	BonusPerStack stats.Stats
}

func MakeStackingAura(character *Character, config StackingStatAura) *Aura {
	bonusPerStack := config.BonusPerStack
	config.Aura.OnStacksChange = func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
		character.AddStatsDynamic(sim, bonusPerStack.Multiply(float64(newStacks-oldStacks)))
	}
	return character.RegisterAura(config.Aura)
}

// Returns the same Aura for chaining.
func MakePermanent(aura *Aura) *Aura {
	aura.Duration = NeverExpires
	if aura.OnReset == nil {
		aura.OnReset = func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		}
	} else {
		oldOnReset := aura.OnReset
		aura.OnReset = func(aura *Aura, sim *Simulation) {
			oldOnReset(aura, sim)
			aura.Activate(sim)
		}
	}
	return aura
}

func (character *Character) NewTemporaryStatBuffWithStacks(auraLabel string, actionID ActionID, bonusPerStack stats.Stats, maxStacks int32, duration time.Duration) *Aura {
	return MakeStackingAura(character, StackingStatAura{
		Aura: Aura{
			Label:     auraLabel,
			ActionID:  actionID,
			Duration:  duration,
			MaxStacks: maxStacks,
		},
		BonusPerStack: bonusPerStack,
	})
}

// Helper for the common case of making an aura that adds stats.
func (character *Character) NewTemporaryStatsAura(auraLabel string, actionID ActionID, tempStats stats.Stats, duration time.Duration) *Aura {
	return character.NewTemporaryStatsAuraWrapped(auraLabel, actionID, tempStats, duration, nil)
}

// Alternative that allows modifying the Aura config.
func (character *Character) NewTemporaryStatsAuraWrapped(auraLabel string, actionID ActionID, buffs stats.Stats, duration time.Duration, modConfig func(*Aura)) *Aura {
	// If one of the stat bonuses is a health bonus, then set up healing metrics for the associated
	// heal, since all temporary max health bonuses also instantaneously heal the player.
	var healthMetrics *ResourceMetrics
	amountHealed := buffs[stats.Health]
	includesHealthBuff := amountHealed > 0

	if includesHealthBuff {
		healthMetrics = character.NewHealthMetrics(actionID)
		buffs[stats.Health] = 0
	}

	config := Aura{
		Label:    auraLabel,
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *Aura, sim *Simulation) {
			if sim.Log != nil {
				character.Log(sim, "Gained %s from %s.", buffs.FlatString(), actionID)
			}
			character.AddStatsDynamic(sim, buffs)

			if includesHealthBuff {
				character.UpdateMaxHealth(sim, amountHealed, healthMetrics)
			}

			for i := range character.OnTemporaryStatsChanges {
				character.OnTemporaryStatsChanges[i](sim, aura, buffs)
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if sim.Log != nil {
				character.Log(sim, "Lost %s from fading %s.", buffs.FlatString(), actionID)
			}
			character.AddStatsDynamic(sim, buffs.Invert())

			if includesHealthBuff {
				character.UpdateMaxHealth(sim, -amountHealed, healthMetrics)
			}

			for i := range character.OnTemporaryStatsChanges {
				character.OnTemporaryStatsChanges[i](sim, aura, buffs.Invert())
			}
		},
	}

	if modConfig != nil {
		modConfig(&config)
	}

	return character.GetOrRegisterAura(config)
}

type ShieldStrengthCalculator func(unit *Unit) float64

type DamageAbsorptionAura struct {
	*Aura
	ShieldStrength                float64
	FreshShieldStrengthCalculator ShieldStrengthCalculator
}

func (aura *DamageAbsorptionAura) Activate(sim *Simulation) {
	aura.Aura.Activate(sim)
	aura.ShieldStrength = aura.FreshShieldStrengthCalculator(aura.Unit)
	stacks := max(1, int32(aura.ShieldStrength))
	aura.Aura.MaxStacks = stacks
	aura.Aura.SetStacks(sim, stacks)
}

func (character *Character) NewDamageAbsorptionAura(auraLabel string, actionID ActionID, duration time.Duration, calculator ShieldStrengthCalculator) *DamageAbsorptionAura {
	return CreateDamageAbsorptionAura(character, auraLabel, actionID, duration, calculator, nil)
}

func (character *Character) NewDamageAbsorptionAuraForSchool(auraLabel string, actionID ActionID, duration time.Duration, school SpellSchool, calculator ShieldStrengthCalculator) *DamageAbsorptionAura {
	return CreateDamageAbsorptionAura(character, auraLabel, actionID, duration, calculator, func(spell *Spell) bool {
		return spell.SpellSchool.Matches(school)
	})
}

func CreateDamageAbsorptionAura(character *Character, auraLabel string, actionID ActionID, duration time.Duration, calculator ShieldStrengthCalculator, extraSpellCheck func(spell *Spell) bool) *DamageAbsorptionAura {
	aura := &DamageAbsorptionAura{
		Aura: character.RegisterAura(Aura{
			Label:    auraLabel,
			ActionID: actionID,
			Duration: duration,
		}),
		FreshShieldStrengthCalculator: calculator,
	}

	character.AddDynamicDamageTakenModifier(func(sim *Simulation, spell *Spell, result *SpellResult) {
		if aura.Aura.IsActive() && result.Damage > 0 && (extraSpellCheck == nil || extraSpellCheck(spell)) {
			absorbedDamage := min(aura.ShieldStrength, result.Damage)
			result.Damage -= absorbedDamage
			aura.ShieldStrength -= absorbedDamage

			if sim.Log != nil {
				character.Log(sim, "%s absorbed %.1f damage, new shield strength: %.1f", auraLabel, absorbedDamage, aura.ShieldStrength)
			}

			aura.Aura.SetStacks(sim, int32(aura.ShieldStrength))
		}
	})

	return aura
}

func ApplyFixedUptimeAura(aura *Aura, uptime float64, tickLength time.Duration, startTime time.Duration) {
	auraDuration := aura.Duration
	ticksPerAura := float64(auraDuration) / float64(tickLength)
	chancePerTick := TernaryFloat64(uptime == 1, 1, 1.0-math.Pow(1-uptime, 1/ticksPerAura))

	aura.Unit.RegisterResetEffect(func(sim *Simulation) {
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period: tickLength,
			OnAction: func(sim *Simulation) {
				if sim.RandomFloat("FixedAura") < chancePerTick {
					aura.Activate(sim)
					if aura.MaxStacks > 0 {
						aura.AddStack(sim)
					}
				}
			},
		})

		// Also try once at the start.
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period:   startTime,
			NumTicks: 1,
			OnAction: func(sim *Simulation) {
				if sim.RandomFloat("FixedAura") < uptime {
					// Use random duration to compensate for increased chance collapsed into single tick.
					randomDur := tickLength + time.Duration(float64(auraDuration-tickLength)*sim.RandomFloat("FixedAuraDur"))

					aura.Duration = randomDur
					aura.Activate(sim)
					aura.Duration = auraDuration
				}
			},
		})
	})
}
