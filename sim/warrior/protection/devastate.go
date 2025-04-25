package protection

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) RegisterDevastate() {
	if !war.Talents.Devastate {
		return
	}

	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20243},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskDevastate | warrior.SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		FlatThreatBonus:  315,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: verify stacking behavior (it probably hasn't changed)
			// Bonus 19.333% weapon damage / stack of sunder. Counts stacks AFTER cast but only if stacks > 0.
			weaponDamageMult := 1.09
			saStacks := war.SunderArmorAuras.Get(target).GetStacks()
			if saStacks != 0 {
				weaponDamageMult += []float64{0.0, 0.19333, 0.38666, 0.58}[min(saStacks+1, 3)]
			}

			baseDamage := (weaponDamageMult * spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()))

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			result.Threat = spell.ThreatFromDamage(result.Outcome, result.Damage+0.05*spell.MeleeAttackPower())
			spell.DealDamage(sim, result)

			if result.Landed() {
				war.TryApplySunderArmorEffect(sim, target)
			} else {
				spell.IssueRefund(sim)
			}
		},

		RelatedAuraArrays: war.SunderArmorAuras.ToMap(),
	})
}
