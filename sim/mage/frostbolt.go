package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerFrostboltSpell() {
	hasPrimeGlyph := mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfFrostbolt)

	replProcChance := float64(mage.Talents.EnduringWinter) / 3
	var replSrc core.ReplenishmentSource
	if replProcChance > 0 {
		replSrc = mage.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 86508})
	}

	mage.Frostbolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 116},
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | ArcaneMissileSpells | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.13,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2, //TODO early winter lowers cast time of frostbolt by 0.6 seconds, then effect is inactive for 15 seconds
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(hasPrimeGlyph, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			core.TernaryFloat64(hasPrimeGlyph, .05, 0) +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.943,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.884 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if replProcChance == 1 || sim.RandomFloat("Enduring Winter") < replProcChance {
					mage.Env.Raid.ProcReplenishment(sim, replSrc)
				}
			})
		},
	})
}
