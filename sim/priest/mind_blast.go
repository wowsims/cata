package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) registerMindBlastSpell() {
	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8092},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellMindBlast,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 17,
			PercentModifier: 100,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ThreatMultiplier: 1,
		BonusCoefficient: 1.104,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := priest.calcBaseDamage(sim, 1.557, 0.055)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := priest.calcBaseDamage(sim, 1.557, 0.055)
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
