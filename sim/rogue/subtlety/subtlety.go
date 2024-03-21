package subtlety

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

const masteryDamagePerPoint = .01
const masteryBaseEffect = 0.2

var hasT6 bool = false

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
	hasT6 = subRogue.HasSetBonus(rogue.Tier6, 2)

	masteryPoints := subRogue.GetStat(stats.Mastery) / core.MasteryRatingPerMasteryPoint
	masteryEffect := masteryPoints*masteryDamagePerPoint + masteryBaseEffect
	subRogue.SliceAndDiceBonus = 1.4
	subRogue.SliceAndDiceBonus *= (1 + masteryEffect)
	if hasT6 {
		subRogue.SliceAndDiceBonus += 0.05
	}
	subRogue.Eviscerate.DamageMultiplierAdditive += masteryEffect
	subRogue.Rupture.DamageMultiplierAdditive += masteryEffect

	subRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryPointsOld := oldMastery / core.MasteryRatingPerMasteryPoint
		masteryPointsNew := newMastery / core.MasteryRatingPerMasteryPoint
		masteryEffect := (masteryPointsNew-masteryPointsOld)*masteryDamagePerPoint + masteryBaseEffect
		subRogue.SliceAndDiceBonus = 1.4
		subRogue.SliceAndDiceBonus *= (1 + masteryEffect)
		if hasT6 {
			subRogue.SliceAndDiceBonus += 0.05
		}
		subRogue.Eviscerate.DamageMultiplierAdditive += masteryEffect
		subRogue.Rupture.DamageMultiplierAdditive += masteryEffect
	})
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
