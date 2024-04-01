package shadow

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/priest"
)

func RegisterShadowPriest() {
	core.RegisterAgentFactory(
		proto.Player_ShadowPriest{},
		proto.Spec_SpecShadowPriest,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewShadowPriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ShadowPriest)
			if !ok {
				panic("Invalid spec value for Shadow Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewShadowPriest(character *core.Character, options *proto.Player) *ShadowPriest {
	shadowOptions := options.GetShadowPriest()

	selfBuffs := priest.SelfBuffs{
		UseShadowfiend: true,
		UseInnerFire:   shadowOptions.Options.ClassOptions.Armor == proto.PriestOptions_InnerFire,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	basePriest.Latency = float64(basePriest.ChannelClipDelay.Milliseconds())
	spriest := &ShadowPriest{
		Priest:  basePriest,
		options: shadowOptions.Options,
	}

	// add spirit -> spell hit conversion for Twisted Faith talent
	if spriest.Talents.TwistedFaith > 0 {
		spriest.AddStatDependency(stats.Spirit, stats.SpellHit, 0.5*float64(spriest.Talents.TwistedFaith))
	}

	spriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	if spriest.Talents.PowerInfusion && shadowOptions.Options.PowerInfusionTarget != nil {
		spriest.SelfBuffs.PowerInfusionTarget = shadowOptions.Options.PowerInfusionTarget
	}

	return spriest
}

type ShadowPriest struct {
	*priest.Priest
	options *proto.ShadowPriest_Options
}

func (spriest *ShadowPriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *ShadowPriest) Initialize() {
	spriest.Priest.Initialize()
}

func (spriest *ShadowPriest) Reset(sim *core.Simulation) {
	spriest.Priest.Reset(sim)
}
