package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (priest *Priest) registerShadowWordDeathSpell() {
	spellBaseDamage := func(spellPower float64, sim *core.Simulation) float64 {
		base := priest.ScalingBaseDamage*0.357 + 0.316*spellPower

		// SW:D does 3x damage <= 25% hp
		if sim.IsExecutePhase25() {
			return base * 3
		}

		return base
	}

	priest.ShadowWordDeath = priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 32379},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: int64(PriestSpellShadowWordDeath),

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1,
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           1,

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
			spell.CalcAndDealDamage(sim, target, spellBaseDamage(spell.SpellPower(), sim), spell.OutcomeMagicHitAndCrit)
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			return spell.CalcDamage(sim, target, spellBaseDamage(spell.SpellPower(), sim), spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
