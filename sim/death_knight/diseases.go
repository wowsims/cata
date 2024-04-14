package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (dk *DeathKnight) AllDiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverSpell.Dot(target).IsActive() && dk.BloodPlagueSpell.Dot(target).IsActive()
}

func (dk *DeathKnight) DiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverSpell.Dot(target).IsActive() || dk.BloodPlagueSpell.Dot(target).IsActive()
}

// func (dk *DeathKnight) DrwDiseasesAreActive(target *core.Unit) bool {
// 	return dk.Talents.DancingRuneWeapon && dk.RuneWeapon.FrostFeverSpell.Dot(target).IsActive() || dk.RuneWeapon.BloodPlagueSpell.Dot(target).IsActive()
// }

func (dk *DeathKnight) dkCountActiveDiseases(target *core.Unit) float64 {
	count := 0
	if dk.FrostFeverSpell.Dot(target).IsActive() {
		count++
	}
	if dk.BloodPlagueSpell.Dot(target).IsActive() {
		count++
	}
	if dk.Talents.EbonPlaguebringer > 0 && dk.EbonPlagueBringerAura.Get(target).IsActive() {
		count++
	}
	return float64(count)
}

func (dk *DeathKnight) registerDiseaseDots() {
	dk.registerFrostFever()
	dk.registerBloodPlague()
}

func (dk *DeathKnight) registerFrostFever() {
	extraEffectAura := dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.FrostFeverAura(&dk.Unit, target, dk.Talents.BrittleBones)
	})

	dk.FrostFeverSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55095},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease,
		ClassSpellMask: DeathKnightSpellFrostFever,

		DamageMultiplier: 1.15,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FrostFever" + dk.Label,
				Tag:   "FrostFever",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					extraEffectAura.Get(aura.Unit).Activate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					extraEffectAura.Get(aura.Unit).Deactivate(sim)
				},
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.MeleeAttackPower()*0.055+0.31999999285*core.CharacterLevel)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if dk.Talents.EbonPlaguebringer > 0 {
				dk.EbonPlagueBringerAura.Get(target).Activate(sim)
			}

			dot := spell.Dot(target)
			dot.Apply(sim)
		},

		RelatedAuras: []core.AuraArray{extraEffectAura},
	})
}

func (dk *DeathKnight) registerBloodPlague() {
	dk.BloodPlagueSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55078},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease,
		ClassSpellMask: DeathKnightSpellBloodPlague,

		DamageMultiplier: 1.15,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "BloodPlague" + dk.Label,
				Tag:   "BloodPlague",
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.MeleeAttackPower()*0.055+0.3939999938*core.CharacterLevel)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if dk.Talents.EbonPlaguebringer > 0 {
				dk.EbonPlagueBringerAura.Get(target).Activate(sim)
			}

			spell.Dot(target).Apply(sim)
		},
	})
}

// func (dk *DeathKnight) drwCountActiveDiseases(target *core.Unit) float64 {
// 	count := 0
// 	if dk.Talents.DancingRuneWeapon {
// 		if dk.RuneWeapon.FrostFeverSpell.Dot(target).IsActive() {
// 			count++
// 		}
// 		if dk.RuneWeapon.BloodPlagueSpell.Dot(target).IsActive() {
// 			count++
// 		}
// 	}
// 	return float64(count)
// }

// func (dk *DeathKnight) registerDrwDiseaseDots() {
// 	dk.registerDrwFrostFever()
// 	dk.registerDrwBloodPlague()
// }

// func (dk *DeathKnight) registerDrwFrostFever() {
// 	dk.RuneWeapon.FrostFeverSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
// 		ActionID:    core.ActionID{SpellID: 55095},
// 		SpellSchool: core.SpellSchoolFrost,
// 		ProcMask:    core.ProcMaskSpellDamage,
// 		Flags:       core.SpellFlagDisease,

// 		DamageMultiplier: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
// 		ThreatMultiplier: 1,

// 		Dot: core.DotConfig{
// 			Aura: core.Aura{
// 				Label: "DrwFrostFever",
// 			},
// 			NumberOfTicks: 5 + dk.Talents.Epidemic,
// 			TickLength:    time.Second * 3,
// 			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
// 				// 80.0 * 0.32 * 1.15 base, 0.055 * 1.15
// 				dot.SnapshotBaseDamage = 29.44 + 0.06325*dk.getImpurityBonus(dot.Spell)

// 				if !isRollover {
// 					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
// 				}
// 			},
// 			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
// 				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			spell.Dot(target).Apply(sim)
// 		},
// 	})
// }

// func (dk *DeathKnight) registerDrwBloodPlague() {
// 	// Tier9 4Piece
// 	canCrit := dk.HasSetBonus(ItemSetThassariansBattlegear, 4)

// 	dk.RuneWeapon.BloodPlagueSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
// 		ActionID:    core.ActionID{SpellID: 55078},
// 		SpellSchool: core.SpellSchoolShadow,
// 		ProcMask:    core.ProcMaskSpellDamage,
// 		Flags:       core.SpellFlagDisease,

// 		DamageMultiplier: 1,
// 		CritMultiplier:   dk.RuneWeapon.DefaultMeleeCritMultiplier(),
// 		ThreatMultiplier: 1,

// 		Dot: core.DotConfig{
// 			Aura: core.Aura{
// 				Label: "DrwBloodPlague",
// 			},
// 			NumberOfTicks: 5 + dk.Talents.Epidemic,
// 			TickLength:    time.Second * 3,

// 			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
// 				// 80.0 * 0.394 * 1.15 for base, 0.055 * 1.15 for ap coeff
// 				dot.SnapshotBaseDamage = 36.248 + 0.06325*dk.getImpurityBonus(dot.Spell)

// 				if !isRollover {
// 					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
// 					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
// 				}
// 			},
// 			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
// 				var result *core.SpellResult
// 				if canCrit {
// 					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
// 				} else {
// 					result = dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
// 				}
// 				dk.doWanderingPlague(sim, dot.Spell, result)
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			spell.Dot(target).Apply(sim)
// 		},
// 	})
// }
