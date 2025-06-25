package blood

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (bdk *BloodDeathKnight) registerRiposte() {
	riposteAura := bdk.RegisterAura(core.Aura{
		Label:     "Riposte" + bdk.Label,
		ActionID:  core.ActionID{SpellID: 145677},
		Duration:  time.Second * 20,
		MaxStacks: math.MaxInt32,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			bdk.AddStatDynamic(sim, stats.CritRating, float64(newStacks-oldStacks))
		},
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:     "Riposte Trigger" + bdk.Label,
		ActionID: core.ActionID{SpellID: 145676},
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeDodge | core.OutcomeParry,
		ICD:      time.Second * 1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bonusCrit := math.Round((bdk.GetStat(stats.DodgeRating) + bdk.GetParryRatingWithoutStrength()) * 0.75)
			riposteAura.Activate(sim)
			riposteAura.SetStacks(sim, int32(bonusCrit))
		},
	})
}
