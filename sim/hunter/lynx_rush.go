package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) RegisterLynxRushSpell() {
	if !hunter.Talents.LynxRush || hunter.Pet == nil {
		return
	}
	hunter.Pet.lynxRushSpell = hunter.Pet.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 120697},
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		ClassSpellMask: HunterSpellLynxRush,
		ProcMask:       core.ProcMaskProc,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagMeleeMetrics,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Lynx Rush",
				MaxStacks: 9,
				Duration:  time.Second * 15,
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				if isRollover {
					panic("Lynx Rush cannot roll over snapshots!")
				}

				dot.SnapshotPhysical(target, (0.038*dot.Spell.MeleeAttackPower()+hunter.GetBaseDamageFromCoeff(0.06899999827))*float64(dot.Aura.GetStacks()))
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		CritMultiplier:   1,
		DamageMultiplier: 1,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 0,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			targetCount := hunter.Env.GetNumTargets()
			idx := int32(sim.RollWithLabel(0, float64(targetCount), "LynxRush"))
			if idx < 0 {
				idx = 0
			} else if idx >= targetCount {
				idx = targetCount - 1
			}
			target := sim.Environment.AllUnits[idx]

			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				dot := spell.Dot(target)
				if dot.IsActive() {
					dot.Refresh(sim)
					dot.AddStack(sim)
					dot.TakeSnapshot(sim, false)
				} else {
					dot.Apply(sim)
					dot.SetStacks(sim, 1)
					dot.TakeSnapshot(sim, false)
				}
			}
		},
	})
	hunter.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 120697},
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		ProcMask:    core.ProcMaskEmpty,
		SpellSchool: core.SpellSchoolNone,
		Flags:       core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 90,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			if hunter.Pet != nil && hunter.Pet.lynxRushSpell != nil {
				hunter.Pet.AutoAttacks.DelayMeleeBy(sim, time.Second*4)
				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					NumTicks: 9,
					Period:   time.Millisecond * 450,
					OnAction: func(sim *core.Simulation) {
						hunter.Pet.lynxRushSpell.Cast(sim, nil)
					},
				})
			}
		},
	})
}
