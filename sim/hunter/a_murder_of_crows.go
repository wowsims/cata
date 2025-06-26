package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerAMOCSpell() {
	if !hunter.Talents.AMurderOfCrows {
		return
	}

	// Add a spell modifier that reduces cooldown by 50% during execute phase
	executePhaseMod := hunter.AddDynamicMod(core.SpellModConfig{
		ClassMask: HunterSpellAMurderOfCrows,
		TimeValue: -60 * time.Second,
		Kind:      core.SpellMod_Cooldown_Flat,
	})

	hunter.RegisterResetEffect(func(sim *core.Simulation) {
		executePhaseMod.Deactivate()
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, executePhase int32) {
			if executePhase == 20 {
				executePhaseMod.Activate()
			}
		})
	})

	hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 131894},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskProc,
		ClassSpellMask: HunterSpellAMurderOfCrows,
		Flags:          core.SpellFlagAPL | core.SpellFlagApplyArmorReduction | core.SpellFlagRanged,
		MinRange:       0,
		MaxRange:       40,
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
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalHitAndCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)

			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + (time.Second * 2),
				OnAction: func(sim *core.Simulation) {
					if result.Landed() {
						spell.Dot(target).Apply(sim)
					}
				},
			})

		},
	})
}
