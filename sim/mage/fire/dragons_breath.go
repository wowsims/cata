package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerDragonsBreathSpell() {

	dragonsBreathVariance := 0.15   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Variance"
	dragonsBreathScaling := 1.97    // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Coefficient"
	dragonsBreathCoefficient := .22 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "BonusCoefficient"

	fire.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31661},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellDragonsBreath,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    fire.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		DamageMultiplier: 1,
		CritMultiplier:   fire.DefaultCritMultiplier(),
		BonusCoefficient: dragonsBreathCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := fire.CalcAndRollDamageRange(sim, dragonsBreathScaling, dragonsBreathVariance)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
