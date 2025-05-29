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
		demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getNormalMasteryBonus()
		demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getMetaMasteryBonus()

	}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getMetaMasteryBonus()
		if scaleAction != nil {
			return
		}

		// Gain of new mastery bonus is dealy and seems to be refrence accurate
		scaleAction = &core.PendingAction{
			NextActionAt: sim.CurrentTime + core.GCDDefault,
			Priority:     core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				scaleAction = nil
				demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()
			},
		}
		sim.AddPendingAction(scaleAction)
	}).ApplyOnReset(func(aura *core.Aura, sim *core.Simulation) {
		// Make sure to execute the pending scale action if we reset the iteration
		if scaleAction != nil {
			scaleAction.OnAction(sim)
		}
	})

	demo.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating, newMasteryRating float64) {
		if demo.Metamorphosis.RelatedSelfBuff.IsActive() {
			demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getMetaMasteryBonusFrom(core.MasteryRatingToMasteryPoints(oldMasteryRating))
			demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getMetaMasteryBonus()
		} else {
			demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getNormalMasteryBonusFrom(core.MasteryRatingToMasteryPoints(oldMasteryRating))
			if scaleAction == nil {
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
