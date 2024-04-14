package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerArcaneExplosionSpell() {

	var GCDReduction time.Duration
	switch mage.Talents.ImprovedArcaneExplosion {
	case 2:
		GCDReduction = time.Millisecond * 200
	case 1:
		GCDReduction = time.Millisecond * 500
	case 0:
		GCDReduction = 0
	}

	mage.ArcaneExplosion = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1449},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.22,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault - GCDReduction,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.186,
		ThreatMultiplier: 1 - 0.4*float64(mage.Talents.ImprovedArcaneExplosion),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := 0.368 * mage.ScalingBaseDamage
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
