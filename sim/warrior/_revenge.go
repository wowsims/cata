package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterRevengeSpell() {
	actionID := core.ActionID{SpellID: 6572}

	revengeReadyAura := warrior.RegisterAura(core.Aura{
		Label:    "Revenge Ready",
		Duration: 5 * time.Second,
		ActionID: actionID.WithTag(1),
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:     "Overpower Trigger",
		ActionID: actionID,
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			revengeReadyAura.Activate(sim)
		},
	})

	extraHit := warrior.Talents.ImprovedRevenge > 0 && warrior.Env.GetNumTargets() > 1
	extraHitMult := 0.5 * float64(warrior.Talents.ImprovedRevenge)

	warrior.Revenge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskRevenge | SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   5,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 5,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance) && revengeReadyAura.IsActive()
		},

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1,
		FlatThreatBonus:  121,
		CritMultiplier:   warrior.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ap := spell.MeleeAttackPower() * 0.31
			baseDamage := sim.Roll(1618.3, 1977.92) + ap
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			if extraHit {
				otherTarget := sim.Environment.NextTargetUnit(target)
				// TODO: Reimplement using scaling coefficients and variance once those stats are available
				baseDamage := sim.Roll(1618.3, 1977.92) + ap
				spell.CalcAndDealDamage(sim, otherTarget, baseDamage*extraHitMult, spell.OutcomeMeleeSpecialHitAndCrit)
			}

			revengeReadyAura.Deactivate(sim)
		},
	})
}
