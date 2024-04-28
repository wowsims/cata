package priest

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var TalentTreeSizes = [3]int{21, 21, 21}

type Priest struct {
	core.Character
	SelfBuffs
	Talents *proto.PriestTalents

	SurgeOfLight bool

	Latency float64

	ShadowfiendAura *core.Aura
	ShadowfiendPet  *Shadowfiend

	// cached cast stuff
	// TODO: aoe multi-target situations will need multiple spells ticking for each target.
	HolyEvangelismProcAura *core.Aura
	DarkEvangelismProcAura *core.Aura

	SurgeOfLightProcAura *core.Aura

	// might want to move these spell / talents into spec specific initialization
	Archangel         *core.Spell
	DarkArchangel     *core.Spell
	BindingHeal       *core.Spell
	CircleOfHealing   *core.Spell
	DevouringPlague   *core.Spell
	FlashHeal         *core.Spell
	GreaterHeal       *core.Spell
	HolyFire          *core.Spell
	InnerFocus        *core.Spell
	ShadowWordPain    *core.Spell
	MindBlast         *core.Spell
	MindFlay          []*core.Spell
	MindFlayAPL       *core.Spell
	MindSear          []*core.Spell
	MindSearAPL       *core.Spell
	Penance           *core.Spell
	PenanceHeal       *core.Spell
	PowerWordShield   *core.Spell
	PrayerOfHealing   *core.Spell
	PrayerOfMending   *core.Spell
	Renew             *core.Spell
	EmpoweredRenew    *core.Spell
	ShadowWordDeath   *core.Spell
	Shadowfiend       *core.Spell
	Smite             *core.Spell
	VampiricTouch     *core.Spell
	Dispersion        *core.Spell
	MindSpike         *core.Spell
	ShadowyApparition *core.Spell

	WeakenedSouls core.AuraArray

	ProcPrayerOfMending core.ApplySpellResults

	ClassSpellScaling    float64
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

func (priest *Priest) Initialize() {
	if priest.SelfBuffs.UseInnerFire {
		priest.AddStat(stats.SpellPower, 531)
		priest.ApplyEquipScaling(stats.Armor, 1.6)
		core.MakePermanent(priest.RegisterAura(core.Aura{
			Label:    "Inner Fire",
			ActionID: core.ActionID{SpellID: 588},
		}))
	}

	// priest.registerSetBonuses()
	priest.registerDevouringPlagueSpell()
	priest.registerShadowWordPainSpell()

	priest.registerMindBlastSpell()
	priest.registerShadowWordDeathSpell()
	priest.registerShadowfiendSpell()
	priest.registerVampiricTouchSpell()
	priest.registerDispersionSpell()
	priest.registerMindSpike()

	priest.registerPowerInfusionSpell()

	priest.MindFlayAPL = priest.newMindFlaySpell(0)
	priest.MindSearAPL = priest.newMindSearSpell(0)

	priest.MindFlay = []*core.Spell{
		nil, // So we can use # of ticks as the index
		priest.newMindFlaySpell(1),
		priest.newMindFlaySpell(2),
		priest.newMindFlaySpell(3),
	}

	priest.MindSear = []*core.Spell{
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
		Character:         *char,
		SelfBuffs:         selfBuffs,
		Talents:           &proto.PriestTalents{},
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassPriest),
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

func (hunter *Priest) HasPrimeGlyph(glyph proto.PriestPrimeGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}
func (hunter *Priest) HasMajorGlyph(glyph proto.PriestMajorGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}
func (hunter *Priest) HasMinorGlyph(glyph proto.PriestMinorGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}

const (
	PriestSpellFlagNone  int64 = 0
	PriestSpellArchangel int64 = 1 << iota
	PriestSpellDarkArchangel
	PriestSpellBindingHeal
	PriestSpellCircleOfHealing
	PriestSpellDevouringPlague
	PriestSpellDesperatePrayer
	PriestSpellDispersion
	PriestSpellDivineAegis
	PriestSpellDivineHymn
	PriestSpellEmpoweredRenew
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

	PriestSpellLast
	PriestSpellsAll    = PriestSpellLast<<1 - 1
	PriestSpellDoT     = PriestSpellDevouringPlague | PriestSpellHolyFire | PriestSpellMindFlay | PriestSpellShadowWordPain | PriestSpellVampiricTouch | PriestSpellImprovedDevouringPlague
	PriestSpellInstant = PriestSpellCircleOfHealing |
		PriestSpellDesperatePrayer |
		PriestSpellDevouringPlague |
		PriestSpellImprovedDevouringPlague |
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
	PriestShadowSpells = PriestSpellImprovedDevouringPlague |
		PriestSpellDevouringPlague |
		PriestSpellShadowWordDeath |
		PriestSpellShadowWordPain |
		PriestSpellMindFlay |
		PriestSpellMindBlast |
		PriestSpellMindSear |
		PriestSpellMindSpike |
		PriestSpellVampiricTouch
)

func (priest *Priest) calcBaseDamage(sim *core.Simulation, coefficient float64, variance float64) float64 {
	baseDamage := priest.ClassSpellScaling * coefficient
	if variance > 0 {
		delta := priest.ClassSpellScaling * variance * 0.5
		baseDamage += sim.Roll(-delta, delta)
	}

	return baseDamage
}
