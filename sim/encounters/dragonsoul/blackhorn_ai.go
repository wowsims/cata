package dragonsoul

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const blackhornMeleeDamageSpread = 0.1
const blackhornID int32 = 56427
const gorionaID int32 = 56781

func addBlackhorn(raidPrefix string) {
	createBlackhornHeroicPreset(raidPrefix, 25, 89_671_248, 295_695, 80_759_953, 249_124)
}

func createBlackhornHeroicPreset(raidPrefix string, raidSize int32, bossHealth float64, bossMinBaseDamage float64, addHealth float64, addMinBaseDamage float64) {
	bossName := fmt.Sprintf("Warmaster Blackhorn %d H", raidSize)
	addName := fmt.Sprintf("Goriona %d H", raidSize)

	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: raidPrefix,

		Config: &proto.Target{
			Id:        blackhornID,
			Name:      bossName,
			Level:     88,
			MobType:   proto.MobType_MobTypeHumanoid,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      bossHealth,
				stats.Armor:       11977,
				stats.AttackPower: 0, // actual value doesn't matter in Cata, as long as damage parameters are fit consistently
			}.ToProtoArray(),

			SpellSchool:   proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:    1.5,
			MinBaseDamage: bossMinBaseDamage,
			DamageSpread:  blackhornMeleeDamageSpread,
			TargetInputs:  blackhornTargetInputs(),
		},

		AI: makeBlackhornAI(raidSize, true),
	})

	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: raidPrefix,

		Config: &proto.Target{
			Id:        gorionaID,
			Name:      addName,
			Level:     88,
			MobType:   proto.MobType_MobTypeDragonkin,
			TankIndex: 1,

			Stats: stats.Stats{
				stats.Health:      addHealth,
				stats.Armor:       11977,
				stats.AttackPower: 0, // actual value doesn't matter in Cata, as long as damage parameters are fit consistently
			}.ToProtoArray(),

			SpellSchool:   proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:    1.5,
			MinBaseDamage: addMinBaseDamage,
			DamageSpread:  0.2,
			TargetInputs:  []*proto.TargetInput{},
		},

		AI: makeBlackhornAI(raidSize, false),
	})

	core.AddPresetEncounter(bossName+" P2", []string{
		raidPrefix + "/" + bossName,
		raidPrefix + "/" + addName,
	})
}

func blackhornTargetInputs() []*proto.TargetInput {
	return []*proto.TargetInput{
		{
			Label:       "Tank swap interval",
			Tooltip:     "Elapsed time (in seconds) between simulated tank swaps",
			InputType:   proto.InputType_Number,
			NumberValue: 30,
		},
		{
			Label:       "Add de-activation time",
			Tooltip:     "Simulation time (in seconds) at which to disable Goriona's attacks",
			InputType:   proto.InputType_Number,
			NumberValue: 96,
		},
		{
			Label:       "Burn phase HP %",
			Tooltip:     "% of boss HP remaining when Goriona is de-activated",
			InputType:   proto.InputType_Number,
			NumberValue: 73,
		},
		{
			Label:       "Nerf state",
			Tooltip:     "Strength of the stacking Power of the Aspects debuff",
			InputType:   proto.InputType_Enum,
			EnumOptions: []string{"0%", "5%", "10%", "15%", "20%", "25%", "30%", "35%"},
			EnumValue:   0,
		},
	}
}

func makeBlackhornAI(raidSize int32, isBoss bool) core.AIFactory {
	return func() core.TargetAI {
		return &BlackhornAI{
			raidSize: raidSize,
			isBoss:   isBoss,
		}
	}
}

type BlackhornAI struct {
	// Unit references
	Target     *core.Target
	BossUnit   *core.Unit
	AddUnit    *core.Unit
	MainTank   *core.Unit
	OffTank    *core.Unit
	ValidTanks []*core.Unit

	// Static parameters associated with a given preset
	raidSize int32
	isBoss   bool

	// Dynamic parameters taken from user inputs
	tankSwapInterval             time.Duration
	disableAddAt                 time.Duration
	cleavePhaseVengeanceInterval time.Duration
	cleavePhaseVengeanceGain     int32
	nerfLevel                    int32

	// Spell + aura references
	Devastate      *core.Spell
	DisruptingRoar *core.Spell
	TwilightBreath *core.Spell
}

func (ai *BlackhornAI) Initialize(target *core.Target, config *proto.Target) {
	// Save unit references
	ai.Target = target
	ai.Target.AutoAttacks.MHConfig().ActionID.Tag = core.TernaryInt32(ai.isBoss, blackhornID, gorionaID)

	if ai.isBoss {
		ai.BossUnit = &target.Unit
		ai.AddUnit = &target.NextTarget().Unit
	} else {
		ai.AddUnit = &target.Unit
		ai.BossUnit = &target.NextTarget().Unit
	}

	ai.MainTank = ai.BossUnit.CurrentTarget
	ai.OffTank = ai.AddUnit.CurrentTarget

	ai.ValidTanks = core.FilterSlice([]*core.Unit{ai.MainTank, ai.OffTank}, func(unit *core.Unit) bool {
		return unit != nil
	})

	// Save user input parameters
	if ai.isBoss {
		ai.tankSwapInterval = core.DurationFromSeconds(config.TargetInputs[0].NumberValue)
		ai.disableAddAt = core.DurationFromSeconds(config.TargetInputs[1].NumberValue)
		ai.cleavePhaseVengeanceGain = 100 - int32(config.TargetInputs[2].NumberValue)
		ai.cleavePhaseVengeanceInterval = ai.disableAddAt / time.Duration(ai.cleavePhaseVengeanceGain)
		ai.nerfLevel = config.TargetInputs[3].EnumValue
	}

	// Register relevant spells and auras
	ai.registerDevastate()
	ai.registerDisruptingRoar()
	ai.registerVengeance()
	ai.registerTwilightBreath()
	ai.registerPowerOfTheAspects()
}

func (ai *BlackhornAI) registerDevastate() {
	if !ai.isBoss {
		return
	}

	sunderActionID := core.ActionID{SpellID: 108043}
	sunderDebuffConfig := core.Aura{
		Label:     "Sunder Armor",
		ActionID:  sunderActionID,
		Duration:  time.Second * 30,
		MaxStacks: 5,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.ArmorMultiplier *= (1.0 - 0.2*float64(newStacks)) / (1.0 - 0.2*float64(oldStacks))
		},
	}

	for _, tankUnit := range ai.ValidTanks {
		tankUnit.GetOrRegisterAura(sunderDebuffConfig)
	}

	ai.Devastate = ai.BossUnit.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 108042},
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		Flags:            core.SpellFlagMeleeMetrics,
		DamageMultiplier: 1.2,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.BossGCD,
			},

			CD: core.Cooldown{
				Timer:    ai.BossUnit.NewTimer(),
				Duration: time.Second * 8,
			},

			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, tankTarget *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.AutoAttacks.MH().EnemyWeaponDamage(sim, spell.MeleeAttackPower(), blackhornMeleeDamageSpread)
			result := spell.CalcAndDealDamage(sim, tankTarget, baseDamage, spell.OutcomeEnemyMeleeWhite)

			if result.Landed() {
				sunderAura := tankTarget.GetAuraByID(sunderActionID)

				if sunderAura != nil {
					sunderAura.Activate(sim)
					sunderAura.AddStack(sim)
				}
			}

			// Devastate resets swing timer whether or not it landed
			spell.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		},
	})

	ai.BossUnit.RegisterResetEffect(func(sim *core.Simulation) {
		ai.Devastate.CD.Set(core.DurationFromSeconds(sim.RandomFloat("Devastate Timing") * ai.Devastate.CD.Duration.Seconds()))
	})
}

func (ai *BlackhornAI) registerDisruptingRoar() {
	if !ai.isBoss {
		return
	}

	// 0 - 10H, 1 - 25H
	scalingIndex := core.TernaryInt(ai.raidSize == 10, 0, 1)

	// https://wago.tools/db2/SpellEffect?build=4.4.2.58947&filter[SpellID]=108044&page=1
	disruptingRoarBase := []float64{71250, 92625}[scalingIndex]
	disruptingRoarVariance := []float64{7500, 9750}[scalingIndex]

	ai.DisruptingRoar = ai.BossUnit.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 108044},
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagIgnoreArmor | core.SpellFlagAPL,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.BossGCD,
			},

			CD: core.Cooldown{
				Timer:    ai.BossUnit.NewTimer(),
				Duration: time.Second * 18,
			},

			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Raid.AllPlayerUnits {
				damageRoll := disruptingRoarBase + disruptingRoarVariance*sim.RandomFloat("Disrupting Roar Damage")
				spell.CalcAndDealDamage(sim, aoeTarget, damageRoll, spell.OutcomeAlwaysHit)
			}

			// Different swing delay behavior from Devastate based on log analysis
			spell.Unit.AutoAttacks.PauseMeleeBy(sim, core.BossGCD+1)
		},
	})

	ai.BossUnit.RegisterResetEffect(func(sim *core.Simulation) {
		ai.DisruptingRoar.CD.Set(core.DurationFromSeconds(sim.RandomFloat("Disrupting Roar Timing") * ai.DisruptingRoar.CD.Duration.Seconds()))
	})
}

func (ai *BlackhornAI) registerVengeance() {
	if !ai.isBoss {
		return
	}

	ai.BossUnit.RegisterAura(core.Aura{
		Label:     "Vengeance",
		ActionID:  core.ActionID{SpellID: 108045},
		MaxStacks: 100,
		Duration:  core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   ai.cleavePhaseVengeanceInterval,
				NumTicks: int(ai.cleavePhaseVengeanceGain),
				Priority: core.ActionPriorityDOT,

				OnAction: func(sim *core.Simulation) {
					aura.AddStack(sim)
				},
			})
		},

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= (1.0 + 0.01*float64(newStacks)) / (1.0 + 0.01*float64(oldStacks))

			if newStacks == ai.cleavePhaseVengeanceGain {
				newNumTicks := int(aura.MaxStacks - newStacks)
				newPeriod := sim.GetRemainingDuration() / time.Duration(newNumTicks)

				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   newPeriod,
					NumTicks: newNumTicks,
					Priority: core.ActionPriorityDOT,

					OnAction: func(sim *core.Simulation) {
						aura.AddStack(sim)
					},
				})
			}
		},
	})
}

func (ai *BlackhornAI) registerTwilightBreath() {
	// 0 - 10H, 1 - 25H
	scalingIndex := core.TernaryInt(ai.raidSize == 10, 0, 1)

	// https://wago.tools/db2/SpellEffect?build=4.4.2.58947&filter[SpellID]=110212&page=1
	twilightBreathBase := []float64{81600, 120000}[scalingIndex]
	twilightBreathVariance := []float64{6800, 10000}[scalingIndex]

	twilightBreathActionID := core.ActionID{SpellID: 110212}
	twilightBreathCastTime := time.Second * 2
	twilightBreathConfig := core.SpellConfig{
		ActionID:         twilightBreathActionID,
		SpellSchool:      core.SpellSchoolShadow,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagAPL,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.BossGCD * 2,
				CastTime: twilightBreathCastTime,
			},

			CD: core.Cooldown{
				Timer:    ai.AddUnit.NewTimer(),
				Duration: time.Second * 18,
			},

			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Conal breath will hit both tanks but no one else.
			for _, tankUnit := range ai.ValidTanks {
				damageRoll := twilightBreathBase + twilightBreathVariance*sim.RandomFloat("Twilight Breath Damage")
				spell.CalcAndDealDamage(sim, tankUnit, damageRoll, spell.OutcomeAlwaysHit)
			}

			// No swing reset/delay logic here, boss can continue to melee while casting based on log analysis.
		},
	}

	if !ai.isBoss {
		ai.TwilightBreath = ai.AddUnit.RegisterSpell(twilightBreathConfig)
	} else {
		ai.AddUnit.RegisterResetEffect(func(sim *core.Simulation) {
			// Hacky work-around to the add AI not having access to user input parameters
			twilightBreathSpell := ai.AddUnit.GetSpell(twilightBreathActionID)
			twilightBreathSpell.CD.Set(core.DurationFromSeconds(sim.RandomFloat("Twilight Breath Timing") * twilightBreathSpell.CD.Duration.Seconds()))

			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt:     ai.disableAddAt - twilightBreathCastTime,
				Priority: core.ActionPriorityDOT,

				OnAction: func(_ *core.Simulation) {
					twilightBreathSpell.CD.Set(core.NeverExpires)
				},
			})
		})
	}
}

func (ai *BlackhornAI) registerPowerOfTheAspects() {
	if !ai.isBoss || (ai.nerfLevel == 0) {
		return
	}

	damageMultiplier := 1.0 - 0.05*float64(ai.nerfLevel)
	auraID := 109250 + ai.nerfLevel

	debuffConfig := core.Aura{
		Label:    "Power of the Aspects",
		ActionID: core.ActionID{SpellID: auraID},
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= damageMultiplier
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= damageMultiplier
		},
	}

	ai.BossUnit.GetOrRegisterAura(debuffConfig)
	ai.AddUnit.GetOrRegisterAura(debuffConfig)
}

func (ai *BlackhornAI) Reset(sim *core.Simulation) {
	// Randomize GCD and swing timings to prevent fake APL-Haste couplings.
	ai.Target.ExtendGCDUntil(sim, core.DurationFromSeconds(sim.RandomFloat("Specials Timing")*core.BossGCD.Seconds()))
	ai.Target.AutoAttacks.RandomizeMeleeTiming(sim)

	if !ai.isBoss {
		return
	}

	// Set up delayed action for disabling add swings.
	core.StartDelayedAction(sim, core.DelayedActionOptions{
		DoAt:     ai.disableAddAt,
		Priority: core.ActionPriorityDOT,

		OnAction: func(sim *core.Simulation) {
			ai.AddUnit.AutoAttacks.CancelAutoSwing(sim)
		},
	})

	// Set up periodic action for tank swaps.
	core.StartPeriodicAction(sim, core.PeriodicActionOptions{
		Period:   ai.tankSwapInterval,
		NumTicks: int(sim.Duration / ai.tankSwapInterval),
		Priority: core.ActionPriorityDOT,

		OnAction: func(sim *core.Simulation) {
			newBossTank := core.Ternary((sim.CurrentTime/ai.tankSwapInterval)%2 == 0, ai.MainTank, ai.OffTank)
			ai.swapTargets(sim, ai.BossUnit, newBossTank, true)
			ai.Devastate.CD.Set(sim.CurrentTime + core.DurationFromSeconds(sim.RandomFloat("Devastate Timing")*ai.Devastate.CD.Duration.Seconds()))
			newAddTank := core.Ternary(newBossTank == ai.MainTank, ai.OffTank, ai.MainTank)
			ai.swapTargets(sim, ai.AddUnit, newAddTank, sim.CurrentTime < ai.disableAddAt)
		},
	})
}

func (ai *BlackhornAI) swapTargets(sim *core.Simulation, npc *core.Unit, newTankTarget *core.Unit, enableAutos bool) {
	npc.AutoAttacks.CancelAutoSwing(sim)
	npc.CurrentTarget = newTankTarget

	if newTankTarget != nil {
		newTankTarget.CurrentTarget = npc
	}

	if enableAutos {
		npc.AutoAttacks.EnableAutoSwing(sim)
		npc.AutoAttacks.RandomizeMeleeTiming(sim)
	}
}

func (ai *BlackhornAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget
	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if ai.isBoss && (target == ai.BossUnit.CurrentTarget) && ai.Devastate.IsReady(sim) {
		ai.Devastate.Cast(sim, target)
		return
	}

	if ai.isBoss && ai.DisruptingRoar.IsReady(sim) && sim.Proc(0.75, "Disrupting Roar AI") {
		ai.DisruptingRoar.Cast(sim, target)
		return
	}

	if !ai.isBoss && ai.TwilightBreath.IsReady(sim) && sim.Proc(0.75, "Twilight Breath AI") {
		ai.TwilightBreath.Cast(sim, target)
		return
	}

	ai.Target.ExtendGCDUntil(sim, sim.CurrentTime+core.BossGCD)
}
