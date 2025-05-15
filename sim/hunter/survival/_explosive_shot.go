package survival

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (svHunter *SurvivalHunter) registerExplosiveShotSpell() {
	actionID := core.ActionID{SpellID: 53301}
	altActionID := core.ActionID{SpellID: 1215485}
	minFlatDamage := 410.708 - (76.8024 / 2)
	maxFlatDamage := 410.708 + (76.8024 / 2)
	var explosiveShotCounter = 0
	var alternateExplosiveShot = svHunter.Hunter.RegisterSpell(core.SpellConfig{
		ActionID:         altActionID,
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskRangedSpecial,
		DamageMultiplier: 1,
		ClassSpellMask:   hunter.HunterSpellExplosiveShot,
		CritMultiplier:   svHunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		Flags:            core.SpellFlagMeleeMetrics,
		Dot: core.DotConfig{
			Aura: core.Aura{
				ActionID: altActionID,
				Label:    "Explosive Shot - Dot (Second)",
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					explosiveShotCounter = 0
				},
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				rap := dot.Spell.RangedAttackPower()
				baseDmg := sim.Roll(minFlatDamage, maxFlatDamage) + (0.273 * rap)
				dot.Snapshot(target, baseDmg)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
			// OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			// 	rap := dot.Spell.RangedAttackPower()
			// 	baseDmg := sim.Roll(minFlatDamage, maxFlatDamage) + (0.273 * rap)
			// 	dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDmg, dot.OutcomeTickPhysicalCrit)
			// },
		},
	})
	svHunter.Hunter.ExplosiveShot = svHunter.Hunter.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ClassSpellMask: hunter.HunterSpellExplosiveShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       5,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 50,
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
		DamageMultiplier: 1,
		CritMultiplier:   svHunter.CritMultiplier(1, 1),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Explosive Shot - Dot (First)",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				rap := dot.Spell.RangedAttackPower()
				baseDmg := sim.Roll(minFlatDamage, maxFlatDamage) + (0.273 * rap)
				dot.Snapshot(target, baseDmg)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			explosiveShotCounter++

			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					var dot *core.Dot
					if explosiveShotCounter%2 == 0 {
						dot = alternateExplosiveShot.Dot(target)
					} else {
						dot = spell.Dot(target)
					}
					dot.Apply(sim)
					dot.TickOnce(sim)
					spell.DealOutcome(sim, result)

				}
			})
		},
	})

}
