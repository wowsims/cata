package protection

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerUnwaveringSentinel() {
	core.MakePermanent(war.GetOrRegisterAura(core.Aura{
		Label:    "Unwavering Sentinel",
		ActionID: core.ActionID{SpellID: 29144},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.ApplyDynamicEquipScaling(sim, stats.Armor, 1.25)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.RemoveDynamicEquipScaling(sim, stats.Armor, 1.25)
		},
	}).AttachStatDependency(
		war.NewDynamicMultiplyStat(stats.Stamina, 1.15),
	).AttachAdditivePseudoStatBuff(
		&war.PseudoStats.ReducedCritTakenChance, 0.06,
	).AttachSpellMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskThunderClap,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -2,
	}))
}

func (war *ProtectionWarrior) registerBastionOfDefense() {
	core.MakePermanent(war.GetOrRegisterAura(core.Aura{
		Label:    "Bastion of Defense",
		ActionID: core.ActionID{SpellID: 84608},
	}).AttachAdditivePseudoStatBuff(
		&war.PseudoStats.BaseBlockChance, 0.1,
	).AttachAdditivePseudoStatBuff(
		&war.PseudoStats.BaseDodgeChance, 0.02,
	).AttachSpellMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskShieldWall,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -1 * time.Minute,
	}))
}

func (war *ProtectionWarrior) registerSwordAndBoard() {
	war.SwordAndBoardAura = war.GetOrRegisterAura(core.Aura{
		Label:    "Sword and Board",
		ActionID: core.ActionID{SpellID: 46953},
		Duration: 5 * time.Second,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Sword and Board - Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: warrior.SpellMaskDevastate,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.3,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.SwordAndBoardAura.Activate(sim)
			war.ShieldSlam.CD.Reset()
		},
	})
}

func (war *ProtectionWarrior) registerUltimatum() {
	war.UltimatumAura = war.GetOrRegisterAura(core.Aura{
		Label:    "Ultimatum",
		ActionID: core.ActionID{SpellID: 122510},
		Duration: 10 * time.Second,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.HeroicStrikeCleaveCostMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if !war.InciteAura.IsActive() {
				war.HeroicStrikeCleaveCostMod.Deactivate()
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskHeroicStrike | warrior.SpellMaskCleave,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 100,
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Ultimatum - Consume",
		ClassSpellMask: warrior.SpellMaskHeroicStrike | warrior.SpellMaskCleave,
		Callback:       core.CallbackOnCastComplete,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.CurCast.Cost <= 0
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.UltimatumAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Ultimatum - Trigger",
		ActionID:       core.ActionID{SpellID: 122509},
		ClassSpellMask: warrior.SpellMaskShieldSlam,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.UltimatumAura.Activate(sim)
		},
	})
}

func (war *ProtectionWarrior) registerRiposte() {
	shared.RegisterRiposteEffect(&war.Character, 145674, 145672)
}
