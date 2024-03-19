package marksmanship

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *MarksmanshipHunter) registerAimedShotSpell(timer *core.Timer) {

	hunter.AimedShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49050},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 50 - (float64(hunter.Talents.Efficiency) * 2),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfAimedShot), time.Second*2, 0),
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(hunter.Talents.TrueshotAura, 10*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1,
		DamageMultiplier: 1,
		CritMultiplier:   hunter.CritMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := hunter.AutoAttacks.Ranged().CalculateWeaponDamage(sim, spell.RangedAttackPower(target))
			rap := spell.RangedAttackPower(target) * 0.724 + 766
			baseDamage := ((wepDmg + rap) * 1.6) + 100
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})
}
