package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerArcaneBlastSpell() {
	abDamageScalar := core.TernaryFloat64(mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfArcaneBlast), 0.13, 0.1)
	abDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellArcaneBlast | MageSpellArcaneExplosion,
		FloatValue: abDamageScalar,
		Kind:       core.SpellMod_DamageDone_Flat,
	})
	abCostMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellArcaneBlast,
		IntValue:  150,
		Kind:      core.SpellMod_PowerCost_Pct,
	})
	abCastMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellArcaneBlast,
		TimeValue: time.Millisecond * -100,
		Kind:      core.SpellMod_CastTime_Flat,
	})

	arcaneBlastAura := mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast Debuff",
		ActionID:  core.ActionID{SpellID: 36032},
		Duration:  time.Second * 6,
		MaxStacks: 4,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			abDamageMod.Activate()
			abCostMod.Activate()
			abCastMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			abDamageMod.Deactivate()
			abCostMod.Deactivate()
			abCastMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			abDamageMod.UpdateFloatValue(abDamageScalar * float64(newStacks))
			abCostMod.UpdateIntValue(150 * newStacks)
			abCastMod.UpdateTimeValue(time.Millisecond * time.Duration(-100*newStacks))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&(MageSpellArcaneBlast|MageSpellArcaneExplosion) > 0 {
				return
			}
			if !spell.SpellSchool.Matches(core.SpellSchoolArcane) {
				return
			}
			if spell.ProcMask != core.ProcMaskSpellDamage {
				return
			}
			aura.Deactivate(sim)
		},
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30451},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneBlast,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: 1.0,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.933 * mage.ClassSpellScaling
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				arcaneBlastAura.Activate(sim)
				arcaneBlastAura.AddStack(sim)
			}
		},
	})
}
