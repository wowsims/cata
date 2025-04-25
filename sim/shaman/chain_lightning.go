package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
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
	spellConfig := shaman.newElectricSpellConfig(core.ActionID{SpellID: 421}, 26, time.Second*2, isElementalOverload, 0.571)
	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, SpellMaskChainLightningOverload, SpellMaskChainLightning)

	if !isElementalOverload {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 3,
		}
	}

	numHits := int32(3)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning) {
		spellConfig.DamageMultiplier *= 0.90
		numHits += 2
	}
	numHits = min(numHits, shaman.Env.GetNumTargets())

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		bounceReduction := core.TernaryFloat64(shaman.DungeonSet3.IsActive() && !isElementalOverload, 0.83, 0.7)
		baseDamage := shaman.CalcAndRollDamageRange(sim, 1.08800005913, 0.13300000131)
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
			if !spell.ProcMask.Matches(core.ProcMaskSpellProc) { //So that procs from DTR does not cast an overload
				if !isElementalOverload && results[hitIndex].Landed() && sim.Proc(shaman.GetOverloadChance()/3, "Chain Lightning Elemental Overload") {
					shaman.ChainLightningOverloads[hitIndex].Cast(sim, results[hitIndex].Target)
				}
			}

			spell.DealDamage(sim, results[hitIndex])
			spell.DamageMultiplier /= bounceReduction
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
