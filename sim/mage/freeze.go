package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerFreezeSpell() {

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 33395},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFreeze,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 25,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		BonusCoefficient:         0.029,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := 0.409 * mage.ClassSpellScaling
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			mage.FingersOfFrostAura.Activate(sim)
			mage.FingersOfFrostAura.SetStacks(sim, 2)
		},
	})
}
