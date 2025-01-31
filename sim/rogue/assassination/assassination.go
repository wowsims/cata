package assassination

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

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

	sinRogue.MasteryBaseValue = 0.28
	sinRogue.MasteryMultiplier = 0.035

	sinRogue.registerMutilateSpell()
	sinRogue.registerOverkill()
	sinRogue.registerColdBloodCD()
	sinRogue.applySealFate()
	sinRogue.registerVenomousWounds()
	sinRogue.registerVendetta()

	// Apply Mastery
	// As far as I am able to find, Asn's Mastery is an additive bonus. To be tested.
	masteryEffect := int64(sinRogue.GetMasteryBonusFromRating(sinRogue.GetStat(stats.MasteryRating)) * 100)

	for _, spell := range sinRogue.InstantPoison {
		spell.ApplyDamageMultiplierAdditive(masteryEffect)
	}
	for _, spell := range sinRogue.WoundPoison {
		spell.ApplyDamageMultiplierAdditive(masteryEffect)
	}
	sinRogue.DeadlyPoison.ApplyDamageMultiplierAdditive(masteryEffect)
	sinRogue.Envenom.ApplyDamageMultiplierAdditive(masteryEffect)
	if sinRogue.Talents.VenomousWounds > 0 {
		sinRogue.VenomousWounds.ApplyDamageMultiplierAdditive(masteryEffect)
	}

	sinRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryEffectOld := int64(sinRogue.GetMasteryBonusFromRating(oldMastery) * 100)
		masteryEffectNew := int64(sinRogue.GetMasteryBonusFromRating(newMastery) * 100)

		for _, spell := range sinRogue.InstantPoison {
			spell.ApplyDamageMultiplierAdditive(masteryEffectNew - masteryEffectOld)
		}
		for _, spell := range sinRogue.WoundPoison {
			spell.ApplyDamageMultiplierAdditive(masteryEffectNew - masteryEffectOld)
		}
		sinRogue.DeadlyPoison.ApplyDamageMultiplierAdditive(masteryEffectNew - masteryEffectOld)
		sinRogue.Envenom.ApplyDamageMultiplierAdditive(masteryEffectNew - masteryEffectOld)
		if sinRogue.Talents.VenomousWounds > 0 {
			sinRogue.VenomousWounds.ApplyDamageMultiplierAdditive(masteryEffectNew - masteryEffectOld)
		}
	})

	// Assassin's Resolve: +20% Multiplicative physical damage (confirmed)
	// +20 Energy handled in base rogue
	if sinRogue.GetMHWeapon() != nil && sinRogue.GetMHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger {
		sinRogue.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.2
	}
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
