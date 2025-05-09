package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
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
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// This is a dummy spell that doesn't actually heal
			// It just triggers Dream of Cenarius if the talent is selected
			spell.CalcAndDealHealing(sim, &druid.Unit, 0, spell.OutcomeHealing)
		},
	})
}
