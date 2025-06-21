package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (fdk *FrostDeathKnight) registerImprovedFrostPresence() {
	fdk.OnSpellRegistered(func(spell *core.Spell) {
		if !spell.Matches(death_knight.DeathKnightSpellFrostPresence) {
			return
		}

		impFrostPresenceAura := fdk.RegisterAura(core.Aura{
			Label:    "Improved Frost Presence" + fdk.Label,
			ActionID: core.ActionID{SpellID: 50385},
			Duration: core.NeverExpires,
		}).AttachSpellMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			ClassMask: death_knight.DeathKnightSpellFrostStrike,
			IntValue:  -15,
		})

		fdk.FrostPresenceSpell.RelatedSelfBuff.AttachDependentAura(impFrostPresenceAura)
	})
}
