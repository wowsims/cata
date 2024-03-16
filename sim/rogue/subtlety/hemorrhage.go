package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func (subRogue *SubtletyRogue) registerHemorrhageSpell() {
	if !subRogue.Talents.Hemorrhage {
		return
	}

	actionID := core.ActionID{SpellID: 16511}
	hasGlyph := subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfHemorrhage)
	hemoAuras := subRogue.NewEnemyAuraArray(core.HemorrhageAura)

	subRogue.Rogue.Hemorrhage = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | rogue.SpellFlagBuilder | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   subRogue.GetGeneratorCostModifier(35 - 2*float64(subRogue.Talents.SlaughterFromTheShadows)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: core.TernaryFloat64(subRogue.HasSetBonus(rogue.Tier9, 4), 5*core.CritRatingPerCritChance, 0),

		DamageMultiplier: core.TernaryFloat64(subRogue.HasDagger(core.MainHand), 3.25, 2.24) +
			core.TernaryFloat64(subRogue.HasSetBonus(rogue.Tier6, 4), 0.06, 0),
		CritMultiplier:   subRogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Hemorrhage DoT",
				Tag:      rogue.RogueBleedTag,
				ActionID: core.ActionID{SpellID: 89775},
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			subRogue.BreakStealth(sim)
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				subRogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				hemoAuras.Get(target).Activate(sim)
				if hasGlyph {
					dot := spell.Dot(target)
					dot.Spell = spell
					dot.SnapshotBaseDamage = result.Damage * .05
					dot.Apply(sim)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},

		RelatedAuras: []core.AuraArray{hemoAuras},
	})
}
