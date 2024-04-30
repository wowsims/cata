package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

var MutilateSpellID int32 = 1329

func (sinRogue *AssassinationRogue) newMutilateHitSpell(isMH bool) *core.Spell {
	actionID := core.ActionID{SpellID: MutilateSpellID, Tag: 1}
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		actionID = core.ActionID{SpellID: MutilateSpellID, Tag: 2}
		procMask = core.ProcMaskMeleeOHSpecial
	}
	mutBaseDamage := sinRogue.ClassSpellScaling * 0.17900000513
	t11Bonus := core.TernaryFloat64(sinRogue.HasSetBonus(rogue.Tier11, 2), 5*core.CritRatingPerCritChance, 0)

	return sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | rogue.SpellFlagBuilder | rogue.SpellFlagColdBlooded,

		BonusCritRating: 5*core.CritRatingPerCritChance*float64(sinRogue.Talents.PuncturingWounds) +
			t11Bonus,

		DamageMultiplier:         1.86, // 84 * 1.3220000267 + 75
		DamageMultiplierAdditive: 1 + 0.1*float64(sinRogue.Talents.Opportunity),
		CritMultiplier:           sinRogue.MeleeCritMultiplier(true),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = mutBaseDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = mutBaseDamage + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}

func (sinRogue *AssassinationRogue) registerMutilateSpell() {
	sinRogue.MutilateMH = sinRogue.newMutilateHitSpell(true)
	sinRogue.MutilateOH = sinRogue.newMutilateHitSpell(false)

	sinRogue.Mutilate = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: MutilateSpellID, Tag: 0},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   sinRogue.GetGeneratorCostModifier(60 - core.TernaryFloat64(sinRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfMutilate), 5, 0)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			sinRogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit) // Miss/Dodge/Parry/Hit
			if result.Landed() {
				sinRogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
				sinRogue.MutilateOH.Cast(sim, target)
				sinRogue.MutilateMH.Cast(sim, target)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
