package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerInfernoBlastSpell() {

	infernoBlastVariance := 0.17   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=108853 Field: "Variance"
	infernoBlastScaling := .60     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=108853 Field: "Coefficient"
	infernoBlastCoefficient := .60 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=108853 Field: "BonusCoefficient"

	fire.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 108853},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellInfernoBlast,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 2,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    fire.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           fire.DefaultCritMultiplier(),
		BonusCoefficient:         infernoBlastCoefficient,
		ThreatMultiplier:         1,
		BonusCritPercent:         100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := fire.CalcAndRollDamageRange(sim, infernoBlastScaling, infernoBlastVariance)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				fire.ApplyIgnite(sim, target, result.Damage)
			}
		},
	})
}
