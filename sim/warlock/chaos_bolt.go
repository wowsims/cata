package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerChaosBoltSpell() {
	// ChaosBolt is affected by level-based partial resists.
	// TODO If there's bosses with elevated fire resistances, we'd need another spell flag,
	//  or add an unlimited amount of "bonusSpellPenetration".
	warlock.ChaosBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 50796},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagAPL,
		ClassSpellMask:   WarlockSpellChaosBolt,
		BonusCoefficient: 0.62800002098,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplierAdditive: 1 + warlock.GrandFirestoneBonus(),
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//TODO: Damage
			spell.CalcAndDealDamage(sim, target, sim.Roll(1429, 1813), spell.OutcomeMagicCrit)
		},
	})
}
