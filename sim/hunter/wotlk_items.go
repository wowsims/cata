package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/common/wotlk"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var ItemSetAhnKaharBloodHuntersBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Ahn'Kahar Blood Hunter's Battlegear",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			hunter := agent.(HunterAgent).GetHunter()

			procAura := hunter.RegisterAura(core.Aura{
				Label:    "AhnKahar 2pc Proc",
				ActionID: core.ActionID{SpellID: 70727},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.15
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.15
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "AhnKahar 2pc",
				ProcChance: 0.05,
				ProcMask:   core.ProcMaskRangedAuto,
				Callback:   core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			hunter := agent.(HunterAgent).GetHunter()

			var curBonus stats.Stats
			procAura := hunter.RegisterAura(core.Aura{
				Label:    "AhnKahar 4pc Proc",
				ActionID: core.ActionID{SpellID: 70730},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					curBonus = stats.Stats{
						stats.AttackPower:       aura.Unit.GetStat(stats.AttackPower) * 0.1,
						stats.RangedAttackPower: aura.Unit.GetStat(stats.RangedAttackPower) * 0.1,
					}

					aura.Unit.AddStatsDynamic(sim, curBonus)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatsDynamic(sim, curBonus.Invert())
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "AhnKahar 4pc",
				ProcChance:     0.05,
				ClassSpellMask: HunterSpellSerpentSting,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
	},
})

func init() {

	wotlk.NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Zod's Repeating Longbow"
		itemID := int32(50034)
		procChance := 0.04
		if isHeroic {
			name += " H"
			itemID = 50638
			procChance = 0.05
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			if agent.GetCharacter().Class != proto.Class_ClassHunter {
				return
			}

			hunter := agent.(HunterAgent).GetHunter()

			var rangedSpell *core.Spell
			initSpell := func() {
				rangedSpell = hunter.RegisterSpell(core.SpellConfig{
					ActionID:    core.ActionID{ItemID: itemID},
					SpellSchool: core.SpellSchoolPhysical,
					ProcMask:    core.ProcMaskRangedAuto,
					Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

					DamageMultiplier: 0.5,
					CritMultiplier:   hunter.AutoAttacks.RangedConfig().CritMultiplier,
					ThreatMultiplier: 1,

					BonusCoefficient: 1,

					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						baseDamage := hunter.RangedWeaponDamage(sim, spell.RangedAttackPower(target))

						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
					},
				})
			}

			triggerAura := core.MakeProcTriggerAura(&hunter.Unit, core.ProcTrigger{
				Name:       name + " Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskRanged,
				Outcome:    core.OutcomeLanded,
				ProcChance: procChance,
				ActionID:   core.ActionID{ItemID: itemID},
				Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
					rangedSpell.Cast(sim, result.Target)
				},
			})
			triggerAura.OnInit = func(aura *core.Aura, sim *core.Simulation) {
				initSpell()
			}

		})
	})

}
