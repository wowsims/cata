package destruction

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warlock"
)

func (destruction *DestructionWarlock) registerConflagrateSpell() {
	destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 17962},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellConflagrate,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.16,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    destruction.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destruction.ImmolateDot.Dot(target).IsActive()
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           destruction.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//TODO: How does conflagrate work?
			baseDamage := 471.0 + 0.6*destruction.ImmolateDot.Dot(target).Spell.SpellPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if !result.Landed() {
				return
			}
		},
	})
}
