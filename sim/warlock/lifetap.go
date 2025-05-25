package warlock

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (warlock *Warlock) registerLifeTap() {
	actionID := core.ActionID{SpellID: 1454}
	manaMetrics := warlock.NewManaMetrics(actionID)

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellLifeTap,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			restore := 0.15 * warlock.GetStat(stats.Health)
			warlock.AddMana(sim, restore, manaMetrics)
		},
	})
}
