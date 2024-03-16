package hunter

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *Hunter) registerAspectOfTheHawkSpell() {
	const improvedHawkProcChance = 0.1
	actionID := core.ActionID{SpellID: 61847}
	hunter.AspectOfTheHawkAura = hunter.NewTemporaryStatsAuraWrapped(
		"Aspect of the Hawk",
		actionID,
		stats.Stats{
			stats.RangedAttackPower: 300,
		},
		core.NeverExpires,
		func(aura *core.Aura) {
		})
	hunter.applySharedAspectConfig(true, hunter.AspectOfTheHawkAura)

	hunter.AspectOfTheHawk = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheHawkAura.Activate(sim)
		},
	})
}

// Todo: Implement Aspect of the Fox?

func (hunter *Hunter) applySharedAspectConfig(isHawk bool, aura *core.Aura) {
	if isHawk {
		aura.OnReset = func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		}
	}

	aura.Duration = core.NeverExpires
	aura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})
}
