package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/paladin"
)

/*
	Tooltip:

Sends bolts of power in all directions, causing ${($m1+$M1)/2+$SPH*0.91} Holy damage

-- Glyph of Focused Wrath --

	divided among all enemies within 10 yards

-- else --

	to your target

--

, stunning Demons

-- Glyph of Holy Wrath --

	, Aberrations, Dragonkin, Elementals

--

	and Undead for 3 sec.

-- Glyph of Final Wrath --

	Causes 50% additional damage to targets with less than 20% health.

--
*/
func (prot *ProtectionPaladin) registerHolyWrath() {
	scalingCoef := 7.53200006485
	baseDamage := prot.CalcScalingSpellDmg(scalingCoef / 2) // The scaling coef is divided by 2 in the game (and sort of matches the tooltip)
	// variance := 0.1099999994 // unused???
	apCoef := 0.91 // This coef is only in the tooltip, where it says it's an SP coef but that's wrong

	hasFinalWrath := prot.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFinalWrath)
	hasFocusedWrath := prot.HasMinorGlyph(proto.PaladinMinorGlyph_GlyphOfFocusedWrath)

	var numTargets int32
	if hasFocusedWrath {
		numTargets = 1
	} else {
		numTargets = prot.Env.GetNumTargets()
	}

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
			damage := baseDamage + apCoef*spell.MeleeAttackPower()

			// Damage is split between all mobs, each hit rolls for hit/crit separately
			damage /= float64(numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				multiplier := spell.DamageMultiplier
				if hasFinalWrath && target.CurrentHealthPercent() < 0.2 {
					spell.DamageMultiplier *= 1.5
				}

				results[idx] = spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

				spell.DamageMultiplier = multiplier

				target = sim.Environment.NextTargetUnit(target)
			}

			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				for idx := int32(0); idx < numTargets; idx++ {
					spell.DealDamage(sim, results[idx])
				}
			})
		},
	})
}
