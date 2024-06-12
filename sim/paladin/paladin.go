package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
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
	SpellMaskHammerOfTheRighteous
	SpellMaskHandOfReckoning
	SpellMaskShieldOfRighteousness
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

	SpellMaskHolyShock
	SpellMaskWordOfGlory

	SpellMaskSealOfTruth
	SpellMaskSealOfInsight
	SpellMaskSealOfRighteousness
	SpellMaskSealOfJustice
)

const SpellMaskJudgement = SpellMaskJudgementOfTruth |
	SpellMaskJudgementOfInsight |
	SpellMaskJudgementOfRighteousness |
	SpellMaskJudgementOfJustice

const SpellMaskCanTriggerSealOfJustice = SpellMaskCrusaderStrike |
	SpellMaskTemplarsVerdict |
	SpellMaskHammerOfWrath

const SpellMaskCanTriggerSealOfInsight = SpellMaskCanTriggerSealOfJustice

const SpellMaskCanTriggerSealOfRighteousness = SpellMaskCanTriggerSealOfJustice |
	SpellMaskDivineStorm

const SpellMaskCanTriggerSealOfTruth = SpellMaskCrusaderStrike |
	SpellMaskTemplarsVerdict |
	SpellMaskExorcism |
	SpellMaskHammerOfWrath |
	SpellMaskJudgement

const SpellMaskCanTriggerAncientPower = SpellMaskCanTriggerSealOfTruth |
	SpellMaskHolyWrath

const SpellMaskModifiedByInquisition = SpellMaskHammerOfWrath |
	SpellMaskConsecration |
	SpellMaskExorcism |
	SpellMaskGlyphOfExorcism |
	SpellMaskJudgement |
	SpellMaskSealOfTruth |
	SpellMaskCensure |
	SpellMaskHandOfLight |
	SpellMaskHolyWrath |
	SpellMaskAncientFury

const SpellMaskCanTriggerDivinePurpose = SpellMaskHammerOfWrath |
	SpellMaskExorcism |
	SpellMaskJudgement |
	SpellMaskHolyWrath |
	SpellMaskTemplarsVerdict |
	SpellMaskDivineStorm |
	SpellMaskInquisition

const SpellMaskCanConsumeDivinePurpose = SpellMaskInquisition |
	SpellMaskTemplarsVerdict |
	SpellMaskZealotry

var TalentTreeSizes = [3]int{20, 20, 20}

type Paladin struct {
	core.Character
	HolyPowerBar

	PaladinAura proto.PaladinAura
	Seal        proto.PaladinSeal

	Talents *proto.PaladinTalents

	CurrentSeal      *core.Aura
	CurrentJudgement *core.Spell

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

	HolyShieldAura          *core.Aura
	RighteousFuryAura       *core.Aura
	DivinePleaAura          *core.Aura
	SealOfTruthAura         *core.Aura
	SealOfInsightAura       *core.Aura
	SealOfRighteousnessAura *core.Aura
	SealOfJusticeAura       *core.Aura
	AvengingWrathAura       *core.Aura
	DivineProtectionAura    *core.Aura
	ForbearanceAura         *core.Aura
	VengeanceAura           *core.Aura
	ZealotryAura            *core.Aura
	InquisitionAura         *core.Aura
	DivinePurposeAura       *core.Aura

	ArtOfWarInstantCast *core.Aura

	SpiritualAttunementMetrics *core.ResourceMetrics
}

// Implemented by each Paladin spec.
type PaladinAgent interface {
	GetPaladin() *Paladin
}

func (paladin *Paladin) GetCharacter() *core.Character {
	return &paladin.Character
}

func (paladin *Paladin) HasPrimeGlyph(glyph proto.PaladinPrimeGlyph) bool {
	return paladin.HasGlyph(int32(glyph))
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

func (paladin *Paladin) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if paladin.PaladinAura == proto.PaladinAura_Devotion {
		raidBuffs.DevotionAura = true
	}
	if paladin.PaladinAura == proto.PaladinAura_Retribution {
		raidBuffs.RetributionAura = true
	}
	if paladin.PaladinAura == proto.PaladinAura_Resistance {
		raidBuffs.ResistanceAura = true
	}
	if paladin.Talents.Communion {
		raidBuffs.Communion = true
	}
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	paladin.applyGlyphs()
	paladin.registerSpells()
	paladin.addBloodthirstyGloves()
}

func (paladin *Paladin) registerSpells() {
	paladin.registerCrusaderStrike()
	paladin.registerExorcism()
	paladin.registerJudgement()
	paladin.registerSealOfTruth()
	paladin.registerSealOfInsight()
	paladin.registerSealOfRighteousness()
	paladin.registerSealOfJustice()
	paladin.registerInquisition()
	paladin.registerHammerOfWrathSpell()
	paladin.registerAvengingWrath()
	paladin.registerDivinePleaSpell()
	paladin.registerConsecrationSpell()
	paladin.registerHolyWrath()
	paladin.registerGuardianOfAncientKings()
}

func (paladin *Paladin) Reset(sim *core.Simulation) {
	switch paladin.Seal {
	case proto.PaladinSeal_Truth:
		paladin.CurrentJudgement = paladin.JudgementOfTruth
		paladin.CurrentSeal = paladin.SealOfTruthAura
		paladin.SealOfTruthAura.Activate(sim)
	case proto.PaladinSeal_Insight:
		paladin.CurrentJudgement = paladin.JudgementOfInsight
		paladin.CurrentSeal = paladin.SealOfInsightAura
		paladin.SealOfInsightAura.Activate(sim)
	case proto.PaladinSeal_Righteousness:
		paladin.CurrentJudgement = paladin.JudgementOfRighteousness
		paladin.CurrentSeal = paladin.SealOfRighteousnessAura
		paladin.SealOfRighteousnessAura.Activate(sim)
	case proto.PaladinSeal_Justice:
		paladin.CurrentJudgement = paladin.JudgementOfJustice
		paladin.CurrentSeal = paladin.SealOfJusticeAura
		paladin.SealOfJusticeAura.Activate(sim)
	}

	paladin.HolyPowerBar.Reset()
}

func NewPaladin(character *core.Character, talentsStr string, options *proto.PaladinOptions) *Paladin {
	paladin := &Paladin{
		Character:   *character,
		Talents:     &proto.PaladinTalents{},
		Seal:        options.Seal,
		PaladinAura: options.Aura,
	}

	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()
	paladin.initializeHolyPowerBar()

	// Only retribution and holy are actually pets performing some kind of action
	if paladin.Spec != proto.Spec_SpecProtectionPaladin {
		paladin.AncientGuardian = paladin.NewAncientGuardian()
	}

	paladin.EnableAutoAttacks(paladin, core.AutoAttackOptions{
		MainHand:       paladin.WeaponFromMainHand(paladin.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	paladin.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	paladin.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)

	// TODO: figure out the exact tanking stat dependencies for prot pala
	// // Paladins get 0.0167 dodge per agi. ~1% per 59.88
	// paladin.AddStatDependency(stats.Agility, stats.Dodge, (1.0/59.88)*core.DodgeRatingPerDodgeChance)
	// // Paladins get more melee haste from haste than other classes
	// paladin.PseudoStats.MeleeHasteRatingPerHastePercent /= 1.3
	// // Base dodge is unaffected by Diminishing Returns
	// paladin.PseudoStats.BaseDodge += 0.034943
	// paladin.PseudoStats.BaseParry += 0.05

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	return paladin
}
