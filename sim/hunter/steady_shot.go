package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerSteadyShotSpell() {
	
	ssMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 49052})
	if hunter.Talents.ImprovedSteadyShot > 0 {
		hunter.ImprovedSteadyShotAura = hunter.RegisterAura(core.Aura{
			Label:    "Improved Steady Shot",
			ActionID: core.ActionID{SpellID: 53220},
			Duration: time.Second * 8,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				// Todo Apply 20% ranged attack speed per point when ss is used two times in a row
				
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {

			},
		})
	}
	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49052},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Second * 2,
			},
			IgnoreHaste: false, // Hunter GCD is locked at 1.5s //Todo: no longer in Cata, its 1sec // should probably be affected by haste now

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},

		// BonusCritRating: 0 +
		// 	2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts),
		DamageMultiplierAdditive: 1,
		// DamageMultiplier: 1 *
		// 	hunter.markedForDeathMultiplier(),
		CritMultiplier:1,//   hunter.critMultiplier(true, true, false), // what is this
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.21 * spell.RangedAttackPower(target) +
				hunter.AutoAttacks.Ranged().BaseDamage(sim)*2.8/hunter.AutoAttacks.Ranged().SwingSpeed + 280
			hunter.AddFocus(sim, 6, ssMetrics)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}