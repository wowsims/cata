package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type AncientGuardianPet struct {
	core.Pet

	paladinOwner *Paladin
}

func (guardian *AncientGuardianPet) Initialize() {
}

const PetExpertiseScale = 3.25 * core.ExpertisePerQuarterPercentReduction / core.MeleeHitRatingPerHitChance // 0.8125

func (paladin *Paladin) NewAncientGuardian() *AncientGuardianPet {
	ancientGuardian := &AncientGuardianPet{
		Pet: core.NewPet("Ancient Guardian", &paladin.Character, stats.Stats{
			stats.Stamina: 100,
		}, func(ownerStats stats.Stats) stats.Stats {
			// Draenei Heroic Presence is not included
			hit := ownerStats[stats.MeleeHit]
			if paladin.Race == proto.Race_RaceDraenei {
				hit -= 1 * core.MeleeHitRatingPerHitChance
			}

			return stats.Stats{
				stats.MeleeHit:  hit,
				stats.Expertise: hit * PetExpertiseScale,

				// Taken from combined logs with > 1600 hits, seems to be around 2%
				stats.MeleeCrit: (5 + 1.8) * core.CritRatingPerCritChance,
			}
		}, false, true),
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
	ancientPowerAura := ancientGuardian.RegisterAura(core.Aura{
		Label:    "Ancient Power",
		ActionID: ancientPowerID,
		Duration: core.NeverExpires,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			ancientGuardian.paladinOwner.GetAuraByID(ancientPowerID).AddStack(sim)
		},
	})

	ancientGuardian.EnableAutoAttacks(ancientGuardian, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:     5576,
			BaseDamageMax:     7265,
			SwingSpeed:        2,
			CritMultiplier:    2,
			AttackPowerPerDPS: 0,
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
