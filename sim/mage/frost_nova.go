package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerfrostNovaSpell() {

	frostNovaCoefficient := 0.19 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A122 Field "EffetBonusCoefficient"
	frostNovaScaling := 0.53     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A122 Field "Coefficient"
	frostNovaVariance := 0.15    // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A122 Field "Variance"

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 122},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagAoE,
		ClassSpellMask: MageSpellFrostNova,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 2,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: 25 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: frostNovaCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := mage.CalcAndRollDamageRange(sim, frostNovaScaling, frostNovaVariance)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
