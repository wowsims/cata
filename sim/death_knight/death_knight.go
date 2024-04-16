package death_knight

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const SpellFlagMercilessCombat = core.SpellFlagAgentReserved1

const (
	PetSpellHitScale  = 17.0 / 8.0 * core.SpellHitRatingPerHitChance / core.MeleeHitRatingPerHitChance    // 1.7
	PetExpertiseScale = 3.25 * core.ExpertisePerQuarterPercentReduction / core.MeleeHitRatingPerHitChance // 0.8125
)

var TalentTreeSizes = [3]int{20, 20, 20}

type DeathKnightInputs struct {
	// Option Vars
	IsDps bool

	UnholyFrenzyTarget *proto.UnitReference

	StartingRunicPower float64
	PetUptime          float64

	// Rotation Vars
	UseAMS            bool
	AvgAMSSuccessRate float64
	AvgAMSHit         float64

	Spec proto.Spec
}

type DeathKnight struct {
	core.Character
	Talents *proto.DeathKnightTalents

	ClassBaseScaling float64

	Inputs DeathKnightInputs

	Ghoul     *GhoulPet
	RaiseDead *core.Spell

	Gargoyle                 *GargoylePet
	OnGargoyleStartFirstCast func()

	//RuneWeapon        *RuneWeaponPet
	DancingRuneWeapon *core.Spell

	ArmyOfTheDead *core.Spell
	ArmyGhoul     []*GhoulPet

	//Bloodworm []*BloodwormPet

	IcyTouch   *core.Spell
	BloodBoil  *core.Spell
	Pestilence *core.Spell

	PlagueStrike    *core.Spell
	FesteringStrike *core.Spell

	DeathStrike      *core.Spell
	DeathStrikeMhHit *core.Spell
	DeathStrikeOhHit *core.Spell
	DeathStrikeHeals []float64

	Obliterate      *core.Spell
	ObliterateMhHit *core.Spell
	ObliterateOhHit *core.Spell

	BloodStrike      *core.Spell
	BloodStrikeMhHit *core.Spell
	BloodStrikeOhHit *core.Spell

	FrostStrike      *core.Spell
	FrostStrikeMhHit *core.Spell
	FrostStrikeOhHit *core.Spell

	HeartStrike       *core.Spell
	HeartStrikeOffHit *core.Spell

	RuneStrikeQueued bool
	RuneStrikeQueue  *core.Spell
	RuneStrike       *core.Spell
	RuneStrikeOh     *core.Spell
	RuneStrikeAura   *core.Aura

	GhoulFrenzy *core.Spell
	// Dummy aura for timeline metrics
	GhoulFrenzyAura *core.Aura

	ScourgeStrike *core.Spell

	DeathCoil *core.Spell

	DeathAndDecay *core.Spell

	HowlingBlast *core.Spell

	HasDraeneiHitAura bool
	HornOfWinter      *core.Spell

	// "CDs"
	RuneTap     *core.Spell
	MarkOfBlood *core.Spell

	BloodTap     *core.Spell
	BloodTapAura *core.Aura

	AntiMagicShell     *core.Spell
	AntiMagicShellAura *core.Aura

	EmpowerRuneWeapon *core.Spell

	VampiricBlood     *core.Spell
	VampiricBloodAura *core.Aura

	BoneShield     *core.Spell
	BoneShieldAura *core.Aura

	IceboundFortitude     *core.Spell
	IceboundFortitudeAura *core.Aura

	DeathPact *core.Spell

	// Used only to proc stuff as its free GCD
	MindFreezeSpell *core.Spell

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell
	EbonPlagueAura   core.AuraArray

	//UnholyBlightSpell *core.Spell

	// Talent Auras
	KillingMachineAura  *core.Aura
	IcyTalonsAura       *core.Aura
	BloodCakedBladeAura *core.Aura
	ButcheryAura        *core.Aura
	ButcheryPA          *core.PendingAction
	FreezingFogAura     *core.Aura
	ScentOfBloodAura    *core.Aura
	WillOfTheNecropolis *core.Aura

	// Presences
	BloodPresence      *core.Spell
	BloodPresenceAura  *core.Aura
	FrostPresence      *core.Spell
	FrostPresenceAura  *core.Aura
	UnholyPresence     *core.Spell
	UnholyPresenceAura *core.Aura
}

func (dk *DeathKnight) GetCharacter() *core.Character {
	return &dk.Character
}

func (dk *DeathKnight) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (dk *DeathKnight) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if dk.Talents.AbominationsMight > 0 {
		raidBuffs.AbominationsMight = true
	}

	if dk.Talents.ImprovedIcyTalons {
		raidBuffs.IcyTalons = true
	}

	raidBuffs.HornOfWinter = true
}

func (dk *DeathKnight) ApplyTalents() {
	// Apply Armor Spec
	dk.EnableArmorSpecialization(stats.Strength, proto.ArmorType_ArmorTypePlate)

	dk.ApplyBloodTalents()
	dk.ApplyFrostTalents()
	dk.ApplyUnholyTalents()

	dk.ApplyGlyphs()
}

func (dk *DeathKnight) Initialize() {
	dk.registerPresences()

	dk.registerHornOfWinterSpell()
	dk.registerDiseaseDots()
	dk.registerIcyTouchSpell()
	dk.registerPlagueStrikeSpell()
	dk.registerDeathCoilSpell()
	dk.registerDeathAndDecaySpell()
	dk.registerFesteringStrikeSpell()
	dk.registerEmpowerRuneWeaponSpell()
	dk.registerUnholyFrenzySpell()
	dk.registerSummonGargoyleSpell()
	dk.registerArmyOfTheDeadSpell()
	dk.registerRaiseDeadSpell()
	dk.registerBloodTapSpell()
	dk.registerObliterateSpell()
	dk.registerHowlingBlastSpell()
	dk.registerPillarOfFrostSpell()
	dk.registerPestilenceSpell()
}

func (dk *DeathKnight) Reset(sim *core.Simulation) {
	dk.DeathStrikeHeals = dk.DeathStrikeHeals[:0]
}

func (dk *DeathKnight) HasPrimeGlyph(glyph proto.DeathKnightPrimeGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}
func (dk *DeathKnight) HasMajorGlyph(glyph proto.DeathKnightMajorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}
func (dk *DeathKnight) HasMinorGlyph(glyph proto.DeathKnightMinorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}

func NewDeathKnight(character *core.Character, inputs DeathKnightInputs, talents string) *DeathKnight {
	dk := &DeathKnight{
		Character:        *character,
		Talents:          &proto.DeathKnightTalents{},
		Inputs:           inputs,
		ClassBaseScaling: 1125.227400,
	}
	core.FillTalentsProto(dk.Talents.ProtoReflect(), talents, TalentTreeSizes)

	maxRunicPower := 100.0 + 15.0*float64(dk.Talents.RunicPowerMastery)
	currentRunicPower := math.Min(maxRunicPower, dk.Inputs.StartingRunicPower)

	dk.EnableRunicPowerBar(
		currentRunicPower,
		maxRunicPower,
		10*time.Second,
		func(sim *core.Simulation, changeType core.RuneChangeType) {
		},
		nil,
	)

	// Runic Focus
	dk.AddStat(stats.SpellHit, 9*core.SpellHitRatingPerHitChance)

	dk.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/243.7)
	dk.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/430.69289874)
	dk.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	dk.AddStatDependency(stats.Strength, stats.Parry, 0.25)
	dk.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	dk.PseudoStats.CanParry = true

	// 	// Base dodge unaffected by Diminishing Returns
	dk.PseudoStats.BaseDodge += 0.03664
	dk.PseudoStats.BaseParry += 0.05

	if dk.Talents.SummonGargoyle {
		dk.Gargoyle = dk.NewGargoyle()
		dk.OnGargoyleStartFirstCast = func() {}
	}

	dk.Ghoul = dk.NewGhoulPet(dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight)

	dk.ArmyGhoul = make([]*GhoulPet, 8)
	for i := 0; i < 8; i++ {
		dk.ArmyGhoul[i] = dk.NewArmyGhoulPet(i)
	}

	// 	if dk.Talents.Bloodworms > 0 {
	// 		dk.Bloodworm = make([]*BloodwormPet, 4)
	// 		for i := 0; i < 4; i++ {
	// 			dk.Bloodworm[i] = dk.NewBloodwormPet(i)
	// 		}
	// 	}

	// 	if dk.Talents.DancingRuneWeapon {
	// 		dk.RuneWeapon = dk.NewRuneWeapon()
	// 	}

	// 	// done here so enchants that modify stats are applied before stats are calculated
	// 	dk.registerItems()

	return dk
}

// Agent is a generic way to access underlying warrior on any of the agents.

func (dk *DeathKnight) GetDeathKnight() *DeathKnight {
	return dk
}

type DeathKnightAgent interface {
	GetDeathKnight() *DeathKnight
}

const (
	DeathKnightSpellFlagNone int64 = 0
	DeathKnightSpellIcyTouch int64 = 1 << iota
	DeathKnightSpellDeathCoil
	DeathKnightSpellDeathAndDecay
	DeathKnightSpellOutbreak
	DeathKnightSpellEmpowerRuneWeapon
	DeathKnightSpellUnholyFrenzy
	DeathKnightSpellDarkTransformation
	DeathKnightSpellSummonGargoyle
	DeathKnightSpellArmyOfTheDead
	DeathKnightSpellRaiseDead
	DeathKnightSpellBloodTap
	DeathKnightSpellObliterate
	DeathKnightSpellFrostStrike
	DeathKnightSpellRuneStrike
	DeathKnightSpellPlagueStrike
	DeathKnightSpellFesteringStrike
	DeathKnightSpellScourgeStrike
	DeathKnightSpellScourgeStrikeShadow
	DeathKnightSpellFrostFever
	DeathKnightSpellBloodPlague
	DeathKnightSpellHowlingBlast
	DeathKnightSpellHornOfWinter
	DeathKnightSpellPillarOfFrost
	DeathKnightSpellPestilence

	DeathKnightSpellLast
	DeathKnightSpellsAll = DeathKnightSpellLast<<1 - 1

	DeathKnightSpellDisease = DeathKnightSpellFrostFever | DeathKnightSpellBloodPlague
)
