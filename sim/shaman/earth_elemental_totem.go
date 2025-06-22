package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) registerEarthElementalTotem(isGuardian bool) {

	actionID := core.ActionID{SpellID: 2062}

	totalDuration := time.Second * 60

	earthElementalAura := shaman.RegisterAura(core.Aura{
		Label:    "Earth Elemental Totem",
		ActionID: actionID,
		Duration: totalDuration,
	})

	shaman.EarthElementalTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: SpellMaskEarthElementalTotem,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 28.1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 5,
			},
			SharedCD: core.Cooldown{
				Timer:    shaman.GetOrInitTimer(&shaman.ElementalSharedCDTimer),
				Duration: time.Minute * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			if shaman.Totems.Earth != proto.EarthTotem_NoEarthTotem {
				shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + totalDuration
			}

			shaman.EarthElemental.EnableWithTimeout(sim, shaman.EarthElemental, totalDuration)

			// Add a dummy aura to show in metrics
			earthElementalAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.EarthElementalTotem,
		Type:  core.CooldownTypeDPS,
	})
}
