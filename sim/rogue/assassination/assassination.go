package assassination

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

// Damage Done By Caster setup
const (
	DDBC_Vendetta = iota

	DDBC_Total
)

func RegisterAssassinationRogue() {
	core.RegisterAgentFactory(
		proto.Player_AssassinationRogue{},
		proto.Spec_SpecAssassinationRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewAssassinationRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_AssassinationRogue)
			if !ok {
				panic("Invalid spec value for Assassination Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func (sinRogue *AssassinationRogue) Initialize() {
	sinRogue.Rogue.Initialize()

	sinRogue.MasteryBaseValue = 0.28
	sinRogue.MasteryMultiplier = 0.035

	sinRogue.registerMutilateSpell()
	sinRogue.registerDispatch()
	sinRogue.applySealFate()
	sinRogue.registerVenomousWounds()
	sinRogue.registerVendetta()
	sinRogue.registerEnvenom()
	sinRogue.registerAllPassives()
}

func NewAssassinationRogue(character *core.Character, options *proto.Player) *AssassinationRogue {
	sinOptions := options.GetAssassinationRogue().Options

	sinRogue := &AssassinationRogue{
		Rogue: rogue.NewRogue(character, sinOptions.ClassOptions, options.TalentsString),
	}
	sinRogue.AssassinationOptions = sinOptions

	return sinRogue
}

type AssassinationRogue struct {
	*rogue.Rogue
}

func (sinRogue *AssassinationRogue) GetRogue() *rogue.Rogue {
	return sinRogue.Rogue
}

func (sinRogue *AssassinationRogue) Reset(sim *core.Simulation) {
	sinRogue.Rogue.Reset(sim)
}
