package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

/* 1 spell bar, no channeling background
func (mage *Mage) registerArcaneMissilesSpell() {

	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)
	var numTicks int32 = 3

	mage.ArcaneMissiles = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5143},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:          SpellFlagMage | core.SpellFlagAPL | core.SpellFlagChanneled,
		ClassSpellMask: MageSpellArcaneMissilesCast,
		MissileSpeed:   20,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCritRating:  core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1,
		DamageMultiplierAdditive: 1 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.278,

		Dot: core.DotConfig{
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: 7268},
				Label:    "Arcane Missile Casting",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if mage.ArcaneMissilesProcAura.IsActive() {
						if !hasT8_4pc || sim.Proc(T84PcProcChance, "MageT84PC") {
							mage.ArcaneMissilesProcAura.Deactivate(sim)
						}
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					// Ensures tick occurs before aura ends after the last tick
					dot := mage.ArcaneMissiles.Dot(aura.Unit)
					if dot.TickCount < dot.NumberOfTicks {
						dot.TickCount++
						dot.TickOnce(sim)
					}
					mage.ArcaneBlastAura.Deactivate(sim)
				},
			},
			NumberOfTicks:       numTicks - 1, // -1 due to forcing first tick to occur on Aura gain
			TickLength:          time.Millisecond * 700,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				mage.ArcaneBlastAura.Deactivate(sim)
				baseDamage := 0.432 * mage.ScalingBaseDamage
				dot.Snapshot(target, baseDamage)
				//dot.TickOnce(sim)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitAndCrit)
			dot.Apply(sim)
			dot.TickOnce(sim)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealOutcome(sim, result)
			})
		},
	})

}
*/

func (mage *Mage) registerArcaneMissilesSpell() {

	mage.ArcaneMissilesTickSpell = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7268},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:          SpellFlagMage,
		ClassSpellMask: MageSpellArcaneMissilesTick,
		MissileSpeed:   20,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.278,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 0.432 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			//spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			spell.DealDamage(sim, result)
			//})
		},
	})

	var numTicks int32 = 3
	mage.ArcaneMissiles = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7268},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagMage | core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneMissilesCast,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.ArcaneMissilesProcAura.IsActive()
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ArcaneMissiles",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					// Currently it doesn't cast the first missile until the time of the first tick
					mage.ArcaneMissilesTickSpell.Cast(sim, mage.CurrentTarget)
					mage.ArcaneMissilesProcAura.Deactivate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if mage.ArcaneMissilesProcAura.IsActive() {
						//mage.ArcaneMissilesProcAura.Deactivate(sim)
					}

					// TODO: This check is necessary to ensure the final tick occurs before
					// Arcane Blast stacks are dropped. To fix this, ticks need to reliably
					// occur before aura expirations.
					dot := mage.ArcaneMissiles.Dot(aura.Unit)
					if dot.TickCount < dot.NumberOfTicks {
						dot.TickCount++
						dot.TickOnce(sim)
					}
					mage.ArcaneBlastAura.Deactivate(sim)
				},
			},
			NumberOfTicks:       numTicks - 1, // subtracting 1 due to autocasting one OnGain
			TickLength:          time.Millisecond * 700,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.ArcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
