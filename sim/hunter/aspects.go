package hunter

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *Hunter) registerAspectOfTheHawkSpell() {
	actionID := core.ActionID{SpellID: 13165}
	hunter.AspectOfTheHawkAura = hunter.NewTemporaryStatsAuraWrapped(
		"Aspect of the Hawk",
		actionID,
		stats.Stats{
			stats.RangedAttackPower: 300, // Base value 20 with coeff 6.0952501297 on RAP https://wago.tools/db2/SpellEffect?build=4.4.0.53750&filter[SpellID]=exact%3A13165&page=1
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
	aura.Duration = core.NeverExpires
	aura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})
}
