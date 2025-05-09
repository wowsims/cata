package windwalker

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/monk"
)

// Damage Done By Caster setup
const (
	DDBC_RisingSunKick int = iota

	DDBC_Total
)

func RegisterWindwalkerMonk() {
	core.RegisterAgentFactory(
		proto.Player_WindwalkerMonk{},
		proto.Spec_SpecWindwalkerMonk,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewWindwalkerMonk(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_WindwalkerMonk)
			if !ok {
				panic("Invalid spec value for Windwalker Monk!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewWindwalkerMonk(character *core.Character, options *proto.Player) *WindwalkerMonk {
	monkOptions := options.GetWindwalkerMonk()

	ww := &WindwalkerMonk{
		Monk: monk.NewMonk(character, monkOptions.Options.ClassOptions, options.TalentsString),
	}

	ww.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	ww.AddStatDependency(stats.Agility, stats.AttackPower, 2)

	return ww
}

type WindwalkerMonk struct {
	*monk.Monk

	TigereyeBrewStackAura *core.Aura

	outstandingChi int32
}

func (ww *WindwalkerMonk) GetMonk() *monk.Monk {
	return ww.Monk
}

func (ww *WindwalkerMonk) Initialize() {
	ww.Monk.Initialize()
	ww.RegisterSpecializationEffects()
}

func (ww *WindwalkerMonk) ApplyTalents() {
	ww.Monk.ApplyTalents()
	ww.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeLeather, 120227)
}

func (ww *WindwalkerMonk) Reset(sim *core.Simulation) {
	ww.Monk.Reset(sim)
}

func (ww *WindwalkerMonk) RegisterSpecializationEffects() {
	ww.registerEnergizingBrew()
	ww.registerFistsOfFury()
	ww.registerPassives()
	ww.registerRisingSunKick()
	ww.registerTigereyeBrew()
	ww.registerSpinningFireBlossom()
}

func (ww *WindwalkerMonk) getMasteryPercent() float64 {
	return (8.0 + ww.GetMasteryPoints()) * 0.025
}
