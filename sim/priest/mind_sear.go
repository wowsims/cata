package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (priest *Priest) getMindSearBaseConfig() core.SpellConfig {
	return core.SpellConfig{
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskProc,
		ClassSpellMask:           PriestSpellMindSear,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           priest.DefaultSpellCritMultiplier(),
		BonusCoefficient:         0.2622,
	}
}

func (priest *Priest) getMindSearTickSpell() *core.Spell {
	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: 48045}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		damage := priest.ClassSpellScaling * 0.23
		spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
	}
	return priest.RegisterSpell(config)
}

func (priest *Priest) newMindSearSpell() *core.Spell {
	mindSearTickSpell := priest.getMindSearTickSpell()

	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: 48045}
	config.Flags = core.SpellFlagChanneled | core.SpellFlagNoMetrics | core.SpellFlagAPL
	config.ManaCost = core.ManaCostOptions{
		BaseCost: 0.28,
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
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				mindSearTickSpell.Cast(sim, aoeTarget)
				mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
			}
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
		if result.Landed() {
			spell.Dot(target).Apply(sim)
			mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts += 1
		}
	}
	config.ExpectedTickDamage = func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
		baseDamage := priest.ClassSpellScaling * 0.23
		return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
	}

	return priest.RegisterSpell(config)
}
