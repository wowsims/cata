package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type BloodwormPet struct {
	core.Pet

	stackAura *core.Aura
	dkOwner   *DeathKnight
}

func (dk *DeathKnight) NewBloodwormPet(_ int) *BloodwormPet {
	bloodworm := &BloodwormPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Bloodworm",
			Owner:                           &dk.Character,
			BaseStats:                       stats.Stats{},
			StatInheritance:                 dk.bloodwormStatInheritance(),
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		dkOwner: dk,
	}

	baseDamage := dk.CalcScalingSpellDmg(0.55)
	bloodworm.EnableAutoAttacks(bloodworm, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  baseDamage,
			BaseDamageMax:  baseDamage,
			SwingSpeed:     2,
			CritMultiplier: dk.DefaultCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	bloodworm.OnPetDisable = bloodworm.disable

	// command doesn't apply to Bloodworms
	if dk.Race == proto.Race_RaceOrc {
		bloodworm.PseudoStats.DamageDealtMultiplier /= 1.02
	}

	dk.AddPet(bloodworm)

	return bloodworm
}

func (bloodworm *BloodwormPet) GetPet() *core.Pet {
	return &bloodworm.Pet
}

func (bloodworm *BloodwormPet) Initialize() {
	bloodworm.stackAura = bloodworm.GetOrRegisterAura(core.Aura{
		Label:     "Blood Gorged Proc",
		ActionID:  core.ActionID{SpellID: 81277},
		Duration:  core.NeverExpires,
		MaxStacks: 12,
	})

	healSpell := bloodworm.dkOwner.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 81280},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellHealing,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
	})

	explosion := bloodworm.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 81280},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellHealing,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healSpell.Cast(sim, target)

			for _, target := range sim.Raid.AllPlayerUnits {
				if target == &bloodworm.Unit {
					continue
				}

				if target.DistanceFromTarget > core.MaxMeleeRange {
					continue
				}

				wormHp := 4000 + (bloodworm.dkOwner.GetStat(stats.Stamina)-10)*3.45100
				healSpell.CalcAndDealHealing(sim, target, float64(bloodworm.stackAura.GetStacks())*0.05*wormHp, healSpell.OutcomeHealing)
			}
			bloodworm.Pet.Disable(sim)
		},
	})

	explodeChances := []float64{0, 0, 0, 0.005, 0.005, 0.11, 0.18, 0.27, 0.36, 0.49, 0.75, 0.85, 1}
	core.MakeProcTriggerAura(&bloodworm.Unit, core.ProcTrigger{
		Name:     "Blood Gorged",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if bloodworm.stackAura.IsActive() && bloodworm.stackAura.GetStacks() >= 3 {
				explodeChance := explodeChances[bloodworm.stackAura.GetStacks()]
				if sim.Proc(explodeChance, "Blood Burst") {
					explosion.Cast(sim, &bloodworm.Unit)
					return
				}
			}

			bloodworm.stackAura.Activate(sim)
			bloodworm.stackAura.AddStack(sim)
		},
	})
}

func (bloodworm *BloodwormPet) Reset(_ *core.Simulation) {
}

func (bloodworm *BloodwormPet) ExecuteCustomRotation(_ *core.Simulation) {
}

func (bloodworm *BloodwormPet) disable(sim *core.Simulation) {
	bloodworm.stackAura.Deactivate(sim)
}

func (dk *DeathKnight) bloodwormStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		hitRating := ownerStats[stats.HitRating]
		expertiseRating := ownerStats[stats.ExpertiseRating]
		combined := (hitRating + expertiseRating) * 0.5

		return stats.Stats{
			stats.Armor:               ownerStats[stats.Armor],
			stats.AttackPower:         ownerStats[stats.AttackPower] * 0.55,
			stats.CritRating:          ownerStats[stats.CritRating],
			stats.ExpertiseRating:     combined,
			stats.HasteRating:         ownerStats[stats.HasteRating],
			stats.Health:              ownerStats[stats.Health],
			stats.HitRating:           combined,
			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
			stats.Stamina:             ownerStats[stats.Stamina],
		}
	}
}
