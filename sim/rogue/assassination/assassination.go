package assassination

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

const masteryDamagePerPoint = 0.035
const masteryBaseEffect = 0.28

func RegisterAssassinationRogue() {
	core.RegisterAgentFactory(
		proto.Player_AssassinationRogue{},
		proto.Spec_SpecAssassinationRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewAssassinationRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_AssassinationRogue)
			if !ok {
				panic("Invalid spec value for Assassination Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func (sinRogue *AssassinationRogue) Initialize() {
	sinRogue.Rogue.Initialize()

	sinRogue.registerMutilateSpell()
	sinRogue.registerOverkill()
	sinRogue.registerColdBloodCD()
	sinRogue.applySealFate()
	sinRogue.registerVenomousWounds()
	sinRogue.registerVendetta()

	// Apply Mastery
	masteryEffect := getMasteryBonus(sinRogue.GetStat(stats.Mastery))
	for _, spell := range sinRogue.InstantPoison {
		spell.DamageMultiplier += masteryEffect
	}
	for _, spell := range sinRogue.WoundPoison {
		spell.DamageMultiplier += masteryEffect
	}
	sinRogue.DeadlyPoison.DamageMultiplier += masteryEffect
	sinRogue.Envenom.DamageMultiplier += masteryEffect
	if sinRogue.Talents.VenomousWounds > 0 {
		sinRogue.VenomousWounds.DamageMultiplier += masteryEffect
	}

	sinRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryEffectOld := getMasteryBonus(oldMastery)
		masteryEffectNew := getMasteryBonus(newMastery)
		for _, spell := range sinRogue.InstantPoison {
			spell.DamageMultiplier -= masteryEffectOld
			spell.DamageMultiplier += masteryEffectNew
		}
		for _, spell := range sinRogue.WoundPoison {
			spell.DamageMultiplier -= masteryEffectOld
			spell.DamageMultiplier += masteryEffectNew
		}
		sinRogue.DeadlyPoison.DamageMultiplier -= masteryEffectOld
		sinRogue.DeadlyPoison.DamageMultiplier += masteryEffectNew
		sinRogue.Envenom.DamageMultiplier -= masteryEffectOld
		sinRogue.Envenom.DamageMultiplier += masteryEffectNew
		if sinRogue.Talents.VenomousWounds > 0 {
			sinRogue.VenomousWounds.DamageMultiplier -= masteryEffectOld
			sinRogue.VenomousWounds.DamageMultiplier += masteryEffectNew
		}
	})

	// Assassin's Resolve: +20% physical damage
	// +20 Energy handled in base rogue
	if sinRogue.GetMHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger {
		sinRogue.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.2
	}
}

func getMasteryBonus(masteryRating float64) float64 {
	return masteryBaseEffect + core.MasteryRatingToMasteryPoints(masteryRating)*masteryDamagePerPoint
}

func NewAssassinationRogue(character *core.Character, options *proto.Player) *AssassinationRogue {
	sinOptions := options.GetAssassinationRogue().Options

	sinRogue := &AssassinationRogue{
		Rogue: rogue.NewRogue(character, sinOptions.ClassOptions, options.TalentsString),
	}
	sinRogue.AssassinationOptions = sinOptions

	return sinRogue
}

type AssassinationRogue struct {
	*rogue.Rogue
}

func (sinRogue *AssassinationRogue) GetRogue() *rogue.Rogue {
	return sinRogue.Rogue
}

func (sinRogue *AssassinationRogue) Reset(sim *core.Simulation) {
	sinRogue.Rogue.Reset(sim)
}
