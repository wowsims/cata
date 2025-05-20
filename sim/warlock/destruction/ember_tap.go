package destruction

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (destruction *DestructionWarlock) registerEmberTap() {
	metric := destruction.NewHealthMetrics(core.ActionID{SpellID: 114635})
	spell := destruction.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 114635},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskEmpty,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ClassSpellMask:   warlock.WarlockSpellEmberTap,
		Flags:            core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			destruction.BurningEmbers.Spend(10, spell.ActionID, sim)
			destruction.GainHealth(sim, destruction.MaxHealth()*(1.15*spell.DamageMultiplier), metric)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destruction.BurningEmbers.CanSpend(10)
		},
	})

	destruction.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return c.CurrentHealthPercent() < 0.5
		},
	})
}
