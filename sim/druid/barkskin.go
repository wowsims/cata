package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerBarkskinCD() {
	actionId := core.ActionID{SpellID: 22812}
	hasGlyph := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfBarkskin)

	druid.BarkskinAura = druid.RegisterAura(core.Aura{
		Label:    "Barkskin",
		ActionID: actionId,
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageTakenMultiplier *= 0.8
			if hasGlyph {
				druid.PseudoStats.ReducedCritTakenChance += 0.25
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageTakenMultiplier /= 0.8
			if hasGlyph {
				druid.PseudoStats.ReducedCritTakenChance -= 0.25
			}
		},
	})

	druid.Barkskin = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionId,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: core.TernaryDuration(druid.Spec == proto.Spec_SpecGuardianDruid, time.Second * 30, time.Second * 60),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.BarkskinAura.Activate(sim)
			druid.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Barkskin.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
