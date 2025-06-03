package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (druid *Druid) registerFrenziedRegenerationSpell() {
	actionID := core.ActionID{SpellID: 22842}
	isGlyphed := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFrenziedRegeneration)
	buffConfig := core.Aura{
		Label:    "Frenzied Regeneration",
		ActionID: actionID,
		Duration: time.Second * 6,

		OnGain: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.HealingTakenMultiplier *= 1.4
		},

		OnExpire: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.HealingTakenMultiplier /= 1.4
		},
	}

	var rageMetrics *core.ResourceMetrics

	if isGlyphed {
		druid.FrenziedRegenerationAura = druid.RegisterAura(buffConfig)
	} else {
		rageMetrics = druid.NewRageMetrics(actionID)
	}

	druid.FrenziedRegeneration = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:         actionID,
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},

		RageCost: core.RageCostOptions{
			Cost: core.TernaryInt32(isGlyphed, 50, 0),
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if isGlyphed {
				druid.FrenziedRegenerationAura.Activate(sim)
			} else {
				const maxRageCost = 60.0
				rageDumped := min(druid.CurrentRage(), maxRageCost)
				healthGained := max((druid.GetStat(stats.AttackPower)-2*druid.GetStat(stats.Agility))*2.2, druid.GetStat(stats.Stamina)*2.5) * rageDumped / maxRageCost
				spell.CalcAndDealHealing(sim, spell.Unit, healthGained, spell.OutcomeHealing)
				druid.SpendRage(sim, rageDumped, rageMetrics)
			}
		},
	})
}
