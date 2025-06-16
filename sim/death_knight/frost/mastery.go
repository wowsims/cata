package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// Increases all Frost damage done by (16 + (<Mastery Rating>/600)*2)%.
func (fdk *FrostDeathKnight) registerMastery() {
	fdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		oldMasteryMultiplier := fdk.currentMasteryMultiplier
		fdk.currentMasteryMultiplier = (1.0 + fdk.getMasteryPercent(newMastery))
		fdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= fdk.currentMasteryMultiplier / oldMasteryMultiplier
	})

	core.MakePermanent(fdk.RegisterAura(core.Aura{
		Label:    "Frozen Heart" + fdk.Label,
		ActionID: core.ActionID{SpellID: 77514},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			oldMasteryMultiplier := fdk.currentMasteryMultiplier
			fdk.currentMasteryMultiplier = 1.0 + fdk.getMasteryPercent(fdk.GetStat(stats.MasteryRating))
			fdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= fdk.currentMasteryMultiplier / oldMasteryMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			fdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] /= fdk.currentMasteryMultiplier
			fdk.currentMasteryMultiplier = 1.0
		},
	}))
}

func (fdk *FrostDeathKnight) getMasteryPercent(masteryRating float64) float64 {
	return 0.16 + 0.02*core.MasteryRatingToMasteryPoints(masteryRating)
}
