package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

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
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	var numTicks int32 = 3
	TickLengthBase := time.Millisecond * time.Duration(700-100*mage.Talents.MissileBarrage)
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
					// Force cast a missile right away
					mage.ArcaneMissilesTickSpell.Cast(sim, mage.CurrentTarget)
					mage.ArcaneMissilesProcAura.Deactivate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if mage.ArcaneMissilesProcAura.IsActive() {
						mage.ArcaneMissilesProcAura.Deactivate(sim)
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
			NumberOfTicks:        numTicks - 1, // subtracting 1 due to autocasting one OnGain
			TickLength:           TickLengthBase,
			HasteAffectsDuration: true,
			AffectedByCastSpeed:  true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.ArcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			//spell.DealOutcome(sim, result)
		},
	})
}
