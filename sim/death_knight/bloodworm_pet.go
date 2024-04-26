package death_knight

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

type BloodwormPet struct {
	core.Pet

	stackAura *core.Aura
	dkOwner   *DeathKnight
}

func (dk *DeathKnight) NewBloodwormPet(_ int) *BloodwormPet {
	bloodworm := &BloodwormPet{
		Pet:     core.NewPet("Bloodworm", &dk.Character, bloodwormPetBaseStats, dk.bloodwormStatInheritance(), false, true),
		dkOwner: dk,
	}

	weapon := dk.WeaponFromMainHand(2)

	bloodworm.EnableAutoAttacks(bloodworm, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  weapon.BaseDamageMin,
			BaseDamageMax:  weapon.BaseDamageMax,
			SwingSpeed:     2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	bloodworm.AddStatDependency(stats.Strength, stats.AttackPower, 1.0+1)
	bloodworm.AddStatDependency(stats.Agility, stats.MeleeCrit, 1.0+(core.CritRatingPerCritChance/83.3))

	bloodworm.OnPetEnable = bloodworm.enable
	bloodworm.OnPetDisable = bloodworm.disable

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
		MaxStacks: 10,
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

				healSpell.CalcAndDealHealing(sim, target, float64(bloodworm.stackAura.GetStacks())*0.3*(bloodworm.dkOwner.MaxHealth()*0.18), healSpell.OutcomeHealing)
			}
			bloodworm.Pet.Disable(sim)
		},
	})

	core.MakeProcTriggerAura(&bloodworm.Unit, core.ProcTrigger{
		Name:     "Blood Gorged",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodworm.stackAura.Activate(sim)
			bloodworm.stackAura.AddStack(sim)

			explodeChance := float64(bloodworm.stackAura.GetStacks()*bloodworm.stackAura.GetStacks()) * 0.01

			if sim.Proc(explodeChance, "Blood Burst") {
				explosion.Cast(sim, &bloodworm.Unit)
			}
		},
	})
}

func (bloodworm *BloodwormPet) Reset(_ *core.Simulation) {
}

func (bloodworm *BloodwormPet) ExecuteCustomRotation(_ *core.Simulation) {
}

func (bloodworm *BloodwormPet) enable(sim *core.Simulation) {
	// Snapshot extra % speed modifiers from dk owner
	bloodworm.PseudoStats.MeleeSpeedMultiplier = 1
	bloodworm.MultiplyMeleeSpeed(sim, bloodworm.dkOwner.PseudoStats.MeleeSpeedMultiplier)
}

func (bloodworm *BloodwormPet) disable(sim *core.Simulation) {
	// Clear snapshot speed
	bloodworm.PseudoStats.MeleeSpeedMultiplier = 1
	bloodworm.MultiplyMeleeSpeed(sim, 1)
	bloodworm.stackAura.Deactivate(sim)
}

var bloodwormPetBaseStats = stats.Stats{
	stats.MeleeCrit: 8 * core.CritRatingPerCritChance,
}

func (dk *DeathKnight) bloodwormStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.AttackPower: ownerStats[stats.AttackPower] * 0.112,
			stats.MeleeHaste:  ownerStats[stats.MeleeHaste],
		}
	}
}
