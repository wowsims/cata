package retribution

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

/*
An area attack that consumes 3 charges of Holy Power to cause 100% weapon damage as Holy damage to all enemies within 8 yards.

-- Glyph of Divine Storm --
Using Divine Storm will also heal you for 5% of your maximum health.
-- /Glyph of Divine Storm --
*/
func (ret *RetributionPaladin) registerDivineStorm() {
	numTargets := ret.Env.GetNumTargets()
	actionID := core.ActionID{SpellID: 53385}

	ret.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagAoE,
		ClassSpellMask: paladin.SpellMaskDivineStorm,

		MaxRange: 8,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return ret.DivineCrusaderAura.IsActive() || ret.HolyPower.CanSpend(3)
		},

		DamageMultiplier: 1,
		CritMultiplier:   ret.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)

			for idx := range numTargets {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				baseDamage := ret.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}

			if !ret.DivineCrusaderAura.IsActive() {
				ret.HolyPower.Spend(sim, 3, actionID)
			}

			for idx := range numTargets {
				spell.DealDamage(sim, results[idx])
			}
		},
	})
}
