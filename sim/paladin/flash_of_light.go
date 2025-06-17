package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerFlashOfLight() {
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 19750},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: SpellMaskFlashOfLight,

		MaxRange: 40,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 37.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 1.12000000477,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if target.IsOpponent(&paladin.Unit) {
				target = &paladin.Unit
			}

			baseHealing := paladin.CalcAndRollDamageRange(sim, 11.03999996185, 0.11500000209)

			damageMultiplier := spell.DamageMultiplier
			if paladin.SelflessHealerAura.IsActive() && (target != &paladin.Unit || paladin.BastionOfGloryAura.IsActive()) {
				spell.DamageMultiplier *= 1.0 + 0.2*float64(paladin.SelflessHealerAura.GetStacks())
			}

			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			spell.DamageMultiplier = damageMultiplier

			spell.DealHealing(sim, result)
		},
	})
}
