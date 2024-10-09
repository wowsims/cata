package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (druid *Druid) registerSavageDefensePassive() {
	if !druid.InForm(Bear) {
		return
	}

	savageDefenseAura := druid.NewDamageAbsorptionAuraForSchool(
		"Savage Defense",
		core.ActionID{SpellID: 62606},
		10*time.Second,
		core.SpellSchoolPhysical,
		func(unit *core.Unit) float64 {
			return 0.35 * druid.GetStat(stats.AttackPower) * (1.32 + 0.04*core.MasteryRatingToMasteryPoints(druid.GetStat(stats.MasteryRating)))
		})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:       "Savage Defense Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		Harmful:    true,
		ProcChance: 0.5,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			savageDefenseAura.Activate(sim)
		},
	})
}
