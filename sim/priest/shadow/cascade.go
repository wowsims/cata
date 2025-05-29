package shadow

import (
	"math"
	"slices"
	"time"

	"github.com/wowsims/mop/sim/core"
)

const cascadeScale = 12
const cascadeCoeff = 1.225

func (shadow *ShadowPriest) registerCascade() {
	if !shadow.Talents.Cascade {
		return
	}

	targets := []*core.Unit{}
	cascadeHandler := func(damageMod float64, bounceSpell *core.Spell, target *core.Unit, sim *core.Simulation) {
		bounceSpell.DamageMultiplier *= damageMod
		bounceSpell.CalcAndDealDamage(sim, target, shadow.CalcScalingSpellDmg(cascadeScale), bounceSpell.OutcomeMagicHitAndCrit)
		bounceSpell.DamageMultiplier /= damageMod

		if len(targets) >= 31 {
			return
		}

		bounceTargets := []*core.Unit{}
		for _, unit := range sim.Encounter.TargetUnits {
			if unit == target {
				continue
			}

			if slices.Contains(targets, unit) {
				continue
			}

			targets = append(targets, unit)
			bounceTargets = append(bounceTargets, unit)
			if len(bounceTargets) == 2 {
				break
			}
		}

		core.StartDelayedAction(sim, core.DelayedActionOptions{
			DoAt: sim.CurrentTime + time.Millisecond*100,
			OnAction: func(s *core.Simulation) {
				for _, unit := range bounceTargets {
					bounceSpell.Cast(sim, unit)
				}
			}})
	}

	bounceSpell := shadow.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 127632}.WithTag(1),
		SpellSchool:      core.SpellSchoolShadow,
		Flags:            core.SpellFlagPassiveSpell,
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1,
		CritMultiplier:   shadow.DefaultCritMultiplier(),
		BonusCoefficient: cascadeCoeff,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damageMod := 0.4 // assume minimal distance for now
			cascadeHandler(damageMod, spell, target, sim)
		},
	})

	shadow.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 127632},
		SpellSchool:  core.SpellSchoolShadow,
		Flags:        core.SpellFlagAPL,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 10,
		},
		DamageMultiplier: 1,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shadow.NewTimer(),
				Duration: time.Second * 25,
			},
		},
		ThreatMultiplier: 1,
		CritMultiplier:   shadow.DefaultCritMultiplier(),
		BonusCoefficient: cascadeCoeff,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damageMod := math.Min(0.4+0.6*(1-(30-shadow.DistanceFromTarget)/30), 1)
			targets = []*core.Unit{target}
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				cascadeHandler(damageMod, bounceSpell, target, sim)
			})
		},
	})
}
