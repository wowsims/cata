package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterOverpowerSpell() {
	actionID := core.ActionID{SpellID: 7384}
	opAura := warrior.RegisterAura(core.Aura{
		ActionID: actionID,
		Tag:      EnableOverpowerTag,
		Label:    "Overpower Ready",
		Duration: time.Second * 5,
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:     "Overpower Trigger",
		ActionID: actionID,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeDodge,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !warrior.HasActiveAuraWithTagExcludingAura(EnableOverpowerTag, opAura) {
				opAura.Activate(sim)
			}
		},
	})

	warrior.Overpower = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7384},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskOverpower | SpellMaskSpecialAttack,
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
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.HasActiveAuraWithTag(EnableOverpowerTag) && warrior.StanceMatches(BattleStance)
		},

		DamageMultiplier: 1.25,
		CritMultiplier:   warrior.DefaultCritMultiplier(),
		ThreatMultiplier: 0.75,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			auras := warrior.GetAurasWithTag(EnableOverpowerTag)
			for _, aura := range auras {
				if aura.IsActive() {
					aura.Deactivate(sim)
				}
			}

			baseDamage := 0 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
