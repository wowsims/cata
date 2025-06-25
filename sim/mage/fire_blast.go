package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerFireBlastSpell() {
	mage.FireBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2136},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFireBlast,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 21,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{ // Note: Impact talent triggers CD refresh on spell *land*, not cast
				Timer:    mage.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: 0.429,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.113 * mage.ClassSpellScaling
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
