package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerMoonfireSpell() {
	druid.registerMoonfireImpactSpell()
	druid.registerMoonfireDoTSpell()
}

func (druid *Druid) registerMoonfireDoTSpell() {
	druid.Moonfire.RelatedDotSpell = druid.Unit.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8921}.WithTag(1),
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellMoonfireDoT,
		Flags:          SpellFlagOmenTrigger | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Moonfire",
			},
			NumberOfTicks:       6,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.18,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := core.CalcScalingSpellAverageEffect(proto.Class_ClassDruid, 0.095)
				dot.Snapshot(target, baseDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)

			spell.Dot(target).Apply(sim)
			spell.DealOutcome(sim, result)
		},
	})
}

func (druid *Druid) registerMoonfireImpactSpell() {
	druid.SetSpellEclipseEnergy(DruidSpellMoonfire, MoonfireBaseEnergyGain, MoonfireBaseEnergyGain)

	druid.Moonfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8921},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellMoonfire,
		Flags:          core.SpellFlagAPL | SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 9,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,

		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.18,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			min, max := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, 0.221, 0.2)
			baseDamage := sim.Roll(min, max)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				if druid.Sunfire.Dot(target).IsActive() {
					druid.Sunfire.Dot(target).Deactivate(sim)
				}

				druid.ExtendingMoonfireStacks = 3
				druid.Moonfire.RelatedDotSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
