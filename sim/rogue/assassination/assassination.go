package assassination

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/rogue"
)

// Damage Done By Caster setup
const (
	DDBC_Vendetta = iota

	DDBC_Total
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

	// sinRogue.registerMutilateSpell()
	// sinRogue.registerOverkill()
	// sinRogue.registerColdBloodCD()
	// sinRogue.applySealFate()
	// sinRogue.registerVenomousWounds()
	// sinRogue.registerVendetta()

	// Apply Mastery
	// As far as I am able to find, Asn's Mastery is an additive bonus. To be tested.
	masteryMod := sinRogue.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  rogue.RogueSpellInstantPoison | rogue.RogueSpellWoundPoison | rogue.RogueSpellDeadlyPoison | rogue.RogueSpellEnvenom | rogue.RogueSpellVenomousWounds,
		FloatValue: sinRogue.GetMasteryBonusFromRating(sinRogue.GetStat(stats.MasteryRating)),
	})
	masteryMod.Activate()

	sinRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryMod.UpdateFloatValue(sinRogue.GetMasteryBonusFromRating(newMastery))
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
