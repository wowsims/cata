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
	spellConfig := shaman.newElectricSpellConfig(core.ActionID{SpellID: 421}, 0.26, time.Second*2, isElementalOverload, 0.571)
	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, SpellMaskChainLightningOverload, SpellMaskChainLightning)

	if !isElementalOverload {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 3,
		}
	}

	baseDamage := shaman.ClassSpellScaling * 1.08800005913
	numHits := int32(3)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning) {
		spellConfig.DamageMultiplier *= 0.90
		numHits += 2
	}
	numHits = min(numHits, shaman.Env.GetNumTargets())

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		bounceReduction := 0.7
		curTarget := target

		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
			result := shaman.calcDamageStormstrikeCritChance(sim, curTarget, baseDamage, spell)
			spell.DealDamage(sim, result)

			if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance()/3, "Chain Lightning Elemental Overload") {
				shaman.ChainLightningOverloads[hitIndex].Cast(sim, result.Target)
			}

			curTarget = sim.Environment.NextTargetUnit(curTarget)
			spell.DamageMultiplier *= bounceReduction
		}

		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
			spell.DamageMultiplier /= bounceReduction
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
