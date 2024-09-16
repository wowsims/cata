package death_knight

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	HitCapRatio             = 17.0 / 8.0 // 2.125
	ExpertiseCapRatio       = 6.5 / 8.0  // 0.8125
	PetExpertiseRatingScale = ExpertiseCapRatio * (4 * core.ExpertisePerQuarterPercentReduction)
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

	ClassSpellScaling float64

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

	// T12 spell
	BurningBloodSpell *core.Spell
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
		Character:         *character,
		Talents:           &proto.DeathKnightTalents{},
		Inputs:            inputs,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassDeathKnight),
	}
	core.FillTalentsProto(dk.Talents.ProtoReflect(), talents, TalentTreeSizes)

	maxRunicPower := 100.0 + 10.0*float64(dk.Talents.RunicPowerMastery)
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
				deathConvertSpell := dk.GetOrRegisterSpell(core.SpellConfig{
					ActionID:       core.ActionID{SpellID: deathRuneConvertSpellId},
					Flags:          core.SpellFlagNoLogs | core.SpellFlagNoMetrics,
					ClassSpellMask: DeathKnightSpellConvertToDeathRune,
				})
				deathConvertSpell.Cast(sim, nil)
			}
		},
		nil,
	)

	// Runic Focus
	dk.AddStat(stats.SpellHitPercent, 9)

	dk.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	dk.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[dk.Class])

	dk.AddStat(stats.ParryRating, -dk.GetBaseStats()[stats.Strength]*0.27)
	dk.AddStatDependency(stats.Strength, stats.ParryRating, 0.27)

	dk.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	dk.PseudoStats.CanParry = true

	// 	// Base dodge unaffected by Diminishing Returns
	dk.PseudoStats.BaseDodgeChance += 0.05
	dk.PseudoStats.BaseParryChance += 0.05

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

	dk.EnableAutoAttacks(dk, core.AutoAttackOptions{
		MainHand:       dk.WeaponFromMainHand(dk.DefaultMeleeCritMultiplier()),
		OffHand:        dk.WeaponFromOffHand(dk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

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
	DeathKnightSpellScourgeStrikeShadow
	DeathKnightSpellHeartStrike
	DeathKnightSpellDeathStrike
	DeathKnightSpellDeathStrikeHeal
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

	DeathKnightSpellKillingMachine     // Used to react to km procs
	DeathKnightSpellConvertToDeathRune // Used to react to death rune gains

	DeathKnightSpellLast
	DeathKnightSpellsAll = DeathKnightSpellLast<<1 - 1

	DeathKnightSpellDisease = DeathKnightSpellFrostFever | DeathKnightSpellBloodPlague
)
