package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerArcaneMissilesSpell() {

	arcaneMissilesTickSpell := mage.GetOrRegisterSpell(core.SpellConfig{
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
			damage := 0.432 * mage.ClassSpellScaling
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	mage.RegisterSpell(core.SpellConfig{
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
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					// Make sure the arcane blast deactivation happens after last tick
					core.StartDelayedAction(sim, core.DelayedActionOptions{
						Priority: core.ActionPriorityDOT - 1,
						DoAt:     sim.CurrentTime,
						OnAction: func(sim *core.Simulation) {
							mage.ArcaneBlastAura.Deactivate(sim)
						},
					})
				},
			},
			NumberOfTicks:        3 - 1, // subtracting 1 due to force tick after apply
			TickLength:           time.Millisecond * 700,
			HasteAffectsDuration: true,
			AffectedByCastSpeed:  true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				arcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				spell.Dot(target).TickOnce(sim)
			}
			//spell.DealOutcome(sim, result)
		},
	})
}
