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

func (hunter *SurvivalHunter) applyMastery() {
	hunter.RegisterAura(core.Aura{
		Label:    "Essence of the Viper",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		//Todo: Change to OnMasteryChanged when available
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			multiplier := 1.08 + (hunter.Hunter.CalculateMasteryPoints() * 0.01)
			hunter.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] = multiplier
			hunter.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] = multiplier
		},
	})
}
func (hunter *SurvivalHunter) Initialize() {
	// Initialize global Hunter spells
	hunter.Hunter.Initialize()

	// Spec specific spells
	hunter.applyMastery()
	hunter.registerExplosiveShotSpell()
	hunter.registerBlackArrowSpell(hunter.FireTrapTimer)
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
