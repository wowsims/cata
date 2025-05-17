package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

func (moonkin *BalanceDruid) RegisterBalancePassives() {
	moonkin.registerMoonkinForm()
	moonkin.registerShootingStars()
	moonkin.registerBalanceOfPower()
	moonkin.registerEuphoria()
	moonkin.registerOwlkinFrenzy()
	moonkin.registerKillerInstinct()
	moonkin.registerLeatherSpecialization()
	moonkin.registerNaturalInsight()
	moonkin.registerTotalEclipse()
	moonkin.registerLunarShower()
	moonkin.registerNaturesGrace()
}

func (moonkin *BalanceDruid) registerMoonkinForm() {
	moonkin.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane | core.SpellSchoolNature,
		FloatValue: 0.1,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	moonkin.MultiplyStat(stats.Armor, 0.6)

	core.MakePermanent(moonkin.RegisterAura(core.Aura{
		Label: "Moonkin Form",
		ActionID: core.ActionID{
			SpellID: 4858,
		},
	}))

	core.MakePermanent(core.MoonkinAura(&moonkin.Unit))
}

func (moonkin *BalanceDruid) registerShootingStars() {
	ssCastTimeMod := moonkin.AddDynamicMod(core.SpellModConfig{
		ClassMask:  druid.DruidSpellStarsurge,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
	})

	ssAura := moonkin.RegisterAura(core.Aura{
		Label:    "Shooting Stars" + moonkin.Label,
		ActionID: core.ActionID{SpellID: 93400},
		Duration: time.Second * 12,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(druid.DruidSpellStarsurge) {
				return
			}

			ssCastTimeMod.Deactivate()
			aura.Deactivate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ssCastTimeMod.Activate()
			moonkin.Starsurge.CD.Reset()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ssCastTimeMod.Deactivate()
		},
	})

	core.MakeProcTriggerAura(&moonkin.Unit, core.ProcTrigger{
		Name:           "Shooting Stars Trigger" + moonkin.Label,
		Callback:       core.CallbackOnPeriodicDamageDealt,
		Outcome:        core.OutcomeCrit,
		ProcChance:     0.3,
		ClassSpellMask: druid.DruidSpellSunfireDoT | druid.DruidSpellMoonfireDoT,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			ssAura.Activate(sim)
		},
	})
}

func (moonkin *BalanceDruid) registerBalanceOfPower() {
	moonkin.AddStat(stats.SpellHitPercent, -moonkin.GetBaseStats()[stats.Spirit])
	moonkin.AddStatDependency(stats.Spirit, stats.SpellHitPercent, 1)
}

func (moonkin *BalanceDruid) registerNaturesGrace() {
	ngAura := moonkin.RegisterAura(core.Aura{
		Label:    "Nature's Grace",
		ActionID: core.ActionID{SpellID: 16886},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			moonkin.MultiplyCastSpeed(1.15)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			moonkin.MultiplyCastSpeed(1 / 1.15)
		},
	})

	if moonkin.HasEclipseBar() {
		moonkin.AddEclipseCallback(func(_ Eclipse, gained bool, sim *core.Simulation) {
			if gained {
				ngAura.Activate(sim)
			}
		})
	}
}

func (moonkin *BalanceDruid) registerEuphoria() {
	moonkin.SetSpellEclipseEnergy(druid.DruidSpellWrath, WrathBaseEnergyGain, WrathBaseEnergyGain*2)
	moonkin.SetSpellEclipseEnergy(druid.DruidSpellStarfire, StarfireBaseEnergyGain, StarfireBaseEnergyGain*2)
	moonkin.SetSpellEclipseEnergy(druid.DruidSpellStarsurge, StarsurgeBaseEnergyGain, StarsurgeBaseEnergyGain*2)
}

func (moonkin *BalanceDruid) registerOwlkinFrenzy() {}

func (moonkin *BalanceDruid) registerKillerInstinct() {}

func (moonkin *BalanceDruid) registerLeatherSpecialization() {}

func (moonkin *BalanceDruid) registerNaturalInsight() {}

func (moonkin *BalanceDruid) registerTotalEclipse() {}

func (moonkin *BalanceDruid) registerLunarShower() {
	lunarShowerDmgMod := moonkin.AddDynamicMod(core.SpellModConfig{
		ClassMask: druid.DruidSpellMoonfire | druid.DruidSpellSunfire,
		Kind:      core.SpellMod_DamageDone_Pct,
	})

	lunarShowerResourceMod := moonkin.AddDynamicMod(core.SpellModConfig{
		ClassMask: druid.DruidSpellMoonfire | druid.DruidSpellSunfire,
		Kind:      core.SpellMod_PowerCost_Pct,
	})

	var lunarShowerAura = moonkin.RegisterAura(core.Aura{
		Label:     "Lunar Shower",
		Duration:  time.Second * 3,
		ActionID:  core.ActionID{SpellID: 81192},
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			lunarShowerDmgMod.UpdateFloatValue(float64(aura.GetStacks()) * 0.45)
			lunarShowerDmgMod.Activate()

			lunarShowerResourceMod.UpdateIntValue(aura.GetStacks() * -30)
			lunarShowerResourceMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			lunarShowerDmgMod.Deactivate()
			lunarShowerResourceMod.Deactivate()
		},
	})

	moonkin.RegisterAura(core.Aura{
		Label:    "Lunar Shower Handler",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(druid.DruidSpellMoonfire | druid.DruidSpellSunfire) {
				return
			}

			// does not proc off procs
			if spell.ProcMask.Matches(core.ProcMaskProc) {
				return
			}

			if lunarShowerAura.IsActive() {
				if lunarShowerAura.GetStacks() < 3 {
					lunarShowerAura.AddStack(sim)
					lunarShowerAura.Refresh(sim)
				}
			} else {
				lunarShowerAura.Activate(sim)
				lunarShowerAura.SetStacks(sim, 1)
			}
		},
	})
}
