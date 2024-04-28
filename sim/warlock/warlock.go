package warlock

import (
	"math"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	PetExpertiseScale = 1.53 * core.ExpertisePerQuarterPercentReduction
)

var TalentTreeSizes = [3]int{18, 19, 19}

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	ShadowBolt           *core.Spell
	Incinerate           *core.Spell
	Immolate             *core.Spell
	ImmolateDot          *core.Spell
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
	BurningEmbers        *core.Spell
	DemonSoul            *core.Spell

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

	SummonFelhunter *core.Spell
	SummonFelguard  *core.Spell
	SummonImp       *core.Spell
	SummonSuccubus  *core.Spell
	SummonDoomguard *core.Spell
	SummonInfernal  *core.Spell

	ActivePet string
	Felhunter *FelhunterPet
	Felguard  *FelguardPet
	Imp       *ImpPet
	Succubus  *SuccubusPet
	Doomguard *DoomguardPet
	Infernal  *InfernalPet
	EbonImp   *EbonImpPet

	SummonGuardianTimer *core.Timer

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

	// base scaling value for a level 85 warlock
	warlock.ScalingBaseDamage = 962.335630
	warlock.SummonGuardianTimer = warlock.NewTimer()

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
	warlock.registerSearingPainSpell()
	warlock.registerSummonInfernalSpell(warlock.SummonGuardianTimer)
	warlock.registerSummonDoomguardSpell(warlock.SummonGuardianTimer)
	warlock.registerSummonImpSpell()
	warlock.registerSummonFelHunterSpell()
	warlock.registerSummonSuccubusSpell()
	warlock.registerDemonSoulSpell()

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
	switch warlock.Options.Summon {
	case proto.WarlockOptions_Felguard:
		warlock.ChangeActivePet(sim, PetFelguard)
	case proto.WarlockOptions_Felhunter:
		warlock.ChangeActivePet(sim, PetFelhunter)
	case proto.WarlockOptions_Imp:
		warlock.ChangeActivePet(sim, PetImp)
	case proto.WarlockOptions_Succubus:
		warlock.ChangeActivePet(sim, PetSuccubus)
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

	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)

	if warlock.Options.Armor == proto.WarlockOptions_FelArmor {
		warlock.AddStat(stats.SpellPower, 638)
	}

	warlock.EbonImp = warlock.NewEbonImp()
	warlock.Infernal = warlock.NewInfernalPet()
	warlock.Doomguard = warlock.NewDoomguardPet()
	warlock.Felguard = warlock.NewFelguardPet()
	warlock.Felhunter = warlock.NewFelhunterPet()
	warlock.Imp = warlock.NewImpPet()
	warlock.Succubus = warlock.NewSuccubusPet()

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
	WarlockSpellBurningEmbers
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

	WarlockShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellHaunt |
		WarlockSpellDrainSoul | WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony |
		WarlockSpellShadowBolt | WarlockSpellSeedOfCorruptionExposion | WarlockSpellHandOfGuldan

	WarlockPeriodicShadowDamage = WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDrainSoul |
		WarlockSpellDrainLife | WarlockSpellBaneOfDoom | WarlockSpellBaneOfAgony

	WarlockFireDamage = WarlockSpellConflagrate | WarlockSpellImmolate | WarlockSpellIncinerate | WarlockSpellSoulFire |
		WarlockSpellImmolationAura | WarlockSpellHandOfGuldan | WarlockSpellSearingPain | WarlockSpellImmolateDot
)

func (warlock *Warlock) CalcBaseDamage(coefficient float64) float64 {
	return warlock.ScalingBaseDamage * coefficient
}

func (warlock *Warlock) CalcBaseDamageWithVariance(sim *core.Simulation, coefficient float64, variance float64) float64 {
	baseDamage := warlock.ScalingBaseDamage * coefficient
	if variance > 0 {
		delta := warlock.ScalingBaseDamage * variance * 0.5
		baseDamage += sim.Roll(-delta, delta)
	}

	return baseDamage
}

func (warlock *Warlock) MakeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		// based on testing for WotLK Classic the following is true:
		// - pets are meele hit capped if and only if the warlock has 210 (8%) spell hit rating or more
		//   - this is unaffected by suppression and by magic hit debuffs like FF
		// - pets gain expertise from 0% to 6.5% relative to the owners hit, reaching cap at 17% spell hit
		//   - this is also unaffected by suppression and by magic hit debuffs like FF
		//   - this is continious, i.e. not restricted to 0.25 intervals
		// - pets gain spell hit from 0% to 17% relative to the owners hit, reaching cap at 12% spell hit
		// spell hit rating is floor'd
		//   - affected by suppression and ff, but in weird ways:
		// 3/3 suppression => 262 hit  (9.99%) results in misses, 263 (10.03%) no misses
		// 2/3 suppression => 278 hit (10.60%) results in misses, 279 (10.64%) no misses
		// 1/3 suppression => 288 hit (10.98%) results in misses, 289 (11.02%) no misses
		// 0/3 suppression => 314 hit (11.97%) results in misses, 315 (12.01%) no misses
		// 3/3 suppression + FF => 209 hit (7.97%) results in misses, 210 (8.01%) no misses
		// 2/3 suppression + FF => 222 hit (8.46%) results in misses, 223 (8.50%) no misses
		//
		// the best approximation of this behaviour is that we scale the warlock's spell hit by `1/12*17` floor
		// the result and then add the hit percent from suppression/ff

		// does correctly not include ff/misery
		ownerHitChance := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance

		// TODO: Account for sunfire/soulfrost
		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      ownerStats[stats.SpellPower] * 0.57,
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			//With demonic tactics gone is there any crit inheritance?
			//stats.SpellCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			//stats.MeleeCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			stats.MeleeHit: ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit: math.Floor(ownerStats[stats.SpellHit] / 12.0 * 17.0),
			// TODO: revisit
			stats.Expertise: (ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance) *
				PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,

			// Resists, 40%
		}
	}
}

func (warlock *Warlock) ChangeActivePet(sim *core.Simulation, newPet string) {
	switch warlock.ActivePet {
	case PetFelguard:
		warlock.Felguard.WarlockPet.Disable(sim)
	case PetFelhunter:
		warlock.Felhunter.WarlockPet.Disable(sim)
	case PetImp:
		warlock.Imp.WarlockPet.Disable(sim)
	case PetSuccubus:
		warlock.Succubus.WarlockPet.Disable(sim)
	}

	switch newPet {
	case PetFelguard:
		warlock.Felguard.WarlockPet.Enable(sim, warlock.Felguard)
	case PetFelhunter:
		warlock.Felhunter.WarlockPet.Enable(sim, warlock.Felhunter)
	case PetImp:
		warlock.Imp.WarlockPet.Enable(sim, warlock.Imp)
	case PetSuccubus:
		warlock.Succubus.WarlockPet.Enable(sim, warlock.Succubus)
	}
}

// func (warlock *Warlock) ChangeActivePet(sim *core.Simulation, newPet *WarlockPet) {
// 		if warlock.ActivePet != nil {
// 			warlock.ActivePet.Disable(sim)
// 		}
// 		warlock.ActivePet = &newPet
// 		warlock.ActivePet.Enable(sim, warlock.ActivePet)
// 	}
