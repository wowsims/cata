package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/rogue"
)

func (asnRogue *AssassinationRogue) registerAllPassives() {
	asnRogue.registerBlindsidePassive()
}

func (asnRogue *AssassinationRogue) registerBlindsidePassive() {
	// Apply Mastery
	masteryMod := asnRogue.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  rogue.RogueSpellWoundPoison | rogue.RogueSpellDeadlyPoison | rogue.RogueSpellEnvenom | rogue.RogueSpellVenomousWounds,
		FloatValue: asnRogue.GetMasteryBonusFromRating(asnRogue.GetStat(stats.MasteryRating)),
	})
	masteryMod.Activate()

	asnRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryMod.UpdateFloatValue(asnRogue.GetMasteryBonusFromRating(newMastery))
	})

	// Assassin's Resolve: +25% Multiplicative all-school damage
	// +20 Energy handled in base rogue
	if asnRogue.HasDagger(core.MainHand) || asnRogue.HasDagger(core.OffHand) {
		asnRogue.PseudoStats.DamageDealtMultiplier *= 1.25
	}

	energyMod := asnRogue.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  rogue.RogueSpellDispatch,
		FloatValue: -2,
	})

	blindsideProc := asnRogue.RegisterAura(core.Aura{
		Label:    "Blindside",
		ActionID: core.ActionID{SpellID: 121153},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			energyMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			energyMod.Deactivate()
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellID == 111240 {
				// Dispatch casted, consume aura
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(core.MakeProcTriggerAura(&asnRogue.Unit, core.ProcTrigger{
		Name:           "Blindside Proc Trigger",
		ActionID:       core.ActionID{ItemID: 121152},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: rogue.RogueSpellMutilate,
		ProcChance:     0.3,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			blindsideProc.Activate(sim)
		},
	}))
}
