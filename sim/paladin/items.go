package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// Tier 11 ret
var ItemSetReinforcedSapphiriumBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battleplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  SpellMaskTemplarsVerdict,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent) {
			// Handled in inquisition.go
		},
	},
})

// Tier 12 ret
var ItemSetBattleplateOfImmolation = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Immolation",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			flamesOfTheFaithful := paladin.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 99092},
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,
				Flags: core.SpellFlagIgnoreModifiers |
					core.SpellFlagBinary |
					core.SpellFlagNoOnCastComplete |
					core.SpellFlagNoOnDamageDealt,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label:    "Flames of the Faithful",
						Duration: time.Second * 4,
					},
					NumberOfTicks:       2,
					AffectedByCastSpeed: false,
					TickLength:          2 * time.Second,

					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
						dot.SnapshotAttackerMultiplier = 1
						dot.SnapshotCritChance = 0
					},

					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).Apply(sim)
				},
			})

			core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
				Name:           "Flames of the Faithful",
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskCrusaderStrike,
				Outcome:        core.OutcomeLanded,

				ProcChance: 1,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					dot := flamesOfTheFaithful.Dot(result.Target)

					outstandingDamage := dot.OutstandingDmg()
					newDamage := result.Damage * 0.15
					totalDamage := outstandingDamage + newDamage

					waitTime := time.Millisecond * time.Duration(sim.Roll(375, 625))
					applyDotAt := sim.CurrentTime + waitTime

					if sim.Log != nil {
						paladin.Log(sim, "Schedule travel (%0.2f) for Flames of the Faithful", waitTime.Seconds())
						if dot.IsActive() && dot.NextTickAt() < applyDotAt {
							paladin.Log(sim, "Flames of the Faithful rolled with %0.3f damage both ticking and rolled into next", outstandingDamage)
						}
					}

					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt:     applyDotAt,
						Priority: core.ActionPriorityDOT,
						OnAction: func(simulation *core.Simulation) {
							ticks := float64(dot.BaseTickCount + core.TernaryInt32(dot.IsActive(), 1, 0))
							dot.SnapshotBaseDamage = totalDamage / ticks

							flamesOfTheFaithful.Cast(sim, result.Target)
						},
					})
				},
			})
		},
		4: func(agent core.Agent) {
			// Handled in talents_retribution.go
		},
	},
})

// Tier 13 ret
var ItemSetBattleplateOfRadiantGlory = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Radiant Glory",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Handled in judgement.go
		},
		4: func(agent core.Agent) {
			// Handled in talents_retribution.go
		},
	},
})

// PvP set
var ItemSetGladiatorsVindication = core.NewItemSet(core.ItemSet{
	ID:   917,
	Name: "Gladiator's Vindication",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStat(stats.Strength, 70)
		},
		4: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStat(stats.Strength, 90)
			paladin.AddStaticMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: SpellMaskJudgementBase,
				TimeValue: -1 * time.Second,
			})
		},
	},
})

func (paladin *Paladin) addBloodthirstyGloves() {
	switch paladin.Hands().ID {
	case 64844, 70649, 60414, 65591, 72379, 70250, 70488, 73707, 73570:
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: 0.05,
		})
	default:
		break
	}
}

// Tier 11 prot
var ItemSetReinforcedSapphiriumBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  SpellMaskCrusaderStrike,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent) {
			// Handled in guardian_of_ancient_kings.go
		},
	},
})

// Tier 12 prot
var ItemSetBattlearmorOfImmolation = core.NewItemSet(core.ItemSet{
	Name: "Battlearmor of Immolation",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			procDamage := 0.0

			righteousFlames := paladin.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 99075},
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,
				Flags: core.SpellFlagIgnoreModifiers |
					core.SpellFlagBinary |
					core.SpellFlagNoOnCastComplete |
					core.SpellFlagNoOnDamageDealt,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, procDamage, spell.OutcomeAlwaysHit)
				},
			})

			core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
				Name:           "Righteous Flames",
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskShieldOfTheRighteous,
				Outcome:        core.OutcomeLanded,

				ProcChance: 1,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procDamage = result.Damage * 0.2
					righteousFlames.Cast(sim, result.Target)
				},
			})
		},
		4: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.FlamingAegis = paladin.GetOrRegisterAura(core.Aura{
				Label:    "Flaming Aegis",
				ActionID: core.ActionID{SpellID: 99090},
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					paladin.AddStatDynamic(sim, stats.ParryRating, 12*core.ParryRatingPerParryPercent)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					paladin.AddStatDynamic(sim, stats.ParryRating, -12*core.ParryRatingPerParryPercent)
				},
			})

			// Trigger in divine_protection.go
		},
	},
})
