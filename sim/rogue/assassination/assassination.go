package assassination

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

const masteryDamagePerPercent = .035
const masteryBaseEffect = .28

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
	masteryPercent := sinRogue.GetStat(stats.Mastery) / core.MasteryRatingPerMasteryPoint
	masteryEffect := masteryBaseEffect + masteryPercent*masteryDamagePerPercent
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
		masteryPercentOld := oldMastery / core.MasteryRatingPerMasteryPoint
		masteryPercentNew := newMastery / core.MasteryRatingPerMasteryPoint
		masteryEffectChange := masteryBaseEffect + (masteryPercentNew-masteryPercentOld)*masteryDamagePerPercent
		for _, spell := range sinRogue.InstantPoison {
			spell.DamageMultiplier += masteryEffectChange
		}
		for _, spell := range sinRogue.WoundPoison {
			spell.DamageMultiplier += masteryEffectChange
		}
		sinRogue.DeadlyPoison.DamageMultiplier += masteryEffectChange
		sinRogue.Envenom.DamageMultiplier += masteryEffectChange
		if sinRogue.Talents.VenomousWounds > 0 {
			sinRogue.VenomousWounds.DamageMultiplier += masteryEffectChange
		}
	})

	// Assassin's Resolve: 20% additive melee damage
	// +20 Energy handled in base rogue
	if sinRogue.GetMHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger &&
		sinRogue.GetOHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger {
		for _, spell := range sinRogue.Spellbook {
			if spell.Flags.Matches(rogue.SpellFlagBuilder | rogue.SpellFlagFinisher) {
				spell.DamageMultiplierAdditive += 0.2
			}
		}
		sinRogue.AutoAttacks.MHConfig().DamageMultiplierAdditive += 0.2
		sinRogue.AutoAttacks.OHConfig().DamageMultiplierAdditive += 0.2
		sinRogue.AutoAttacks.RangedConfig().DamageMultiplierAdditive += 0.2
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
