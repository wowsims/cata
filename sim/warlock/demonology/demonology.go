package demonology

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
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

	return demonology
}

type DemonologyWarlock struct {
	*warlock.Warlock
}

func (demonology DemonologyWarlock) getMasteryBonus() float64 {
	return 0.18 + 0.023*demonology.GetMasteryPoints()
}

func (demonology *DemonologyWarlock) GetWarlock() *warlock.Warlock {
	return demonology.Warlock
}

func (demonology *DemonologyWarlock) Initialize() {
	demonology.Warlock.Initialize()

	demonology.registerHandOfGuldanSpell()
}

func (demonology *DemonologyWarlock) ApplyTalents() {
	demonology.Warlock.ApplyTalents()

	// TODO: Mastery: Master Demonologist
	// This needs to affect your own damage during metamorphosis and your pets damage at all times
	masteryMod := demonology.AddDynamicMod(core.SpellModConfig{})

	demonology.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		masteryMod.UpdateFloatValue(demonology.getMasteryBonus())
	})

	core.MakePermanent(demonology.GetOrRegisterAura(core.Aura{
		Label:    "Mastery: Master Demonologist",
		ActionID: core.ActionID{SpellID: 77219},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.UpdateFloatValue(demonology.getMasteryBonus())
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

	// Demonic Knowledge
	demonology.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  warlock.WarlockShadowDamage | warlock.WarlockFireDamage,
		FloatValue: 0.15,
	})
}

func (demonology *DemonologyWarlock) Reset(sim *core.Simulation) {
	demonology.Warlock.Reset(sim)
}
