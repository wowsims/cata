package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	SpellFlagSecondaryJudgement = core.SpellFlagAgentReserved1
	SpellFlagPrimaryJudgement   = core.SpellFlagAgentReserved2
)

const (
	SpellMaskSpecialAttack int64 = 1 << iota

	SpellMaskTemplarsVerdict
	SpellMaskCrusaderStrike
	SpellMaskDivineStorm
	SpellMaskExorcism
	SpellMaskHammerOfWrath
	SpellMaskJudgement
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

	SpellMaskHolyShock
	SpellMaskWordOfGlory

	SpellMaskSealOfTruth
	SpellMaskSealOfInsight
	SpellMaskSealOfRighteousness
	SpellMaskSealOfJustice
)

const SpellMaskSingleTarget = SpellMaskCrusaderStrike | SpellMaskTemplarsVerdict

var TalentTreeSizes = [3]int{20, 20, 20}

type Paladin struct {
	core.Character
	HolyPowerBar

	PaladinAura proto.PaladinAura

	Talents *proto.PaladinTalents

	SharedBuilderCooldown *core.Cooldown // Used for CS/DS

	CurrentSeal      *core.Aura
	CurrentJudgement *core.Spell

	DivinePlea            *core.Spell
	DivineStorm           *core.Spell
	HolyWrath             *core.Spell
	Consecration          *core.Spell
	CrusaderStrike        *core.Spell
	Exorcism              *core.Spell
	HolyShield            *core.Spell
	HammerOfTheRighteous  *core.Spell
	HandOfReckoning       *core.Spell
	ShieldOfRighteousness *core.Spell
	AvengersShield        *core.Spell
	HammerOfWrath         *core.Spell
	AvengingWrath         *core.Spell
	DivineProtection      *core.Spell
	TemplarsVerdict       *core.Spell
	Zealotry              *core.Spell
	Inquisition           *core.Spell
	SealsOfCommand        *core.Spell
	SealOfTruth           *core.Spell

	HolyShieldAura          *core.Aura
	RighteousFuryAura       *core.Aura
	DivinePleaAura          *core.Aura
	SealOfTruthAura         *core.Aura
	SealOfRighteousnessAura *core.Aura
	AvengingWrathAura       *core.Aura
	DivineProtectionAura    *core.Aura
	ForbearanceAura         *core.Aura
	VengeanceAura           *core.Aura
	ZealotryAura            *core.Aura
	InquisitionAura         *core.Aura

	ArtOfWarInstantCast *core.Aura
	DivinePurposeProc   *core.Aura

	SpiritualAttunementMetrics *core.ResourceMetrics

	HasTuralyonsOrLiadrinsBattlegear2Pc bool

	DemonAndUndeadTargetCount int32

	mutualLockoutDPAW *core.Timer
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
	if paladin.PaladinAura == proto.PaladinAura_DevotionAura {
		raidBuffs.DevotionAura = true
	}
	if paladin.PaladinAura == proto.PaladinAura_RetributionAura {
		raidBuffs.RetributionAura = true
	}
	if paladin.PaladinAura == proto.PaladinAura_ResistanceAura {
		raidBuffs.ResistanceAura = true
	}
	if paladin.Talents.Communion {
		raidBuffs.Communion = true
	}
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	paladin.ApplyGlyphs()
	paladin.RegisterSpells()
}

func (paladin *Paladin) RegisterSpells() {
	paladin.RegisterCrusaderStrike()
	paladin.RegisterExorcism()
	paladin.RegisterJudgement()
	paladin.RegisterSealOfTruth()
}

func (paladin *Paladin) Reset(_ *core.Simulation) {
	paladin.CurrentSeal = nil
	paladin.CurrentJudgement = nil
	paladin.HolyPowerBar.Reset()
}

func NewPaladin(character *core.Character, talentsStr string) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
	}

	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr, TalentTreeSizes)

	// // This is used to cache its effect in talents.go
	// paladin.HasTuralyonsOrLiadrinsBattlegear2Pc = paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 2)

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()
	paladin.InitializeHolyPowerBar()

	paladin.SharedBuilderCooldown = &core.Cooldown{
		// TODO: needs to interrogate ret talents for Sanctity of Battle
		// and have this cooldown conditionally be reduced based on haste rating
		Timer:    paladin.NewTimer(),
		Duration: time.Millisecond * 4500,
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

// Shared 30sec cooldown for Divine Protection and Avenging Wrath
func (paladin *Paladin) GetMutualLockoutDPAW() *core.Timer {
	return paladin.Character.GetOrInitTimer(&paladin.mutualLockoutDPAW)
}
