package death_knight

import (
	"github.com/wowsims/cata/sim/core"
)

var BloodBoilActionID = core.ActionID{SpellID: 48721}

func (dk *DeathKnight) registerBloodBoilSpell() {
	rpMetric := dk.NewRunicPowerMetrics(BloodBoilActionID)
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
			results := make([]*core.SpellResult, sim.GetNumTargets())
			anyHit := false
			for idx, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := dk.ClassBaseScaling*0.31700000167 + 0.08*spell.MeleeAttackPower()
				baseDamage *= core.TernaryFloat64(dk.DiseasesAreActive(aoeTarget), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				results[idx] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				anyHit = anyHit || results[idx].Landed()
			}

			if anyHit {
				dk.AddRunicPower(sim, 10, rpMetric)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}

// func (dk *DeathKnight) registerDrwBloodBoilSpell() {
// 	dk.RuneWeapon.BloodBoil = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
// 		ActionID:    BloodBoilActionID,
// 		SpellSchool: core.SpellSchoolShadow,
// 		ProcMask:    core.ProcMaskSpellDamage,

// 		DamageMultiplier: dk.bloodyStrikesBonus(BloodyStrikesBB),
// 		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			for _, aoeTarget := range sim.Encounter.TargetUnits {
// 				baseDamage := (sim.Roll(180, 220) + 0.06*dk.RuneWeapon.getImpurityBonus(spell)) * core.TernaryFloat64(dk.DrwDiseasesAreActive(aoeTarget), 1.5, 1.0)
// 				baseDamage *= sim.Encounter.AOECapMultiplier()

// 				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
// 			}
// 		},
// 	})
// }
