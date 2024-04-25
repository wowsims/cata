package warlock

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var TalentTreeSizes = [3]int{18, 19, 19}

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	Pet *WarlockPet

	ShadowBolt           *core.Spell
	Incinerate           *core.Spell
	Immolate             *core.Spell
	UnstableAffliction   *core.Spell
	Corruption           *core.Spell
	Haunt                *core.Spell
	LifeTap              *core.Spell
	DarkPact             *core.Spell
	ChaosBolt            *core.Spell
	SoulFire             *core.Spell
	Conflagrate          *core.Spell
	DrainSoul            *core.Spell
	Shadowburn           *core.Spell
	SearingPain          *core.Spell
	HandOfGuldan         *core.Spell
	CurseOfElements      *core.Spell
	CurseOfElementsAuras core.AuraArray
	CurseOfWeakness      *core.Spell
	CurseOfWeaknessAuras core.AuraArray
	CurseOfTongues       *core.Spell
	CurseOfTonguesAuras  core.AuraArray
	BaneOfAgony          *core.Spell
	BaneOfDoom           *core.Spell
	BaneOfHavoc          *core.Spell
	Seed                 *core.Spell
	SeedDamageTracker    []float64

	ShadowEmbraceAuras     core.AuraArray
	NightfallProcAura      *core.Aura
	EradicationAura        *core.Aura
	DemonicSoulAura        *core.Aura
	Metamorphosis          *core.Spell
	MetamorphosisAura      *core.Aura
	ImmolationAura         *core.Spell
	HauntDebuffAuras       core.AuraArray
	MoltenCoreAura         *core.Aura
	DecimationAura         *core.Aura
	PyroclasmAura          *core.Aura
	BackdraftAura          *core.Aura
	EmpoweredImpAura       *core.Aura
	SpiritsoftheDamnedAura *core.Aura

	Infernal *InfernalPet
	Inferno  *core.Spell

	petStmBonusSP float64

	FireAndBrimstoneAura core.AuraArray

	ScalingBaseDamage float64
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) ApplyTalents() {
	// Apply Armor Spec
	warlock.EnableArmorSpecialization(stats.Intellect, proto.ArmorType_ArmorTypeCloth)

	warlock.ApplyAfflictionTalents()
	warlock.ApplyDemonologyTalents()
	warlock.ApplyDestructionTalents()

	warlock.ApplyGlyphs()
}

func (warlock *Warlock) Initialize() {

	// base scaling value for a level 85 priest
	warlock.ScalingBaseDamage = 962.335630

	warlock.registerIncinerateSpell()
	warlock.registerShadowBoltSpell()
	warlock.registerImmolateSpell()
	warlock.registerCorruptionSpell()
	warlock.registerCurseOfElementsSpell()
	warlock.registerCurseOfWeaknessSpell()
	warlock.registerCurseOfTonguesSpell()
	warlock.registerBaneOfAgonySpell()
	warlock.registerBaneOfDoomSpell()
	warlock.registerLifeTapSpell()
	warlock.registerSeedSpell()
	warlock.registerSoulFireSpell()
	warlock.registerDrainSoulSpell()
	warlock.registerMetamorphosisSpell()
	// warlock.registerShadowBurnSpell()
	warlock.registerSearingPainSpell()
	// warlock.registerInfernoSpell()
	// warlock.registerBlackBook()

	// // Do this post-finalize so cast speed is updated with new stats
	// warlock.Env.RegisterPostFinalizeEffect(func() {
	// 	// if itemswap is enabled, correct for any possible haste changes
	// 	var correction stats.Stats
	// 	if warlock.ItemSwap.IsEnabled() {
	// 		correction = warlock.ItemSwap.CalcStatChanges([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand,
	// 			proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged})

	// 		warlock.AddStats(correction)
	// 		warlock.MultiplyCastSpeed(1.0)
	// 	}

	// 	if warlock.Options.Summon != proto.WarlockOptions_NoSummon && warlock.Talents.DemonicKnowledge > 0 {
	// 		warlock.RegisterPrepullAction(-999*time.Second, func(sim *core.Simulation) {
	// 			// TODO: investigate a better way of handling this like a "reverse inheritance" for pets.
	// 			// TODO: this will break if we ever get stamina/intellect from procs, but there aren't
	// 			// many such effects and none that we care about
	// 			bonus := (warlock.Pet.GetStat(stats.Stamina) + warlock.Pet.GetStat(stats.Intellect)) *
	// 				(0.04 * float64(warlock.Talents.DemonicKnowledge))
	// 			if bonus != warlock.petStmBonusSP {
	// 				warlock.AddStatDynamic(sim, stats.SpellPower, bonus-warlock.petStmBonusSP)
	// 				warlock.petStmBonusSP = bonus
	// 			}
	// 		})
	// 	}
	// })
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	// raidBuffs.BloodPact = max(raidBuffs.BloodPact, core.MakeTristateValue(
	// 	warlock.Options.Summon == proto.WarlockOptions_Imp,
	// 	warlock.Talents.ImprovedImp == 2,
	// ))

	// raidBuffs.FelIntelligence = max(raidBuffs.FelIntelligence, core.MakeTristateValue(
	// 	warlock.Options.Summon == proto.WarlockOptions_Felhunter,
	// 	warlock.Talents.ImprovedFelhunter == 2,
	// ))
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	if sim.CurrentTime == 0 {
		warlock.petStmBonusSP = 0
	}
}

func NewWarlock(character *core.Character, options *proto.Player, warlockOptions *proto.WarlockOptions) *Warlock {
	warlock := &Warlock{
		Character: *character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions,
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	warlock.EnableManaBar()

	// warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)

	// if warlock.Options.Armor == proto.WarlockOptions_FelArmor {
	// 	demonicAegisMultiplier := 1 + float64(warlock.Talents.DemonicAegis)*0.1
	// 	amount := 180.0 * demonicAegisMultiplier
	// 	warlock.AddStat(stats.SpellPower, amount)
	// 	warlock.AddStatDependency(stats.Spirit, stats.SpellPower, 0.3*demonicAegisMultiplier)
	// }

	// if warlock.Options.Summon != proto.WarlockOptions_NoSummon {
	// 	warlock.Pet = warlock.NewWarlockPet()
	// }

	// warlock.Infernal = warlock.NewInfernal()

	return warlock
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}

func (warlock *Warlock) HasPrimeGlyph(glyph proto.WarlockPrimeGlyph) bool {
	return warlock.HasGlyph(int32(glyph))
}

func (warlock *Warlock) HasMajorGlyph(glyph proto.WarlockMajorGlyph) bool {
	return warlock.HasGlyph(int32(glyph))
}

func (warlock *Warlock) HasMinorGlyph(glyph proto.WarlockMinorGlyph) bool {
	return warlock.HasGlyph(int32(glyph))
}

const (
	WarlockSpellFlagNone    int64 = 0
	WarlockSpellConflagrate int64 = 1 << iota
	WarlockSpellShadowBolt
	WarlockSpellChaosBolt
	WarlockSpellImmolate
	WarlockSpellIncinerate
	WarlockSpellSoulFire
	WarlockSpellLifeTap
	WarlockSpellCorruption
	WarlockSpellHaunt
	WarlockSpellUnstableAffliction
	WarlockSpellCurseOfElements
	WarlockSpellCurseOfWeakness
	WarlockSpellCurseOfTongues
	WarlockSpellBaneOfAgony
	WarlockSpellBaneOfDoom
	WarlockSpellDrainSoul
	WarlockSpellDrainLife
	WarlockSpellMetamorphosis
	WarlockSpellSuccubusLashOfPain
	WarlockSpellFelGuardLegionStrike
	WarlockSpellFelGuardFelstorm
	WarlockSpellImpFireBolt
	WarlockSpellFelHunterShadowBite
	WarlockSpellSeedOfCorruption
	WarlockSpellSeedOfCorruptionExposion
	WarlockSpellHandOfGuldan
	WarlockSpellImmolationAura

	WarlockShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellHaunt |
		WarlockSpellDrainSoul | WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony |
		WarlockSpellShadowBolt | WarlockSpellSeedOfCorruptionExposion | WarlockSpellHandOfGuldan

	WarlockPeriodicShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony

	WarlockFireDamage = WarlockSpellConflagrate | WarlockSpellImmolate | WarlockSpellIncinerate | WarlockSpellSoulFire |
		WarlockSpellImmolationAura | WarlockSpellHandOfGuldan
)

func (warlock *Warlock) calcBaseDamage(sim *core.Simulation, coefficient float64, variance float64) float64 {
	baseDamage := warlock.ScalingBaseDamage * coefficient
	if variance > 0 {
		delta := warlock.ScalingBaseDamage * variance * 0.5
		baseDamage += sim.Roll(-delta, delta)
	}

	return baseDamage
}
