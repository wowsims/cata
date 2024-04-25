package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO: Check spell damage and coefficients
func (warlock *Warlock) registerSeedSpell() {
	actionID := core.ActionID{SpellID: 27243}

	seedExplosion := warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1), // actually 47834
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagNoLogs,
		ClassSpellMask: WarlockSpellSeedOfCorruptionExposion,

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.17159999907,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := (862 + 0.17159999907*spell.SpellPower()) * sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	warlock.SeedDamageTracker = make([]float64, len(warlock.Env.AllUnits))
	trySeedPop := func(sim *core.Simulation, target *core.Unit, dmg float64) {
		warlock.SeedDamageTracker[target.UnitIndex] += dmg
		if warlock.SeedDamageTracker[target.UnitIndex] > 2378 {
			warlock.Seed.Dot(target).Deactivate(sim)
			seedExplosion.Cast(sim, target)
		}
	}

	warlock.Seed = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagAPL,
		MissileSpeed:   28,
		ClassSpellMask: WarlockSpellSeedOfCorruption,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.34,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Seed",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}
					if spell.ActionID.SpellID == actionID.SpellID {
						return // Seed can't pop seed.
					}
					trySeedPop(sim, aura.Unit, result.Damage)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					trySeedPop(sim, aura.Unit, result.Damage)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.SeedDamageTracker[aura.Unit.UnitIndex] = 0
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.SeedDamageTracker[aura.Unit.UnitIndex] = 0
				},
			},

			NumberOfTicks:    6,
			TickLength:       time.Second * 3,
			BonusCoefficient: 0.30,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 2042/6)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					// seed is mutually exclusive with corruption
					warlock.Corruption.Dot(target).Deactivate(sim)

					if warlock.Options.DetonateSeed {
						seedExplosion.Cast(sim, target)
					} else {
						spell.Dot(target).Apply(sim)
					}
				}
			})
		},
	})
}
