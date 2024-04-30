package fire

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/mage"
)

func (fire *FireMage) registerPyroblastSpell() {
	fire.PyroblastDotImpact = fire.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11366},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          mage.SpellFlagMage | mage.HotStreakSpells | core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellPyroblast,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: core.TernaryFloat64(fire.HotStreakAura.IsActive(), 0, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 3500,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if fire.HotStreakAura.IsActive() {
					cast.CastTime = 0
					fire.HotStreakAura.Deactivate(sim)
				}
			},
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           fire.DefaultSpellCritMultiplier(),
		BonusCoefficient:         1.545,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.5 * fire.ClassSpellScaling
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					fire.PyroblastDot.Dot(target).Apply(sim)
					spell.DealDamage(sim, result)
				}
			})
		},
	})

	fire.PyroblastDot = fire.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11366}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          mage.SpellFlagMage,
		ClassSpellMask: mage.MageSpellPyroblastDot,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fire.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "PyroblastDoT",
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 3,
			BonusCoefficient:    0.180,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 0.175*fire.ClassSpellScaling)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
		},
	})
}
