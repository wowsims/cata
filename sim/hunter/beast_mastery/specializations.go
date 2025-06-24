package beast_mastery

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/hunter"
)

func (bmHunter *BeastMasteryHunter) ApplyTalents() {
	bmHunter.applyFrenzy()
	bmHunter.applyGoForTheThroat()
	bmHunter.applyCobraStrikes()
	bmHunter.applyInvigoration()
	bmHunter.applyBeastCleave()
	bmHunter.Hunter.ApplyTalents()
}

func (bmHunter *BeastMasteryHunter) applyFrenzy() {
	if bmHunter.Pet == nil {
		return
	}
	actionId := core.ActionID{SpellID: 19623}
	bmHunter.Pet.FrenzyAura = bmHunter.Pet.RegisterAura(core.Aura{
		Label:     "Frenzy",
		Duration:  time.Second * 30,
		ActionID:  actionId,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/(1+0.04*float64(oldStacks)))
			aura.Unit.MultiplyMeleeSpeed(sim, 1+0.04*float64(newStacks))
		},
	})

	procChance := 0.4
	bmHunter.Pet.RegisterAura(core.Aura{
		Label:    "FrenzyHandler",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(hunter.HunterPetFocusDump) {
				return
			}
			if sim.RandomFloat("Frenzy") >= procChance {
				return
			}
			if bmHunter.Pet.FrenzyAura.IsActive() {
				if bmHunter.Pet.FrenzyAura.GetStacks() != 5 {
					bmHunter.Pet.FrenzyAura.AddStack(sim)
				}
				bmHunter.Pet.FrenzyAura.Refresh(sim)
			} else {
				bmHunter.Pet.FrenzyAura.Activate(sim)
				bmHunter.Pet.FrenzyAura.SetStacks(sim, 1)
			}
		},
	})
}

func (bmHunter *BeastMasteryHunter) applyGoForTheThroat() {
	if bmHunter.Pet == nil {
		return
	}

	focusMetrics := bmHunter.Pet.NewFocusMetrics(core.ActionID{SpellID: 34953})

	bmHunter.RegisterAura(core.Aura{
		Label:    "Go for the Throat",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !(spell.OtherID == proto.OtherAction_OtherActionShoot && result.Outcome.Matches(core.OutcomeCrit)) {
				return
			}

			bmHunter.Pet.AddFocus(sim, 15, focusMetrics)
		},
	})
}

func (bmHunter *BeastMasteryHunter) applyCobraStrikes() {
	if bmHunter.Pet == nil {
		return
	}

	basicAttackCritMod := bmHunter.Pet.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ProcMask:   core.ProcMaskMeleeMHSpecial,
		FloatValue: 100,
	})

	actionID := core.ActionID{SpellID: 53260}
	procChance := 0.15
	csAura := bmHunter.Pet.RegisterAura(core.Aura{
		Label:     "Cobra Strikes",
		ActionID:  actionID,
		Duration:  time.Second * 15,
		MaxStacks: 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			basicAttackCritMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			basicAttackCritMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial) {
				aura.RemoveStack(sim)
			}
		},
	})

	bmHunter.RegisterAura(core.Aura{
		Label:    "Cobra Strikes",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(hunter.HunterSpellArcaneShot) {
				return
			}

			if sim.RandomFloat("Cobra Strikes") < procChance {
				if csAura.IsActive() {
					csAura.SetStacks(sim, csAura.GetStacks()+2)
					csAura.Refresh(sim) // TODO: Confirm how stacking works
				} else {
					csAura.Activate(sim)
					csAura.SetStacks(sim, 2)
				}
			}
		},
	})
}

func (bmHunter *BeastMasteryHunter) applyInvigoration() {
	if bmHunter.Pet == nil {
		return
	}

	focusMetrics := bmHunter.NewFocusMetrics(core.ActionID{SpellID: 53253})

	procChance := 0.15
	bmHunter.Pet.RegisterAura(core.Aura{
		Label:    "Invigoration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(hunter.HunterPetFocusDump) {
				if sim.RandomFloat("Invigoration") < procChance {
					bmHunter.AddFocus(sim, 20, focusMetrics)
				}
			}
		},
	})
}

func (bmHunter *BeastMasteryHunter) applyBeastCleave() {
	if bmHunter.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 115939}

	var copyDamage float64
	hitSpell := bmHunter.Pet.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: hunter.HunterPetBeastCleaveHit,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 0.75,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, copyDamage, spell.OutcomeAlwaysHit)
		},
	})

	beastCleaveAura := core.MakeProcTriggerAura(&bmHunter.Pet.Unit, core.ProcTrigger{
		Name:     "Beast Cleave",
		ActionID: actionID,
		Duration: time.Second * 4,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if bmHunter.Env.GetNumTargets() < 2 || result.Damage <= 0 || spell.Matches(hunter.HunterPetBeastCleaveHit) {
				return
			}

			// FIXME: This ignores target modifiers and assumes they are the same for the original target and the cleaved target
			copyDamage = result.Damage / result.ArmorMultiplier

			nextTarget := bmHunter.Env.NextTargetUnit(result.Target)
			for nextTarget != nil && nextTarget.Index != result.Target.Index {
				hitSpell.Cast(sim, nextTarget)
				nextTarget = bmHunter.Env.NextTargetUnit(nextTarget)
			}
		},
	})

	bmHunter.RegisterAura(core.Aura{
		Label:    "Beast Cleave",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(hunter.HunterSpellMultiShot) {
				if beastCleaveAura.IsActive() {
					beastCleaveAura.Refresh(sim)
				} else {
					beastCleaveAura.Activate(sim)
				}
			}
		},
	})
}
