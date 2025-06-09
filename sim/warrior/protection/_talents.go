package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) ApplyTalents() {

	// Vigilance is not implemented as it requires a second friendly target
	// We can probably fake it or make it configurable or something, but I expect it wouldn't
	// make much of a difference as I think tanks getting hit by bosses will be sitting at or near
	// their max vengeance bonus pretty much all the time

	war.Warrior.ApplyCommonTalents()

	war.RegisterConcussionBlow()
	war.RegisterDevastate()

	war.applyBastionOfDefense()
	war.applyHeavyRepercussions()
	war.applyImprovedRevenge()
	war.applySwordAndBoard()
	war.applyThunderstruck()

	war.ApplyGlyphs()
}

func (war *ProtectionWarrior) applyImprovedRevenge() {
	if war.Talents.ImprovedRevenge == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskRevenge,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.3 * float64(war.Talents.ImprovedRevenge),
	})

	// extra hit is implemented inside of revenge
}

func (war *ProtectionWarrior) applyHeavyRepercussions() {
	if war.Talents.HeavyRepercussions == 0 {
		return
	}

	damageMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskShieldSlam,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.5 * float64(war.Talents.HeavyRepercussions),
	})

	buff := war.RegisterAura(core.Aura{
		Label:    "Heavy Repercussions",
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Heavy Repercussions Trigger",
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: warrior.SpellMaskShieldBlock,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buff.Activate(sim)
		},
	})
}

func (war *ProtectionWarrior) applySwordAndBoard() {
	if war.Talents.SwordAndBoard == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskDevastate,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5 * float64(war.Talents.SwordAndBoard),
	})

	costMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskShieldSlam,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	})

	actionID := core.ActionID{SpellID: 50227}
	buffAura := war.RegisterAura(core.Aura{
		Label:    "Sword and Board",
		ActionID: actionID,
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Activate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell.ClassSpellMask & warrior.SpellMaskShieldSlam) != 0 {
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Deactivate()
		},
	})

	procChance := 0.1 * float64(war.Talents.SwordAndBoard)
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Heavy Repercussions Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: warrior.SpellMaskDevastate | warrior.SpellMaskRevenge,
		Outcome:        core.OutcomeLanded,
		ProcChance:     procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.shieldSlam.CD.Reset()
			buffAura.Activate(sim)
		},
	})
}
