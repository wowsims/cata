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
	SpellScalingConstant  = 986.626460 // for level 85
)

var TalentTreeSizes = [3]int{20, 22, 21}

type Druid struct {
	core.Character
	SelfBuffs
	Talents *proto.DruidTalents

	StartingForm DruidForm

	RebirthUsed       bool
	RebirthTiming     float64
	BleedsActive      int
	AssumeBleedActive bool
	LeatherSpecActive bool
	RipTfSnapshot     bool

	MHAutoSpell       *core.Spell
	ReplaceBearMHFunc core.ReplaceMHSwing

	Barkskin             *DruidSpell
	Berserk              *DruidSpell
	DemoralizingRoar     *DruidSpell
	Enrage               *DruidSpell
	FaerieFire           *DruidSpell
	FerociousBite        *DruidSpell
	ForceOfNature        *DruidSpell
	FrenziedRegeneration *DruidSpell
	Hurricane            *DruidSpell
	HurricaneTickSpell   *DruidSpell
	InsectSwarm          *DruidSpell
	GiftOfTheWild        *DruidSpell
	Lacerate             *DruidSpell
	Languish             *DruidSpell
	MangleBear           *DruidSpell
	MangleCat            *DruidSpell
	Maul                 *DruidSpell
	MaulQueueSpell       *DruidSpell
	Moonfire             *DruidSpell
	Pulverize            *DruidSpell
	Rebirth              *DruidSpell
	Rake                 *DruidSpell
	Rip                  *DruidSpell
	SavageRoar           *DruidSpell
	Shred                *DruidSpell
	Starfire             *DruidSpell
	Starfall             *DruidSpell
	StarfallSplash       *DruidSpell
	SurvivalInstincts    *DruidSpell
	SwipeBear            *DruidSpell
	SwipeCat             *DruidSpell
	TigersFury           *DruidSpell
	Thrash               *DruidSpell
	Typhoon              *DruidSpell
	Wrath                *DruidSpell

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
	MaulQueueAura            *core.Aura
	MoonkinT84PCAura         *core.Aura
	NaturesGraceProcAura     *core.Aura
	PredatoryInstinctsAura   *core.Aura
	PrimalMadnessAura        *core.Aura
	PulverizeAura            *core.Aura
	SavageDefenseAura        *core.Aura
	SurvivalInstinctsAura    *core.Aura
	TigersFuryAura           *core.Aura
	SavageRoarAura           *core.Aura
	SolarEclipseProcAura     *core.Aura
	LunarEclipseProcAura     *core.Aura
	OwlkinFrenzyAura         *core.Aura

	BleedCategories core.ExclusiveCategoryArray

	PrimalMadnessRageMetrics       *core.ResourceMetrics
	PrimalPrecisionRecoveryMetrics *core.ResourceMetrics
	SavageRoarDurationTable        [6]time.Duration

	ProcOoc func(sim *core.Simulation)

	ExtendingMoonfireStacks int
	LunarICD                core.Cooldown
	SolarICD                core.Cooldown
	// Treant1                 *TreantPet
	// Treant2                 *TreantPet
	// Treant3                 *TreantPet

	form         DruidForm
	disabledMCDs []*core.MajorCooldown
}

type SelfBuffs struct {
	InnervateTarget *proto.UnitReference
}

func (druid *Druid) GetCharacter() *core.Character {
	return &druid.Character
}

func (druid *Druid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	// raidBuffs.GiftOfTheWild = max(raidBuffs.GiftOfTheWild, proto.TristateEffect_TristateEffectRegular)
	// if druid.Talents.ImprovedMarkOfTheWild == 2 { // probably could work on actually calculating the fraction effect later if we care.
	// 	raidBuffs.GiftOfTheWild = proto.TristateEffect_TristateEffectImproved
	// }

	// raidBuffs.Thorns = max(raidBuffs.Thorns, proto.TristateEffect_TristateEffectRegular)
	// if druid.Talents.Brambles == 3 {
	// 	raidBuffs.Thorns = proto.TristateEffect_TristateEffectImproved
	// }

	// if druid.InForm(Moonkin) && druid.Talents.MoonkinForm {
	// 	raidBuffs.MoonkinAura = max(raidBuffs.MoonkinAura, proto.TristateEffect_TristateEffectRegular)
	// 	if druid.Talents.ImprovedMoonkinForm > 0 {
	// 		// For now, we assume Improved Moonkin Form is maxed-out
	// 		raidBuffs.MoonkinAura = proto.TristateEffect_TristateEffectImproved
	// 	}
	// }
	if druid.InForm(Cat|Bear) && druid.Talents.LeaderOfThePack {
		raidBuffs.LeaderOfThePack = true
	}
}

// func (druid *Druid) BalanceCritMultiplier() float64 {
// 	return druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))
// }

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
	druid.LeatherSpecActive = druid.MeetsArmorSpecializationRequirement()
	druid.BleedCategories = druid.GetEnemyExclusiveCategories(core.BleedEffectCategory)

	druid.Env.RegisterPostFinalizeEffect(func() {
		druid.MHAutoSpell = druid.AutoAttacks.MHAuto()
	})

	// if druid.Talents.PrimalPrecision > 0 {
	// 	druid.PrimalPrecisionRecoveryMetrics = druid.NewEnergyMetrics(core.ActionID{SpellID: 48410})
	// }
	druid.registerFaerieFireSpell()
	// druid.registerRebirthSpell()
	// druid.registerInnervateCD()
	// druid.registerFakeGotw()
	druid.applyOmenOfClarity()
}

// func (druid *Druid) RegisterBalanceSpells() {
// 	druid.registerHurricaneSpell()
// 	druid.registerInsectSwarmSpell()
// 	druid.registerMoonfireSpell()
// 	druid.registerStarfireSpell()
// 	druid.registerWrathSpell()
// 	druid.registerStarfallSpell()
// 	druid.registerTyphoonSpell()
// 	druid.registerForceOfNatureCD()
// }

func (druid *Druid) RegisterFeralCatSpells() {
	druid.registerBerserkCD()
	druid.registerCatFormSpell()
	// 	druid.registerBearFormSpell()
	// 	druid.registerEnrageSpell()
	druid.registerFerociousBiteSpell()
	// 	druid.registerMangleBearSpell()
	druid.registerMangleCatSpell()
	// 	druid.registerMaulSpell()
	druid.registerLacerateSpell()
	druid.registerRakeSpell()
	druid.registerRipSpell()
	druid.registerSavageRoarSpell()
	druid.registerShredSpell()
	// 	druid.registerSwipeBearSpell()
	// 	druid.registerSwipeCatSpell()
	druid.registerTigersFurySpell()
	druid.registerThrashBearSpell()
}

// func (druid *Druid) RegisterFeralTankSpells() {
// 	druid.registerBarkskinCD()
// 	druid.registerBerserkCD()
// 	druid.registerBearFormSpell()
// 	druid.registerDemoralizingRoarSpell()
// 	druid.registerEnrageSpell()
// 	druid.registerFrenziedRegenerationCD()
// 	druid.registerMangleBearSpell()
// 	druid.registerMaulSpell()
// 	druid.registerLacerateSpell()
// 	druid.registerRakeSpell()
// 	druid.registerRipSpell()
// 	druid.registerSavageDefensePassive()
// 	druid.registerSurvivalInstinctsCD()
// 	druid.registerSwipeBearSpell()
//  druid.registerThrashBearSpell()
// }

func (druid *Druid) Reset(_ *core.Simulation) {
	// druid.BleedsActive = 0
	// druid.form = druid.StartingForm
	// druid.disabledMCDs = []*core.MajorCooldown{}
	// druid.RebirthUsed = false
	// druid.LunarICD.Timer.Reset()
	// druid.SolarICD.Timer.Reset()
}

func New(char *core.Character, form DruidForm, selfBuffs SelfBuffs, talents string) *Druid {
	druid := &Druid{
		Character:    *char,
		SelfBuffs:    selfBuffs,
		Talents:      &proto.DruidTalents{},
		StartingForm: form,
		form:         form,
	}
	core.FillTalentsProto(druid.Talents.ProtoReflect(), talents, TalentTreeSizes)
	druid.EnableManaBar()

	druid.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	druid.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	druid.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[char.Class]*core.CritRatingPerCritChance)
	// 	// // Druid get 0.0209 dodge per agi (before dr), roughly 1 per 47.846
	// 	// druid.AddStatDependency(stats.Agility, stats.Dodge, (0.0209)*core.DodgeRatingPerDodgeChance)

	// 	// // Druids get extra melee haste
	// 	// druid.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3

	// Base dodge is unaffected by Diminishing Returns
	druid.PseudoStats.BaseDodge += 0.056097 // TODO: Check if this is different in Cata

	// 	// if druid.Talents.ForceOfNature {
	// 	// 	druid.Treant1 = druid.NewTreant()
	// 	// 	druid.Treant2 = druid.NewTreant()
	// 	// 	druid.Treant3 = druid.NewTreant()
	// 	// }

	return druid
}

type DruidSpell struct {
	*core.Spell
	FormMask DruidForm
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

// Agent is a generic way to access underlying druid on any of the agents (for example balance druid.)
type DruidAgent interface {
	GetDruid() *Druid
}
