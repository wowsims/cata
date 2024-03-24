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

	cdReduction := time.Duration(comRogue.Talents.RestlessBlades) * time.Second
	comRogue.RestlessBladesAura = comRogue.RegisterAura(core.Aura{
		Label:    "Restless Blades",
		ActionID: core.ActionID{SpellID: 79096},
		Duration: core.NeverExpires,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit == &comRogue.Unit && spell.Flags.Matches(rogue.SpellFlagFinisher) {
				ksNewTime := comRogue.KillingSpree.CD.Timer.ReadyAt() - cdReduction
				arNewTime := comRogue.AdrenalineRush.CD.Timer.ReadyAt() - cdReduction
				comRogue.KillingSpree.CD.Timer.Set(ksNewTime)
				comRogue.AdrenalineRush.CD.Timer.Set(arNewTime)
			}
		},
	})
}
