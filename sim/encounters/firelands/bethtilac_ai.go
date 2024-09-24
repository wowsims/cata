package firelands

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func addBethtilac(raidPrefix string) {
	createBethtilacPreset(raidPrefix, 25, true, 52498, 83_658_808, 151_332)
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
			DamageSpread:  0.4711,
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
			Label:     "Model pre-nerf damage",
			Tooltip:   "Extrapolate pre-nerf boss damage parameters rather than using nerfed values from PTR logs.",
			InputType: proto.InputType_Bool,
			BoolValue: false,
		},
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
	preNerf       bool
	includeFrenzy bool

	// Spell + aura references
	emberFlame *core.Spell
	frenzyAura *core.Aura
}

func (ai *BethtilacAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
	ai.preNerf = config.TargetInputs[0].BoolValue
	ai.includeFrenzy = config.TargetInputs[1].BoolValue
	ai.registerEmberFlameSpell()
	ai.registerFrenzySpell()

	// "Undo" the 15% damage nerf on the PTR patch state if requested
	if ai.preNerf {
		ai.Target.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 0.85
	}
}

func (ai *BethtilacAI) Reset(sim *core.Simulation) {
	ai.emberFlame.CD.Use(sim)
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
	phase1BaseDamageValues := []float64{14152, 15725, 20229, 25858}
	phase2BaseDamageValues := []float64{15660, 17400, 23809, 30421}
	emberFlameBase := core.TernaryFloat64(ai.includeFrenzy, phase2BaseDamageValues[scalingIndex], phase1BaseDamageValues[scalingIndex])
	emberFlameVariance := emberFlameBase * 0.1622
	emberFlameSpellID := core.TernaryInt32(ai.includeFrenzy, 99859, 98934)

	ai.emberFlame = ai.Target.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: emberFlameSpellID},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagIgnoreResists,
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
