package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const (
	SpellFlagBuilder     = core.SpellFlagAgentReserved2
	SpellFlagFinisher    = core.SpellFlagAgentReserved3
	SpellFlagColdBlooded = core.SpellFlagAgentReserved4
)

const RogueBleedTag = "RogueBleed"

type Rogue struct {
	core.Character

	ClassSpellScaling float64

	Talents              *proto.RogueTalents
	Options              *proto.RogueOptions
	AssassinationOptions *proto.AssassinationRogue_Options
	CombatOptions        *proto.CombatRogue_Options
	SubtletyOptions      *proto.SubtletyRogue_Options

	MasteryBaseValue  float64
	MasteryMultiplier float64

	SliceAndDiceBonusFlat    float64 // The flat bonus Attack Speed bonus before Mastery is applied
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

	deadlyPoisonPPHM  [core.NumItemSlots]*core.DynamicProcManager
	instantPoisonPPMM [core.NumItemSlots]*core.DynamicProcManager
	woundPoisonPPMM   [core.NumItemSlots]*core.DynamicProcManager

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

	MasterPoisonerDebuffAuras core.AuraArray
	SavageCombatDebuffAuras   core.AuraArray
	WoundPoisonDebuffAuras    core.AuraArray

	T12ToTLastBuff int

	ruthlessnessMetrics      *core.ResourceMetrics
	relentlessStrikesMetrics *core.ResourceMetrics
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
	// numPoints := rogue.ComboPoints()
	// rogue.SpendComboPoints(sim, spell.ComboPointMetrics())

	// TODO: Fix this to work with the new talent system.
	// if rogue.Talents.Ruthlessness > 0 && (spell.ClassSpellMask&RogueSpellDamagingFinisher != 0) {
	// 	procChance := 0.2 * float64(rogue.Talents.Ruthlessness)
	// 	if sim.Proc(procChance, "Ruthlessness") {
	// 		rogue.AddComboPoints(sim, 1, rogue.ruthlessnessMetrics)
	// 	}
	// }
	// if rogue.Talents.RelentlessStrikes > 0 {
	// 	procChance := []float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.RelentlessStrikes] * float64(numPoints)
	// 	if sim.Proc(procChance, "Relentless Strikes") {
	// 		rogue.AddEnergy(sim, 25, rogue.relentlessStrikesMetrics)
	// 	}
	// }
	// if rogue.Talents.RestlessBlades > 0 && (spell.ClassSpellMask&RogueSpellDamagingFinisher != 0) {
	// 	cdReduction := time.Duration(rogue.Talents.RestlessBlades) * time.Second * time.Duration(numPoints)

	// 	if rogue.KillingSpree != nil {
	// 		ksNewTime := rogue.KillingSpree.CD.Timer.ReadyAt() - cdReduction
	// 		rogue.KillingSpree.CD.Timer.Set(ksNewTime)
	// 	}
	// 	if rogue.AdrenalineRush != nil {
	// 		arNewTime := rogue.AdrenalineRush.CD.Timer.ReadyAt() - cdReduction
	// 		rogue.AdrenalineRush.CD.Timer.Set(arNewTime)
	// 	}
	// }
	// if rogue.Talents.SerratedBlades > 0 && spell == rogue.Eviscerate {
	// 	chancePerPoint := 0.1 * float64(rogue.Talents.SerratedBlades)
	// 	procChance := float64(numPoints) * chancePerPoint
	// 	if sim.Proc(procChance, "Serrated Blades") {
	// 		rupAura := rogue.Rupture.Dot(spell.Unit.CurrentTarget)
	// 		if rupAura.IsActive() {
	// 			rupAura.Activate(sim)
	// 		}
	// 	}
	// }
}

func (rogue *Rogue) HasMajorGlyph(glyph proto.RogueMajorGlyph) bool {
	return rogue.HasGlyph(int32(glyph))
}

func (rogue *Rogue) HasMinorGlyph(glyph proto.RogueMinorGlyph) bool {
	return rogue.HasGlyph(int32(glyph))
}

func (rogue *Rogue) Initialize() {
	// Update auto crit multipliers now that we have the targets.
	rogue.AutoAttacks.MHConfig().CritMultiplier = rogue.CritMultiplier(false)
	rogue.AutoAttacks.OHConfig().CritMultiplier = rogue.CritMultiplier(false)
	rogue.AutoAttacks.RangedConfig().CritMultiplier = rogue.CritMultiplier(false)

	// rogue.registerStealthAura()
	// rogue.registerVanishSpell()
	rogue.registerFeintSpell()
	// rogue.registerAmbushSpell()
	// rogue.registerGarrote()
	// rogue.registerSinisterStrikeSpell()
	// rogue.registerBackstabSpell()
	// rogue.registerRupture()
	// rogue.registerSliceAndDice()
	// rogue.registerEviscerate()
	// rogue.registerEnvenom()
	// rogue.registerExposeArmorSpell()
	// rogue.registerRecuperate()
	// rogue.registerFanOfKnives()
	// rogue.registerTricksOfTheTradeSpell()
	// rogue.registerDeadlyPoisonSpell()
	// rogue.registerInstantPoisonSpell()
	// rogue.registerWoundPoisonSpell()
	// rogue.registerPoisonAuras()
	// rogue.registerShivSpell()
	rogue.registerThistleTeaCD()
	// rogue.registerGougeSpell()

	rogue.T12ToTLastBuff = 3

	// re-configure poisons when performing an item swap
	rogue.RegisterItemSwapCallback(core.MeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		if !rogue.Options.ApplyPoisonsManually {
			if rogue.MainHand() == nil || rogue.OffHand() == nil {
				return
			}
			mhWeaponSpeed := rogue.MainHand().SwingSpeed
			ohWeaponSpeed := rogue.OffHand().SwingSpeed
			if mhWeaponSpeed <= ohWeaponSpeed {
				rogue.Options.MhImbue = proto.RogueOptions_DeadlyPoison
				rogue.Options.OhImbue = proto.RogueOptions_InstantPoison
				rogue.lastDeadlyPoisonProcMask = core.ProcMaskMeleeMH
			} else {
				rogue.Options.MhImbue = proto.RogueOptions_InstantPoison
				rogue.Options.OhImbue = proto.RogueOptions_DeadlyPoison
				rogue.lastDeadlyPoisonProcMask = core.ProcMaskMeleeOH
			}
			// rogue.UpdateInstantPoisonPPM(0)
		}
	})
}

func (rogue *Rogue) ApplyTalents() {}

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

	rogue.T12ToTLastBuff = 3
}

func (rogue *Rogue) CritMultiplier(applyLethality bool) float64 {
	secondaryModifier := 0.0
	// TODO: Fix this to work with the new talent system.
	// if applyLethality {
	// 	secondaryModifier += 0.1 * float64(rogue.Talents.Lethality)
	// }
	return rogue.GetCharacter().CritMultiplier(1.0, secondaryModifier)
}

func NewRogue(character *core.Character, options *proto.RogueOptions, talents string) *Rogue {
	rogue := &Rogue{
		Character:         *character,
		Talents:           &proto.RogueTalents{},
		Options:           options,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassRogue),
	}

	core.FillTalentsProto(rogue.Talents.ProtoReflect(), talents)

	// Passive rogue threat reduction: https://wotlk.wowhead.com/spell=21184/rogue-passive-dnd
	rogue.PseudoStats.ThreatMultiplier *= 0.71
	rogue.PseudoStats.CanParry = true

	maxEnergy := 100.0

	if rogue.Spec == proto.Spec_SpecAssassinationRogue &&
		rogue.GetMHWeapon() != nil &&
		rogue.GetMHWeapon().WeaponType == proto.WeaponType_WeaponTypeDagger {
		maxEnergy += 20
	}

	rogue.EnableEnergyBar(core.EnergyBarOptions{
		MaxComboPoints:      5,
		MaxEnergy:           maxEnergy,
		StartingComboPoints: options.StartingComboPoints,
		UnitClass:           proto.Class_ClassRogue,
	})

	rogue.EnableAutoAttacks(rogue, core.AutoAttackOptions{
		MainHand:       rogue.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		OffHand:        rogue.WeaponFromOffHand(0),  // Set crit multiplier later when we have targets.
		Ranged:         rogue.WeaponFromRanged(0),
		AutoSwingMelee: true,
	})
	// rogue.applyPoisons()

	rogue.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	rogue.AddStatDependency(stats.Agility, stats.AttackPower, 2)
	rogue.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])

	return rogue
}

// Apply the effects of the Cut to the Chase talent
// TODO: Put a fresh instance of SnD rather than use the original as per client
// TODO (TheBackstabi, 3/16/2024) - Assassination only talent, to be moved?
// func (rogue *Rogue) ApplyCutToTheChase(sim *core.Simulation) {
// 	if rogue.Talents.CutToTheChase > 0 && rogue.SliceAndDiceAura.IsActive() {
// 		procChance := []float64{0.0, 0.33, 0.67, 1.0}[rogue.Talents.CutToTheChase]
// 		if procChance == 1 || sim.Proc(procChance, "Cut to the Chase") {
// 			rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
// 			rogue.SliceAndDiceAura.Activate(sim)
// 		}
// 	}
// }

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
	weapon := rogue.Ranged()
	return weapon != nil && weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeThrown
}

// Check if the rogue is considered in "stealth" for the purpose of casting abilities
func (rogue *Rogue) IsStealthed() bool {
	if rogue.StealthAura.IsActive() {
		return true
	}
	// TODO: Fix this to work with the new talent system.
	// if rogue.Talents.ShadowDance && rogue.ShadowDanceAura.IsActive() {
	// 	return true
	// }
	return false
}

func (rogue *Rogue) GetMasteryBonus() float64 {
	return rogue.GetMasteryBonusFromRating(rogue.GetStat(stats.MasteryRating))
}
func (rogue *Rogue) GetMasteryBonusFromRating(masteryRating float64) float64 {
	return rogue.MasteryBaseValue + core.MasteryRatingToMasteryPoints(masteryRating)*rogue.MasteryMultiplier
}

// Agent is a generic way to access underlying rogue on any of the agents.
type RogueAgent interface {
	GetRogue() *Rogue
}

const (
	RogueSpellFlagNone int64 = 0
	RogueSpellAmbush   int64 = 1 << iota
	RogueSpellBackstab
	RogueSpellEnvenom
	RogueSpellEviscerate
	RogueSpellExposeArmor
	RogueSpellFanOfKnives
	RogueSpellFeint
	RogueSpellGarrote
	RogueSpellGouge
	RogueSpellRecuperate
	RogueSpellRupture
	RogueSpellShiv
	RogueSpellSinisterStrike
	RogueSpellSliceAndDice
	RogueSpellStealth
	RogueSpellTricksOfTheTrade
	RogueSpellTricksOfTheTradeThreat
	RogueSpellVanish
	RogueSpellHemorrhage
	RogueSpellPremeditation
	RogueSpellPreparation
	RogueSpellShadowDance
	RogueSpellShadowstep
	RogueSpellAdrenalineRush
	RogueSpellBladeFlurry
	RogueSpellKillingSpree
	RogueSpellKillingSpreeHit
	RogueSpellMainGauche
	RogueSpellRevealingStrike
	RogueSpellColdBlood
	RogueSpellMutilate
	RogueSpellVendetta
	RogueSpellVenomousWounds
	RogueSpellWoundPoison
	RogueSpellInstantPoison
	RogueSpellDeadlyPoison

	RogueSpellLast
	RogueSpellsAll = RogueSpellLast<<1 - 1

	RogueSpellPoisons          = RogueSpellVenomousWounds | RogueSpellWoundPoison | RogueSpellInstantPoison | RogueSpellDeadlyPoison
	RogueSpellDamagingFinisher = RogueSpellEnvenom | RogueSpellEviscerate | RogueSpellRupture
	RogueSpellWeightedBlades   = RogueSpellSinisterStrike | RogueSpellRevealingStrike
)
