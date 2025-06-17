package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type Xuen struct {
	core.Pet

	owner                   *Monk
	CracklingTigerLightning *core.Spell
}

var baseStats = stats.Stats{
	stats.Strength:    0,
	stats.Agility:     0,
	stats.Stamina:     0,
	stats.Intellect:   0,
	stats.AttackPower: 0,
	stats.Mana:        0,
}

func (monk *Monk) NewXuen() *Xuen {
	xuen := &Xuen{
		Pet: core.NewPet(core.PetConfig{
			Name:      "Xuen, The White Tiger",
			Owner:     &monk.Character,
			BaseStats: baseStats,
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {

				hitRating := ownerStats[stats.HitRating]
				expertiseRating := ownerStats[stats.ExpertiseRating]
				combinedHitExp := (hitRating + expertiseRating) * 0.5

				return stats.Stats{
					stats.Stamina:     ownerStats[stats.Stamina],
					stats.AttackPower: ownerStats[stats.AttackPower] * 0.5,

					stats.HitRating:       combinedHitExp,
					stats.ExpertiseRating: combinedHitExp,
					stats.DodgeRating:     ownerStats[stats.DodgeRating],
					stats.ParryRating:     ownerStats[stats.ParryRating],

					stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
					stats.SpellCritPercent:    ownerStats[stats.PhysicalCritPercent],
				}
			},
			EnabledOnStart:                  false,
			IsGuardian:                      false,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		owner: monk,
	}

	xuen.OnPetEnable = xuen.enable

	actionID := core.ActionID{SpellID: 123996}
	xuen.CracklingTigerLightning = xuen.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		DamageMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Millisecond * 500,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    xuen.NewTimer(),
				Duration: time.Second * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalcScalingSpellDmg(0.293) + xuen.GetStat(stats.AttackPower)*0.505
			for index, target := range sim.Encounter.TargetUnits {
				if index > 3 {
					break
				}
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	xuen.PseudoStats.DamageTakenMultiplier *= 0.1

	// Observed values for Xuen's auto attack damage
	// This could be either:
	// ClassBaseScaling * 1.05853604195
	// CreatureDPS (201.889276) * 5.73988599855808
	// or something completely different
	baseWeaponDamage := 1157.9
	xuen.EnableAutoAttacks(xuen, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        baseWeaponDamage,
			BaseDamageMax:        baseWeaponDamage + 1,
			SwingSpeed:           1,
			NormalizedSwingSpeed: 1,
			CritMultiplier:       monk.DefaultCritMultiplier(),
			SpellSchool:          core.SpellSchoolNature,
		},
		AutoSwingMelee: true,
	})

	xuen.AutoAttacks.MHConfig().BonusCoefficient = 0
	xuen.AutoAttacks.MHConfig().Flags |= core.SpellFlagIgnoreTargetModifiers

	monk.AddPet(xuen)

	return xuen
}

func (xuen *Xuen) Initialize() {
}

func (xuen *Xuen) ExecuteCustomRotation(sim *core.Simulation) {
	if xuen.CracklingTigerLightning.CanCast(sim, xuen.CurrentTarget) {
		xuen.CracklingTigerLightning.Cast(sim, xuen.CurrentTarget)
	}
}

func (xuen *Xuen) Reset(sim *core.Simulation) {
	xuen.Disable(sim)
}

func (xuen *Xuen) enable(sim *core.Simulation) {
}

func (xuen *Xuen) GetPet() *core.Pet {
	return &xuen.Pet
}
