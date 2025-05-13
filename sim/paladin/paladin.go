package paladin

import (
	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Paladin struct {
	core.Character

	Seal      proto.PaladinSeal
	HolyPower core.SecondaryResourceBar

	Talents *proto.PaladinTalents

	// Used for CS/HotR
	sharedBuilderTimer *core.Timer

	CurrentSeal       *core.Aura
	StartingHolyPower int32

	// Pets
	AncientGuardian    *AncientGuardianPet
	GurthalakTentacles []*cata.TentacleOfTheOldOnesPet

	Judgment       *core.Spell
	Exorcism       *core.Spell
	AvengersShield *core.Spell
	HammerOfWrath  *core.Spell

	RighteousFuryAura       *core.Aura
	SealOfTruthAura         *core.Aura
	SealOfInsightAura       *core.Aura
	SealOfRighteousnessAura *core.Aura
	SealOfJusticeAura       *core.Aura
	AvengingWrathAura       *core.Aura
	DivineProtectionAura    *core.Aura
	GoakAura                *core.Aura
	AncientPowerAura        *core.Aura
	DivinePurposeAura       *core.Aura

	// Item sets
	T11Ret4pc *core.Aura

	HolyAvengerActionIDFilter []*core.ActionID

	DynamicHolyPowerSpent                        int32
	BastionOfGloryMultiplier                     float64
	ShieldOfTheRighteousAdditiveMultiplier       float64
	ShieldOfTheRighteousMultiplicativeMultiplier float64
}

func (paladin *Paladin) GetTentacles() []*cata.TentacleOfTheOldOnesPet {
	return paladin.GurthalakTentacles
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

func (paladin *Paladin) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BlessingOfMight = true
	raidBuffs.BlessingOfKings = true
}

func (paladin *Paladin) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (paladin *Paladin) Initialize() {
	paladin.applyGlyphs()
	paladin.registerSpells()
	paladin.addCataclysmPvpGloves()
	paladin.addMistsPvpGloves()
	paladin.applySanctityOfBattle()
}

func (paladin *Paladin) registerSpells() {
	paladin.registerAvengingWrath()
	paladin.registerCrusaderStrike()
	paladin.registerDevotionAura()
	paladin.registerDivineProtectionSpell()
	paladin.registerGuardianOfAncientKings()
	paladin.registerHammerOfTheRighteous()
	paladin.registerHammerOfWrathSpell()
	paladin.registerJudgment()
	paladin.registerSealOfInsight()
	paladin.registerSealOfRighteousness()
	paladin.registerSealOfTruth()
	paladin.registerShieldOfTheRighteous()
	paladin.registerTemplarsVerdict()
	paladin.registerWordOfGlory()
}

func (paladin *Paladin) Reset(sim *core.Simulation) {
	switch paladin.Seal {
	case proto.PaladinSeal_Truth:
		paladin.CurrentSeal = paladin.SealOfTruthAura
		paladin.SealOfTruthAura.Activate(sim)
	case proto.PaladinSeal_Insight:
		paladin.CurrentSeal = paladin.SealOfInsightAura
		paladin.SealOfInsightAura.Activate(sim)
	case proto.PaladinSeal_Righteousness:
		paladin.CurrentSeal = paladin.SealOfRighteousnessAura
		paladin.SealOfRighteousnessAura.Activate(sim)
	case proto.PaladinSeal_Justice:
		paladin.CurrentSeal = paladin.SealOfJusticeAura
		paladin.SealOfJusticeAura.Activate(sim)
	}
}

func NewPaladin(character *core.Character, talentsStr string, options *proto.PaladinOptions) *Paladin {
	paladin := &Paladin{
		Character: *character,
		Talents:   &proto.PaladinTalents{},
		Seal:      options.Seal,
	}

	core.FillTalentsProto(paladin.Talents.ProtoReflect(), talentsStr)

	paladin.PseudoStats.CanParry = true

	paladin.EnableManaBar()
	paladin.HolyPower = paladin.NewDefaultSecondaryResourceBar(core.SecondaryResourceConfig{
		Type:    proto.SecondaryResourceType_SecondaryResourceTypeHolyPower,
		Max:     5,
		Default: paladin.StartingHolyPower,
	})
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

	paladin.PseudoStats.BaseBlockChance += 0.05
	paladin.PseudoStats.BaseDodgeChance += 0.05
	paladin.PseudoStats.BaseParryChance += 0.05

	// Bonus Armor and Armor are treated identically for Paladins
	paladin.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	if mh := paladin.MainHand(); mh.Name == "Gurthalak, Voice of the Deeps" {
		paladin.GurthalakTentacles = make([]*cata.TentacleOfTheOldOnesPet, 10)

		for i := 0; i < 10; i++ {
			paladin.GurthalakTentacles[i] = paladin.NewTentacleOfTheOldOnesPet()
		}
	}

	return paladin
}

func (paladin *Paladin) applySanctityOfBattle() {
	var classMask int64
	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		classMask = SpellMaskSanctityOfBattleProt
	} else if paladin.Spec == proto.Spec_SpecHolyPaladin {
		classMask = SpellMaskSanctityOfBattleHoly
	} else {
		classMask = SpellMaskSanctityOfBattleRet
	}

	cooldownMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Multiplier,
		ClassMask: classMask,
	})

	updateFloatValue := func(castSpeed float64) {
		cooldownMod.UpdateFloatValue(castSpeed)
	}

	paladin.AddOnCastSpeedChanged(func(_ float64, castSpeed float64) {
		updateFloatValue(castSpeed)
	})

	core.MakePermanent(paladin.GetOrRegisterAura(core.Aura{
		Label:    "Sanctity of Battle",
		ActionID: core.ActionID{SpellID: 25956},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			updateFloatValue(paladin.CastSpeed)
			cooldownMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cooldownMod.Deactivate()
		},
	}))
}

func (paladin *Paladin) CanTriggerHolyAvengerHpGain(actionID core.ActionID) {
	paladin.HolyAvengerActionIDFilter = append(paladin.HolyAvengerActionIDFilter, &actionID)
}

// Shared cooldown for CS and HotR
func (paladin *Paladin) BuilderCooldown() *core.Timer {
	return paladin.Character.GetOrInitTimer(&paladin.sharedBuilderTimer)
}

func (paladin *Paladin) SpendableHolyPower() int32 {
	return min(paladin.HolyPower.Value(), 3)
}
