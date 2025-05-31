package survival

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/hunter"
)

func RegisterSurvivalHunter() {
	core.RegisterAgentFactory(
		proto.Player_SurvivalHunter{},
		proto.Spec_SpecSurvivalHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewSurvivalHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SurvivalHunter)
			if !ok {
				panic("Invalid spec value for Survival Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func (hunter *SurvivalHunter) Initialize() {
	// Initialize global Hunter spells
	hunter.Hunter.Initialize()
	hunter.registerExplosiveShotSpell()
	hunter.registerBlackArrowSpell(hunter.FireTrapTimer)
	// Apply SV Hunter mastery
	schoolsAffectedBySurvivalMastery := []stats.SchoolIndex{
		stats.SchoolIndexNature,
		stats.SchoolIndexFire,
		stats.SchoolIndexArcane,
		stats.SchoolIndexFrost,
		stats.SchoolIndexShadow,
	}
	baseMasteryRating := hunter.GetStat(stats.MasteryRating)
	for _, school := range schoolsAffectedBySurvivalMastery {
		hunter.PseudoStats.SchoolDamageDealtMultiplier[school] *= hunter.getMasteryBonus(baseMasteryRating)
	}

	hunter.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
		for _, school := range schoolsAffectedBySurvivalMastery {
			hunter.PseudoStats.SchoolDamageDealtMultiplier[school] /= hunter.getMasteryBonus(oldMasteryRating)
			hunter.PseudoStats.SchoolDamageDealtMultiplier[school] *= hunter.getMasteryBonus(newMasteryRating)
		}
	})

	// Survival Spec Bonus
	//hunter.MultiplyStat(stats.Agility, 1.1)
}
func (hunter *SurvivalHunter) getMasteryBonus(masteryRating float64) float64 {
	return 1.08 + ((masteryRating / core.MasteryRatingPerMasteryPoint) * 0.01)
}

func NewSurvivalHunter(character *core.Character, options *proto.Player) *SurvivalHunter {
	survivalOptions := options.GetSurvivalHunter().Options

	svHunter := &SurvivalHunter{
		Hunter: hunter.NewHunter(character, options, survivalOptions.ClassOptions),
	}

	svHunter.SurvivalOptions = survivalOptions

	return svHunter
}

type SurvivalHunter struct {
	*hunter.Hunter
}

func (svHunter *SurvivalHunter) GetHunter() *hunter.Hunter {
	return svHunter.Hunter
}

func (svHunter *SurvivalHunter) Reset(sim *core.Simulation) {
	svHunter.Hunter.Reset(sim)
}
