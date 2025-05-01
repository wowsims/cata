package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerPurifyingBrew() {
	actionID := core.ActionID{SpellID: 119582}
	chiMetrics := bm.NewChiMetrics(actionID)

	bm.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          monk.SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellPurifyingBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    bm.NewTimer(),
				Duration: time.Second * 1,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bm.StanceMatches(monk.SturdyOx) && bm.ComboPoints() >= 1
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bm.RefreshStagger(sim, &bm.Unit, 0.0)
			bm.SpendChi(sim, 1, chiMetrics)
		},
	})
}
