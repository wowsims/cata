package windwalker

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

/*
Tooltip:
You kick upwards, dealing ${14.4*0.89*$<low>} to ${14.4*0.89*$<high>} damage and applying Mortal Wounds to the target.
Also causes all targets within 8 yards to take an increased 20% damage from your abilities for 15 sec.

-- Mortal Wounds --
Grievously wounds the target, reducing the effectiveness of any healing received for 10 sec.
-- Mortal Wounds --
*/

func (ww *WindwalkerMonk) registerRisingSunKick() {
	actionID := core.ActionID{SpellID: 130320}
	chiMetrics := ww.NewChiMetrics(actionID)

	risingSunKickDamageBonus := func(_ *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
		if !spell.Matches(monk.MonkSpellsAll) {
			return 1.0
		}
		return 1.2
	}

	risingSunKickDebuff := ww.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Rising Sun Kick" + target.Label,
			ActionID: actionID,
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, _ *core.Simulation) {
				core.EnableDamageDoneByCaster(DDBC_RisingSunKick, DDBC_Total, ww.AttackTables[aura.Unit.UnitIndex], risingSunKickDamageBonus)
			},
			OnExpire: func(aura *core.Aura, _ *core.Simulation) {
				core.DisableDamageDoneByCaster(DDBC_RisingSunKick, ww.AttackTables[aura.Unit.UnitIndex])
			},
		})
	})

	ww.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | monk.SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellRisingSunKick,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    ww.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplier: 14.4 * 0.89,
		ThreatMultiplier: 1,
		CritMultiplier:   ww.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return ww.ComboPoints() >= 2
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := ww.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				ww.SpendChi(sim, 2, chiMetrics)
				for _, target := range sim.Encounter.TargetUnits {
					risingSunKickDebuff.Get(target).Activate(sim)
				}
			}
		},
		RelatedAuraArrays: risingSunKickDebuff.ToMap(),
	})
}
