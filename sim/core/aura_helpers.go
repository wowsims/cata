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
	Name              string
	ActionID          ActionID
	Duration          time.Duration
	Callback          AuraCallback
	ProcMask          ProcMask
	ProcMaskExclude   ProcMask
	SpellFlags        SpellFlag
	SpellFlagsExclude SpellFlag
	Outcome           HitOutcome
	Harmful           bool
	ProcChance        float64
	PPM               float64
	DPM               *DynamicProcManager
	ICD               time.Duration
	Handler           ProcHandler
	ClassSpellMask    int64
	ExtraCondition    ProcExtraCondition
}

func ApplyProcTriggerCallback(unit *Unit, procAura *Aura, config ProcTrigger) {
	var icd Cooldown
	if config.ICD != 0 {
		icd = Cooldown{
			Timer:    unit.NewTimer(),
			Duration: config.ICD,
		}
		procAura.Icd = &icd
	}

	var dpm *DynamicProcManager
	if config.DPM != nil {
		dpm = config.DPM
	} else if config.PPM > 0 {
		dpm = unit.AutoAttacks.NewPPMManager(config.PPM, config.ProcMask)
	}

	if dpm != nil {
		procAura.Dpm = dpm
	}

	handler := config.Handler
	callback := func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
		if config.SpellFlags != SpellFlagNone && !spell.Flags.Matches(config.SpellFlags) {
			return
		}
		if config.SpellFlagsExclude != SpellFlagNone && spell.Flags.Matches(config.SpellFlagsExclude) {
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
		} else if dpm != nil && !dpm.Proc(sim, spell.ProcMask, config.Name) {
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
		procAura.OnSpellHitDealt = callback
	}
	if config.Callback.Matches(CallbackOnSpellHitTaken) {
		procAura.OnSpellHitTaken = callback
	}
	if config.Callback.Matches(CallbackOnPeriodicDamageDealt) {
		procAura.OnPeriodicDamageDealt = callback
	}
	if config.Callback.Matches(CallbackOnHealDealt) {
		procAura.OnHealDealt = callback
	}
	if config.Callback.Matches(CallbackOnPeriodicHealDealt) {
		procAura.OnPeriodicHealDealt = callback
	}
	if config.Callback.Matches(CallbackOnCastComplete) {
		procAura.OnCastComplete = func(aura *Aura, sim *Simulation, spell *Spell) {
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
		procAura.OnApplyEffects = func(aura *Aura, sim *Simulation, target *Unit, spell *Spell) {
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

type CustomStatBuffProcCondition func(sim *Simulation, aura *Aura) bool

// Analog to an Aura "sub-class" that additionally links the Aura to one or more
// Stats. Used within APL snapshotting wrappers.
type StatBuffAura struct {
	*Aura

	// All Stat types that are buffed (before dependencies) when this Aura
	// is activated.
	BuffedStatTypes []stats.Stat

	// Any special conditions (beyond standard ICD checks etc.) that must be
	// satisfied before this Aura can be activated.
	CustomProcCondition CustomStatBuffProcCondition

	// Whether the aura is currently swapped (in another item set) out or not.
	IsSwapped bool
}

func (aura *StatBuffAura) BuffsMatchingStat(statTypesToMatch []stats.Stat) bool {
	if aura == nil {
		return false
	}

	return CheckSliceOverlap(aura.BuffedStatTypes, statTypesToMatch)
}

func (aura *StatBuffAura) CanProc(sim *Simulation) bool {
	return !aura.IsSwapped && ((aura.CustomProcCondition == nil) || aura.CustomProcCondition(sim, aura.Aura))
}

func (aura *StatBuffAura) InferCDType() CooldownType {
	cdType := CooldownTypeUnknown

	if aura.BuffsMatchingStat([]stats.Stat{stats.Armor, stats.BlockPercent, stats.DodgeRating, stats.ParryRating, stats.Health, stats.ArcaneResistance, stats.FireResistance, stats.FrostResistance, stats.NatureResistance, stats.ShadowResistance}) {
		cdType |= CooldownTypeSurvival
	} else {
		cdType |= CooldownTypeDPS
	}

	return cdType
}

type StackingStatAura struct {
	Aura          Aura
	BonusPerStack stats.Stats
}

func MakeStackingAura(character *Character, config StackingStatAura) *StatBuffAura {
	bonusPerStack := config.BonusPerStack
	config.Aura.OnStacksChange = func(aura *Aura, sim *Simulation, oldStacks int32, newStacks int32) {
		character.AddStatsDynamic(sim, bonusPerStack.Multiply(float64(newStacks-oldStacks)))
	}
	return &StatBuffAura{
		Aura:            character.GetOrRegisterAura(config.Aura),
		BuffedStatTypes: bonusPerStack.GetBuffedStatTypes(),
	}
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

func (character *Character) NewTemporaryStatBuffWithStacks(auraLabel string, actionID ActionID, bonusPerStack stats.Stats, maxStacks int32, duration time.Duration) *StatBuffAura {
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
func (character *Character) NewTemporaryStatsAura(auraLabel string, actionID ActionID, tempStats stats.Stats, duration time.Duration) *StatBuffAura {
	return character.NewTemporaryStatsAuraWrapped(auraLabel, actionID, tempStats, duration, nil)
}

// Alternative that allows modifying the Aura config.
func (character *Character) NewTemporaryStatsAuraWrapped(auraLabel string, actionID ActionID, buffs stats.Stats, duration time.Duration, modConfig func(*Aura)) *StatBuffAura {
	// If one of the stat bonuses is a health bonus, then set up healing metrics for the associated
	// heal, since all temporary max health bonuses also instantaneously heal the player.
	var healthMetrics *ResourceMetrics
	amountHealed := buffs[stats.Health]
	includesHealthBuff := amountHealed > 0

	buffedStatTypes := buffs.GetBuffedStatTypes()

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

	return &StatBuffAura{
		Aura:            character.GetOrRegisterAura(config),
		BuffedStatTypes: buffedStatTypes,
	}
}

// Creates a new ProcTriggerAura that is dependent on a parent Aura being active
// This should only be used if the dependent Aura is:
// 1. On the a different Unit than parent Aura is registered to (usually the Character)
// 2. You need to register multiple dependent Aura's for the same Unit
func (parentAura *Aura) MakeDependentProcTriggerAura(unit *Unit, config ProcTrigger) *Aura {
	oldExtraCondition := config.ExtraCondition
	config.ExtraCondition = func(sim *Simulation, spell *Spell, result *SpellResult) bool {
		return parentAura.IsActive() && ((oldExtraCondition == nil) || oldExtraCondition(sim, spell, result))
	}

	aura := MakeProcTriggerAura(unit, config)

	return aura
}

// Attaches a ProcTrigger to a parent Aura
// Preffered use-case.
// For non standard use-cases see: MakeDependentProcTriggerAura
// Returns parent aura for chaining
func (parentAura *Aura) AttachProcTrigger(config ProcTrigger) *Aura {
	ApplyProcTriggerCallback(parentAura.Unit, parentAura, config)

	return parentAura
}

// Attaches a SpellMod to a parent Aura
// Returns parent aura for chaining
func (parentAura *Aura) AttachSpellMod(spellModConfig SpellModConfig) *Aura {
	parentAuraDep := parentAura.Unit.AddDynamicMod(spellModConfig)

	parentAura.ApplyOnGain(func(_ *Aura, _ *Simulation) {
		parentAuraDep.Activate()
	})

	parentAura.ApplyOnExpire(func(_ *Aura, _ *Simulation) {
		parentAuraDep.Deactivate()
	})

	return parentAura
}

// Attaches a StatDependency to a parent Aura
// Returns parent aura for chaining
func (parentAura *Aura) AttachStatDependency(statDep *stats.StatDependency) *Aura {

	parentAura.ApplyOnGain(func(_ *Aura, sim *Simulation) {
		parentAura.Unit.EnableBuildPhaseStatDep(sim, statDep)
	})

	parentAura.ApplyOnExpire(func(_ *Aura, sim *Simulation) {
		parentAura.Unit.DisableBuildPhaseStatDep(sim, statDep)
	})

	return parentAura
}

// Adds Stats to a parent Aura
// Returns parent aura for chaining
func (parentAura *Aura) AttachStatsBuff(stats stats.Stats) *Aura {
	parentAura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		aura.Unit.AddStatsDynamic(sim, stats)
	})

	parentAura.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		aura.Unit.AddStatsDynamic(sim, stats.Invert())
	})

	if parentAura.IsActive() {
		parentAura.Unit.AddStats(stats)
	}

	return parentAura
}

// Adds a Stat to a parent Aura
// Returns parent aura for chaining
func (parentAura *Aura) AttachStatBuff(stat stats.Stat, value float64) *Aura {
	statsToAdd := stats.Stats{}
	statsToAdd[stat] = value
	parentAura.AttachStatsBuff(statsToAdd)

	return parentAura
}

// Attaches a multiplicative PseudoStat buff to a parent Aura
// Returns parent aura for chaining
func (parentAura *Aura) AttachMultiplicativePseudoStatBuff(fieldPointer *float64, multiplier float64) *Aura {
	parentAura.ApplyOnGain(func(_ *Aura, _ *Simulation) {
		*fieldPointer *= multiplier
	})

	parentAura.ApplyOnExpire(func(_ *Aura, _ *Simulation) {
		*fieldPointer /= multiplier
	})

	if parentAura.IsActive() {
		*fieldPointer *= multiplier
	}

	return parentAura
}

// Attaches an additive PseudoStat buff to a parent Aura
// Returns parent aura for chaining
func (parentAura *Aura) AttachAdditivePseudoStatBuff(fieldPointer *float64, bonus float64) *Aura {
	parentAura.ApplyOnGain(func(_ *Aura, _ *Simulation) {
		*fieldPointer += bonus
	})

	parentAura.ApplyOnExpire(func(_ *Aura, _ *Simulation) {
		*fieldPointer -= bonus
	})

	if parentAura.IsActive() {
		*fieldPointer += bonus
	}

	return parentAura
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

	aura.ApplyOnExpire(func(_ *Aura, _ *Simulation) {
		aura.ShieldStrength = 0
	})

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
