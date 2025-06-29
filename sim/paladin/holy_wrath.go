package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (paladin *Paladin) registerHolyWrath() {
	hwAvgDamage := core.CalcScalingSpellAverageEffect(proto.Class_ClassPaladin, 2.33299994469)

	paladin.HolyWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2812},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHolyWrath,

		MissileSpeed: 20,
		MaxRange:     10,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 15 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numTargets := sim.Environment.ActiveTargetCount()
			results := make([]*core.SpellResult, numTargets)
			baseDamage := hwAvgDamage + .61*spell.SpellPower()

			// Damage is split between all mobs, each hit rolls for hit/crit separately
			baseDamage /= float64(numTargets)

			for idx, currentTarget := range sim.Environment.GetActiveTargets() {
				results[idx] = spell.CalcDamage(sim, &currentTarget.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				for idx := int32(0); idx < numTargets; idx++ {
					spell.DealDamage(sim, results[idx])
				}
			})
		},
	})
}
