package warlock

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	Corruption           *core.Spell
	CurseOfElementsAuras core.AuraArray
	Immolate             *core.Spell
	Metamorphosis        *core.Spell
	Seed                 *core.Spell
	ShadowEmbraceAuras   core.AuraArray
	Shadowburn           *core.Spell
	Hellfire             *core.Spell
	DrainLife            *core.Spell

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

	serviceTimer *core.Timer

	// Item sets
	T13_4pc      *core.Aura
	T15_2pc      *core.Aura
	T15_4pc      *core.Aura
	T16_2pc_buff *core.Aura
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

	warlock.registerCurseOfElements()
	doomguardInfernalTimer := warlock.NewTimer()
	warlock.registerSummonDoomguard(doomguardInfernalTimer)
	warlock.registerSummonInfernal(doomguardInfernalTimer)
	warlock.registerLifeTap()

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

	warlock.serviceTimer = character.NewTimer()
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
	WarlockSpellAgony
	WarlockSpellDrainSoul
	WarlockSpellDrainLife
	WarlockSpellMetamorphosis
	WarlockSpellSeedOfCorruption
	WarlockSpellSeedOfCorruptionExposion
	WarlockSpellHandOfGuldan
	WarlockSpellHellfire
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
	WarlockSpellDarkSoulKnowledge
	WarlockSpellDarkSoulMisery
	WarlockSpellMaleficGrasp
	WarlockSpellDemonicSlash
	WarlockSpellTouchOfChaos
	WarlockSpellChaosWave
	WarlockSpellCarrionSwarm
	WarlockSpellDoom
	WarlockSpellVoidray
	WarlockSpellAll int64 = 1<<iota - 1

	WarlockShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellHaunt |
		WarlockSpellDrainSoul | WarlockSpellDrainLife | WarlockSpellAgony |
		WarlockSpellShadowBolt | WarlockSpellSeedOfCorruptionExposion | WarlockSpellHandOfGuldan |
		WarlockSpellShadowflame | WarlockSpellFelFlame | WarlockSpellChaosBolt | WarlockSpellShadowBurn

	WarlockPeriodicShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellAgony

	WarlockFireDamage = WarlockSpellConflagrate | WarlockSpellImmolate | WarlockSpellIncinerate | WarlockSpellSoulFire |
		WarlockSpellHandOfGuldan | WarlockSpellSearingPain | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellFelFlame | WarlockSpellChaosBolt | WarlockSpellShadowBurn | WarlockSpellFaBConflagrate |
		WarlockSpellFaBIncinerate

	WarlockDoT = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellAgony | WarlockSpellImmolateDot |
		WarlockSpellShadowflameDot | WarlockSpellBurningEmbers

	WarlockSummonSpells = WarlockSpellSummonImp | WarlockSpellSummonSuccubus | WarlockSpellSummonFelhunter |
		WarlockSpellSummonFelguard

	WarlockDarkSoulSpell             = WarlockSpellDarkSoulInsanity | WarlockSpellDarkSoulKnowledge | WarlockSpellDarkSoulMisery
	WarlockAllSummons                = WarlockSummonSpells | WarlockSpellSummonInfernal | WarlockSpellSummonDoomguard
	WarlockSpellsChaoticEnergyDestro = WarlockSpellAll &^ WarlockAllSummons &^ WarlockSpellDrainLife
)

// Pandemic - For now a Warlock only ability. Might be moved into core support in late expansions
func (warlock *Warlock) ApplyDotWithPandemic(dot *core.Dot, sim *core.Simulation) {

	// if DoT was not active before, there is nothing we need to do for pandemic
	if !dot.IsActive() {
		dot.Apply(sim)
		return
	}

	// MoP Pandemic is a warlock only ability
	// It allows for the extension of up to 50% of the unhasted base duration
	// So we need to determine which is shorter base + remaining or base + maxExtend
	remaining := dot.RemainingDuration(sim)
	extend := time.Duration(math.Min(
		float64(dot.BaseDuration()+remaining),
		float64(dot.BaseDuration()+dot.BaseDuration()/2),
	))

	// First do usual dot carry over
	dot.Apply(sim)
	for dot.RemainingDuration(sim)-dot.TimeUntilNextTick(sim)+dot.TickPeriod() <= extend {
		dot.AddTick()
	}
}

// Called to handle custom resources
type WarlockSpellCastedCallback func(resultList []core.SpellResult, spell *core.Spell, sim *core.Simulation)

type SecondaryResourceCost struct {
	SecondaryCost int
	Name          string
}

// CostFailureReason implements core.ResourceCostImpl.
func (s *SecondaryResourceCost) CostFailureReason(_ *core.Simulation, spell *core.Spell) string {
	return fmt.Sprintf(
		"Not enough %s (Current %s = %0.03f, %s Cost = %0.03f)",
		s.Name,
		s.Name,
		spell.Unit.GetSecondaryResourceBar(),
		s.Name,
		spell.CurCast.Cost,
	)
}

// IssueRefund implements core.ResourceCostImpl.
func (s *SecondaryResourceCost) IssueRefund(sim *core.Simulation, spell *core.Spell) {
	curCost := spell.Cost.PercentModifier * float64(s.SecondaryCost)
	spell.Unit.GetSecondaryResourceBar().Gain(sim, int32(curCost), spell.ActionID)
}

// MeetsRequirement implements core.ResourceCostImpl.
func (s *SecondaryResourceCost) MeetsRequirement(_ *core.Simulation, spell *core.Spell) bool {
	spell.CurCast.Cost = spell.Cost.PercentModifier * float64(s.SecondaryCost)
	return spell.Unit.GetSecondaryResourceBar().CanSpend(int32(spell.CurCast.Cost))
}

// SpendCost implements core.ResourceCostImpl.
func (s *SecondaryResourceCost) SpendCost(sim *core.Simulation, spell *core.Spell) {

	// during some hard casts resourc might tick down, make sure spells don't execute on exhaustion
	if spell.Unit.GetSecondaryResourceBar().CanSpend(int32(spell.CurCast.Cost)) {
		spell.Unit.GetSecondaryResourceBar().Spend(sim, int32(spell.CurCast.Cost), spell.ActionID)
	}
}
