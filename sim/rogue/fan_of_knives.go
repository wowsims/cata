package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (rogue *Rogue) registerFanOfKnives() {
	fokSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51723},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagColdBlooded,

		DamageMultiplier: 0.8 * (1 +
			core.TernaryFloat64(rogue.Spec == proto.Spec_SpecCombatRogue, 0.75, 0.0)),
		CritMultiplier:   rogue.MeleeCritMultiplier(false), // TODO (TheBackstabi, 3/16/2024) - Verify what crit table FoK is on
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
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.HasThrown()
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			// Calc and apply all OH hits first, because MH hits can benefit from an OH felstriker proc.
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := fokSpell.Unit.RangedWeaponDamage(sim, fokSpell.RangedAttackPower(aoeTarget)) //ohSpell.Unit.OHWeaponDamage(sim, ohSpell.MeleeAttackPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				// TODO (TheBackstabi 3/16/2024) - Proc Thrown poison + Vile Poisons proc MH/OH poison

				results[i] = fokSpell.CalcDamage(sim, aoeTarget, baseDamage, fokSpell.OutcomeRangedHitAndCrit)
			}
			for i := range sim.Encounter.TargetUnits {
				fokSpell.DealDamage(sim, results[i])
			}
		},
	})
}
