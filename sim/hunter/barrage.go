package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) getBarrageConfig() core.SpellConfig {
	return core.SpellConfig{
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskSpellProc,
		//ClassSpellMask:           PriestSpellMindSear,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           hunter.DefaultCritMultiplier(),
	}
}

func (hunter *Hunter) getBarrageTickSpell() *core.Spell {
	config := hunter.getBarrageConfig()
	config.ActionID = core.ActionID{SpellID: 120361}
	config.MissileSpeed = 30
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		sharedDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) * 0.4
		for _, aoeTarget := range sim.Encounter.TargetUnits {

			// Calc spell damage but deal as periodic for metric purposes
			result := spell.CalcDamage(sim, aoeTarget, sharedDmg, spell.OutcomeRangedHitAndCrit)
			// Will see if this works in dot effect
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealPeriodicDamage(sim, result)

				// Adjust metrics just for Mind Sear as it is a edgecase and needs to be handled manually
				if result.DidCrit() {
					spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
				} else {
					spell.SpellMetrics[result.Target.UnitIndex].Ticks++
				}
			})
		}

		spell.SpellMetrics[target.UnitIndex].Casts--
	}
	return hunter.RegisterSpell(config)
}

func (hunter *Hunter) registerBarrageSpell() *core.Spell {
	barrageTickSpell := hunter.getBarrageTickSpell()

	config := hunter.getBarrageConfig()
	config.ActionID = core.ActionID{SpellID: 120360}
	config.Flags = core.SpellFlagChanneled | core.SpellFlagAPL
	config.FocusCost = core.FocusCostOptions{
		Cost: 30,
	}
	config.MaxRange = 40
	config.Cast = core.CastConfig{
		DefaultCast: core.Cast{
			GCD: core.GCDDefault,
		},
		CD: core.Cooldown{
			Timer:    hunter.NewTimer(),
			Duration: 30 * time.Second,
		},
	}

	config.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: "Barrage-" + hunter.Label,
		},
		NumberOfTicks:       15,
		TickLength:          time.Millisecond * 200,
		AffectedByCastSpeed: true,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			barrageTickSpell.Cast(sim, target)
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)
		if result.Landed() {
			spell.Dot(target).Apply(sim)
		}
	}
	config.ExpectedTickDamage = func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
		sharedDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) * 0.4
		return spell.CalcDamage(sim, target, sharedDmg, spell.OutcomeRangedHitAndCrit)
	}

	return hunter.RegisterSpell(config)
}
