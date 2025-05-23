package shadow

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const divineStarScale = 4.495
const divineStarCoeff = 0.455
const divineStarVariance = 0.5

func (shadow *ShadowPriest) registerDivineStar() {
	if !shadow.Talents.DivineStar {
		return
	}

	shadow.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 122128},
		SpellSchool:      core.SpellSchoolShadow,
		Flags:            core.SpellFlagAPL,
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1,
		CritMultiplier:   shadow.DefaultCritMultiplier(),
		BonusCoefficient: divineStarCoeff,
		ThreatMultiplier: 1,
		MaxRange:         30,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4.5,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shadow.NewTimer(),
				Duration: time.Second * 15,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hit1 := shadow.DistanceFromTarget / 24
			hit2 := 2.5 - hit1

			// first hit
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + time.Second*time.Duration(hit1),
				OnAction: func(s *core.Simulation) {
					for _, unit := range sim.Encounter.TargetUnits {
						spell.CalcAndDealDamage(
							sim,
							unit,
							shadow.CalcAndRollDamageRange(sim, divineStarScale, divineStarVariance),
							spell.OutcomeMagicHitAndCrit,
						)
					}
				}})

			// second hit
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + time.Second*time.Duration(hit2),
				OnAction: func(s *core.Simulation) {
					for _, unit := range sim.Encounter.TargetUnits {
						spell.CalcAndDealDamage(
							sim,
							unit,
							shadow.CalcAndRollDamageRange(sim, divineStarScale, divineStarVariance),
							spell.OutcomeMagicHitAndCrit,
						)
					}
				}})
		},
	})
}
