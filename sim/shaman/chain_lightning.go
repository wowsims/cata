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

func (shaman *Shaman) NewChainSpellConfig(spellID int32, isElementalOverload bool, bounceReduction float64, bonusCoeff float64, coeff float64, variance float64, spellSchool core.SpellSchool, manaCost float64, overloads []*core.Spell) core.SpellConfig {
	spellConfig := shaman.newElectricSpellConfig(core.ActionID{SpellID: spellID}, manaCost, time.Second*2, isElementalOverload, bonusCoeff)
	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, SpellMaskChainLightningOverload, SpellMaskChainLightning)
	if !isElementalOverload {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 3,
		}
	}
	spellConfig.SpellSchool = spellSchool

	numHits := int32(3)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning) {
		numHits += 2
	}
	numHits = min(numHits, shaman.Env.GetNumTargets())

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, coeff, variance)
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
					overloads[hitIndex].Cast(sim, results[hitIndex].Target)
				}
			}

			spell.DealDamage(sim, results[hitIndex])
			spell.DamageMultiplier /= bounceReduction
		}
	}
	return spellConfig
}

func (shaman *Shaman) newChainLightningSpell(isElementalOverload bool) *core.Spell {
	spellConfig := shaman.NewChainSpellConfig(421, isElementalOverload, 1.0, 0.51800000668, 0.98900002241, 0.13300000131, core.SpellSchoolNature, 30.5, shaman.ChainLightningOverloads)

	return shaman.RegisterSpell(spellConfig)
}
