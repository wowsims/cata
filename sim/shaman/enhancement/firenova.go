package enhancement

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/shaman"
)

func (enh *EnhancementShaman) registerFireNovaSpell() {

	results := make([][]*core.SpellResult, enh.Env.GetNumTargets())
	for i := range enh.Env.GetNumTargets() {
		results[i] = make([]*core.SpellResult, enh.Env.GetNumTargets())
	}

	for range enh.Env.GetNumTargets() {
		nova := enh.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: 1535},
			SpellSchool:    core.SpellSchoolFire,
			ProcMask:       core.ProcMaskSpellDamage,
			Flags:          shaman.SpellFlagShamanSpell | core.SpellFlagAoE,
			ClassSpellMask: shaman.SpellMaskFireNova,

			ApplyEffects: func(sim *core.Simulation, mainTarget *core.Unit, spell *core.Spell) {
				for j, target := range sim.Encounter.TargetUnits {
					if target != mainTarget {
						spell.DealDamage(sim, results[mainTarget.Index][j])
					}
				}
			},
		})
		enh.FireNovas = append(enh.FireNovas, nova)
	}

	enh.FireNova = enh.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1535},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          shaman.SpellFlagShamanSpell | core.SpellFlagAPL | core.SpellFlagAoE,
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
			for i, mainTarget := range sim.Encounter.TargetUnits {
				//need to calculate damage even from non flame shocked target in case echo procs from it
				for j, target := range sim.Encounter.TargetUnits {
					if mainTarget != target {
						baseDamage := enh.CalcAndRollDamageRange(sim, 1.43599998951, 0.15000000596)
						results[i][j] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
					}
				}
			}
			for i, mainTarget := range sim.Encounter.TargetUnits {
				if enh.FlameShock.Dot(mainTarget).IsActive() {
					enh.FireNovas[i].Cast(sim, mainTarget)
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
