package mage

import (
	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerArcaneExplosionSpell() {

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1449},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneExplosion,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.186,
		ThreatMultiplier: 1 - 0.4*float64(mage.Talents.ImprovedArcaneExplosion),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.368 * mage.ClassSpellScaling
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
