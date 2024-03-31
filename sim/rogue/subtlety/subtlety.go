package subtlety

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

const masteryDamagePerPoint = .025
const masteryBaseEffect = 0.2

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

	subRogue.registerHemorrhageSpell()
	subRogue.registerSanguinaryVein()
	subRogue.registerPremeditation()
	subRogue.registerHonorAmongThieves()

	subRogue.applyInitiative()
	subRogue.applySerratedBlades()
	subRogue.applyFindWeakness()

	subRogue.registerMasterOfSubtletyCD()
	subRogue.registerShadowDanceCD()
	subRogue.registerPreparationCD()
	subRogue.registerShadowstepCD()

	// Apply Mastery
	// From all I can find, Sub's Mastery is Additive. Will need to test.
	masteryEffect := getMasteryBonus(subRogue.GetStat(stats.Mastery))

	subRogue.SliceAndDiceBonus *= (1 + masteryEffect)
	subRogue.Eviscerate.DamageMultiplierAdditive += masteryEffect
	subRogue.Rupture.DamageMultiplierAdditive += masteryEffect

	subRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryEffectOld := getMasteryBonus(oldMastery)
		masteryEffectNew := getMasteryBonus(newMastery)

		subRogue.SliceAndDiceBonus /= (1 + masteryEffectOld)
		subRogue.SliceAndDiceBonus *= (1 + masteryEffectNew)
		subRogue.Eviscerate.DamageMultiplierAdditive -= masteryEffectOld
		subRogue.Eviscerate.DamageMultiplierAdditive += masteryEffectNew
		subRogue.Rupture.DamageMultiplierAdditive -= masteryEffectOld
		subRogue.Rupture.DamageMultiplierAdditive += masteryEffectNew
	})
}

func getMasteryBonus(masteryRating float64) float64 {
	return masteryBaseEffect + core.MasteryRatingToMasteryPoints(masteryRating)*masteryDamagePerPoint
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
