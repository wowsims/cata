package destruction

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warlock"
)

func (destruction *DestructionWarlock) registerConflagrate() {
	destruction.Conflagrate = destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 17962},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellConflagrate,

		ManaCost: core.ManaCostOptions{BaseCost: 0.16},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    destruction.NewTimer(),
				Duration: 10 * time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destruction.ImmolateDot.Dot(target).IsActive()
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   destruction.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.17599999905,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destruction.CalcScalingSpellDmg(warlock.Coefficient_ImmolateDot)
			immoDot := destruction.ImmolateDot.Dot(target)
			if !immoDot.IsActive() {
				panic("Casted conflagrate without active immolation on the target")
			}
			spell.DamageMultiplier *= float64(immoDot.NumberOfTicks) * 0.6
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier /= float64(immoDot.NumberOfTicks) * 0.6
		},
	})
}
