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

func (hunter *Hunter) NewStampedePet() *HunterPet {
	if hunter.Options.PetType == proto.HunterOptions_PetNone {
		return nil
	}
	if hunter.Options.PetUptime <= 0 {
		return nil
	}
	petConfig := DefaultPetConfigs[hunter.Options.PetType]

	stampedePet := &HunterPet{
		Pet:         core.NewPet(petConfig.Name, &hunter.Character, hunterPetBaseStats, hunter.makeStatInheritance(), true, false),
		config:      petConfig,
		hunterOwner: hunter,

		//hasOwnerCooldown: petConfig.SpecialAbility == FuriousHowl || petConfig.SpecialAbility == SavageRend,
	}

	return stampedePet
}
func (hunter *Hunter) NewDireBeastPet() *HunterPet {
	if hunter.Options.PetType == proto.HunterOptions_PetNone {
		return nil
	}
	if hunter.Options.PetUptime <= 0 {
		return nil
	}
	petConfig := DefaultPetConfigs[hunter.Options.PetType]

	stampedePet := &HunterPet{
		Pet:         core.NewPet(petConfig.Name, &hunter.Character, hunterPetBaseStats, hunter.makeStatInheritance(), true, false),
		config:      petConfig,
		hunterOwner: hunter,

		//hasOwnerCooldown: petConfig.SpecialAbility == FuriousHowl || petConfig.SpecialAbility == SavageRend,
	}

	return stampedePet
}
func (hunter *Hunter) NewHunterPet() *HunterPet {
	if hunter.Options.PetType == proto.HunterOptions_PetNone {
		return nil
	}
	if hunter.Options.PetUptime <= 0 {
		return nil
	}
	petConfig := DefaultPetConfigs[hunter.Options.PetType]

	hp := &HunterPet{
		Pet:         core.NewPet(petConfig.Name, &hunter.Character, hunterPetBaseStats, hunter.makeStatInheritance(), true, false),
		config:      petConfig,
		hunterOwner: hunter,

		//hasOwnerCooldown: petConfig.SpecialAbility == FuriousHowl || petConfig.SpecialAbility == SavageRend,
	}

	//Todo: Verify this
	// base_focus_regen_per_second  = ( 24.5 / 4.0 );
	// base_focus_regen_per_second *= 1.0 + o -> talents.bestial_discipline -> effect1().percent();
	baseFocusPerSecond := 4.0 // As observed on logs

	// WHFocusIncreaseMod := hp.AddDynamicMod(core.SpellModConfig{
	// 	Kind:     core.SpellMod_PowerCost_Pct,
	// 	ProcMask: core.ProcMaskMeleeMHSpecial,
	// 	IntValue: hp.Talents().WildHunt * 50,
	// })

	// WHDamageMod := hp.AddDynamicMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_DamageDone_Flat,
	// 	ProcMask:   core.ProcMaskMeleeMHSpecial,
	// 	FloatValue: float64(hp.Talents().WildHunt) * 0.6,
	// })

	hp.EnableFocusBar(100+(core.TernaryFloat64(hp.hunterOwner.Spec == proto.Spec_SpecBeastMasteryHunter, 20, 0)), baseFocusPerSecond, false, func(sim *core.Simulation, focus float64) {

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

func (hp *HunterPet) Initialize() {
	cfg := DefaultPetConfigs[hp.hunterOwner.Options.PetType]

	// Primary active ability (often a cooldown)
	if cfg.SpecialAbility != Unknown {
		hp.specialAbility = hp.NewPetAbility(cfg.SpecialAbility, true)
	}

	// Focus-generating basic attack
	if cfg.FocusDump != Unknown {
		hp.focusDump = hp.NewPetAbility(cfg.FocusDump, false)
	}

	// Optional exotic ability
	if cfg.ExoticAbility != Unknown {
		hp.exoticAbility = hp.NewPetAbility(cfg.ExoticAbility, false)
	}
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

var DefaultPetConfigs = [...]PetConfig{
	proto.HunterOptions_PetNone:      {},
	proto.HunterOptions_Bat:          {Name: "Bat", FocusDump: Smack, SpecialAbility: SonicBlast},
	proto.HunterOptions_Bear:         {Name: "Bear", FocusDump: Claw, SpecialAbility: DemoralizingRoar},
	proto.HunterOptions_BirdOfPrey:   {Name: "Bird of Prey", FocusDump: Claw, SpecialAbility: Snatch},
	proto.HunterOptions_Boar:         {Name: "Boar", FocusDump: Bite, SpecialAbility: Gore},
	proto.HunterOptions_CarrionBird:  {Name: "Carrion Bird", FocusDump: Bite, SpecialAbility: DemoralizingScreech},
	proto.HunterOptions_Cat:          {Name: "Cat", FocusDump: Claw},
	proto.HunterOptions_Chimaera:     {Name: "Chimaera", FocusDump: Bite, ExoticAbility: FroststormBreathAoE},
	proto.HunterOptions_CoreHound:    {Name: "Core Hound", FocusDump: Bite, ExoticAbility: LavaBreath},
	proto.HunterOptions_Crab:         {Name: "Crab", FocusDump: Claw, SpecialAbility: Pin},
	proto.HunterOptions_Crocolisk:    {Name: "Crocolisk", FocusDump: Bite},
	proto.HunterOptions_Devilsaur:    {Name: "Devilsaur", FocusDump: Bite, SpecialAbility: TerrifyingRoar, ExoticAbility: MonstrousBite},
	proto.HunterOptions_Dragonhawk:   {Name: "Dragonhawk", FocusDump: Bite, SpecialAbility: FireBreathDebuff},
	proto.HunterOptions_Fox:          {Name: "Fox", FocusDump: Bite, SpecialAbility: TailSpin},
	proto.HunterOptions_Gorilla:      {Name: "Gorilla", FocusDump: Smack, SpecialAbility: Pummel},
	proto.HunterOptions_Hyena:        {Name: "Hyena", FocusDump: Bite, SpecialAbility: CacklingHowl},
	proto.HunterOptions_Moth:         {Name: "Moth", FocusDump: Smack, SpecialAbility: SerenityDust},
	proto.HunterOptions_NetherRay:    {Name: "Nether Ray", FocusDump: Bite, SpecialAbility: NetherShock},
	proto.HunterOptions_Raptor:       {Name: "Raptor", FocusDump: Claw, SpecialAbility: TearArmor},
	proto.HunterOptions_Ravager:      {Name: "Ravager", FocusDump: Bite, SpecialAbility: Ravage},
	proto.HunterOptions_Rhino:        {Name: "Rhino", FocusDump: Bite, SpecialAbility: StampedeDebuff, ExoticAbility: HornToss},
	proto.HunterOptions_Scorpid:      {Name: "Scorpid", FocusDump: Bite, SpecialAbility: Clench},
	proto.HunterOptions_Serpent:      {Name: "Serpent", FocusDump: Bite, SpecialAbility: SerpentsSwiftness},
	proto.HunterOptions_Silithid:     {Name: "Silithid", FocusDump: Claw, SpecialAbility: QirajiFortitude, ExoticAbility: VenomWebSpray},
	proto.HunterOptions_Spider:       {Name: "Spider", FocusDump: Bite, SpecialAbility: Web},
	proto.HunterOptions_SpiritBeast:  {Name: "Spirit Beast", FocusDump: Claw, SpecialAbility: SpiritBeastBlessing, ExoticAbility: SpiritMend},
	proto.HunterOptions_SporeBat:     {Name: "Spore Bat", FocusDump: Smack, SpecialAbility: SporeCloud},
	proto.HunterOptions_Tallstrider:  {Name: "Tallstrider", FocusDump: Claw, SpecialAbility: DustCloud},
	proto.HunterOptions_Turtle:       {Name: "Turtle", FocusDump: Bite, SpecialAbility: ShellShield},
	proto.HunterOptions_WarpStalker:  {Name: "Warp Stalker", FocusDump: Bite, SpecialAbility: TimeWarp},
	proto.HunterOptions_Wasp:         {Name: "Wasp", FocusDump: Smack, SpecialAbility: Sting},
	proto.HunterOptions_WindSerpent:  {Name: "Wind Serpent", FocusDump: Bite, SpecialAbility: LightningBreath},
	proto.HunterOptions_Wolf:         {Name: "Wolf", FocusDump: Bite, SpecialAbility: FuriousHowl},
	proto.HunterOptions_Worm:         {Name: "Worm", FocusDump: Bite, SpecialAbility: AcidSpitDebuff, ExoticAbility: BurrowAttack},
	proto.HunterOptions_ShaleSpider:  {Name: "Shale Spider", FocusDump: Bite, SpecialAbility: EmbraceOfTheShaleSpider, ExoticAbility: WebWrap},
	proto.HunterOptions_Goat:         {Name: "Goat", FocusDump: Bite, SpecialAbility: Trample},
	proto.HunterOptions_Porcupine:    {Name: "Porcupine", FocusDump: Bite, SpecialAbility: ParalyzingQuill},
	proto.HunterOptions_Monkey:       {Name: "Monkey", FocusDump: Bite, SpecialAbility: BadManners},
	proto.HunterOptions_Basilisk:     {Name: "Basilisk", FocusDump: Bite, SpecialAbility: PetrifyingGaze},
	proto.HunterOptions_Crane:        {Name: "Crane", FocusDump: Bite, SpecialAbility: Lullaby},
	proto.HunterOptions_Dog:          {Name: "Dog", FocusDump: Bite, SpecialAbility: LockJaw},
	proto.HunterOptions_Beetle:       {Name: "Beetle", FocusDump: Bite, SpecialAbility: HardenCarapace},
	proto.HunterOptions_Quilen:       {Name: "Quilen", FocusDump: Bite, SpecialAbility: FearlessRoar, ExoticAbility: EternalGuardian},
	proto.HunterOptions_WaterStrider: {Name: "Water Strider", FocusDump: Claw, SpecialAbility: StillWater, ExoticAbility: SurfaceTrot},
}
