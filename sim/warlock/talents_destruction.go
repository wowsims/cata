package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) ApplyDestructionTalents() {
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellShadowBolt | WarlockSpellChaosBolt | WarlockSpellImmolate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -1 * time.Duration([]int{0, 100, 300, 500}[warlock.Talents.Bane]) * time.Millisecond,
	})

	//TODO: Add/Mult?
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellChaosBolt,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{0.0, 1.04, 1.08, 1.12}[warlock.Talents.ShadowAndFlame],
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellImmolate,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{1.0, 1.1, 1.2}[warlock.Talents.ImprovedImmolate],
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellSoulFire,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -1 * time.Duration([]int{0, 500, 1000}[warlock.Talents.Emberstorm]) * time.Millisecond,
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellIncinerate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -1 * time.Duration([]int{0, 130, 250}[warlock.Talents.Emberstorm]) * time.Millisecond,
	})

	warlock.registerImprovedSearingPain()
	warlock.registerImprovedSoulFire()
	warlock.registerBackdraft()

	if warlock.Talents.ChaosBolt {
		warlock.registerShadowBurnSpell()
	}

	//TODO: Burning Embers

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

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			//TODO: Does this need to deactivate somewhere?
			if isExecute == 25 {
				improvedSearingPain.Activate()
			}
		})
	})
}

func (warlock *Warlock) registerImprovedSoulFire() {
	if warlock.Talents.ImprovedSoulFire <= 0 {
		return
	}

	damageBonus := 1 + .04*float64(warlock.Talents.ImprovedSoulFire)

	improvedSoulFire := warlock.RegisterAura(core.Aura{
		Label:    "Improved Soul Fire",
		ActionID: core.ActionID{SpellID: 18120},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			//TODO: Add or mult?
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= damageBonus
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= damageBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			//TODO: Add or mult?
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= damageBonus
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= damageBonus
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Improved Soul Fire Hidden Aura",
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && (spell == warlock.SoulFire) {
				improvedSoulFire.Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) registerBackdraft() {
	if warlock.Talents.Backdraft <= 0 {
		return
	}

	castReduction := 0.10 * float64(warlock.Talents.Backdraft)

	backdraft := warlock.RegisterAura(core.Aura{
		Label:     "Backdraft",
		ActionID:  core.ActionID{SpellID: 47260},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier -= castReduction
			warlock.Incinerate.CastTimeMultiplier -= castReduction
			warlock.ChaosBolt.CastTimeMultiplier -= castReduction
			warlock.ShadowBolt.DefaultCast.GCD = time.Duration(float64(warlock.ShadowBolt.DefaultCast.GCD) * (1 - castReduction))
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) * (1 - castReduction))
			warlock.ChaosBolt.DefaultCast.GCD = time.Duration(float64(warlock.ChaosBolt.DefaultCast.GCD) * (1 - castReduction))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier += castReduction
			warlock.Incinerate.CastTimeMultiplier += castReduction
			warlock.ChaosBolt.CastTimeMultiplier += castReduction
			warlock.ShadowBolt.DefaultCast.GCD = time.Duration(float64(warlock.ShadowBolt.DefaultCast.GCD) / (1 - castReduction))
			warlock.Incinerate.DefaultCast.GCD = time.Duration(float64(warlock.Incinerate.DefaultCast.GCD) / (1 - castReduction))
			warlock.ChaosBolt.DefaultCast.GCD = time.Duration(float64(warlock.ChaosBolt.DefaultCast.GCD) / (1 - castReduction))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.ShadowBolt || spell == warlock.Incinerate || spell == warlock.ChaosBolt {
				aura.RemoveStack(sim)
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Backdraft Hidden Aura",
		ActionID: core.ActionID{SpellID: 47260},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == warlock.Conflagrate {
				backdraft.Activate(sim)
				backdraft.SetStacks(sim, 3)
			}
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

	warlock.RegisterAura(core.Aura{
		Label:    "Soul Leech Hidden Aura",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Shadowburn || spell == warlock.SoulFire || spell == warlock.ChaosBolt {
				warlock.AddMana(sim, restore*warlock.MaxMana(), manaMetrics)
				// also restores health but probably NA
			}
		},
	})
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

	warlock.EmpoweredImpAura = warlock.RegisterAura(core.Aura{
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

	warlock.Pet.RegisterAura(core.Aura{
		Label:    "Empowered Imp Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask == WarlockSpellImpFireBolt && sim.Proc(procChance, "Empowered Imp") {
				warlock.EmpoweredImpAura.Activate(sim)
			}
		},
	})
}
