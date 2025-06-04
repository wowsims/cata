package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (dk *DeathKnight) ApplyUnholyTalents() {
	// Epidemic
	if dk.Talents.Epidemic > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			IntValue:  []int32{0, 1, 2, 4}[dk.Talents.Epidemic],
			ClassMask: DeathKnightSpellDisease,
		})
	}

	// Virulence
	if dk.Talents.Virulence > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDisease,
			FloatValue: 0.1 * float64(dk.Talents.Virulence),
		})
	}

	// Morbidity
	if dk.Talents.Morbidity > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDeathCoil | DeathKnightSpellDeathCoilHeal,
			FloatValue: 0.05 * float64(dk.Talents.Morbidity),
		})

		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDeathAndDecay,
			FloatValue: 0.1 * float64(dk.Talents.Morbidity),
		})
	}

	// Contagion
	dk.applyContagion()

	// Rage of Rivendare
	if dk.Talents.RageOfRivendare > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellPlagueStrike | DeathKnightSpellScourgeStrike | DeathKnightSpellFesteringStrike,
			FloatValue: 0.15 * float64(dk.Talents.RageOfRivendare),
		})
	}

	// Ebon Plaguebringer
	dk.applyEbonPlaguebringer()

	// Sudden Doom
	dk.applySuddenDoom()

	// Shadow Infusion
	shadowInfusionAura := dk.applyShadowInfusion()

	// Dark Transformation
	dk.applyDarkTransformation(shadowInfusionAura)
}

func (dk *DeathKnight) applyContagion() {
	contagionMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.5 * float64(dk.Talents.Contagion),
		ClassMask:  DeathKnightSpellDisease,
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Contagion Activate",
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: DeathKnightSpellPestilence,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			contagionMod.Activate()
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Contagion Deactivate",
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: DeathKnightSpellPestilence,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			contagionMod.Deactivate()
		},
	})
}

func (dk *DeathKnight) ebonPlaguebringerDiseaseMultiplier(_ *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
	return core.TernaryFloat64(spell.Matches(DeathKnightSpellDisease), 1.0+0.15*float64(dk.Talents.EbonPlaguebringer), 1.0)
}

func (dk *DeathKnight) applyEbonPlaguebringer() {
	if dk.Talents.EbonPlaguebringer == 0 {
		return
	}

	dk.EbonPlagueAura = dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := core.EbonPlaguebringerAura(dk.GetCharacter(), target, dk.Talents.Epidemic, dk.Talents.EbonPlaguebringer)
		aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			core.EnableDamageDoneByCaster(DDBC_EbonPlaguebringer, DDBC_Total, dk.AttackTables[aura.Unit.UnitIndex], dk.ebonPlaguebringerDiseaseMultiplier)
		})
		aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
			core.DisableDamageDoneByCaster(DDBC_EbonPlaguebringer, dk.AttackTables[aura.Unit.UnitIndex])
		})
		return aura
	})
	dk.Env.RegisterPreFinalizeEffect(func() {
		dk.FrostFeverSpell.RelatedAuraArrays = dk.FrostFeverSpell.RelatedAuraArrays.Append(dk.EbonPlagueAura)
		dk.BloodPlagueSpell.RelatedAuraArrays = dk.BloodPlagueSpell.RelatedAuraArrays.Append(dk.EbonPlagueAura)
	})

	var lastDiseaseTarget *core.Unit = nil

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Ebon Plague Triggers",
		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !spell.Matches(DeathKnightSpellDisease) {
				return
			}

			lastDiseaseTarget = target
			dk.EbonPlagueAura.Get(target).Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(DeathKnightSpellDisease) {
				return
			}

			dk.EbonPlagueAura.Get(lastDiseaseTarget).UpdateExpires(spell.Dot(lastDiseaseTarget).ExpiresAt())
		},
	}))
}

func (dk *DeathKnight) applySuddenDoom() {
	if dk.Talents.SuddenDoom == 0 {
		return
	}

	mod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: DeathKnightSpellDeathCoil | DeathKnightSpellDeathCoilHeal,
		IntValue:  -100,
	})

	suddenDoomProcAura := dk.GetOrRegisterAura(core.Aura{
		Label:     "Sudden Doom Proc",
		ActionID:  core.ActionID{SpellID: 81340},
		Duration:  time.Second * 10,
		MaxStacks: 0,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(DeathKnightSpellDeathCoil) {
				return
			}

			if spell.CurCast.Cost > 0 {
				return
			}

			if dk.T13Dps2pc.IsActive() {
				aura.RemoveStack(sim)
			} else {
				aura.Deactivate(sim)
			}
		},
	})

	ppm := 1.0 * float64(dk.Talents.SuddenDoom)
	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:     "Sudden Doom",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto,
		Outcome:  core.OutcomeLanded,
		PPM:      ppm,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			suddenDoomProcAura.Activate(sim)

			// T13 2pc: Sudden Doom has a 30% chance to grant 2 charges when triggered instead of 1.
			suddenDoomProcAura.MaxStacks = core.TernaryInt32(dk.T13Dps2pc.IsActive(), 2, 0)
			if dk.T13Dps2pc.IsActive() {
				stacks := core.TernaryInt32(sim.Proc(0.3, "T13 2pc"), 2, 1)
				suddenDoomProcAura.SetStacks(sim, stacks)
			}
		},
	})
}

func (dk *DeathKnight) applyShadowInfusion() *core.Aura {
	if dk.Talents.ShadowInfusion == 0 || dk.Ghoul == nil {
		return nil
	}

	trackingAura := dk.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Infusion Dk",
		ActionID:  core.ActionID{SpellID: 91342},
		Duration:  time.Second * 30,
		MaxStacks: 5,
	})

	damageMod := dk.Ghoul.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.06,
	})

	aura := dk.Ghoul.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Infusion",
		ActionID:  core.ActionID{SpellID: 91342},
		Duration:  time.Second * 30,
		MaxStacks: 5,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			trackingAura.Activate(sim)
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			trackingAura.Deactivate(sim)
			damageMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			trackingAura.Activate(sim)
			trackingAura.SetStacks(sim, newStacks)
			damageMod.UpdateFloatValue(float64(newStacks) * 0.06)
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Shadow Infusion",
		Callback:       core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ClassSpellMask: DeathKnightSpellDeathCoil | DeathKnightSpellDeathCoilHeal,
		Outcome:        core.OutcomeLanded,
		ProcChance:     []float64{0.0, 0.33, 0.66, 1.0}[dk.Talents.ShadowInfusion],

		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			if dk.Ghoul.DarkTransformationAura.IsActive() {
				return
			}
			aura.Activate(sim)
			aura.AddStack(sim)
		},
	})

	return aura
}

func (dk *DeathKnight) applyDarkTransformation(shadowInfusionAura *core.Aura) {
	if !dk.Talents.DarkTransformation {
		return
	}

	actionID := core.ActionID{SpellID: 63560}

	trackingAura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Dark Transformation Dk",
		ActionID: actionID,
		Duration: time.Second * 30,
	})

	damageMod := dk.Ghoul.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.6,
	})

	clawMod := dk.Ghoul.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  GhoulSpellClaw,
		FloatValue: 0.2,
	})

	dk.Ghoul.DarkTransformationAura = dk.Ghoul.GetOrRegisterAura(core.Aura{
		Label:    "Dark Transformation",
		ActionID: actionID,
		Duration: time.Second * 30,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			trackingAura.Activate(sim)
			damageMod.Activate()
			clawMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			trackingAura.Deactivate(sim)
			damageMod.Deactivate()
			clawMod.Deactivate()
		},
	})

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellDarkTransformation,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return shadowInfusionAura.GetStacks() == shadowInfusionAura.MaxStacks
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shadowInfusionAura.Deactivate(sim)
			dk.Ghoul.DarkTransformationAura.Activate(sim)
		},
	})
}
