package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) DiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverSpell.Dot(target).IsActive() || dk.BloodPlagueSpell.Dot(target).IsActive() || dk.BurningBloodSpell != nil && dk.BurningBloodSpell.Dot(target).IsActive()
}

func (dk *DeathKnight) GetDiseaseMulti(target *core.Unit, base float64, increase float64) float64 {
	count := 0
	if dk.FrostFeverSpell.Dot(target).IsActive() {
		count++
	}
	if dk.BloodPlagueSpell.Dot(target).IsActive() {
		count++
	}
	if count < 2 && dk.BurningBloodSpell != nil && dk.BurningBloodSpell.Dot(target).IsActive() {
		count = 2
	}
	return base + increase*float64(count)
}

func (dk *DeathKnight) getFrostFeverConfig(character *core.Character) core.SpellConfig {
	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55095},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellFrostFever,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Frost Fever" + character.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseTickDamage := dk.CalcScalingSpellDmg(0.13300000131) + dot.Spell.MeleeAttackPower()*0.15800000727
				dot.SnapshotPhysical(target, baseTickDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
	}
}

// A disease dealing (166 + 0.158 * <AP>) Frost damage every 3 sec for 30 sec.
func (dk *DeathKnight) registerFrostFever() {
	config := dk.getFrostFeverConfig(dk.GetCharacter())
	config.DamageMultiplier = 1
	config.CritMultiplier = dk.DefaultCritMultiplier()
	config.ThreatMultiplier = 1
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.Dot(target).Apply(sim)
	}
	config.ExpectedTickDamage = func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
		dot := spell.Dot(target)
		if useSnapshot {
			return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
		} else {
			baseTickDamage := dk.CalcScalingSpellDmg(0.13300000131) + dot.Spell.MeleeAttackPower()*0.15800000727
			return spell.CalcPeriodicDamage(sim, target, baseTickDamage, spell.OutcomeExpectedMagicCrit)
		}
	}

	dk.FrostFeverSpell = dk.RegisterSpell(config)
}

func (dk *DeathKnight) getBloodPlagueConfig(character *core.Character) core.SpellConfig {
	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55078},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellBloodPlague,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Blood Plague" + character.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseTickDamage := dk.CalcScalingSpellDmg(0.15800000727) + dot.Spell.MeleeAttackPower()*0.15800000727
				dot.SnapshotPhysical(target, baseTickDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
	}
}

// A disease dealing (197 + 0.158 * <AP>) Shadow damage every 3 sec for 30 sec.
func (dk *DeathKnight) registerBloodPlague() {
	config := dk.getBloodPlagueConfig(dk.GetCharacter())
	config.DamageMultiplier = 1
	config.CritMultiplier = dk.DefaultCritMultiplier()
	config.ThreatMultiplier = 1
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.Dot(target).Apply(sim)
	}
	config.ExpectedTickDamage = func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
		dot := spell.Dot(target)
		if useSnapshot {
			return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
		} else {
			baseTickDamage := dk.CalcScalingSpellDmg(0.15800000727) + dot.Spell.MeleeAttackPower()*0.15800000727
			return spell.CalcPeriodicDamage(sim, target, baseTickDamage, spell.OutcomeExpectedMagicCrit)
		}
	}

	dk.BloodPlagueSpell = dk.RegisterSpell(config)
}

func (dk *DeathKnight) registerDrwFrostFever() {
	config := dk.getFrostFeverConfig(dk.RuneWeapon.GetCharacter())
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		CopySpellMultipliers(dk.FrostFeverSpell, dk.RuneWeapon.FrostFeverSpell, target)
		spell.Dot(target).Apply(sim)
	}

	dk.RuneWeapon.FrostFeverSpell = dk.RuneWeapon.RegisterSpell(config)
}

func (dk *DeathKnight) registerDrwBloodPlague() {
	config := dk.getBloodPlagueConfig(dk.RuneWeapon.GetCharacter())
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		CopySpellMultipliers(dk.BloodPlagueSpell, dk.RuneWeapon.BloodPlagueSpell, target)
		spell.Dot(target).Apply(sim)
	}

	dk.RuneWeapon.BloodPlagueSpell = dk.RuneWeapon.RegisterSpell(config)
}
