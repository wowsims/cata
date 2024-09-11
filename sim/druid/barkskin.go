package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) registerBarkskinCD() {
	if !druid.InForm(Bear) {
		return
	}

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

			if druid.Feral4pT12Active {
				druid.SmokescreenAura.Activate(sim)
			}
		},
	})

	druid.Barkskin = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionId,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 60.0,
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
