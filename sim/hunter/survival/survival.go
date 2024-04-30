package survival

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/hunter"
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
	baseMastery := hunter.GetStat(stats.Mastery)
	for _, school := range schoolsAffectedBySurvivalMastery {
		hunter.PseudoStats.SchoolDamageDealtMultiplier[school] *= hunter.getMasteryBonus(baseMastery)
	}

	hunter.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		for _, school := range schoolsAffectedBySurvivalMastery {
			hunter.PseudoStats.SchoolDamageDealtMultiplier[school] /= hunter.getMasteryBonus(oldMastery)
			hunter.PseudoStats.SchoolDamageDealtMultiplier[school] *= hunter.getMasteryBonus(newMastery)
		}
	})

	// Survival Spec Bonus
	hunter.MultiplyStat(stats.Agility, 1.1)
}
func (hunter *SurvivalHunter) getMasteryBonus(mastery float64) float64 {
	return 1.08 + ((mastery / core.MasteryRatingPerMasteryPoint) * 0.01)
}

func NewSurvivalHunter(character *core.Character, options *proto.Player) *SurvivalHunter {
	survivalOptions := options.GetSurvivalHunter().Options

	svHunter := &SurvivalHunter{
		Hunter: hunter.NewHunter(character, options, survivalOptions.ClassOptions),
	}

	svHunter.SurvivalOptions = survivalOptions
	// Todo: Is there a better way to do this?

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
