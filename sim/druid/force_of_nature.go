package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerForceOfNature() {
	if !druid.Talents.ForceOfNature {
		return
	}

	forceOfNatureAura := druid.RegisterAura(core.Aura{
		Label:    "Force of Nature",
		ActionID: core.ActionID{SpellID: 33831},
		Duration: time.Second * 30,
	})

	druid.ForceOfNature = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 33831},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.Treants.EnableWithTimeout(sim)
			druid.Treants.CancelGCDTimer(sim)

			forceOfNatureAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell:    druid.ForceOfNature.Spell,
		Priority: core.CooldownPriorityDrums + 1000, // Always prefer to use wolves before bloodlust/drums so wolves gain haste buff
		Type:     core.CooldownTypeDPS,
	})
}
