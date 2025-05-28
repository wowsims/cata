package demonology

import (
	"github.com/wowsims/mop/sim/core"
)

// Caster Form + Pet Damage = 1% per Masterypoint
func (demo *DemonologyWarlock) getNormalMasteryBonus() float64 {
	return demo.getNormalMasteryBonusFrom(demo.GetMasteryPoints())
}

func (demo *DemonologyWarlock) getNormalMasteryBonusFrom(points float64) float64 {
	return (points + 8) / 100
}

// Meta Damage = 3% per Mastery Point
func (demo *DemonologyWarlock) getMetaMasteryBonus() float64 {
	return demo.getMetaMasteryBonusFrom(demo.GetMasteryPoints())
}

func (demo DemonologyWarlock) getMetaMasteryBonusFrom(points float64) float64 {
	return (points + 8.0) * 3 / 100
}

func (demo *DemonologyWarlock) registerMasterDemonologist() {
	var scaleAction *core.PendingAction
	demo.Metamorphosis.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		if scaleAction != nil {
			scaleAction.Cancel(sim)
		}

		demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getNormalMasteryBonus()
		demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getMetaMasteryBonus()
	})
	demo.Metamorphosis.RelatedSelfBuff.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getMetaMasteryBonus()
		if scaleAction != nil {
			return
		}

		// Gain of new mastery bonus is dealyed and seems to be refrence accurate
		scaleAction = &core.PendingAction{
			NextActionAt: sim.CurrentTime + core.GCDDefault,
			Priority:     core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()
				scaleAction = nil
			},
			CleanUp: func(sim *core.Simulation) {
				demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()
				scaleAction = nil
			},
		}
		sim.AddPendingAction(scaleAction)
	})

	demo.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating, newMasteryRating float64) {
		if demo.Metamorphosis.RelatedSelfBuff.IsActive() {
			demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getMetaMasteryBonusFrom(core.MasteryRatingToMasteryPoints(oldMasteryRating))
			demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getMetaMasteryBonus()
		} else {
			if scaleAction == nil {
				demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getNormalMasteryBonusFrom(core.MasteryRatingToMasteryPoints(oldMasteryRating))
				demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()
			}
		}

		for _, pet := range demo.Pets {
			if pet.IsActive() {
				pet.PseudoStats.DamageDealtMultiplier /= 1 + demo.getNormalMasteryBonusFrom(core.MasteryRatingToMasteryPoints(oldMasteryRating))
				pet.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()
			}
		}
	})

	demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()

	for _, pet := range demo.Pets {
		oldEnable := pet.OnPetEnable
		pet.OnPetEnable = func(sim *core.Simulation) {
			if oldEnable != nil {
				oldEnable(sim)
			}

			pet.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()
		}

		oldDisable := pet.OnPetDisable
		pet.OnPetDisable = func(sim *core.Simulation) {
			if oldDisable != nil {
				oldDisable(sim)
			}

			pet.PseudoStats.DamageDealtMultiplier /= 1 + demo.getNormalMasteryBonus()
		}
	}
}
