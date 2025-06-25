package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const (
	HealingTouchBonusCoeff = 1.86
	HealingTouchCoeff      = 18.388
	HealingTouchVariance   = 0.166
)

func (druid *Druid) registerHealingTouchSpell() {
	actionID := core.ActionID{SpellID: 5185}

	druid.HealingTouch = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		ClassSpellMask: DruidSpellHealingTouch,
		Flags:          core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 28.9,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		BonusCoefficient: HealingTouchBonusCoeff,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			baseHealing := druid.CalcAndRollDamageRange(sim, HealingTouchCoeff, HealingTouchVariance)
			spell.CalcAndDealHealing(sim, spell.Unit, baseHealing, spell.OutcomeHealingCrit)
		},
	})
}
