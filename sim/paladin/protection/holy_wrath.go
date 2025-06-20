package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/paladin"
)

/*
	Tooltip:

Sends bolts of power in all directions, causing ((8127 + 9075) / 2) / 2 + <AP> * 0.91 Holy damage

-- Glyph of Focused Wrath --
divided among all enemies within 10 yards
-- else --
to your target
----------

, stunning Demons

-- Glyph of Holy Wrath --
, Aberrations, Dragonkin, Elementals
-- /Glyph of Holy Wrath --

and Undead for 3 sec.

-- Glyph of Final Wrath --
Causes 50% additional damage to targets with less than 20% health.
-- /Glyph of Final Wrath --
*/
func (prot *ProtectionPaladin) registerHolyWrath() {
	hasGlyphOfFinalWrath := prot.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFinalWrath)
	hasGlyphOfFocusedWrath := prot.HasMinorGlyph(proto.PaladinMinorGlyph_GlyphOfFocusedWrath)

	numTargets := core.TernaryInt32(hasGlyphOfFocusedWrath, 1, prot.Env.GetNumTargets())

	prot.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 119072},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskHolyWrath,

		MissileSpeed: 40,
		MaxRange:     10,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    prot.NewTimer(),
				Duration: 9 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   prot.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)

			// Ingame tooltip is ((<MIN> + <MAX>) / 2) / 2
			// This is the same as, <AVG> / 2 which is the same as just halving the coef
			baseDamage := prot.CalcScalingSpellDmg(7.53200006485/2) + 0.91*spell.MeleeAttackPower()

			// Damage is split between all mobs, each hit rolls for hit/crit separately
			baseDamage /= float64(numTargets)

			for idx := range numTargets {
				multiplier := spell.DamageMultiplier
				if hasGlyphOfFinalWrath && sim.IsExecutePhase20() {
					spell.DamageMultiplier *= 1.5
				}

				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

				spell.DamageMultiplier = multiplier

				target = sim.Environment.NextTargetUnit(target)
			}

			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				for idx := range numTargets {
					spell.DealDamage(sim, results[idx])
				}
			})
		},
	})
}
