package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) registerShadowWordDeathSpell() {
	priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 32379},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellShadowWordDeath,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 12,
			PercentModifier: 100,
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.316,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 12,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.IsExecutePhase25() {
				spell.DamageMultiplier *= 3
			}
			spell.CalcAndDealDamage(sim, target, priest.ClassSpellScaling*0.357, spell.OutcomeMagicHitAndCrit)
			if sim.IsExecutePhase25() {
				spell.DamageMultiplier /= 3
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			if sim.IsExecutePhase25() {
				spell.DamageMultiplier *= 3
			}
			result := spell.CalcDamage(sim, target, priest.ClassSpellScaling*0.357, spell.OutcomeExpectedMagicHitAndCrit)
			if sim.IsExecutePhase25() {
				spell.DamageMultiplier /= 3
			}
			return result
		},
	})
}
