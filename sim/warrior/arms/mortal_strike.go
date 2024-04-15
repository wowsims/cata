package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (war *ArmsWarrior) RegisterMortalStrikeSpell() {
	war.mortalStrike = war.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 12294},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagMeleeMetrics,

		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Millisecond * 4500,
			},
		},

		CritMultiplier:  war.DefaultMeleeCritMultiplier() + (0.1 * float64(war.Talents.Impale)),
		BonusCritRating: (5.0 * float64(war.Talents.Cruelty)) * core.CritRatingPerCritChance,
		DamageMultiplier: 1.0 + core.TernaryFloat64(war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfMortalStrike), 0.1, 0.0) +
			0.05*float64(war.Talents.WarAcademy),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 423 +
				0.8*(spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())+spell.BonusWeaponDamage())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			} else {
				war.TriggerSlaughter(sim, target)
				if result.DidCrit() {
					war.TriggerWreckingCrew(sim)
				}
			}
		},
	})
}
