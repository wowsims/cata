package death_knight

import (
	//"github.com/wowsims/cata/sim/core/proto"

	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (dk *DeathKnight) ApplyFrostTalents() {
	// Nerves Of Cold Steel
	if dk.Talents.NervesOfColdSteel > 0 && dk.HasMHWeapon() && dk.HasOHWeapon() {
		dk.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*float64(dk.Talents.NervesOfColdSteel))

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

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Merciless Combat Proc",
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: DeathKnightSpellMercilessCombat,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.IsExecutePhase35() {
				debuffs.Get(result.Target).Activate(sim)
			}
		},
	})
}

func (dk *DeathKnight) applyRime() {
	if dk.Talents.Rime == 0 {
		return
	}

	rimeMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -1,
		ClassMask:  DeathKnightSpellIcyTouch | DeathKnightSpellHowlingBlast,
	})

	freezingFogAura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Freezing Fog",
		ActionID: core.ActionID{SpellID: 59052},
		Duration: time.Second * 15,

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
			aura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Rime",
		Callback:       core.CallbackOnSpellHitDealt,
		ProcMask:       core.ProcMaskMeleeMH,
		ClassSpellMask: DeathKnightSpellObliterate,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.45,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			freezingFogAura.Activate(sim)
		},
	})
}

func (dk *DeathKnight) applyKillingMachine() {
	if dk.Talents.KillingMachine == 0 {
		return
	}

	kmMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 100 * core.CritRatingPerCritChance,
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
			if spell.ClassSpellMask&(DeathKnightSpellObliterate|DeathKnightSpellFrostStrike) == 0 {
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

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:     "Killing Machine",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeLanded,
		PPM:      2.0 * float64(dk.Talents.KillingMachine),
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
