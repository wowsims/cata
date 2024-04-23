package fire

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/mage"
)

func (Mage *FireMage) registerPyroblastSpell() {

	/* implement when debuffs updated
	var CMProcChance float64
	if mage.Talents.CriticalMass > 0 {
		CMProcChance = float64(mage.Talents.CriticalMass) / 3.0
		//TODO double check how this works
		mage.CriticalMassAuras = mage.NewEnemyAuraArray(core.CriticalMassAura)
		mage.CritDebuffCategories = mage.GetEnemyExclusiveCategories(core.SpellCritEffectCategory)
		mage.Pyroblast.RelatedAuras = append(mage.Pyroblast.RelatedAuras, mage.CriticalMassAuras)
	} */

	Mage.Pyroblast = Mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11366},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          mage.SpellFlagMage | mage.HotStreakSpells | core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellPyroblast,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: core.TernaryFloat64(Mage.HotStreakAura.IsActive(), 0, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 3500,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if Mage.HotStreakAura.IsActive() {
					cast.CastTime = 0
					Mage.HotStreakAura.Deactivate(sim)
				}
			},
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           Mage.DefaultSpellCritMultiplier(),
		BonusCoefficient:         1.545,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.5 * Mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					Mage.PyroblastDot.Cast(sim, target)
					spell.DealDamage(sim, result)
					//pyroblastDot.SpellMetrics[target.UnitIndex].Hits++
					//pyroblastDot.SpellMetrics[target.UnitIndex].Casts = 0
					/* The 2 above metric changes should show how many ticks land
					without affecting the overall pyroblast cast metric
					*/
				}
			})
		},
	})

	Mage.PyroblastDot = Mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11366}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          mage.SpellFlagMage,
		ClassSpellMask: mage.MageSpellPyroblastDot,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   Mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "PyroblastDoT",
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 3,
			BonusCoefficient:    0.180,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 0.175*Mage.ScalingBaseDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
			Mage.PyroblastDot.SpellMetrics[target.UnitIndex].Casts = 0
		},
	})

	Mage.PyroblastDotImpact = Mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11366}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          mage.SpellFlagMage,
		ClassSpellMask: mage.MageSpellPyroblastDot,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   Mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "PyroblastDoT Fake Impact",
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				//
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
			Mage.PyroblastDot.SpellMetrics[target.UnitIndex].Casts = 0
		},
	})
}
