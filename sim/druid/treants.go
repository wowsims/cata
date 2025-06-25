package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// Extension of PetAgent interface, for treants.
type TreantAgent interface {
	core.PetAgent

	Enable(sim *core.Simulation)
}

// Embed this in spec-specific treant structs.
type DefaultTreantImpl struct {
	core.Pet
}

// Overwrite these for spec variants that register spells.
func (treant *DefaultTreantImpl) Initialize()                              {}
func (treant *DefaultTreantImpl) ExecuteCustomRotation(_ *core.Simulation) {}

func (treant *DefaultTreantImpl) Reset(sim *core.Simulation) {
	treant.Disable(sim)
}

func (treant *DefaultTreantImpl) GetPet() *core.Pet {
	return &treant.Pet
}

func (treant *DefaultTreantImpl) Enable(sim *core.Simulation) {
	treant.EnableWithTimeout(sim, treant, time.Second*15)
}

type TreantConfig struct {
	NonHitExpStatInheritance core.PetStatInheritance
	EnableAutos              bool
	WeaponDamageCoefficient  float64
}

func (druid *Druid) NewDefaultTreant(config TreantConfig) *DefaultTreantImpl {
	treant := &DefaultTreantImpl{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Treant",
			Owner:                           &druid.Character,
			NonHitExpStatInheritance:        config.NonHitExpStatInheritance,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
	}

	if !config.EnableAutos {
		return treant
	}

	baseWeaponDamage := config.WeaponDamageCoefficient * druid.ClassSpellScaling

	treant.EnableAutoAttacks(treant, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        baseWeaponDamage,
			BaseDamageMax:        baseWeaponDamage,
			SwingSpeed:           2,
			NormalizedSwingSpeed: 2,
			CritMultiplier:       druid.DefaultCritMultiplier(),
			SpellSchool:          core.SpellSchoolPhysical,
		},

		AutoSwingMelee: true,
	})

	treant.OnPetEnable = func(sim *core.Simulation) {
		// Treant spawns in front of boss but moves behind after first swing.
		treant.PseudoStats.InFrontOfTarget = true
		pa := sim.GetConsumedPendingActionFromPool()
		pa.NextActionAt = sim.CurrentTime + time.Millisecond*500

		pa.OnAction = func(_ *core.Simulation) {
			treant.PseudoStats.InFrontOfTarget = false
		}

		sim.AddPendingAction(pa)
	}

	return treant
}

type TreantAgents [3]TreantAgent
