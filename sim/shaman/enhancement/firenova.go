package enhancement

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/shaman"
)

func (enh *EnhancementShaman) registerFireNovaSpell() {
	enh.FireNova = enh.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1535},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          shaman.SpellFlagFocusable | core.SpellFlagAPL,
		ClassSpellMask: shaman.SpellMaskFireNova,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 13.7,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    enh.NewTimer(),
				Duration: time.Second * time.Duration(4),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   enh.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.30000001192,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([][]*core.SpellResult, enh.Env.GetNumTargets())
			baseDamage := enh.CalcAndRollDamageRange(sim, 1.43599998951, 0.15000000596)
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				if enh.FlameShock.Dot(aoeTarget).IsActive() {
					results[i] = make([]*core.SpellResult, enh.Env.GetNumTargets())
					for j, newTarget := range sim.Encounter.TargetUnits {
						if newTarget != aoeTarget {
							results[i][j] = spell.CalcDamage(sim, newTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
						}
					}
				}
			}
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				if enh.FlameShock.Dot(aoeTarget).IsActive() {
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
				if enh.FlameShock.Dot(aoeTarget).IsActive() {
					return true
				}
			}
			return false
		},
	})
}
