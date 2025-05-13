package shadow

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

const haloScale = 19.266
const haloVariance = 0.5
const haloCoeff = 1.95

func (shadow *ShadowPriest) registerHalo() {
	if !shadow.Talents.Halo {
		return
	}
	shadow.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 120696},
		SpellSchool:      core.SpellSchoolShadow,
		Flags:            core.SpellFlagAPL,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: haloCoeff,
		ClassSpellMask:   priest.PriestSpellHalo,
		CritMultiplier:   shadow.DefaultCritMultiplier(),
		MissileSpeed:     10,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 13.5,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shadow.NewTimer(),
				Duration: time.Second * 40,
			},
		},
		ProcMask: core.ProcMaskSpellDamage,
		MaxRange: 30,
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				baseDamage := shadow.CalcAndRollDamageRange(sim, haloScale, haloVariance)
				distMod := calcHaloMod(shadow.DistanceFromTarget)
				spell.DamageMultiplier *= distMod
				for _, target := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
				spell.DamageMultiplier /= distMod
			})
		},
	})
}

// https://web.archive.org/web/20120626065654/http://us.battle.net/wow/en/forum/topic/5889309137?page=5#97
func calcHaloMod(distance float64) float64 {
	return 0.5*math.Pow(1.01, -1*math.Pow(((distance-25)/2), 4)) + 0.1 + 0.015*distance
}
