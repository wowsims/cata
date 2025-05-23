package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

var coneOfColdCoefficient = 0.32 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A120 Field "EffetBonusCoefficient"
var coneOfColdScaling = 0.38     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A120 Field "Coefficient"
var coneOfColdVariance = 0.0     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A120 Field "Variance"

func (mage *Mage) registerConeOfColdSpell() {

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 120},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellConeOfCold,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: 10 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: coneOfColdCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mage.CalcAndRollDamageRange(sim, coneOfColdScaling, coneOfColdVariance)
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
