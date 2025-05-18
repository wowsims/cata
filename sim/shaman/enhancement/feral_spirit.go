package enhancement

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/shaman"
)

func (enh *EnhancementShaman) registerFeralSpirit() {
	spiritWolvesActiveAura := enh.RegisterAura(core.Aura{
		Label:    "Feral Spirit",
		ActionID: core.ActionID{SpellID: 51533},
		Duration: time.Second * 30,
	})

	enh.FeralSpirit = enh.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51533},
		ClassSpellMask: shaman.SpellMaskFeralSpirit,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    enh.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			enh.SpiritWolves.EnableWithTimeout(sim)
			enh.SpiritWolves.CancelGCDTimer(sim)

			// Add a dummy aura to show in metrics
			spiritWolvesActiveAura.Activate(sim)

			// https://github.com/JamminL/wotlk-classic-bugs/issues/280
			// instant casts (e.g. shocks) usually don't reset a shaman's swing timer
			enh.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		},
	})

	enh.AddMajorCooldown(core.MajorCooldown{
		Spell:    enh.FeralSpirit,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
