package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

// While in Unholy Presence, your haste is increased by an additional 10%.
func (uhdk *UnholyDeathKnight) registerImprovedUnholyPresence() {
	// Actual effect handled in presences.go
	impUnholyPresenceAura := uhdk.RegisterAura(core.Aura{
		Label:    "Improved Unholy Presence" + uhdk.Label,
		ActionID: core.ActionID{SpellID: 50392},
		Duration: core.NeverExpires,
	})

	uhdk.OnSpellRegistered(func(spell *core.Spell) {
		if !spell.Matches(death_knight.DeathKnightSpellUnholyPresence) {
			return
		}

		uhdk.UnholyPresenceSpell.RelatedSelfBuff.AttachDependentAura(impUnholyPresenceAura)
	})
}
