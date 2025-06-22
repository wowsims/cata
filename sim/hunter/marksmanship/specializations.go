package marksmanship

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (mm *MarksmanshipHunter) ApplySpecialization() {
	mm.SteadyFocusAura()
	mm.PiercingShotsAura()
	mm.MasterMarksmanAura()
	// Hotfix only applies to MM
	mm.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  hunter.HunterSpellBarrage,
		FloatValue: 0.15,
	})

	//Careful Aim
	caCritMod := mm.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  hunter.HunterSpellAimedShot | hunter.HunterSpellSteadyShot,
		FloatValue: 75,
	})

	mm.RegisterResetEffect(func(sim *core.Simulation) {
		caCritMod.Activate()
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			caCritMod.Deactivate()
		})
	})
	//bombardmentActionId := core.ActionID{SpellID: 35110}
	//focusMetrics := mm.NewFocusMetrics(bombardmentActionId)
	dmgMod := mm.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  hunter.HunterSpellMultiShot,
		FloatValue: 0.6,
	})
	costMod := mm.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: hunter.HunterSpellMultiShot,
		IntValue:  -20,
	})

	bombardmentAura := mm.RegisterAura(core.Aura{
		Label:    "Bombardment",
		ActionID: core.ActionID{SpellID: 35110},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Activate()
			dmgMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Deactivate()
			dmgMod.Deactivate()

		},
	})
	core.MakeProcTriggerAura(&mm.Unit, core.ProcTrigger{
		Name:           "Bombardment",
		ActionID:       core.ActionID{ItemID: 35110},
		Callback:       core.CallbackOnSpellHitDealt,
		ProcChance:     1,
		ClassSpellMask: hunter.HunterSpellMultiShot,
		Outcome:        core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if bombardmentAura.IsActive() {
				bombardmentAura.Refresh(sim)
			} else {
				bombardmentAura.Activate(sim)
			}
		},
	})
}
func (mm *MarksmanshipHunter) MasterMarksmanAura() {
	var counter *core.Aura
	procChance := 0.5
	mmAura := mm.RegisterAura(core.Aura{
		Label:    "Ready, Set, Aim...",
		ActionID: core.ActionID{SpellID: 82925},
		Duration: time.Second * 8,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(hunter.HunterSpellAimedShot) && spell.CurCast.Cost == 0 {
				aura.Deactivate(sim) // Consume effect
			}
		},
	})

	counter = mm.RegisterAura(core.Aura{
		Label:     "Master Marksman",
		Duration:  time.Second * 30,
		ActionID:  core.ActionID{SpellID: 34486},
		MaxStacks: 2,
	})

	core.MakePermanent(mm.RegisterAura(core.Aura{
		Label: "Master Marksman Internal",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(hunter.HunterSpellSteadyShot) {
				return
			}
			if procChance == 1 || sim.Proc(procChance, "Master Marksman Proc") {
				if counter.GetStacks() == 2 {
					mmAura.Activate(sim)
					counter.Deactivate(sim)
				} else {
					if !counter.IsActive() {
						counter.Activate(sim)
					}
					counter.AddStack(sim)
				}
			}
		},
	}))
}
func (mm *MarksmanshipHunter) SteadyFocusAura() {
	attackspeedMultiplier := core.TernaryFloat64(mm.CouldHaveSetBonus(hunter.YaunGolSlayersBattlegear, 4), 1.25, 1.15)
	steadyFocusAura := mm.RegisterAura(core.Aura{
		Label:     "Steady Focus",
		ActionID:  core.ActionID{SpellID: 53224, Tag: 1},
		Duration:  time.Second * 20,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, attackspeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, 1/attackspeedMultiplier)
		},
	})

	core.MakePermanent(mm.RegisterAura(core.Aura{
		Label:     "Steady Focus Counter",
		ActionID:  core.ActionID{SpellID: 53224, Tag: 2},
		MaxStacks: 2,
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.ProcMask.Matches(core.ProcMaskRangedAuto) || spell.ActionID.SpellID == 0 || !spell.Flags.Matches(core.SpellFlagAPL) {
				return
			}
			if !spell.Matches(hunter.HunterSpellSteadyShot) {
				aura.SetStacks(sim, 1)
			} else {
				if aura.GetStacks() == 2 {
					steadyFocusAura.Activate(sim)
					aura.SetStacks(sim, 1)
				} else {
					aura.SetStacks(sim, 2)
				}
			}
		},
	}))
}

func (mm *MarksmanshipHunter) PiercingShotsAura() {
	psSpell := mm.RegisterSpell(core.SpellConfig{
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

	mm.RegisterAura(core.Aura{
		Label:    "Piercing Shots Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			if !spell.Matches(hunter.HunterSpellAimedShot) && !spell.Matches(hunter.HunterSpellSteadyShot) && !spell.Matches(hunter.HunterSpellChimeraShot) {
				return
			}

			dot := psSpell.Dot(result.Target)
			newDamage := result.Damage * 0.3

			dot.SnapshotBaseDamage = (dot.OutstandingDmg() + newDamage) / float64(dot.BaseTickCount+core.TernaryInt32(dot.IsActive(), 1, 0))
			psSpell.Cast(sim, result.Target)
		},
	})
}
