package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerAMOCSpell() {
	if !hunter.Talents.AMurderOfCrows {
		return
	}

	hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 131894},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskProc,
		//ClassSpellMask: HunterSpellSerpentSting,
		Flags:    core.SpellFlagAPL | core.SpellFlagApplyArmorReduction,
		MinRange: 0,
		MaxRange: 40,
		FocusCost: core.FocusCostOptions{
			Cost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		DamageMultiplierAdditive: 1,

		CritMultiplier:   hunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: 131900},
				Label:    "Peck",
				Tag:      "Peck",
			},

			NumberOfTicks: 30,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDmg := hunter.GetBaseDamageFromCoeff(0.63) + (0.288 * dot.Spell.RangedAttackPower())
				dot.Snapshot(target, baseDmg)
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + (time.Second * 2),
				OnAction: func(sim *core.Simulation) {
					if result.Landed() {
						spell.Dot(target).Apply(sim)
					}
					if sim.IsExecutePhase20() {
						spell.CD.Duration = time.Second * 30
					}
					spell.DealOutcome(sim, result)
				},
			})

		},
	})
}
