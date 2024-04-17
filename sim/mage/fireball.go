package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerFireballSpell() {

	mage.Fireball = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 133},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagMage | ArcaneMissileSpells | HotStreakSpells | core.SpellFlagAPL,
		ClassSpellMask: MageSpellFireball,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.09,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				CastTime: time.Millisecond *
					time.Duration((2500 * core.TernaryFloat64(mage.HasSetBonus(ItemSetFirelordsVestments, 4), 0.9, 1))),
			},
		},
		BonusCritRating: 0 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 1.236,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.20 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
