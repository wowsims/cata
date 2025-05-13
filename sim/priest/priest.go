package priest

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Priest struct {
	core.Character
	SelfBuffs
	Talents *proto.PriestTalents

	SurgeOfLight bool

	Latency float64

	ShadowfiendAura *core.Aura
	ShadowfiendPet  *Shadowfiend
	MindbenderPet   *MindBender
	MindbenderAura  *core.Aura

	ShadowOrbsAura      *core.Aura
	EmpoweredShadowAura *core.Aura

	// cached cast stuff
	// TODO: aoe multi-target situations will need multiple spells ticking for each target.
	HolyEvangelismProcAura *core.Aura
	DarkEvangelismProcAura *core.Aura

	SurgeOfLightProcAura *core.Aura

	// might want to move these spell / talents into spec specific initialization
	BindingHeal       *core.Spell
	CircleOfHealing   *core.Spell
	FlashHeal         *core.Spell
	GreaterHeal       *core.Spell
	Penance           *core.Spell
	PenanceHeal       *core.Spell
	PowerWordShield   *core.Spell
	PrayerOfHealing   *core.Spell
	PrayerOfMending   *core.Spell
	Renew             *core.Spell
	EmpoweredRenew    *core.Spell
	InnerFocus        *core.Spell
	HolyFire          *core.Spell
	Smite             *core.Spell
	ShadowWordPain    *core.Spell
	Shadowfiend       *core.Spell
	VampiricTouch     *core.Spell
	MindBender        *core.Spell
	ShadowyApparition *core.Spell

	WeakenedSouls core.AuraArray

	ProcPrayerOfMending core.ApplySpellResults
}

type SelfBuffs struct {
	UseShadowfiend bool
	UseInnerFire   bool

	PowerInfusionTarget *proto.UnitReference
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

func (priest *Priest) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {
	if priest.SelfBuffs.UseInnerFire {
		priest.MultiplyStat(stats.SpellPower, 1.1)
		priest.ApplyEquipScaling(stats.Armor, 1.1)
		core.MakePermanent(priest.RegisterAura(core.Aura{
			Label:    "Inner Fire",
			ActionID: core.ActionID{SpellID: 588},
		}))
	}

	priest.MultiplyStat(stats.Intellect, 1.05)
	priest.registerShadowWordPainSpell()
	priest.registerShadowfiendSpell()
	priest.registerVampiricTouchSpell()

	// priest.registerDispersionSpell()

	priest.registerPowerInfusionSpell()
	priest.newMindSearSpell()

	priest.ApplyGlyphs()
}

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

func (priest *Priest) ApplyTalents() {
	priest.registerMindbenderSpell()
}

func (priest *Priest) Reset(_ *core.Simulation) {
}

func New(char *core.Character, selfBuffs SelfBuffs, talents string) *Priest {
	priest := &Priest{
		Character: *char,
		SelfBuffs: selfBuffs,
		Talents:   &proto.PriestTalents{},
	}

	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents)
	priest.EnableManaBar()
	priest.ShadowfiendPet = priest.NewShadowfiend()

	if priest.Talents.Mindbender {
		priest.MindbenderPet = priest.NewMindBender()
	}

	return priest
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}

func (priest *Priest) HasMajorGlyph(glyph proto.PriestMajorGlyph) bool {
	return priest.HasGlyph(int32(glyph))
}
func (priest *Priest) HasMinorGlyph(glyph proto.PriestMinorGlyph) bool {
	return priest.HasGlyph(int32(glyph))
}

const (
	PriestSpellFlagNone  int64 = 0
	PriestSpellArchangel int64 = 1 << iota
	PriestSpellDarkArchangel
	PriestSpellBindingHeal
	PriestSpellCascade
	PriestSpellCircleOfHealing
	PriestSpellDevouringPlague
	PriestSpellDevouringPlagueDoT
	PriestSpellDesperatePrayer
	PriestSpellDispersion
	PriestSpellDivineAegis
	PriestSpellDivineHymn
	PriestSpellDivineStar
	PriestSpellEmpoweredRenew
	PriestSpellFade
	PriestSpellFlashHeal
	PriestSpellGreaterHeal
	PriestSpellGuardianSpirit
	PriestSpellHalo
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
	PriestSpellMindBender
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
	PriestSpellShadowyRecall
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
