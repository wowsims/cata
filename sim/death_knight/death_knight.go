package death_knight

import (
	"math"
	"time"

	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

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

	Spec proto.Spec
}

type DeathKnight struct {
	core.Character
	Talents *proto.DeathKnightTalents

	Inputs DeathKnightInputs

	// Pets
	Ghoul      *GhoulPet
	Gargoyle   *GargoylePet
	ArmyGhoul  []*GhoulPet
	RuneWeapon *RuneWeaponPet
	Bloodworm  []*BloodwormPet

	PestilenceSpell *core.Spell
	RuneTapSpell    *core.Spell

	ConversionAura     *core.Aura
	UnholyPresenceAura *core.Aura

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell
	ScarletFeverAura core.AuraArray

	// Runic power decay, used during pre pull
	RunicPowerDecayAura *core.Aura

	// Cached Gurthalak tentacles
	gurthalakTentacles []*cata.TentacleOfTheOldOnesPet

	// T12 spell
	BurningBloodSpell *core.Spell

	// Item sets
	T12Tank4pc *core.Aura
	T13Dps2pc  *core.Aura
	T13Dps4pc  *core.Aura

	// Used for T13 Tank 4pc
	VampiricBloodBonusHealth float64

	// Modified by T14 Tank 4pc
	deathStrikeHealingMultiplier float64
}

func (deathKnight *DeathKnight) GetTentacles() []*cata.TentacleOfTheOldOnesPet {
	return deathKnight.gurthalakTentacles
}

func (dk *DeathKnight) NewTentacleOfTheOldOnesPet() *cata.TentacleOfTheOldOnesPet {
	pet := cata.NewTentacleOfTheOldOnesPet(&dk.Character)
	dk.AddPet(pet)
	return pet
}

func (dk *DeathKnight) GetCharacter() *core.Character {
	return &dk.Character
}

func (dk *DeathKnight) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

// func (dk *DeathKnight) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
// 	if dk.Talents.AbominationsMight > 0 {
// 		raidBuffs.AbominationsMight = true
// 	}

// 	if dk.Talents.ImprovedIcyTalons {
// 		raidBuffs.IcyTalons = true
// 	}

// 	// TODO: Make horn of winter dynamic
// 	raidBuffs.HornOfWinter = true
// }

func (dk *DeathKnight) Initialize() {
	// dk.registerAntiMagicShell()
	dk.registerArmyOfTheDead()
	// dk.registerBloodBoil()
	dk.registerBloodPlague()
	dk.registerDeathAndDecay()
	// dk.registerDeathCoil()
	dk.registerDeathStrike()
	dk.registerEmpowerRuneWeapon()
	dk.registerFrostFever()
	dk.registerHornOfWinter()
	// dk.registerIceboundFortitude()
	dk.registerIcyTouch()
	dk.registerOutbreak()
	dk.registerPestilence()
	dk.registerPlagueStrike()
	// dk.registerPresences()
	// If talented as permanent pet skip this spell
	if dk.Inputs.Spec != proto.Spec_SpecUnholyDeathKnight {
		dk.registerRaiseDead()
	}
	dk.registerRunicPowerDecay()
}

func (dk *DeathKnight) Reset(sim *core.Simulation) {
}

func (dk *DeathKnight) HasMajorGlyph(glyph proto.DeathKnightMajorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}
func (dk *DeathKnight) HasMinorGlyph(glyph proto.DeathKnightMinorGlyph) bool {
	return dk.HasGlyph(int32(glyph))
}

func NewDeathKnight(character *core.Character, inputs DeathKnightInputs, talents string, deathRuneConvertSpellId int32) *DeathKnight {
	dk := &DeathKnight{
		Character: *character,
		Talents:   &proto.DeathKnightTalents{},
		Inputs:    inputs,
	}
	core.FillTalentsProto(dk.Talents.ProtoReflect(), talents)

	maxRunicPower := 100.0
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
		func(sim *core.Simulation) {
			if sim.CurrentTime >= 0 || dk.RunicPowerDecayAura.IsActive() {
				return
			}

			dk.RunicPowerDecayAura.Activate(sim)
		},
	)

	dk.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	dk.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[dk.Class])

	strengthToParryRating := (1 / 951.158596) * core.ParryRatingPerParryPercent
	dk.AddStat(stats.ParryRating, -dk.GetBaseStats()[stats.Strength]*strengthToParryRating) // Does not apply to base Strength
	dk.AddStatDependency(stats.Strength, stats.ParryRating, strengthToParryRating)

	dk.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	dk.PseudoStats.CanParry = true

	// 	// Base dodge unaffected by Diminishing Returns
	dk.PseudoStats.BaseDodgeChance += 0.03
	dk.PseudoStats.BaseParryChance += 0.03

	dk.Ghoul = dk.NewGhoulPet(dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight)

	dk.ArmyGhoul = make([]*GhoulPet, 8)
	for i := range 8 {
		dk.ArmyGhoul[i] = dk.NewArmyGhoulPet(i)
	}

	dk.EnableAutoAttacks(dk, core.AutoAttackOptions{
		MainHand:       dk.WeaponFromMainHand(dk.DefaultCritMultiplier()),
		OffHand:        dk.WeaponFromOffHand(dk.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	if mh := dk.MainHand(); mh.Name == "Gurthalak, Voice of the Deeps" {
		dk.gurthalakTentacles = make([]*cata.TentacleOfTheOldOnesPet, 10)

		for i := range 10 {
			dk.gurthalakTentacles[i] = dk.NewTentacleOfTheOldOnesPet()
		}
	}

	dk.deathStrikeHealingMultiplier = 0.2

	return dk
}

func (dk *DeathKnight) GetDeathKnight() *DeathKnight {
	return dk
}

type DeathKnightAgent interface {
	GetDeathKnight() *DeathKnight
}

const (
	DeathKnightSpellFlagNone      int64 = 0
	DeathKnightSpellArmyOfTheDead int64 = 1 << iota
	DeathKnightSpellBloodBoil
	DeathKnightSpellBloodPlague
	DeathKnightSpellBloodStrike
	DeathKnightSpellBloodTap
	DeathKnightSpellBoneShield
	DeathKnightSpellConversion
	DeathKnightSpellDancingRuneWeapon
	DeathKnightSpellDarkCommand
	DeathKnightSpellDarkTransformation
	DeathKnightSpellDeathAndDecay
	DeathKnightSpellDeathCoil
	DeathKnightSpellDeathCoilHeal
	DeathKnightSpellDeathPact
	DeathKnightSpellDeathSiphon
	DeathKnightSpellDeathStrike
	DeathKnightSpellDeathStrikeHeal
	DeathKnightSpellEmpowerRuneWeapon
	DeathKnightSpellFesteringStrike
	DeathKnightSpellFrostFever
	DeathKnightSpellFrostStrike
	DeathKnightSpellHeartStrike
	DeathKnightSpellHornOfWinter
	DeathKnightSpellHowlingBlast
	DeathKnightSpellIceboundFortitude
	DeathKnightSpellIcyTouch
	DeathKnightSpellLichborne
	DeathKnightSpellObliterate
	DeathKnightSpellOutbreak
	DeathKnightSpellPestilence
	DeathKnightSpellPillarOfFrost
	DeathKnightSpellPlagueLeech
	DeathKnightSpellPlagueStrike
	DeathKnightSpellRaiseDead
	DeathKnightSpellRuneStrike
	DeathKnightSpellRuneTap
	DeathKnightSpellScourgeStrike
	DeathKnightSpellScourgeStrikeShadow
	DeathKnightSpellSummonGargoyle
	DeathKnightSpellUnholyBlight
	DeathKnightSpellUnholyFrenzy
	DeathKnightSpellVampiricBlood

	DeathKnightSpellKillingMachine     // Used to react to km procs
	DeathKnightSpellConvertToDeathRune // Used to react to death rune gains

	DeathKnightSpellLast
	DeathKnightSpellsAll = DeathKnightSpellLast<<1 - 1

	DeathKnightSpellDisease = DeathKnightSpellFrostFever | DeathKnightSpellBloodPlague
)
