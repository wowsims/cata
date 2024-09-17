package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	SpellFlagNaturesGrace = core.SpellFlagAgentReserved1
	SpellFlagOmenTrigger  = core.SpellFlagAgentReserved2
)

var TalentTreeSizes = [3]int{20, 22, 21}

type Druid struct {
	core.Character
	SelfBuffs
	eclipseEnergyBar
	Talents *proto.DruidTalents

	ClassSpellScaling float64

	StartingForm DruidForm

	RebirthUsed       bool
	RebirthTiming     float64
	BleedsActive      int
	AssumeBleedActive bool
	LeatherSpecActive bool
	Feral4pT12Active  bool

	MHAutoSpell       *core.Spell
	ReplaceBearMHFunc core.ReplaceMHSwing

	Barkskin              *DruidSpell
	Berserk               *DruidSpell
	CatCharge             *DruidSpell
	DemoralizingRoar      *DruidSpell
	Enrage                *DruidSpell
	FaerieFire            *DruidSpell
	FerociousBite         *DruidSpell
	ForceOfNature         *DruidSpell
	FrenziedRegeneration  *DruidSpell
	Hurricane             *DruidSpell
	HurricaneTickSpell    *DruidSpell
	InsectSwarm           *DruidSpell
	GiftOfTheWild         *DruidSpell
	Lacerate              *DruidSpell
	Languish              *DruidSpell
	MangleBear            *DruidSpell
	MangleCat             *DruidSpell
	Maul                  *DruidSpell
	MaulQueueSpell        *DruidSpell
	Moonfire              *DruidSpell
	Pulverize             *DruidSpell
	Rebirth               *DruidSpell
	Rake                  *DruidSpell
	Ravage                *DruidSpell
	Rip                   *DruidSpell
	SavageRoar            *DruidSpell
	Shred                 *DruidSpell
	Starfire              *DruidSpell
	Starfall              *DruidSpell
	Starsurge             *DruidSpell
	Sunfire               *DruidSpell
	SurvivalInstincts     *DruidSpell
	SwipeBear             *DruidSpell
	SwipeCat              *DruidSpell
	TigersFury            *DruidSpell
	Thrash                *DruidSpell
	Typhoon               *DruidSpell
	Wrath                 *DruidSpell
	WildMushrooms         *DruidSpell
	WildMushroomsDetonate *DruidSpell

	CatForm  *DruidSpell
	BearForm *DruidSpell

	BarkskinAura             *core.Aura
	BearFormAura             *core.Aura
	BerserkAura              *core.Aura
	BerserkProcAura          *core.Aura
	CatFormAura              *core.Aura
	ClearcastingAura         *core.Aura
	DemoralizingRoarAuras    core.AuraArray
	EnrageAura               *core.Aura
	FaerieFireAuras          core.AuraArray
	FrenziedRegenerationAura *core.Aura
	LunarEclipseProcAura     *core.Aura
	MaulQueueAura            *core.Aura
	MoonkinT84PCAura         *core.Aura
	NaturesGraceProcAura     *core.Aura
	OwlkinFrenzyAura         *core.Aura
	PredatoryInstinctsAura   *core.Aura
	PrimalMadnessAura        *core.Aura
	PulverizeAura            *core.Aura
	SavageDefenseAura        *core.Aura
	SavageRoarAura           *core.Aura
	SmokescreenAura          *core.Aura
	SolarEclipseProcAura     *core.Aura
	StampedeCatAura          *core.Aura
	StampedeBearAura         *core.Aura
	StrengthOfThePantherAura *core.Aura
	SurvivalInstinctsAura    *core.Aura
	TigersFuryAura           *core.Aura

	BleedCategories core.ExclusiveCategoryArray

	PrimalMadnessRageMetrics       *core.ResourceMetrics
	PrimalPrecisionRecoveryMetrics *core.ResourceMetrics
	SavageRoarDurationTable        [6]time.Duration

	ProcOoc func(sim *core.Simulation)

	ExtendingMoonfireStacks int

	Treants       *Treants
	BurningTreant *BurningTreant

	form         DruidForm
	disabledMCDs []*core.MajorCooldown
}

const (
	WrathBaseEnergyGain     float64 = 13 + 1.0/3
	StarsurgeBaseEnergyGain float64 = 15
	StarfireBaseEnergyGain  float64 = 20
	MoonfireBaseEnergyGain  float64 = 0
	SunfireBaseEnergyGain   float64 = 0

	MoonfireLunarShowerEnergyGain float64 = MoonfireBaseEnergyGain + 8
	SunfireLunarShowerEnergyGain  float64 = SunfireBaseEnergyGain + 8

	Wrath4PT12EnergyGain    float64 = WrathBaseEnergyGain + 3
	Starfire4PT12EnergyGain float64 = StarfireBaseEnergyGain + 5
)

const (
	DruidSpellFlagNone int64 = 0
	DruidSpellBarkskin int64 = 1 << iota
	DruidSpellCyclone
	DruidSpellEntanglingRoots
	DruidSpellFearieFire
	DruidSpellHibernate
	DruidSpellHurricane
	DruidSpellInnervate
	DruidSpellInsectSwarm
	DruidSpellMoonfire
	DruidSpellMoonfireDoT
	DruidSpellNaturesGrasp
	DruidSpellStarfall
	DruidSpellStarfire
	DruidSpellStarsurge
	DruidSpellSunfire
	DruidSpellSunfireDoT
	DruidSpellThorns
	DruidSpellTyphoon
	DruidSpellWildMushroom
	DruidSpellWildMushroomDetonate
	DruidSpellWrath

	DruidSpellHealingTouch
	DruidSpellRegrowth
	DruidSpellLifebloom
	DruidSpellRejuvenation
	DruidSpellNourish
	DruidSpellTranquility
	DruidSpellMarkOfTheWild
	DruidSpellSwiftmend
	DruidSpellWildGrowth

	DruidSpellLast
	DruidSpellsAll      = DruidSpellLast<<1 - 1
	DruidSpellDoT       = DruidSpellInsectSwarm | DruidSpellMoonfireDoT | DruidSpellSunfireDoT
	DruidSpellHoT       = DruidSpellRejuvenation | DruidSpellLifebloom | DruidSpellRegrowth | DruidSpellWildGrowth
	DruidSpellInstant   = DruidSpellBarkskin | DruidSpellInsectSwarm | DruidSpellMoonfire | DruidSpellStarfall | DruidSpellSunfire | DruidSpellFearieFire | DruidSpellBarkskin
	DruidArcaneSpells   = DruidSpellMoonfire | DruidSpellMoonfireDoT | DruidSpellStarfire | DruidSpellStarsurge | DruidSpellStarfall
	DruidNatureSpells   = DruidSpellInsectSwarm | DruidSpellStarsurge | DruidSpellSunfire | DruidSpellSunfireDoT | DruidSpellTyphoon | DruidSpellHurricane
	DruidHealingSpells  = DruidSpellHealingTouch | DruidSpellRegrowth | DruidSpellRejuvenation | DruidSpellLifebloom | DruidSpellNourish | DruidSpellSwiftmend
	DruidDamagingSpells = DruidArcaneSpells | DruidNatureSpells
)

type SelfBuffs struct {
	InnervateTarget *proto.UnitReference
}

func (druid *Druid) GetCharacter() *core.Character {
	return &druid.Character
}

func (druid *Druid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
		raidBuffs.LeaderOfThePack = true
	}

	if druid.InForm(Moonkin) {
		raidBuffs.MoonkinForm = true
	}

	raidBuffs.MarkOfTheWild = true
}

func (druid *Druid) BalanceCritMultiplier() float64 {
	return druid.SpellCritMultiplier(1, 0)
}

func (druid *Druid) HasPrimeGlyph(glyph proto.DruidPrimeGlyph) bool {
	return druid.HasGlyph(int32(glyph))
}

func (druid *Druid) HasMajorGlyph(glyph proto.DruidMajorGlyph) bool {
	return druid.HasGlyph(int32(glyph))
}

func (druid *Druid) HasMinorGlyph(glyph proto.DruidMinorGlyph) bool {
	return druid.HasGlyph(int32(glyph))
}

// func (druid *Druid) TryMaul(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
// 	return druid.MaulReplaceMH(sim, mhSwingSpell)
// }

func (druid *Druid) RegisterSpell(formMask DruidForm, config core.SpellConfig) *DruidSpell {
	prev := config.ExtraCastCondition
	prevModify := config.Cast.ModifyCast

	ds := &DruidSpell{FormMask: formMask}
	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		// Check if we're in allowed form to cast
		// Allow 'humanoid' auto unshift casts
		if (ds.FormMask != Any && !druid.InForm(ds.FormMask)) && !ds.FormMask.Matches(Humanoid) {
			if sim.Log != nil {
				sim.Log("Failed cast to spell %s, wrong form", ds.ActionID)
			}
			return false
		}
		return prev == nil || prev(sim, target)
	}
	config.Cast.ModifyCast = func(sim *core.Simulation, s *core.Spell, c *core.Cast) {
		if !druid.InForm(ds.FormMask) && ds.FormMask.Matches(Humanoid) {
			druid.ClearForm(sim)
		}
		if prevModify != nil {
			prevModify(sim, s, c)
		}
	}

	ds.Spell = druid.Unit.RegisterSpell(config)

	return ds
}

func (druid *Druid) Initialize() {
	druid.LeatherSpecActive = druid.MeetsArmorSpecializationRequirement(proto.ArmorType_ArmorTypeLeather)
	druid.BleedCategories = druid.GetEnemyExclusiveCategories(core.BleedEffectCategory)

	druid.Env.RegisterPostFinalizeEffect(func() {
		druid.MHAutoSpell = druid.AutoAttacks.MHAuto()
	})

	// Leather spec would always provide 5% intellect regardless of the Druid spec or form
	if druid.LeatherSpecActive {
		druid.MultiplyStat(stats.Intellect, 1.05)
	}

	druid.registerFaerieFireSpell()
	// druid.registerRebirthSpell()
	druid.registerInnervateCD()
	druid.applyOmenOfClarity()
}

func (druid *Druid) RegisterBalanceSpells() {
	druid.registerHurricaneSpell()
	druid.registerInsectSwarmSpell()
	druid.registerMoonfireSpell()
	druid.registerSunfireSpell()
	druid.registerStarfireSpell()
	druid.registerWrathSpell()
	druid.registerStarfallSpell()
	druid.registerTyphoonSpell()
	druid.registerForceOfNature()
	druid.registerStarsurgeSpell()
	druid.registerWildMushrooms()
}

func (druid *Druid) RegisterFeralCatSpells() {
	druid.registerBearFormSpell()
	druid.registerBerserkCD()
	druid.registerCatCharge()
	druid.registerCatFormSpell()
	druid.registerEnrageSpell()
	druid.registerFerociousBiteSpell()
	druid.registerLacerateSpell()
	druid.registerMangleBearSpell()
	druid.registerMangleCatSpell()
	druid.registerMaulSpell()
	druid.registerRakeSpell()
	druid.registerRavageSpell()
	druid.registerRipSpell()
	druid.registerSavageRoarSpell()
	druid.registerShredSpell()
	druid.registerSwipeBearSpell()
	druid.registerSwipeCatSpell()
	druid.registerThrashBearSpell()
	druid.registerTigersFurySpell()
}

func (druid *Druid) RegisterFeralTankSpells() {
	druid.registerBarkskinCD()
	druid.registerBerserkCD()
	druid.registerBearFormSpell()
	druid.registerDemoralizingRoarSpell()
	druid.registerEnrageSpell()
	druid.registerFrenziedRegenerationCD()
	druid.registerMangleBearSpell()
	druid.registerMaulSpell()
	druid.registerLacerateSpell()
	druid.registerPulverizeSpell()
	druid.registerRakeSpell()
	druid.registerRipSpell()
	druid.registerSavageDefensePassive()
	druid.registerSurvivalInstinctsCD()
	druid.registerSwipeBearSpell()
	druid.registerThrashBearSpell()
}

func (druid *Druid) Reset(_ *core.Simulation) {
	druid.eclipseEnergyBar.reset()
	druid.BleedsActive = 0
	druid.form = druid.StartingForm
	druid.disabledMCDs = []*core.MajorCooldown{}
	druid.RebirthUsed = false
}

func New(char *core.Character, form DruidForm, selfBuffs SelfBuffs, talents string) *Druid {
	druid := &Druid{
		Character:         *char,
		SelfBuffs:         selfBuffs,
		Talents:           &proto.DruidTalents{},
		StartingForm:      form,
		form:              form,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassDruid),
	}

	core.FillTalentsProto(druid.Talents.ProtoReflect(), talents, TalentTreeSizes)
	druid.EnableManaBar()

	druid.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	druid.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	druid.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[char.Class])
	// Druids get 0.0041 dodge per agi (before dr), roughly 1% per 244
	druid.AddStatDependency(stats.Agility, stats.DodgeRating, 0.00410000*core.DodgeRatingPerDodgePercent)

	// Base dodge is unaffected by Diminishing Returns
	druid.PseudoStats.BaseDodgeChance += 0.04951

	if druid.Talents.ForceOfNature {
		druid.Treants = &Treants{
			Treant1: druid.NewTreant(),
			Treant2: druid.NewTreant(),
			Treant3: druid.NewTreant(),
		}
	}

	if druid.HasSetBonus(ItemSetObsidianArborweaveRegalia, 2) {
		druid.BurningTreant = druid.NewBurningTreant()
	}

	return druid
}

type DruidSpell struct {
	*core.Spell
	FormMask DruidForm

	// Optional fields used in snapshotting calculations
	CurrentSnapshotPower float64
	NewSnapshotPower     float64
	ShortName            string
}

func (ds *DruidSpell) IsReady(sim *core.Simulation) bool {
	if ds == nil {
		return false
	}
	return ds.Spell.IsReady(sim)
}

func (ds *DruidSpell) CanCast(sim *core.Simulation, target *core.Unit) bool {
	if ds == nil {
		return false
	}
	return ds.Spell.CanCast(sim, target)
}

func (ds *DruidSpell) IsEqual(s *core.Spell) bool {
	if ds == nil || s == nil {
		return false
	}
	return ds.Spell == s
}

func (druid *Druid) UpdateBleedPower(bleedSpell *DruidSpell, sim *core.Simulation, target *core.Unit, updateCurrent bool, updateNew bool) {
	snapshotPower := bleedSpell.ExpectedTickDamage(sim, target)

	// Assume that Mangle will be up soon if not currently active.
	if !druid.BleedCategories.Get(target).AnyActive() {
		snapshotPower *= 1.3
	}

	if updateCurrent {
		bleedSpell.CurrentSnapshotPower = snapshotPower

		if sim.Log != nil {
			druid.Log(sim, "%s Snapshot Power: %.1f", bleedSpell.ShortName, snapshotPower)
		}
	}

	if updateNew {
		bleedSpell.NewSnapshotPower = snapshotPower

		if (sim.Log != nil) && !updateCurrent {
			druid.Log(sim, "%s Projected Power: %.1f", bleedSpell.ShortName, snapshotPower)
		}
	}
}

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type DruidAgent interface {
	GetDruid() *Druid
}

func (druid *Druid) GetDefaultSpellValueProvider() core.DefaultSpellValueProvider { return druid }
