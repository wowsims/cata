package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerScorchSpell() {

	/* 	implement when debuffs updated
	   	var CMProcChance float64
	   	if mage.Talents.CriticalMass > 0 {
	   		CMProcChance := float64(mage.Talents.CriticalMass) / 3.0
	   		//TODO double check how this works
	   		mage.CriticalMassAuras = mage.NewEnemyAuraArray(core.CriticalMassAura)
	   		mage.CritDebuffCategories = mage.GetEnemyExclusiveCategories(core.SpellCritEffectCategory)
	   		mage.Scorch.RelatedAuras = append(mage.Scorch.RelatedAuras, mage.CriticalMassAuras)
	   	} */

	mage.Scorch = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2948},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagMage | HotStreakSpells | core.SpellFlagAPL,
		ClassSpellMask: MageSpellScorch,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08 -
				0.04*float64(mage.Talents.ImprovedScorch),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		DamageMultiplier: 1,

		DamageMultiplierAdditive: 1,

		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.512,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.781 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			/*implement when debuffs updated
			if sim.Proc(CMProcChance, "Critical Mass") {
			  	mage.CriticalMassAura.Get(target).Activate(sim)
			} */
			spell.DealDamage(sim, result)
		},
	})
}
