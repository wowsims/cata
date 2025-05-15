package windwalker

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/monk"
)

func CopySpellMultipliers(sourceSpell *core.Spell, targetSpell *core.Spell, target *core.Unit) {
	targetSpell.DamageMultiplier = sourceSpell.DamageMultiplier
	targetSpell.DamageMultiplierAdditive = sourceSpell.DamageMultiplierAdditive
	targetSpell.BonusCritPercent = sourceSpell.BonusCritPercent
	targetSpell.BonusHitPercent = sourceSpell.BonusHitPercent
	targetSpell.CritMultiplier = sourceSpell.CritMultiplier
	targetSpell.ThreatMultiplier = sourceSpell.ThreatMultiplier

	if sourceSpell.Dot(target) != nil {
		sourceDot := sourceSpell.Dot(target)
		targetDot := targetSpell.Dot(target)

		targetDot.BaseTickCount = sourceDot.BaseTickCount
		targetDot.BaseTickLength = sourceDot.BaseTickLength
	}
}

func (ww *WindwalkerMonk) registerStormEarthAndFire() {
	var sefTarget *core.Unit
	damageMultiplier := []float64{1, 0.70, 0.55}

	sefAura := ww.RegisterAura(core.Aura{
		Label:     "Storm, Earth, and Fire",
		ActionID:  core.ActionID{SpellID: 137639},
		Duration:  core.NeverExpires,
		MaxStacks: 2,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			ww.SefController.CastCopySpell(sim, spell)
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			// We only care if stacks are increasing
			// as decreasing would mean disabling SEF
			if newStacks > oldStacks {
				ww.SefController.PickClone(sim, sefTarget)
				ww.PseudoStats.DamageDealtMultiplier /= damageMultiplier[oldStacks]
				ww.PseudoStats.DamageDealtMultiplier *= damageMultiplier[newStacks]
				for _, pet := range ww.SefController.pets {
					pet.PseudoStats.DamageDealtMultiplier = ww.PseudoStats.DamageDealtMultiplier
				}
			} else {
				ww.SefController.Reset(sim)
				ww.PseudoStats.DamageDealtMultiplier /= damageMultiplier[oldStacks]
				for _, pet := range ww.SefController.pets {
					pet.PseudoStats.DamageDealtMultiplier = ww.PseudoStats.DamageDealtMultiplier
				}
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ww.SefController.Reset(sim)
		},
	})

	ww.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 137639},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellStormEarthAndFire,

		EnergyCost: core.EnergyCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			sefTarget = target
			sefAura.Activate(sim)
			sefAura.AddStack(sim)
		},
	})
}

type StormEarthAndFireController struct {
	owner      *WindwalkerMonk
	pets       []*StormEarthAndFirePet
	spells     map[core.ActionID]*core.Spell
	sefTargets map[int32]*StormEarthAndFirePet
}

func (controller *StormEarthAndFireController) AddCopySpell(actionId core.ActionID, spell *core.Spell) {
	controller.spells[actionId] = spell
}

func (controller *StormEarthAndFireController) GetCopySpell(actionId core.ActionID) *core.Spell {
	return controller.spells[actionId]
}

func (controller *StormEarthAndFireController) CastCopySpell(sim *core.Simulation, spell *core.Spell) {
	copySpell := controller.GetCopySpell(spell.ActionID)
	if copySpell == nil {
		return
	}

	for _, pet := range controller.pets {
		copySpell.Cast(sim, pet.CurrentTarget)
		CopySpellMultipliers(spell, copySpell, pet.CurrentTarget)
	}
}

func (controller *StormEarthAndFireController) PickClone(sim *core.Simulation, target *core.Unit) {
	clone := controller.GetCloneFromTarget(sim, target)
	// If the target already has an active clone, disable it
	if clone != nil {
		controller.Deactivate(sim, clone, target)
		return
	}

	// Pick a random clone to spawn from the clones that are not already enabled
	if controller.GetActiveCloneCount() == 3 {
		controller.Deactivate(sim, clone, target)
		return
	}

	validClones := controller.GetInactiveClones()
	cloneIndex := int32(math.Round(sim.Roll(0, float64(len(validClones)-1))))
	controller.EnableClone(sim, validClones[cloneIndex], target)
}

func (controller *StormEarthAndFireController) GetCloneFromTarget(sim *core.Simulation, target *core.Unit) *StormEarthAndFirePet {
	targetUnixIndex := target.UnitIndex
	return controller.sefTargets[targetUnixIndex]
}

func (controller *StormEarthAndFireController) EnableClone(sim *core.Simulation, clone *StormEarthAndFirePet, target *core.Unit) {
	clone.CurrentTarget = target
	clone.EnableWithStartAttackDelay(sim, clone, core.DurationFromSeconds(sim.RollWithLabel(2, 2.3, "SEF Spawn Delay")))
	controller.sefTargets[target.UnitIndex] = clone
}

func (controller *StormEarthAndFireController) Deactivate(sim *core.Simulation, pet *StormEarthAndFirePet, target *core.Unit) {
	pet.Disable(sim)
	controller.sefTargets[target.UnitIndex] = nil
}

func (controller *StormEarthAndFireController) Reset(sim *core.Simulation) {
	for _, pet := range controller.pets {
		pet.Disable(sim)
	}
	controller.sefTargets = make(map[int32]*StormEarthAndFirePet)
}

func (controller *StormEarthAndFireController) GetInactiveClones() []*StormEarthAndFirePet {
	return core.FilterSlice(controller.pets, func(pet *StormEarthAndFirePet) bool {
		return !pet.IsEnabled()
	})
}

func (controller *StormEarthAndFireController) GetActiveClones() []*StormEarthAndFirePet {
	return core.FilterSlice(controller.pets, func(pet *StormEarthAndFirePet) bool {
		return pet.IsEnabled()
	})
}

func (controller *StormEarthAndFireController) GetActiveCloneCount() int32 {
	return int32(len(controller.GetActiveClones()) - 1)
}

func (ww *WindwalkerMonk) registerSEFPets() {
	ww.SefController = &StormEarthAndFireController{
		owner:      ww,
		spells:     make(map[core.ActionID]*core.Spell),
		pets:       make([]*StormEarthAndFirePet, 0, 3),
		sefTargets: make(map[int32]*StormEarthAndFirePet),
	}

	ww.SefController.pets = append(ww.SefController.pets, ww.NewSEFPet("Storm Spirit", 2.7))
	ww.SefController.pets = append(ww.SefController.pets, ww.NewSEFPet("Earth Spirit", 3.6))
	ww.SefController.pets = append(ww.SefController.pets, ww.NewSEFPet("Fire Spirit", 2.7))
}

type StormEarthAndFirePet struct {
	core.Pet

	owner *WindwalkerMonk
}

func (sefClone *StormEarthAndFirePet) Initialize() {
}

func (ww *WindwalkerMonk) NewSEFPet(name string, swingSpeed float64) *StormEarthAndFirePet {
	sefClone := &StormEarthAndFirePet{
		Pet: core.NewPet(core.PetConfig{
			Name:      name,
			Owner:     &ww.Character,
			BaseStats: stats.Stats{},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				return stats.Stats{
					stats.Stamina:     ownerStats[stats.Stamina] * 0.1,
					stats.AttackPower: ownerStats[stats.AttackPower] * 0,
					stats.HasteRating: ownerStats[stats.HasteRating],

					stats.PhysicalHitPercent: ownerStats[stats.PhysicalHitPercent],
					stats.SpellHitPercent:    ownerStats[stats.PhysicalHitPercent],

					stats.ExpertiseRating: ownerStats[stats.PhysicalHitPercent],

					stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
					stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
				}
			},
			EnabledOnStart:                  false,
			IsGuardian:                      false,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
		owner: ww,
	}

	isDualWielding := swingSpeed == 2.7
	mhWeapon := ww.WeaponFromMainHand(ww.DefaultCritMultiplier())
	mhAvgDPS := mhWeapon.DPS()

	// This number is derived from naked Dummy testing and using multiple
	// other weapons. This was the constant difference between them.
	baseCloneDamage := 266.0

	avgMhDamage := 0.0

	var cloneOhWeapon core.Weapon
	if isDualWielding {
		ohWeapon := ww.WeaponFromOffHand(ww.DefaultCritMultiplier())
		ohAvgDPS := ohWeapon.DPS()

		avgMhDamage = (mhAvgDPS + (ohAvgDPS / 2)) * swingSpeed * core.TernaryFloat64(ohAvgDPS > 0, monk.DualWieldModifier, 1.0)
		cloneOhWeapon = core.Weapon{
			// The clone has a tiny variance in auto attack damage
			BaseDamageMin:  baseCloneDamage - 1 + avgMhDamage,
			BaseDamageMax:  baseCloneDamage + 1 + avgMhDamage,
			SwingSpeed:     swingSpeed,
			CritMultiplier: ww.DefaultCritMultiplier(),
		}
	} else {
		avgMhDamage = mhAvgDPS * swingSpeed
	}

	cloneMhWeapon := core.Weapon{
		BaseDamageMin:  baseCloneDamage + avgMhDamage,
		BaseDamageMax:  baseCloneDamage + avgMhDamage,
		SwingSpeed:     swingSpeed,
		CritMultiplier: ww.DefaultCritMultiplier(),
	}

	sefClone.EnableAutoAttacks(sefClone, core.AutoAttackOptions{
		MainHand:       cloneMhWeapon,
		OffHand:        cloneOhWeapon,
		AutoSwingMelee: true,
	})

	sefClone.OnPetEnable = sefClone.enable
	sefClone.OnPetDisable = sefClone.disable

	ww.AddPet(sefClone)

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
	sefClone.MultiplyMeleeSpeed(sim, sefClone.owner.PseudoStats.MeleeSpeedMultiplier)
	sefClone.PseudoStats.DamageDealtMultiplier = sefClone.owner.PseudoStats.DamageDealtMultiplier

	if sefClone.AutoAttacks.IsDualWielding {
		sefClone.AutoAttacks.DesyncOffHand(sim, sim.CurrentTime)
	}

	sefClone.owner.RegisterOnStanceChanged(func(sim *core.Simulation, _ monk.Stance) {
		sefClone.PseudoStats.DamageDealtMultiplier = sefClone.owner.PseudoStats.DamageDealtMultiplier
	})
}

func (sefClone *StormEarthAndFirePet) disable(sim *core.Simulation) {
}
