package encounters

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func addMovementAI() {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: "Default",
		Config: &proto.Target{
			Id:        31147,
			Name:      "Movement",
			Level:     88,
			MobType:   proto.MobType_MobTypeMechanical,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      120_016_403,
				stats.Armor:       11977,
				stats.AttackPower: 650,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.5,
			MinBaseDamage:    210000,
			DamageSpread:     0.4,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs: []*proto.TargetInput{
				{
					Label:       "Movement Interval",
					Tooltip:     "How often the player will move in seconds",
					InputType:   proto.InputType_Number,
					NumberValue: 10.0,
				},
				{
					Label:       "Reaction Time",
					Tooltip:     "How long the player can wait for casts to finish before moving in seconds",
					InputType:   proto.InputType_Number,
					NumberValue: 1.5,
				},
				{
					Label:       "Yards",
					Tooltip:     "How many yards the player moves",
					InputType:   proto.InputType_Number,
					NumberValue: 5,
				},
			},
		},
		AI: NewMovementAI(),
	})
	core.AddPresetEncounter("Movement", []string{
		"Default/Movement",
	})
}

type MovementAI struct {
	Target       *core.Target
	LastMoveTime time.Duration
	MoveInterval time.Duration // How often moves happen
	ReactionTime time.Duration // Time available to react before area should be cleared
	MoveYards    float64       // Duration of the move
}

func NewMovementAI() core.AIFactory {
	return func() core.TargetAI {
		return &MovementAI{}
	}
}

func (ai *MovementAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target

	if len(config.TargetInputs) > 0 {
		ai.MoveInterval = core.DurationFromSeconds(config.TargetInputs[0].NumberValue)
	} else {
		ai.MoveInterval = core.DurationFromSeconds(10)
	}

	if len(config.TargetInputs) > 1 {
		ai.ReactionTime = core.DurationFromSeconds(config.TargetInputs[1].NumberValue)
	} else {
		ai.ReactionTime = core.DurationFromSeconds(1.5)
	}

	if len(config.TargetInputs) > 2 {
		ai.MoveYards = config.TargetInputs[2].NumberValue
	} else {
		ai.MoveYards = 5.0
	}
}

func (ai *MovementAI) Reset(sim *core.Simulation) {
	ai.LastMoveTime = 0
}

func (ai *MovementAI) ExecuteCustomRotation(sim *core.Simulation) {
	players := sim.Raid.AllPlayerUnits

	if !ai.ShouldMove(sim) {
		return
	}
	for i := 0; i < len(players); i++ {
		player := players[i]
		duration := ai.TimeToMove(ai.MoveYards, player)
		if player.Hardcast.Expires > sim.CurrentTime && !player.Hardcast.CanMove {
			castEndsAt := player.Hardcast.Expires - sim.CurrentTime
			// if castEndsAt < ai.ReactionTime {
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt:     sim.CurrentTime + castEndsAt,
				Priority: core.ActionPriorityPrePull + 1,
				OnAction: func(s *core.Simulation) {
					player.MoveDuration(duration, sim)
				},
			})
			// } else {
			// 	// Cancel casted spell and move immediately
			// 	// For now we do nothing in this scenario
			// 	return
			// }
		} else {
			player.MoveDuration(duration, sim)
		}
	}
}
func (ai *MovementAI) TimeToMove(distance float64, unit *core.Unit) time.Duration {
	return core.DurationFromSeconds(distance / unit.GetMovementSpeed())
}
func (ai *MovementAI) ShouldMove(sim *core.Simulation) bool {
	if sim.CurrentTime-ai.LastMoveTime >= ai.MoveInterval {
		ai.LastMoveTime = sim.CurrentTime
		return true
	}
	return false
}
