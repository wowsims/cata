package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

var SEFSpellID = int32(138228)

func (monk *Monk) registerStormEarthAndFire() {
	if monk.Spec != proto.Spec_SpecWindwalkerMonk {
		return
	}

	var sefTarget *core.Unit
	damageMultiplier := []float64{1, 0.70, 0.55}

	sefAura := monk.RegisterAura(core.Aura{
		Label:     "Storm, Earth, and Fire",
		ActionID:  core.ActionID{SpellID: SEFSpellID},
		Duration:  core.NeverExpires,
		MaxStacks: 2,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			monk.SefController.CastCopySpell(sim, spell)
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			// We only care if stacks are increasing
			// as decreasing would mean disabling SEF
			if newStacks > oldStacks {
				monk.SefController.PickClone(sim, sefTarget)
				monk.PseudoStats.DamageDealtMultiplier /= damageMultiplier[oldStacks]
				monk.PseudoStats.DamageDealtMultiplier *= damageMultiplier[newStacks]
				for _, pet := range monk.SefController.pets {
					pet.PseudoStats.DamageDealtMultiplier /= damageMultiplier[oldStacks]
					pet.PseudoStats.DamageDealtMultiplier *= damageMultiplier[newStacks]
				}
				return
			}

			aura.Deactivate(sim)
			monk.PseudoStats.DamageDealtMultiplier /= damageMultiplier[oldStacks]
			for _, pet := range monk.SefController.pets {
				pet.PseudoStats.DamageDealtMultiplier /= damageMultiplier[oldStacks]
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			sefTarget = nil
			monk.SefController.Reset(sim)
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: SEFSpellID},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MonkSpellStormEarthAndFire,

		EnergyCost: core.EnergyCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			// If we have 2 clones already, disable SEF
			if sefAura.GetStacks() == 2 {
				sefAura.Deactivate(sim)
				return
			}
			sefTarget = target
			sefAura.Activate(sim)
			sefAura.AddStack(sim)
		},
	})
}

type StormEarthAndFireController struct {
	owner          *Monk
	pets           []*StormEarthAndFirePet
	activeClones   []*StormEarthAndFirePet
	inactiveClones []*StormEarthAndFirePet
}

// Modifies the spell that should be copied with
// Damage Multipliers / Tags etc.
func (pet *StormEarthAndFirePet) modifyCopySpell(sourceSpell *core.Spell, targetSpell *core.Spell) {
	targetSpell.DamageMultiplier = sourceSpell.DamageMultiplier
	targetSpell.DamageMultiplierAdditive = sourceSpell.DamageMultiplierAdditive
	targetSpell.BonusCritPercent = sourceSpell.BonusCritPercent
	targetSpell.BonusHitPercent = sourceSpell.BonusHitPercent
	targetSpell.CritMultiplier = sourceSpell.CritMultiplier
	targetSpell.ThreatMultiplier = sourceSpell.ThreatMultiplier
	targetSpell.BonusCoefficient = sourceSpell.BonusCoefficient

	if sourceSpell.Dot(pet.CurrentTarget) != nil {
		sourceDot := sourceSpell.Dot(pet.CurrentTarget)
		targetDot := targetSpell.Dot(pet.CurrentTarget)

		targetDot.BaseTickCount = sourceDot.BaseTickCount
		targetDot.BaseTickLength = sourceDot.BaseTickLength
	}
}

func (controller *StormEarthAndFireController) CastCopySpell(sim *core.Simulation, spell *core.Spell) {
	for _, pet := range controller.activeClones {
		petSpellActionID := spell.ActionID.WithTag(SEFSpellID)
		copySpell := pet.spells[petSpellActionID]
		if copySpell == nil {
			copySpell = pet.GetSpell(petSpellActionID)
			pet.spells[petSpellActionID] = copySpell
		}
		if copySpell != nil {
			if pet.CurrentTarget == pet.owner.CurrentTarget {
				continue
			}
			pet.modifyCopySpell(spell, copySpell)
			copySpell.Cast(sim, pet.CurrentTarget)
		} else {
			// Break the loop early because the spell doesn't exist
			break
		}
	}
}

func (controller *StormEarthAndFireController) PickClone(sim *core.Simulation, target *core.Unit) {
	clone := controller.getCloneFromTarget(target)
	// If the target already has an active clone, disable it
	if clone != nil {
		controller.deactivateClone(sim, clone)
		return
	}

	// If we have 2 clones already, disable SEF
	if controller.getActiveCloneCount() == 2 {
		controller.Reset(sim)
		return
	}

	// Pick a random clone to spawn from the clones that are not already enabled
	validClones := controller.getInactiveClones()
	cloneIndex := int32(sim.RollWithLabel(0, float64(len(validClones)), "Pick Random Clone"))
	clone = validClones[cloneIndex]

	controller.enableClone(sim, clone, target)
}

func (controller *StormEarthAndFireController) getCloneFromTarget(target *core.Unit) *StormEarthAndFirePet {
	for _, pet := range controller.activeClones {
		if pet.CurrentTarget == target {
			return pet
		}
	}

	return nil
}

func (controller *StormEarthAndFireController) enableClone(sim *core.Simulation, clone *StormEarthAndFirePet, target *core.Unit) {
	clone.CurrentTarget = target
	clone.EnableWithStartAttackDelay(sim, clone, core.DurationFromSeconds(sim.RollWithLabel(2, 2.3, "SEF Spawn Delay")))
	controller.updateActiveClones()
}

func (controller *StormEarthAndFireController) deactivateClone(sim *core.Simulation, clone *StormEarthAndFirePet) {
	clone.Disable(sim)
	controller.updateActiveClones()
}

func (controller *StormEarthAndFireController) updateActiveClones() {
	controller.activeClones = make([]*StormEarthAndFirePet, 0, 3)
	controller.inactiveClones = make([]*StormEarthAndFirePet, 0, 3)
	for _, pet := range controller.pets {
		if pet.IsEnabled() {
			controller.activeClones = append(controller.activeClones, pet)
		} else {
			controller.inactiveClones = append(controller.inactiveClones, pet)
		}
	}
}

func (controller *StormEarthAndFireController) getActiveCloneCount() int32 {
	return int32(len(controller.activeClones) - 1)
}

func (controller *StormEarthAndFireController) getInactiveClones() []*StormEarthAndFirePet {
	return controller.inactiveClones
}

func (controller *StormEarthAndFireController) Reset(sim *core.Simulation) {
	for _, pet := range controller.pets {
		pet.Disable(sim)
	}
	controller.updateActiveClones()
}

func (monk *Monk) RegisterSEFPets() {
	monk.SefController = &StormEarthAndFireController{
		owner: monk,
		pets:  make([]*StormEarthAndFirePet, 0, 3),
	}

	monk.SefController.pets = append(monk.SefController.pets, monk.NewSEFPet("Storm Spirit", 138121, 2.7))
	monk.SefController.pets = append(monk.SefController.pets, monk.NewSEFPet("Earth Spirit", 138122, 3.6))
	monk.SefController.pets = append(monk.SefController.pets, monk.NewSEFPet("Fire Spirit", 138123, 2.7))

	monk.SefController.updateActiveClones()
}

type StormEarthAndFirePet struct {
	core.Pet
	cloneID int32
	owner   *Monk
	spells  map[core.ActionID]*core.Spell
}

func (sefClone *StormEarthAndFirePet) Initialize() {
	// Talents
	sefClone.registerSEFChiBurst()
	sefClone.registerSEFChiWave()
	sefClone.registerSEFRushingJadeWind()

	// Passives - Windwalker
	sefClone.registerSEFCombatConditioning()
	sefClone.registerSEFTigerStrikes()

	// Spells - Monk
	sefClone.registerSEFJab()
	sefClone.registerSEFTigerPalm()
	sefClone.registerSEFBlackoutKick()
	sefClone.registerSEFSpinningCraneKick()

	// Spells - Windwalker
	sefClone.registerSEFRisingSunKick()
	sefClone.registerSEFFistsOfFury()
	sefClone.registerSpinningFireBlossom()
}

func (monk *Monk) NewSEFPet(name string, cloneID int32, swingSpeed float64) *StormEarthAndFirePet {
	sefClone := &StormEarthAndFirePet{
		Pet: core.NewPet(core.PetConfig{
			Name:      name,
			Owner:     &monk.Character,
			BaseStats: stats.Stats{},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {

				hitRating := ownerStats[stats.HitRating]
				expertiseRating := ownerStats[stats.ExpertiseRating]
				combinedHitExp := (hitRating + expertiseRating) * 0.5

				return stats.Stats{
					stats.Stamina:     ownerStats[stats.Stamina] * 0.1,
					stats.AttackPower: ownerStats[stats.AttackPower],

					stats.HitRating:       combinedHitExp,
					stats.ExpertiseRating: combinedHitExp,
					stats.DodgeRating:     ownerStats[stats.DodgeRating],
					stats.ParryRating:     ownerStats[stats.ParryRating],

					stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
					stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
				}
			},
			EnabledOnStart:                  false,
			IsGuardian:                      false,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
		cloneID: cloneID,
		spells:  make(map[core.ActionID]*core.Spell),
		owner:   monk,
	}

	isDualWielding := swingSpeed == 2.7
	mhWeapon := monk.WeaponFromMainHand(monk.DefaultCritMultiplier())
	mhAvgDPS := mhWeapon.DPS()

	// This number is derived from naked Dummy testing and using multiple
	// other weapons. This was the constant difference between them.
	baseCloneDamage := 266.0

	avgWeaponDamageMin := 0.0
	avgWeaponDamageMax := 0.0

	var cloneOhWeapon core.Weapon
	if isDualWielding {
		ohWeapon := monk.WeaponFromOffHand(monk.DefaultCritMultiplier())
		ohAvgDPS := ohWeapon.DPS()
		baseDamage := (mhAvgDPS + (ohAvgDPS / 2)) * swingSpeed * core.TernaryFloat64(ohAvgDPS > 0, DualWieldModifier, 1.0)
		// The DW clone has a tiny variance in auto attack damage
		avgWeaponDamageMin = baseDamage - 1
		avgWeaponDamageMax = baseDamage + 1
		cloneOhWeapon = core.Weapon{
			BaseDamageMin:     baseCloneDamage + avgWeaponDamageMin,
			BaseDamageMax:     baseCloneDamage + avgWeaponDamageMax,
			SwingSpeed:        swingSpeed,
			CritMultiplier:    monk.DefaultCritMultiplier(),
			AttackPowerPerDPS: 0,
		}
	} else {
		avgWeaponDamageMin = mhAvgDPS * swingSpeed
		avgWeaponDamageMax = avgWeaponDamageMin
	}

	cloneMhWeapon := core.Weapon{
		BaseDamageMin:     baseCloneDamage + avgWeaponDamageMin,
		BaseDamageMax:     baseCloneDamage + avgWeaponDamageMax,
		SwingSpeed:        swingSpeed,
		CritMultiplier:    monk.DefaultCritMultiplier(),
		AttackPowerPerDPS: 0,
	}

	sefClone.EnableAutoAttacks(sefClone, core.AutoAttackOptions{
		MainHand:       cloneMhWeapon,
		OffHand:        cloneOhWeapon,
		AutoSwingMelee: true,
	})

	sefClone.OnPetEnable = sefClone.enable
	sefClone.OnPetDisable = sefClone.disable

	monk.AddPet(sefClone)

	return sefClone
}

func (sefClone *StormEarthAndFirePet) GetPet() *core.Pet {
	return &sefClone.Pet
}

func (sefClone *StormEarthAndFirePet) Reset(_ *core.Simulation) {
}

func (sefClone *StormEarthAndFirePet) ExecuteCustomRotation(_ *core.Simulation) {
}

func (sefClone *StormEarthAndFirePet) enable(sim *core.Simulation) {
	if sefClone.AutoAttacks.IsDualWielding {
		sefClone.AutoAttacks.DesyncOffHand(sim, sim.CurrentTime)
	}
}

func (sefClone *StormEarthAndFirePet) disable(sim *core.Simulation) {

}
