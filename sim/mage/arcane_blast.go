package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerArcaneBlastSpell() {

	abDamageScalar := core.TernaryFloat64(mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfArcaneBlast), 0.13, 0.1)
	abDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellArcaneBlast | MageSpellArcaneExplosion,
		FloatValue: abDamageScalar,
		Kind:       core.SpellMod_DamageDone_Flat,
	})
	abCostMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellArcaneBlast,
		FloatValue: 1.5,
		Kind:       core.SpellMod_PowerCost_Pct,
	})
	abCastMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask: MageSpellArcaneBlast,
		TimeValue: time.Millisecond * -100,
		Kind:      core.SpellMod_CastTime_Flat,
	})

	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
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
			abCostMod.UpdateFloatValue(1.5 * float64(newStacks))
			abCastMod.UpdateTimeValue(time.Millisecond * time.Duration(-100*newStacks))
		},
	})

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30451},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagMage | ArcaneMissileSpells | core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneBlast,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.05,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 1.0,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.933 * mage.ScalingBaseDamage
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			mage.ArcaneBlastAura.Activate(sim)
			mage.ArcaneBlastAura.AddStack(sim)
		},
	})
}
