package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (sinRogue *AssassinationRogue) registerDispatch() {
	addedDamage := sinRogue.GetBaseDamageFromCoefficient(0.62900000811)
	weaponPercent := 6.45

	sinRogue.Dispatch = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 111240},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellDispatch,

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 700,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sinRogue.HasDagger(core.MainHand) && (sim.IsExecutePhase35() || sinRogue.HasActiveAura("Blindside"))
		},

		DamageMultiplier:         weaponPercent,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           sinRogue.CritMultiplier(true),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := addedDamage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			outcome := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)
			if outcome.Landed() {
				sinRogue.AddComboPointsOrAnticipation(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
