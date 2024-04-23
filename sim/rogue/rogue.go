package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	SpellFlagBuilder     = core.SpellFlagAgentReserved2
	SpellFlagFinisher    = core.SpellFlagAgentReserved3
	SpellFlagColdBlooded = core.SpellFlagAgentReserved4
)

var TalentTreeSizes = [3]int{19, 19, 19}

const RogueBaseDamageScalar = 1125.23
const RogueBleedTag = "RogueBleed"

type Rogue struct {
	core.Character

	Talents              *proto.RogueTalents
	Options              *proto.RogueOptions
	AssassinationOptions *proto.AssassinationRogue_Options
	CombatOptions        *proto.CombatRogue_Options
	SubtletyOptions      *proto.SubtletyRogue_Options

	SliceAndDiceBonus        float64
	AdditiveEnergyRegenBonus float64

	sliceAndDiceDurations [6]time.Duration
	exposeArmorDurations  [6]time.Duration

	Backstab         *core.Spell
	BladeFlurry      *core.Spell
	DeadlyPoison     *core.Spell
	FanOfKnives      *core.Spell
	Feint            *core.Spell
	Garrote          *core.Spell
	Ambush           *core.Spell
	Hemorrhage       *core.Spell
	GhostlyStrike    *core.Spell
	HungerForBlood   *core.Spell
	InstantPoison    [4]*core.Spell
	WoundPoison      [4]*core.Spell
	Mutilate         *core.Spell
	MutilateMH       *core.Spell
	MutilateOH       *core.Spell
	Shiv             *core.Spell
	SinisterStrike   *core.Spell
	TricksOfTheTrade *core.Spell
	Shadowstep       *core.Spell
	Preparation      *core.Spell
	Premeditation    *core.Spell
	ShadowDance      *core.Spell
	ColdBlood        *core.Spell
	Vanish           *core.Spell
	VenomousWounds   *core.Spell
	Vendetta         *core.Spell
	RevealingStrike  *core.Spell
	KillingSpree     *core.Spell
	AdrenalineRush   *core.Spell
	Gouge            *core.Spell

	Envenom      *core.Spell
	Eviscerate   *core.Spell
	ExposeArmor  *core.Spell
	Rupture      *core.Spell
	SliceAndDice *core.Spell
	Recuperate   *core.Spell

	lastDeadlyPoisonProcMask core.ProcMask

	deadlyPoisonProcChanceBonus float64
	instantPoisonPPMM           core.PPMManager
	woundPoisonPPMM             core.PPMManager

	AdrenalineRushAura   *core.Aura
	BladeFlurryAura      *core.Aura
	EnvenomAura          *core.Aura
	ExposeArmorAuras     core.AuraArray
	HungerForBloodAura   *core.Aura
	KillingSpreeAura     *core.Aura
	OverkillAura         *core.Aura
	SliceAndDiceAura     *core.Aura
	RecuperateAura       *core.Aura
	MasterOfSubtletyAura *core.Aura
	ShadowstepAura       *core.Aura
	ShadowDanceAura      *core.Aura
	DirtyDeedsAura       *core.Aura
	HonorAmongThieves    *core.Aura
	StealthAura          *core.Aura
	BanditsGuileAura     *core.Aura
	RestlessBladesAura   *core.Aura

	MasterPoisonerDebuffAuras core.AuraArray
	SavageCombatDebuffAuras   core.AuraArray
	WoundPoisonDebuffAuras    core.AuraArray

	generatorCostModifier      func(float64) float64
	finishingMoveEffectApplier func(sim *core.Simulation, numPoints int32)
}

func (rogue *Rogue) GetCharacter() *core.Character {
	return &rogue.Character
}

func (rogue *Rogue) GetRogue() *Rogue {
	return rogue
}

func (rogue *Rogue) AddRaidBuffs(_ *proto.RaidBuffs)   {}
func (rogue *Rogue) AddPartyBuffs(_ *proto.PartyBuffs) {}

// Apply the effect of successfully casting a finisher to combo points
func (rogue *Rogue) ApplyFinisher(sim *core.Simulation, spell *core.Spell) {
	numPoints := rogue.ComboPoints()
	rogue.SpendComboPoints(sim, spell.ComboPointMetrics())
	rogue.finishingMoveEffectApplier(sim, numPoints)
}

func (rogue *Rogue) HasPrimeGlyph(glyph proto.RoguePrimeGlyph) bool {
	return rogue.HasGlyph(int32(glyph))
}

func (rogue *Rogue) HasMajorGlyph(glyph proto.RogueMajorGlyph) bool {
	return rogue.HasGlyph(int32(glyph))
}

func (rogue *Rogue) HasMinorGlyph(glyph proto.RogueMinorGlyph) bool {
	return rogue.HasGlyph(int32(glyph))
}

func (rogue *Rogue) GetGeneratorCostModifier(cost float64) float64 {
	return rogue.generatorCostModifier(cost)
}

func (rogue *Rogue) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	rogue.AutoAttacks.MHConfig().CritMultiplier = rogue.MeleeCritMultiplier(false)
	rogue.AutoAttacks.OHConfig().CritMultiplier = rogue.MeleeCritMultiplier(false)
	rogue.AutoAttacks.RangedConfig().CritMultiplier = rogue.MeleeCritMultiplier(false)

	rogue.generatorCostModifier = rogue.makeGeneratorCostModifier()

	rogue.registerStealthAura()
	rogue.registerVanishSpell()
	rogue.registerFeintSpell()
	rogue.registerAmbushSpell()
	rogue.registerGarrote()
	rogue.registerSinisterStrikeSpell()
	rogue.registerBackstabSpell()
	rogue.registerRupture()
	rogue.registerSliceAndDice()
	rogue.registerEviscerate()
	rogue.registerEnvenom()
	rogue.registerExposeArmorSpell()
	rogue.registerRecuperate()
	rogue.registerFanOfKnives()
	rogue.registerTricksOfTheTradeSpell()
	rogue.registerDeadlyPoisonSpell()
	rogue.registerInstantPoisonSpell()
	rogue.registerWoundPoisonSpell()
	rogue.registerPoisonAuras()
	rogue.registerShivSpell()
	rogue.registerThistleTeaCD()
	rogue.registerGougeSpell()

	rogue.finishingMoveEffectApplier = rogue.makeFinishingMoveEffectApplier()

	rogue.SliceAndDiceBonus = 0.4
}

func (rogue *Rogue) ApplyAdditiveEnergyRegenBonus(sim *core.Simulation, increment float64) {
	oldBonus := rogue.AdditiveEnergyRegenBonus
	newBonus := oldBonus + increment
	rogue.AdditiveEnergyRegenBonus = newBonus
	rogue.MultiplyEnergyRegenSpeed(sim, (1.0+newBonus)/(1.0+oldBonus))
}

func (rogue *Rogue) Reset(sim *core.Simulation) {
	for _, mcd := range rogue.GetMajorCooldowns() {
		mcd.Disable()
	}

	rogue.MultiplyEnergyRegenSpeed(sim, 1.0+rogue.AdditiveEnergyRegenBonus)
}

func (rogue *Rogue) MeleeCritMultiplier(applyLethality bool) float64 {
	secondaryModifier := 0.0
	if applyLethality {
		secondaryModifier += 0.1 * float64(rogue.Talents.Lethality)
	}
	return rogue.Character.MeleeCritMultiplier(1.0, secondaryModifier)
}
func (rogue *Rogue) SpellCritMultiplier() float64 {
	return rogue.Character.SpellCritMultiplier(1, 0)
}

func NewRogue(character *core.Character, options *proto.RogueOptions, talents string) *Rogue {
	rogue := &Rogue{
		Character: *character,
		Talents:   &proto.RogueTalents{},
		Options:   options,
	}
	core.FillTalentsProto(rogue.Talents.ProtoReflect(), talents, TalentTreeSizes)

	// Passive rogue threat reduction: https://wotlk.wowhead.com/spell=21184/rogue-passive-dnd
	rogue.PseudoStats.ThreatMultiplier *= 0.71
	rogue.PseudoStats.CanParry = true

	maxEnergy := 100.0
	if rogue.HasSetBonus(Arena, 4) {
		maxEnergy += 10
	}
	if rogue.Spec == proto.Spec_SpecAssassinationRogue &&
		rogue.GetMHWeapon() != nil &&
		rogue.GetMHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger {
		maxEnergy += 20
	}
	rogue.EnableEnergyBar(maxEnergy)

	rogue.EnableAutoAttacks(rogue, core.AutoAttackOptions{
		MainHand:       rogue.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		OffHand:        rogue.WeaponFromOffHand(0),  // Set crit multiplier later when we have targets.
		Ranged:         rogue.WeaponFromRanged(0),
		AutoSwingMelee: true,
	})
	rogue.applyPoisons()

	rogue.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.AttackPower, 2)
	rogue.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)
	// Make an assumption we're wearing leather for Leather Armor Spec
	rogue.MultiplyStat(stats.Agility, 1.05)

	return rogue
}

// Apply the effects of the Cut to the Chase talent
// TODO: Put a fresh instance of SnD rather than use the original as per client
// TODO (TheBackstabi, 3/16/2024) - Assassination only talent, to be moved?
func (rogue *Rogue) ApplyCutToTheChase(sim *core.Simulation) {
	if rogue.Talents.CutToTheChase > 0 && rogue.SliceAndDiceAura.IsActive() {
		procChance := []float64{0.0, 0.33, 0.67, 1.0}[rogue.Talents.CutToTheChase]
		if procChance == 1 || sim.Proc(procChance, "Cut to the Chase") {
			rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
			rogue.SliceAndDiceAura.Activate(sim)
		}
	}
}

// Deactivate Stealth if it is active. This must be added to all abilities that cause Stealth to fade.
func (rogue *Rogue) BreakStealth(sim *core.Simulation) {
	if rogue.StealthAura.IsActive() {
		rogue.StealthAura.Deactivate(sim)
		rogue.AutoAttacks.EnableAutoSwing(sim)
	}
}

// Does the rogue have a dagger equipped in the specified hand (main or offhand)?
func (rogue *Rogue) HasDagger(hand core.Hand) bool {
	if hand == core.MainHand {
		return rogue.MainHand().WeaponType == proto.WeaponType_WeaponTypeDagger
	}
	return rogue.OffHand().WeaponType == proto.WeaponType_WeaponTypeDagger
}

// Does the rogue have a thrown weapon equipped in the ranged slot?
func (rogue *Rogue) HasThrown() bool {
	return rogue.Ranged().RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeThrown
}

// Check if the rogue is considered in "stealth" for the purpose of casting abilities
func (rogue *Rogue) IsStealthed() bool {
	if rogue.StealthAura.IsActive() {
		return true
	}
	if rogue.Talents.ShadowDance && rogue.ShadowDanceAura.IsActive() {
		return true
	}
	return false
}

// Agent is a generic way to access underlying rogue on any of the agents.
type RogueAgent interface {
	GetRogue() *Rogue
}
