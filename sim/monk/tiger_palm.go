package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Tooltip:
Attack with the palm of your hand, dealing

-- Teachings of the Monastery --

	${6*$<low>} to ${6*$<high>}

-- else --

	${3*$<low>} to ${3*$<high>}

--

	damage.

Also grants you Tiger Power, causing your attacks to ignore 30% of enemies' armor for 20 sec.
*/
func (monk *Monk) registerTigerPalm() {
	actionID := core.ActionID{SpellID: 100787}
	chiMetrics := monk.NewChiMetrics(actionID)
	isBrewmaster := monk.Spec == proto.Spec_SpecBrewmasterMonk

	tigerPowerBuff := monk.RegisterAura(core.Aura{
		Label:    "Tiger Power" + monk.Label,
		ActionID: actionID.WithTag(2),
		Duration: time.Second * 20,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, target := range sim.Encounter.TargetUnits {
				monk.AttackTables[target.UnitIndex].ArmorIgnoreFactor += 0.3
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, target := range sim.Encounter.TargetUnits {
				monk.AttackTables[target.UnitIndex].ArmorIgnoreFactor -= 0.3
			}
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellTigerPalm,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 3.0,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return monk.ComboPoints() >= 1 || monk.ComboBreakerTigerPalmAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				tigerPowerBuff.Activate(sim)

				if monk.ComboBreakerTigerPalmAura.IsActive() || isBrewmaster {
					monk.SpendChi(sim, 0, chiMetrics)
				} else {
					monk.SpendChi(sim, 1, chiMetrics)
				}
			}

			spell.DealOutcome(sim, result)
		},
	})
}
