// Implements the shadow priest's mastery
// Every tick of a priest's DoT can be replicated
// The chance is based on the mastery a priest has
package shadow

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

func (shadow *ShadowPriest) registerShadowyRecall() {
	swpDupe := shadow.buildSingleTickSpell(
		589,
		priest.SwpScaleCoeff,
		priest.SwpSpellCoeff,
		0,
		priest.PriestSpellShadowWordPain,
		nil,
	)
	mfDupe := shadow.buildSingleTickSpell(
		15407,
		MfScale,
		MfCoeff,
		0,
		priest.PriestSpellMindFlay,
		nil,
	)
	vtDupe := shadow.buildSingleTickSpell(
		34914,
		priest.VtScaleCoeff,
		priest.VtSpellCoeff,
		0,
		priest.PriestSpellVampiricTouch,
		nil,
	)
	searDupe := shadow.buildSingleTickSpell(
		48045,
		priest.SearScale,
		priest.SearCoeff,
		priest.SearVariance,
		priest.PriestSpellMindSear,
		nil,
	)
	dpDupe := shadow.buildSingleTickSpell(
		2944,
		DpDotScale,
		DpDotCoeff,
		0,
		priest.PriestSpellDevouringPlagueDoT,
		func() float64 {
			return float64(shadow.orbsConsumed)
		},
	)

	spellList := []*core.Spell{swpDupe, mfDupe, vtDupe, dpDupe}
	core.MakePermanent(shadow.RegisterAura(core.Aura{
		Label: "Shadowy Recall (Mastery)",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || result.Damage == 0 {
				return
			}

			for _, dupeSpell := range spellList {
				classMask := dupeSpell.ClassSpellMask &^ priest.PriestSpellShadowyRecall
				if classMask == spell.ClassSpellMask && spell != dupeSpell {
					if sim.Proc((shadow.GetMasteryPoints()*1.8+8*1.8)/100, "Shadowy Recall (Proc)") {
						dupeSpell.Cast(sim, result.Target)
					}

					return
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || result.Damage == 0 || spell.ClassSpellMask != priest.PriestSpellMindSear || spell == searDupe {
				return
			}

			if sim.Proc((shadow.GetMasteryPoints()*1.8+8*1.8)/100, "Shadowy Recall (Proc)") {
				searDupe.Cast(sim, result.Target)
			}
		},
	}))
}

func (shadow *ShadowPriest) buildSingleTickSpell(spellId int32, scale float64, coeff float64, variance float64, classMask int64, customDamageMod func() float64) *core.Spell {
	return shadow.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: spellId}.WithTag(77486),
		SpellSchool:      core.SpellSchoolShadow,
		Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		BonusCoefficient: coeff,
		CritMultiplier:   shadow.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		ClassSpellMask:   classMask | priest.PriestSpellShadowyRecall,
		ProcMask:         core.ProcMaskSpellDamage,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if customDamageMod != nil {
				spell.DamageMultiplier *= customDamageMod()
			}
			baseDamage := core.TernaryFloat64(
				variance > 0,
				shadow.CalcAndRollDamageRange(sim, scale, variance),
				shadow.CalcScalingSpellDmg(scale),
			)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCritNoHitCounter)
			if customDamageMod != nil {
				spell.DamageMultiplier /= customDamageMod()
			}
			if result.Landed() {
				if result.DidCrit() {
					spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
				} else {
					spell.SpellMetrics[result.Target.UnitIndex].Ticks++
				}
			}

			spell.DealPeriodicDamage(sim, result)
		},
	})
}
