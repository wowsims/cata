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

	SealOfTruth *core.Spell

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
	// raidBuffs.DevotionAura = max(raidBuffs.DevotionAura, core.MakeTristateValue(
	// 	paladin.PaladinAura == proto.PaladinAura_DevotionAura,
	// 	paladin.Talents.ImprovedDevotionAura == 5))

	// if paladin.PaladinAura == proto.PaladinAura_RetributionAura {
	// 	raidBuffs.RetributionAura = true
	// }

	// if paladin.Talents.SanctifiedRetribution {
	// 	raidBuffs.SanctifiedRetribution = true
	// }

	// if paladin.Talents.SwiftRetribution == 3 {
	// 	raidBuffs.SwiftRetribution = paladin.Talents.SwiftRetribution == 3 // TODO: Fix-- though having something between 0/3 and 3/3 is unlikely
	// }

	// TODO: Figure out a way to just start with 1 DG cooldown available without making a redundant Spell
	//if paladin.Talents.DivineGuardian == 2 {
	//	raidBuffs.divineGuardians++
	//}
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	paladin.ApplyGlyphs()
	paladin.RegisterJudgement()
	paladin.RegisterSealOfTruth()
	// // Update auto crit multipliers now that we have the targets.
	// paladin.AutoAttacks.MHConfig().CritMultiplier = paladin.MeleeCritMultiplier()
	// paladin.registerSealOfVengeanceSpellAndAura()
	// paladin.registerSealOfRighteousnessSpellAndAura()
	// paladin.registerSealOfCommandSpellAndAura()
	// // paladin.setupSealOfTheCrusader()
	// // paladin.setupSealOfWisdom()
	// // paladin.setupSealOfLight()
	// // paladin.setupSealOfRighteousness()
	// // paladin.setupJudgementRefresh()

	paladin.RegisterCrusaderStrike()

	// paladin.registerConsecrationSpell()
	// paladin.registerHammerOfWrathSpell()
	// paladin.registerHolyWrathSpell()

	paladin.RegisterExorcism()
	// paladin.registerHolyShieldSpell()
	// paladin.registerHammerOfTheRighteousSpell()
	// paladin.registerHandOfReckoningSpell()
	// paladin.registerShieldOfRighteousnessSpell()
	// paladin.registerAvengersShieldSpell()
	// paladin.registerJudgements()

	// paladin.registerSpiritualAttunement()
	// paladin.registerDivinePleaSpell()
	// paladin.registerDivineProtectionSpell()
	// paladin.registerForbearanceDebuff()

	// for i := int32(0); i < paladin.Env.GetNumTargets(); i++ {
	// 	unit := paladin.Env.GetTargetUnit(i)
	// 	if unit.MobType == proto.MobType_MobTypeDemon || unit.MobType == proto.MobType_MobTypeUndead {
	// 		paladin.DemonAndUndeadTargetCount += 1
	// 	}
	// }
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
