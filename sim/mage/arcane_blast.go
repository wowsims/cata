package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const ArcaneBlastBaseCastTime = time.Millisecond * 2500

func (mage *Mage) registerArcaneBlastSpell() {
	abAuraMultiplierPerStack := core.TernaryFloat64(mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfArcaneBlast), .13, .10)
	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  core.ActionID{SpellID: 30451},
		Duration:  time.Second * 6,
		MaxStacks: 4,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := 1 + float64(oldStacks)*abAuraMultiplierPerStack
			newMultiplier := 1 + float64(newStacks)*abAuraMultiplierPerStack
			mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= newMultiplier / oldMultiplier
			mage.ArcaneBlast.CostMultiplier += 1.5 * float64(newStacks-oldStacks)
		},
	})

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30451},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | BarrageSpells | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: ArcaneBlastBaseCastTime,
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 *
			(1 + .02*float64(mage.Talents.TormentTheWeak)),

		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.933*mage.ScalingBaseDamage + 1.0*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			mage.ArcaneBlastAura.Activate(sim)
			mage.ArcaneBlastAura.AddStack(sim)
		},
	})
}
