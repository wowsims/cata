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

	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty, // application burst might have a proc that might necessitate separating
		Flags:          SpellFlagMage,
		ClassSpellMask: MageSpellCombustion,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
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
				/* 				for _, spell := range dotSpells {
					dots := spell.Dot(mage.CurrentTarget)
					// EJ states that combustion double dips on mastery for LB + Pyro, but not ignite
					// https://web.archive.org/web/20120207223126/http://elitistjerks.com/f75/t110187-cataclysm_mage_simulators_formulators/p3/
					if spell != mage.Ignite {
						if dots != nil && dots.IsActive() {
							normalizedDPS = 1000000000 * spell.Dot(mage.CurrentTarget).SnapshotBaseDamage / float64(spell.Dot(mage.CurrentTarget).TickPeriod()) * (0.22 + 0.28*mage.GetMasteryPoints())
						}
					} else {
						if dots != nil && dots.IsActive() {
							normalizedDPS = 1000000000 * spell.Dot(mage.CurrentTarget).SnapshotBaseDamage / float64(spell.Dot(mage.CurrentTarget).TickPeriod())

						}
					}
					combustionDotDamage += normalizedDPS
				} */
				var spellDPS float64
				for _, spell := range dotSpells {
					dots := spell.Dot(mage.CurrentTarget)
					// EJ states that combustion double dips on mastery for LB + Pyro, but not ignite
					// https://web.archive.org/web/20120207223126/http://elitistjerks.com/f75/t110187-cataclysm_mage_simulators_formulators/p3/
					if spell != mage.Ignite {
						if dots != nil && dots.IsActive() {
							spellDPS = spell.Dot(mage.CurrentTarget).CalcSnapshotDamage(sim, mage.CurrentTarget, dots.OutcomeTick).Damage / 3
							fmt.Println(dots.Label, " snapped for : ", spellDPS)
						}
					} else {
						//This part is for ignite. The denominator will probably be variable if it works as intended in cata.
						if dots != nil && dots.IsActive() {
							spellDPS = spell.Dot(mage.CurrentTarget).SnapshotBaseDamage / 2
							fmt.Println(dots.Label, " snapped for : ", spellDPS)
						}
					}
					combustionDotDamage += spellDPS
				}
				dot.Snapshot(mage.CurrentTarget, combustionDotDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		BonusCoefficient:         1.113,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			spell.DealDamage(sim, result)
			spell.Dot(target).Apply(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.Combustion,
		Type:  core.CooldownTypeDPS,
	})

	mage.CombustionImpact = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlag(MageSpellFlagNone),
		ClassSpellMask: MageSpellCombustion,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Combustion Impact Fake Dot",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {

			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
		},
	})

}
