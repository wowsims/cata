package beast_mastery

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/hunter"
)

func RegisterBeastMasteryHunter() {
	core.RegisterAgentFactory(
		proto.Player_BeastMasteryHunter{},
		proto.Spec_SpecBeastMasteryHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewBeastMasteryHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BeastMasteryHunter)
			if !ok {
				panic("Invalid spec value for Beast Mastery Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewBeastMasteryHunter(character *core.Character, options *proto.Player) *BeastMasteryHunter {
	bmOptions := options.GetBeastMasteryHunter().Options

	bmHunter := &BeastMasteryHunter{
		Hunter: hunter.NewHunter(character, options, bmOptions.ClassOptions),
	}
	bmHunter.BeastMasteryOptions = bmOptions

	return bmHunter
}

func (hunter *BeastMasteryHunter) Initialize() {
	// Initialize global Hunter spells

	hunter.Hunter.Initialize()

	// Apply BM Hunter mastery
	baseMastery := hunter.GetStat(stats.Mastery)
	if hunter.Pet != nil {
		hunter.Pet.PseudoStats.DamageDealtMultiplier *= hunter.getMasteryBonus(baseMastery)
	}

	hunter.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		if hunter.Pet != nil {
			hunter.Pet.PseudoStats.DamageDealtMultiplier /= hunter.getMasteryBonus(oldMastery)
			hunter.Pet.PseudoStats.DamageDealtMultiplier *= hunter.getMasteryBonus(newMastery)
		}
	})

	// BM Hunter Spec Bonus
	hunter.MultiplyStat(stats.AttackPower, 1.30)
}
func (hunter *BeastMasteryHunter) getMasteryBonus(mastery float64) float64 {
	return 1.134 + ((mastery / core.MasteryRatingPerMasteryPoint) * 0.0167)
}

type BeastMasteryHunter struct {
	*hunter.Hunter
}

func (bmHunter *BeastMasteryHunter) GetHunter() *hunter.Hunter {
	return bmHunter.Hunter
}

func (bmHunter *BeastMasteryHunter) Reset(sim *core.Simulation) {
	bmHunter.Hunter.Reset(sim)
}
