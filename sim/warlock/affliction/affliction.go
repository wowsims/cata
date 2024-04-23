package affliction

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func RegisterAfflictionWarlock() {
	core.RegisterAgentFactory(
		proto.Player_AfflictionWarlock{},
		proto.Spec_SpecAfflictionWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewAfflictionWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_AfflictionWarlock)
			if !ok {
				panic("Invalid spec value for Affliction Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewAfflictionWarlock(character *core.Character, options *proto.Player) *AfflictionWarlock {
	affOptions := options.GetAfflictionWarlock().Options
	affliction := &AfflictionWarlock{
		Warlock: warlock.NewWarlock(character, options, affOptions.ClassOptions),
	}

	return affliction
}

type AfflictionWarlock struct {
	*warlock.Warlock
}

func (affliction AfflictionWarlock) getMasteryBonus() float64 {
	return 0.25 + 0.0163*affliction.GetMasteryPoints()
}

func (affliction *AfflictionWarlock) GetWarlock() *warlock.Warlock {
	return affliction.Warlock
}

func (affliction *AfflictionWarlock) Initialize() {
	affliction.Warlock.Initialize()

	//affliction.registerUnstableAfflictionSpell()
}

func (affliction *AfflictionWarlock) ApplyTalents() {
	affliction.Warlock.ApplyTalents()

	// TODO: Mastery: Potent Afflictions

	// Shadow Mastery
	affliction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  warlock.WarlockShadowDamage,
		FloatValue: 0.30,
	})
}

func (affliction *AfflictionWarlock) Reset(sim *core.Simulation) {
	affliction.Warlock.Reset(sim)
}
