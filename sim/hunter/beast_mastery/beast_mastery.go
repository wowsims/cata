package beast_mastery

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/hunter"
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

func (bmHunter *BeastMasteryHunter) Initialize() {
	// Initialize global Hunter spells

	bmHunter.Hunter.Initialize()

	// Apply BM Hunter mastery
	// baseMasteryRating := bmHunter.GetStat(stats.MasteryRating)
	// kcMod := bmHunter.AddDynamicMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_DamageDone_Pct,
	// 	ClassMask:  hunter.HunterSpellKillCommand,
	// 	FloatValue: bmHunter.getMasteryBonus(baseMasteryRating),
	// })

	// if bmHunter.Pet != nil {
	// 	bmHunter.Pet.PseudoStats.DamageDealtMultiplier *= bmHunter.getMasteryBonus(baseMasteryRating)
	// 	kcMod.Activate()
	// }

	// bmHunter.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
	// 	if bmHunter.Pet != nil {
	// 		bmHunter.Pet.PseudoStats.DamageDealtMultiplier /= bmHunter.getMasteryBonus(oldMasteryRating)
	// 		bmHunter.Pet.PseudoStats.DamageDealtMultiplier *= bmHunter.getMasteryBonus(newMasteryRating)
	// 		kcMod.UpdateFloatValue(bmHunter.getMasteryBonus(newMasteryRating))
	// 	}
	// })

	// BM Hunter Spec Bonus
	bmHunter.MultiplyStat(stats.RangedAttackPower, 1.30)
}
func (hunter *BeastMasteryHunter) ApplyTalents() {}

func (hunter *BeastMasteryHunter) getMasteryBonus(masteryRating float64) float64 {
	return 1.134 + ((masteryRating / core.MasteryRatingPerMasteryPoint) * 0.0167)
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
