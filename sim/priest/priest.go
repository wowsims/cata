package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

var TalentTreeSizes = [3]int{21, 21, 21}

type PriestSpell struct {
	*core.Spell
	ClassSpell PriestSpellFlag
}

type Priest struct {
	core.Character
	SelfBuffs
	Talents *proto.PriestTalents

	SurgeOfLight bool

	Latency float64

	ShadowfiendAura *core.Aura
	ShadowfiendPet  *Shadowfiend

	// Aura Mods
	DamageDonePercentMods    []*PriestAuraMod[float64]
	DamageDonePercentAddMods []*PriestAuraMod[float64]
	PowerCostPercentMods     []*PriestAuraMod[float64]
	CastTimePercentMods      []*PriestAuraMod[float64]
	CooldownMods             []*PriestAuraMod[time.Duration]
	ProcChanceMods           []*PriestAuraMod[float64]

	// cached cast stuff
	// TODO: aoe multi-target situations will need multiple spells ticking for each target.
	InnerFocusAura         *core.Aura
	HolyEvangelismProcAura *core.Aura
	DarkEvangelismProcAura *core.Aura
	ShadowyInsightAura     *core.Aura
	DispersionAura         *core.Aura
	MindMeltProcAura       *core.Aura

	SurgeOfLightProcAura *core.Aura

	// might want to move these spell / talents into spec specific initialization
	Archangel         *PriestSpell
	DarkArchangel     *PriestSpell
	BindingHeal       *PriestSpell
	CircleOfHealing   *PriestSpell
	DevouringPlague   *PriestSpell
	FlashHeal         *PriestSpell
	GreaterHeal       *PriestSpell
	HolyFire          *PriestSpell
	InnerFocus        *PriestSpell
	ShadowWordPain    *PriestSpell
	MindBlast         *PriestSpell
	MindFlay          []*PriestSpell
	MindFlayAPL       *PriestSpell
	MindSear          []*PriestSpell
	MindSearAPL       *PriestSpell
	Penance           *PriestSpell
	PenanceHeal       *PriestSpell
	PowerWordShield   *PriestSpell
	PrayerOfHealing   *PriestSpell
	PrayerOfMending   *PriestSpell
	Renew             *PriestSpell
	EmpoweredRenew    *PriestSpell
	ShadowWordDeath   *PriestSpell
	Shadowfiend       *PriestSpell
	Smite             *PriestSpell
	VampiricTouch     *PriestSpell
	Dispersion        *PriestSpell
	MindSpike         *PriestSpell
	ShadowyApparition *PriestSpell

	WeakenedSouls core.AuraArray

	ProcPrayerOfMending core.ApplySpellResults

	ScalingBaseDamage    float64
	ShadowCritMultiplier float64

	// set bonus cache
	// The mana cost of your Mind Blast is reduced by 10%.
	T7TwoSetBonus bool
	// Your Shadow Word: Death has an additional 10% chance to critically strike.
	T7FourSetBonus bool
	// Increases the damage done by your Devouring Plague by 15%.
	T8TwoSetBonus bool
	// Your Mind Blast also grants you 240 haste for 4 sec.
	T8FourSetBonus bool
	// Increases the duration of your Vampiric Touch spell by 6 sec.
	T9TwoSetBonus bool
	// Increases the critical strike chance of your Mind Flay spell by 5%.
	T9FourSetBonus bool
	// The critical strike chance of your Shadow Word: Pain, Devouring Plague, and Vampiric Touch spells is increased by 5%
	T10TwoSetBonus bool
	// Reduces the channel duration by 0.51 sec and period by 0.17 sec on your Mind Flay spell
	T10FourSetBonus bool
}

type SelfBuffs struct {
	UseShadowfiend bool
	UseInnerFire   bool

	PowerInfusionTarget *proto.UnitReference
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

// func (priest *Priest) HasMajorGlyph(glyph proto.PriestMajorGlyph) bool {
// 	return priest.HasGlyph(int32(glyph))
// }
// func (priest *Priest) HasMinorGlyph(glyph proto.PriestMinorGlyph) bool {
// 	return priest.HasGlyph(int32(glyph))
// }

// func (priest *Priest) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
// 	raidBuffs.ShadowProtection = true
// 	raidBuffs.DivineSpirit = true

// 	raidBuffs.PowerWordFortitude = max(raidBuffs.PowerWordFortitude, core.MakeTristateValue(
// 		true,
// 		priest.Talents.ImprovedPowerWordFortitude == 2))
// }

func (priest *Priest) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (priest *Priest) RegisterSpell(flag PriestSpellFlag, config core.SpellConfig) *PriestSpell {

	spell := &PriestSpell{ClassSpell: flag}
	apply := config.ApplyEffects
	snapShot := config.Dot.OnSnapshot

	if config.ManaCost.Multiplier == 0 {
		config.ManaCost.Multiplier = 1
	}

	if config.CritMultiplier == 0 {
		config.CritMultiplier = 1
	}

	if config.ThreatMultiplier == 0 {
		config.ThreatMultiplier = 1
	}

	config.ManaCost.Multiplier *= priest.GetClassSpellModPowerPercent(flag, config.SpellSchool)
	config.CritMultiplier *= priest.SpellCritMultiplier(1, priest.ShadowCritMultiplier)
	config.Cast.CD.Duration += priest.GetClassSpellModCooldown(flag, config.SpellSchool)
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.DamageMultiplier = priest.GetClassSpellDamageDonePercent(flag, config.SpellSchool)
		spell.DamageMultiplierAdditive = priest.GetClassSpellDamageDoneAddPercent(flag, config.SpellSchool)
		if apply != nil {
			apply(sim, target, spell)
		}
	}

	if snapShot != nil {
		config.Dot.OnSnapshot = func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.Spell.DamageMultiplier = priest.GetClassSpellDamageDonePercent(PriestSpellMindFlay, core.SpellSchoolShadow)
			dot.Spell.DamageMultiplierAdditive = priest.GetClassSpellDamageDoneAddPercent(PriestSpellMindFlay, core.SpellSchoolShadow)
			snapShot(sim, target, dot, isRollover)
		}
	}

	spell.Spell = priest.Unit.RegisterSpell(config)
	return spell
}

func (priest *Priest) Initialize() {

	// base scaling value for a level 85 priest
	priest.ScalingBaseDamage = 945.188842773437500

	// priest.registerSetBonuses()
	priest.registerDevouringPlagueSpell()
	priest.registerShadowWordPainSpell()

	priest.registerMindBlastSpell()
	priest.registerShadowWordDeathSpell()
	priest.registerShadowfiendSpell()
	priest.registerVampiricTouchSpell()
	priest.registerDispersionSpell()

	// priest.registerPowerInfusionCD()

	priest.MindFlayAPL = priest.newMindFlaySpell(0)
	priest.MindSearAPL = priest.newMindSearSpell(0)

	priest.MindFlay = []*PriestSpell{
		nil, // So we can use # of ticks as the index
		priest.newMindFlaySpell(1),
		priest.newMindFlaySpell(2),
		priest.newMindFlaySpell(3),
	}

	priest.MindSear = []*PriestSpell{
		nil, // So we can use # of ticks as the index
		priest.newMindSearSpell(1),
		priest.newMindSearSpell(2),
		priest.newMindSearSpell(3),
		priest.newMindSearSpell(4),
		priest.newMindSearSpell(5),
	}
}

// func (priest *Priest) RegisterHealingSpells() {
// 	priest.registerPenanceHealSpell()
// 	priest.registerBindingHealSpell()
// 	priest.registerCircleOfHealingSpell()
// 	priest.registerFlashHealSpell()
// 	priest.registerGreaterHealSpell()
// 	priest.registerPowerWordShieldSpell()
// 	priest.registerPrayerOfHealingSpell()
// 	priest.registerPrayerOfMendingSpell()
// 	priest.registerRenewSpell()
// }

func (priest *Priest) AddHolyEvanglismStack(sim *core.Simulation) {
	if priest.HolyEvangelismProcAura != nil {
		priest.HolyEvangelismProcAura.Activate(sim)
		priest.HolyEvangelismProcAura.AddStack(sim)
	}
}

func (priest *Priest) AddDarkEvangelismStack(sim *core.Simulation) {
	if priest.DarkEvangelismProcAura != nil {
		priest.DarkEvangelismProcAura.Activate(sim)
		priest.DarkEvangelismProcAura.AddStack(sim)
	}
}

func (priest *Priest) Reset(_ *core.Simulation) {
}

func New(char *core.Character, selfBuffs SelfBuffs, talents string) *Priest {
	priest := &Priest{
		Character: *char,
		SelfBuffs: selfBuffs,
		Talents:   &proto.PriestTalents{},
	}

	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents, TalentTreeSizes)
	priest.EnableManaBar()
	priest.ShadowfiendPet = priest.NewShadowfiend()

	return priest
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}

// Priest Specific Aura MOD Handling
// Other flags
type PriestSpellFlag uint64

// Returns whether there is any overlap between the given masks.
func (se PriestSpellFlag) Matches(other PriestSpellFlag) bool {
	return (se & other) != 0
}

const (
	PriestSpellFlagNone  PriestSpellFlag = 0
	PriestSpellArchangel PriestSpellFlag = 1 << iota
	PriestSpellDarkArchangel
	PriestSpellBindingHeal
	PriestSpellCircleOfHealing
	PriestSpellDevouringPlague
	PriestSpellDesperatePrayer
	PriestSpellDispersion
	PriestSpellDivineHymn
	PriestSpellFade
	PriestSpellFlashHeal
	PriestSpellGreaterHeal
	PriestSpellGuardianSpirit
	PriestSpellHolyFire
	PriestSpellHolyNova
	PriestSpellHolyWordChastise
	PriestSpellHolyWordSanctuary
	PriestSpellHolyWordSerenity
	PriestSpellHymnOfHope
	PriestSpellImprovedDevouringPlague
	PriestSpellInnerFire
	PriestSpellInnerFocus
	PriestSpellInnerWill
	PriestSpellManaBurn
	PriestSpellMindBlast
	PriestSpellMindFlay
	PriestSpellMindSear
	PriestSpellMindSpike
	PriestSpellMindTrauma
	PriestSpellPainSuppresion
	PriestSpellPenance
	PriestSpellPowerInfusion
	PriestSpellPowerWordBarrier
	PriestSpellPowerWordShield
	PriestSpellPrayerOfHealing
	PriestSpellPrayerOfMending
	PriestSpellPsychicScream
	PriestSpellRenew
	PriestSpellShadowOrbPassive
	PriestSpellShadowWordDeath
	PriestSpellShadowWordPain
	PriestSpellShadowFiend
	PriestSpellShadowyApparation
	PriestSpellSmite
	PriestSpellVampiricEmbrace
	PriestSpellVampiricTouch

	PriestSpellDoT     = PriestSpellDevouringPlague | PriestSpellHolyFire | PriestSpellMindFlay | PriestSpellShadowWordPain | PriestSpellVampiricTouch
	PriestSpellInstant = PriestSpellCircleOfHealing |
		PriestSpellDesperatePrayer |
		PriestSpellDevouringPlague |
		PriestSpellFade |
		PriestSpellGuardianSpirit |
		PriestSpellHolyNova |
		PriestSpellHolyWordChastise |
		PriestSpellHolyWordSanctuary |
		PriestSpellHolyWordSerenity |
		PriestSpellInnerFire |
		PriestSpellPainSuppresion |
		PriestSpellPowerInfusion |
		PriestSpellPowerWordBarrier |
		PriestSpellPowerWordShield |
		PriestSpellRenew |
		PriestSpellShadowWordDeath |
		PriestSpellShadowWordPain |
		PriestSpellVampiricEmbrace
)

type PriestAuraMod[T any] struct {
	ClassSpell PriestSpellFlag
	School     core.SpellSchool
	BaseValue  T

	// dynamic evaluation will be added to the given BaseValue
	DynamicValue func(*Priest) T

	// Stacks * (BaseValue + DynamicValue)
	Stacks  int32
	SpellID int32
}

// Adds a PriestAuraMod to the list of modifiers or replaces an
// existing mod if it has the same spellID.
func AddOrReplaceMod[T any](modList *[]*PriestAuraMod[T], mod *PriestAuraMod[T]) {
	if mod.SpellID == 0 {
		panic("mod.SpellID should never be 0")
	}

	if mod.Stacks == 0 {
		mod.Stacks = 1
	}

	for key, val := range *modList {
		if val.SpellID == mod.SpellID && val.ClassSpell == mod.ClassSpell {
			(*modList)[key] = mod
			return
		}
	}

	*modList = append(*modList, mod)
}

// Removes all mods for a specific spellID
func RemoveMod[T any](modList *[]*PriestAuraMod[T], spellID int32) {
	removeIdx := []int{}
	for key, val := range *modList {
		if val.SpellID == spellID {
			removeIdx = append(removeIdx, key)
		}
	}

	// order of operation is not significant for mod lists
	// move last index to remove idex and shorten slice
	for i := 0; i < len(removeIdx); i++ {
		idx := removeIdx[len(removeIdx)-1-i]
		(*modList)[idx] = (*modList)[len(*modList)-1-i]
	}

	*modList = (*modList)[:len(*modList)-len(removeIdx)]
}

func (priest *Priest) GetClassSpellDamageDonePercent(spell PriestSpellFlag, school core.SpellSchool) float64 {
	return applyMod(priest, &priest.DamageDonePercentMods, 1, spell, school, multiplyOp)
}

func (priest *Priest) GetClassSpellDamageDoneAddPercent(spell PriestSpellFlag, school core.SpellSchool) float64 {
	return applyMod(priest, &priest.DamageDonePercentAddMods, 1, spell, school, multiplyOp)
}

func (priest *Priest) GetClassSpellModPowerPercent(spell PriestSpellFlag, school core.SpellSchool) float64 {
	return 1 - applyMod(priest, &priest.PowerCostPercentMods, 1, spell, school, addOp)
}

func (priest *Priest) GetClassSpellModCooldown(spell PriestSpellFlag, school core.SpellSchool) time.Duration {
	return applyMod(priest, &priest.CooldownMods, time.Duration(0), spell, school, addOpTime)
}

func (priest *Priest) GetClassSpellProcChance(base float64, spell PriestSpellFlag, school core.SpellSchool) float64 {
	return max(0, base+applyMod(priest, &priest.ProcChanceMods, 0, spell, school, addOp))
}

func addOp(base float64, value float64, stacks int32) float64 {
	return base + (value * float64(stacks))
}

func addOpTime(base time.Duration, value time.Duration, stacks int32) time.Duration {
	return base - time.Duration(int32(value.Milliseconds())*stacks)*time.Millisecond
}

func multiplyOp(base float64, value float64, stacks int32) float64 {
	return base * (1 + value*float64(stacks))
}

func applyMod[T float64 | time.Duration](priest *Priest, modList *[]*PriestAuraMod[T], base T, spell PriestSpellFlag, school core.SpellSchool, op func(T, T, int32) T) T {
	for _, mod := range *modList {
		if mod.ClassSpell.Matches(spell) &&
			(mod.School == core.SpellSchoolNone || mod.School.Matches(school)) {
			baseValue := mod.BaseValue
			if mod.DynamicValue != nil {
				baseValue += mod.DynamicValue(priest)
			}

			base = op(base, baseValue, mod.Stacks)
		}
	}

	return base
}

func (ps *PriestSpell) IsEqual(other *core.Spell) bool {
	if ps == nil || other == nil {
		return false
	}

	return ps.Spell == other
}
