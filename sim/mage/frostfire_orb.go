package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerFrostfireOrbSpell() {
	// For future documenting...
	// Talent point 1 = benefit from frost specialization
	// Talent point 2 = slows more, so not important here
	// Need logs/spell IDs

	mage.FrostfireOrbTickSpell = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 82739},
		SpellSchool: core.SpellSchoolFire,
		// no idea yet what it procs, likely nothing
		ProcMask:       core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:          SpellFlagMage | core.SpellFlagNoLogs,
		ClassSpellMask: MageSpellFrostfireOrb,
		MissileSpeed:   20,

		DamageMultiplier: 1 + .01*float64(mage.Talents.TormentTheWeak),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Seems like the most straightforward way to split damage onto different units for results page
			damage := (0.278*mage.ScalingBaseDamage + 0.134*spell.SpellPower()) / float64(len(sim.Encounter.TargetUnits))
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	mage.FrostfireOrb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 84714}, // Assumed
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Frostfire Orb",
				Duration: time.Second * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for i := 0; i < 15; i++ {
				mage.FrostfireOrbTickSpell.Cast(sim, target)
			}
		},
	})
}
