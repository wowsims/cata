package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerScorchSpell() {
	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2948},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          HotStreakSpells | core.SpellFlagAPL,
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
			baseDamage := 0.781 * mage.ClassSpellScaling
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			/*implement when debuffs updated
			if sim.Proc(CMProcChance, "Critical Mass") {
			  	mage.CriticalMassAura.Get(target).Activate(sim)
			} */
			spell.DealDamage(sim, result)
		},
	})
}
