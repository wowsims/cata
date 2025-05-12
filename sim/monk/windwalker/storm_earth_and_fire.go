package windwalker

import (
	"fmt"
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
	sefAura := ww.RegisterAura(core.Aura{
		Label:    "Storm, Earth, and Fire",
		ActionID: core.ActionID{SpellID: 137639},
		Duration: core.NeverExpires,
		// Casts copy
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			copySpell := ww.SefController.GetCopySpell(spell.ActionID)
			if copySpell == nil {
				return
			}

			CopySpellMultipliers(spell, copySpell, ww.CurrentTarget)

			copySpell.Cast(sim, ww.CurrentTarget)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ww.SefController.Activate(sim, ww.CurrentTarget)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ww.SefController.Deactivate(sim)
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

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			sefAura.Activate(sim)
		},
	})
}

type StormEarthAndFireController struct {
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

func (controller *StormEarthAndFireController) Activate(sim *core.Simulation, target *core.Unit) {
	targetUnixIndex := target.UnitIndex
	activeSef := controller.sefTargets[targetUnixIndex]
	// If the target already has an active clone, disable it
	if activeSef != nil {
		activeSef.Disable(sim)
		return
	}
	// Pick a random pet to spawn
	petIndex := int32(math.Round(sim.Roll(0, 2)))
	pet := controller.pets[petIndex]
	fmt.Println("Activating SEF for target", len(controller.sefTargets), targetUnixIndex, petIndex)
	// petUnitIndex := pet.UnitIndex
	pet.EnableWithStartAttackDelay(sim, pet, core.DurationFromSeconds(sim.RollWithLabel(2, 2.3, "SEF Spawn Delay")))
	pet.AutoAttacks.SwapTarget(sim, target)

}

func (controller *StormEarthAndFireController) Deactivate(sim *core.Simulation) {
	for _, pet := range controller.pets {
		pet.Disable(sim)
	}
}

func (ww *WindwalkerMonk) registerSEFPets() {
	ww.SefController = &StormEarthAndFireController{
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
		Pet: core.NewPet(name, &ww.Character, stats.Stats{}, func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{
				stats.Stamina:     ownerStats[stats.Stamina] * 0.1,
				stats.AttackPower: ownerStats[stats.AttackPower],
				stats.HasteRating: ownerStats[stats.HasteRating],

				stats.PhysicalHitPercent: ownerStats[stats.PhysicalHitPercent],
				stats.SpellHitPercent:    ownerStats[stats.PhysicalHitPercent],

				stats.ExpertiseRating: ownerStats[stats.PhysicalHitPercent],

				stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
				stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
			}
		}, false, false),
		owner: ww,
	}

	isDualWielding := swingSpeed == 2.7
	mhWeapon := ww.WeaponFromMainHand(ww.DefaultCritMultiplier())
	mhAvgDPS := mhWeapon.AverageDamage()

	// This number is derived from naked Dummy testing and using multiple
	// other weapons. This was the constant difference between them.
	baseCloneDamage := 161.0

	avgMhDamage := 0.0
	avgOhDamage := 0.0

	var cloneOhWeapon core.Weapon
	if isDualWielding {
		ohWeapon := ww.WeaponFromOffHand(ww.DefaultCritMultiplier())
		ohAvgDPS := ohWeapon.AverageDamage()

		avgMhDamage = (mhAvgDPS + (ohAvgDPS / 2)) * swingSpeed * monk.DualWieldModifier
		avgOhDamage = avgMhDamage / 2
		cloneOhWeapon = core.Weapon{
			// The clone has a tiny variance in auto attack damage
			BaseDamageMin:  baseCloneDamage - 1 + avgOhDamage,
			BaseDamageMax:  baseCloneDamage + 1 + avgOhDamage,
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
	sefClone.owner.RegisterOnStanceChanged(func(sim *core.Simulation, _ monk.Stance) {
		sefClone.PseudoStats.DamageDealtMultiplier = sefClone.owner.PseudoStats.DamageDealtMultiplier
	})

	sefClone.EnableDynamicMeleeSpeed(func(amount float64) {
		sefClone.MultiplyCastSpeed(amount)
		sefClone.MultiplyMeleeSpeed(sim, amount)
	})
}

func (sefClone *StormEarthAndFirePet) disable(sim *core.Simulation) {
}
