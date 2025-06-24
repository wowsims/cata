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
	bmHunter.registerBestialWrathCD()
	bmHunter.registerKillCommandSpell()
	bmHunter.registerFocusFireSpell()

	// Apply BM Hunter mastery
	baseMasteryRating := bmHunter.GetStat(stats.MasteryRating)
	baseMasteryBonus := bmHunter.getMasteryBonus(baseMasteryRating)

	var petMod *core.SpellMod
	if bmHunter.Pet != nil {
		petMod = bmHunter.Pet.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: baseMasteryBonus,
		})
		petMod.Activate()
	}

	var stampedePetMods []*core.SpellMod = make([]*core.SpellMod, len(bmHunter.StampedePet))
	for i := range bmHunter.StampedePet {
		if bmHunter.StampedePet[i] != nil {
			stampedePetMods[i] = bmHunter.StampedePet[i].AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: baseMasteryBonus,
			})
			stampedePetMods[i].Activate()
		}
	}

	var direBeastPetMod *core.SpellMod
	if bmHunter.DireBeastPet != nil {
		direBeastPetMod = bmHunter.DireBeastPet.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: baseMasteryBonus,
		})
		direBeastPetMod.Activate()
	}

	var amocDamageMod *core.SpellMod
	if bmHunter.Talents.AMurderOfCrows {
		amocDamageMod = bmHunter.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  hunter.HunterSpellAMurderOfCrows,
			FloatValue: baseMasteryBonus,
		})
		amocDamageMod.Activate()
	}

	bmHunter.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
		masteryBonus := bmHunter.getMasteryBonus(newMasteryRating)
		if petMod != nil {
			petMod.UpdateFloatValue(masteryBonus)
		}

		for i := range bmHunter.StampedePet {
			if stampedePetMods[i] != nil {
				stampedePetMods[i].UpdateFloatValue(masteryBonus)
			}
		}

		if direBeastPetMod != nil {
			direBeastPetMod.UpdateFloatValue(masteryBonus)
		}

		if amocDamageMod != nil {
			amocDamageMod.UpdateFloatValue(masteryBonus)
		}
	})
}

func (hunter *BeastMasteryHunter) getMasteryBonus(masteryRating float64) float64 {
	return 0.16 + ((masteryRating / core.MasteryRatingPerMasteryPoint) * 0.02)
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
