package paladin

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type AncientGuardianPet struct {
	core.Pet

	paladinOwner *Paladin
}

func (guardian *AncientGuardianPet) Initialize() {
}

func (paladin *Paladin) NewAncientGuardian() *AncientGuardianPet {
	ancientGuardian := &AncientGuardianPet{
		Pet: core.NewPet(core.PetConfig{
			Name:      "Ancient Guardian",
			Owner:     &paladin.Character,
			BaseStats: stats.Stats{},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				// Draenei Heroic Presence is not included, so inherit HitRating
				// rather than PhysicalHitPercent.
				hitRating := ownerStats[stats.HitRating]
				expertiseRating := ownerStats[stats.ExpertiseRating]
				combined := (hitRating + expertiseRating) * 0.5

				return stats.Stats{
					stats.Armor:               ownerStats[stats.Armor],
					stats.AttackPower:         ownerStats[stats.AttackPower] * 6.1,
					stats.CritRating:          ownerStats[stats.CritRating],
					stats.DodgeRating:         ownerStats[stats.DodgeRating],
					stats.ExpertiseRating:     combined,
					stats.HasteRating:         ownerStats[stats.HasteRating],
					stats.Health:              ownerStats[stats.Health],
					stats.HitRating:           combined,
					stats.ParryRating:         ownerStats[stats.ParryRating],
					stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
					stats.Stamina:             ownerStats[stats.Stamina],
				}
			},
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		paladinOwner: paladin,
	}

	if paladin.Spec == proto.Spec_SpecRetributionPaladin {
		ancientGuardian.registerRetributionVariant()
	} else if paladin.Spec == proto.Spec_SpecHolyPaladin {
		ancientGuardian.registerHolyVariant()
	}

	ancientGuardian.PseudoStats.DamageTakenMultiplier = 0

	paladin.AddPet(ancientGuardian)

	return ancientGuardian
}

func (ancientGuardian *AncientGuardianPet) GetPet() *core.Pet {
	return &ancientGuardian.Pet
}

func (ancientGuardian *AncientGuardianPet) Reset(_ *core.Simulation) {
}

func (ancientGuardian *AncientGuardianPet) ExecuteCustomRotation(sim *core.Simulation) {
	ancientGuardian.WaitUntil(sim, ancientGuardian.AutoAttacks.NextAttackAt())
}

func (ancientGuardian *AncientGuardianPet) registerRetributionVariant() {
	ancientPowerID := core.ActionID{SpellID: 86700}
	ancientPowerAura := core.MakeProcTriggerAura(&ancientGuardian.Unit, core.ProcTrigger{
		Name:     "Ancient Power" + ancientGuardian.Label,
		ActionID: ancientPowerID,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ancientGuardian.paladinOwner.GetAuraByID(ancientPowerID).AddStack(sim)
		},
	})

	baseDamage := ancientGuardian.paladinOwner.CalcScalingSpellDmg(6.1)
	ancientGuardian.EnableAutoAttacks(ancientGuardian, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  baseDamage,
			BaseDamageMax:  baseDamage,
			SwingSpeed:     2,
			CritMultiplier: ancientGuardian.DefaultCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	ancientGuardian.OnPetEnable = func(sim *core.Simulation) {
		ancientPowerAura.Activate(sim)
	}
	ancientGuardian.OnPetDisable = func(sim *core.Simulation) {
		ancientPowerAura.Deactivate(sim)
	}
}

func (ancientGuardian *AncientGuardianPet) registerHolyVariant() {
	// TODO: Implement this when Holy spec is in place

	// // Heals the target of your last single-target heal and allies within 10 yards of the target.
	// lightOfTheAncientKings := ancientGuardian.RegisterSpell(core.SpellConfig{
	// 	ActionID:    core.ActionID{SpellID: 86678},
	// 	SpellSchool: core.SpellSchoolHoly,
	// 	Flags:       core.SpellFlagHelpful,
	// 	ProcMask:    core.ProcMaskEmpty,
	//
	// 	MaxRange: 100,
	//
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			NonEmpty: true,
	// 		},
	// 	},
	// })
}
