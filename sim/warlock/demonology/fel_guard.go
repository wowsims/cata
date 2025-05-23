package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
)

func (demo *DemonologyWarlock) registerFelguard() *warlock.WarlockPet {
	name := proto.WarlockOptions_Summon_name[int32(proto.WarlockOptions_Felguard)]
	enabledOnStart := proto.WarlockOptions_Felguard == demo.Options.Summon
	pet := demo.RegisterPet(proto.WarlockOptions_Felguard, 2, 3.5, name, enabledOnStart)
	registerLegionStrikeSpell(pet, demo)
	return pet
}

func registerLegionStrikeSpell(pet *warlock.WarlockPet, demo *DemonologyWarlock) {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30213},
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
			baseDmg := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * 1.3
			baseDmg /= float64(sim.Environment.GetNumTargets())

			for _, target := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, target, baseDmg, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			demo.DemonicFury.Gain(12, core.ActionID{SpellID: 30213}, sim)
		},
	}))
}
