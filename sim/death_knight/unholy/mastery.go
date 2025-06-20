package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/death_knight"
)

// Increases all Shadow damage done by (20 + (<Mastery Rating>/600)*2.5)%.
func (uhdk *UnholyDeathKnight) registerMastery() {
	// This needs to be a spell mod since the shadow part of SS ignores all multipliers except for SpellMods.
	masteryMod := uhdk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: death_knight.DeathKnightSpellScourgeStrikeShadow,
	})

	uhdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		oldMasteryMultiplier := uhdk.getMasteryPercent(oldMastery)
		newMasteryMultiplier := uhdk.getMasteryPercent(newMastery)
		uhdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= (1.0 + newMasteryMultiplier) / (1.0 + oldMasteryMultiplier)
		masteryMod.UpdateFloatValue(newMasteryMultiplier)
	})

	core.MakePermanent(uhdk.RegisterAura(core.Aura{
		Label:    "Dreadblade" + uhdk.Label,
		ActionID: core.ActionID{SpellID: 77514},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMultiplier := uhdk.getMasteryPercent(uhdk.GetStat(stats.MasteryRating))
			uhdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.0 + masteryMultiplier
			masteryMod.UpdateFloatValue(masteryMultiplier)
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))
}

func (uhdk *UnholyDeathKnight) getMasteryPercent(masteryRating float64) float64 {
	return 0.2 + 0.025*core.MasteryRatingToMasteryPoints(masteryRating)
}
