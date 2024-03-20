package combat

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/rogue"
)

func (comRogue *CombatRogue) applyRestlessBlades() {
	if comRogue.Talents.RestlessBlades == 0 {
		return
	}

	cdReduction := core.Ternary(comRogue.Talents.RestlessBlades == 2, time.Second*2, time.Second)
	comRogue.RestlessBladesAura = comRogue.RegisterAura(core.Aura{
		Label:    "Restless Blades",
		ActionID: core.ActionID{SpellID: 79096},
		Duration: core.NeverExpires,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit == &comRogue.Unit && spell.Flags.Matches(rogue.SpellFlagFinisher) {
				*comRogue.KillingSpree.CD.Timer -= core.Timer(cdReduction)
				*comRogue.AdrenalineRush.CD.Timer -= core.Timer(cdReduction)
			}
		},
	})
}
