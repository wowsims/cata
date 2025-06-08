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

func getFrostFeverConfig(character *core.Character) core.SpellConfig {
	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55095},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellFrostFever,

		DamageMultiplier: 1.15,
		CritMultiplier:   character.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Frost Fever" + character.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseTickDamage := character.CalcScalingSpellDmg(0.13300000131) + dot.Spell.MeleeAttackPower()*0.15800000727
				dot.Snapshot(target, baseTickDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	}
}

func (dk *DeathKnight) registerFrostFever() {
	config := getFrostFeverConfig(dk.GetCharacter())
	dk.FrostFeverSpell = dk.RegisterSpell(config)
}

func getBloodPlagueConfig(character *core.Character) core.SpellConfig {
	return core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55078},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellBloodPlague,

		DamageMultiplier: 1.15,
		CritMultiplier:   character.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Blood Plague" + character.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseTickDamage := character.CalcScalingSpellDmg(0.3939999938) + dot.Spell.MeleeAttackPower()*0.15800000727
				dot.Snapshot(target, baseTickDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	}
}

func (dk *DeathKnight) registerBloodPlague() {
	dk.BloodPlagueSpell = dk.RegisterSpell(getBloodPlagueConfig(dk.GetCharacter()))
}

func (dk *DeathKnight) registerDrwFrostFever() {
	config := getFrostFeverConfig(dk.RuneWeapon.GetCharacter())
	oldApplyEffects := config.ApplyEffects
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		CopySpellMultipliers(dk.FrostFeverSpell, dk.RuneWeapon.FrostFeverSpell, target)
		oldApplyEffects(sim, target, spell)
	}

	dk.RuneWeapon.FrostFeverSpell = dk.RuneWeapon.RegisterSpell(config)
}

func (dk *DeathKnight) registerDrwBloodPlague() {
	config := getBloodPlagueConfig(dk.RuneWeapon.GetCharacter())
	oldApplyEffects := config.ApplyEffects
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		CopySpellMultipliers(dk.BloodPlagueSpell, dk.RuneWeapon.BloodPlagueSpell, target)
		oldApplyEffects(sim, target, spell)
	}

	dk.RuneWeapon.BloodPlagueSpell = dk.RuneWeapon.RegisterSpell(config)
}
