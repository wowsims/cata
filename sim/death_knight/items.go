package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// TODO: T13

// T11 - DPS
var ItemSetMagmaPlatedBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Magma Plated Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of your Death Coil and Frost Strike abilities by 5%.
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  DeathKnightSpellDeathCoil | DeathKnightSpellFrostStrike,
				FloatValue: 5,
			})
		},
		4: func(agent core.Agent) {
			// Each time you gain a Death Rune, you also gain 1% increased attack power for 30 sec. Stacks up to 3 times.
			// Also activated whenever KM procs
			character := agent.GetCharacter()

			apDep := make([]*stats.StatDependency, 3)
			for i := 1; i <= 3; i++ {
				apDep[i-1] = character.NewDynamicMultiplyStat(stats.AttackPower, 1.0+float64(i)*0.01)
			}

			aura := character.GetOrRegisterAura(core.Aura{
				Label:     "Death Eater",
				ActionID:  core.ActionID{SpellID: 90507},
				Duration:  time.Second * 30,
				MaxStacks: 3,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if oldStacks > 0 {
						character.DisableDynamicStatDep(sim, apDep[oldStacks-1])
					}
					if newStacks > 0 {
						character.EnableDynamicStatDep(sim, apDep[newStacks-1])
					}
				},
			})

			core.MakeProcTriggerAura(&agent.GetCharacter().Unit, core.ProcTrigger{
				Name:           "Magma Plated Battlegear",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: DeathKnightSpellConvertToDeathRune | DeathKnightSpellKillingMachine,
				ICD:            time.Millisecond * 10, // Batch together double rune converts
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
					aura.AddStack(sim)
				},
			})
		},
	},
})

// T11 - Tank
var ItemSetMagmaPlatedBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Magma Plated Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage done by your Death Strike ability by 5%.
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  DeathKnightSpellDeathStrike,
				FloatValue: 0.05,
			})
		},
		4: func(agent core.Agent) {
			// Increases the duration of your Icebound Fortitude ability by 50%.
			// Implemented in icebound_fortitude.go
		},
	},
})

// T12 - DPS
var ItemSetElementiumDeathplateBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Elementium Deathplate Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			actionID := core.ActionID{SpellID: 98971}
			rpMetrics := dk.NewRunicPowerMetrics(actionID)
			var pa *core.PendingAction

			buff := dk.RegisterAura(core.Aura{
				Label:    "Smoldering Rune",
				ActionID: actionID,
				Duration: time.Minute * 2,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
						Period: time.Second * 5,
						OnAction: func(sim *core.Simulation) {
							// Make sure the multiplier is always 1 as this effect doesn't seem to scale
							runicMulti := dk.GetRunicRegenMultiplier()
							dk.MultiplyRunicRegen(1 / runicMulti)
							dk.AddRunicPower(sim, 3, rpMetrics)
							dk.MultiplyRunicRegen(runicMulti)
						},
					})
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					pa.Cancel(sim)
				},
			})

			core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
				Name:           "Smolering Rune Trigger",
				ActionID:       actionID,
				ClassSpellMask: DeathKnightSpellHornOfWinter,
				Callback:       core.CallbackOnCastComplete,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					buff.Activate(sim)
				},
			})

		},
		4: func(agent core.Agent) {

			dk := agent.(DeathKnightAgent).GetDeathKnight()
			damage := 0.0

			newFlamingTormentSpell := func(spellID int32) core.SpellConfig {
				actionID := core.ActionID{SpellID: spellID} // actually 99000

				return core.SpellConfig{
					ActionID:    actionID.WithTag(3),
					SpellSchool: core.SpellSchoolFire,
					ProcMask:    core.ProcMaskEmpty,
					Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagIgnoreModifiers | core.SpellFlagPassiveSpell,

					DamageMultiplier: 1,
					ThreatMultiplier: 1,

					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						spell.CalcAndDealDamage(sim, spell.Unit.CurrentTarget, damage, spell.OutcomeAlwaysHit)
					},
				}
			}

			var flamingTormentSpellForObliterate = dk.RegisterSpell(newFlamingTormentSpell(49020))
			var flamingTormentSpellForScourgeStrike = dk.RegisterSpell(newFlamingTormentSpell(55090))

			core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
				Name:           "Flaming Torment Trigger",
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: DeathKnightSpellObliterate | DeathKnightSpellScourgeStrike | DeathKnightSpellScourgeStrikeShadow,
				Outcome:        core.OutcomeLanded,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					damage = result.Damage * 0.06
					if spell.ClassSpellMask == DeathKnightSpellObliterate {
						flamingTormentSpellForObliterate.Cast(sim, result.Target)
					} else {
						flamingTormentSpellForScourgeStrike.Cast(sim, result.Target)
					}
				},
			})

		},
	},
})

// T12 - Tank
var ItemSetElementiumDeathplateBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Elementium Deathplate Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {

			dk := agent.(DeathKnightAgent).GetDeathKnight()

			dk.BurningBloodSpell = dk.RegisterSpell(core.SpellConfig{
				ActionID:         core.ActionID{SpellID: 98957},
				SpellSchool:      core.SpellSchoolFire,
				Flags:            core.SpellFlagIgnoreModifiers | core.SpellFlagAPL | core.SpellFlagPassiveSpell,
				ProcMask:         core.ProcMaskEmpty,
				DamageMultiplier: 1,
				CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Burning Blood" + dk.Label,
					},
					NumberOfTicks: 3,
					TickLength:    time.Second * 2,

					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
						baseDamage := 800.0
						if target.HasActiveAuraWithTag(core.SpellDamageEffectAuraTag) {
							baseDamage *= 1.08
						}
						dot.Snapshot(target, baseDamage)
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).Apply(sim)
				},
			})

			core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
				Name:       "Burning Blood Trigger",
				ActionID:   core.ActionID{SpellID: 98956},
				ProcMask:   core.ProcMaskMelee,
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcChance: 1,
				Harmful:    true,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					dk.BurningBloodSpell.Cast(sim, result.Target)
				},
			})

		},
		4: func(agent core.Agent) {
			// When your Dancing Rune Weapon expires, you gain 15% additional parry chance for 12 sec.
			// Implemented in dancing_rune_weapon.go
		},
	},
})

func init() {
}
