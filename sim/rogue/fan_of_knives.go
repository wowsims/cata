package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerFanOfKnives() {
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
				baseDamage := fokSpell.Unit.RangedWeaponDamage(sim, fokSpell.RangedAttackPower(aoeTarget))
				baseDamage *= sim.Encounter.AOECapMultiplier()

				results[i] = fokSpell.CalcAndDealDamage(sim, aoeTarget, baseDamage, fokSpell.OutcomeRangedHitAndCrit)
			}
		},
	})
}
