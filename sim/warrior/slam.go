package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warrior *Warrior) RegisterSlamSpell() {
	warrior.Slam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1464},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*1500 - time.Millisecond*500*time.Duration(warrior.Talents.ImprovedSlam),
			},
			IgnoreHaste: false, // Slam now has a "Haste Affects Melee Ability Casttime" flag in cata
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if cast.CastTime > 0 {
					warrior.AutoAttacks.DelayMeleeBy(sim, cast.CastTime)
				}
			},
		},

		BonusCritRating:  core.TernaryFloat64(warrior.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfSlam), 5, 0) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.1*float64(warrior.Talents.ImprovedSlam) + 0.05*float64(warrior.Talents.WarAcademy),
		CritMultiplier:   warrior.DefaultMeleeCritMultiplier() + (0.1 * float64(warrior.Talents.Impale)),
		ThreatMultiplier: 1,
		FlatThreatBonus:  140,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bloodsurgeActive := warrior.BloodsurgeAura != nil && warrior.BloodsurgeAura.IsActive()
			bloodsurgeMultiplier := core.TernaryFloat64(bloodsurgeActive, 1.2, 1.0)
			baseDamage := 431 +
				1.1*(spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())+spell.BonusWeaponDamage())

			result := spell.CalcAndDealDamage(sim, target, baseDamage*bloodsurgeMultiplier, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() && !bloodsurgeActive {
				spell.IssueRefund(sim)
			}

			// SMF adds an OH hit to slam
			if warrior.Talents.SingleMindedFury {
				baseDamage := 431 +
					1.1*(spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())+spell.BonusWeaponDamage())

				spell.CalcAndDealDamage(sim, target, baseDamage*bloodsurgeMultiplier, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	})
}
