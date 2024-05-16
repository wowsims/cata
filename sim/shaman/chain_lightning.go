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

		// Damage calculation and DealDamage are in separate loops so that e.g. a spell power proc
		// can't proc on the first target and apply to the second
		results := make([]*core.SpellResult, numHits)
		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
			results[hitIndex] = shaman.calcDamageStormstrikeCritChance(sim, curTarget, baseDamage, spell)

			curTarget = sim.Environment.NextTargetUnit(curTarget)
			spell.DamageMultiplier *= bounceReduction
		}

		for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
			if !isElementalOverload && results[hitIndex].Landed() && sim.Proc(shaman.GetOverloadChance()/3, "Chain Lightning Elemental Overload") {
				shaman.ChainLightningOverloads[hitIndex].Cast(sim, results[hitIndex].Target)
			}

			spell.DealDamage(sim, results[hitIndex])
			spell.DamageMultiplier /= bounceReduction
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
