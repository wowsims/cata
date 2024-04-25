package death_knight

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	PetSpellHitScale  = 17.0 / 8.0 * core.SpellHitRatingPerHitChance / core.MeleeHitRatingPerHitChance    // 1.7
	PetExpertiseScale = 3.25 * core.ExpertisePerQuarterPercentReduction / core.MeleeHitRatingPerHitChance // 0.8125
)

var TalentTreeSizes = [3]int{20, 20, 20}

// Damage Done By Caster setup
const (
	DDBC_MercilessCombat   int = 0
	DDBC_EbonPlaguebringer     = iota
	DDBC_RuneOfRazorice

	DDBC_Total
)

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

	// Pets
	Ghoul      *GhoulPet
	Gargoyle   *GargoylePet
	ArmyGhoul  []*GhoulPet
	RuneWeapon *RuneWeaponPet
	Bloodworm  []*BloodwormPet

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell
	EbonPlagueAura   core.AuraArray
	ScarletFeverAura core.AuraArray
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

	// TODO: Make horn of winter dynamic
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
	dk.registerFrostFever()
	dk.registerBloodPlague()
	dk.registerOutbreak()
	dk.registerHornOfWinterSpell()
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
	dk.registerBloodBoilSpell()
	dk.registerRuneStrikeSpell()
	dk.registerDeathStrikeSpell()
	dk.registerRuneTapSpell()
	dk.registerVampiricBloodSpell()
	dk.registerIceboundFortitudeSpell()
	dk.registerBoneShieldSpell()
	dk.registerDancingRuneWeaponSpell()
	dk.registerDeathPactSpell()
	dk.registerAntiMagicShellSpell()
}

func (dk *DeathKnight) Reset(sim *core.Simulation) {
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

func NewDeathKnight(character *core.Character, inputs DeathKnightInputs, talents string, deathRuneConvertSpellId int32) *DeathKnight {
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
		func(sim *core.Simulation, changeType core.RuneChangeType, runeRegen []int8) {
			if deathRuneConvertSpellId == 0 {
				return
			}
			if changeType.Matches(core.ConvertToDeath) {
				spell := dk.GetOrRegisterSpell(core.SpellConfig{
					ActionID:       core.ActionID{SpellID: deathRuneConvertSpellId},
					Flags:          core.SpellFlagNoLogs | core.SpellFlagNoMetrics,
					ClassSpellMask: DeathKnightSpellConvertToDeathRune,
				})
				spell.Cast(sim, nil)
			}
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
	}

	dk.Ghoul = dk.NewGhoulPet(dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight)

	dk.ArmyGhoul = make([]*GhoulPet, 8)
	for i := 0; i < 8; i++ {
		dk.ArmyGhoul[i] = dk.NewArmyGhoulPet(i)
	}

	if dk.Talents.BloodParasite > 0 {
		dk.Bloodworm = make([]*BloodwormPet, 5)
		for i := 0; i < 5; i++ {
			dk.Bloodworm[i] = dk.NewBloodwormPet(i)
		}
	}

	if dk.Talents.DancingRuneWeapon {
		dk.RuneWeapon = dk.NewRuneWeapon()
	}

	return dk
}

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
	DeathKnightSpellHeartStrike
	DeathKnightSpellDeathStrike
	DeathKnightSpellScourgeStrikeShadow
	DeathKnightSpellFrostFever
	DeathKnightSpellBloodPlague
	DeathKnightSpellHowlingBlast
	DeathKnightSpellHornOfWinter
	DeathKnightSpellPillarOfFrost
	DeathKnightSpellPestilence
	DeathKnightSpellBloodBoil
	DeathKnightSpellRuneTap
	DeathKnightSpellVampiricBlood
	DeathKnightSpellIceboundFortitude
	DeathKnightSpellBoneShield
	DeathKnightSpellDancingRuneWeapon
	DeathKnightSpellDeathPact

	DeathKnightSpellDeathStrikeHeal // Heal spell for DS

	DeathKnightSpellKillingMachine     // Used to react to km procs
	DeathKnightSpellConvertToDeathRune // Used to react to death rune gains

	DeathKnightSpellLast
	DeathKnightSpellsAll = DeathKnightSpellLast<<1 - 1

	DeathKnightSpellDisease = DeathKnightSpellFrostFever | DeathKnightSpellBloodPlague
)
