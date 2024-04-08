package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (shaman *Shaman) registerChainLightningSpell() {
	numHits := min(core.TernaryInt32(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning), 5, 3), shaman.Env.GetNumTargets())
	shaman.ChainLightning = shaman.newChainLightningSpell(false)
	shaman.ChainLightningOverloads = []*core.Spell{}
	for i := int32(0); i < numHits; i++ {
		shaman.ChainLightningOverloads = append(shaman.ChainLightningOverloads, shaman.newChainLightningSpell(true))
	}
}

func (shaman *Shaman) newChainLightningSpell(isElementalOverload bool) *core.Spell {
	castTime := time.Millisecond * 2500
	spellCoeff := 0.571
	canOverload := false
	overloadChance := shaman.GetOverloadChance()
	if shaman.Spec == proto.Spec_SpecElementalShaman {
		castTime -= 500
		// 0.36 is shamanism bonus
		spellCoeff += 0.36
		canOverload = true
	}

	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: 421},
		0.26,
		castTime,
		isElementalOverload)

	baseDamage := 1093.0
	numHits := int32(3)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning) {
		baseDamage *= 0.90
		numHits += 2
	}
	numHits = min(numHits, shaman.Env.GetNumTargets())

	dmgReductionPerBounce := 0.7

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		bounceCoeff := 1.0
		curTarget := target
		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
			totalDamage := baseDamage + spellCoeff*spell.SpellPower()
			totalDamage *= bounceCoeff
			result := spell.CalcDamage(sim, curTarget, totalDamage, spell.OutcomeMagicHitAndCrit)

			if canOverload && result.Landed() && sim.RandomFloat("Chain Lightning Elemental Overload") <= overloadChance {
				shaman.ChainLightningOverloads[hitIndex].Cast(sim, curTarget)
			}

			spell.DealDamage(sim, result)

			bounceCoeff *= dmgReductionPerBounce
			curTarget = sim.Environment.NextTargetUnit(curTarget)
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
