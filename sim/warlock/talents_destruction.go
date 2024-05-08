package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warlock *Warlock) ApplyDestructionTalents() {
	//Bane
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellShadowBolt | WarlockSpellChaosBolt | WarlockSpellImmolate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Duration([]int{0, -100, -300, -500}[warlock.Talents.Bane]) * time.Millisecond,
	})

	// Shadow And Flame
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellIncinerate,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{0.0, 0.04, 0.08, 0.12}[warlock.Talents.ShadowAndFlame],
	})

	// Improved Immolate
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellImmolate | WarlockSpellImmolateDot,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{0.0, 0.1, 0.2}[warlock.Talents.ImprovedImmolate],
	})

	// Emberstorm
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellSoulFire,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(-500*warlock.Talents.Emberstorm),
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellIncinerate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration([]float64{0, -130, -250}[warlock.Talents.Emberstorm]),
	})

	warlock.registerImprovedSearingPain()
	warlock.registerBackdraft()
	warlock.registerShadowBurnSpell()

	warlock.registerBurningEmbers()

	warlock.registerSoulLeech()

	//FireAndBrimstoneDamage mod is in Immolate
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellConflagrate,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5.0 * float64(warlock.Talents.FireAndBrimstone) * core.CritRatingPerCritChance,
	})

	warlock.registerEmpoweredImp()

	// TODO: BANE OF HAVOC

	if warlock.Talents.ChaosBolt {
		warlock.registerChaosBoltSpell()
	}
}

func (warlock *Warlock) registerImprovedSearingPain() {
	if warlock.Talents.ImprovedSearingPain <= 0 {
		return
	}

	improvedSearingPain := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Rating,
		ClassMask:  WarlockSpellSearingPain,
		FloatValue: 20 * float64(warlock.Talents.ImprovedSearingPain) * core.CritRatingPerCritChance,
	})

	improvedSearingPainAura := warlock.RegisterAura(core.Aura{
		Label:    "Improved Searing Pain",
		ActionID: core.ActionID{SpellID: 17927},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			improvedSearingPain.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			improvedSearingPain.Deactivate()
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				improvedSearingPainAura.Activate(sim)
			}
		})
	})
}

func (warlock *Warlock) registerBackdraft() {
	if warlock.Talents.Backdraft <= 0 {
		return
	}

	castReduction := -0.10 * float64(warlock.Talents.Backdraft)

	castTimeMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellIncinerate | WarlockSpellChaosBolt,
		FloatValue: castReduction,
	})

	backdraft := warlock.RegisterAura(core.Aura{
		Label:     "Backdraft",
		ActionID:  core.ActionID{SpellID: 47260},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
			warlock.ShadowBolt.DefaultCast.GCD = time.Duration(float64(warlock.ShadowBolt.DefaultCast.GCD) * (1 - castReduction))
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) * (1 - castReduction))
			warlock.ChaosBolt.DefaultCast.GCD = time.Duration(float64(warlock.ChaosBolt.DefaultCast.GCD) * (1 - castReduction))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
			warlock.ShadowBolt.DefaultCast.GCD = time.Duration(float64(warlock.ShadowBolt.DefaultCast.GCD) / (1 - castReduction))
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) / (1 - castReduction))
			warlock.ChaosBolt.DefaultCast.GCD = time.Duration(float64(warlock.ChaosBolt.DefaultCast.GCD) / (1 - castReduction))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&(WarlockSpellShadowBolt|WarlockSpellIncinerate|WarlockSpellChaosBolt) > 0 {
				aura.RemoveStack(sim)
			}
		},
	})

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Backdraft Hidden Aura",
			ActionID: core.ActionID{SpellID: 47260},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.ClassSpellMask == WarlockSpellConflagrate {
					backdraft.Activate(sim)
					backdraft.SetStacks(sim, 3)
				}
			},
		}))
}

func (warlock *Warlock) registerBurningEmbers() {
	if warlock.Talents.BurningEmbers <= 0 {
		return
	}

	damageGainPerHit := 0.25 * float64(warlock.Talents.BurningEmbers)
	spellPowerMultiplier := 0.7 * float64(warlock.Talents.BurningEmbers)
	additionalDamage := warlock.ClassSpellScaling * []float64{0.0, Coefficient_BurningEmbers_1, Coefficient_BurningEmbers_2}[warlock.Talents.BurningEmbers]
	var burningEmberTicks int32 = 7

	burningEmbers := warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 85112},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellBurningEmbers,

		DamageMultiplier: 1,
		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Burning Embers",
			},
			NumberOfTicks: burningEmberTicks,
			TickLength:    time.Second * 1,
			//?AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})

	//TODO: Can't find a good way to duplicate this effect across two units, warlock and imp
	core.MakeProcTriggerAura(&warlock.Imp.Unit, core.ProcTrigger{
		Name:           "Burning Embers",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: WarlockSpellImpFireBolt,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dot := burningEmbers.Dot(result.Target)
			// Max damage is based on the formula which changes per talent point where x is based off of:
			//1: [(Spell power * 0.7 + x) / 7]
			//2: [(Spell power * 1.4 + x) / 7]
			maxDamagePerTick := ((dot.Spell.SpellPower() * spellPowerMultiplier) + additionalDamage) / float64(burningEmberTicks)

			// The damage per tick gain is then based off the damage of the Imp Firebolt or Soulfire that just hit
			dot.SnapshotBaseDamage = min(dot.SnapshotBaseDamage+result.Damage/float64(burningEmberTicks)*damageGainPerHit, maxDamagePerTick)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[result.Target.UnitIndex], true)
			burningEmbers.Cast(sim, result.Target)
		},
	})

	core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
		Name:           "Burning Embers",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: WarlockSpellSoulFire,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dot := burningEmbers.Dot(result.Target)
			// Max damage is based on the formula which changes per talent point where x is based off of:
			//1: [(Spell power * 0.7 + x) / 7]
			//2: [(Spell power * 1.4 + x) / 7]
			maxDamagePerTick := ((dot.Spell.SpellPower() * spellPowerMultiplier) + additionalDamage) / 7

			// The damage per tick gain is then based off the damage of the Imp Firebolt or Soulfire that just hit
			dot.SnapshotBaseDamage = min(dot.SnapshotBaseDamage+result.Damage*damageGainPerHit, maxDamagePerTick)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[result.Target.UnitIndex], true)
			burningEmbers.Cast(sim, result.Target)
		},
	})
}

func (warlock *Warlock) registerSoulLeech() {
	if warlock.Talents.SoulLeech <= 0 {
		return
	}

	actionID := core.ActionID{SpellID: 30295}
	restore := 0.02 * float64(warlock.Talents.SoulLeech)
	manaMetrics := warlock.NewManaMetrics(actionID)

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Soul Leech Hidden Aura",
			ActionID: actionID,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ClassSpellMask&(WarlockSpellShadowBurn|WarlockSpellSoulFire|WarlockSpellChaosBolt) > 0 {
					warlock.AddMana(sim, restore*warlock.MaxMana(), manaMetrics)
					// also restores health but probably NA
				}
			},
		}))
}

func (warlock *Warlock) registerEmpoweredImp() {
	if warlock.Talents.EmpoweredImp <= 0 || warlock.Options.Summon != proto.WarlockOptions_Imp {
		return
	}

	procChance := 0.02 * float64(warlock.Talents.EmpoweredImp)

	castTimeMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellSoulFire,
		FloatValue: -1,
	})

	empoweredImpAura := warlock.RegisterAura(core.Aura{
		Label:    "Empowered Imp",
		ActionID: core.ActionID{SpellID: 47221},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.SoulFire {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(
		warlock.Imp.RegisterAura(core.Aura{
			Label: "Empowered Imp Hidden Aura",

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ClassSpellMask == WarlockSpellImpFireBolt && sim.Proc(procChance, "Empowered Imp") {
					empoweredImpAura.Activate(sim)
				}
			},
		}))
}
