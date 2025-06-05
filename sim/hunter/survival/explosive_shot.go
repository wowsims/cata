package survival

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (svHunter *SurvivalHunter) registerExplosiveShotSpell() {
	actionID := core.ActionID{SpellID: 53301}

	svHunter.Hunter.ExplosiveShot = svHunter.Hunter.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ClassSpellMask: hunter.HunterSpellExplosiveShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       0,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    svHunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,

		CritMultiplier:   svHunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Explosive Shot",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickMagicCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rap := spell.RangedAttackPower()
			baseDamage := svHunter.CalcAndRollDamageRange(sim, 0.391, 1) + (0.391 * rap)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			svHunter.GetHunter().HuntersMarkSpell.Cast(sim, target)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					dot := spell.Dot(target)
					remaining := dot.RemainingTicks()
					// If only the last tick is still queued, add one more so we end up with two.
					if remaining == 1 {
						dot.AddTick()
					}

					// When more than once tick remain, add it to the outstanding damage pile
					outstandingDamage := 0.0
					if remaining > 1 {
						outstandingDamage = dot.OutstandingDmg() / 4
					}

					// Tick before outstanding damage is added
					dot.Snapshot(target, baseDamage+outstandingDamage)

					dot.Apply(sim)
					spell.DealDamage(sim, result)
				}
			})
		},
	})

}
