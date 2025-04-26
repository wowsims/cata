package windwalker

import (
	"time"

	"github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/monk"
)

func (ww *WindwalkerMonk) registerPassives() {
	ww.registerCombatConditioning()
	ww.registerComboBreaker()
	ww.registerTigerStrikes()
}

func (ww *WindwalkerMonk) registerCombatConditioning() {
	if !ww.HasMinorGlyph(proto.MonkMinorGlyph_GlyphOfBlackoutKick) && ww.PseudoStats.InFrontOfTarget {
		return
	}

	cata.RegisterIgniteEffect(&ww.Unit, cata.IgniteConfig{
		ActionID:           core.ActionID{SpellID: 100784}.WithTag(2), // actual 128531
		DotAuraLabel:       "Blackout Kick (DoT)" + ww.Label,
		DisableCastMetrics: true,
		IncludeAuraDelay:   true,
		SpellSchool:        core.SpellSchoolPhysical,
		NumberOfTicks:      4,
		TickLength:         time.Second,

		ProcTrigger: core.ProcTrigger{
			Name:           "Combat Conditioning" + ww.Label,
			Callback:       core.CallbackOnSpellHitDealt,
			ClassSpellMask: monk.MonkSpellBlackoutKick,
			Outcome:        core.OutcomeLanded,
		},

		DamageCalculator: func(result *core.SpellResult) float64 {
			return result.Damage * 0.2
		},
	})
}

func (ww *WindwalkerMonk) registerComboBreaker() {
	ww.ComboBreakerBlackoutKickAura = ww.RegisterAura(core.Aura{
		Label:    "Combo Breaker: Blackout Kick" + ww.Label,
		ActionID: core.ActionID{SpellID: 116768},
		Duration: time.Second * 20,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&monk.MonkSpellBlackoutKick == 0 || !result.Landed() {
				return
			}

			ww.ComboBreakerBlackoutKickAura.Deactivate(sim)
		},
	})

	ww.ComboBreakerTigerPalmAura = ww.RegisterAura(core.Aura{
		Label:    "Combo Breaker: Tiger Palm" + ww.Label,
		ActionID: core.ActionID{SpellID: 118864},
		Duration: time.Second * 20,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&monk.MonkSpellTigerPalm == 0 || !result.Landed() {
				return
			}

			ww.ComboBreakerTigerPalmAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&ww.Unit, core.ProcTrigger{
		Name:           "Combo Breaker: Blackout Kick Trigger" + ww.Label,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: monk.MonkSpellJab,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.12,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ww.ComboBreakerBlackoutKickAura.Activate(sim)
		},
	})

	core.MakeProcTriggerAura(&ww.Unit, core.ProcTrigger{
		Name:           "Combo Breaker: Tiger Palm Trigger" + ww.Label,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: monk.MonkSpellJab,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.12,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ww.ComboBreakerTigerPalmAura.Activate(sim)
		},
	})
}

func (ww *WindwalkerMonk) registerTigerStrikes() {
	tigerStrikesMHID := core.ActionID{SpellID: 120274}
	tigerStrikesOHID := core.ActionID{SpellID: 120278}

	var tigerStrikesMHSpell *core.Spell
	var tigerStrikesOHSpell *core.Spell
	var tigerStrikesBuff *core.Aura
	tigerStrikesBuff = ww.RegisterAura(core.Aura{
		Label:     "Tiger Strikes" + ww.Label,
		ActionID:  core.ActionID{SpellID: 120273},
		Duration:  time.Second * 15,
		MaxStacks: 4,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			mhConfig := *ww.AutoAttacks.MHConfig()
			mhConfig.ActionID = tigerStrikesMHID
			mhConfig.ClassSpellMask = monk.MonkSpellTigerStrikes
			mhConfig.Flags |= core.SpellFlagPassiveSpell
			mhConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhiteNoGlance)
			}
			tigerStrikesMHSpell = ww.GetOrRegisterSpell(mhConfig)

			if ww.HasOHWeapon() {
				ohConfig := *ww.AutoAttacks.OHConfig()
				ohConfig.ActionID = tigerStrikesOHID
				ohConfig.ClassSpellMask = monk.MonkSpellTigerStrikes
				ohConfig.Flags |= core.SpellFlagPassiveSpell
				ohConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhiteNoGlance)
				}
				tigerStrikesOHSpell = ww.GetOrRegisterSpell(ohConfig)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ww.MultiplyMeleeSpeed(sim, 1.5)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ww.MultiplyMeleeSpeed(sim, 1/1.5)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// TODO: Verify if it actually procs even on misses, seems so based on logs and simc
			if !tigerStrikesBuff.IsActive() /*|| !result.Landed()*/ || spell.ClassSpellMask&monk.MonkSpellTigerStrikes != 0 {
				return
			}

			if spell == ww.AutoAttacks.MHAuto() {
				tigerStrikesBuff.RemoveStack(sim)
				tigerStrikesMHSpell.Cast(sim, result.Target)
			} else if spell == ww.AutoAttacks.OHAuto() {
				tigerStrikesBuff.RemoveStack(sim)
				tigerStrikesOHSpell.Cast(sim, result.Target)
			}

			// TODO: Simc has delays but some SoO logs I looked at didn't...
			/*if spellToCast != nil {
				delaySeconds := sim.RollWithLabel(0.8, 1.2, "Tiger Strikes Delay")
				sim.AddPendingAction(&core.PendingAction{
					NextActionAt: sim.CurrentTime + core.DurationFromSeconds(delaySeconds),
					Priority:     core.ActionPriorityAuto,
					OnAction: func(sim *core.Simulation) {
						spellToCast.Cast(sim, result.Target)
					},
				})
			}*/
		},
	})

	core.MakeProcTriggerAura(&ww.Unit, core.ProcTrigger{
		Name:       "Tiger Strikes Buff Trigger" + ww.Label,
		ActionID:   core.ActionID{SpellID: 120272},
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskWhiteHit,
		ProcChance: 1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&monk.MonkSpellTigerStrikes != 0 {
				return
			}

			if sim.Proc(0.08, "Tiger Strikes") {
				tigerStrikesBuff.Activate(sim)
				tigerStrikesBuff.SetStacks(sim, 4)
			}
		},
	})
}
