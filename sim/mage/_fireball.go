package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerFireballSpell() {
	hasPrimeGlyph := mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfFireball)

	mage.Fireball = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 133},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells | HotStreakSpells | core.SpellFlagAPL,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.09,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				CastTime: time.Millisecond*2500 *
				core.TernaryFloat64(mage.HasSetBonus(ItemSetFirelordsVestments, 4), 0.9, 1)
			},
		},

		BonusCritRating: 0 +
			// waiting to see how buff talents will be implemented
			//float64(mage.Talents.PiercingIce)*core.CritRatingPerCritChance + 
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0) + 
			core.TernaryFloat64(hasPrimeGlyph, 5*core.CritRatingPerCritChance, 0)

		DamageMultiplier: 1

		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower) +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),

		CritMultiplier: mage.DefaultSpellCritMultiplier();

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(mage.ScalingBaseDamage*1.20, mage.ScalingBaseDamage*1.20+13) + 1.2359*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
