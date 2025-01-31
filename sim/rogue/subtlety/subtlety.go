package subtlety

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

func RegisterSubtletyRogue() {
	core.RegisterAgentFactory(
		proto.Player_SubtletyRogue{},
		proto.Spec_SpecSubtletyRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewSubtletyRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SubtletyRogue)
			if !ok {
				panic("Invalid spec value for Subtlety Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func (subRogue *SubtletyRogue) Initialize() {
	subRogue.Rogue.Initialize()

	subRogue.MasteryBaseValue = 0.2
	subRogue.MasteryMultiplier = .025

	subRogue.registerHemorrhageSpell()
	subRogue.registerSanguinaryVein()
	subRogue.registerPremeditation()
	subRogue.registerHonorAmongThieves()

	subRogue.applyInitiative()
	subRogue.applyFindWeakness()

	subRogue.registerMasterOfSubtletyCD()
	subRogue.registerShadowDanceCD()
	subRogue.registerPreparationCD()
	subRogue.registerShadowstepCD()

	// Apply Mastery
	masteryMod := subRogue.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: rogue.RogueSpellRupture | rogue.RogueSpellEviscerate,
		IntValue:  0,
	})

	subRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryMod.UpdateIntValue(int64(subRogue.GetMasteryBonus() * 100))
	})

	core.MakePermanent(subRogue.GetOrRegisterAura(core.Aura{
		Label:    "Executioner",
		ActionID: core.ActionID{SpellID: 76808},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.UpdateIntValue(int64(subRogue.GetMasteryBonus() * 100))
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

}

func NewSubtletyRogue(character *core.Character, options *proto.Player) *SubtletyRogue {
	subOptions := options.GetSubtletyRogue().Options

	subRogue := &SubtletyRogue{
		Rogue: rogue.NewRogue(character, subOptions.ClassOptions, options.TalentsString),
	}
	subRogue.SubtletyOptions = subOptions

	subRogue.MultiplyStat(stats.Agility, 1.30)

	return subRogue
}

type SubtletyRogue struct {
	*rogue.Rogue
}

func (subRogue *SubtletyRogue) GetRogue() *rogue.Rogue {
	return subRogue.Rogue
}

func (subRogue *SubtletyRogue) Reset(sim *core.Simulation) {
	subRogue.Rogue.Reset(sim)
}
