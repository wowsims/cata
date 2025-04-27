package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerEnvenom() {
	coefficient := 0.21400000155
	apScalingPerComboPoint := 0.09

	baseDamage := coefficient * rogue.ClassSpellScaling

	rogue.EnvenomAura = rogue.RegisterAura(core.Aura{
		Label:    "Envenom",
		ActionID: core.ActionID{SpellID: 32645},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.UpdateDeadlyPoisonPPH(0.15)
			rogue.UpdateInstantPoisonPPM(0.75)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.UpdateDeadlyPoisonPPH(0.0)
			rogue.UpdateInstantPoisonPPM(0.0)
		},
	})

	chanceToRetainStacks := core.TernaryFloat64(rogue.Talents.MasterPoisoner, 1, 0)

	rogue.Envenom = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 32645},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskMeleeMHSpecial, // not core.ProcMaskSpellDamage
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | SpellFlagColdBlooded | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellEnvenom,

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.8,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0 && rogue.DeadlyPoison.Dot(target).IsActive()
		},

		DamageMultiplier: 1,
		DamageMultiplierAdditive: 1 +
			0.12*float64(rogue.Talents.VilePoisons) +
			[]float64{0, 0.07, .14, .20}[rogue.Talents.CoupDeGrace],
		CritMultiplier:   rogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			comboPoints := rogue.ComboPoints()
			// - the aura is active even if the attack fails to land
			// - the aura is applied before the hit effect
			// See: https://github.com/where-fore/rogue-wotlk/issues/32
			rogue.EnvenomAura.Duration = time.Second * time.Duration(1+comboPoints)
			rogue.EnvenomAura.Activate(sim)

			dp := rogue.DeadlyPoison.Dot(target)
			// - 215 base is scaled by consumed doses (<= comboPoints)
			// - apRatio is independent of consumed doses (== comboPoints)
			consumed := min(dp.GetStacks(), comboPoints)
			baseDamage := baseDamage*float64(consumed) +
				apScalingPerComboPoint*float64(comboPoints)*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
				rogue.ApplyCutToTheChase(sim)
				if !sim.Proc(chanceToRetainStacks, "Master Poisoner") {
					if newStacks := dp.GetStacks() - comboPoints; newStacks > 0 {
						dp.SetStacks(sim, newStacks)
					} else {
						dp.Deactivate(sim)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (rogue *Rogue) EnvenomDuration(comboPoints int32) time.Duration {
	return time.Second * (1 + time.Duration(comboPoints))
}
