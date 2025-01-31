package combat

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func (comRogue *CombatRogue) registerRevealingStrike() {
	if !comRogue.Talents.RevealingStrike {
		return
	}

	hasGlyph := comRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfRevealingStrike)
	multiplier := core.TernaryFloat64(hasGlyph, .45, .35)
	actionID := core.ActionID{SpellID: 84617}
	isApplied := false

	damageMultiMod := comRogue.AddDynamicMod(core.SpellModConfig{
		ClassMask:  rogue.RogueSpellEviscerate | rogue.RogueSpellEnvenom | rogue.RogueSpellRupture,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: multiplier,
	})

	// Enemy Debuff Aura for Finisher Damage
	rvsAura := comRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Revealing Strike",
			ActionID: actionID,
			Duration: 15 * time.Second,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if !isApplied {
					damageMultiMod.Activate()
					isApplied = true
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				damageMultiMod.Deactivate()
				isApplied = false
				aura.Deactivate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.Flags.Matches(rogue.SpellFlagFinisher) {
					aura.Deactivate(sim)
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
				GCD: time.Second,
			},
		},
		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: 0.8,
		},

		DamageMultiplier: 1.29,
		CritMultiplier:   comRogue.MeleeCritMultiplier(false),
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
