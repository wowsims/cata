package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerPyroblastSpell() {
	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)

	var pyroblastDot *core.Spell
	/* implement when debuffs updated
	var CMProcChance float64
	if mage.Talents.CriticalMass > 0 {
		CMProcChance = float64(mage.Talents.CriticalMass) / 3.0
		//TODO double check how this works
		mage.CriticalMassAuras = mage.NewEnemyAuraArray(core.CriticalMassAura)
		mage.CritDebuffCategories = mage.GetEnemyExclusiveCategories(core.SpellCritEffectCategory)
		mage.Pyroblast.RelatedAuras = append(mage.Pyroblast.RelatedAuras, mage.CriticalMassAuras)
	} */

	pyroConfig := core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 11366},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | HotStreakSpells | core.SpellFlagAPL,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: core.TernaryFloat64(mage.HotStreakAura.IsActive(), 0, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 3500,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if mage.HotStreakAura.IsActive() {
					cast.CastTime = 0
					if !hasT8_4pc || sim.RandomFloat("MageT84PC") > T84PcProcChance {
						mage.HotStreakAura.Deactivate(sim)
					}
				}
			},
		},

		DamageMultiplier: 1,
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 1.545,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Pyroblast",
				Tag:   "FireMasteryDot",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 0.175 * mage.ScalingBaseDamage
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
			BonusCoefficient: 0.180,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.5 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					pyroblastDot.Dot(target).Apply(sim)
					//pyroblastDot.SpellMetrics[target.UnitIndex].Casts++
					/* The above line is used to count the dot in the cpm chart.
					This is misleading since it ends up doubling the overall pyroblast cpm,
					when most users probably just care how many times they press pyroblast.
					Example: 5 pyroblast casts in 1 minute end up showing as 10 cpm (5 pyro, 5 pyro dot)
					Should delete in my opinion.*/
				}
				spell.DealDamage(sim, result)
			})
		},
	}
	// Unsure about the implementation of the below, but just trusting it since it existed here
	mage.Pyroblast = mage.RegisterSpell(pyroConfig)

	dotConfig := pyroConfig
	dotConfig.ActionID = dotConfig.ActionID.WithTag(1)
	pyroblastDot = mage.RegisterSpell(dotConfig)
}
