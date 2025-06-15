package elemental

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/shaman"
)

func (ele *ElementalShaman) registerLavaBeamSpell() {
	numHits := min(core.TernaryInt32(ele.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning), 5, 3), ele.Env.GetNumTargets())
	ele.LavaBeam = ele.newLavaBeamSpell(false)
	ele.LavaBeamOverloads = [2][]*core.Spell{}
	for i := int32(0); i < numHits; i++ {
		ele.LavaBeamOverloads[0] = append(ele.LavaBeamOverloads[0], ele.newLavaBeamSpell(true))
		ele.LavaBeamOverloads[1] = append(ele.LavaBeamOverloads[1], ele.newLavaBeamSpell(true))
	}
}

func (ele *ElementalShaman) newLavaBeamSpell(isElementalOverload bool) *core.Spell {
	shamConfig := shaman.ShamSpellConfig{
		ActionID:            core.ActionID{SpellID: 114074},
		IsElementalOverload: isElementalOverload,
		BaseCostPercent:     8.3,
		BonusCoefficient:    0.57099997997,
		Coeff:               1.08800005913,
		Variance:            0.13300000131,
		SpellSchool:         core.SpellSchoolFire,
		Overloads:           &ele.LavaBeamOverloads,
		BounceReduction:     1.1,
	}
	spellConfig := ele.NewChainSpellConfig(shamConfig)
	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, shaman.SpellMaskLavaBeamOverload, shaman.SpellMaskLavaBeam)
	spellConfig.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return ele.AscendanceAura.IsActive()
	}
	return ele.RegisterSpell(spellConfig)
}
