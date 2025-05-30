package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const (
	SpellFlagNaturesGrace = core.SpellFlagAgentReserved1
	SpellFlagOmenTrigger  = core.SpellFlagAgentReserved2
)

type Druid struct {
	core.Character
	SelfBuffs
	// eclipseEnergyBar
	Talents *proto.DruidTalents

	ClassSpellScaling float64

	StartingForm DruidForm

	EclipseEnergyMap EclipseEnergyMap

	RebirthUsed       bool
	RebirthTiming     float64
	BleedsActive      int
	AssumeBleedActive bool
	CannotShredTarget bool

	MHAutoSpell *core.Spell

	Barkskin              *DruidSpell
	Berserk               *DruidSpell
	CatCharge             *DruidSpell
	DemoralizingRoar      *DruidSpell
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
	MightOfUrsoc          *DruidSpell
	Moonfire              *DruidSpell
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
	ThrashBear            *DruidSpell
	ThrashCat             *DruidSpell
	Typhoon               *DruidSpell
	Wrath                 *DruidSpell
	WildMushrooms         *DruidSpell
	WildMushroomsDetonate *DruidSpell

	CatForm  *DruidSpell
	BearForm *DruidSpell

	BarkskinAura             *core.Aura
	BlazeOfGloryAura         *core.Aura
	BearFormAura             *core.Aura
	BerserkBearAura          *core.Aura
	BerserkCatAura           *core.Aura
	CatFormAura              *core.Aura
	ClearcastingAura         *core.Aura
	WeakenedBlowsAuras       core.AuraArray
	FaerieFireAuras          core.AuraArray
	FrenziedRegenerationAura *core.Aura
	LunarEclipseProcAura     *core.Aura
	MightOfUrsocAura         *core.Aura
	NaturesGraceProcAura     *core.Aura
	OwlkinFrenzyAura         *core.Aura
	SavageDefenseAura        *core.DamageAbsorptionAura
	SolarEclipseProcAura     *core.Aura
	SurvivalInstinctsAura    *core.Aura

	SavageRoarDurationTable [6]time.Duration

	ProcOoc func(sim *core.Simulation)

	ExtendingMoonfireStacks int

	Treants       *Treants
	BurningTreant *BurningTreant

	form         DruidForm
	disabledMCDs []*core.MajorCooldown

	// Guardian leather specialization is form-specific
	GuardianLeatherSpecTracker *core.Aura
	GuardianLeatherSpecDep     *stats.StatDependency

	// Item sets
	T13Feral4pBonus *core.Aura
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
	DruidSpellMangleBear
	DruidSpellMangleCat
	DruidSpellMaul
	DruidSpellMoonfire
	DruidSpellMoonfireDoT
	DruidSpellNaturesGrasp
	DruidSpellRavage
	DruidSpellShred
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
	DruidSpellMangle    = DruidSpellMangleBear | DruidSpellMangleCat
	DruidArcaneSpells   = DruidSpellMoonfire | DruidSpellMoonfireDoT | DruidSpellStarfire | DruidSpellStarsurge | DruidSpellStarfall
	DruidNatureSpells   = DruidSpellWrath | DruidSpellInsectSwarm | DruidSpellStarsurge | DruidSpellSunfire | DruidSpellSunfireDoT | DruidSpellTyphoon | DruidSpellHurricane
	DruidHealingSpells  = DruidSpellHealingTouch | DruidSpellRegrowth | DruidSpellRejuvenation | DruidSpellLifebloom | DruidSpellNourish | DruidSpellSwiftmend
	DruidDamagingSpells = DruidArcaneSpells | DruidNatureSpells
)

type SelfBuffs struct {
	InnervateTarget *proto.UnitReference
}

func (druid *Druid) GetCharacter() *core.Character {
	return &druid.Character
}

// func (druid *Druid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
// 	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
// 		raidBuffs.LeaderOfThePack = true
// 	}

// 	if druid.InForm(Moonkin) {
// 		raidBuffs.MoonkinForm = true
// 	}

// 	raidBuffs.MarkOfTheWild = true
// }

func (druid *Druid) BalanceCritMultiplier() float64 {
	return druid.CritMultiplier(1, 0)
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
	druid.form = druid.StartingForm
	//druid.BleedCategories = druid.GetEnemyExclusiveCategories(core.BleedEffectCategory)

	druid.Env.RegisterPostFinalizeEffect(func() {
		druid.MHAutoSpell = druid.AutoAttacks.MHAuto()
		druid.BlazeOfGloryAura = druid.GetAura("Blaze of Glory")
	})

	druid.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(sim *core.Simulation, slot proto.ItemSlot) {
		switch {
		case druid.InForm(Cat):
			druid.AutoAttacks.SetMH(druid.GetCatWeapon())
		case druid.InForm(Bear):
			druid.AutoAttacks.SetMH(druid.GetBearWeapon())
		}
	})

	druid.WeakenedBlowsAuras = druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.WeakenedBlowsAura(target)
	})

	druid.registerFaerieFireSpell()
	// druid.registerRebirthSpell()
	// druid.registerInnervateCD()
	druid.registerTranquilityCD()
}

func (druid *Druid) RegisterBalanceSpells() {
	druid.registerHurricaneSpell()
	druid.registerInsectSwarmSpell()
	druid.registerMoonfireSpell()
	druid.registerSunfireSpell()
	// druid.registerStarfireSpell()
	druid.registerWrathSpell()
	// druid.registerStarfallSpell()
	// druid.registerTyphoonSpell()
	// druid.registerForceOfNature()
	druid.registerStarsurgeSpell()
	druid.registerWildMushrooms()
}

func (druid *Druid) RegisterFeralCatSpells() {
	druid.registerBearFormSpell()
	// druid.registerBerserkCD()
	// druid.registerCatCharge()
	druid.registerCatFormSpell()
	druid.registerFerociousBiteSpell()
	druid.registerLacerateSpell()
	druid.registerMangleBearSpell()
	druid.registerMangleCatSpell()
	druid.registerMaulSpell()
	druid.registerRakeSpell()
	// druid.registerRavageSpell()
	druid.registerRipSpell()
	// druid.registerSavageRoarSpell()
	// druid.registerShredSpell()
	druid.registerSwipeBearSpell()
	druid.registerSwipeCatSpell()
	druid.registerThrashBearSpell()
	druid.registerThrashCatSpell()
	// druid.registerTigersFurySpell()
}

func (druid *Druid) RegisterFeralTankSpells() {
	druid.registerBarkskinCD()
	druid.registerBearFormSpell()
	// druid.registerBerserkCD()
	//druid.registerDemoralizingRoarSpell()
	// druid.registerEnrageSpell()
	druid.registerFrenziedRegenerationSpell()
	druid.registerMangleBearSpell()
	druid.registerMaulSpell()
	druid.registerMightOfUrsocCD()
	druid.registerLacerateSpell()
	//druid.registerSavageDefensePassive()
	// druid.registerSurvivalInstinctsCD()
	druid.registerSwipeBearSpell()
	druid.registerThrashBearSpell()
}

func (druid *Druid) Reset(_ *core.Simulation) {
	// druid.eclipseEnergyBar.reset()
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
		EclipseEnergyMap:  make(EclipseEnergyMap),
	}

	core.FillTalentsProto(druid.Talents.ProtoReflect(), talents)
	druid.EnableManaBar()

	druid.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	druid.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	druid.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[char.Class])

	// Druids get roughly 1% Dodge per 951.16 Agi at level 90
	druid.AddStatDependency(stats.Agility, stats.DodgeRating, 0.00105135 * core.DodgeRatingPerDodgePercent)

	// Base dodge is unaffected by Diminishing Returns
	druid.PseudoStats.BaseDodgeChance += 0.03

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
