package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) RegisterJudgement() {
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20271},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc,
		Flags:       SpellFlagPrimaryJudgement | core.SpellFlagAPL,

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
			// if paladin.CurrentSeal.IsActive() {
			// 	paladin.CurrentSeal.Refresh(sim)
			// } else {
			// 	paladin.CurrentSeal.Activate(sim)
			// }

			paladin.CurrentJudgement.Cast(sim, target)
		},
	})
}
