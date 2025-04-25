package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) registerFireElementalTotem() {

	actionID := core.ActionID{SpellID: 2894}

	totalDuration := time.Second * time.Duration(120*(1.0+0.20*float64(shaman.Talents.TotemicFocus)))

	fireElementalAura := shaman.RegisterAura(core.Aura{
		Label:    "Fire Elemental Totem",
		ActionID: actionID,
		Duration: totalDuration,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			shaman.FireElemental.ChangeStatInheritance(shaman.FireElemental.shamanOwner.fireElementalStatInheritance())
		},
	})

	shaman.FireElementalTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskFireElementalTotem,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 23,
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
			if shaman.Totems.Fire != proto.FireTotem_NoFireTotem {
				shaman.TotemExpirations[FireTotem] = sim.CurrentTime + totalDuration
			}

			shaman.MagmaTotem.AOEDot().Deactivate(sim)
			searingTotemDot := shaman.SearingTotem.Dot(shaman.CurrentTarget)
			if searingTotemDot != nil {
				searingTotemDot.Deactivate(sim)
			}

			shaman.FireElemental.Disable(sim)
			shaman.FireElemental.EnableWithTimeout(sim, shaman.FireElemental, totalDuration)

			// Add a dummy aura to show in metrics
			fireElementalAura.Activate(sim)
		},
	})

	/*shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.FireElementalTotem,
		Type:  core.CooldownTypeDPS,
	})*/
}
