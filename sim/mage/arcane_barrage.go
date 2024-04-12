package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerArcaneBarrageSpell() {

	mage.ArcaneBarrage = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 44781},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells | core.SpellFlagAPL,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.11,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 4,
			},
		},

		DamageMultiplier: 1 +
			.02*float64(mage.Talents.TormentTheWeak) +
			core.TernaryFloat64(mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfArcaneBarrage), .04, 0),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.413*mage.ScalingBaseDamage + 0.907*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
			mage.ArcaneBlastAura.Deactivate(sim)
		},
	})
}
