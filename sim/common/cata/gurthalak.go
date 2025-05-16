package cata

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// Gurthalak, Voice of the Deeps
// Equip: Your melee attacks have a chance to cause you to summon a Tentacle of the Old Ones to fight by your side for 12 sec.
// (Proc chance: 2%)
func init() {
	// These spells ignore the slot the weapon is in.
	// Any other ability should only trigger the proc if the weapon is in the right slot.
	ignoresSlot := map[int32]bool{
		23881: true, // Bloodthirst
		6544:  true, // Heroic Leap
	}

	gurthalakItemIDs := []int32{78487, 77191, 78478}
	for _, v := range []ItemVersion{ItemVersionLFR, ItemVersionNormal, ItemVersionHeroic} {
		version := v // Gotta scope this for the closure
		labelSuffix := []string{" (LFR)", "", " (Heroic)"}[version]

		gurthalakItemID := gurthalakItemIDs[version]
		summonSpellID := []int32{109838, 107818, 109840}[version]

		core.NewItemEffect(gurthalakItemID, func(agent core.Agent) {
			var gurthalak GurthalakAgent
			if gurth, canEquip := agent.(GurthalakAgent); canEquip {
				gurthalak = gurth
			} else {
				return
			}

			character := agent.GetCharacter()
			meleeWeaponSlots := core.MeleeWeaponSlots()
			label := fmt.Sprintf("Gurthalak Trigger %s", labelSuffix)
			summonSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: summonSpellID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					tentacles := gurthalak.GetTentacles()

					for _, tentacle := range tentacles {
						if tentacle.IsActive() {
							continue
						}

						tentacle.ProccedItemVersion = version
						tentacle.EnableWithTimeout(sim, tentacle, time.Second*12)

						return
					}

					if sim.Log != nil {
						character.Log(sim, "No tentacles available for Gurthalak to proc, this is unreasonable.")
					}
				},
			})

			makeProcTrigger := func(character *core.Character, isMH bool) {
				itemSlot := core.Ternary(isMH, meleeWeaponSlots[:1], meleeWeaponSlots[1:])
				procMask := core.Ternary(isMH, core.ProcMaskMeleeMH|core.ProcMaskMeleeProc, core.ProcMaskMeleeOH)

				aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
					Name:     fmt.Sprintf("%s %s", label, core.Ternary(isMH, "MH", "OH")),
					Callback: core.CallbackOnSpellHitDealt,
					ProcMask: core.ProcMaskMelee,
					Outcome:  core.OutcomeLanded,
					Harmful:  true,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if _, ignore := ignoresSlot[spell.ActionID.SpellID]; !spell.ProcMask.Matches(procMask) && !ignore {
							return
						}

						if sim.Proc(0.02, label) {
							summonSpell.Cast(sim, result.Target)
						}
					},
				})

				character.ItemSwap.RegisterProcWithSlots(gurthalakItemID, aura, itemSlot)
			}

			if character.ItemSwap.CouldHaveItemEquippedInSlot(gurthalakItemID, proto.ItemSlot_ItemSlotMainHand) {
				makeProcTrigger(character, true)
			}
			if character.ItemSwap.CouldHaveItemEquippedInSlot(gurthalakItemID, proto.ItemSlot_ItemSlotOffHand) {
				makeProcTrigger(character, false)
			}
		})
	}

}

type TentacleOfTheOldOnesPet struct {
	core.Pet
	mindFlay           map[ItemVersion]*core.Spell
	ProccedItemVersion ItemVersion
}

type GurthalakAgent interface {
	GetTentacles() []*TentacleOfTheOldOnesPet
}

func NewTentacleOfTheOldOnesPet(character *core.Character) *TentacleOfTheOldOnesPet {
	pet := &TentacleOfTheOldOnesPet{
		Pet: core.NewPet(core.PetConfig{
			Name:  "Tentacle of the Old Ones",
			Owner: character,
			BaseStats: stats.Stats{
				stats.Stamina: 100,
			},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				return stats.Stats{
					stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
					stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
				}
			},
			EnabledOnStart: false,
			IsGuardian:     true,
		}),
	}

	pet.Pet.OnPetEnable = pet.enable
	pet.Pet.OnPetDisable = pet.disable
	pet.PseudoStats.DamageTakenMultiplier = 0

	return pet
}

func (pet *TentacleOfTheOldOnesPet) enable(sim *core.Simulation) {
	pet.SetStartAttackDelay(core.DurationFromSeconds(sim.RollWithLabel(1.0, core.BossGCD.Seconds(), "Tentacle Summon Delay")))

	// It also inherits the owner's spell hit and crit, on each cast of Mind Flay.
	// TODO: Verify this on the PTR
	pet.EnableDynamicStats(func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
		}
	})
}
func (pet *TentacleOfTheOldOnesPet) disable(sim *core.Simulation) {
	if pet.ChanneledDot != nil {
		pet.ChanneledDot.Deactivate(sim)
	}
}

func (pet *TentacleOfTheOldOnesPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *TentacleOfTheOldOnesPet) Initialize() {
	pet.registerMindFlay()
}

func (pet *TentacleOfTheOldOnesPet) registerMindFlay() {
	actionID := core.ActionID{SpellID: 52586}

	pet.mindFlay = make(map[ItemVersion]*core.Spell)
	for _, v := range []ItemVersion{ItemVersionLFR, ItemVersionNormal, ItemVersionHeroic} {
		version := v // Gotta scope this for the closure
		labelSuffix := []string{" (LFR)", "", " (Heroic)"}[version]

		baseTickDamage := []float64{9881, 11155, 12591}[version]

		pet.mindFlay[version] = pet.RegisterSpell(core.SpellConfig{
			ActionID:    actionID.WithTag(int32(version)),
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagChanneled,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Mind Flay " + labelSuffix,
				},
				NumberOfTicks:       3,
				TickLength:          time.Second,
				AffectedByCastSpeed: false,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, baseTickDamage)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				},
			},

			DamageMultiplier:         1,
			DamageMultiplierAdditive: 1,
			ThreatMultiplier:         1,
			CritMultiplier:           pet.DefaultCritMultiplier(),

			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return pet.ChanneledDot == nil
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)

				if result.Landed() {
					dot := spell.Dot(target)
					dot.Apply(sim)
					// Delaying two NPC GCDs minus the total cast time of Mind Flay (3s)
					pet.ExtendGCDUntil(sim, sim.CurrentTime+core.BossGCD*2)
				} else {
					spell.DealOutcome(sim, result)
					// No delay between a miss and the retry apparently
					pet.QueueSpell(sim, spell, target, sim.CurrentTime)
				}
			},
		})
	}
}

func (pet *TentacleOfTheOldOnesPet) Reset(sim *core.Simulation) {
}

func (pet *TentacleOfTheOldOnesPet) ExecuteCustomRotation(sim *core.Simulation) {
	if !pet.IsEnabled() {
		if pet.Unit.ChanneledDot != nil {
			pet.Unit.ChanneledDot.Deactivate(sim)
		}

		return
	}

	mindFlay := pet.mindFlay[pet.ProccedItemVersion]
	if !pet.GCD.IsReady(sim) || !mindFlay.CanCast(sim, pet.Owner.CurrentTarget) {
		return
	}

	pet.PseudoStats.DamageDealtMultiplier = pet.Owner.PseudoStats.DamageDealtMultiplier
	pet.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] = pet.Owner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow]

	mindFlay.Cast(sim, pet.Owner.CurrentTarget)
}
