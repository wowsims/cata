package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/death_knight"
)

// Increases all Frost damage done by (16 + (<Mastery Rating>/600)*2)%.
func (fdk *FrostDeathKnight) registerMastery() {
	// Beta changes 2025-06-21: https://eu.forums.blizzard.com/en/wow/t/feedback-mists-of-pandaria-class-changes/576939/51
	// - Soul Reaper now scales with your Mastery. [New]
	// - Obliterate now scales with your Mastery. [New]
	// Undocummented: only 20% effective when using a Two-handed weapon.
	masteryMod := fdk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: death_knight.DeathKnightSpellObliterate | death_knight.DeathKnightSpellSoulReaper,
	})

	extraMultiplier := 1.0

	fdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		oldMasteryMultiplier := fdk.getMasteryPercent(oldMastery)
		newMasteryMultiplier := fdk.getMasteryPercent(newMastery)
		fdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= (1.0 + newMasteryMultiplier) / (1.0 + oldMasteryMultiplier)
		masteryMod.UpdateFloatValue(newMasteryMultiplier * extraMultiplier)
	})

	core.MakePermanent(fdk.RegisterAura(core.Aura{
		Label:    "Frozen Heart" + fdk.Label,
		ActionID: core.ActionID{SpellID: 77514},

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if mh := fdk.GetMHWeapon(); mh != nil && mh.HandType == proto.HandType_HandTypeTwoHand {
				extraMultiplier = 0.2
			} else {
				extraMultiplier = 1.0
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMultiplier := fdk.getMasteryPercent(fdk.GetStat(stats.MasteryRating))
			fdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1.0 + masteryMultiplier
			masteryMod.UpdateFloatValue(masteryMultiplier * extraMultiplier)
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

	fdk.RegisterItemSwapCallback(core.AllWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		if mh := fdk.GetMHWeapon(); mh != nil && mh.HandType == proto.HandType_HandTypeTwoHand {
			extraMultiplier = 0.2
		} else {
			extraMultiplier = 1.0
		}
	})
}

func (fdk *FrostDeathKnight) getMasteryPercent(masteryRating float64) float64 {
	return 0.16 + 0.02*core.MasteryRatingToMasteryPoints(masteryRating)
}
