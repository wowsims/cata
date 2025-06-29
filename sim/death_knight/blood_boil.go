package death_knight

import (
	"github.com/wowsims/cata/sim/core"
)

var BloodBoilActionID = core.ActionID{SpellID: 48721}

func (dk *DeathKnight) registerBloodBoilSpell() {
	rpMetric := dk.NewRunicPowerMetrics(BloodBoilActionID)
	results := make([]*core.SpellResult, dk.Env.TotalTargetCount())
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       BloodBoilActionID,
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellBloodBoil,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			anyHit := false
			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				baseDamage := dk.ClassSpellScaling*0.31700000167 + 0.08*spell.MeleeAttackPower()
				baseDamage *= core.TernaryFloat64(dk.DiseasesAreActive(aoeTarget), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				results[aoeTarget.Index] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				anyHit = anyHit || results[aoeTarget.Index].Landed()
			}

			if anyHit {
				dk.AddRunicPower(sim, 10, rpMetric)
			}

			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				spell.DealDamage(sim, results[aoeTarget.Index])
			}
		},
	})
}

func (dk *DeathKnight) registerDrwBloodBoilSpell() *core.Spell {
	results := make([]*core.SpellResult, dk.Env.TotalTargetCount())
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    BloodBoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				baseDamage := dk.ClassSpellScaling*0.31700000167 + 0.08*spell.MeleeAttackPower()
				baseDamage *= core.TernaryFloat64(dk.RuneWeapon.DiseasesAreActive(aoeTarget), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				results[aoeTarget.Index] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			for _, aoeTarget := range sim.Encounter.ActiveTargetUnits {
				spell.DealDamage(sim, results[aoeTarget.Index])
			}
		},
	})
}
