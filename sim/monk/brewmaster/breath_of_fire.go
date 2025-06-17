package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerBreathOfFire() {
	actionID := core.ActionID{SpellID: 115181}
	dotActionID := core.ActionID{SpellID: 123725}
	chiMetrics := bm.NewChiMetrics(actionID)

	bm.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | monk.SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellBreathOfFire,
		MaxRange:       8,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    bm.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   bm.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Breath Of Fire" + bm.Label,
				ActionID: dotActionID,
			},
			NumberOfTicks:       4,
			TickLength:          time.Millisecond * 2000,
			AffectedByCastSpeed: false,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := bm.CalcAndRollDamageRange(sim, 0.475, 0.242) + 0.1626*dot.Spell.MeleeAttackPower()
				dot.Snapshot(target, baseDamage)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bm.StanceMatches(monk.SturdyOx) && bm.GetChi() >= 2
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, enemyTarget := range sim.Encounter.TargetUnits {
				baseDamage := bm.CalcAndRollDamageRange(sim, 1.475, 0.242) + 0.3626*spell.MeleeAttackPower()
				result := spell.CalcOutcome(sim, enemyTarget, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter)

				if result.Landed() {
					spell.CalcAndDealDamage(sim, enemyTarget, baseDamage, spell.OutcomeMagicCrit)

					if bm.DizzyingHazeAuras.Get(enemyTarget).IsActive() {
						spell.Dot(enemyTarget).Apply(sim)
					}
				}
			}

			bm.SpendChi(sim, 2, chiMetrics)
		},
	})
}
