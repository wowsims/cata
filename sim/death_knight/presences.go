package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const presenceEffectCategory = "Presence"

func (dk *DeathKnight) registerBloodPresenceAura(timer *core.Timer) {
	threatMult := 5.0
	armorScaling := 1.55
	damageTakenMult := 0.92
	stamDep := dk.NewDynamicMultiplyStat(stats.Stamina, 1.08)
	runicMulti := 1.0 + 0.02*float64(dk.Talents.ImprovedFrostPresence)
	critReduction := 0.03 * float64(dk.Talents.ImprovedBloodPresence)

	actionID := core.ActionID{SpellID: 48263}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)
	healthMetrics := dk.NewHealthMetrics(actionID)
	runeRegenSpeed := 1.0 + 0.1*float64(dk.Talents.ImprovedBloodPresence)

	presenceAura := dk.GetOrRegisterAura(core.Aura{
		Label:      "Blood Presence",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(dk.Spec == proto.Spec_SpecBloodDeathKnight, core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ReducedCritTakenChance += critReduction
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMult
			aura.Unit.EnableDynamicStatDep(sim, stamDep)

			if sim.CurrentTime < 0 {
				// Fully heal the DK if the presence was activated before combat
				aura.Unit.GainHealth(sim, aura.Unit.MaxHealth()-aura.Unit.CurrentHealth(), healthMetrics)
			}

			dk.ApplyDynamicEquipScaling(sim, stats.Armor, armorScaling)
			dk.MultiplyRunicRegen(runicMulti)
			dk.MultiplyRuneRegenSpeed(sim, runeRegenSpeed)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ReducedCritTakenChance -= critReduction
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMult
			aura.Unit.DisableDynamicStatDep(sim, stamDep)
			dk.RemoveDynamicEquipScaling(sim, stats.Armor, armorScaling)
			dk.MultiplyRunicRegen(1 / runicMulti)
			dk.MultiplyRuneRegenSpeed(sim, 1/runeRegenSpeed)
		},
	})
	presenceAura.NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})

	dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: dk.activatePresence(presenceAura, rpMetrics),
	})
}

func (dk *DeathKnight) registerFrostPresenceAura(timer *core.Timer) {
	actionID := core.ActionID{SpellID: 48266}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	damageMulti := 1.1
	runicMulti := 1.1

	if dk.Talents.ImprovedFrostPresence > 0 {
		damageMulti += []float64{0, 0.02, 0.05}[dk.Talents.ImprovedFrostPresence]
	}

	presenceAura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Frost Presence",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.DamageDealtMultiplier *= damageMulti
			dk.MultiplyRunicRegen(runicMulti)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.DamageDealtMultiplier /= damageMulti
			dk.MultiplyRunicRegen(1 / runicMulti)
		},
	})
	presenceAura.NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})

	dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: dk.activatePresence(presenceAura, rpMetrics),
	})
}

func (dk *DeathKnight) registerUnholyPresenceAura(timer *core.Timer) {
	actionID := core.ActionID{SpellID: 48265}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	hasteMulti := 1.1
	if dk.Talents.ImprovedUnholyPresence > 0 {
		hasteMulti += []float64{0, 0.02, 0.05}[dk.Talents.ImprovedUnholyPresence]
	}
	runicMulti := 1.0 + 0.02*float64(dk.Talents.ImprovedFrostPresence)

	unholyPresenceMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		ClassMask: DeathKnightSpellsAll,
		TimeValue: time.Millisecond * -500,
	})

	presenceAura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Unholy Presence",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyMeleeSpeed(sim, hasteMulti)
			dk.MultiplyRuneRegenSpeed(sim, hasteMulti)
			dk.MultiplyRunicRegen(runicMulti)
			unholyPresenceMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyMeleeSpeed(sim, 1/hasteMulti)
			dk.MultiplyRuneRegenSpeed(sim, 1/hasteMulti)
			dk.MultiplyRunicRegen(1 / runicMulti)
			unholyPresenceMod.Deactivate()
		},
	})
	presenceAura.NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})
	presenceAura.NewMovementSpeedEffect(0.15)

	dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: dk.activatePresence(presenceAura, rpMetrics),
	})
}

func (dk *DeathKnight) activatePresence(presence *core.Aura, rpMetrics *core.ResourceMetrics) core.ApplySpellResults {
	return func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
		presence.Activate(sim)
		if rp := dk.CurrentRunicPower(); rp > 0 && sim.CurrentTime >= 0 {
			dk.SpendRunicPower(sim, rp, rpMetrics)
		}
	}
}

func (dk *DeathKnight) registerPresences() {
	presenceTimer := dk.NewTimer()
	dk.registerBloodPresenceAura(presenceTimer)
	dk.registerUnholyPresenceAura(presenceTimer)
	dk.registerFrostPresenceAura(presenceTimer)
}
