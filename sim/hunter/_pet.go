package hunter

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type HunterPet struct {
	core.Pet

	config PetConfig

	hunterOwner *Hunter

	CobraStrikesAura     *core.Aura
	KillCommandAura      *core.Aura
	FrenzyStacksSnapshot float64
	FrenzyAura           *core.Aura

	specialAbility *core.Spell
	KillCommand    *core.Spell
	focusDump      *core.Spell
	exoticAbility  *core.Spell

	uptimePercent    float64
	wolverineBite    *core.Spell
	frostStormBreath *core.Spell
	hasOwnerCooldown bool
}

func (hunter *Hunter) NewHunterPet() *HunterPet {
	if hunter.Options.PetType == proto.HunterOptions_PetNone {
		return nil
	}
	if hunter.Options.PetUptime <= 0 {
		return nil
	}
	petConfig := PetConfigs[hunter.Options.PetType]

	hp := &HunterPet{
		Pet: core.NewPet(core.PetConfig{
			Name:            petConfig.Name,
			Owner:           &hunter.Character,
			BaseStats:       hunterPetBaseStats,
			StatInheritance: hunter.makeStatInheritance(),
			EnabledOnStart:  true,
			IsGuardian:      false,
		}),
		config:      petConfig,
		hunterOwner: hunter,
		//hasOwnerCooldown: petConfig.SpecialAbility == FuriousHowl || petConfig.SpecialAbility == SavageRend,
	}

	//Todo: Verify this
	// base_focus_regen_per_second  = ( 24.5 / 4.0 );
	// base_focus_regen_per_second *= 1.0 + o -> talents.bestial_discipline -> effect1().percent();
	baseFocusPerSecond := 4.0 // As observed on logs
	baseFocusPerSecond *= 1.0 + (0.10 * float64(hunter.Talents.BestialDiscipline))

	WHFocusIncreaseMod := hp.AddDynamicMod(core.SpellModConfig{
		Kind:     core.SpellMod_PowerCost_Pct,
		ProcMask: core.ProcMaskMeleeMHSpecial,
		IntValue: hp.Talents().WildHunt * 50,
	})

	WHDamageMod := hp.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ProcMask:   core.ProcMaskMeleeMHSpecial,
		FloatValue: float64(hp.Talents().WildHunt) * 0.6,
	})

	hp.EnableFocusBar(100+(float64(hunter.Talents.KindredSpirits)*5), baseFocusPerSecond, false, func(sim *core.Simulation, focus float64) {
		if hp.Talents().WildHunt > 0 {
			if focus >= 50 {
				WHFocusIncreaseMod.Activate()
				WHDamageMod.Activate()
			} else {
				WHFocusIncreaseMod.Deactivate()
				WHDamageMod.Deactivate()
			}
		}
	})

	atkSpd := 2.0
	// Todo: Change for Cataclysm
	hp.EnableAutoAttacks(hp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  73,
			BaseDamageMax:  110,
			CritMultiplier: 2,
			SwingSpeed:     atkSpd,
		},
		AutoSwingMelee: true,
	})

	// Happiness
	// Todo:
	hp.PseudoStats.DamageDealtMultiplier *= 1.25

	// Pet family bonus is now the same for all pets.
	//Todo: Should this stay?
	hp.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.05

	hp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	hp.AddStatDependency(stats.Strength, stats.RangedAttackPower, 2)
	hp.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, 1/324.72)

	hunter.AddPet(hp)

	return hp
}

func (hp *HunterPet) GetPet() *core.Pet {
	return &hp.Pet
}

func (hp *HunterPet) Talents() *proto.HunterPetTalents {
	if talents := hp.hunterOwner.Options.PetTalents; talents != nil {
		return talents
	}
	return &proto.HunterPetTalents{}
}

func (hp *HunterPet) Initialize() {
	hp.specialAbility = hp.NewPetAbility(hp.config.SpecialAbility, true)
	hp.focusDump = hp.NewPetAbility(hp.config.FocusDump, false)
	hp.exoticAbility = hp.NewPetAbility(hp.config.ExoticAbility, false)
	hp.KillCommand = hp.RegisterKillCommandSpell()
}

func (hp *HunterPet) Reset(_ *core.Simulation) {
	hp.uptimePercent = min(1, max(0, hp.hunterOwner.Options.PetUptime))
}

func (hp *HunterPet) ExecuteCustomRotation(sim *core.Simulation) {
	percentRemaining := sim.GetRemainingDurationPercent()
	if percentRemaining < 1.0-hp.uptimePercent { // once fight is % completed, disable pet.
		hp.Disable(sim)
		return
	}

	if hp.hasOwnerCooldown && hp.CurrentFocus() < 50 {
		// When a major ability (Furious Howl or Savage Rend) is ready, pool enough
		// energy to use on-demand.
		return
	}

	target := hp.CurrentTarget

	if hp.frostStormBreath != nil && hp.frostStormBreath.CanCast(sim, target) && len(sim.Encounter.TargetUnits) > 4 {
		hp.frostStormBreath.Cast(sim, target)
	}

	if hp.wolverineBite.CanCast(sim, target) {
		hp.wolverineBite.Cast(sim, target)
	}

	if hp.focusDump == nil {
		hp.specialAbility.Cast(sim, target)
		return
	}
	if hp.specialAbility == nil {
		hp.focusDump.Cast(sim, target)
		return
	}

	if hp.config.RandomSelection {
		if sim.RandomFloat("Hunter Pet Ability") < 0.5 {
			_ = hp.specialAbility.Cast(sim, target) || hp.focusDump.Cast(sim, target)
		} else {
			_ = hp.focusDump.Cast(sim, target) || hp.specialAbility.Cast(sim, target)
		}
	} else {
		_ = hp.specialAbility.Cast(sim, target) || hp.focusDump.Cast(sim, target)
	}
}

var hunterPetBaseStats = stats.Stats{
	stats.Agility:     438,
	stats.Strength:    476,
	stats.AttackPower: -20, // Apparently pets and warriors have a AP penalty.

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.PhysicalCritPercent: 3.2 + 1.8,
}

const PetExpertiseRatingScale = 3.25 * core.PhysicalHitRatingPerHitPercent

func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {

	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:           ownerStats[stats.Stamina] * 0.3,
			stats.Armor:             ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:       ownerStats[stats.RangedAttackPower] * 0.425,
			stats.RangedAttackPower: ownerStats[stats.RangedAttackPower],

			stats.PhysicalHitPercent: ownerStats[stats.PhysicalHitPercent],
			stats.ExpertiseRating:    ownerStats[stats.PhysicalHitPercent] * PetExpertiseRatingScale,
			stats.SpellHitPercent:    ownerStats[stats.PhysicalHitPercent],

			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
			stats.SpellCritPercent:    ownerStats[stats.PhysicalCritPercent],

			stats.HasteRating: ownerStats[stats.HasteRating],
		}
	}
}

type PetConfig struct {
	Name string

	SpecialAbility PetAbilityType
	FocusDump      PetAbilityType
	ExoticAbility  PetAbilityType

	// Randomly select between abilities instead of using a prio.
	RandomSelection bool
}

var PetConfigs = map[proto.HunterOptions_PetType]PetConfig{
	proto.HunterOptions_Bat: {
		Name:      "Bat",
		FocusDump: Claw,
	},
	proto.HunterOptions_Bear: {
		Name:           "Bear",
		FocusDump:      Claw,
		SpecialAbility: DemoralizingScreech,
	},
	proto.HunterOptions_BirdOfPrey: {
		Name:           "Bird of Prey",
		FocusDump:      Claw,
		SpecialAbility: DemoralizingScreech,
	},
	proto.HunterOptions_Boar: {
		Name:           "Boar",
		FocusDump:      Bite,
		SpecialAbility: Stampede,
	},
	proto.HunterOptions_CarrionBird: {
		Name:      "Carrion Bird",
		FocusDump: Bite,
	},
	proto.HunterOptions_Cat: {
		Name:      "Cat",
		FocusDump: Claw,
	},
	proto.HunterOptions_Chimaera: {
		Name:          "Chimaera",
		FocusDump:     Bite,
		ExoticAbility: FrostStormBreath,
	},
	proto.HunterOptions_CoreHound: {
		Name:      "Core Hound",
		FocusDump: Bite,
	},
	proto.HunterOptions_Crab: {
		Name:      "Crab",
		FocusDump: Claw,
	},
	proto.HunterOptions_Crocolisk: {
		Name: "Crocolisk",
		//SpecialAbility: BadAttitude,
		FocusDump: Bite,
	},
	proto.HunterOptions_Devilsaur: {
		Name:      "Devilsaur",
		FocusDump: Bite,
	},
	proto.HunterOptions_Fox: {
		Name:      "Fox",
		FocusDump: Claw,
	},
	proto.HunterOptions_ShaleSpider: {
		Name:      "Shale Spider",
		FocusDump: Bite,
	},
	proto.HunterOptions_Dragonhawk: {
		Name:           "Dragonhawk",
		FocusDump:      Bite,
		SpecialAbility: FireBreath,
	},
	proto.HunterOptions_Gorilla: {
		Name: "Gorilla",
		//SpecialAbility: Pummel,
		FocusDump: Smack,
	},
	proto.HunterOptions_Hyena: {
		Name:           "Hyena",
		SpecialAbility: Stampede,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Moth: {
		Name:      "Moth",
		FocusDump: Smack,
	},
	proto.HunterOptions_NetherRay: {
		Name:      "Nether Ray",
		FocusDump: Bite,
	},
	proto.HunterOptions_Raptor: {
		Name:           "Raptor",
		FocusDump:      Claw,
		SpecialAbility: CorrosiveSpit,
	},
	proto.HunterOptions_Ravager: {
		Name:           "Ravager",
		FocusDump:      Bite,
		SpecialAbility: AcidSpit,
	},
	proto.HunterOptions_Rhino: {
		Name:           "Rhino",
		FocusDump:      Bite,
		SpecialAbility: Stampede,
	},
	proto.HunterOptions_Scorpid: {
		Name:      "Scorpid",
		FocusDump: Bite,
	},
	proto.HunterOptions_Serpent: {
		Name:           "Serpent",
		FocusDump:      Bite,
		SpecialAbility: CorrosiveSpit,
	},
	proto.HunterOptions_Silithid: {
		Name:      "Silithid",
		FocusDump: Claw,
	},
	proto.HunterOptions_Spider: {
		Name: "Spider",
		//SpecialAbility:   Web,
		FocusDump: Bite,
	},
	proto.HunterOptions_SpiritBeast: {
		Name:      "Spirit Beast",
		FocusDump: Claw,
	},
	proto.HunterOptions_SporeBat: {
		Name:      "Spore Bat",
		FocusDump: Smack,
	},
	proto.HunterOptions_Tallstrider: {
		Name: "Tallstrider",
		//SpecialAbility:   DustCloud,
		FocusDump: Claw,
	},
	proto.HunterOptions_Turtle: {
		Name: "Turtle",
		//SpecialAbility: ShellShield,
		FocusDump: Bite,
	},
	proto.HunterOptions_WarpStalker: {
		Name: "Warp Stalker",
		//SpecialAbility:   Warp,
		FocusDump: Bite,
	},
	proto.HunterOptions_Wasp: {
		Name:      "Wasp",
		FocusDump: Smack,
	},
	proto.HunterOptions_WindSerpent: {
		Name:           "Wind Serpent",
		FocusDump:      Bite,
		SpecialAbility: FireBreath,
	},
	proto.HunterOptions_Wolf: {
		Name:      "Wolf",
		FocusDump: Bite,
	},
	proto.HunterOptions_Worm: {
		Name:           "Worm",
		FocusDump:      Bite,
		SpecialAbility: AcidSpit,
	},
}
