package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerJudgment() {
	apCoef := 0.32800000906
	spCoef := 0.54600000381
	baseDamage := paladin.CalcScalingSpellDmg(spCoef)

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
			damage := baseDamage +
				apCoef*spell.MeleeAttackPower() +
				spCoef*spell.SpellPower()

			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
}
