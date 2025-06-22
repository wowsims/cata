package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const drainSoulScale = 0.257 * 1.5 // 2025.06.13 Changes to Beta - Drain Soul Damage increased by 50%
const drainSoulCoeff = 0.257 * 1.5

func (affliction *AfflictionWarlock) registerDrainSoul() {
	affliction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1120},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagChanneled,
		ClassSpellMask: warlock.WarlockSpellDrainSoul,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 1.5},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           affliction.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura:                 core.Aura{Label: "DrainSoul"},
			NumberOfTicks:        6,
			TickLength:           2 * time.Second,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			BonusCoefficient:     drainSoulCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, affliction.CalcScalingSpellDmg(drainSoulScale))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

				// Every 2nd tick grants 1 soul shard
				if dot.TickCount()%2 == 0 {
					affliction.SoulShards.Gain(sim, 1, dot.Spell.ActionID)
				}

				if !result.Landed() || !sim.IsExecutePhase20() {
					return
				}

				// 2025.06.13 Changes to Beta - Drain Soul DoT damage increased to 100%
				affliction.ProcMaleficEffect(target, 1, sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	dmgMode := affliction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1,
		ClassMask:  warlock.WarlockSpellDrainSoul,
	})

	affliction.RegisterResetEffect(func(s *core.Simulation) {
		dmgMode.Deactivate()
		s.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute > 20 {
				return
			}

			dmgMode.Activate()
		})
	})
}
