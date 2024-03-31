package assassination

import "github.com/wowsims/cata/sim/core"

func (sinRogue *AssassinationRogue) registerVenomousWounds() {
	if sinRogue.Talents.VenomousWounds == 0 {
		return
	}

	vwSpellID := 79132 + sinRogue.Talents.VenomousWounds
	vwActionID := core.ActionID{SpellID: vwSpellID}

	// https://web.archive.org/web/20111128070437/http://elitistjerks.com/f78/t105429-cataclysm_mechanics_testing/  Ctrl-F "Venomous Wounds"
	vwBaseTickDamage := 675.0
	vwMetrics := sinRogue.NewEnergyMetrics(vwActionID)

	sinRogue.VenomousWounds = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:         vwActionID,
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellDamage,
		CritMultiplier:   sinRogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			vwDamage := vwBaseTickDamage + 0.176*spell.MeleeAttackPower()
			spell.CalcAndDealDamage(sim, target, vwDamage, spell.OutcomeMagicCrit)
			sinRogue.AddEnergy(sim, 10, vwMetrics)
		},
	})
}
