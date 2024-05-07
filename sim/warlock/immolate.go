package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerImmolateSpell() {
	fireAndBrimstoneMod := warlock.AddDynamicMod(core.SpellModConfig{
		ClassMask:  WarlockSpellIncinerate | WarlockSpellChaosBolt,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.05 * float64(warlock.Talents.FireAndBrimstone),
	})

	warlock.ImmolateDot = warlock.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 348}.WithTag(1),
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		ClassSpellMask:   WarlockSpellImmolateDot,
		DamageMultiplier: 1,
		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate (DoT)",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					fireAndBrimstoneMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					fireAndBrimstoneMod.Deactivate()
				},
			},
			NumberOfTicks:       5,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.17599999905,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, warlock.CalcScalingSpellDmg(Coefficient_ImmolateDot))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 348},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellImmolate,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.08,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2000 * time.Millisecond,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.21999999881,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, warlock.CalcScalingSpellDmg(0.69199997187), spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				warlock.ImmolateDot.Cast(sim, target)
			}
			spell.DealDamage(sim, result)
		},
	})
}
