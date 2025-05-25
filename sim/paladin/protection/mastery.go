package protection

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (prot *ProtectionPaladin) registerMastery() {
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Mastery: Divine Bulwark" + prot.Label,
		ActionID:   core.ActionID{SpellID: 76671},
		BuildPhase: core.CharacterBuildPhaseBuffs,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.ShieldOfTheRighteousAdditiveMultiplier = prot.getMasteryPercent()
		},
	})).AttachStatBuff(stats.BlockPercent, prot.getMasteryPercent())

	// Keep it updated when mastery changes
	prot.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
		prot.AddStatDynamic(sim, stats.BlockPercent, core.MasteryRatingToMasteryPoints(newMasteryRating-oldMasteryRating))
		prot.ShieldOfTheRighteousAdditiveMultiplier = prot.getMasteryPercent()
	})
}

func (prot *ProtectionPaladin) getMasteryPercent() float64 {
	return (8.0 + prot.GetMasteryPoints()) / 100.0
}
