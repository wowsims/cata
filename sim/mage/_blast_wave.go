package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerBlastWaveSpell() {
	if !mage.Talents.BlastWave {
		return
	}

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11113},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: MageSpellBlastWave,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 7,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Millisecond * 500,
				GCDMin: time.Millisecond * 500,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 15,
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultCritMultiplier(),
		BonusCoefficient:         0.193,
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var targetCount int32
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				targetCount++
				baseDamage := sim.Roll(1047, 1233)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
