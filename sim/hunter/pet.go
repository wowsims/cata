package hunter

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
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
	focusDump      *core.Spell

	uptimePercent    float64
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
		Pet:         core.NewPet(petConfig.Name, &hunter.Character, hunterPetBaseStats, hunter.makeStatInheritance(), true, false),
		config:      petConfig,
		hunterOwner: hunter,

		//hasOwnerCooldown: petConfig.SpecialAbility == FuriousHowl || petConfig.SpecialAbility == SavageRend,
	}

	//Todo: Verify this
	// base_focus_regen_per_second  = ( 24.5 / 4.0 );
	// base_focus_regen_per_second *= 1.0 + o -> talents.bestial_discipline -> effect1().percent();
	baseFocusPerSecond := 24.5 / 4.0
	baseFocusPerSecond *= 1.0 + (0.10 * float64(hunter.Talents.BestialDiscipline))

	WHFocusIncreaseMod := hp.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Flat,
		ProcMask:   core.ProcMaskMeleeMHSpecial,
		FloatValue: float64(hp.Talents().WildHunt) * 12.5,
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

	atkSpd := 2 / (1 + 0.05*float64(hp.Talents().SerpentSwiftness))
	// Todo: Change for Cataclysm
	hp.EnableAutoAttacks(hp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  73,
			BaseDamageMax:  110,
			SwingSpeed:     atkSpd,
			CritMultiplier: 2,
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
	hp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/324.72)

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
	stats.MeleeCrit: (3.2 + 1.8) * core.CritRatingPerCritChance,
}

const PetExpertiseScale = 3.25

func (hunter *Hunter) makeStatInheritance() core.PetStatInheritance {

	return func(ownerStats stats.Stats) stats.Stats {
		// EJ posts claim this value is passed through math.Floor, but in-game testing
		// shows pets benefit from each point of owner hit rating in WotLK Classic.
		// https://web.archive.org/web/20120112003252/http://elitistjerks.com/f80/t100099-demonology_releasing_demon_you

		return stats.Stats{
			stats.Stamina:           ownerStats[stats.Stamina] * 0.3,
			stats.Armor:             ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:       ownerStats[stats.RangedAttackPower] * 0.425,
			stats.RangedAttackPower: ownerStats[stats.RangedAttackPower] * 0.40,

			stats.MeleeHit:  ownerStats[stats.MeleeHit],
			stats.Expertise: ownerStats[stats.MeleeHit] * PetExpertiseScale,
			stats.SpellHit:  ownerStats[stats.MeleeHit],

			stats.MeleeCrit: ownerStats[stats.MeleeCrit],
			stats.SpellCrit: ownerStats[stats.MeleeCrit],

			stats.MeleeHaste: ownerStats[stats.MeleeHaste],
			stats.SpellHaste: ownerStats[stats.MeleeHaste],
		}
	}
}

type PetConfig struct {
	Name string

	SpecialAbility PetAbilityType
	FocusDump      PetAbilityType

	// Randomly select between abilities instead of using a prio.
	RandomSelection bool
}

// Abilities reference: https://wotlk.wowhead.com/hunter-pets
// https://wotlk.wowhead.com/guides/hunter-dps-best-pets-taming-loyalty-burning-crusade-classic
var PetConfigs = map[proto.HunterOptions_PetType]PetConfig{
	proto.HunterOptions_Bat: {
		Name:           "Bat",
		SpecialAbility: SonicBlast,
		FocusDump:      Claw,
	},
	proto.HunterOptions_Bear: {
		Name:           "Bear",
		SpecialAbility: Swipe,
		FocusDump:      Claw,
	},
	proto.HunterOptions_BirdOfPrey: {
		Name:           "Bird of Prey",
		SpecialAbility: Snatch,
		FocusDump:      Claw,
	},
	proto.HunterOptions_Boar: {
		Name:           "Boar",
		SpecialAbility: Gore,
		FocusDump:      Bite,
	},
	proto.HunterOptions_CarrionBird: {
		Name:           "Carrion Bird",
		SpecialAbility: DemoralizingScreech,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Cat: {
		Name:           "Cat",
		SpecialAbility: Rake,
		FocusDump:      Claw,
	},
	proto.HunterOptions_Chimaera: {
		Name:           "Chimaera",
		SpecialAbility: FroststormBreath,
		FocusDump:      Bite,
	},
	proto.HunterOptions_CoreHound: {
		Name:           "Core Hound",
		SpecialAbility: LavaBreath,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Crab: {
		Name:           "Crab",
		SpecialAbility: Pin,
		FocusDump:      Claw,
	},
	proto.HunterOptions_Crocolisk: {
		Name: "Crocolisk",
		//SpecialAbility: BadAttitude,
		FocusDump: Bite,
	},
	proto.HunterOptions_Devilsaur: {
		Name:           "Devilsaur",
		SpecialAbility: MonstrousBite,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Dragonhawk: {
		Name:           "Dragonhawk",
		SpecialAbility: FireBreath,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Gorilla: {
		Name: "Gorilla",
		//SpecialAbility: Pummel,
		FocusDump: Smack,
	},
	proto.HunterOptions_Hyena: {
		Name:           "Hyena",
		SpecialAbility: TendonRip,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Moth: {
		Name: "Moth",
		//SpecialAbility:   SerentiyDust,
		FocusDump: Smack,
	},
	proto.HunterOptions_NetherRay: {
		Name:           "Nether Ray",
		SpecialAbility: NetherShock,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Raptor: {
		Name:           "Raptor",
		SpecialAbility: SavageRend,
		FocusDump:      Claw,
	},
	proto.HunterOptions_Ravager: {
		Name:           "Ravager",
		SpecialAbility: Ravage,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Rhino: {
		Name:           "Rhino",
		SpecialAbility: Stampede,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Scorpid: {
		Name:           "Scorpid",
		SpecialAbility: ScorpidPoison,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Serpent: {
		Name:           "Serpent",
		SpecialAbility: PoisonSpit,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Silithid: {
		Name:           "Silithid",
		SpecialAbility: VenomWebSpray,
		FocusDump:      Claw,
	},
	proto.HunterOptions_Spider: {
		Name: "Spider",
		//SpecialAbility:   Web,
		FocusDump: Bite,
	},
	proto.HunterOptions_SpiritBeast: {
		Name:           "Spirit Beast",
		SpecialAbility: SpiritStrike,
		FocusDump:      Claw,
	},
	proto.HunterOptions_SporeBat: {
		Name:           "Spore Bat",
		SpecialAbility: SporeCloud,
		FocusDump:      Smack,
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
		Name:           "Wasp",
		SpecialAbility: Sting,
		FocusDump:      Smack,
	},
	proto.HunterOptions_WindSerpent: {
		Name:           "Wind Serpent",
		SpecialAbility: LightningBreath,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Wolf: {
		Name:           "Wolf",
		SpecialAbility: FuriousHowl,
		FocusDump:      Bite,
	},
	proto.HunterOptions_Worm: {
		Name:           "Worm",
		SpecialAbility: AcidSpit,
		FocusDump:      Bite,
	},
}
