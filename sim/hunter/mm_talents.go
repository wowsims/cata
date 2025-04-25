package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) ApplyMMTalents() {
	if hunter.Talents.Efficiency > 0 {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			ClassMask: HunterSpellArcaneShot,
			IntValue:  -hunter.Talents.Efficiency,
		})
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			ClassMask: HunterSpellExplosiveShot | HunterSpellChimeraShot,
			IntValue:  -hunter.Talents.Efficiency * 2,
		})
	}
	if hunter.Talents.CarefulAim > 0 {
		caCritMod := hunter.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  HunterSpellAimedShot | HunterSpellCobraShot | HunterSpellSteadyShot,
			FloatValue: 30 * float64(hunter.Talents.CarefulAim),
		})

		hunter.RegisterResetEffect(func(sim *core.Simulation) {
			caCritMod.Activate()
			sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
				caCritMod.Deactivate()
			})
		})
	}

	if hunter.Talents.Posthaste > 0 {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			ClassMask: HunterSpellRapidFire,
			TimeValue: -(time.Minute * time.Duration(hunter.Talents.Posthaste)),
		})
	}
	hunter.registerSicEm()
	hunter.applyPiercingShots()
	hunter.applyGoForTheThroat()
	hunter.applyImprovedSteadyShot()
	hunter.registerReadinessCD()
	hunter.applyMasterMarksman()
	hunter.applyTermination()
	hunter.applyBombardment()
}

func (hunter *Hunter) applyTermination() {
	if hunter.Talents.Termination == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 83490}

	focusMetrics := hunter.NewFocusMetrics(actionID)
	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Termination",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.IsExecutePhase25() && spell == hunter.SteadyShot || spell == hunter.CobraShot {
				hunter.AddFocus(sim, float64(hunter.Talents.Termination)*3, focusMetrics)
			}
		},
	}))
}
func (hunter *Hunter) applyBombardment() {
	if hunter.Talents.Bombardment == 0 {
		return
	}
	costMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: HunterSpellMultiShot,
		IntValue:  -25 * hunter.Talents.Bombardment,
	})

	bombardmentAura := hunter.RegisterAura(core.Aura{
		Label:    "Bombardment",
		ActionID: core.ActionID{SpellID: 35110},
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Deactivate()
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Bombardment Proc",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == hunter.MultiShot && result.DidCrit() {
				bombardmentAura.Activate(sim)
			}
		},
	}))
}
func (hunter *Hunter) applyMasterMarksman() {
	if hunter.Talents.MasterMarksman == 0 {
		return
	}
	costMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: HunterSpellAimedShot,
		IntValue:  -100,
	})
	procChance := float64(hunter.Talents.MasterMarksman) * 0.2
	hunter.MasterMarksmanAura = hunter.RegisterAura(core.Aura{
		Label:    "Ready, Set, Aim...",
		ActionID: core.ActionID{SpellID: 82925},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.AimedShot != nil {
				costMod.Activate()
				hunter.AimedShot.DefaultCast.CastTime = 0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.AimedShot != nil {
				costMod.Deactivate()
				hunter.AimedShot.DefaultCast.CastTime = time.Second * 3
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == hunter.AimedShot {
				hunter.MasterMarksmanCounterAura.SetStacks(sim, 0)
				hunter.MasterMarksmanCounterAura.Activate(sim)
				aura.Deactivate(sim) // Consume effect
			}

		},
	})
	hunter.MasterMarksmanCounterAura = hunter.RegisterAura(core.Aura{
		Label:     "Master Marksman",
		Duration:  time.Second * 30,
		ActionID:  core.ActionID{SpellID: 34486},
		MaxStacks: 4,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != hunter.SteadyShot {
				return
			}
			if procChance == 1 || sim.Proc(procChance, "Master Marksman Proc") {
				if aura.GetStacks() == 4 {
					hunter.MasterMarksmanAura.Activate(sim)
				} else {
					aura.AddStack(sim)
				}
			}
		},
	})
}
func (hunter *Hunter) applyPiercingShots() {
	if hunter.Talents.PiercingShots == 0 {
		return
	}

	psSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53238},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "PiercingShots",
				Duration: time.Second * 8,
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Specifically account for bleed modifiers, since it still affects the spell, but we're ignoring all modifiers.
				dot.SnapshotAttackerMultiplier = target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Piercing Shots Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			if spell != hunter.AimedShot && spell != hunter.SteadyShot && spell != hunter.ChimeraShot {
				return
			}

			dot := psSpell.Dot(result.Target)
			newDamage := result.Damage * 0.1 * float64(hunter.Talents.PiercingShots)

			dot.SnapshotBaseDamage = (dot.OutstandingDmg() + newDamage) / float64(dot.BaseTickCount+core.TernaryInt32(dot.IsActive(), 1, 0))
			psSpell.Cast(sim, result.Target)
		},
	})
}

func (hunter *Hunter) applyGoForTheThroat() {
	if hunter.Talents.GoForTheThroat == 0 {
		return
	}
	if hunter.Pet == nil {
		return
	}

	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34950})

	amount := 5 * float64(hunter.Talents.GoForTheThroat)

	hunter.RegisterAura(core.Aura{
		Label:    "Go for the Throat",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskRangedAuto) || !result.DidCrit() {
				return
			}
			if !hunter.Pet.IsEnabled() {
				return
			}
			hunter.Pet.AddFocus(sim, amount, focusMetrics)
		},
	})
}
func (hunter *Hunter) applyImprovedSteadyShot() {
	if hunter.Talents.ImprovedSteadyShot == 0 {
		return
	}

	attackspeedMultiplier := 1 + (float64(hunter.Talents.ImprovedSteadyShot) * 0.05)
	hunter.ImprovedSteadyShotAura = hunter.RegisterAura(core.Aura{
		Label:     "Improved Steady Shot",
		ActionID:  core.ActionID{SpellID: 53221, Tag: 1},
		Duration:  time.Second * 8,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, attackspeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, 1/attackspeedMultiplier)
		},
	})
	hunter.ImprovedSteadyShotAuraCounter = core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label:     "Imp SS Counter",
		ActionID:  core.ActionID{SpellID: 53221, Tag: 2},
		MaxStacks: 2,
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.ProcMask.Matches(core.ProcMaskRangedAuto) || spell.ActionID.SpellID == 0 || !spell.Flags.Matches(core.SpellFlagAPL) {
				return
			}
			if spell != hunter.SteadyShot {
				aura.SetStacks(sim, 1)
			} else {
				if aura.GetStacks() == 2 {
					hunter.ImprovedSteadyShotAura.Activate(sim)
					aura.SetStacks(sim, 1)
				} else {
					aura.SetStacks(sim, 2)
				}
			}
		},
	}))
}
func (hunter *Hunter) registerSicEm() {
	if hunter.Talents.SicEm == 0 || hunter.Pet == nil {
		return
	}

	actionId := core.ActionID{SpellID: 83356}

	sicEmMod := hunter.Pet.AddDynamicMod(core.SpellModConfig{
		Kind:     core.SpellMod_PowerCost_Pct,
		IntValue: -hunter.Talents.SicEm * 50,
		ProcMask: core.ProcMaskMeleeMHSpecial,
	})

	sicEmAura := hunter.Pet.RegisterAura(core.Aura{
		ActionID: actionId,
		Label:    "Sic'Em",
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			sicEmMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			sicEmMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeMHSpecial) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Sic'Em Mod",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(HunterSpellArcaneShot | HunterSpellAimedShot | HunterSpellExplosiveShot) {
				if result.DidCrit() {
					sicEmAura.Activate(sim)
				}
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(HunterSpellExplosiveShot) {
				if result.DidCrit() {
					sicEmAura.Activate(sim)
				}
			}
		},
	}))
}
func (hunter *Hunter) registerReadinessCD() {
	if !hunter.Talents.Readiness {
		return
	}

	actionID := core.ActionID{SpellID: 23989}

	readinessSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 1,
			},
			IgnoreHaste: true, // Hunter GCD is locked
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use if there are no cooldowns to reset.
			return !hunter.RapidFire.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFire.CD.Reset()
			hunter.KillShot.CD.Reset()
			hunter.RaptorStrike.CD.Reset()
			hunter.ExplosiveTrap.CD.Reset()
			if hunter.KillCommand != nil {
				hunter.KillCommand.CD.Reset()
			}
			if hunter.ChimeraShot != nil {
				hunter.ChimeraShot.CD.Reset()
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.CD.Reset()
			}
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: readinessSpell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// If RF is about to become ready naturally, wait so we can get 2x usages.
			if !hunter.RapidFire.IsReady(sim) && hunter.RapidFire.TimeToReady(sim) < time.Second*10 {
				return false
			}
			return !hunter.RapidFireAura.IsActive() || hunter.RapidFireAura.RemainingDuration(sim) < time.Second*10
		},
	})
}
