package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) getBarrageConfig() core.SpellConfig {
	return core.SpellConfig{
		SpellSchool:              core.SpellSchoolPhysical,
		ProcMask:                 core.ProcMaskRangedSpecial,
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
		sharedDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower()) * 0.2
		for _, aoeTarget := range sim.Encounter.TargetUnits {
			if aoeTarget == target {
				sharedDmg *= 2
			}
			result := spell.CalcDamage(sim, aoeTarget, sharedDmg, spell.OutcomeRangedHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealPeriodicDamage(sim, result)

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

func (hunter *Hunter) registerBarrageSpell() {
	if !hunter.Talents.Barrage {
		return
	}
	barrageTickSpell := hunter.getBarrageTickSpell()

	config := hunter.getBarrageConfig()
	config.ActionID = core.ActionID{SpellID: 120360}
	config.ProcMask = core.ProcMaskRangedSpecial
	config.Flags = core.SpellFlagChanneled | core.SpellFlagAPL | core.SpellFlagRanged
	config.ClassSpellMask = HunterSpellBarrage
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
		NumberOfTicks:       16,
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
		sharedDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower()) * 0.2
		return spell.CalcDamage(sim, target, sharedDmg, spell.OutcomeRangedHitAndCrit)
	}

	hunter.RegisterSpell(config)
}
