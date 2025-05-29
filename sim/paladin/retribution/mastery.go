package retribution

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

// Your Crusader Strike, Hammer of the Righteous, Hammer of Wrath, Templar's Verdict and Divine Storm deal ((8 + <Mastery Rating> / 600) * 1.85)% additional damage as Holy damage.
func (ret *RetributionPaladin) registerMastery() {
	handOfLight := ret.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 96172},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1.0,
		CritMultiplier:   0.0,
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := ret.HoLDamage
			if target.HasActiveAuraWithTag(core.SpellDamageEffectAuraTag) {
				baseDamage *= 1.05
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Mastery: Hand of Light" + ret.Label,
		ActionID:       core.ActionID{SpellID: 76672},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskCanTriggerHandOfLight,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HoLDamage = ret.getMasteryPercent() * result.Damage
			handOfLight.Cast(sim, result.Target)
		},
	})
}

func (ret *RetributionPaladin) getMasteryPercent() float64 {
	return ((8.0 + ret.GetMasteryPoints()) * 1.85) / 100.0
}
