package arms

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) RegisterMortalStrike() {
	weaponDamageConfig := warrior.SpellEffectWeaponDmgPctConfig{
		BaseWeapon_Pct:    0.8,
		Coefficient:       0.37599998713,
		EffectPerLevel:    1,
		BaseSpellLevel:    10,
		MaxSpellLevel:     80,
		ClassSpellScaling: war.ClassSpellScaling,
	}

	war.mortalStrike = war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12294},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: warrior.SpellMaskMortalStrike | warrior.SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

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

		DamageMultiplier: weaponDamageConfig.CalcSpellDamagePct(),
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := weaponDamageConfig.CalcAddedSpellDamage() +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
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
