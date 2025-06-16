package demonology

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (demonlogy *DemonologyWarlock) registerGrimoireOfSupremacy() {
	if !demonlogy.Talents.GrimoireOfSupremacy {
		return
	}

	// Pimp my demo pet
	felGuard := demonlogy.Felguard
	demonlogy.Felguard.PseudoStats.DamageDealtMultiplier *= 1.2
	felGuard.Name = "Wrathguard"
	felGuard.Label = fmt.Sprintf("%s - %s", demonlogy.Label, "Wrathguard")

	// Now dualwield with 1.5x less base damage
	weaponConfig := warlock.ScaledAutoAttackConfig(2)
	weaponConfig.MainHand.BaseDamageMax /= 1.5
	weaponConfig.MainHand.BaseDamageMin /= 1.5
	weaponConfig.OffHand = weaponConfig.MainHand

	felGuard.EnableAutoAttacks(felGuard, *weaponConfig)
	felGuard.ChangeStatInheritance(demonlogy.SimplePetStatInheritanceWithScale(2 + 1.0/3.0))
	felGuard.PseudoStats.DisableDWMissPenalty = true

	mortalCleave := felGuard.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 115625},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellFelGuardLegionStrike,

		EnergyCost: core.EnergyCostOptions{
			Cost: 60,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 1,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * 1.95
			baseDmg /= float64(sim.Environment.GetNumTargets())

			for _, target := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, target, baseDmg, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			demonlogy.DemonicFury.Gain(sim, 12, core.ActionID{SpellID: 30213})
		},
	})

	felGuard.AutoCastAbilities = []*core.Spell{mortalCleave}
}

func (demonology *DemonologyWarlock) registerGrimoireOfService() {
	if !demonology.Talents.GrimoireOfService {
		return
	}

	felGuard := demonology.registerFelguardWithName("Grimoire: Felguard", false, true, true)
	felGuard.MinEnergy = 90
	demonology.BuildAndRegisterSummonSpell(111898, felGuard)
}

func (demonology *DemonologyWarlock) registerGrimoireOfSacrifice() {
	if !demonology.Talents.GrimoireOfSacrifice {
		return
	}

	// rest handle din talents.go of warlock
	for _, pet := range demonology.WildImps {
		pet.Fireball.DamageMultiplier *= 1.25
	}
}
