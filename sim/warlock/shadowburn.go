package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerShadowBurnSpell() {
	if !warlock.Talents.Shadowburn {
		return
	}

	warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 17877},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellShadowBurn,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.15,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 15 * time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20()
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         1.056,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, Coefficient_Shadowburn, Variance_Shadowburn)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
