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
	if dk.EbonPlagueAura.Get(target).IsActive() {
		count++
	}
	if count < 2 && dk.BurningBloodSpell != nil && dk.BurningBloodSpell.Dot(target).IsActive() {
		count = 2
	}
	return base + increase*float64(count)
}

var OutbreakActionID = core.ActionID{SpellID: 77575}

func (dk *DeathKnight) registerOutbreak() {
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       OutbreakActionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellOutbreak,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				dk.FrostFeverSpell.Cast(sim, target)
				dk.BloodPlagueSpell.Cast(sim, target)
			}
		},
	})
}

func (dk *DeathKnight) registerFrostFever() {
	extraEffectAura := dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.FrostFeverAura(target, dk.Talents.BrittleBones)
	})

	dk.FrostFeverSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55095},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellFrostFever,

		DamageMultiplier: 1.15,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FrostFever" + dk.Label,
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
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},

		RelatedAuraArrays: extraEffectAura.ToMap(),
	})
}

func (dk *DeathKnight) registerBloodPlague() {
	dk.BloodPlagueSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55078},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellBloodPlague,

		DamageMultiplier: 1.15,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "BloodPlague" + dk.Label,
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.MeleeAttackPower()*0.055+0.3939999938*core.CharacterLevel)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}

func (runeWeapon *RuneWeaponPet) DiseasesAreActive(target *core.Unit) bool {
	return runeWeapon.FrostFeverSpell.Dot(target).IsActive() || runeWeapon.BloodPlagueSpell.Dot(target).IsActive()
}

func (runeWeapon *RuneWeaponPet) GetDiseaseMulti(target *core.Unit, base float64, increase float64) float64 {
	count := 0
	if runeWeapon.FrostFeverSpell.Dot(target).IsActive() {
		count++
	}
	if runeWeapon.BloodPlagueSpell.Dot(target).IsActive() {
		count++
	}
	return base + increase*float64(count)
}

func (dk *DeathKnight) registerDrwOutbreakSpell() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:       OutbreakActionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellOutbreak,

		CritMultiplier: dk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				dk.RuneWeapon.FrostFeverSpell.Cast(sim, target)
				dk.RuneWeapon.BloodPlagueSpell.Cast(sim, target)
			}
		},
	})
}

func (dk *DeathKnight) registerDrwFrostFever() {
	dk.RuneWeapon.FrostFeverSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55095},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagDisease | core.SpellFlagPassiveSpell,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FrostFever" + dk.RuneWeapon.Label,
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.MeleeAttackPower()*0.055+0.31999999285*core.CharacterLevel)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			CopySpellMultipliers(dk.FrostFeverSpell, dk.RuneWeapon.FrostFeverSpell, target)
			spell.Dot(target).Apply(sim)
		},
	})
}

func (dk *DeathKnight) registerDrwBloodPlague() {
	dk.RuneWeapon.BloodPlagueSpell = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55078},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagDisease | core.SpellFlagPassiveSpell,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "BloodPlague" + dk.RuneWeapon.Label,
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.MeleeAttackPower()*0.055+0.3939999938*core.CharacterLevel)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			CopySpellMultipliers(dk.BloodPlagueSpell, dk.RuneWeapon.BloodPlagueSpell, target)
			spell.Dot(target).Apply(sim)
		},
	})
}
