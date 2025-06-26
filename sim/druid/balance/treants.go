package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

type BalanceTreant struct {
	*druid.DefaultTreantImpl

	Wrath *core.Spell
}

const (
	TreantWrathBonusCoeff = 0.375
	TreantWrathCoeff      = 1.875
	TreantWrathVariance   = 0.12
)

func (moonkin *BalanceDruid) newTreant() *BalanceTreant {
	treant := &BalanceTreant{
		DefaultTreantImpl: moonkin.NewDefaultTreant(druid.TreantConfig{
			NonHitExpStatInheritance: func(ownerStats stats.Stats) stats.Stats {
				return stats.Stats{
					stats.Health:           0.4 * ownerStats[stats.Health],
					stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
					stats.SpellPower:       ownerStats[stats.SpellPower],
					stats.HasteRating:      ownerStats[stats.HasteRating],
				}
			},

			EnableAutos: false,
		}),
	}

	treant.PseudoStats.DamageDealtMultiplier *= 1.091
	moonkin.AddPet(treant)

	return treant
}

func (moonkin *BalanceDruid) registerTreants() {
	for idx := range moonkin.Treants {
		moonkin.Treants[idx] = moonkin.newTreant()
	}
}

func (treant *BalanceTreant) registerWrathSpell() {
	treant.Wrath = treant.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 113769},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 20,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Millisecond * 500,
				CastTime: time.Millisecond * 2500,
			},
		},

		BonusCoefficient: TreantWrathBonusCoeff,

		DamageMultiplier: 1,

		CritMultiplier: treant.DefaultCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := treant.CalcAndRollDamageRange(sim, TreantWrathCoeff, TreantWrathVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

func (treant *BalanceTreant) Initialize() {
	treant.registerWrathSpell()
}

func (treant *BalanceTreant) ExecuteCustomRotation(sim *core.Simulation) {
	if treant.Wrath.CanCast(sim, treant.CurrentTarget) {
		treant.Wrath.Cast(sim, treant.CurrentTarget)
		return
	}
}
