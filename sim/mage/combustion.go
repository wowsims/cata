package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerCombustionSpell() {
	if !mage.Talents.Combustion {
		return
	}

	mage.CombustionImpact = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty, // need to check proc mask for impact damage
		ClassSpellMask: MageSpellCombustionApplication,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultMageCritMultiplier(),
		BonusCoefficient:         1.113,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429 * mage.ClassSpellScaling
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.DealDamage(sim, result)
				mage.Combustion.Cast(sim, target)
			}
		},
	})

	var LBContribution float64
	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: MageSpellCombustion,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultMageCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Combustion Dot",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite, mage.PyroblastDot}
				combustionDotDamage := 0.0
				for _, spell := range dotSpells {
					dot := spell.Dot(target)
					if dot.IsActive() {
						if spell == mage.LivingBomb {
							// Living Bomb uses Critical Mass multiplicative in the combustion calculation
							LBContribution = (234 + 0.258*spell.SpellPower()) * 1.25 * (1 + 0.01*float64(mage.Talents.FirePower)) * (1.224 + 0.028*mage.GetMasteryPoints()) * (1 + 0.05*float64(mage.Talents.CriticalMass))
							LBContribution /= 3
							combustionDotDamage += LBContribution
						} else if spell == mage.Ignite {
							//Ignite's Contribution. Multiply by mastery again.
							combustionDotDamage += dot.SnapshotBaseDamage / 2 * (1.224 + 0.028*mage.GetMasteryPoints())
						} else if spell == mage.PyroblastDot {
							combustionDotDamage += dot.SnapshotBaseDamage * 1.25 * (1 + 0.01*float64(mage.Talents.FirePower)) * (1.224 + 0.028*mage.GetMasteryPoints())
						}
					}
				}
				dot.Snapshot(target, combustionDotDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.CombustionImpact,
		Type:  core.CooldownTypeDPS,
	})
}
