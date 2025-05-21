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

const PetExpertiseScale = 3.25 * core.ExpertisePerQuarterPercentReduction / core.PhysicalHitRatingPerHitPercent // 0.8125

func (paladin *Paladin) NewAncientGuardian() *AncientGuardianPet {
	ancientGuardian := &AncientGuardianPet{
		Pet: core.NewPet(core.PetConfig{
			Name:  "Ancient Guardian",
			Owner: &paladin.Character,
			BaseStats: stats.Stats{
				stats.Stamina: 100,

				// Taken from combined logs with > 1600 hits, seems to
				// be around 2% final Crit chance after the 4.8%
				// suppression from boss level mobs.
				stats.PhysicalCritPercent: 6.8,
			},
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				// Draenei Heroic Presence is not included, so inherit HitRating
				// rather than PhysicalHitPercent.
				ownerHitRating := ownerStats[stats.HitRating]

				return stats.Stats{
					stats.HitRating:       ownerHitRating,
					stats.ExpertiseRating: ownerHitRating * PetExpertiseScale,
				}
			},
			EnabledOnStart: false,
			IsGuardian:     true,
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
	ancientPowerAura := ancientGuardian.RegisterAura(core.Aura{
		Label:    "Ancient Power" + ancientGuardian.Label,
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
