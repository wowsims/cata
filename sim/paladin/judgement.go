package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) registerJudgement() {
	hpMetrics := paladin.NewHolyPowerMetrics(core.ActionID{SpellID: 105767})
	hasT132pc := paladin.HasSetBonus(ItemSetBattleplateOfRadiantGlory, 2)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20271},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc,
		Flags:       SpellFlagPrimaryJudgement | core.SpellFlagAPL | core.SpellFlagNoLogs | core.SpellFlagNoMetrics,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.05,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.paladin.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal == nil {
				return
			}

			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)

			if result.Landed() {
				paladin.CurrentJudgement.Cast(sim, target)
				if hasT132pc {
					// TODO: Measure the aura update delay distribution on PTR.
					waitTime := time.Millisecond * time.Duration(sim.RollWithLabel(150, 750, "T13 2pc"))
					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt:     sim.CurrentTime + waitTime,
						Priority: core.ActionPriorityRegen,

						OnAction: func(_ *core.Simulation) {
							paladin.GainHolyPower(sim, 1, hpMetrics)
						},
					})
				}
			}
		},
	})
}
