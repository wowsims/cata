package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (asnRogue *AssassinationRogue) registerEnvenom() {
	baseDamage := asnRogue.GetBaseDamageFromCoefficient(0.38499999046)
	apScalingPerComboPoint := 0.112

	// Envenom has a DoT-like clipping window, where it adds up to 1 seconds to the new duration.
	// This functions exactly like DoT clipping, just for a standard aura
	clipInterval := time.Second * 1

	asnRogue.EnvenomAura = asnRogue.RegisterAura(core.Aura{
		Label:    "Envenom",
		ActionID: core.ActionID{SpellID: 32645},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			asnRogue.UpdateLethalPoisonPPH(0.15)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			asnRogue.UpdateLethalPoisonPPH(0.0)
		},
	})

	asnRogue.Envenom = asnRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 32645},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskMeleeMHSpecial, // not core.ProcMaskSpellDamage
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: rogue.RogueSpellEnvenom,

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.8,
			RefundMetrics: asnRogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 700,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(asnRogue.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return asnRogue.ComboPoints() > 0
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           asnRogue.CritMultiplier(false),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			asnRogue.BreakStealth(sim)
			comboPoints := asnRogue.ComboPoints()

			bonusDuration := time.Duration(0)
			if asnRogue.EnvenomAura.IsActive() {
				bonusDuration = asnRogue.EnvenomAura.RemainingDuration(sim) % clipInterval
				bonusDuration = core.TernaryDuration(bonusDuration == 0, clipInterval, bonusDuration)
			}
			asnRogue.EnvenomAura.Duration = time.Second*time.Duration(1+comboPoints) + bonusDuration
			if asnRogue.Has2PT15 {
				asnRogue.EnvenomAura.Duration += time.Second * 1
			}
			asnRogue.EnvenomAura.Activate(sim)

			baseDamage := baseDamage*float64(comboPoints) +
				apScalingPerComboPoint*float64(comboPoints)*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				asnRogue.ApplyFinisher(sim, spell)
				asnRogue.ApplyCutToTheChase(sim)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (asnRogue *AssassinationRogue) EnvenomDuration(comboPoints int32) time.Duration {
	return time.Second * (1 + time.Duration(comboPoints))
}
