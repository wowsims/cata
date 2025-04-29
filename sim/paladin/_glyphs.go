package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) applyGlyphs() {
	// Prime Glyphs
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfCrusaderStrike) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: 5,
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
			FloatValue: 10 * core.ExpertisePerQuarterPercentReduction,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfExorcism) {
		registerGlyphOfExorcism(paladin)
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfHammerOfTheRighteous) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  SpellMaskHammerOfTheRighteous,
			FloatValue: 0.1,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfShieldOfTheRighteous) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  SpellMaskShieldOfTheRighteous,
			FloatValue: 0.1,
		})
	}

	// Major Glyphs
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfWrath) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Pct,
			ClassMask: SpellMaskHammerOfWrath,
			IntValue:  -100,
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
			Kind:      core.SpellMod_PowerCost_Pct,
			ClassMask: SpellMaskCrusaderStrike,
			IntValue:  -30,
		})
	}
}

func registerGlyphOfExorcism(paladin *Paladin) {
	glyphOfExoBaseDamage :=
		core.CalcScalingSpellAverageEffect(proto.Class_ClassPaladin, 0.17900000513)

	// Used for checking "Is Aura Known" in the APL
	paladin.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 54934},
		Label:    "Glyph of Exorcism (DoT)" + paladin.Label,
	})

	glyphOfExorcismDot := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 879}.WithTag(3), // actual 54934
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: SpellMaskGlyphOfExorcism,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			NumberOfTicks:       3,
			AffectedByCastSpeed: false,
			TickLength:          2 * time.Second,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := glyphOfExoBaseDamage + 0.344*max(dot.Spell.SpellPower(), dot.Spell.MeleeAttackPower())*0.2/3

				bonusCritPercent := dot.Spell.BonusCritPercent
				if target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead {
					dot.Spell.BonusCritPercent += 100
				}

				dot.Snapshot(target, baseDamage)

				dot.Spell.BonusCritPercent = bonusCritPercent
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Exorcism" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 54934},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskExorcism,
		Outcome:        core.OutcomeLanded,
		ProcChance:     1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			glyphOfExorcismDot.Cast(sim, result.Target)
		},
	})
}
