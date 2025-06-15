package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// Shared logic for all shocks.
func (shaman *Shaman) newShockSpellConfig(spellID int32, spellSchool core.SpellSchool, baseCostPercent float64, shockTimer *core.Timer, bonusCoefficient float64) core.SpellConfig {
	actionID := core.ActionID{SpellID: spellID}

	return core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: spellSchool,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagShamanSpell | SpellFlagShock | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: baseCostPercent,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shockTimer,
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: bonusCoefficient,
	}
}

func (shaman *Shaman) registerEarthShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8042, core.SpellSchoolNature, 14.4, shockTimer, 0.58099997044)
	config.ClassSpellMask = SpellMaskEarthShock
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 1.92200005054, 0.1000000014)
		result := shaman.calcDamageStormstrikeCritChance(sim, target, baseDamage, spell)
		spell.DealDamage(sim, result)
	}

	shaman.EarthShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8050, core.SpellSchoolFire, 11.9, shockTimer, 0.44900000095)

	config.ClassSpellMask = SpellMaskFlameShockDirect

	config.RelatedDotSpell = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 8050, Tag: 1},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            config.Flags & ^core.SpellFlagAPL & ^SpellFlagShamanSpell | core.SpellFlagPassiveSpell,
		ClassSpellMask:   SpellMaskFlameShockDot,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FlameShock",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.20999999344,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				baseDamage := shaman.CalcScalingSpellDmg(0.26100000739)
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
		baseDamage := shaman.CalcScalingSpellDmg(0.97399997711)
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		if result.Landed() {
			spell.RelatedDotSpell.Cast(sim, target)
		}
		spell.DealDamage(sim, result)
	}

	shaman.FlameShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerFrostShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8056, core.SpellSchoolFrost, 21.1, shockTimer, 0.50999999046)
	config.ClassSpellMask = SpellMaskFrostShock
	config.ThreatMultiplier *= 2
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 1.1210000515, 0.05600000173)
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
