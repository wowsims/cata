package firelands

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func addBethtilac(raidPrefix string) {
	createBethtilacPreset(raidPrefix, 25, true, 52498, 98_518_124, 170_101)
}

func createBethtilacPreset(raidPrefix string, raidSize int32, isHeroic bool, bossNpcId int32, bossHealth float64, bossMinBaseDamage float64) {
	targetName := fmt.Sprintf("Beth'tilac %d", raidSize)

	if isHeroic {
		targetName += " H"
	}

	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: raidPrefix,

		Config: &proto.Target{
			Id:        bossNpcId,
			Name:      targetName,
			Level:     88,
			MobType:   proto.MobType_MobTypeBeast,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      bossHealth,
				stats.Armor:       11977,
				stats.AttackPower: 0, // actual value doesn't matter in Cata, as long as damage parameters are fit consistently
			}.ToProtoArray(),

			SpellSchool:   proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:    2.0,
			MinBaseDamage: bossMinBaseDamage,
			DamageSpread:  0.4832,
			TargetInputs:  bethtilacTargetInputs(),
		},

		AI: makeBethtilacAI(raidSize, isHeroic),
	})

	core.AddPresetEncounter(targetName, []string{
		raidPrefix + "/" + targetName,
	})
}

func bethtilacTargetInputs() []*proto.TargetInput {
	return []*proto.TargetInput{
		{
			Label:     "Include Frenzy",
			Tooltip:   "Model the burn phase of the encounter with the stacking boss damage amp. Note that The Widow's Kiss will not be modeled even if this option is selected, since it is assumed that an off-tank will be pre-taunting to soak the debuff.",
			InputType: proto.InputType_Bool,
			BoolValue: false,
		},
	}
}

func makeBethtilacAI(raidSize int32, isHeroic bool) core.AIFactory {
	return func() core.TargetAI {
		return &BethtilacAI{
			raidSize: raidSize,
			isHeroic: isHeroic,
		}
	}
}

type BethtilacAI struct {
	Target *core.Target

	// Static parameters associated with a given preset
	raidSize int32
	isHeroic bool

	// Dynamic parameters taken from user inputs
	includeFrenzy bool

	// Spell + aura references
	emberFlame *core.Spell
	frenzyAura *core.Aura
}

func (ai *BethtilacAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.includeFrenzy = config.TargetInputs[0].BoolValue
	ai.registerEmberFlameSpell()
	ai.registerFrenzySpell()
}

func (ai *BethtilacAI) Reset(sim *core.Simulation) {
	// Add random auto delay to avoid artificial Haste breakpoints coming from APL evaluations after tank autos synchronizing with
	// damage taken events.
	randomAutoOffset := core.DurationFromSeconds(sim.RandomFloat("Melee Timing") * ai.Target.AutoAttacks.MainhandSwingSpeed().Seconds())
	ai.Target.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime-randomAutoOffset, false)

	// Do the same for boss GCD as well.
	randomSpecialsOffset := core.DurationFromSeconds(sim.RandomFloat("Specials Timing") * core.BossGCD.Seconds())
	ai.Target.ExtendGCDUntil(sim, sim.CurrentTime+randomSpecialsOffset)
	ai.emberFlame.CD.Set(core.DurationFromSeconds(sim.RandomFloat("Ember Flame Timing") * ai.emberFlame.CD.Duration.Seconds()))
}

func (ai *BethtilacAI) registerFrenzySpell() {
	if !ai.includeFrenzy {
		return
	}

	ai.frenzyAura = ai.Target.RegisterAura(core.Aura{
		Label:     "Frenzy",
		ActionID:  core.ActionID{SpellID: 99497},
		MaxStacks: math.MaxInt32,
		Duration:  time.Second * 299,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= (1.0 + 0.05*float64(newStacks)) / (1.0 + 0.05*float64(oldStacks))
		},

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			period := time.Second * 5
			numTicks := int(sim.GetRemainingDuration() / period)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   period,
				NumTicks: numTicks,
				Priority: core.ActionPriorityDOT,

				OnAction: func(sim *core.Simulation) {
					aura.Refresh(sim)
					aura.AddStack(sim)
				},
			})
		},
	})
}

func (ai *BethtilacAI) registerEmberFlameSpell() {
	// 0 - 10N, 1 - 25N, 2 - 10H, 3 - 25H
	scalingIndex := core.TernaryInt(ai.raidSize == 10, core.TernaryInt(ai.isHeroic, 2, 0), core.TernaryInt(ai.isHeroic, 3, 1))

	// Update from second PTR test: it looks like Blizz is now using the damage parameters for the Phase 2 version of the spell in Phase
	// 1 as well. The spell ID for Phase 1 Ember Flame is still logged as 98934, but the damage values are consistent with 99859 instead,
	// suggesting that the tooltip for 98934 just hasn't been updated yet. To avoid confusion, the sim model will use 99859 for both
	// phases.
	//
	//phase1BaseDamageValues := []float64{14152, 15725, 20229, 25858}
	phase2BaseDamageValues := []float64{15660, 17400, 23809, 30421}
	//emberFlameBase := core.TernaryFloat64(ai.includeFrenzy, phase2BaseDamageValues[scalingIndex], phase1BaseDamageValues[scalingIndex])
	emberFlameBase := phase2BaseDamageValues[scalingIndex]
	emberFlameVariance := emberFlameBase * 0.1622
	//emberFlameSpellID := core.TernaryInt32(ai.includeFrenzy, 99859, 98934)
	emberFlameSpellID := int32(99859)

	ai.emberFlame = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: emberFlameSpellID},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagIgnoreArmor,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    ai.Target.NewTimer(),
				Duration: time.Second * 6,
			},

			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Raid.AllPlayerUnits {
				damageRoll := emberFlameBase + emberFlameVariance*sim.RandomFloat("Ember Flame Damage")
				spell.CalcAndDealDamage(sim, aoeTarget, damageRoll, spell.OutcomeAlwaysHit)
			}
		},
	})
}

func (ai *BethtilacAI) ExecuteCustomRotation(sim *core.Simulation) {
	target := ai.Target.CurrentTarget
	if target == nil {
		// For individual non tank sims we still want abilities to work
		target = &ai.Target.Env.Raid.Parties[0].Players[0].GetCharacter().Unit
	}

	if ai.emberFlame.IsReady(sim) {
		ai.emberFlame.Cast(sim, target)
	}

	ai.Target.ExtendGCDUntil(sim, sim.CurrentTime+core.BossGCD)
}
