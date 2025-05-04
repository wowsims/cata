package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) registerRevealingStrike() {
	multiplier := 1.35
	actionID := core.ActionID{SpellID: 84617}
	cpMetric := comRogue.NewComboPointMetrics(core.ActionID{SpellID: 1752}) // Sinister Strike spell ID

	wepDamage := 1.6

	// Enemy Debuff Aura for Finisher Damage
	rvsAura := comRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Revealing Strike",
			ActionID: actionID,
			Duration: 15 * time.Second,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				core.EnableDamageDoneByCaster(DDBC_RevealingStrike, DDBC_Total, comRogue.AttackTables[aura.Unit.UnitIndex], func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
					if spell.Matches(rogue.RogueSpellDamagingFinisher) {
						return multiplier
					}
					return 1.0
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				core.DisableDamageDoneByCaster(DDBC_RevealingStrike, comRogue.AttackTables[aura.Unit.UnitIndex])
				aura.Deactivate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.ClassSpellMask == rogue.RogueSpellSinisterStrike {
					if sim.Proc(0.2, "Revealing Strike Extra Combo Point") {
						comRogue.AddComboPoints(sim, 1, cpMetric)
					}
				}
			},
		})
	})

	// Attack
	comRogue.RevealingStrike = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL | rogue.SpellFlagBuilder,
		ClassSpellMask: rogue.RogueSpellRevealingStrike,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 800,
			},
		},
		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: 0.8,
		},

		DamageMultiplier: wepDamage,
		CritMultiplier:   comRogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			comRogue.BreakStealth(sim)

			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if result.Landed() {
				comRogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				rvsAura.Get(target).Activate(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
