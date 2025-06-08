package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (moonkin *BalanceDruid) ApplyBalanceTalents() {
	moonkin.registerIncarnation()
	moonkin.registerDreamOfCenarius()
	moonkin.registerSoulOfTheForest()
}

func (moonkin *BalanceDruid) registerIncarnation() {
	actionID := core.ActionID{SpellID: 102560}

	incarnationSpellMod := moonkin.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane | core.SpellSchoolNature,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.25,
	})

	incarnationAura := moonkin.RegisterAura(core.Aura{
		Label:    "Incarnation: Chosen of Elune",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Only apply the damage bonus when in Eclipse
			if moonkin.HasEclipseBar() && moonkin.IsInEclipse() {
				incarnationSpellMod.Activate()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			incarnationSpellMod.Deactivate()
		},
	})

	// Add Eclipse callback to apply/remove damage bonus when entering/exiting Eclipse
	moonkin.AddEclipseCallback(func(eclipse Eclipse, gained bool, sim *core.Simulation) {
		if incarnationAura.IsActive() {
			if gained {
				incarnationSpellMod.Activate()
			} else {
				incarnationSpellMod.Deactivate()
			}
		}
	})

	moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    moonkin.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			incarnationAura.Activate(sim)
		},
	})
}

func (moonkin *BalanceDruid) registerDreamOfCenarius() {
	if !moonkin.Talents.DreamOfCenarius {
		return
	}

	moonkin.DreamOfCenarius = moonkin.RegisterAura(core.Aura{
		Label:    "Dream of Cenarius",
		ActionID: core.ActionID{SpellID: 145151},
		Duration: time.Second * 30,
	})

	core.MakeProcTriggerAura(&moonkin.Unit, core.ProcTrigger{
		Name:           "Dream of Cenarius Trigger",
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: druid.DruidSpellHealingTouch,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			moonkin.DreamOfCenarius.Activate(sim)
		},
	})
}

func (moonkin *BalanceDruid) registerSoulOfTheForest() {
	if !moonkin.Talents.SoulOfTheForest {
		return
	}

	moonkin.AstralInsight = moonkin.RegisterAura(core.Aura{
		Label:    "Astral Insight (SotF)",
		ActionID: core.ActionID{SpellID: 145138},
		Duration: time.Second * 30,
	})

	core.MakeProcTriggerAura(&moonkin.Unit, core.ProcTrigger{
		Name:           "Astral Insight (SotF) Trigger",
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: druid.DruidSpellWrath | druid.DruidSpellStarfire | druid.DruidSpellStarsurge,
		ProcChance:     0.08,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			moonkin.AstralInsight.Activate(sim)
		},
	})
}
