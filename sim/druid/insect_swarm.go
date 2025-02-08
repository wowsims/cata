package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) registerInsectSwarmSpell() {
	druid.InsectSwarm = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5570},
		SpellSchool:    core.SpellSchoolArcane | core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellInsectSwarm,
		Flags:          core.SpellFlagAPL | SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.08,
			Multiplier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   druid.BalanceCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Insect Swarm",
			},

			NumberOfTicks:       6,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.13,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := core.CalcScalingSpellAverageEffect(proto.Class_ClassDruid, 0.138)
				dot.Snapshot(target, baseDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}

			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			baseDamage := core.CalcScalingSpellAverageEffect(proto.Class_ClassDruid, 0.138)
			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
		},
	})
}
