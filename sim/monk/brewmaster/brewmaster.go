package brewmaster

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/monk"
)

func RegisterBrewmasterMonk() {
	core.RegisterAgentFactory(
		proto.Player_BrewmasterMonk{},
		proto.Spec_SpecBrewmasterMonk,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewBrewmasterMonk(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BrewmasterMonk)
			if !ok {
				panic("Invalid spec value for Brewmaster Monk!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewBrewmasterMonk(character *core.Character, options *proto.Player) *BrewmasterMonk {
	monkOptions := options.GetBrewmasterMonk()

	bm := &BrewmasterMonk{
		Monk:           monk.NewMonk(character, monkOptions.Options.ClassOptions, options.TalentsString),
		StartingStance: monkOptions.Options.Stance,
	}
	bm.SetStartingStance()

	bm.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	bm.AddStatDependency(stats.Agility, stats.AttackPower, 2)

	return bm
}

func (bm *BrewmasterMonk) SetStartingStance() {
	switch bm.StartingStance {
	case proto.MonkStance_SturdyOx:
		bm.Stance = monk.SturdyOx
	case proto.MonkStance_FierceTiger:
		bm.Stance = monk.FierceTiger
	}
}

type BrewmasterMonk struct {
	*monk.Monk
	StartingStance proto.MonkStance
}

func (bm *BrewmasterMonk) GetMonk() *monk.Monk {
	return bm.Monk
}

func (bm *BrewmasterMonk) Initialize() {
	bm.Monk.Initialize()
	bm.RegisterSpecializationEffects()
}

func (bm *BrewmasterMonk) ApplyTalents() {
	bm.Monk.ApplyTalents()
	bm.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypeLeather, 120225)
}

func (bm *BrewmasterMonk) Reset(sim *core.Simulation) {
	bm.SetStartingStance()
	bm.Monk.Reset(sim)
}

func (bm *BrewmasterMonk) RegisterSpecializationEffects() {
	bm.RegisterMastery()
}

func (bm *BrewmasterMonk) RegisterMastery() {
}
