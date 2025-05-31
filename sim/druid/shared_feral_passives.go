package druid

import (
	"github.com/wowsims/mop/sim/core"
)

const RendAndTearBonusCritPercent = 35.0
const RendAndTearDamageMultiplier = 1.2

// Modifies the Bleed aura to apply the bonus.
func (druid *Druid) applyRendAndTear(aura core.Aura) core.Aura {
	if druid.AssumeBleedActive {
		return aura
	}

	affectedSpells := []*DruidSpell{druid.Maul, druid.SwipeBear, druid.SwipeCat}

	aura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) {
		if druid.BleedsActive == 0 {
			if druid.FerociousBite != nil {
				druid.FerociousBite.BonusCritPercent += RendAndTearBonusCritPercent
			}

			for _, spell := range affectedSpells {
				if spell != nil {
					spell.DamageMultiplier *= RendAndTearDamageMultiplier
				}
			}
		}
		druid.BleedsActive++
	})
	aura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
		druid.BleedsActive--
		if druid.BleedsActive == 0 {
			if druid.FerociousBite != nil {
				druid.FerociousBite.BonusCritPercent -= RendAndTearBonusCritPercent
			}

			for _, spell := range affectedSpells {
				if spell != nil {
					spell.DamageMultiplier /= RendAndTearDamageMultiplier
				}
			}
		}
	})

	return aura
}

func (druid *Druid) ApplyPrimalFury() {
	actionID := core.ActionID{SpellID: 16961}
	rageMetrics := druid.NewRageMetrics(actionID)
	cpMetrics := druid.NewComboPointMetrics(actionID)

	druid.RegisterAura(core.Aura{
		Label:    "Primal Fury",
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}

			if druid.InForm(Bear) {
				if (spell == druid.MHAutoSpell) || druid.MangleBear.IsEqual(spell) {
					druid.AddRage(sim, 15, rageMetrics)
				}
			} else if druid.InForm(Cat) {
				if druid.MangleCat.IsEqual(spell) || druid.Shred.IsEqual(spell) || druid.Rake.IsEqual(spell) || druid.Ravage.IsEqual(spell) {
					druid.AddComboPoints(sim, 1, cpMetrics)
				}
			}
		},
	})
}
