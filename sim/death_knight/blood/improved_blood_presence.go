package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (bdk *BloodDeathKnight) registerImprovedBloodPresence() {
	multi := 1.2
	impBloodPresenceAura := bdk.RegisterAura(core.Aura{
		Label:    "Improved Blood Presence" + bdk.Label,
		ActionID: core.ActionID{SpellID: 50371},
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bdk.MultiplyRuneRegenSpeed(sim, multi)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			bdk.MultiplyRuneRegenSpeed(sim, 1/multi)
		},
	}).AttachAdditivePseudoStatBuff(
		&bdk.PseudoStats.ReducedCritTakenChance, 0.06,
	)

	bdk.OnSpellRegistered(func(spell *core.Spell) {
		if !spell.Matches(death_knight.DeathKnightSpellBloodPresence) {
			return
		}

		bdk.BloodPresenceSpell.RelatedSelfBuff.AttachDependentAura(impBloodPresenceAura)
	})
}
