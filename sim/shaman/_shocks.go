package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// Shared logic for all shocks.
func (shaman *Shaman) newShockSpellConfig(spellID int32, spellSchool core.SpellSchool, baseCostPercent int32, shockTimer *core.Timer, bonusCoefficient float64) core.SpellConfig {
	actionID := core.ActionID{SpellID: spellID}

	return core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: spellSchool,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagShock | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: baseCostPercent,
			PercentModifier: 100 - (5 * shaman.Talents.Convection) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shockTimer,
				Duration: time.Second*6 - time.Millisecond*500*time.Duration(shaman.Talents.Reverberation),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: bonusCoefficient,
	}
}

func (shaman *Shaman) registerEarthShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8042, core.SpellSchoolNature, 18, shockTimer, 0.386)
	config.ClassSpellMask = SpellMaskEarthShock
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 0.92699998617, 0.05299999937)
		result := shaman.calcDamageStormstrikeCritChance(sim, target, baseDamage, spell)
		spell.DealDamage(sim, result)

		if shaman.Talents.Fulmination && shaman.LightningShieldAura.GetStacks() > 3 {
			shaman.Fulmination.Cast(sim, target)
			shaman.LightningShieldAura.SetStacks(sim, 3)
		}
	}

	shaman.EarthShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8050, core.SpellSchoolFire, 17, shockTimer, 0.214)

	config.ClassSpellMask = SpellMaskFlameShockDirect

	lavaBurstCritMod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskLavaBurst | SpellMaskLavaBurstOverload,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 100,
	})

	config.RelatedDotSpell = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 8050, Tag: 1},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            config.Flags & ^core.SpellFlagAPL | core.SpellFlagPassiveSpell,
		ClassSpellMask:   SpellMaskFlameShockDot,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FlameShock",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					lavaBurstCritMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					lavaBurstCritMod.Deactivate()
				},
			},
			NumberOfTicks:       6,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				baseDamage := shaman.ClassSpellScaling * 0.14200000465
				dot.Snapshot(target, baseDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
			spell.Dot(target).Apply(sim)
		},
	})

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.ClassSpellScaling * 0.52899998426
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		if result.Landed() {
			spell.RelatedDotSpell.Cast(sim, target)
		}
		spell.DealDamage(sim, result)
	}

	shaman.FlameShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerFrostShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8056, core.SpellSchoolFrost, 18, shockTimer, 0.386)
	config.ClassSpellMask = SpellMaskFrostShock
	config.ThreatMultiplier *= 2
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 0.86900001764, 0.05600000173)
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}

	shaman.FrostShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerShocks() {
	shockTimer := shaman.NewTimer()
	shaman.registerEarthShockSpell(shockTimer)
	shaman.registerFlameShockSpell(shockTimer)
	shaman.registerFrostShockSpell(shockTimer)
}
