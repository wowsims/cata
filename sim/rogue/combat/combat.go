package combat

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func RegisterCombatRogue() {
	core.RegisterAgentFactory(
		proto.Player_CombatRogue{},
		proto.Spec_SpecCombatRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewCombatRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_CombatRogue)
			if !ok {
				panic("Invalid spec value for Combat Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewCombatRogue(character *core.Character, options *proto.Player) *CombatRogue {
	combatOptions := options.GetCombatRogue().Options

	combatRogue := &CombatRogue{
		Rogue: rogue.NewRogue(character, combatOptions.ClassOptions, options.TalentsString),
	}
	combatRogue.CombatOptions = combatOptions

	return combatRogue
}

type CombatRogue struct {
	*rogue.Rogue
}

func (combatRogue *CombatRogue) GetRogue() *rogue.Rogue {
	return combatRogue.Rogue
}

func (combatRogue *CombatRogue) Reset(sim *core.Simulation) {
	combatRogue.Rogue.Reset(sim)
}
