package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerRevenge() {
	actionID := core.ActionID{SpellID: 6572}
	rageMetrics := war.NewRageMetrics(actionID)

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskRevenge,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 9,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			chainMultiplier := 1.0
			aoeTarget := target
			hitLanded := false

			for idx := range war.Env.GetNumTargets() {
				if idx >= 3 {
					break
				}
				baseDamage := chainMultiplier * (war.CalcAndRollDamageRange(sim, 7.5, 0.20000000298) + spell.MeleeAttackPower()*0.63999998569)
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				chainMultiplier *= 0.5
				war.Env.NextTargetUnit(aoeTarget)

				if result.Landed() && !hitLanded {
					hitLanded = true
				}
			}

			if hitLanded {
				if war.StanceMatches(warrior.DefensiveStance) {
					war.AddRage(sim, 20*war.GetRageMultiplier(target), rageMetrics)
				}
			}
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Revenge Reset Trigger",
		ActionID: actionID,
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeDodge | core.OutcomeParry,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			spell.CD.Reset()
		},
	})
}
