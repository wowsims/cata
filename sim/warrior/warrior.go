package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var TalentTreeSizes = [3]int{20, 21, 20}

type WarriorInputs struct {
	StanceSnapshot bool
}

const (
	SpellFlagBleed = core.SpellFlagAgentReserved1
	ArmsTree       = 0
	FuryTree       = 1
	ProtTree       = 2
)

const SpellMaskNone int64 = 0
const (
	SpellMaskSpecialAttack int64 = 1 << iota

	// Abilities that don't cost rage and aren't attacks
	SpellMaskBattleShout
	SpellMaskBerserkerRage
	SpellMaskCommandingShout
	SpellMaskRecklessness
	SpellMaskShieldWall
	SpellMaskLastStand
	SpellMaskDeadlyCalm

	// Abilities that cost rage but aren't attacks
	SpellMaskDemoShout
	SpellMaskInnerRage
	SpellMaskShieldBlock
	SpellMaskDeathWish
	SpellMaskSweepingStrikes

	// Special attacks
	SpellMaskCleave
	SpellMaskColossusSmash
	SpellMaskExecute
	SpellMaskHeroicStrike
	SpellMaskHeroicThrow
	SpellMaskOverpower
	SpellMaskRend
	SpellMaskRevenge
	SpellMaskShatteringThrow
	SpellMaskSlam
	SpellMaskSunderArmor
	SpellMaskThunderClap
	SpellMaskWhirlwind
	SpellMaskShieldSlam
	SpellMaskConcussionBlow
	SpellMaskDevastate
	SpellMaskShockwave
	SpellMaskVictoryRush
	SpellMaskBloodthirst
	SpellMaskRagingBlow
	SpellMaskMortalStrike
	SpellMaskBladestorm
	SpellMaskHeroicLeap
)

const EnableOverpowerTag = "EnableOverpower"
const EnrageTag = "EnrageEffect"

type Warrior struct {
	core.Character

	ClassSpellScaling float64

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance                 Stance
	EnrageEffectMultiplier float64
	CriticalBlockChance    float64 // Can be gained as non-prot via certain talents and spells
	BlockDamageReduction   float64

	BattleShout     *core.Spell
	CommandingShout *core.Spell
	BattleStance    *core.Spell
	DefensiveStance *core.Spell
	BerserkerStance *core.Spell

	BerserkerRage     *core.Spell
	ColossusSmash     *core.Spell
	DemoralizingShout *core.Spell
	Execute           *core.Spell
	Overpower         *core.Spell
	Rend              *core.Spell
	Revenge           *core.Spell
	ShieldBlock       *core.Spell
	Slam              *core.Spell
	SunderArmor       *core.Spell
	ThunderClap       *core.Spell
	Whirlwind         *core.Spell
	DeepWounds        *core.Spell

	recklessnessDeadlyCalmCD *core.Timer
	hsCleaveCD               *core.Timer
	HeroicStrike             *core.Spell
	Cleave                   *core.Spell

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	BerserkerRageAura *core.Aura
	BloodsurgeAura    *core.Aura
	SuddenDeathAura   *core.Aura
	ShieldBlockAura   *core.Aura
	ThunderstruckAura *core.Aura
	InnerRageAura     *core.Aura

	DemoralizingShoutAuras core.AuraArray
	SunderArmorAuras       core.AuraArray
	ThunderClapAuras       core.AuraArray
	ColossusSmashAuras     core.AuraArray
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	if warrior.Talents.Rampage {
		raidBuffs.Rampage = true
	}
}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.registerStances()
	warrior.EnrageEffectMultiplier = 1.0
	warrior.hsCleaveCD = warrior.NewTimer()
	warrior.BlockDamageReduction = 0.3

	warrior.RegisterBerserkerRageSpell()
	warrior.RegisterColossusSmash()
	warrior.RegisterDemoralizingShoutSpell()
	warrior.RegisterExecuteSpell()
	warrior.RegisterHeroicStrikeSpell()
	warrior.RegisterCleaveSpell()
	warrior.RegisterHeroicLeap()
	warrior.RegisterHeroicThrow()
	warrior.RegisterInnerRage()
	warrior.RegisterOverpowerSpell()
	warrior.RegisterRecklessnessCD()
	warrior.RegisterRendSpell()
	warrior.RegisterRevengeSpell()
	warrior.RegisterShatteringThrowCD()
	warrior.RegisterShieldBlockCD()
	warrior.RegisterShieldWallCD()
	warrior.RegisterShouts()
	warrior.RegisterSlamSpell()
	warrior.RegisterSunderArmor()
	warrior.RegisterThunderClapSpell()
	warrior.RegisterWhirlwindSpell()
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
}

func NewWarrior(character *core.Character, talents string, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:         *character,
		Talents:           &proto.WarriorTalents{},
		WarriorInputs:     inputs,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassWarrior),
	}
	core.FillTalentsProto(warrior.Talents.ProtoReflect(), talents, TalentTreeSizes)

	warrior.PseudoStats.CanParry = true

	warrior.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)
	// Dodge no longer granted from agility
	//warrior.AddStatDependency(stats.Agility, stats.Dodge, core.DodgeRatingPerDodgeChance/84.746)
	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warrior.AddStatDependency(stats.Strength, stats.Parry, 0.25) // Change from block to pary in cata
	warrior.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// Base dodge unaffected by Diminishing Returns
	warrior.PseudoStats.BaseDodge += 0.03664
	warrior.PseudoStats.BaseParry += 0.05

	return warrior
}

func (warrior *Warrior) HasPrimeGlyph(glyph proto.WarriorPrimeGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) HasMajorGlyph(glyph proto.WarriorMajorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) HasMinorGlyph(glyph proto.WarriorMinorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) IntensifyRageCooldown(baseCd time.Duration) time.Duration {
	baseCd /= 100
	return []time.Duration{baseCd * 100, baseCd * 90, baseCd * 80}[warrior.Talents.IntensifyRage]
}

// Shared cooldown for Deadly Calm and Recklessness Activation
func (warrior *Warrior) RecklessnessDeadlyCalmLock() *core.Timer {
	return warrior.Character.GetOrInitTimer(&warrior.recklessnessDeadlyCalmCD)
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
