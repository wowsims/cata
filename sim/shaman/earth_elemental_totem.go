package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) registerEarthElementalTotem() {

	actionID := core.ActionID{SpellID: 2062}

	totalDuration := time.Second * time.Duration(120*(1.0+0.20*float64(shaman.Talents.TotemicFocus)))

	earthElementalAura := shaman.RegisterAura(core.Aura{
		Label:    "Earth Elemental Totem",
		ActionID: actionID,
		Duration: totalDuration,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			shaman.EarthElemental.ChangeStatInheritance(shaman.EarthElemental.shamanOwner.earthElementalStatInheritance())
		},
	})

	shaman.EarthElementalTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskEarthElementalTotem,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 24,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 10,
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

	/*shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.EarthElementalTotem,
		Type:  core.CooldownTypeDPS,
	})*/
}
