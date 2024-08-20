package rogue

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

var Tier11 = core.NewItemSet(core.ItemSet{
	Name: "Wind Dancer's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// +5% Crit to Backstab, Mutilate, and Sinister Strike
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 5,
				ClassMask:  RogueSpellBackstab | RogueSpellMutilate | RogueSpellSinisterStrike,
			})
		},
		4: func(agent core.Agent) {
			// 1% Chance on Auto Attack to increase crit of next Evis or Envenom by +100% for 15 seconds
			rogue := agent.(RogueAgent).GetRogue()

			t11Proc := rogue.RegisterAura(core.Aura{
				Label:    "Deadly Scheme Proc",
				ActionID: core.ActionID{SpellID: 90472},
				Duration: time.Second * 15,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.Envenom.BonusCritPercent += 100
					rogue.Eviscerate.BonusCritPercent += 100
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.Envenom.BonusCritPercent -= 100
					rogue.Eviscerate.BonusCritPercent -= 100
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == rogue.Envenom || spell == rogue.Eviscerate {
						aura.Deactivate(sim)
					}
				},
			})

			core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
				Name:       "Deadly Scheme Aura",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.01,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					t11Proc.Activate(sim)
				},
			})
		},
	},
})

func MakeT12StatAura(action core.ActionID, stat stats.Stat, name string) core.Aura {
	var lastRatingGain float64
	return core.Aura{
		Label:    name,
		ActionID: action,
		Duration: 30 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			lastRatingGain = aura.Unit.GetStat(stat) * 0.25
			aura.Unit.AddStatDynamic(sim, stat, lastRatingGain)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stat, -lastRatingGain)
		},
	}
}

var Tier12 = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Dark Phoenix",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your melee critical strikes deal 6% additional damage as Fire over 4 sec.
			// Rolls like ignite
			// Tentatively, this is just Ignite. Testing required to validate behavior.
			rogue := agent.GetCharacter()

			burningWounds := rogue.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 99173},
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskProc,
				Flags:       core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label:     "Burning Wounds",
						MaxStacks: 1000000,
					},
					NumberOfTicks:       2,
					TickLength:          time.Second * 2,
					AffectedByCastSpeed: false,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {

					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						result := dot.Spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
						dot.Spell.DealPeriodicDamage(sim, result)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).Apply(sim)
				},
			})

			var applyBurningWounds = func(sim *core.Simulation, result *core.SpellResult) {
				const IgniteTicksFresh = 2
				const IgniteTicksRefresh = 3

				newDamage := result.Damage * .06
				dot := burningWounds.Dot(result.Target)

				// Cata Ignite
				// 1st ignite application = 4s, split into 2 ticks (2s, 0s)
				// Ignite refreshes: Duration = 4s + MODULO(remaining duration, 2), max 6s. Split damage over 3 ticks at 4s, 2s, 0s.
				if dot.IsActive() {
					dot.SnapshotBaseDamage = (dot.OutstandingDmg() + newDamage) / float64(IgniteTicksRefresh)
				} else {
					dot.SnapshotBaseDamage = newDamage / IgniteTicksFresh
				}
				burningWounds.Cast(sim, result.Target)
				dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
			}

			core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
				Name:     "Rogue T12 2P Bonus",
				Callback: core.CallbackOnSpellHitDealt,
				ProcMask: core.ProcMaskMelee,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && result.DidCrit() {
						applyBurningWounds(sim, result)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// Your Tricks of the Trade ability also causes you to gain a 25% increase to Haste, Mastery, or Critical Strike chosen at random for 30 sec.
			// Cannot pick the same stat twice in a row. No other logic appears to exist
			// Not a dynamic 1.25% mod; snapshots stats and applies that much as bonus rating for duration
			// Links to all buffs: https://www.wowhead.com/spell=99175/item-rogue-t12-4p-bonus#comments:id=1507073
			rogue := agent.GetCharacter()

			// Aura for adding 25% of current rating as extra rating
			hasteAura := rogue.GetOrRegisterAura(MakeT12StatAura(core.ActionID{SpellID: 99186}, stats.HasteRating, "Future on Fire"))
			critAura := rogue.GetOrRegisterAura(MakeT12StatAura(core.ActionID{SpellID: 99187}, stats.CritRating, "Fiery Devastation"))
			mastAura := rogue.GetOrRegisterAura(MakeT12StatAura(core.ActionID{SpellID: 99188}, stats.MasteryRating, "Master of the Flames"))
			auraArray := [3]*core.Aura{hasteAura, critAura, mastAura}

			// Proc aura watching for ToT threat transfer start
			lastStat := 3
			core.MakeProcTriggerAura(&agent.GetCharacter().Unit, core.ProcTrigger{
				Name:           "Rogue T12 4P Bonus",
				Callback:       core.CallbackOnApplyEffects,
				ClassSpellMask: RogueSpellTricksOfTheTradeThreat,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if lastStat == 3 { // any of 3
						randomStat := int(math.Mod(sim.RandomFloat("Rogue T12 4P Bonus")*10, 3))
						lastStat = randomStat
						auraArray[lastStat].Activate(sim)
					} else { // cannot re-roll same
						randomStat := int(math.Mod(sim.RandomFloat("Rogue T12 4P Bonus")*10, 1)) + 1
						lastStat = (lastStat + randomStat) % 3
						auraArray[lastStat].Activate(sim)
					}
				},
			})
		},
	},
})

var CataPVPSet = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Vestments",
	ID:   914,
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Agility, 70)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Agility, 90)
			// 10 maximum energy added in rogue.go
		},
	},
})
