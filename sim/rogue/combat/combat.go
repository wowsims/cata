package combat

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/rogue"
)

// Damage Done By Caster setup
const (
	DDBC_BanditsGuile    int = 0
	DDBC_RevealingStrike     = iota

	DDBC_Total
)

func RegisterCombatRogue() {
	core.RegisterAgentFactory(
		proto.Player_CombatRogue{},
		proto.Spec_SpecCombatRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewCombatRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_CombatRogue)
			if !ok {
				panic("Invalid spec value for Combat Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewCombatRogue(character *core.Character, options *proto.Player) *CombatRogue {
	combatOptions := options.GetCombatRogue().Options

	combatRogue := &CombatRogue{
		Rogue: rogue.NewRogue(character, combatOptions.ClassOptions, options.TalentsString),
	}
	combatRogue.CombatOptions = combatOptions

	return combatRogue
}

func (combatRogue *CombatRogue) Initialize() {
	combatRogue.Rogue.Initialize()

	combatRogue.MasteryBaseValue = 0.16
	combatRogue.MasteryMultiplier = 0.02

	// Ambidexterity Passive
	combatRogue.AutoAttacks.OHConfig().DamageMultiplier *= 1.75
	// Vitality Passive
	combatRogue.AdditiveEnergyRegenBonus += 0.20
	combatRogue.MultiplyStat(stats.AttackPower, 1.4)

	combatRogue.registerSinisterStrikeSpell()
	combatRogue.registerRevealingStrike()
	combatRogue.registerBladeFlurry()
	combatRogue.registerBanditsGuile()

	combatRogue.applyCombatPotency()

	combatRogue.registerKillingSpreeCD()
	combatRogue.registerAdrenalineRushCD()

	combatRogue.applyMastery()
}

type CombatRogue struct {
	*rogue.Rogue
}

func (combatRogue *CombatRogue) GetRogue() *rogue.Rogue {
	return combatRogue.Rogue
}

func (combatRogue *CombatRogue) Reset(sim *core.Simulation) {
	combatRogue.Rogue.Reset(sim)

	combatRogue.BanditsGuileAura.Activate(sim)
}
