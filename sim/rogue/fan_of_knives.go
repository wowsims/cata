package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerFanOfKnives() {
	baseDamage := rogue.GetBaseDamageFromCoefficient(1.25)
	apScaling := 0.17499999702
	damageSpread := baseDamage * 0.40000000596
	minDamage := baseDamage - damageSpread/2

	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 51723})

	fokSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51723},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagColdBlooded,
		ClassSpellMask: RogueSpellFanOfKnives,

		DamageMultiplier: 1,
		CritMultiplier:   rogue.CritMultiplier(false),
		ThreatMultiplier: 1,
	})

	results := make([]*core.SpellResult, len(rogue.Env.Encounter.TargetUnits))

	rogue.FanOfKnives = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51723},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost: 35,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				damage := minDamage +
					sim.RandomFloat("Fan of Knives")*damageSpread +
					spell.MeleeAttackPower()*apScaling

				damage *= sim.Encounter.AOECapMultiplier()

				results[i] = fokSpell.CalcAndDealDamage(sim, aoeTarget, damage, fokSpell.OutcomeMeleeSpecialNoBlockDodgeParry)
				if results[i].Landed() && aoeTarget == rogue.CurrentTarget {
					rogue.AddComboPointsOrAnticipation(sim, 1, cpMetrics)
				}
			}
		},
	})
}
