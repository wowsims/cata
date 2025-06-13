package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerMightOfUrsocCD() {
	actionID := core.ActionID{SpellID: 106922}
	healthMetrics := druid.NewHealthMetrics(actionID)
	isGlyphed := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMightOfUrsoc)
	bonusHealthFrac := core.TernaryFloat64(isGlyphed, 0.5, 0.3)

	var bonusHealth float64

	druid.MightOfUrsocAura = druid.RegisterAura(core.Aura{
		Label:    "Might of Ursoc",
		ActionID: actionID,
		Duration: time.Second * 20,

		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			bonusHealth = druid.MaxHealth() * bonusHealthFrac
			druid.UpdateMaxHealth(sim, bonusHealth, healthMetrics)
		},

		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			druid.UpdateMaxHealth(sim, -bonusHealth, healthMetrics)
		},
	})

	druid.MightOfUrsoc = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: core.TernaryDuration(isGlyphed, time.Minute * 5, time.Minute * 3),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if !druid.InForm(Bear) {
				druid.BearFormAura.Activate(sim)
			}

			druid.MightOfUrsocAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.MightOfUrsoc.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
