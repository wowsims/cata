package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Hurls a magical hammer that strikes an enemy for (<1746-1931> + 1.61 * <SP>) Holy damage

-- Sword of Light --
and generates a charge of Holy Power
-- /Sword of Light --

.
Only usable on enemies that have 20% or less health

-- Sword of Light --
or during Avenging Wrath
-- /Sword of Light --

.
*/
func (paladin *Paladin) registerHammerOfWrath() {
	paladin.HammerOfWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 24275},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHammerOfWrath,

		MissileSpeed: 50,
		MaxRange:     30,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 6 * time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20()
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 1.61000001431,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcAndRollDamageRange(sim, 1.61000001431, 0.10000000149)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
