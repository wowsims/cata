package death_knight

import (
	"math"
	"time"

	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const (
	HitCapRatio             = 17.0 / 8.0 // 2.125
	ExpertiseCapRatio       = 6.5 / 8.0  // 0.8125
	PetExpertiseRatingScale = ExpertiseCapRatio * (4 * core.ExpertisePerQuarterPercentReduction)
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
	// Ghoul      *GhoulPet
	// Gargoyle   *GargoylePet
	// ArmyGhoul  []*GhoulPet
	// RuneWeapon *RuneWeaponPet
	Bloodworm []*BloodwormPet

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell
	EbonPlagueAura   core.AuraArray
	ScarletFeverAura core.AuraArray

	// T12 spell
	BurningBloodSpell *core.Spell

	// Runic power decay, used during pre pull
	RunicPowerDecayAura *core.Aura

	// Cached Gurthalak tentacles
	gurthalakTentacles []*cata.TentacleOfTheOldOnesPet

	// Item sets
	T12Tank4pc *core.Aura
	T13Dps2pc  *core.Aura
	T13Dps4pc  *core.Aura
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

func (dk *DeathKnight) ApplyTalents() {
	// dk.ApplyBloodTalents()
	// dk.ApplyFrostTalents()
	// dk.ApplyUnholyTalents()

	// dk.ApplyGlyphs()
}

func (dk *DeathKnight) Initialize() {
	// dk.registerPresences()
	// dk.registerFrostFever()
	// dk.registerBloodPlague()
	// dk.registerOutbreak()
	// dk.registerHornOfWinterSpell()
	// dk.registerIcyTouchSpell()
	// dk.registerPlagueStrikeSpell()
	// dk.registerDeathCoilSpell()
	// dk.registerDeathAndDecaySpell()
	// dk.registerFesteringStrikeSpell()
	// dk.registerEmpowerRuneWeaponSpell()
	// dk.registerUnholyFrenzySpell()
	// dk.registerSummonGargoyleSpell()
	// dk.registerArmyOfTheDeadSpell()
	// dk.registerRaiseDeadSpell()
	// dk.registerBloodTapSpell()
	// dk.registerObliterateSpell()
	// dk.registerHowlingBlastSpell()
	// dk.registerPillarOfFrostSpell()
	// dk.registerPestilenceSpell()
	// dk.registerBloodBoilSpell()
	// dk.registerRuneStrikeSpell()
	// dk.registerDeathStrikeSpell()
	// dk.registerRuneTapSpell()
	// dk.registerVampiricBloodSpell()
	// dk.registerIceboundFortitudeSpell()
	// dk.registerBoneShieldSpell()
	// dk.registerDancingRuneWeaponSpell()
	// dk.registerDeathPactSpell()
	// dk.registerAntiMagicShellSpell()
	// dk.registerRunicPowerDecay()
	// dk.registerBloodStrikeSpell()
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
		Character:         *character,
		Talents:           &proto.DeathKnightTalents{},
		Inputs:            inputs,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassDeathKnight),
	}
	core.FillTalentsProto(dk.Talents.ProtoReflect(), talents)

	// TODO: Fix this to work with the new talent system.
	// maxRunicPower := 100.0 + 10.0*float64(dk.Talents.RunicPowerMastery)
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

	// if dk.Talents.SummonGargoyle {
	// 	dk.Gargoyle = dk.NewGargoyle()
	// }

	// dk.Ghoul = dk.NewGhoulPet(dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight)

	// dk.ArmyGhoul = make([]*GhoulPet, 8)
	// for i := 0; i < 8; i++ {
	// 	dk.ArmyGhoul[i] = dk.NewArmyGhoulPet(i)
	// }

	// if dk.Talents.BloodParasite > 0 {
	// 	dk.Bloodworm = make([]*BloodwormPet, 5)
	// 	for i := 0; i < 5; i++ {
	// 		dk.Bloodworm[i] = dk.NewBloodwormPet(i)
	// 	}
	// }

	// if dk.Talents.DancingRuneWeapon {
	// 	dk.RuneWeapon = dk.NewRuneWeapon()
	// }

	dk.EnableAutoAttacks(dk, core.AutoAttackOptions{
		MainHand:       dk.WeaponFromMainHand(dk.DefaultCritMultiplier()),
		OffHand:        dk.WeaponFromOffHand(dk.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	if mh := dk.MainHand(); mh.Name == "Gurthalak, Voice of the Deeps" {
		dk.gurthalakTentacles = make([]*cata.TentacleOfTheOldOnesPet, 10)

		for i := 0; i < 10; i++ {
			dk.gurthalakTentacles[i] = dk.NewTentacleOfTheOldOnesPet()
		}
	}

	return dk
}

func (dk *DeathKnight) registerRunicPowerDecay() {
	decayMetrics := dk.NewRunicPowerMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionPrepull})

	// TODO: Fix this to work with the new talent system.
	// Base decay rate is about 1.25/s
	// For some reason Butchery works out of combat which reduces this by 1/5 or 2/5 respectively
	// decayRate := []float64{1.25, 1.05, 0.85}[dk.Talents.Butchery]
	decayRate := 1.25

	var decay *core.PendingAction
	dk.RunicPowerDecayAura = dk.GetOrRegisterAura(core.Aura{
		Label:    "Runic Power Decay",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if sim.CurrentTime >= 0 || dk.CurrentRunicPower() <= 0 {
				dk.RunicPowerDecayAura.Deactivate(sim)
				return
			}

			dk.SpendRunicPower(sim, decayRate, decayMetrics)

			decay = &core.PendingAction{
				Priority:     core.ActionPriorityPrePull,
				NextActionAt: sim.CurrentTime + time.Second,
				OnAction: func(sim *core.Simulation) {
					if dk.CurrentRunicPower() <= 0 {
						aura.Deactivate(sim)
						return
					}

					dk.SpendRunicPower(sim, decayRate, decayMetrics)

					nextTick := sim.CurrentTime + time.Second
					if nextTick >= 0 {
						aura.Deactivate(sim)
						return
					}

					decay.NextActionAt = nextTick
					sim.AddPendingAction(decay)
				},
			}

			sim.AddPendingAction(decay)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if decay != nil {
				decay.Cancel(sim)
			}
		},
	})
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
	DeathKnightSpellDeathCoilHeal
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
	DeathKnightSpellUnholyBlight
	DeathKnightSpellBloodStrike

	DeathKnightSpellKillingMachine     // Used to react to km procs
	DeathKnightSpellConvertToDeathRune // Used to react to death rune gains

	DeathKnightSpellLast
	DeathKnightSpellsAll = DeathKnightSpellLast<<1 - 1

	DeathKnightSpellDisease = DeathKnightSpellFrostFever | DeathKnightSpellBloodPlague
)
