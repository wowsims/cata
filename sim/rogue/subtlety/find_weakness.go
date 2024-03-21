package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (subRogue *SubtletyRogue) applyFindWeakness() {
	if subRogue.Talents.FindWeakness == 0 {
		return
	}

	debuffPower := .35 * float64(subRogue.Talents.FindWeakness)

	fwDebuff := subRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Find Weakness",
			Duration: time.Second * 10,
			ActionID: core.ActionID{SpellID: 91023},

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				// for _, spell := range subRogue.Spellbook {
				//
				// }
				// TODO Thebackstabi 3/20/2024 -- Update to AttackTables once completed
				aura.Unit.PseudoStats.ArmorMultiplier -= debuffPower
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.ArmorMultiplier += debuffPower
			},
		})
	})

	subRogue.RegisterAura(core.Aura{
		Label:    "Find Weakness",
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit == &subRogue.Unit && result.Landed() && (spell == subRogue.Garrote || spell == subRogue.Ambush) {
				aura := fwDebuff.Get(result.Target)
				aura.Activate(sim)
			}
		},
	})
}
