package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerScorchSpell() {

	if mage.Talents.CriticalMass > 0 {
		CMProcChance := float64(mage.Talents.CriticalMass) / 3.0
		//TODO double check how this works
		mage.CriticalMassAuras = mage.NewEnemyAuraArray(core.CriticalMassAura)
		mage.CritDebuffCategories = mage.GetEnemyExclusiveCategories(core.SpellCritEffectCategory)
		mage.Scorch.RelatedAuras = append(mage.Scorch.RelatedAuras, mage.CriticalMassAuras)
	}

	mage.Scorch = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 2948},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | HotStreakSpells | core.SpellFlagAPL,

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

		DamageMultipler: 1,

		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower),

		CritMultiplier: mage.DefaultSpellCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(mage.ScalingBaseDamage*0.781, mage.ScalingBaseDamage*0.781+13) + 0.512*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if sim.Proc(CMProcChance, "Critical Mass") {
				CriticalMassAuras.Get(target).Activate(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}
