package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const SearCoeff = 0.3
const SearVariance = 0.08
const SearScale = 0.3

func (priest *Priest) getMindSearBaseConfig() core.SpellConfig {
	return core.SpellConfig{
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskSpellProc,
		ClassSpellMask:           PriestSpellMindSear,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		BonusCoefficient:         SearCoeff,
	}
}

func (priest *Priest) getMindSearTickSpell() *core.Spell {
	config := priest.getMindSearBaseConfig()
	config.Flags = core.SpellFlagNoOnDamageDealt | core.SpellFlagAoE
	config.ActionID = core.ActionID{SpellID: 48045}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		damage := priest.CalcAndRollDamageRange(sim, SearScale, SearVariance)
		for _, aoeTarget := range sim.Encounter.TargetUnits {

			// Calc spell damage but deal as periodic for metric purposes
			result := spell.CalcDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCritNoHitCounter)

			// TODO: Verify actual proc behaviour
			// Damage is logged as a tick, has the Flag 'Treat as Periodic' and 'Not a Proc'
			// However, i.E. Trinkets proccing from 'Periodic Damage Dealt' do not trigger
			spell.DealPeriodicDamage(sim, result)

			// For now Sear seems to trigger damage dealt and not periodic dealt for procs
			spell.Unit.OnSpellHitDealt(sim, spell, result)
			result.Target.OnSpellHitTaken(sim, spell, result)

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
		BaseCostPercent: 3,
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
		baseDamage := priest.CalcAndRollDamageRange(sim, SearScale, SearVariance)
		return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
	}

	return priest.RegisterSpell(config)
}
