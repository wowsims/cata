package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) ApplyGlyphs() {
	if druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfHealingTouch) {
		druid.RegisterAura(core.Aura{
			Label:    "Glyph of Healing Touch",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.Matches(DruidSpellHealingTouch) && !druid.NaturesSwiftness.CD.IsReady(sim) {
					*druid.NaturesSwiftness.CD.Timer = core.Timer(time.Duration(*druid.NaturesSwiftness.CD.Timer) - time.Second*3)
					druid.UpdateMajorCooldowns()
				}
			},
		})
	}
}
