package affliction

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
)

func RegisterAfflictionWarlock() {
	core.RegisterAgentFactory(
		proto.Player_AfflictionWarlock{},
		proto.Spec_SpecAfflictionWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewAfflictionWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_AfflictionWarlock)
			if !ok {
				panic("Invalid spec value for Affliction Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewAfflictionWarlock(character *core.Character, options *proto.Player) *AfflictionWarlock {
	affOptions := options.GetAfflictionWarlock().Options
	affliction := &AfflictionWarlock{
		Warlock: warlock.NewWarlock(character, options, affOptions.ClassOptions),
	}

	return affliction
}

type AfflictionWarlock struct {
	*warlock.Warlock

	SoulShards         core.SecondaryResourceBar
	Agony              *core.Spell
	UnstableAffliction *core.Spell
	SoulBurnAura       *core.Aura
	HauntDebuffAuras   core.AuraArray
	LastCorruption     *core.Dot // Tracks the last corruption we've applied
	ProcMaleficEffect  func(target *core.Unit, coeff float64, sim *core.Simulation)
	HauntImpactTime    time.Duration
}

func (affliction AfflictionWarlock) getMasteryBonus() float64 {
	return (8 + affliction.GetMasteryPoints()) * 3.1
}

func (affliction *AfflictionWarlock) GetWarlock() *warlock.Warlock {
	return affliction.Warlock
}

func (affliction *AfflictionWarlock) Initialize() {
	affliction.Warlock.Initialize()

	affliction.SoulShards = affliction.RegisterNewDefaultSecondaryResourceBar(core.SecondaryResourceConfig{
		Type:    proto.SecondaryResourceType_SecondaryResourceTypeSoulShards,
		Max:     4,
		Default: 4,
	})

	affliction.registerPotentAffliction()
	affliction.registerHaunt()
	corruption := affliction.RegisterCorruption(func(resultList []core.SpellResult, spell *core.Spell, sim *core.Simulation) {
		if resultList[0].Landed() {
			affliction.LastCorruption = spell.Dot(resultList[0].Target)
		}
	})

	// June 16th Beta Changes +33% for affliction
	corruption.DamageMultiplier *= 1.33

	affliction.registerAgony()
	affliction.registerNightfall()
	affliction.registerUnstableAffliction()
	affliction.registerMaleficEffect()
	affliction.registerMaleficGrasp()
	affliction.registerDrainSoul()
	affliction.registerDarkSoulMisery()
	affliction.registerSoulburn()
	affliction.registerSeed()
	affliction.registerSoulSwap()

	affliction.registerGlpyhs()
}

func (affliction *AfflictionWarlock) ApplyTalents() {
	affliction.Warlock.ApplyTalents()
}

func (affliction *AfflictionWarlock) Reset(sim *core.Simulation) {
	affliction.Warlock.Reset(sim)

	affliction.LastCorruption = nil
	affliction.HauntImpactTime = 0
}

func calculateDoTBaseTickDamage(dot *core.Dot) float64 {
	stacks := math.Max(float64(dot.Aura.GetStacks()), 1)
	return dot.SnapshotBaseDamage * dot.SnapshotAttackerMultiplier * stacks
}
