package warlock

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

type EbonImpPet struct {
	core.Pet

	owner *Warlock
}

// TODO: This file is one big TODO
func (warlock *Warlock) NewEbonImp() *EbonImpPet {
	imp := &EbonImpPet{
		Pet:   core.NewPet("Ebon Imp", &warlock.Character, warlock.ebonImpBaseStats(), warlock.ebonImpStatInheritance(), false, true),
		owner: warlock,
	}

	imp.EnableAutoAttacks(imp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			// Base 240 DPS with observed around 300 range
			BaseDamageMin:     (240 - 75) * 2,
			BaseDamageMax:     (240 + 75) * 2,
			SwingSpeed:        2,
			CritMultiplier:    2,
			AttackPowerPerDPS: 1, //TODO,
		},
		AutoSwingMelee: true,
	})

	imp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warlock.AddPet(imp)
	imp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/324.72)

	return imp
}

func (warlock *Warlock) SetupEbonImp(imp *EbonImpPet) {

}

func (imp *EbonImpPet) GetPet() *core.Pet {
	return &imp.Pet
}

func (imp *EbonImpPet) Initialize() {
}

func (imp *EbonImpPet) Reset(_ *core.Simulation) {

}

func (imp *EbonImpPet) ExecuteCustomRotation(sim *core.Simulation) {
}

// TODO: copied from ghoul pet
func (warlock *Warlock) ebonImpBaseStats() stats.Stats {
	return stats.Stats{
		stats.Stamina:     388,
		stats.Agility:     3343 - 10, // We remove 10 to not mess with crit conversion
		stats.Strength:    476,
		stats.AttackPower: -20,
	}
}

// TODO: copied from ghoul pet
func (warlock *Warlock) ebonImpStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:  ownerStats[stats.Stamina] * 0.75,
			stats.Strength: ownerStats[stats.Strength] * 1.0,

			stats.MeleeHit:  ownerStats[stats.MeleeHit],
			stats.Expertise: ownerStats[stats.MeleeHit] * PetExpertiseScale,

			stats.MeleeHaste: ownerStats[stats.MeleeHaste],
			stats.MeleeCrit:  ownerStats[stats.MeleeCrit],
		}
	}
}
