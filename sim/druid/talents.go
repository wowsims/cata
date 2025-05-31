package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (druid *Druid) ApplyTalents() {
	druid.registerHeartOfTheWild()
	druid.registerNaturesVigil()
}

func (druid *Druid) registerHeartOfTheWild() {
	if !druid.Talents.HeartOfTheWild {
		return
	}

	// Apply 6% increase to Stamina, Agility, and Intellect
	statMultiplier := 1.06
	druid.MultiplyStat(stats.Stamina, statMultiplier)
	druid.MultiplyStat(stats.Agility, statMultiplier)
	druid.MultiplyStat(stats.Intellect, statMultiplier)

	// The activation spec specific effects are not implemented - most likely irrelevant for the sim unless proven otherwise
}

func (druid *Druid) registerNaturesVigil() {
	if !druid.Talents.NaturesVigil {
		return
	}

	actionID := core.ActionID{SpellID: 124974}

	naturesVigilAura := druid.RegisterAura(core.Aura{
		Label:    "Nature's Vigil",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageDealtMultiplier *= 1.12
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageDealtMultiplier /= 1.12
		},
	})

	druid.RegisterSpell(Humanoid|Moonkin|Cat|Bear, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			naturesVigilAura.Activate(sim)
		},
	})
}
