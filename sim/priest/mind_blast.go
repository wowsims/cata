package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (priest *Priest) registerMindBlastSpell() {
	priest.MindBlast = priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8092},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellMindBlast,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultSpellCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: 1,
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
			baseDamage := priest.ScalingBaseDamage * 1.557
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := priest.ScalingBaseDamage * 1.557
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
