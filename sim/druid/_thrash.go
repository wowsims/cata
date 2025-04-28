package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerThrashBearSpell() {
	// Raw parameters from spell database
	hitCoefficient := 1.05599999428
	hitVariance := 0.20999999344
	bleedCoefficient := 0.58899998665

	// Scaled parameters for spell code
	avgBaseDamage := hitCoefficient * druid.ClassSpellScaling // second factor is defined in druid.go
	damageSpread := hitVariance * avgBaseDamage
	flatBaseDamage := avgBaseDamage - damageSpread/2
	flatBleedDamage := bleedCoefficient * druid.ClassSpellScaling

	druid.Thrash = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77758},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label:    "Thrash",
				Duration: time.Second * 6,
			}),
			NumberOfTicks: 3,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = flatBleedDamage + 0.0167*dot.Spell.MeleeAttackPower()
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatBaseDamage + 0.0982*spell.MeleeAttackPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				perTargetDamage := (baseDamage + (sim.RandomFloat("Thrash") * damageSpread)) * sim.Encounter.AOECapMultiplier()
				if druid.BleedCategories.Get(aoeTarget).AnyActive() {
					perTargetDamage *= 1.3
				}
				result := spell.CalcAndDealDamage(sim, aoeTarget, perTargetDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				if result.Landed() {
					spell.Dot(aoeTarget).Apply(sim)
				}
			}
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := avgBaseDamage + 0.0982*spell.MeleeAttackPower()
			initial := spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)

			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := (critChance * (spell.CritMultiplier - 1))
			initial.Damage *= 1 + critMod
			return initial
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			tickBase := (flatBleedDamage + 0.0167*spell.MeleeAttackPower())
			ticks := spell.CalcPeriodicDamage(sim, target, tickBase, spell.OutcomeExpectedMagicAlwaysHit)

			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := (critChance * (spell.CritMultiplier - 1))
			ticks.Damage *= 1 + critMod

			return ticks
		},
	})
}
