package elemental

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (ele *ElementalShaman) registerLavaBeamSpell() {
	numHits := min(core.TernaryInt32(ele.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning), 5, 3), ele.Env.GetNumTargets())
	ele.LavaBeam = ele.newLavaBeamSpell(false)
	ele.LavaBeamOverloads = []*core.Spell{}
	for i := int32(0); i < numHits; i++ {
		ele.LavaBeamOverloads = append(ele.LavaBeamOverloads, ele.newLavaBeamSpell(true))
	}
}

func (ele *ElementalShaman) newLavaBeamSpell(isElementalOverload bool) *core.Spell {
	spellConfig := ele.NewChainSpellConfig(114074, isElementalOverload, 1.1, 0.57099997997, 1.08800005913, 0.13300000131, core.SpellSchoolFire, 8.3, ele.LavaBeamOverloads)
	spellConfig.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return ele.GetAura("Ascendance").IsActive()
	}
	return ele.RegisterSpell(spellConfig)
}
