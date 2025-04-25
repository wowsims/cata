package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const (
	// This could be value or bitflag if we ended up needing multiple flags at the same time.
	//1 to 5 are used by MaelstromWeapon Stacks
	CastTagLightningOverload int32 = 6
)

// Shared precomputation logic for LB and CL.
func (shaman *Shaman) newElectricSpellConfig(actionID core.ActionID, baseCostPercent int32, baseCastTime time.Duration, isElementalOverload bool, bonusCoefficient float64) core.SpellConfig {
	mask := core.ProcMaskSpellDamage
	flags := SpellFlagElectric | SpellFlagFocusable
	if isElementalOverload {
		mask = core.ProcMaskSpellProc
		flags |= core.SpellFlagPassiveSpell
	} else {
		flags |= core.SpellFlagAPL
	}

	spell := core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     mask,
		Flags:        flags,
		MetricSplits: 6,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryInt32(isElementalOverload, 0, baseCostPercent),
			PercentModifier: 100 - (5 * shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: baseCastTime,
				GCD:      core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(shaman.MaelstromWeaponAura.GetStacks())
				castTime := shaman.ApplyCastSpeedForSpell(cast.CastTime, spell)
				if sim.CurrentTime+castTime > shaman.AutoAttacks.NextAttackAt() {
					shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
				}
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),
		BonusCoefficient: bonusCoefficient,
	}

	if isElementalOverload {
		spell.ActionID.Tag = CastTagLightningOverload
		spell.Cast.DefaultCast.CastTime = 0
		spell.Cast.DefaultCast.GCD = 0
		spell.Cast.DefaultCast.Cost = 0
		spell.Cast.ModifyCast = nil
		spell.MetricSplits = 0
		spell.DamageMultiplier *= 0.75
		spell.ThreatMultiplier = 0
	}

	return spell
}
