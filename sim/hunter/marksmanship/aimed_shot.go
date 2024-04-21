package marksmanship

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mmHunter *MarksmanshipHunter) registerAimedShotSpell() {
	if mmHunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfAimedShot) {
		focusMetrics := mmHunter.NewFocusMetrics(core.ActionID{SpellID: 42897})
		mmHunter.RegisterAura(core.Aura{
			Label: "Glyph of Aimed Shot",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == mmHunter.AimedShot && result.DidCrit() {
					mmHunter.AddFocus(sim, 5, focusMetrics)
				}
			},
		})
	}
	mmHunter.AimedShot = mmHunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 19434},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		MissileSpeed: 40,
		FocusCost: core.FocusCostOptions{
			Cost: 50 - (float64(mmHunter.Talents.Efficiency) * 2),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Second * 3,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / mmHunter.RangedSwingSpeed())
			},
		},
		BonusCritRating:  0,
		DamageMultiplier: 1.32,
		CritMultiplier:   mmHunter.CritMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := mmHunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target))
			rap := spell.RangedAttackPower(target) * 0.724
			baseDamage := (wepDmg + rap) + 821
			if sim.IsExecutePhase90() {
				spell.BonusCritRating = (30.0 * float64(mmHunter.Talents.CarefulAim)) * core.CritRatingPerCritChance
			}
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			if sim.IsExecutePhase90() {
				spell.BonusCritRating = 0
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
