package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) registerFireNovaSpell() {
	shaman.FireNova = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1535},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagFocusable | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskFireNova,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 22,
			PercentModifier: 100 - (5 * shaman.Talents.Convection) - shaman.GetMentalQuicknessBonus(),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * time.Duration(4),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.164,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([][]*core.SpellResult, shaman.Env.GetNumTargets())
			baseDamage := shaman.CalcAndRollDamageRange(sim, 0.78500002623, 0.11200000346)
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				if shaman.FlameShock.Dot(aoeTarget).IsActive() {
					results[i] = make([]*core.SpellResult, shaman.Env.GetNumTargets())
					for j, newTarget := range sim.Encounter.TargetUnits {
						if newTarget != aoeTarget {
							results[i][j] = spell.CalcDamage(sim, newTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
							results[i][j].Damage *= sim.Encounter.AOECapMultiplier()
						}
					}
				}
			}
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				if shaman.FlameShock.Dot(aoeTarget).IsActive() {
					for j, newTarget := range sim.Encounter.TargetUnits {
						if newTarget != aoeTarget {
							spell.DealDamage(sim, results[i][j])
						}
					}
				}
			}
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				if shaman.FlameShock.Dot(aoeTarget).IsActive() {
					return true
				}
			}
			return false
		},
	})
}
