package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerDragonsBreathSpell() {
	/* 	if !mage.Talents.DragonsBreath {
		return
	} */

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31661},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: MageSpellDragonsBreath,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 7,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: 0.193,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := 1.378 * mage.ClassSpellScaling
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
