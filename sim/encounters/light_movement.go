package encounters

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func addLightMovementAI() {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: "Movement",
		Config: &proto.Target{
			Id:        31147,
			Name:      "Light",
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
			TargetInputs:     []*proto.TargetInput{},
		},
		AI: NewLightMovementAI(),
	})
	core.AddPresetEncounter("Light", []string{
		"Movement/Light",
	})
}

type LightMovementAI struct {
	Target       *core.Target
	LastMoveTime time.Duration
}

func NewLightMovementAI() core.AIFactory {
	return func() core.TargetAI {
		return &LightMovementAI{}
	}
}

func (ai *LightMovementAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target
}

func (ai *LightMovementAI) Reset(sim *core.Simulation) {
	ai.LastMoveTime = 0
}

func (ai *LightMovementAI) ExecuteCustomRotation(sim *core.Simulation) {
	players := sim.Raid.AllPlayerUnits

	if !ai.ShouldMove(sim) {
		return
	}

	for i := 0; i < len(players); i++ {
		player := players[i]
		moveOut := player.DistanceFromTarget + 10
		moveIn := player.DistanceFromTarget - 10

		if moveOut <= 40 && moveIn >= 5 {
			player.MoveTo(moveOut, sim)
		} else if moveIn >= 5 {
			player.MoveTo(moveIn, sim)
		}
	}
}

func (ai *LightMovementAI) ShouldMove(sim *core.Simulation) bool {
	if sim.CurrentTime-ai.LastMoveTime >= 10*time.Second {
		ai.LastMoveTime = sim.CurrentTime
		return true
	}
	return false
}
