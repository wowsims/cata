package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// Increases all Shadow damage done by (20 + (<Mastery Rating>/600)*2.5)%.
func (uhdk *UnholyDeathKnight) registerMastery() {
	uhdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		oldMasteryMultiplier := uhdk.currentMasteryMultiplier
		uhdk.currentMasteryMultiplier = (1.0 + uhdk.getMasteryPercent(newMastery))
		uhdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= uhdk.currentMasteryMultiplier / oldMasteryMultiplier
	})

	core.MakePermanent(uhdk.RegisterAura(core.Aura{
		Label:    "Dreadblade" + uhdk.Label,
		ActionID: core.ActionID{SpellID: 77514},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			oldMasteryMultiplier := uhdk.currentMasteryMultiplier
			uhdk.currentMasteryMultiplier = 1.0 + uhdk.getMasteryPercent(uhdk.GetStat(stats.MasteryRating))
			uhdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= uhdk.currentMasteryMultiplier / oldMasteryMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			uhdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= uhdk.currentMasteryMultiplier
			uhdk.currentMasteryMultiplier = 1.0
		},
	}))
}

func (uhdk *UnholyDeathKnight) getMasteryPercent(masteryRating float64) float64 {
	return 0.2 + 0.025*core.MasteryRatingToMasteryPoints(masteryRating)
}
