package combat

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) applyMastery() {
	mgAttack := comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 86392},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty, // MG Appears to be unable to proc anything EXCEPT poisons. This specific case is handled by Poisons directly.
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		ClassSpellMask: rogue.RogueSpellMainGauche,

		DamageMultiplier:         1.2,
		DamageMultiplierAdditive: 1.0,
		CritMultiplier:           comRogue.CritMultiplier(false),
		ThreatMultiplier:         1.0,

		BonusCoefficient: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	core.MakeProcTriggerAura(&comRogue.Unit, core.ProcTrigger{
		Name:     "Mastery: Main Gauche",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		ProcMask: core.ProcMaskMeleeMH | core.ProcMaskMeleeProc,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			if spell == comRogue.Rupture {
				return false
			}

			// Implement the proc in here so we can get the most up to date proc chance from mastery
			return sim.Proc(comRogue.GetMasteryBonus(), "Main Gauche")
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			mgAttack.Cast(sim, result.Target)
		},
	})
}
