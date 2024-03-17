package marksmanship

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *MarksmanshipHunter) registerChimeraShotSpell() {
	if !hunter.Talents.ChimeraShot {
		return
	}

	ssProcSpell := hunter.chimeraShotSerpentStingSpell()

	hunter.ChimeraShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53209},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 50 - (float64(hunter.Talents.Efficiency) * 2),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfChimeraShot), time.Second*1, 0),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   hunter.CritMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.2*spell.RangedAttackPower(target) +
				hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				spell.BonusWeaponDamage()
			baseDamage *= 1.25

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			if result.Landed() {
				if hunter.SerpentSting.Dot(target).IsActive() {
					hunter.SerpentSting.Dot(target).Rollover(sim)
					ssProcSpell.Cast(sim, target)
				} else if hunter.ScorpidStingAuras.Get(target).IsActive() {
					hunter.ScorpidStingAuras.Get(target).Refresh(sim)
				}
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (hunter *MarksmanshipHunter) chimeraShotSerpentStingSpell() *core.Spell {
	return hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53353},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplierAdditive: 1 +
			0.15*float64(hunter.Talents.ImprovedSerpentSting),
		DamageMultiplier: 1 *
			(2.0 + core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSerpentSting), 0.8, 0)),
		CritMultiplier:   hunter.CritMultiplier(true, false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 242 + 0.04*spell.RangedAttackPower(target)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedCritOnly)
		},
	})
}
