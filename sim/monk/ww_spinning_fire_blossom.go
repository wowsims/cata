package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Tooltip:
Deals ${2.1*$<low>} to ${2.1*$<high>} Fire damage to the first enemy target in front of you within 50 yards.

If Spinning Fire Blossom travels further than 10 yards, the damage is increased by 50% and you root the target for 2 sec.
*/

var spinningFireBlossomActionID = core.ActionID{SpellID: 115073}

func spinningFireBlossomSpellConfig(monk *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       spinningFireBlossomActionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellSpinningFireBlossom,
		MissileSpeed:   20,
		MaxRange:       50,

		Cast: overrides.Cast,

		DamageMultiplier: 2.1,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(), // TODO: Spell or melee?

		ExtraCastCondition: overrides.ExtraCastCondition,

		ApplyEffects: overrides.ApplyEffects,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
		config.Flags &= ^(core.SpellFlagAPL | SpellFlagSpender)
	}

	return config
}

func (monk *Monk) registerSpinningFireBlossom() {
	chiMetrics := monk.NewChiMetrics(spinningFireBlossomActionID)

	monk.RegisterSpell(spinningFireBlossomSpellConfig(monk, false, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return monk.GetChi() >= 1
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			if target.DistanceFromTarget >= 10 {
				baseDamage *= 1.5
			}

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)
				if result.Landed() {
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
				}
			})

			monk.SpendChi(sim, 1, chiMetrics)
		},
	}))
}

func (pet *StormEarthAndFirePet) registerSpinningFireBlossom() {
	pet.RegisterSpell(spinningFireBlossomSpellConfig(pet.owner, false, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := pet.owner.CalculateMonkStrikeDamage(sim, spell)

			if target.DistanceFromTarget >= 10 {
				baseDamage *= 1.5
			}

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)
				if result.Landed() {
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
				}
			})
		},
	}))
}
