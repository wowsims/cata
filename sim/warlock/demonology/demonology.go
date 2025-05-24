package demonology

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
)

func RegisterDemonologyWarlock() {
	core.RegisterAgentFactory(
		proto.Player_DemonologyWarlock{},
		proto.Spec_SpecDemonologyWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDemonologyWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DemonologyWarlock)
			if !ok {
				panic("Invalid spec value for Demonology Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDemonologyWarlock(character *core.Character, options *proto.Player) *DemonologyWarlock {
	demoOptions := options.GetDemonologyWarlock().Options

	demonology := &DemonologyWarlock{
		Warlock: warlock.NewWarlock(character, options, demoOptions.ClassOptions),
	}

	demonology.Felguard = demonology.registerFelguard()
	return demonology
}

type DemonologyWarlock struct {
	*warlock.Warlock

	DemonicFury   core.SecondaryResourceBar
	Metamorphosis *core.Spell
	Felguard      *warlock.WarlockPet
}

func (demonology *DemonologyWarlock) GetWarlock() *warlock.Warlock {
	return demonology.Warlock
}

func (demonology *DemonologyWarlock) Initialize() {
	demonology.Warlock.Initialize()

	demonology.DemonicFury = demonology.RegisterNewDefaultSecondaryResourceBar(core.SecondaryResourceConfig{
		Type:    proto.SecondaryResourceType_SecondaryResourceTypeDemonicFury,
		Max:     1000,
		Default: 200,
	})

	demonology.registerMetamorphosis()
	demonology.registerMasterDemonologist()
	demonology.registerShadowBolt()
	demonology.registerFelFlame()
	demonology.registerCorruption()
	demonology.registerDrainLife()
	demonology.registerHandOfGuldan()
	demonology.registerHellfire()

	// demonology.registerHandOfGuldan()
}

func (demonology *DemonologyWarlock) ApplyTalents() {
	demonology.Warlock.ApplyTalents()
}

func (demonology *DemonologyWarlock) Reset(sim *core.Simulation) {
	demonology.Warlock.Reset(sim)
}

func NewDemonicFuryCost(cost int) *warlock.SecondaryResourceCost {
	return &warlock.SecondaryResourceCost{
		SecondaryCost: cost,
		Name:          "Demonic Fury",
	}
}

func (demo *DemonologyWarlock) IsInMeta() bool {
	return demo.Metamorphosis.RelatedSelfBuff.IsActive()
}
