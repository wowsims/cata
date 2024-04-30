package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerKillShotSpell() {
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfKillShot) {
		icd := core.Cooldown{
			Timer:    hunter.NewTimer(),
			Duration: time.Second * 6,
		}
		hunter.RegisterAura(core.Aura{
			Label:    "Kill Shot Glyph",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.KillShot {
					if icd.IsReady(sim) {
						icd.Use(sim)
						hunter.KillShot.CD.Reset()
					}
				}
			},
		})
	}

	hunter.KillShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 53351},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		MissileSpeed: 40,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20()
		},

		BonusCritRating:  0 + 5*core.CritRatingPerCritChance*float64(hunter.Talents.SniperTraining),
		DamageMultiplier: 1.5, //
		CritMultiplier:   hunter.CritMultiplier(true, true, false),
		ThreatMultiplier: 1,
		// https://web.archive.org/web/20120207222124/http://elitistjerks.com/f74/t110306-hunter_faq_cataclysm_edition_read_before_asking_questions/
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// (100% weapon dmg + 45% RAP + 543) * 150%
			normalizedWeaponDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target))
			rapBonusDamage := spell.RangedAttackPower(target) * 0.45
			flatBonus := 543.0

			baseDamage := normalizedWeaponDamage + rapBonusDamage + flatBonus
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
