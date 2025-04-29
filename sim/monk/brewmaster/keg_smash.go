package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerKegSmash() {
	actionID := core.ActionID{SpellID: 121253}
	chiMetrics := bm.NewChiMetrics(actionID)

	bm.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | monk.SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellKegSmash,
		MaxRange:       core.MaxMeleeRange,
		MissileSpeed:   30,

		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    bm.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplier: 10.0,
		ThreatMultiplier: 1,
		CritMultiplier:   bm.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bm.StanceMatches(monk.SturdyOx)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bm.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
				return
			}

			bm.AddChi(sim, spell, 2, chiMetrics)
			spell.DealOutcome(sim, result)

			bm.DizzyingHazeAuras.Get(target).Activate(sim)
			for _, otherTarget := range sim.Encounter.TargetUnits {
				if otherTarget != target {
					result := spell.CalcAndDealDamage(sim, otherTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
					if result.Landed() {
						bm.DizzyingHazeAuras.Get(otherTarget).Activate(sim)
					}
				}
			}
		},
	})
}
