package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var tigerStrikesMHID = core.ActionID{SpellID: 120274}
var tigerStrikesOHID = core.ActionID{SpellID: 120278}

func tigerStrikesBuffAura(unit *core.Unit) {
	var tigerStrikesMHSpell *core.Spell
	var tigerStrikesOHSpell *core.Spell
	var tigerStrikesBuff *core.Aura
	tigerStrikesBuff = unit.RegisterAura(core.Aura{
		Label:     "Tiger Strikes" + unit.Label,
		ActionID:  core.ActionID{SpellID: 120273},
		Duration:  time.Second * 15,
		MaxStacks: 4,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			mhConfig := *unit.AutoAttacks.MHConfig()
			mhConfig.ActionID = tigerStrikesMHID
			mhConfig.ClassSpellMask = MonkSpellTigerStrikes
			mhConfig.Flags |= core.SpellFlagPassiveSpell
			mhConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, target, baseDamage*unit.AutoAttacks.MHAuto().DamageMultiplier, spell.OutcomeMeleeWhiteNoGlance)
			}
			tigerStrikesMHSpell = unit.GetOrRegisterSpell(mhConfig)

			if unit.AutoAttacks.OH() != nil {
				ohConfig := *unit.AutoAttacks.OHConfig()
				ohConfig.ActionID = tigerStrikesOHID
				ohConfig.ClassSpellMask = MonkSpellTigerStrikes
				ohConfig.Flags |= core.SpellFlagPassiveSpell
				ohConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
					spell.CalcAndDealDamage(sim, target, baseDamage*unit.AutoAttacks.OHAuto().DamageMultiplier, spell.OutcomeMeleeWhiteNoGlance)
				}
				tigerStrikesOHSpell = unit.GetOrRegisterSpell(ohConfig)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			unit.MultiplyMeleeSpeed(sim, 1.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			unit.MultiplyMeleeSpeed(sim, 1/1.5)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// TODO: Verify if it actually procs even on misses, seems so based on logs and simc
			if !tigerStrikesBuff.IsActive() || spell.Matches(MonkSpellTigerStrikes) {
				return
			}

			if spell == unit.AutoAttacks.MHAuto() {
				tigerStrikesBuff.RemoveStack(sim)
				tigerStrikesMHSpell.Cast(sim, result.Target)
			} else if spell == unit.AutoAttacks.OHAuto() {
				tigerStrikesBuff.RemoveStack(sim)
				tigerStrikesOHSpell.Cast(sim, result.Target)
			}
		},
	})

	core.MakeProcTriggerAura(unit, core.ProcTrigger{
		Name:       "Tiger Strikes Buff Trigger" + unit.Label,
		ActionID:   core.ActionID{SpellID: 120272},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskWhiteHit,
		ProcChance: 1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(MonkSpellTigerStrikes) {
				return
			}

			if sim.Proc(0.08, "Tiger Strikes") {
				tigerStrikesBuff.Activate(sim)
				tigerStrikesBuff.SetStacks(sim, 4)
			}
		},
	})
}

func (monk *Monk) registerTigerStrikes() {
	if monk.Spec != proto.Spec_SpecWindwalkerMonk {
		return
	}

	tigerStrikesBuffAura(&monk.Unit)
}

func (pet *StormEarthAndFirePet) registerSEFTigerStrikes() {
	if pet.owner.Spec != proto.Spec_SpecWindwalkerMonk {
		return
	}

	tigerStrikesBuffAura(&pet.Unit)
}
