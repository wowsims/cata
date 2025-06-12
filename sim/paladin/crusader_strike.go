package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
An instant strike that causes 125% weapon damage plus 633 and grants a charge of Holy Power.

Applies the Weakened Blows effect.

Weakened Blows
Demoralizes the target, reducing their physical damage dealt by 10% for 30 sec.
*/
func (paladin *Paladin) registerCrusaderStrike() {
	actionID := core.ActionID{SpellID: 35395}
	paladin.CanTriggerHolyAvengerHpGain(actionID)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCrusaderStrike,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 10,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.BuilderCooldown(),
				Duration: time.Millisecond * 4500,
			},
		},

		DamageMultiplier: 1.25,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + paladin.CalcScalingSpellDmg(0.55400002003)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.HolyPower.Gain(sim, 1, actionID)
			}

			spell.DealDamage(sim, result)
		},
	})
}
