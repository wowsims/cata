package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (shaman *Shaman) RegisterElementalBlastSpell() {

	masteryAura := shaman.NewTemporaryStatsAura("Elemental Blast Mastery", core.ActionID{SpellID: 118522}, stats.Stats{stats.MasteryRating: 3500}, time.Second*8)
	hasteAura := shaman.NewTemporaryStatsAura("Elemental Blast Mastery", core.ActionID{SpellID: 118522}, stats.Stats{stats.HasteRating: 3500}, time.Second*8)
	critAura := shaman.NewTemporaryStatsAura("Elemental Blast Mastery", core.ActionID{SpellID: 118522}, stats.Stats{stats.CritRating: 3500}, time.Second*8)

	shaman.ElementalBlast = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 117014},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   40,
		ClassSpellMask: SpellMaskElementalBlast,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Second * 2,
				GCD:      time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 12,
			},
		},
		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: 2.11199998856,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shaman.CalcAndRollDamageRange(sim, 4.23999977112, 0.15000000596)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			rand := sim.RandomFloat("Elemental Blast buff")
			if rand < 1.0/3.0 {
				masteryAura.Activate(sim)
			} else {
				if rand < 2.0/3.0 {
					hasteAura.Activate(sim)
				} else {
					critAura.Activate(sim)
				}
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
