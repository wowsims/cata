package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (earthElemental *EarthElemental) registerPulverize() {
	earthElemental.Pulverize = earthElemental.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 118345},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    earthElemental.NewTimer(),
				Duration: time.Second * 40,
			},
		},

		DamageMultiplier: 1.5,
		CritMultiplier:   earthElemental.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !earthElemental.IsGuardian()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
