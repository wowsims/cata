package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"time"
)

func (paladin *Paladin) applyGlyphs() {
	// Prime Glyphs
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfCrusaderStrike) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Rating,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: 5 * core.CritRatingPerCritChance,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfJudgement) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  SpellMaskJudgement,
			FloatValue: 0.1,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfTemplarSVerdict) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  SpellMaskTemplarsVerdict,
			FloatValue: 0.15,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfSealOfTruth) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusExpertise_Rating,
			ClassMask:  SpellMaskSealOfTruth,
			FloatValue: 10 * core.ExpertisePerQuarterPercentReduction,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfExorcism) {
		registerGlyphOfExorcism(paladin)
	}

	// Major Glyphs
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfWrath) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Pct,
			ClassMask:  SpellMaskHammerOfWrath,
			FloatValue: -1,
		})
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_Cooldown_Multiplier,
			ClassMask:  SpellMaskConsecration,
			FloatValue: 0.2,
		})
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DotNumberOfTicks_Flat,
			ClassMask:  SpellMaskConsecration,
			FloatValue: 2,
		})
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfTheAsceticCrusader) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Pct,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: -0.3,
		})
	}
}

func registerGlyphOfExorcism(paladin *Paladin) {
	exorcismAverageDamage :=
		core.CalcScalingSpellAverageEffect(proto.Class_ClassPaladin, 2.663)

	glyphOfExorcismDot := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 54934},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: SpellMaskGlyphOfExorcism,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Exorcism (DoT)",
			},
			NumberOfTicks:        3,
			HasteAffectsDuration: false,
			TickLength:           2 * time.Second,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := exorcismAverageDamage +
					0.344*max(dot.Spell.SpellPower(), dot.Spell.MeleeAttackPower())

				// Total damage is 20% of an average hit
				baseDamage *= 0.2

				// Damage is spread over 3 ticks
				dot.SnapshotBaseDamage = baseDamage / 3
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Exorcism",
		ActionID:       core.ActionID{SpellID: 54934},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskExorcism,

		ProcChance: 1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				glyphOfExorcismDot.Dot(result.Target).Apply(sim)
			}
		},
	})
}
