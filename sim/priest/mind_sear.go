package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) getMindSearBaseConfig() core.SpellConfig {
	return core.SpellConfig{
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskSpellProc,
		ClassSpellMask:           PriestSpellMindSear,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		BonusCoefficient:         0.2622,
	}
}

func (priest *Priest) getMindSearTickSpell() *core.Spell {
	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: 48045}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		damage := priest.ClassSpellScaling * 0.23
		for _, aoeTarget := range sim.Encounter.TargetUnits {

			// Calc spell damage but deal as periodic for metric purposes
			result := spell.CalcDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCritNoHitCounter)
			spell.DealPeriodicDamage(sim, result)

			// Adjust metrics just for Mind Sear as it is a edgecase and needs to be handled manually
			if result.DidCrit() {
				spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
			} else {
				spell.SpellMetrics[result.Target.UnitIndex].Ticks++
			}
		}

		spell.SpellMetrics[target.UnitIndex].Casts--
	}
	return priest.RegisterSpell(config)
}

func (priest *Priest) newMindSearSpell() *core.Spell {
	mindSearTickSpell := priest.getMindSearTickSpell()

	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: 48045}
	config.Flags = core.SpellFlagChanneled | core.SpellFlagAPL
	config.ManaCost = core.ManaCostOptions{
		BaseCostPercent: 28,
	}

	config.Cast = core.CastConfig{
		DefaultCast: core.Cast{
			GCD: core.GCDDefault,
		},
	}

	config.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: "MindSear-" + priest.Label,
		},
		NumberOfTicks:       5,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			mindSearTickSpell.Cast(sim, target)
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
		if result.Landed() {
			spell.Dot(target).Apply(sim)
		}
	}
	config.ExpectedTickDamage = func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
		baseDamage := priest.ClassSpellScaling * 0.23
		return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
	}

	return priest.RegisterSpell(config)
}
