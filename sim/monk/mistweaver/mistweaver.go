package mistweaver

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/monk"
)

func RegisterMistweaverMonk() {
	core.RegisterAgentFactory(
		proto.Player_MistweaverMonk{},
		proto.Spec_SpecMistweaverMonk,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewMistweaverMonk(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_MistweaverMonk)
			if !ok {
				panic("Invalid spec value for Mistweaver Monk!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewMistweaverMonk(character *core.Character, options *proto.Player) *MistweaverMonk {
	monkOptions := options.GetMistweaverMonk()

	mw := &MistweaverMonk{
		Monk:           monk.NewMonk(character, monkOptions.Options.ClassOptions, options.TalentsString),
		StartingStance: monkOptions.Options.Stance,
	}
	mw.SetStartingStance()
	mw.EnableManaBar()

	strAPDep := mw.NewDynamicStatDependency(stats.Strength, stats.AttackPower, 1)
	agiAPDep := mw.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 2)

	mw.RegisterOnStanceChanged(func(sim *core.Simulation, newStance monk.Stance) {
		if newStance == monk.FierceTiger {
			mw.EnableDynamicStatDep(sim, strAPDep)
			mw.EnableDynamicStatDep(sim, agiAPDep)
		} else {
			mw.DisableDynamicStatDep(sim, strAPDep)
			mw.DisableDynamicStatDep(sim, agiAPDep)
		}
	})

	return mw
}

func (mw *MistweaverMonk) SetStartingStance() {
	switch mw.StartingStance {
	case proto.MonkStance_WiseSerpent:
		mw.Stance = monk.WiseSerpent
	case proto.MonkStance_FierceTiger:
		mw.Stance = monk.FierceTiger
	}
}

type MistweaverMonk struct {
	*monk.Monk
	StartingStance proto.MonkStance
}

func (mw *MistweaverMonk) GetMonk() *monk.Monk {
	return mw.Monk
}

func (mw *MistweaverMonk) Initialize() {
	mw.Monk.Initialize()
	mw.RegisterSpecializationEffects()
}

func (mw *MistweaverMonk) ApplyTalents() {
	mw.Monk.ApplyTalents()
	mw.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeLeather, 120224)
}

func (mw *MistweaverMonk) Reset(sim *core.Simulation) {
	mw.SetStartingStance()
	mw.Monk.Reset(sim)
}

func (mw *MistweaverMonk) RegisterSpecializationEffects() {
	mw.RegisterMastery()
}

func (mw *MistweaverMonk) RegisterMastery() {
}
