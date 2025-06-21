package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (uhdk *UnholyDeathKnight) registerImprovedUnholyPresence() {
	uhdk.OnSpellRegistered(func(spell *core.Spell) {
		if !spell.Matches(death_knight.DeathKnightSpellUnholyPresence) {
			return
		}

		impUnholyPresenceAura := uhdk.RegisterAura(core.Aura{
			Label:    "Improved Unholy Presence" + uhdk.Label,
			ActionID: core.ActionID{SpellID: 50392},
			Duration: core.NeverExpires,
		})

		uhdk.UnholyPresenceSpell.RelatedSelfBuff.AttachDependentAura(impUnholyPresenceAura)
	})
}
