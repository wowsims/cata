package warlock

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

const PetFelhunter string = "Felhunter"
const PetFelguard string = "Felguard"
const PetSuccubus string = "Succubus"
const PetImp string = "Imp"

type WarlockPet struct {
	core.Pet

	CurseOfGuldanDebuffs core.AuraArray
}

func NewWarlockPet(owner *Warlock, name string, baseStats stats.Stats, autoAttackOptions *core.AutoAttackOptions) *WarlockPet {
	warlockPet := &WarlockPet{
		Pet: core.NewPet(name, &owner.Character, baseStats, owner.MakeStatInheritance(), false, false),
	}

	warlockPet.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warlockPet.AddStat(stats.AttackPower, -20)

	warlockPet.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/52.0833)

	if autoAttackOptions != nil {
		warlockPet.EnableAutoAttacks(warlockPet, *autoAttackOptions)
	}

	core.ApplyPetConsumeEffects(&warlockPet.Character, owner.Consumes)

	owner.AddPet(warlockPet)

	return warlockPet
}

func (warlockPet *WarlockPet) GetPet() *core.Pet {
	return &warlockPet.Pet
}

func (warlockPet *WarlockPet) Initialize() {

}

func (warlockPet *WarlockPet) Reset(_ *core.Simulation) {
}
