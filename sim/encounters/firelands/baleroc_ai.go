package firelands

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func addBaleroc(raidPrefix string) {
	createBalerocPreset(raidPrefix, 25, true, 53494, 195_576_084, 323_321)
}

func createBalerocPreset(raidPrefix string, raidSize int32, isHeroic bool, bossNpcId int32, bossHealth float64, bossMinBaseDamage float64) {
	targetName := fmt.Sprintf("Baleroc %d", raidSize)

	if isHeroic {
		targetName += " H"
	}

	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: raidPrefix,

		Config: &proto.Target{
			Id:      bossNpcId,
			Name:    targetName,
			Level:   88,
			MobType: proto.MobType_MobTypeElemental,

			//By default, the off-tank will start the pull in order
			// to accumulate a few Blaze of Glory stacks for living
			// the first Decimation Blade window.
			TankIndex:       1,
			SecondTankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      bossHealth,
				stats.Armor:       11977,
				stats.AttackPower: 0, // actual value doesn't matter in Cata, as long as damage parameters are fit consistently
			}.ToProtoArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.0,
			MinBaseDamage:    bossMinBaseDamage,
			DualWield:        true,
			DualWieldPenalty: false,
			DamageSpread:     0.4738,
			TargetInputs:     balerocTargetInputs(),
		},

		AI: makeBalerocAI(raidSize, isHeroic),
	})

	core.AddPresetEncounter(targetName, []string{
		raidPrefix + "/" + targetName,
	})
}

func balerocTargetInputs() []*proto.TargetInput {
	return []*proto.TargetInput{
		{
			Label:     "Tank swap for Decimation Blade",
			Tooltip:   "If checked, the boss will be tanked by Tank 2 rather than the Main Tank during Decimation Blade windows.",
			InputType: proto.InputType_Bool,
			BoolValue: true,
		},
		{
			Label:       "Initial OT Blaze of Glory stacks",
			Tooltip:     "If non-zero, Tank 2 will initiate the pull until the specified stack count is accumulated. Larger values allow the OT to gear more aggressively while still meeting the 250k HP threshold for the first Decimation Blade window, but make life more difficult for the MT during Inferno Blade windows. This input is ignored if tank swaps are disabled.",
			InputType:   proto.InputType_Number,
			NumberValue: 2,
		},
		{
			Label:       "Initial Vital Spark ramp rate",
			Tooltip:     "For the first minute of the pull, the simulated healing agent will accumulate this many Vital Spark stacks on average between Vital Flame refreshes on the tank.",
			InputType:   proto.InputType_Number,
			NumberValue: 7.5,
		},
		{
			Label:       "Steady state Vital Spark accumulation rate",
			Tooltip:     "Slower rate for the remainder of the pull.",
			InputType:   proto.InputType_Number,
			NumberValue: 4,
		},
	}
}

func makeBalerocAI(raidSize int32, isHeroic bool) core.AIFactory {
	return func() core.TargetAI {
		return &BalerocAI{
			raidSize: raidSize,
			isHeroic: isHeroic,
		}
	}
}

type BalerocAI struct {
	// Unit references
	Target   *core.Target
	MainTank *core.Unit
	OffTank  *core.Unit

	// Static parameters associated with a given preset
	raidSize int32
	isHeroic bool

	// Dynamic parameters taken from user inputs
	tankSwap               bool
	stackCountForFirstSwap int32
	initialHealerStackGain float64
	steadyHealerStackGain  float64

	// Spell + aura references
	blazeOfGlory     *core.Spell
	infernoBlade     *core.Spell
	decimationBlade  *core.Spell
	sharedBladeTimer *core.Timer
}

func (ai *BalerocAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target

	// OT starts the pull by default unless configured otherwise.
	ai.OffTank = target.CurrentTarget
	ai.MainTank = target.SecondaryTarget
	ai.tankSwap = config.TargetInputs[0].BoolValue
	ai.stackCountForFirstSwap = int32(config.TargetInputs[1].NumberValue)

	if !ai.tankSwap {
		ai.OffTank = nil
		ai.stackCountForFirstSwap = 0
	}

	if ai.stackCountForFirstSwap <= 0 {
		target.CurrentTarget = ai.MainTank
		target.SecondaryTarget = ai.OffTank
	}

	ai.initialHealerStackGain = config.TargetInputs[2].NumberValue
	ai.steadyHealerStackGain = config.TargetInputs[3].NumberValue

	ai.registerBlazeOfGlory()
	ai.registerBlades()
	ai.registerVitalSpark()
}

func (ai *BalerocAI) randomizeAutoTiming(sim *core.Simulation) {
	// Add random auto delay to avoid artificial Haste breakpoints from APL evaluations after tank autos
	// synchronizing with damage events.
	swingDur := ai.Target.AutoAttacks.MainhandSwingSpeed()
	randomAutoOffset := core.DurationFromSeconds(sim.RandomFloat("Melee Timing") * swingDur.Seconds() / 2)
	ai.Target.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime-swingDur+randomAutoOffset, true)
}

func (ai *BalerocAI) Reset(sim *core.Simulation) {
	// Set starting CDs on abilities.
	ai.sharedBladeTimer.Set(time.Second * 30)
	ai.blazeOfGlory.CD.Set(ai.blazeOfGlory.CD.Duration)

	// Randomize melee and cast timings to prevent fake APL-Haste couplings.
	ai.randomizeAutoTiming(sim)
	ai.Target.ExtendGCDUntil(sim, sim.CurrentTime+core.DurationFromSeconds(sim.RandomFloat("Specials Timing")*core.BossGCD.Seconds()))
}

func (ai *BalerocAI) swapTargets(sim *core.Simulation, newTankTarget *core.Unit) {
	ai.Target.AutoAttacks.CancelAutoSwing(sim)
	ai.Target.CurrentTarget = newTankTarget
	ai.Target.AutoAttacks.EnableAutoSwing(sim)
	ai.randomizeAutoTiming(sim)
}

func (ai *BalerocAI) registerBlazeOfGlory() {
	// Boss buff aura setup
	const maxPossibleStacks int32 = 45 // limited by Berserk timer

	incendiarySoulAura := ai.Target.RegisterAura(core.Aura{
		Label:     "Incendiary Soul",
		ActionID:  core.ActionID{SpellID: 99369},
		MaxStacks: maxPossibleStacks,
		Duration:  core.NeverExpires,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= (1.0 + 0.2*float64(newStacks)) / (1.0 + 0.2*float64(oldStacks))

			if (newStacks > 0) && (newStacks == ai.stackCountForFirstSwap) {
				ai.swapTargets(sim, ai.MainTank)
			}
		},
	})

	// Tank debuff aura setup
	blazeOfGloryActionID := core.ActionID{SpellID: 99252}

	for _, tankUnit := range []*core.Unit{ai.MainTank, ai.OffTank} {
		if tankUnit == nil {
			continue
		}

		// Set up HP multiplier stat dependencies for each stack level.
		hpDepByStackCount := map[int32]*stats.StatDependency{}

		for i := int32(1); i <= maxPossibleStacks; i++ {
			hpDepByStackCount[i] = tankUnit.NewDynamicMultiplyStat(stats.Health, 1.0+0.2*float64(i))
		}

		// Blaze of Glory applications also heal the player, just like
		// most other temporary max health increases.
		healthMetrics := tankUnit.NewHealthMetrics(blazeOfGloryActionID)

		tankUnit.GetOrRegisterAura(core.Aura{
			Label:     "Blaze of Glory",
			ActionID:  blazeOfGloryActionID,
			MaxStacks: maxPossibleStacks,
			Duration:  core.NeverExpires,

			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= (1.0 + 0.2*float64(newStacks)) / (1.0 + 0.2*float64(oldStacks))

				// Cache max HP prior to processing multipliers.
				oldMaxHp := aura.Unit.MaxHealth()

				if oldStacks > 0 {
					aura.Unit.DisableDynamicStatDep(sim, hpDepByStackCount[oldStacks])
				}

				if newStacks > 0 {
					aura.Unit.EnableDynamicStatDep(sim, hpDepByStackCount[newStacks])
				}

				hpGain := aura.Unit.MaxHealth() - oldMaxHp

				if hpGain > 0 {
					aura.Unit.GainHealth(sim, hpGain, healthMetrics)
				}
			},
		})
	}

	// Stack accumlation spell
	ai.blazeOfGlory = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID: blazeOfGloryActionID,
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, tankTarget *core.Unit, _ *core.Spell) {
			if (tankTarget != nil) && (tankTarget == ai.Target.CurrentTarget) {
				blazeOfGloryAura := tankTarget.GetAuraByID(blazeOfGloryActionID)

				if blazeOfGloryAura != nil {
					blazeOfGloryAura.Activate(sim)
					blazeOfGloryAura.AddStack(sim)
				}
			}

			incendiarySoulAura.Activate(sim)
			incendiarySoulAura.AddStack(sim)
		},
	})
}

func (ai *BalerocAI) registerBlades() {
	// First register the blade auras and activation spells.
	const bladeDuration = time.Second * 15
	const bladeCooldown = time.Second * 45 // very first one is special cased as 30s
	const bladeCastTime = time.Millisecond * 1500

	sharedBladeCastHandler := func(sim *core.Simulation) {
		// First, schedule a swing timer reset to fire on cast completion.
		ai.Target.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+bladeCastTime, true)

		// Then delay the off-hand until the aura expires.
		swingDur := ai.Target.AutoAttacks.MainhandSwingSpeed()
		first2hSwing := ai.Target.AutoAttacks.NextAttackAt()
		num2hSwings := bladeDuration / swingDur
		ai.Target.AutoAttacks.SetOffhandSwingAt(first2hSwing + num2hSwings*swingDur + swingDur/2)

		// Finally, reset the CD on Blaze of Glory at start of cast.
		ai.blazeOfGlory.CD.Set(sim.CurrentTime + ai.blazeOfGlory.CD.Duration)
	}

	infernoBladeActionID := core.ActionID{SpellID: 99350}
	infernoBladeAura := ai.Target.RegisterAura(core.Aura{
		Label:    "Inferno Blade",
		ActionID: infernoBladeActionID,
		Duration: bladeDuration,
	})

	ai.infernoBlade = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID: infernoBladeActionID,
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.BossGCD,
				CastTime: bladeCastTime,
			},

			IgnoreHaste: true,

			SharedCD: core.Cooldown{
				Timer:    ai.Target.GetOrInitTimer(&ai.sharedBladeTimer),
				Duration: bladeCooldown,
			},

			ModifyCast: func(sim *core.Simulation, _ *core.Spell, _ *core.Cast) {
				sharedBladeCastHandler(sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			infernoBladeAura.Activate(sim)
		},
	})

	decimationBladeActionID := core.ActionID{SpellID: 99352}
	decimationBladeAura := ai.Target.RegisterAura(core.Aura{
		Label:    "Decimation Blade",
		ActionID: decimationBladeActionID,
		Duration: bladeDuration,

		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			if ai.tankSwap && (ai.Target.CurrentTarget == ai.OffTank) {
				ai.swapTargets(sim, ai.MainTank)
			}
		},
	})

	ai.decimationBlade = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID: decimationBladeActionID,
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.BossGCD,
				CastTime: bladeCastTime,
			},

			IgnoreHaste: true,

			SharedCD: core.Cooldown{
				Timer:    ai.Target.GetOrInitTimer(&ai.sharedBladeTimer),
				Duration: bladeCooldown,
			},

			ModifyCast: func(sim *core.Simulation, _ *core.Spell, _ *core.Cast) {
				if ai.tankSwap {
					ai.swapTargets(sim, ai.OffTank)
				}

				sharedBladeCastHandler(sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			decimationBladeAura.Activate(sim)
		},
	})

	// Then register the strikes that replace boss melees during each blade.
	// 0 - 10N, 1 - 25N, 2 - 10H, 3 - 25H
	scalingIndex := core.TernaryInt(ai.raidSize == 10, core.TernaryInt(ai.isHeroic, 2, 0), core.TernaryInt(ai.isHeroic, 3, 1))

	// https://wago.tools/db2/SpellEffect?build=4.4.1.57294&filter[SpellID]=99351&page=1&sort[SpellID]=asc
	infernoStrikeBase := []float64{97499, 165749, 136499, 232049}[scalingIndex]
	infernoStrikeVariance := []float64{5000, 8500, 7000, 11900}[scalingIndex]

	infernoStrike := ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 99351},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagMeleeMetrics,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damageRoll := infernoStrikeBase + infernoStrikeVariance*sim.RandomFloat("Inferno Strike Damage")
			spell.CalcAndDealDamage(sim, target, damageRoll, spell.OutcomeEnemyMeleeWhite)
		},
	})

	decimatingStrikeActionID := core.ActionID{SpellID: 99353}
	decimatingStrikeDebuffConfig := core.Aura{
		Label:    "Decimating Strike",
		ActionID: decimatingStrikeActionID,
		Duration: time.Second * 4,

		OnGain: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.HealingDealtMultiplier *= 0.1
		},

		OnExpire: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.HealingDealtMultiplier /= 0.1
		},
	}

	for _, tankUnit := range []*core.Unit{ai.MainTank, ai.OffTank} {
		if tankUnit != nil {
			tankUnit.GetOrRegisterAura(decimatingStrikeDebuffConfig)
		}
	}

	decimatingStrike := ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         decimatingStrikeActionID,
		SpellSchool:      core.SpellSchoolShadow,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreModifiers | core.SpellFlagIgnoreArmor,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, tankTarget *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, tankTarget, max(0.9*tankTarget.MaxHealth(), 250000), spell.OutcomeEnemyMeleeWhite)

			if result.Landed() {
				debuffAura := tankTarget.GetAuraByID(decimatingStrikeActionID)

				if debuffAura != nil {
					debuffAura.Activate(sim)
				}
			}

			// MT should taunt as soon as the final Decimating Strike goes out in order to maximize their Blaze of Glory stack count.
			if ai.tankSwap && (ai.stackCountForFirstSwap > 0) && (decimationBladeAura.ExpiresAt() < ai.Target.AutoAttacks.NextAttackAt()) {
				ai.swapTargets(sim, ai.MainTank)
			}
		},
	})

	ai.Target.AutoAttacks.SetReplaceMHSwing(func(_ *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
		if infernoBladeAura.IsActive() {
			return infernoStrike
		} else if decimationBladeAura.IsActive() {
			return decimatingStrike
		} else {
			return mhSwingSpell
		}
	})
}

func (ai *BalerocAI) registerVitalSpark() {
	// Since the tank sim healing model collapses the behavior of several
	// real healers onto one aggregate agent, we will approximate the effect
	// of the Vital Spark mechanic by registering a healing taken buff on
	// the tank that ramps over time based on the encounter settings.
	const vitalFlameDuration = time.Second * 15

	calcStackGain := func(sim *core.Simulation) int32 {
		if sim.CurrentTime <= vitalFlameDuration {
			return 0
		}

		var avgGain float64

		if sim.CurrentTime <= vitalFlameDuration*2 {
			avgGain = ai.initialHealerStackGain * 2
		} else if sim.CurrentTime <= vitalFlameDuration*4 {
			avgGain = ai.initialHealerStackGain
		} else {
			avgGain = ai.steadyHealerStackGain
		}

		return int32(math.Round(sim.Roll(avgGain-1, avgGain+1)))
	}

	vitalFlameConfig := core.Aura{
		Label:     "Vital Flame",
		ActionID:  core.ActionID{SpellID: 99263},
		MaxStacks: math.MaxInt32,
		Duration:  vitalFlameDuration,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.ExternalHealingTakenMultiplier *= (1.0 + 0.05*float64(newStacks)) / (1.0 + 0.05*float64(oldStacks))
		},

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			period := vitalFlameDuration - aura.Unit.ReactionTime
			numTicks := int(sim.GetRemainingDuration() / period)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   period,
				NumTicks: numTicks,
				Priority: core.ActionPriorityDOT,

				OnAction: func(sim *core.Simulation) {
					aura.Refresh(sim)
					aura.SetStacks(sim, aura.GetStacks()+calcStackGain(sim))
				},
			})
		},
	}

	for _, tankUnit := range []*core.Unit{ai.MainTank, ai.OffTank} {
		if tankUnit != nil {
			tankUnit.GetOrRegisterAura(vitalFlameConfig)
		}
	}
}

func (ai *BalerocAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget
	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if ai.sharedBladeTimer.IsReady(sim) {
		// First Blade is always Inferno, subsequent ones are randomized
		if sim.CurrentTime < time.Minute {
			ai.infernoBlade.Cast(sim, target)
		} else if sim.Proc(0.5, "Baleroc Blade Selection") {
			ai.infernoBlade.Cast(sim, target)
		} else {
			ai.decimationBlade.Cast(sim, target)
		}
	} else if ai.blazeOfGlory.IsReady(sim) {
		ai.blazeOfGlory.Cast(sim, target)
	}

	ai.Target.ExtendGCDUntil(sim, sim.CurrentTime+core.BossGCD)
}
