package monk

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Tooltip:
You kick upwards, dealing ${14.4*0.89*$<low>} to ${14.4*0.89*$<high>} damage and applying Mortal Wounds to the target.
Also causes all targets within 8 yards to take an increased 20% damage from your abilities for 15 sec.

-- Mortal Wounds --
Grievously wounds the target, reducing the effectiveness of any healing received for 10 sec.
-- Mortal Wounds --
*/

var risingSunKickActionID = core.ActionID{SpellID: 130320}

func risingSunKickDamageBonus(_ *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
	if !spell.Matches(MonkSpellsAll ^ MonkSpellTigerStrikes) {
		return 1.0
	}
	return 1.2
}

func risingSunKickSpellConfig(monk *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       risingSunKickActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellRisingSunKick,
		MaxRange:       core.MaxMeleeRange,

		Cast: overrides.Cast,

		DamageMultiplier: 14.4 * 0.89,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ExtraCastCondition: overrides.ExtraCastCondition,

		ApplyEffects: overrides.ApplyEffects,

		RelatedAuraArrays: overrides.RelatedAuraArrays,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
		config.Flags &= ^(core.SpellFlagAPL | SpellFlagSpender)
	}

	return config
}

func (monk *Monk) registerRisingSunKick() {
	chiMetrics := monk.NewChiMetrics(risingSunKickActionID)

	risingSunKickDebuff := monk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    fmt.Sprintf("Rising Sun Kick %s", target.Label),
			ActionID: risingSunKickActionID,
			Duration: time.Second * 15,
		}).AttachDDBC(DDBC_RisingSunKick, DDBC_Total, &monk.AttackTables, risingSunKickDamageBonus)
	})

	monk.RegisterSpell(risingSunKickSpellConfig(monk, false, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return monk.GetChi() >= 2
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				monk.SpendChi(sim, 2, chiMetrics)
				for _, target := range sim.Encounter.TargetUnits {
					risingSunKickDebuff.Get(target).Activate(sim)
				}
			}
		},
		RelatedAuraArrays: risingSunKickDebuff.ToMap(),
	}))
}

func (pet *StormEarthAndFirePet) registerSEFRisingSunKick() {

	risingSunKickDebuff := pet.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    fmt.Sprintf("Rising Sun Kick - Clone %s", target.Label),
			ActionID: risingSunKickActionID.WithTag(SEFSpellID),
			Duration: time.Second * 15,
		}).AttachDDBC(DDBC_RisingSunKickSEF, DDBC_Total, &pet.AttackTables, risingSunKickDamageBonus)
	})

	pet.RegisterSpell(risingSunKickSpellConfig(pet.owner, true, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := pet.owner.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				for _, target := range sim.Encounter.TargetUnits {
					risingSunKickDebuff.Get(target).Activate(sim)
				}
			}
		},
		RelatedAuraArrays: risingSunKickDebuff.ToMap(),
	}))

}
