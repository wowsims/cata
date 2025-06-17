package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

const (
	HealingTouchBonusCoeff = 1.86
	HealingTouchCoeff      = 18.388
	HealingTouchVariance   = 0.166
)

func (druid *Druid) registerHealingTouchSpell() {
	actionID := core.ActionID{SpellID: 5185}

	druid.HealingTouch = druid.RegisterSpell(Humanoid|Moonkin|Bear|Cat, core.SpellConfig{
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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := sim.Roll(core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, HealingTouchCoeff, HealingTouchVariance))
			baseHealing += spell.HealingPower(target) * HealingTouchBonusCoeff

			spell.CalcAndDealHealing(sim, &druid.Unit, baseHealing, spell.OutcomeHealing)
		},
	})
}
