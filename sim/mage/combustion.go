package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerCombustionSpell() {
	if !mage.Talents.Combustion {
		return
	}

	var combustionDotDamage float64
	mage.CombustionImpact = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty, // need to check proc mask for impact damage
		ClassSpellMask: MageSpellCombustionApplication,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
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
		ActionID:       core.ActionID{SpellID: 83853},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: MageSpellCombustion,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagIgnoreAttackerModifiers,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

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
				combustionDotDamage = 0.0

				dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite, mage.PyroblastDot}

				var spellDPS float64
				for _, spell := range dotSpells {
					dots := spell.Dot(target)
					if spell == mage.LivingBomb {
						if dots != nil && dots.IsActive() {
							// Living Bomb uses Critical Mass multiplicative in the combustion calculation
							LBContribution = (234 + 0.258*spell.SpellPower()) * 1.25 * (1 + 0.01*float64(mage.Talents.FirePower)) * (1.224 + 0.028*mage.GetMasteryPoints()) * (1 + 0.05*float64(mage.Talents.CriticalMass))
							LBContribution /= 3
							spellDPS = LBContribution
						} else {
							spellDPS = 0
						}
					} else if spell == mage.Ignite {
						//Ignite's Contribution. Multiply by mastery again.
						if dots != nil && dots.IsActive() {
							spellDPS = spell.Dot(target).SnapshotBaseDamage / 2 * (1.224 + 0.028*mage.GetMasteryPoints())
						} else {
							spellDPS = 0
						}
					} else if spell == mage.PyroblastDot {
						if dots != nil && dots.IsActive() {
							spellDPS = dots.SnapshotBaseDamage * 1.25 * (1 + 0.01*float64(mage.Talents.FirePower)) * (1.224 + 0.028*mage.GetMasteryPoints())
						} else {
							spellDPS = 0
						}
					}
					fmt.Println("Adding ", spellDPS, " damage to combustion from ", dots.Label)
					combustionDotDamage += spellDPS
				}
				dot.Snapshot(target, combustionDotDamage)
				dot.SnapshotAttackerMultiplier = 1
				fmt.Println("Combustion SnapshotBaseDamage: ", combustionDotDamage)

			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				fmt.Println("Tick Damage: ", result.Damage)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mage.Combustion.Dot(target).Apply(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.CombustionImpact,
		Type:  core.CooldownTypeDPS,
	})
}
