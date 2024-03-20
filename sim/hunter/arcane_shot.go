package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerArcaneShotSpell() {

	dmgMultiplier := 0.61
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfArcaneShot) {
		dmgMultiplier *= 0.12
	}
	hunter.ArcaneShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 3044},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 25 - float64(hunter.Talents.Efficiency),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritRating:          0,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         dmgMultiplier,
		CritMultiplier:           hunter.CritMultiplier(true, true, false),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target))
			baseDamage := wepDmg + (0.0483 * spell.RangedAttackPower(target))
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.DealDamage(sim, result)
		},
	})
}
