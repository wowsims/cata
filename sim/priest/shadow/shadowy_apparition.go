package shadow

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

func (shadow ShadowPriest) registerShadowyApparition() {
	const apparitionScaling = 0.375
	const apparitionCoeff = 0.375

	shadow.Priest.ShadowyApparition = shadow.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 148859},
		MissileSpeed:             7,
		ProcMask:                 core.ProcMaskEmpty, // summoned guardian, should not be able to proc stuff - verify
		ClassSpellMask:           priest.PriestSpellShadowyApparation,
		Flags:                    core.SpellFlagPassiveSpell,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           shadow.DefaultCritMultiplier(),
		SpellSchool:              core.SpellSchoolShadow,
		BonusCoefficient:         apparitionCoeff,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				baseDamage := shadow.CalcScalingSpellDmg(apparitionScaling)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			})
		},
	})

	core.MakeProcTriggerAura(&shadow.Unit, core.ProcTrigger{
		Name:           "Shadowy Apparition Aura",
		Callback:       core.CallbackOnPeriodicDamageDealt,
		Outcome:        core.OutcomeCrit,
		ClassSpellMask: priest.PriestSpellShadowWordPain,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			shadow.Priest.ShadowyApparition.Cast(sim, result.Target)
		},
	})
}
