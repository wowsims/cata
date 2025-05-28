package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const hellFireScale = 1
const hellFireCoeff = 1

func (warlock *Warlock) RegisterHellfire(callback WarlockSpellCastedCallback) *core.Spell {
	var baseDamage = warlock.CalcScalingSpellDmg(hellFireScale)
	hellfireTick := warlock.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 5857},
		SpellSchool:      core.SpellSchoolFire,
		Flags:            core.SpellFlagNoOnCastComplete,
		ProcMask:         core.ProcMaskSpellDamage,
		ClassSpellMask:   WarlockSpellHellfire,
		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		BonusCoefficient: hellFireCoeff,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]core.SpellResult, len(sim.Encounter.TargetUnits))
			for idx := 0; idx < len(sim.Encounter.TargetUnits); idx++ {
				result := spell.CalcDamage(sim, sim.Encounter.TargetUnits[idx], baseDamage, spell.OutcomeMagicHit)
				spell.DealPeriodicDamage(sim, result)
				results[idx] = *result
			}

			if callback != nil {
				callback(results, spell, sim)
			}
		},
	})

	hellfireActionID := core.ActionID{SpellID: 1949}
	manaMetric := warlock.NewManaMetrics(hellfireActionID)
	warlock.Hellfire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:         hellfireActionID,
		SpellSchool:      core.SpellSchoolFire,
		Flags:            core.SpellFlagChanneled | core.SpellFlagAPL,
		ProcMask:         core.ProcMaskSpellDamage,
		ClassSpellMask:   WarlockSpellHellfire,
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura:                 core.Aura{Label: "Hellfire"},
			TickLength:           time.Second,
			NumberOfTicks:        14,
			HasteReducesDuration: true,
			AffectedByCastSpeed:  true,
			IsAOE:                true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if warlock.CurrentHealthPercent() <= 0.04 {
					return
				}

				hellfireTick.Cast(sim, target)

				warlock.SpendMana(sim, warlock.MaxMana()*0.02, manaMetric)
				warlock.RemoveHealth(sim, warlock.MaxHealth()*0.02)

				if warlock.CurrentHealthPercent() < 0.04 {
					sim.AddPendingAction(&core.PendingAction{
						NextActionAt: sim.CurrentTime,
						Priority:     core.ActionPriorityAuto,
						OnAction: func(sim *core.Simulation) {
							dot.Deactivate(sim)
						},
					})
				}
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.CurrentHealthPercent() > 0.02
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Activate(sim)
		},

		RelatedDotSpell: hellfireTick,
	})

	return warlock.Hellfire
}
