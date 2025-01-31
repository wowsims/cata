package cata

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
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
		6343:  true, // Thunder Clap
	}

	// TODO: Check if there's any restrictions on certain off-hand attacks like there are for Apparatus and Vessel.

	for _, v := range []ItemVersion{ItemVersionLFR, ItemVersionNormal, ItemVersionHeroic} {
		version := v // Gotta scope this for the closure
		labelSuffix := []string{" (LFR)", "", " (Heroic)"}[version]

		gurthalakItemID := []int32{78487, 77191, 78478}[version]
		core.NewItemEffect(gurthalakItemID, func(agent core.Agent) {
			var gurthalak GurthalakAgent
			if gurth, canEquip := agent.(GurthalakAgent); canEquip {
				gurthalak = gurth
			} else {
				return
			}

			character := agent.GetCharacter()

			summonSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 109840},
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

			procMask := character.GetDynamicProcMaskForWeaponEffect(gurthalakItemID)
			aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     "Gurthalak Trigger" + labelSuffix,
				ActionID: core.ActionID{ItemID: gurthalakItemID},
				Callback: core.CallbackOnSpellHitDealt,
				ProcMask: core.ProcMaskMelee, // TODO: Verify correct proc masks on the PTR
				Outcome:  core.OutcomeLanded,
				Harmful:  true,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if _, ignore := ignoresSlot[spell.ActionID.SpellID]; !spell.ProcMask.Matches(*procMask) && !ignore {
						return
					}

					if !sim.Proc(0.02, "Gurthalak, Voice of the Deeps"+labelSuffix) {
						return
					}

					if sim.Log != nil {
						slotLabel := core.Ternary(spell.IsMH(), "Main Hand", "Off Hand")
						character.Log(sim, "Gurthalak (%s) procced by %s", slotLabel, spell.ActionID)
					}

					summonSpell.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(gurthalakItemID, aura)
		})
	}
}

type TentacleOfTheOldOnesPet struct {
	core.Pet
	mindFlay           map[ItemVersion]*core.Spell
	castsLeft          int32
	castDelay          time.Duration
	ProccedItemVersion ItemVersion
}

type GurthalakAgent interface {
	GetTentacles() []*TentacleOfTheOldOnesPet
}

func NewTentacleOfTheOldOnesPet(character *core.Character) *TentacleOfTheOldOnesPet {
	pet := &TentacleOfTheOldOnesPet{
		Pet: core.NewPet("Tentacle of the Old Ones", character, stats.Stats{
			stats.Stamina: 100,
		}, func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{
				stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
				stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
			}
		}, false, true),
	}

	pet.Pet.OnPetEnable = pet.enable
	pet.PseudoStats.DamageTakenMultiplier = 0

	return pet
}

func (pet *TentacleOfTheOldOnesPet) enable(sim *core.Simulation) {
	pet.castDelay = sim.CurrentTime + time.Second
	pet.castsLeft = 3

	// The tentacles inherit the owner's damage dealt multipliers (Avenging Wrath, Berserker Stance, Communion etc.)
	// as well as Shadow School damage multipliers like Unholy DKs Dreadblade mastery.
	// TODO: Verify this on the PTR
	pet.PseudoStats.DamageDealtMultiplier = pet.Owner.PseudoStats.DamageDealtMultiplier
	pet.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] = pet.Owner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow]

	// It also inherits the owner's spell hit and crit, on each cast of Mind Flay.
	// TODO: Verify this on the PTR
	pet.EnableDynamicStats(func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
		}
	})
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

		// TODO: Verify these damage numbers on the PTR, they do match all of the videos on YouTube at least.
		baseTickDamage := []float64{9881, 11155, 12591}[version]

		pet.mindFlay[version] = pet.RegisterSpell(core.SpellConfig{
			ActionID:    actionID.WithTag(int32(version)),
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskSpellDamage,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagChanneled,

			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
			},

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

			DamageMultiplier: 1,

			ThreatMultiplier: 1,
			CritMultiplier:   pet.DefaultSpellCritMultiplier(),

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				pet.castsLeft--

				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
				dot := spell.Dot(target)
				// The last cast had a 50/50 chance of only getting 2 ticks back in the day.
				// This could've been because of aura delays or something not present in Classic.
				// TODO: Verify this on the PTR
				if pet.castsLeft == 0 {
					dot.BaseTickCount = []int32{2, 3}[int(sim.RollWithLabel(0, 2, "Tentacle of the Old Ones"+labelSuffix))]
				} else {
					dot.BaseTickCount = 3
				}

				if result.Landed() {
					dot.Apply(sim)
					// In old videos there seems to be a delay between each cast, we simulate this by setting an explicit wait time.
					// TODO: Verify this on the PTR
					pet.castDelay = sim.CurrentTime + time.Second*4
				}
			},
		})
	}
}

func (pet *TentacleOfTheOldOnesPet) Reset(sim *core.Simulation) {
	pet.castDelay = sim.CurrentTime
	pet.castsLeft = 3
}

func (pet *TentacleOfTheOldOnesPet) ExecuteCustomRotation(sim *core.Simulation) {
	if sim.CurrentTime < pet.castDelay {
		pet.WaitUntil(sim, pet.castDelay)
		return
	}

	mindFlay := pet.mindFlay[pet.ProccedItemVersion]
	if !pet.GCD.IsReady(sim) || pet.Unit.ChanneledDot != nil && pet.Unit.ChanneledDot.Spell == mindFlay {
		return
	}

	if pet.castsLeft <= 0 {
		if pet.IsEnabled() {
			pet.Disable(sim)
		}
		return
	}

	// These should update at the time of cast. Popping wings mid spawn should apply on the next cast for example.
	// TODO: Verify this on the PTR
	pet.PseudoStats.DamageDealtMultiplier = pet.Owner.PseudoStats.DamageDealtMultiplier
	pet.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] = pet.Owner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow]

	mindFlay.Cast(sim, pet.CurrentTarget)
}
