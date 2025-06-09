package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerUnwaveringSentinel() {
	staminaDep := war.NewDynamicMultiplyStat(stats.Stamina, 1.15)

	core.MakePermanent(war.GetOrRegisterAura(core.Aura{
		Label:    "Unwavering Sentinel",
		ActionID: core.ActionID{SpellID: 29144},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.ApplyDynamicEquipScaling(sim, stats.Armor, 0.25)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.RemoveDynamicEquipScaling(sim, stats.Armor, 0.25)
		},
	}).AttachStatDependency(
		staminaDep,
	).AttachAdditivePseudoStatBuff(
		&war.PseudoStats.ReducedCritTakenChance, 0.06,
	).AttachSpellMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskThunderClap,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
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
		Name:           "Sword and Board Trigger",
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
