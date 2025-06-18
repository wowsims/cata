package assassination

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (sinRogue *AssassinationRogue) registerVenomousWounds() {
	vwActionID := core.ActionID{SpellID: 79134}

	vwBaseTickDamage := sinRogue.GetBaseDamageFromCoefficient(0.55000001192)
	vwAPCoeff := 0.15999999642
	vwMetrics := sinRogue.NewEnergyMetrics(vwActionID)
	vwProcChance := 0.75

	// VW tracked via Aura instead of each bleed to keep the functionality in one place
	core.MakePermanent(sinRogue.RegisterAura(core.Aura{
		Label: "Venomous Wounds Aura",

		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// If the target has both Rupture and Garrote, Garrote cannot trigger VW
			if spell == sinRogue.Garrote && result.Target.HasActiveAura("Rupture") {
				return
			}

			if spell == sinRogue.Rupture || spell == sinRogue.Garrote {
				if sim.Proc(vwProcChance, "Venomous Wounds") {
					// Trigger VW after small delay to prevent aura refresh loops
					// https://i.gyazo.com/dc845a371102294abfb207c6fd586bfa.png
					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt:     sim.CurrentTime + 1,
						Priority: core.ActionPriorityDOT,
						OnAction: func(s *core.Simulation) {
							sinRogue.VenomousWounds.Cast(sim, result.Target)
						},
					})
				}
			}
		},
	}))

	sinRogue.VenomousWounds = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:       vwActionID,
		ClassSpellMask: rogue.RogueSpellVenomousWounds,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell,

		CritMultiplier:   sinRogue.CritMultiplier(false),
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			vwDamage := vwBaseTickDamage + vwAPCoeff*spell.MeleeAttackPower()
			result := spell.CalcAndDealDamage(sim, target, vwDamage, spell.OutcomeMeleeSpecialCritOnly)
			if result.Landed() {
				sinRogue.AddEnergy(sim, 10, vwMetrics)
			}
		},
	})
}
