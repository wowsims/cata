package demonology

import (
	"slices"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
)

func (demo *DemonologyWarlock) registerFelguard() *warlock.WarlockPet {
	name := proto.WarlockOptions_Summon_name[int32(proto.WarlockOptions_Felguard)]
	enabledOnStart := proto.WarlockOptions_Felguard == demo.Options.Summon
	return demo.registerFelguardWithName(name, enabledOnStart, false, false)
}

func (demo *DemonologyWarlock) registerFelguardWithName(name string, enabledOnStart bool, autoCastFelstorm bool, isGuardian bool) *warlock.WarlockPet {
	pet := demo.RegisterPet(proto.WarlockOptions_Felguard, 2, 3.5, name, enabledOnStart, isGuardian)
	registerLegionStrikeSpell(pet, demo)
	felStorm := registerFelstorm(pet, demo, autoCastFelstorm)
	pet.MinEnergy = 120

	if !isGuardian {
		demo.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 89751},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagAPL | core.SpellFlagNoMetrics,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    demo.NewTimer(),
					Duration: time.Second * 45,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				pet.AutoCastAbilities = slices.Insert(pet.AutoCastAbilities, 0, felStorm)
			},
		})
	}

	return pet
}

var legionStrikePetAction = core.ActionID{SpellID: 30213}

func registerLegionStrikeSpell(pet *warlock.WarlockPet, demo *DemonologyWarlock) {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       legionStrikePetAction,
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

			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: time.Millisecond * 1300, // add small cooldown to allow for proper rotation of abilities
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * 1.3
			baseDmg /= float64(sim.Environment.GetNumTargets())

			for _, target := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, target, baseDmg, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			demo.DemonicFury.Gain(sim, 12, core.ActionID{SpellID: 30213})
		},
	}))
}

func registerFelstorm(pet *warlock.WarlockPet, _ *DemonologyWarlock, autoCast bool) *core.Spell {
	felStorm := pet.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 89751},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagAoE | core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagChanneled,
		EnergyCost: core.EnergyCostOptions{
			Cost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: time.Second * 45,
			},
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   2,
		Dot: core.DotConfig{
			IsAOE:         true,
			Aura:          core.Aura{Label: "Felstorm"},
			NumberOfTicks: 6,
			TickLength:    time.Second,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := dot.Spell.Unit.MHWeaponDamage(sim, dot.Spell.MeleeAttackPower()) + dot.Spell.Unit.OHWeaponDamage(sim, dot.Spell.MeleeAttackPower())
				for _, enemy := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, enemy, baseDamage, dot.Spell.OutcomeMeleeSpecialBlockAndCritNoHitCounter)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
			spell.AOEDot().TickOnce(sim)
			pet.AutoAttacks.DelayMeleeBy(sim, spell.AOEDot().BaseDuration())

			// remove from auto cast again to trigger it once
			if !pet.IsGuardian() {
				pet.AutoCastAbilities = pet.AutoCastAbilities[1:]
			}
		},
	})

	if autoCast {
		pet.AutoCastAbilities = append(pet.AutoCastAbilities, felStorm)
	}

	return felStorm
}
