package warlock

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
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
	Metamorphosis        *core.Spell
	Seed                 *core.Spell
	ShadowEmbraceAuras   core.AuraArray
	Shadowburn           *core.Spell
	UnstableAffliction   *core.Spell

	ActivePet *WarlockPet
	Felhunter *WarlockPet
	// Felguard  *WarlockPet
	Imp        *WarlockPet
	Succubus   *WarlockPet
	Voidwalker *WarlockPet

	Doomguard *DoomguardPet
	Infernal  *InfernalPet
	// EbonImp   *EbonImpPet
	FieryImp *FieryImpPet

	// Item sets
	T13_4pc *core.Aura
	T15_2pc *core.Aura
	T15_4pc *core.Aura
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) ApplyTalents() {
	warlock.registerHarvestLife()
	warlock.registerArchimondesDarkness()
	warlock.registerKilJaedensCunning()
	warlock.registerMannarothsFury()
	warlock.registerGrimoireOfSupremacy()
	warlock.registerGrimoireOfSacrifice()
}

func (warlock *Warlock) Initialize() {

	warlock.registerDarkSoulInstability()
	warlock.registerCurseOfElements()
	warlock.registerDrainLife()

	// warlock.registerBaneOfAgony()
	// warlock.registerBaneOfDoom()
	// warlock.registerCorruption()
	// warlock.registerCurseOfElements()
	// warlock.registerCurseOfTongues()
	// warlock.registerCurseOfWeakness()
	// warlock.registerDemonSoul()
	// warlock.registerDrainLife()
	// warlock.registerDrainSoul()
	// warlock.registerFelFlame()
	// warlock.registerImmolate()
	// warlock.registerIncinerate()
	// warlock.registerLifeTap()
	// warlock.registerSearingPain()
	// warlock.registerSeed()
	// warlock.registerShadowBolt()
	// warlock.registerShadowflame()
	// warlock.registerSoulFire()
	// warlock.registerSoulHarvest()
	// warlock.registerSoulburn()
	// warlock.registerSummonDemon()

	doomguardInfernalTimer := warlock.NewTimer()
	warlock.registerSummonDoomguard(doomguardInfernalTimer)
	warlock.registerSummonInfernal(doomguardInfernalTimer)

	// Fel Armor 10% Stamina
	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label:    "Fel Armor",
			ActionID: core.ActionID{SpellID: 104938},
		}))
	warlock.MultiplyStat(stats.Stamina, 1.1)
	warlock.MultiplyStat(stats.Health, 1.1)

	// 5% int passive
	warlock.MultiplyStat(stats.Intellect, 1.05)
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {

}

func (warlock *Warlock) Reset(sim *core.Simulation) {
}

func NewWarlock(character *core.Character, options *proto.Player, warlockOptions *proto.WarlockOptions) *Warlock {
	warlock := &Warlock{
		Character: *character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions,
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString)
	warlock.EnableManaBar()
	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)

	// warlock.EbonImp = warlock.NewEbonImp()
	warlock.Infernal = warlock.NewInfernalPet()
	warlock.Doomguard = warlock.NewDoomguardPet()
	warlock.FieryImp = warlock.NewFieryImp()

	warlock.registerPets()
	warlock.registerGrimoireOfService()

	return warlock
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
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
	WarlockSpellFaBConflagrate
	WarlockSpellShadowBolt
	WarlockSpellChaosBolt
	WarlockSpellImmolate
	WarlockSpellImmolateDot
	WarlockSpellIncinerate
	WarlockSpellFaBIncinerate
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
	WarlockSpellHellfire
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
	WarlockSpellVoidwalkerTorment
	WarlockSpellSummonInfernal
	WarlockSpellDemonSoul
	WarlockSpellShadowflame
	WarlockSpellShadowflameDot
	WarlockSpellSoulBurn
	WarlockSpellFelFlame
	WarlockSpellBurningEmbers
	WarlockSpellEmberTap
	WarlockSpellRainOfFire
	WarlockSpellFireAndBrimstone
	WarlockSpellDarkSoulInsanity
	WarlockSpellMaleficGrasp
	WarlockSpellDemonicSlash
	WarlockSpellTouchOfChaos
	WarlockSpellAll int64 = 1<<iota - 1

	WarlockShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellHaunt |
		WarlockSpellDrainSoul | WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony |
		WarlockSpellShadowBolt | WarlockSpellSeedOfCorruptionExposion | WarlockSpellHandOfGuldan |
		WarlockSpellShadowflame | WarlockSpellFelFlame | WarlockSpellChaosBolt | WarlockSpellShadowBurn

	WarlockPeriodicShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony

	WarlockFireDamage = WarlockSpellConflagrate | WarlockSpellImmolate | WarlockSpellIncinerate | WarlockSpellSoulFire |
		WarlockSpellImmolationAura | WarlockSpellHandOfGuldan | WarlockSpellSearingPain | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellFelFlame | WarlockSpellChaosBolt | WarlockSpellShadowBurn | WarlockSpellFaBConflagrate |
		WarlockSpellFaBIncinerate

	WarlockDoT = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellBurningEmbers

	WarlockSummonSpells = WarlockSpellSummonImp | WarlockSpellSummonSuccubus | WarlockSpellSummonFelhunter |
		WarlockSpellSummonFelguard

	WarlockDarkSoulSpell             = WarlockSpellDarkSoulInsanity
	WarlockAllSummons                = WarlockSummonSpells | WarlockSpellSummonInfernal | WarlockSpellSummonDoomguard
	WarlockSpellsChaoticEnergyDestro = WarlockSpellAll &^ WarlockAllSummons
)

const (
	PetExpertiseScale = 1.53 * core.ExpertisePerQuarterPercentReduction
)
