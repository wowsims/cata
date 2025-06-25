package bwd

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// Log used for fitting damage parameters: https://classic.warcraftlogs.com/reports/NTgLfqc2atyFh8BX#fight=26&type=damage-taken&target=325&view=events&pins=0%24Off%24%23244F4B%24auras-gained%241%240.0.0.Any%240.0.0.Any%24true%240.0.0.Any%24true%2479330%2463%5E2%24Off%24%23909049%24auras-gained%241%240.0.0.Any%240.0.0.Any%24true%240.0.0.Any%24true%241160%24true%24true%2495%24and%24auras-gained%241%240.0.0.Any%240.0.0.Any%24true%240.0.0.Any%24true%246343%24true%24true%2495
// Assumes that this logging bug with Demo Shout is still present: https://github.com/JamminL/cata-classic-bugs/issues/1163

func addNefarian(bossPrefix string) {
	// TODO: Add support for 10-man and Normal variants
	createNefarianPreset(bossPrefix, 25, true, 41918, 38_745_000, 2394)
}

func createNefarianPreset(bossPrefix string, raidSize int, isHeroic bool, addNpcId int32, addHealth float64, addMinBaseDamage float64) {
	// TODO: Add support for tanking boss instead of adds
	targetName := fmt.Sprintf("Nefarian %d", raidSize)
	targetNameAdd := fmt.Sprintf("Animated Bone Warrior %d", raidSize)

	if isHeroic {
		targetName += " H"
		targetNameAdd += " H"
	}

	var targetPathNames []string

	for addIdx := int32(1); addIdx <= 12; addIdx++ {
		currentAddName := targetNameAdd + fmt.Sprintf(" - %d", addIdx)
		targetInputs := []*proto.TargetInput{}

		if addIdx == 1 {
			targetInputs = append(targetInputs, &proto.TargetInput{
				Label:       "Electrocute Count",
				Tooltip:     "Number of Electrocute casts to model. Total count will be spread evenly over the encounter duration with a randomized offset.",
				InputType:   proto.InputType_Number,
				NumberValue: 6,
			})
		}

		core.AddPresetTarget(&core.PresetTarget{
			PathPrefix: bossPrefix,

			Config: &proto.Target{
				Id:        addNpcId*100 + addIdx, // hack to guarantee distinct IDs for each add
				Name:      currentAddName,
				Level:     85,
				MobType:   proto.MobType_MobTypeUndead,
				TankIndex: 0, // change if boss tanking support is added

				Stats: stats.Stats{
					stats.Health:      addHealth,
					stats.Armor:       11977, // TODO: verify add armor
					stats.AttackPower: 0,     // actual value doesn't matter in Cata, as long as damage parameters are fit consistently
				}.ToProtoArray(),

				SpellSchool:   proto.SpellSchool_SpellSchoolPhysical,
				SwingSpeed:    2.0,
				MinBaseDamage: addMinBaseDamage,
				DamageSpread:  0.34,
				TargetInputs:  targetInputs,
			},

			AI: makeNefarianAddAI(raidSize, isHeroic, addIdx),
		})

		targetPathNames = append(targetPathNames, bossPrefix+"/"+currentAddName)
	}

	core.AddPresetEncounter(targetName+" Adds", targetPathNames)
}

func makeNefarianAddAI(raidSize int, isHeroic bool, addIdx int32) core.AIFactory {
	return func() core.TargetAI {
		return &NefarianAddAI{
			raidSize: raidSize,
			isHeroic: isHeroic,
			addIdx:   addIdx,
		}
	}
}

type NefarianAddAI struct {
	Target *core.Target

	raidSize int
	isHeroic bool
	addIdx   int32

	empowerAura      *core.Aura
	shadowblazeSpark *core.Spell

	isController     bool // designate one "add" to cast raid-wide mechanics
	numElectrocutes  int32
	electrocuteSpell *core.Spell
}

func (ai *NefarianAddAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.Target.AutoAttacks.MHConfig().ActionID.Tag = 4191800 + ai.addIdx // hack for UI results parsing
	ai.isController = (ai.addIdx == 1)

	if ai.isController {
		ai.numElectrocutes = int32(config.TargetInputs[0].NumberValue)
	}

	ai.registerSpells()
}

func (ai *NefarianAddAI) Reset(sim *core.Simulation) {
}

func (ai *NefarianAddAI) registerSpells() {
	// Empower Aura
	empowerDamageMod := core.TernaryFloat64(ai.isHeroic, 2.0, 1.0)
	empowerActionID := core.ActionID{SpellID: 79330}
	ai.empowerAura = ai.Target.GetOrRegisterAura(core.Aura{
		Label:     "Empower",
		ActionID:  empowerActionID,
		MaxStacks: 13,
		Duration:  time.Second * 52,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= (1.0 + empowerDamageMod*float64(newStacks)) / (1.0 + empowerDamageMod*float64(oldStacks))
		},

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, 1)
			aura.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime-aura.Unit.AutoAttacks.MainhandSwingSpeed()+1)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 4,
				NumTicks: 12,
				Priority: core.ActionPriorityDOT,

				OnAction: func(sim *core.Simulation) {
					aura.AddStack(sim)
				},
			})
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+time.Second*26)
		},
	})

	// Add re-activation via Shadowblaze Spark
	if !ai.isController {
		return
	}

	ai.shadowblazeSpark = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 81031},
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 26,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, addUnit := range sim.Encounter.TargetUnits {
				empowerAura := addUnit.GetAuraByID(empowerActionID)

				// Assume that the tank is always pre-moving adds before the spark hits them, so that Empower is never refreshed on already active adds.
				if (empowerAura != nil) && !empowerAura.IsActive() {
					empowerAura.Activate(sim)
				}
			}
		},
	})

	// Electrocute raid mechanic
	electrocuteBase := core.TernaryFloat64(ai.isHeroic, 128700, 72765)
	electrocuteVariance := core.TernaryFloat64(ai.isHeroic, 2600, 4470)
	ai.electrocuteSpell = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 81272},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagAPL,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: 1, // Placeholder value, will be set on reset
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Raid.AllPlayerUnits {
				damageRoll := electrocuteBase + electrocuteVariance*sim.RandomFloat("Electrocute Damage")
				spell.CalcAndDealDamage(sim, aoeTarget, damageRoll, spell.OutcomeAlwaysHit)
			}
		},
	})

	ai.Target.RegisterResetEffect(func(sim *core.Simulation) {
		// Randomize Shadowblaze timer to desync from de-activation timer
		ai.shadowblazeSpark.CD.Set(core.DurationFromSeconds(sim.RandomFloat("Shadowblaze Timing") * ai.shadowblazeSpark.CD.Duration.Seconds()))

		// Set a "cooldown" for Electrocute to match user input
		ai.electrocuteSpell.CD.Duration = sim.Duration/time.Duration(ai.numElectrocutes) - BossGCD/time.Duration(2)
		ai.electrocuteSpell.CD.Set(core.DurationFromSeconds(sim.RandomFloat("Electrocute Timing") * ai.electrocuteSpell.CD.Duration.Seconds()))
	})
}

func (ai *NefarianAddAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget
	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if ai.isController && ai.electrocuteSpell.IsReady(sim) {
		ai.electrocuteSpell.Cast(sim, target)
	}

	if ai.isController && ai.shadowblazeSpark.IsReady(sim) {
		ai.shadowblazeSpark.Cast(sim, target)
	}

	ai.Target.ExtendGCDUntil(sim, sim.CurrentTime+BossGCD)
}
