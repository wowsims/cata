package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerHemorrhageSpell() {
	hemoActionID := core.ActionID{SpellID: 16511, Tag: 1}
	hemoDotActionID := core.ActionID{SpellID: 16511}
	hasMinorGlyph := subRogue.HasMinorGlyph(proto.RogueMinorGlyph_GlyphOfHemorrhage)
	var lastHemoDamage float64

	weaponDamage := 1.6

	hemoDot := subRogue.RegisterSpell(core.SpellConfig{
		ActionID:    hemoDotActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagPassiveSpell, // From initial testing, Hemo DoT only benefits from debuffs on target, such as 30% bleed damage

		ThreatMultiplier: 1,
		CritMultiplier:   subRogue.CritMultiplier(false), // Per WoWHead data, Lethality does not boost the DoT directly,
		DamageMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Hemorrhage DoT",
				Tag:      rogue.RogueBleedTag,
				ActionID: hemoDotActionID,
				Duration: time.Second * 24,
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotPhysical(target, lastHemoDamage*0.0625)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	subRogue.Rogue.Hemorrhage = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:       hemoActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellHemorrhage,

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 700,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: core.TernaryFloat64(subRogue.HasDagger(core.MainHand), weaponDamage*1.45, weaponDamage),
		CritMultiplier:   subRogue.CritMultiplier(true),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			subRogue.BreakStealth(sim)
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				subRogue.AddComboPointsOrAnticipation(sim, 1, spell.ComboPointMetrics())
				lastHemoDamage = result.Damage

				if hasMinorGlyph { // Prevents triggering the DoT unless Garrote/Rupture/Crimson Tempest are active
					if target.HasAuraWithTag(rogue.RogueBleedTag) {
						hemoDot.Cast(sim, target)
					}
				} else {
					hemoDot.Cast(sim, target)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})

	subRogue.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(s *core.Simulation, slot proto.ItemSlot) {
		// Recalculate Hemorrhage's multiplier in case the MH weapon changed.
		subRogue.Hemorrhage.DamageMultiplier = core.TernaryFloat64(subRogue.HasDagger(core.MainHand), weaponDamage*1.45, weaponDamage)
	})
}
