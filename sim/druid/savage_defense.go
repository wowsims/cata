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

	druid.SavageDefenseAura = druid.RegisterAura(core.Aura{
		Label:    "Savage Defense",
		ActionID: core.ActionID{SpellID: 62606},
		Duration: 10 * time.Second,
	})

	var shieldStrength float64

	druid.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if druid.SavageDefenseAura.IsActive() && (result.Damage > 0) && spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
			absorbedDamage := min(shieldStrength, result.Damage)
			result.Damage -= absorbedDamage
			shieldStrength -= absorbedDamage

			if sim.Log != nil {
				druid.Log(sim, "Savage Defense absorbed %.1f damage, new Shield Strength: %.1f", absorbedDamage, shieldStrength)
			}

			if shieldStrength == 0 {
				druid.SavageDefenseAura.Deactivate(sim)
			}
		}
	})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:       "Savage Defense Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		Harmful:    true,
		ProcChance: 0.5,
		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			druid.SavageDefenseAura.Activate(sim)
			shieldStrength = 0.35 * druid.GetStat(stats.AttackPower) * (1.32 + 0.04*core.MasteryRatingToMasteryPoints(druid.GetStat(stats.MasteryRating)))

			if sim.Log != nil {
				druid.Log(sim, "Savage Defense Shield Strength: %.1f", shieldStrength)
			}
		},
	})
}
