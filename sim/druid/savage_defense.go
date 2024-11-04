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
		func(_ *core.Unit) float64 {
			freshShieldStrength := 0.35 * druid.GetStat(stats.AttackPower) * (1.32 + 0.04*druid.GetMasteryPoints())

			if druid.BlazeOfGloryAura.IsActive() {
				freshShieldStrength *= 1.0 + 0.2*float64(druid.BlazeOfGloryAura.GetStacks())
			}

			return freshShieldStrength
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
