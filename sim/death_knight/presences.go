package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

const presenceEffectCategory = "Presence"

func (dk *DeathKnight) registerBloodPresenceAura(timer *core.Timer) {
	threatMult := 4.0
	armorScaling := 1.55
	damageTakenMult := 1 / 1.08
	stamDep := dk.NewDynamicMultiplyStat(stats.Stamina, 1.08)

	actionID := core.ActionID{SpellID: 48263}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	dk.BloodPresence = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.BloodPresenceAura.Activate(sim)
			if dk.CurrentRunicPower() > 0 {
				dk.SpendRunicPower(sim, dk.CurrentRunicPower(), rpMetrics)
			}
		},
	})

	aura := core.Aura{
		Label:    "Blood Presence",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageTakenMult
			aura.Unit.EnableDynamicStatDep(sim, stamDep)
			dk.ApplyDynamicEquipScaling(sim, stats.Armor, armorScaling)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageTakenMult
			aura.Unit.DisableDynamicStatDep(sim, stamDep)
			dk.RemoveDynamicEquipScaling(sim, stats.Armor, armorScaling)
		},
	}

	dk.BloodPresenceAura = dk.GetOrRegisterAura(aura)
	dk.BloodPresenceAura.NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})
}

func (dk *DeathKnight) registerFrostPresenceAura(timer *core.Timer) {
	actionID := core.ActionID{SpellID: 48266}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	damageMulti := 1.1

	dk.FrostPresence = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost: 1,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.FrostPresenceAura.Activate(sim)
			if dk.CurrentRunicPower() > 0 {
				dk.SpendRunicPower(sim, dk.CurrentRunicPower(), rpMetrics)
			}
		},
	})

	// TODO: Runic Power Gen

	dk.FrostPresenceAura = dk.GetOrRegisterAura(core.Aura{
		Label:    "Frost Presence",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.DamageDealtMultiplier *= damageMulti
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.DamageDealtMultiplier /= damageMulti
		},
	})
	dk.FrostPresenceAura.NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})
}

func (dk *DeathKnight) registerUnholyPresenceAura(timer *core.Timer) {
	actionID := core.ActionID{SpellID: 48265}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	hasteMulti := 1.1
	if dk.Talents.ImprovedUnholyPresence > 0 {
		hasteMulti += []float64{0, 0.02, 0.05}[dk.Talents.ImprovedUnholyPresence]
	}

	dk.UnholyPresence = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.UnholyPresenceAura.Activate(sim)
			if dk.CurrentRunicPower() > 0 {
				dk.SpendRunicPower(sim, dk.CurrentRunicPower(), rpMetrics)
			}
		},
	})

	unholyPresenceMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		ClassMask: DeathKnightSpellsAll,
		TimeValue: time.Millisecond * -500,
	})

	dk.UnholyPresenceAura = dk.GetOrRegisterAura(core.Aura{
		Label:    "Unholy Presence",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyMeleeSpeed(sim, hasteMulti)
			dk.MultiplyRuneRegenSpeed(sim, hasteMulti)
			unholyPresenceMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyMeleeSpeed(sim, 1/hasteMulti)
			dk.MultiplyRuneRegenSpeed(sim, 1/hasteMulti)
			unholyPresenceMod.Deactivate()
		},
	})
	dk.UnholyPresenceAura.NewExclusiveEffect(presenceEffectCategory, true, core.ExclusiveEffect{})
}

func (dk *DeathKnight) registerPresences() {
	presenceTimer := dk.NewTimer()
	dk.registerBloodPresenceAura(presenceTimer)
	dk.registerUnholyPresenceAura(presenceTimer)
	dk.registerFrostPresenceAura(presenceTimer)
}
