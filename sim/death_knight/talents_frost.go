package death_knight

import (
	//"github.com/wowsims/mop/sim/core/proto"

	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (dk *DeathKnight) ApplyFrostTalents() {
	// Nerves Of Cold Steel
	if dk.Talents.NervesOfColdSteel > 0 && dk.HasMHWeapon() && dk.HasOHWeapon() {
		dk.AddStat(stats.PhysicalHitPercent, float64(dk.Talents.NervesOfColdSteel))

		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: []float64{0.0, 0.08, 0.16, 0.25}[dk.Talents.NervesOfColdSteel],
			ProcMask:   core.ProcMaskMeleeOH,
		})
	}

	// Annihilation
	if dk.Talents.Annihilation > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellObliterate,
			FloatValue: 0.15 * float64(dk.Talents.Annihilation),
		})
	}

	// Chill of the Grave
	dk.applyChillOfTheGrave()

	// Killing Machine
	dk.applyKillingMachine()

	// Merciless Combat
	dk.applyMercilessCombat()

	// Rime
	dk.applyRime()

	// Improved Icy Talons
	if dk.Talents.ImprovedIcyTalons {
		dk.PseudoStats.MeleeSpeedMultiplier *= 1.05
	}

	// Brittle Bones
	if dk.Talents.BrittleBones > 0 {
		dk.MultiplyStat(stats.Strength, 1.0+0.02*float64(dk.Talents.BrittleBones))
	}

	// Might of the Frozen Wastes
	dk.applyMightOfTheFrozenWastes()
}

const DeathKnightChillOfTheGrave = DeathKnightSpellIcyTouch | DeathKnightSpellHowlingBlast | DeathKnightSpellObliterate

func (dk *DeathKnight) applyChillOfTheGrave() {
	if dk.Talents.ChillOfTheGrave == 0 {
		return
	}

	rpAmount := 5.0 * float64(dk.Talents.ChillOfTheGrave)
	rpMetric := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 50115})
	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:            "Chill of the Grave",
		Callback:        core.CallbackOnSpellHitDealt,
		ClassSpellMask:  DeathKnightChillOfTheGrave,
		Outcome:         core.OutcomeLanded,
		ProcMaskExclude: core.ProcMaskMeleeOH, // Dont trigger on Obliterate Off hand
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dk.AddRunicPower(sim, rpAmount, rpMetric)
		},
	})
}

const DeathKnightSpellMercilessCombat = DeathKnightSpellIcyTouch | DeathKnightSpellObliterate | DeathKnightSpellFrostStrike | DeathKnightSpellHowlingBlast

func (dk *DeathKnight) mercilessCombatMultiplier(sim *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
	if spell.ClassSpellMask&(DeathKnightSpellMercilessCombat) == 0 {
		return 1.0
	}
	return 1.0 + 0.06*float64(dk.Talents.MercilessCombat)
}

func (dk *DeathKnight) applyMercilessCombat() {
	if dk.Talents.MercilessCombat == 0 {
		return
	}

	debuffs := dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := target.GetOrRegisterAura(core.Aura{
			Label:    "Merciless Combat" + dk.Label,
			ActionID: core.ActionID{SpellID: 49538},
			Duration: core.NeverExpires,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				core.EnableDamageDoneByCaster(DDBC_MercilessCombat, DDBC_Total, dk.AttackTables[aura.Unit.UnitIndex], dk.mercilessCombatMultiplier)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				core.DisableDamageDoneByCaster(DDBC_MercilessCombat, dk.AttackTables[aura.Unit.UnitIndex])
			},
		})
		return aura
	})

	dk.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 35 {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					debuffs.Get(aoeTarget).Activate(sim)
				}
			}
		})
	})
}

func (dk *DeathKnight) applyRime() {
	if dk.Talents.Rime == 0 {
		return
	}

	rimeMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
		ClassMask: DeathKnightSpellIcyTouch | DeathKnightSpellHowlingBlast,
	})

	freezingFogAura := dk.GetOrRegisterAura(core.Aura{
		Label:     "Freezing Fog",
		ActionID:  core.ActionID{SpellID: 59052},
		Duration:  time.Second * 15,
		MaxStacks: 0,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rimeMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&(DeathKnightSpellIcyTouch|DeathKnightSpellHowlingBlast) == 0 {
				return
			}

			if dk.T13Dps2pc.IsActive() {
				aura.RemoveStack(sim)
			} else {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Rime",
		Callback:       core.CallbackOnSpellHitDealt,
		ProcMask:       core.ProcMaskMeleeMH,
		ClassSpellMask: DeathKnightSpellObliterate,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.15 * float64(dk.Talents.Rime),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			freezingFogAura.Activate(sim)

			// T13 2pc: Rime has a 60% chance to grant 2 charges when triggered instead of 1.
			freezingFogAura.MaxStacks = core.TernaryInt32(dk.T13Dps2pc.IsActive(), 2, 0)
			if dk.T13Dps2pc.IsActive() {
				stacks := core.TernaryInt32(sim.Proc(0.6, "T13 2pc"), 2, 1)
				freezingFogAura.SetStacks(sim, stacks)
			}
		},
	})
}

func (dk *DeathKnight) applyKillingMachine() {
	if dk.Talents.KillingMachine == 0 {
		return
	}

	kmMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 100,
		ClassMask:  DeathKnightSpellObliterate | DeathKnightSpellFrostStrike,
	})

	kmAura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Killing Machine Proc",
		ActionID: core.ActionID{SpellID: 51124},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			kmMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			kmMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(DeathKnightSpellObliterate | DeathKnightSpellFrostStrike) {
				return
			}
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				return
			}
			aura.Deactivate(sim)
		},
	})

	// Dummy spell to react with triggers
	kmProcSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51124},
		Flags:          core.SpellFlagNoLogs | core.SpellFlagNoMetrics,
		ClassSpellMask: DeathKnightSpellKillingMachine,
	})

	ppm := 2.0 * float64(dk.Talents.KillingMachine)
	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:     "Killing Machine",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeLanded,
		PPM:      ppm,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			kmAura.Activate(sim)
			kmProcSpell.Cast(sim, nil)
		},
	})

}

func (dk *DeathKnight) applyMightOfTheFrozenWastes() {
	if dk.Talents.MightOfTheFrozenWastes == 0 || dk.Equipment.MainHand().HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	dk.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{0.0, 0.03, 0.6, 0.10}[dk.Talents.MightOfTheFrozenWastes],
		ProcMask:   core.ProcMaskMelee,
	})

	rpMetric := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 81331})
	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:       "Might of the Frozen Wastes",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15 * float64(dk.Talents.MightOfTheFrozenWastes),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dk.AddRunicPower(sim, 10, rpMetric)
		},
	})
}

func (dk *DeathKnight) ThreatOfThassarianProc(sim *core.Simulation, result *core.SpellResult, ohSpell *core.Spell) {
	if dk.Talents.ThreatOfThassarian == 0 || dk.GetOHWeapon() == nil {
		return
	}
	if sim.Proc([]float64{0.0, 0.3, 0.6, 1.0}[dk.Talents.ThreatOfThassarian], "Threat of Thassarian") {
		ohSpell.Cast(sim, result.Target)
	}
}
