package demonology

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
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

func (demo *DemonologyWarlock) getMetaMasteryBonusFrom(points float64) float64 {
	return (points + 8.0) * 3 / 100
}

func (demo *DemonologyWarlock) registerMasterDemonologist() {
	var scaleAction *core.PendingAction
	corruptionMod := demo.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: -1 + 1/demo.getMetaMasteryBonus(),
		ClassMask:  warlock.WarlockSpellCorruption,
	})

	corruptionModCaster := demo.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: demo.getNormalMasteryBonus(),
		ClassMask:  warlock.WarlockSpellCorruption,
	})

	demo.Metamorphosis.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		if scaleAction != nil {
			scaleAction.Cancel(sim)
		}

		corruptionMod.UpdateFloatValue(-1 + 1/(1+demo.getMetaMasteryBonus()))
		corruptionModCaster.UpdateFloatValue(demo.getNormalMasteryBonus())
		corruptionMod.Activate()
		corruptionModCaster.Activate()
		demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getNormalMasteryBonus()
		demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getMetaMasteryBonus()

	}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		demo.PseudoStats.DamageDealtMultiplier /= 1 + demo.getMetaMasteryBonus()
		corruptionMod.Deactivate()
		corruptionModCaster.Deactivate()
		if scaleAction != nil {
			return
		}

		// Gain of new mastery bonus is dealyed and seems to be refrence accurate
		scaleAction = &core.PendingAction{
			NextActionAt: sim.CurrentTime + core.GCDDefault,
			Priority:     core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				scaleAction = nil
				demo.PseudoStats.DamageDealtMultiplier *= 1 + demo.getNormalMasteryBonus()
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
			corruptionMod.UpdateFloatValue(-1 + 1/(1+demo.getMetaMasteryBonus()))
			corruptionModCaster.UpdateFloatValue(demo.getNormalMasteryBonus())
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
