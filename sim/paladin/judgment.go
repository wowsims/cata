package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
A magic attack that unleashes the energy of a Seal to cause (623 + 0.328 * <AP> + 0.546 * <SP>) Holy damage

-- Judgments of the Wise --
and generates one charge of Holy Power
-- /Judgments of the Wise --

-- Judgments of the Bold --
generate one charge of Holy Power, and apply the Physical Vulnerability debuff to a target

Physical Vulnerability
Weakens the constitution of an enemy target, increasing their physical damage taken by 4% for 30 sec
-- /Judgments of the Bold --

.
*/
func (paladin *Paladin) registerJudgment() {
	paladin.Judgment = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20271},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: SpellMaskJudgment,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,

		MaxRange: 30,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcScalingSpellDmg(0.54600000381) +
				0.32800000906*spell.MeleeAttackPower() +
				0.54600000381*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
}
