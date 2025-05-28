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
				dot.SnapshotAttackerMultiplier = 1
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)

			svHunter.GetHunter().HuntersMarkSpell.Cast(sim, target)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					dot := spell.Dot(target)
					rap := dot.Spell.RangedAttackPower()
					dot.SnapshotBaseDamage = ((dot.OutstandingDmg() / 3) + ((243.35 + sim.RandomFloat("Explosive Shot")*487) + (0.39 * rap)))
					dot.Apply(sim)
					dot.TickOnce(sim)
					spell.DealOutcome(sim, result)

				}
			})
		},
	})

}
