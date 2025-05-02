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
		Flags:    core.SpellFlagAPL,
		MinRange: 0,
		MaxRange: 40,
		FocusCost: core.FocusCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplierAdditive: 1,

		CritMultiplier:   hunter.CritMultiplier(1, 0),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: 131900},
				Label:    "AMurderOfCrowsDot",
				Tag:      "AMurderOfCrows",
			},

			NumberOfTicks: 30,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := hunter.GetBaseDamageFromCoeff(0.63) + 0.288*dot.Spell.RangedAttackPower(target)
				dot.Snapshot(target, baseDmg)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
				spell.DealOutcome(sim, result)
			})
		},
	})
}
