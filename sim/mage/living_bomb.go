package mage

import (
	"sort"
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerLivingBombSpell() {
	// Cata version has a cap of 3 active dots at once
	// activeLivingBombs should only ever be 3 LBs long
	// When a dot is trying to be applied,
	// 1) it should remove the dot with the longest remaining duration
	// 2) to do this, when a dot is applied, it checks the length of the array
	//   2a) if the array is longer than 3, remove the last element
	// 3) append the dot to the array
	// 4) sort the array by remaining duration, such that the longest remaining duration is LAST, to fit step 1
	// When a dot expires, remove the 1st element.
	var activeLivingBombs []*core.Dot
	const maxLivingBombs int = 3

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44461},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellLivingBombExplosion,

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		BonusCoefficient:         0.516,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.5 * mage.ClassSpellScaling
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

		},
	})

	mage.LivingBomb = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44457},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellLivingBombDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.17,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return len(activeLivingBombs) < maxLivingBombs
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					//TODO make it not cast the explosion if expiration was due to target cap
					livingBombExplosionSpell.Cast(sim, aura.Unit)
					if len(activeLivingBombs) != 0 {
						activeLivingBombs = activeLivingBombs[1:]
					}
				},
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.258,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 0.25*mage.ClassSpellScaling)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			spell.DealOutcome(sim, result)

			if result.Landed() {
				if len(activeLivingBombs) >= maxLivingBombs {
					activeLivingBombs[len(activeLivingBombs)-1].Deactivate(sim)
					if len(activeLivingBombs) != 0 {
						activeLivingBombs = activeLivingBombs[:1]
					}
				}
				spell.Dot(target).Apply(sim)
				activeLivingBombs = append(activeLivingBombs, mage.LivingBomb.Dot(mage.CurrentTarget))
				sort.Slice(activeLivingBombs, func(i, j int) bool {
					return activeLivingBombs[i].Duration < activeLivingBombs[j].Duration
				})
			}
		},
	})

	mage.LivingBombImpact = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44457},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellLivingBombDot,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.17,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					livingBombExplosionSpell.Cast(sim, aura.Unit)
				},
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)

		},
	})
}
