package hunter

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *Hunter) newFlamingArrowSpell(spellID int32) core.SpellConfig {
	actionID := core.ActionID{SpellID: spellID} // actually 99058

	return core.SpellConfig{
		ActionID:    actionID.WithTag(3),
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 0.8,
		CritMultiplier:   hunter.CritMultiplier(false, false, false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower(target))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	}
}

var ItemSetFlameWakersBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Flamewaker's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()

			var flamingArrowSpellForSteadyShot = hunter.RegisterSpell(hunter.newFlamingArrowSpell(56641))
			var flamingArrowSpellForCobraShot = hunter.RegisterSpell(hunter.newFlamingArrowSpell(77767))

			hunter.RegisterAura(core.Aura{
				Label:    "T12 2-set",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell != hunter.SteadyShot && spell != hunter.CobraShot {
						return
					}
					procChance := 0.1
					if sim.RandomFloat("Flaming Arrow") < procChance {
						fmt.Printf("Flaming Arrow procced on %s\n", spell.ActionID.SpellID)
						if spell == hunter.SteadyShot {
							flamingArrowSpellForSteadyShot.Cast(sim, result.Target)
						} else {
							flamingArrowSpellForCobraShot.Cast(sim, result.Target)
						}
					}
				},
			})
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			var baMod = hunter.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				ClassMask:  HunterSpellsAll,
				FloatValue: -1,
			})
			var burningAdrenaline = hunter.RegisterAura(core.Aura{
				Label:    "Burning Adrenaline",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 99060},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					baMod.Activate()
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ClassSpellMask&HunterSpellsAll != spell.ClassSpellMask {
						return
					}
					baMod.Deactivate()
					aura.Deactivate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					baMod.Deactivate()
				},
			})
			hunter.RegisterAura(core.Aura{
				Label:    "T12 4-set",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell != hunter.AutoAttacks.RangedAuto() {
						return
					}
					procChance := 0.1
					if sim.RandomFloat("Burning Adrenaline") < procChance {
						burningAdrenaline.Activate(sim)
					}
				},
			})
		},
	},
})
var ItemSetWyrmstalkerBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Wyrmstalker Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Handled in Cobra and Steady code respectively
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			var chronoHunter = hunter.RegisterAura(core.Aura{ // 105919
				Label:    "Chronohunter",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 105919},
				Icd: &core.Cooldown{ // Assume 113sec icd for now
					Duration: time.Second * 113,
					Timer:    hunter.NewTimer(),
				},
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyRangedSpeed(sim, 1.3)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.MultiplyRangedSpeed(sim, 1/1.3)
				},
			})
			hunter.RegisterAura(core.Aura{
				Label:    "T13 4-set",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell != hunter.ArcaneShot {
						return
					}
					procChance := 0.4
					if sim.RandomFloat("Chronohunter") < procChance {
						chronoHunter.Activate(sim)
					}
				},
			})
		},
	},
})
var ItemSetLightningChargedBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Lightning-Charged Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// 5% Crit on SS
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  HunterSpellSerpentSting,
				FloatValue: 5,
			})
		},
		4: func(agent core.Agent) {
			// Cobra & Steady Shot < 0.2s cast time
			// Cannot be spell modded for now
		},
	},
})

var ItemSetGladiatorsPursuit = core.NewItemSet(core.ItemSet{
	ID:   920,
	Name: "Gladiator's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.AddStats(stats.Stats{
				stats.Agility: 70,
			})
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			// Multiply focus regen 1.05
			hunter.AddStats(stats.Stats{
				stats.Agility: 90,
			})
		},
	},
})

func (hunter *Hunter) addBloodthirstyGloves() {
	switch hunter.Hands().ID {
	case 64991, 64709, 60424, 65544, 70534, 70260, 70441, 72369, 73717, 73583:
		hunter.AddStaticMod(core.SpellModConfig{
			ClassMask: HunterSpellExplosiveTrap | HunterSpellBlackArrow,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: -time.Second * 2,
		})
	default:
		break
	}
}
