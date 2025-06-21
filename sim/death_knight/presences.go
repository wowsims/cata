package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const presenceEffectCategory = "Presence"

func (dk *DeathKnight) registerBloodPresence() {
	actionID := core.ActionID{SpellID: 48263}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)
	healthMetrics := dk.NewHealthMetrics(actionID)
	buildPhase := core.Ternary(dk.Spec == proto.Spec_SpecBloodDeathKnight, core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone)

	presenceAura := dk.RegisterAura(core.Aura{
		Label:      "Blood Presence" + dk.Label,
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: buildPhase,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if sim.CurrentTime < 0 {
				// Fully heal the DK if the presence was activated before combat
				aura.Unit.GainHealth(sim, aura.Unit.MaxHealth()-aura.Unit.CurrentHealth(), healthMetrics)
			}

			dk.ApplyDynamicEquipScaling(sim, stats.Armor, 1.55)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.RemoveDynamicEquipScaling(sim, stats.Armor, 1.55)
		},
	}).AttachMultiplicativePseudoStatBuff(
		&dk.PseudoStats.DamageTakenMultiplier, 0.9,
	).AttachMultiplicativePseudoStatBuff(
		&dk.PseudoStats.ThreatMultiplier, 7.0,
	).AttachStatDependency(
		dk.NewDynamicMultiplyStat(stats.Stamina, 1.25),
	).NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})

	dk.BloodPresenceSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		ApplyEffects: dk.activatePresence(presenceAura.Aura, rpMetrics),

		RelatedSelfBuff: presenceAura.Aura,
	})
}

func (dk *DeathKnight) registerFrostPresence() {
	actionID := core.ActionID{SpellID: 48266}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)
	buildPhase := core.Ternary(dk.Spec == proto.Spec_SpecFrostDeathKnight, core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone)

	presenceAura := dk.GetOrRegisterAura(core.Aura{
		Label:      "Frost Presence" + dk.Label,
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: buildPhase,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyRunicRegen(1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyRunicRegen(1 / 1.2)
		},
	}).NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})

	dk.FrostPresenceSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		ApplyEffects: dk.activatePresence(presenceAura.Aura, rpMetrics),

		RelatedSelfBuff: presenceAura.Aura,
	})
}

func (dk *DeathKnight) registerUnholyPresence() {
	actionID := core.ActionID{SpellID: 48265}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)
	buildPhase := core.Ternary(dk.Spec == proto.Spec_SpecUnholyDeathKnight, core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone)

	hasteMulti := 1.1
	if dk.Spec == proto.Spec_SpecUnholyDeathKnight {
		hasteMulti += 0.1
	}

	presenceAura := dk.GetOrRegisterAura(core.Aura{
		Label:      "Unholy Presence" + dk.Label,
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: buildPhase,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyAttackSpeed(sim, hasteMulti)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyAttackSpeed(sim, 1/hasteMulti)
		},
	}).NewMovementSpeedEffect(0.15)
	presenceAura.Aura.NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})

	dk.UnholyPresenceSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		ApplyEffects: dk.activatePresence(presenceAura.Aura, rpMetrics),

		RelatedSelfBuff: presenceAura.Aura,
	})
}

func (dk *DeathKnight) activatePresence(presence *core.Aura, rpMetrics *core.ResourceMetrics) core.ApplySpellResults {
	hasGlyphOfShiftingPresences := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfShiftingPresences)

	return func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
		presence.Activate(sim)
		if rp := dk.CurrentRunicPower(); rp > 0 && sim.CurrentTime >= 0 {
			if hasGlyphOfShiftingPresences {
				rp *= 0.3
			}

			dk.SpendRunicPower(sim, rp, rpMetrics)
		}
	}
}

func (dk *DeathKnight) registerPresences() {
	dk.registerBloodPresence()
	dk.registerUnholyPresence()
	dk.registerFrostPresence()
}
