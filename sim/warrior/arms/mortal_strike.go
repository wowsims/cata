package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ArmsWarrior) RegisterMortalStrike() {
	war.mortalStrike = war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12294},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagMeleeMetrics,
		ClassSpellMask: warrior.SpellMaskMortalStrike | warrior.SpellMaskSpecialAttack,

		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Millisecond * 4500,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 423 + 0.8*(spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()))
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			} else {
				war.TriggerSlaughter(sim, target)
				if result.DidCrit() {
					war.TriggerWreckingCrew(sim)
				}
			}
		},
	})
}
