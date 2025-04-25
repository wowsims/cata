package assassination

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (sinRogue *AssassinationRogue) registerVenomousWounds() {
	if sinRogue.Talents.VenomousWounds == 0 {
		return
	}

	vwSpellID := 79132 + sinRogue.Talents.VenomousWounds
	vwActionID := core.ActionID{SpellID: vwSpellID}

	// https://web.archive.org/web/20111128070437/http://elitistjerks.com/f78/t105429-cataclysm_mechanics_testing/  Ctrl-F "Venomous Wounds"
	vwBaseTickDamage := 675.0
	vwMetrics := sinRogue.NewEnergyMetrics(vwActionID)
	vwProcChance := 0.3 * float64(sinRogue.Talents.VenomousWounds)

	// VW tracked via Aura instead of each bleed to keep the functionality in one place
	core.MakePermanent(sinRogue.RegisterAura(core.Aura{
		Label: "Venomous Wounds Aura",

		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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

		CritMultiplier:   sinRogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			vwDamage := vwBaseTickDamage + 0.176*spell.MeleeAttackPower()
			result := spell.CalcAndDealDamage(sim, target, vwDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				sinRogue.AddEnergy(sim, 10, vwMetrics)
			}
		},
	})
}
