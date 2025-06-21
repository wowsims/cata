package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerProwlSpell() {
	actionID := core.ActionID{SpellID: 5215}
	movementSpeedMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfProwl), 1.0, 0.7)

	icd := core.Cooldown{
		Timer:    druid.NewTimer(),
		Duration: time.Second * 10,
	}

	druid.ProwlAura = druid.RegisterAura(core.Aura{
		Label:    "Prowl",
		ActionID: actionID,
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMovementSpeed(sim, movementSpeedMultiplier)
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			aura.Deactivate(sim)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			icd.Use(sim)
			aura.Unit.MultiplyMovementSpeed(sim, 1.0/movementSpeedMultiplier)
		},
	})

	druid.CatFormAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		if druid.ProwlAura.IsActive() {
			druid.ProwlAura.Deactivate(sim)
		}
	})

	druid.Prowl = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: icd,
		},

		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			return (sim.CurrentTime < 0)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if !druid.InForm(Cat) {
				druid.CatFormAura.Activate(sim)
			}

			druid.ProwlAura.Activate(sim)
		},
	})
}
