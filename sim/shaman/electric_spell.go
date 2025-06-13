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

type ShamSpellConfig struct {
	ActionID            core.ActionID
	BaseCostPercent     float64
	BaseCastTime        time.Duration
	IsElementalOverload bool
	BonusCoefficient    float64
	BounceReduction     float64
	Coeff               float64
	Variance            float64
	SpellSchool         core.SpellSchool
	Overloads           *[2][]*core.Spell
}

// Shared precomputation logic for LB and CL.
// Needs isElementalOverload, actionID, baseCostPercent, baseCastTime, bonusCoefficient fields of the shamSpellConfig
func (shaman *Shaman) newElectricSpellConfig(config ShamSpellConfig) core.SpellConfig {
	mask := core.ProcMaskSpellDamage
	flags := SpellFlagShamanSpell | SpellFlagFocusable
	if config.IsElementalOverload {
		mask = core.ProcMaskSpellProc
		flags |= core.SpellFlagPassiveSpell
	} else {
		flags |= core.SpellFlagAPL
	}

	spell := core.SpellConfig{
		ActionID:     config.ActionID,
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     mask,
		Flags:        flags,
		MetricSplits: 6,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(config.IsElementalOverload, 0, config.BaseCostPercent),
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: config.BaseCastTime,
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
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: config.BonusCoefficient,
	}

	if config.IsElementalOverload {
		spell.ActionID.Tag = CastTagLightningOverload
		spell.ManaCost.BaseCostPercent = 0
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
