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
	multiplier := 1 + core.TernaryFloat64(hasGlyph, .45, .35)
	actionID := core.ActionID{SpellID: 84617}

	// Enemy Debuff Aura for Finisher Damage
	rvsAura := comRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Revealing Strike",
			ActionID: actionID,
			Duration: 15 * time.Second,

			// Technically this _could_ cause problems in a target swapping situation, but it's good enough.
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				comRogue.Eviscerate.DamageMultiplier *= multiplier
				comRogue.Envenom.DamageMultiplier *= multiplier
				comRogue.Rupture.DamageMultiplier *= multiplier
			},

			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.Flags.Matches(rogue.SpellFlagFinisher) {
					comRogue.Eviscerate.DamageMultiplier /= multiplier
					comRogue.Envenom.DamageMultiplier /= multiplier
					comRogue.Rupture.DamageMultiplier /= multiplier
					aura.Deactivate(sim)
				}
			},
		})
	})

	// Attack
	comRogue.RevealingStrike = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL | rogue.SpellFlagBuilder,
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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			comRogue.BreakStealth(sim)

			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if result.Landed() {
				comRogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				aura := rvsAura.Get(target)
				aura.Activate(sim)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
