package warlock

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	BaneOfAgony          *core.Spell
	BaneOfDoom           *core.Spell
	BurningEmbers        *core.Spell
	Corruption           *core.Spell
	CurseOfElementsAuras core.AuraArray
	CurseOfTonguesAuras  core.AuraArray
	CurseOfWeaknessAuras core.AuraArray
	HauntDebuffAuras     core.AuraArray
	Immolate             *core.Spell
	ImmolateDot          *core.Spell
	Metamorphosis        *core.Spell
	Seed                 *core.Spell
	ShadowEmbraceAuras   core.AuraArray
	Shadowburn           *core.Spell
	UnstableAffliction   *core.Spell

	Felhunter *WarlockPet
	Felguard  *WarlockPet
	Imp       *WarlockPet
	Succubus  *WarlockPet

	Doomguard *DoomguardPet
	Infernal  *InfernalPet
	EbonImp   *EbonImpPet

	SoulShards   int32
	SoulBurnAura *core.Aura
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) ApplyTalents() {
	warlock.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeCloth)

	warlock.ApplyAfflictionTalents()
	warlock.ApplyDemonologyTalents()
	warlock.ApplyDestructionTalents()

	warlock.ApplyGlyphs()
}

func (warlock *Warlock) Initialize() {
	warlock.registerBaneOfAgony()
	warlock.registerBaneOfDoom()
	warlock.registerCorruption()
	warlock.registerCurseOfElements()
	warlock.registerCurseOfTongues()
	warlock.registerCurseOfWeakness()
	warlock.registerDemonSoul()
	warlock.registerDrainLife()
	warlock.registerDrainSoul()
	warlock.registerFelFlame()
	warlock.registerImmolate()
	warlock.registerIncinerate()
	warlock.registerLifeTap()
	warlock.registerSearingPain()
	warlock.registerSeed()
	warlock.registerShadowBolt()
	warlock.registerShadowflame()
	warlock.registerSoulFire()
	warlock.registerSoulburn()
	warlock.registerSummonFelHunter()
	warlock.registerSummonImp()
	warlock.registerSummonSuccubus()

	doomguardInfernalTimer := warlock.NewTimer()
	warlock.registerSummonDoomguard(doomguardInfernalTimer)
	warlock.registerSummonInfernal(doomguardInfernalTimer)

	// TODO: vile hack to make the APLs work for now ...
	if !warlock.HasSetBonus(ItemSetMaleficRaiment, 4) {
		warlock.RegisterAura(core.Aura{
			Label:    "Fel Spark",
			ActionID: core.ActionID{SpellID: 89937},
		})
	}

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Fel Armor",
			ActionID: core.ActionID{SpellID: 28176},
		}))

	warlock.registerPetAbilities()

	// warlock.registerBlackBook()

	// Do this post-finalize so cast speed is updated with new stats
	warlock.Env.RegisterPostFinalizeEffect(func() {
		// if itemswap is enabled, correct for any possible haste changes
		var correction stats.Stats
		if warlock.ItemSwap.IsEnabled() {
			correction = warlock.ItemSwap.CalcStatChanges([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand,
				proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged})

			warlock.AddStats(correction)
			warlock.MultiplyCastSpeed(1.0)
		}
	})
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = warlock.Options.Summon == proto.WarlockOptions_Imp
	raidBuffs.FelIntelligence = warlock.Options.Summon == proto.WarlockOptions_Felhunter
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	warlock.SoulShards = 4
}

func NewWarlock(character *core.Character, options *proto.Player, warlockOptions *proto.WarlockOptions) *Warlock {
	warlock := &Warlock{
		Character: *character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions,
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString, [3]int{18, 19, 19})
	warlock.EnableManaBar()

	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)

	// Add Fel Armor SP by default
	warlock.AddStat(stats.SpellPower, 638)

	warlock.EbonImp = warlock.NewEbonImp()
	warlock.Infernal = warlock.NewInfernalPet()
	warlock.Doomguard = warlock.NewDoomguardPet()

	warlock.registerPets()

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
	WarlockSpellImmolateDot
	WarlockSpellIncinerate
	WarlockSpellSoulFire
	WarlockSpellShadowBurn
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
	WarlockSpellSeedOfCorruption
	WarlockSpellSeedOfCorruptionExposion
	WarlockSpellHandOfGuldan
	WarlockSpellImmolationAura
	WarlockSpellSearingPain
	WarlockSpellSummonDoomguard
	WarlockSpellDoomguardDoomBolt
	WarlockSpellSummonFelguard
	WarlockSpellFelGuardLegionStrike
	WarlockSpellFelGuardFelstorm
	WarlockSpellSummonImp
	WarlockSpellImpFireBolt
	WarlockSpellSummonFelhunter
	WarlockSpellFelHunterShadowBite
	WarlockSpellSummonSuccubus
	WarlockSpellSuccubusLashOfPain
	WarlockSpellDemonSoul
	WarlockSpellShadowflame
	WarlockSpellShadowflameDot
	WarlockSpellSoulBurn
	WarlockSpellFelFlame
	WarlockSpellBurningEmbers

	WarlockShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellHaunt |
		WarlockSpellDrainSoul | WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony |
		WarlockSpellShadowBolt | WarlockSpellSeedOfCorruptionExposion | WarlockSpellHandOfGuldan |
		WarlockSpellShadowflame | WarlockSpellFelFlame

	WarlockPeriodicShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony

	WarlockFireDamage = WarlockSpellConflagrate | WarlockSpellImmolate | WarlockSpellIncinerate | WarlockSpellSoulFire |
		WarlockSpellImmolationAura | WarlockSpellHandOfGuldan | WarlockSpellSearingPain | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellFelFlame

	WarlockDoT = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellBurningEmbers

	WarlockBasicPetSpells = WarlockSpellFelGuardLegionStrike | WarlockSpellSuccubusLashOfPain |
		WarlockSpellSuccubusLashOfPain | WarlockSpellFelHunterShadowBite | WarlockSpellImpFireBolt

	WarlockSummonSpells = WarlockSpellSummonImp | WarlockSpellSummonSuccubus | WarlockSpellSummonFelhunter |
		WarlockSpellSummonFelguard
)

const (
	PetExpertiseScale = 1.53 * core.ExpertisePerQuarterPercentReduction
)

const Coefficient_BaneOfDoom float64 = 2.024
const Coefficient_Immolate float64 = 0.692
const Coefficient_ImmolateDot float64 = 0.43900001049
const Coefficient_SeedExplosion float64 = 2.113
const Coefficient_SeedDot float64 = 0.3024
const Coefficient_ChaosBolt float64 = 1.547
const Coefficient_Infernal float64 = 0.485
const Coefficient_ShadowBolt float64 = 0.62
const Coefficient_HandOfGuldan float64 = 1.593
const Coefficient_Incinerate float64 = 0.573
const Coefficient_Shadowburn float64 = 0.714
const Coefficient_Shadowflame float64 = 0.72699999809
const Coefficient_ShadowflameDot float64 = 0.16899999976

const Variance_ChaosBolt float64 = 0.238
const Variance_ShadowBolt float64 = 0.1099999994
const Variance_HandOfGuldan float64 = 0.166
const Variance_Infernal float64 = 0.119
const Variance_Incinerate float64 = 0.15
const Variance_Shadowburn float64 = 0.1099999994
const Variance_Shadowflame float64 = 0.09000000358

func (warlock *Warlock) DefaultSpellCritMultiplier() float64 {
	return warlock.SpellCritMultiplier(1.33, 0.0)
}
