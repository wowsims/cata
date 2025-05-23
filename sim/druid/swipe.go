package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerSwipeBearSpell() {
	flatBaseDamage := 0.94199997187 * druid.ClassSpellScaling // ~929

	druid.SwipeBear = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 779},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		MaxRange:         8,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatBaseDamage + 0.123*spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		},
	})
}

func (druid *Druid) registerSwipeCatSpell() {
	weaponMulti := 5.25

	druid.SwipeCat = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 62078},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost: 45,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: weaponMulti,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())
			baseDamage *= sim.Encounter.AOECapMultiplier()
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWeaponSpecialHitAndCrit)
		},
	})
}

func (druid *Druid) CurrentSwipeCatCost() float64 {
	return druid.SwipeCat.Cost.GetCurrentCost()
}

func (druid *Druid) IsSwipeSpell(spell *core.Spell) bool {
	return druid.SwipeBear.IsEqual(spell) || druid.SwipeCat.IsEqual(spell)
}
