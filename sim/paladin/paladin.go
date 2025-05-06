package paladin

import (
	"time"

	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const (
	SpellFlagSecondaryJudgement = core.SpellFlagAgentReserved1
	SpellFlagPrimaryJudgement   = core.SpellFlagAgentReserved2
)

const (
	SpellMaskTemplarsVerdict int64 = 1 << iota
	SpellMaskCrusaderStrike
	SpellMaskDivineStorm
	SpellMaskExorcism
	SpellMaskGlyphOfExorcism
	SpellMaskHammerOfWrath
	SpellMaskJudgementBase
	SpellMaskJudgementOfTruth
	SpellMaskJudgementOfInsight
	SpellMaskJudgementOfRighteousness
	SpellMaskJudgementOfJustice
	SpellMaskHolyWrath
	SpellMaskConsecration
	SpellMaskHammerOfTheRighteousMelee
	SpellMaskHammerOfTheRighteousAoe
	SpellMaskAvengersShield
	SpellMaskDivinePlea
	SpellMaskDivineProtection
	SpellMaskAvengingWrath
	SpellMaskCensure
	SpellMaskInquisition
	SpellMaskHandOfLight
	SpellMaskZealotry
	SpellMaskGuardianOfAncientKings
	SpellMaskAncientFury
	SpellMaskSealsOfCommand
	SpellMaskShieldOfTheRighteous
	SpellMaskHolyShield
	SpellMaskArdentDefender

	SpellMaskHolyShock
	SpellMaskWordOfGlory

	SpellMaskSealOfTruth
	SpellMaskSealOfInsight
	SpellMaskSealOfRighteousness
	SpellMaskSealOfJustice
)

const SpellMaskBuilder = SpellMaskCrusaderStrike |
	SpellMaskDivineStorm |
	SpellMaskHammerOfTheRighteousMelee

const SpellMaskHammerOfTheRighteous = SpellMaskHammerOfTheRighteousMelee | SpellMaskHammerOfTheRighteousAoe

const SpellMaskJudgement = SpellMaskJudgementOfTruth |
	SpellMaskJudgementOfInsight |
	SpellMaskJudgementOfRighteousness |
	SpellMaskJudgementOfJustice

const SpellMaskCanTriggerSealOfJustice = SpellMaskCrusaderStrike |
	SpellMaskTemplarsVerdict |
	SpellMaskHammerOfWrath

const SpellMaskCanTriggerSealOfInsight = SpellMaskCanTriggerSealOfJustice

const SpellMaskCanTriggerSealOfRighteousness = SpellMaskCanTriggerSealOfJustice |
	SpellMaskDivineStorm |
	SpellMaskHammerOfTheRighteousMelee

const SpellMaskCanTriggerSealOfTruth = SpellMaskCrusaderStrike |
	SpellMaskTemplarsVerdict |
	SpellMaskExorcism |
	SpellMaskHammerOfWrath |
	SpellMaskJudgement |
	SpellMaskHammerOfTheRighteousMelee |
	SpellMaskShieldOfTheRighteous

const SpellMaskCanTriggerAncientPower = SpellMaskCanTriggerSealOfTruth |
	SpellMaskHolyWrath

const SpellMaskCanTriggerDivinePurpose = SpellMaskHammerOfWrath |
	SpellMaskExorcism |
	SpellMaskJudgement |
	SpellMaskHolyWrath |
	SpellMaskTemplarsVerdict |
	SpellMaskDivineStorm |
	SpellMaskInquisition

const SpellMaskCanConsumeDivinePurpose = SpellMaskInquisition |
	SpellMaskTemplarsVerdict

const SpellMaskModifiedByTwoHandedSpec = SpellMaskJudgement |
	SpellMaskSealOfTruth |
	SpellMaskSealsOfCommand |
	SpellMaskHammerOfWrath

const SpellMaskModifiedByZealOfTheCrusader = SpellMaskTemplarsVerdict |
	SpellMaskCrusaderStrike |
	SpellMaskDivineStorm |
	SpellMaskExorcism |
	SpellMaskGlyphOfExorcism |
	SpellMaskHammerOfWrath |
	SpellMaskJudgementOfTruth |
	SpellMaskJudgementOfRighteousness |
	SpellMaskHolyWrath |
	SpellMaskConsecration |
	SpellMaskHammerOfTheRighteousMelee |
	SpellMaskHammerOfTheRighteousAoe |
	SpellMaskAvengersShield |
	SpellMaskCensure |
	SpellMaskSealsOfCommand |
	SpellMaskHolyShock |
	SpellMaskSealOfTruth |
	SpellMaskSealOfRighteousness |
	SpellMaskSealOfJustice

type Paladin struct {
	core.Character

	PaladinAura proto.PaladinAura
	Seal        proto.PaladinSeal
	HolyPower   core.SecondaryResourceBar

	Talents *proto.PaladinTalents

	// Used for CS/DS/HotR
	sharedBuilderTimer  *core.Timer
	sharedBuilderBaseCD time.Duration

	CurrentSeal       *core.Aura
	CurrentJudgement  *core.Spell
	StartingHolyPower int32

	// Pets
	AncientGuardian *AncientGuardianPet

	DivinePlea               *core.Spell
	DivineStorm              *core.Spell
	HolyWrath                *core.Spell
	Consecration             *core.Spell
	CrusaderStrike           *core.Spell
	Exorcism                 *core.Spell
	HolyShield               *core.Spell
	HammerOfTheRighteous     *core.Spell
	HandOfReckoning          *core.Spell
	ShieldOfRighteousness    *core.Spell
	AvengersShield           *core.Spell
	HammerOfWrath            *core.Spell
	AvengingWrath            *core.Spell
	DivineProtection         *core.Spell
	TemplarsVerdict          *core.Spell
	Zealotry                 *core.Spell
	Inquisition              *core.Spell
	HandOfLight              *core.Spell
	JudgementOfTruth         *core.Spell
	JudgementOfInsight       *core.Spell
	JudgementOfRighteousness *core.Spell
	JudgementOfJustice       *core.Spell
	ShieldOfTheRighteous     *core.Spell

	HolyShieldAura          *core.Aura
	RighteousFuryAura       *core.Aura
	DivinePleaAura          *core.Aura
	SealOfTruthAura         *core.Aura
	SealOfInsightAura       *core.Aura
	SealOfRighteousnessAura *core.Aura
	SealOfJusticeAura       *core.Aura
	AvengingWrathAura       *core.Aura
	DivineProtectionAura    *core.Aura
	ZealotryAura            *core.Aura
	InquisitionAura         *core.Aura
	DivinePurposeAura       *core.Aura
	JudgementsOfThePureAura *core.Aura
	GrandCrusaderAura       *core.Aura
	SacredDutyAura          *core.Aura
	GoakAura                *core.Aura
	AncientPowerAura        *core.Aura

	// Cached Gurthalak tentacles
	gurthalakTentacles []*cata.TentacleOfTheOldOnesPet

	// Item sets
	T11Ret4pc *core.Aura
}

func (paladin *Paladin) GetTentacles() []*cata.TentacleOfTheOldOnesPet {
	return paladin.gurthalakTentacles
}

func (paladin *Paladin) NewTentacleOfTheOldOnesPet() *cata.TentacleOfTheOldOnesPet {
	pet := cata.NewTentacleOfTheOldOnesPet(&paladin.Character)
	paladin.AddPet(pet)
	return pet
}

// Implemented by each Paladin spec.
type PaladinAgent interface {
	GetPaladin() *Paladin
}

func (paladin *Paladin) GetCharacter() *core.Character {
	return &paladin.Character
}

func (paladin *Paladin) HasMajorGlyph(glyph proto.PaladinMajorGlyph) bool {
	return paladin.HasGlyph(int32(glyph))
}
func (paladin *Paladin) HasMinorGlyph(glyph proto.PaladinMinorGlyph) bool {
	return paladin.HasGlyph(int32(glyph))
}

func (paladin *Paladin) GetPaladin() *Paladin {
	return paladin
}

// func (paladin *Paladin) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
// 	if paladin.PaladinAura == proto.PaladinAura_Devotion {
// 		raidBuffs.DevotionAura = true
// 	}
// 	if paladin.PaladinAura == proto.PaladinAura_Retribution {
// 		raidBuffs.RetributionAura = true
// 	}
// 	if paladin.PaladinAura == proto.PaladinAura_Resistance {
// 		raidBuffs.ResistanceAura = true
// 	}
// 	if paladin.Talents.Communion {
// 		raidBuffs.Communion = true
// 	}
// }

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	// paladin.applyGlyphs()
	paladin.registerSpells()
	paladin.addBloodthirstyGloves()
}

func (paladin *Paladin) ApplyTalents() {}

func (paladin *Paladin) registerSpells() {
	paladin.registerCrusaderStrike()
	paladin.registerExorcism()
	paladin.registerJudgement()
	// paladin.registerSealOfTruth()
	paladin.registerSealOfInsight()
	// paladin.registerSealOfRighteousness()
	paladin.registerSealOfJustice()
	// paladin.registerInquisition()
	// paladin.registerHammerOfWrathSpell()
	paladin.registerAvengingWrath()
	paladin.registerDivinePleaSpell()
	paladin.registerConsecrationSpell()
	paladin.registerHolyWrath()
	paladin.registerGuardianOfAncientKings()
	paladin.registerDivineProtectionSpell()
}

func (paladin *Paladin) Reset(sim *core.Simulation) {
	// switch paladin.Seal {
	// case proto.PaladinSeal_Truth:
	// 	paladin.CurrentJudgement = paladin.JudgementOfTruth
	// 	paladin.CurrentSeal = paladin.SealOfTruthAura
	// 	paladin.SealOfTruthAura.Activate(sim)
	// case proto.PaladinSeal_Insight:
	// 	paladin.CurrentJudgement = paladin.JudgementOfInsight
	// 	paladin.CurrentSeal = paladin.SealOfInsightAura
	// 	paladin.SealOfInsightAura.Activate(sim)
	// case proto.PaladinSeal_Righteousness:
	// 	paladin.CurrentJudgement = paladin.JudgementOfRighteousness
	// 	paladin.CurrentSeal = paladin.SealOfRighteousnessAura
	// 	paladin.SealOfRighteousnessAura.Activate(sim)
	// case proto.PaladinSeal_Justice:
	// 	paladin.CurrentJudgement = paladin.JudgementOfJustice
	// 	paladin.CurrentSeal = paladin.SealOfJusticeAura
	// 	paladin.SealOfJusticeAura.Activate(sim)
	// }
}

func NewPaladin(character *core.Character, talentsStr string, options *proto.PaladinOptions) *Paladin {
	paladin := &Paladin{
		Character:           *character,
		Talents:             &proto.PaladinTalents{},
		Seal:                options.Seal,
		PaladinAura:         options.Aura,
		sharedBuilderBaseCD: time.Millisecond * core.TernaryDuration(character.Spec == proto.Spec_SpecProtectionPaladin, 3000, 4500),
	}

	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr)

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()
	paladin.HolyPower = HolyPowerBar{
		DefaultSecondaryResourceBarImpl: paladin.NewDefaultSecondaryResourceBar(core.SecondaryResourceConfig{
			Type:    proto.SecondaryResourceType_SecondaryResourceTypeHolyPower,
			Max:     3,
			Default: paladin.StartingHolyPower,
		}),
		paladin: paladin,
	}
	paladin.RegisterSecondaryResourceBar(paladin.HolyPower)

	// Only retribution and holy are actually pets performing some kind of action
	if paladin.Spec != proto.Spec_SpecProtectionPaladin {
		paladin.AncientGuardian = paladin.NewAncientGuardian()
	}

	paladin.EnableAutoAttacks(paladin, core.AutoAttackOptions{
		MainHand:       paladin.WeaponFromMainHand(paladin.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	paladin.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])
	paladin.AddStat(stats.ParryRating, -paladin.GetBaseStats()[stats.Strength]*0.27) // Does not apply to base Strength
	paladin.AddStatDependency(stats.Strength, stats.ParryRating, 0.27)

	paladin.PseudoStats.BaseDodgeChance += 0.05
	paladin.PseudoStats.BaseParryChance += 0.05

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	if mh := paladin.MainHand(); mh.Name == "Gurthalak, Voice of the Deeps" {
		paladin.gurthalakTentacles = make([]*cata.TentacleOfTheOldOnesPet, 10)

		for i := 0; i < 10; i++ {
			paladin.gurthalakTentacles[i] = paladin.NewTentacleOfTheOldOnesPet()
		}
	}

	return paladin
}

// Shared cooldown for builders
func (paladin *Paladin) BuilderCooldown() *core.Timer {
	return paladin.Character.GetOrInitTimer(&paladin.sharedBuilderTimer)
}
