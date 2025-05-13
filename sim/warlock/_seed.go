package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) registerSeed() {
	actionID := core.ActionID{SpellID: 27243}

	seedExplosion := warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1), // actually 27285
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagNoLogs | core.SpellFlagPassiveSpell,
		ClassSpellMask: WarlockSpellSeedOfCorruptionExposion,

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.22920000553,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := warlock.CalcAndRollDamageRange(sim, 0.76560002565, 0.15000000596) * sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	seedDamageTracker := make([]float64, len(warlock.Env.AllUnits))
	trySeedPop := func(sim *core.Simulation, target *core.Unit, dmg float64) {
		seedDamageTracker[target.UnitIndex] += dmg
		// TODO: this is probably calculated by 0.17159999907*SP + warlock.CalcScalingSpellDmg(2.11299991608)
		if seedDamageTracker[target.UnitIndex] > 2378 {
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

		ManaCost: core.ManaCostOptions{BaseCostPercent: 34},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2000 * time.Millisecond,
			},
		},

		CritMultiplier:           warlock.DefaultCritMultiplier(),
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
					seedDamageTracker[aura.Unit.UnitIndex] = 0
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					seedDamageTracker[aura.Unit.UnitIndex] = 0
				},
			},

			NumberOfTicks:    6,
			TickLength:       3 * time.Second,
			BonusCoefficient: 0.30000001192,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, warlock.CalcScalingSpellDmg(0.30239999294))
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
