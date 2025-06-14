package warrior

import (
	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

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
	SpellMaskCharge

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
	SpellMaskWhirlwindOh
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

	SpellMaskShouts = SpellMaskCommandingShout | SpellMaskBattleShout
)

const EnableOverpowerTag = "EnableOverpower"
const EnrageTag = "EnrageEffect"

type Warrior struct {
	core.Character

	ClassSpellScaling float64

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	// Stance                 Stance
	EnrageEffectMultiplier float64
	CriticalBlockChance    []float64 // Can be gained as non-prot via certain talents and spells
	PrecisionKnown         bool

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
	Charge            *core.Spell
	ChargeAura        *core.Aura

	shoutsCD                 *core.Timer
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

	// Cached Gurthalak tentacles
	gurthalakTentacles []*cata.TentacleOfTheOldOnesPet
}

func (warrior *Warrior) GetTentacles() []*cata.TentacleOfTheOldOnesPet {
	return warrior.gurthalakTentacles
}

func (warrior *Warrior) NewTentacleOfTheOldOnesPet() *cata.TentacleOfTheOldOnesPet {
	pet := cata.NewTentacleOfTheOldOnesPet(&warrior.Character)
	warrior.AddPet(pet)
	return pet
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

// func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
// 	if warrior.Talents.Rampage {
// 		raidBuffs.Rampage = true
// 	}
// }

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	// warrior.registerStances()
	warrior.EnrageEffectMultiplier = 1.0
	warrior.hsCleaveCD = warrior.NewTimer()
	warrior.shoutsCD = warrior.NewTimer()

	// warrior.RegisterBerserkerRageSpell()
	// warrior.RegisterColossusSmash()
	// warrior.RegisterDemoralizingShoutSpell()
	// warrior.RegisterExecuteSpell()
	// warrior.RegisterHeroicStrikeSpell()
	// warrior.RegisterCleaveSpell()
	warrior.RegisterHeroicLeap()
	// warrior.RegisterHeroicThrow()
	warrior.RegisterInnerRage()
	// warrior.RegisterOverpowerSpell()
	warrior.RegisterRecklessnessCD()
	// warrior.RegisterRendSpell()
	// warrior.RegisterRevengeSpell()
	// warrior.RegisterShatteringThrowCD()
	// warrior.RegisterShieldBlockCD()
	warrior.RegisterShieldWallCD()
	// warrior.RegisterShouts()
	// warrior.RegisterSlamSpell()
	// warrior.RegisterSunderArmor()
	// warrior.RegisterThunderClapSpell()
	// warrior.RegisterWhirlwindSpell()
	// warrior.RegisterCharge()
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
	// warrior.Stance = StanceNone
}

func NewWarrior(character *core.Character, talents string, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:         *character,
		Talents:           &proto.WarriorTalents{},
		WarriorInputs:     inputs,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassWarrior),
	}
	core.FillTalentsProto(warrior.Talents.ProtoReflect(), talents)

	warrior.PseudoStats.CanParry = true
	warrior.PrecisionKnown = false

	warrior.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])
	// Dodge no longer granted from agility
	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	warrior.AddStat(stats.ParryRating, -warrior.GetBaseStats()[stats.Strength]*0.27) // Does not apply to base Strength
	warrior.AddStatDependency(stats.Strength, stats.ParryRating, 0.27)               // Change from block to pary in mop (4.2 Changed from 25->27 percent)
	warrior.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// Base dodge unaffected by Diminishing Returns
	warrior.PseudoStats.BaseBlockChance += 0.05
	warrior.PseudoStats.BaseDodgeChance += 0.03664
	warrior.PseudoStats.BaseParryChance += 0.05
	warrior.CriticalBlockChance = append(warrior.CriticalBlockChance, 0.0, 0.0)

	if mh, oh := warrior.MainHand(), warrior.OffHand(); mh.Name == "Gurthalak, Voice of the Deeps" || oh.Name == "Gurthalak, Voice of the Deeps" {
		warrior.gurthalakTentacles = make([]*cata.TentacleOfTheOldOnesPet, 10)

		for i := 0; i < 10; i++ {
			warrior.gurthalakTentacles[i] = warrior.NewTentacleOfTheOldOnesPet()
		}
	}

	return warrior
}

func (warrior *Warrior) HasMajorGlyph(glyph proto.WarriorMajorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) HasMinorGlyph(glyph proto.WarriorMinorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

// Shared cooldown for Deadly Calm and Recklessness Activation
func (warrior *Warrior) RecklessnessDeadlyCalmLock() *core.Timer {
	return warrior.Character.GetOrInitTimer(&warrior.recklessnessDeadlyCalmCD)
}

func (warrior *Warrior) GetCriticalBlockChance() float64 {
	return warrior.CriticalBlockChance[0] + warrior.CriticalBlockChance[1]
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
