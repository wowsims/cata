package warrior

import (
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
	SpellMaskCommandingShout
	// SpellMaskBerserkerRage
	// SpellMaskRecklessness
	// SpellMaskShieldWall
	// SpellMaskLastStand
	// SpellMaskDeadlyCalm
	// SpellMaskCharge

	// Abilities that cost rage but aren't attacks
	// SpellMaskDemoShout
	// SpellMaskInnerRage
	// SpellMaskShieldBlock
	// SpellMaskDeathWish
	// SpellMaskSweepingStrikes

	// Special attacks
	// SpellMaskCleave
	SpellMaskColossusSmash
	// SpellMaskExecute
	// SpellMaskHeroicStrike
	// SpellMaskHeroicThrow
	// SpellMaskOverpower
	// SpellMaskRend
	// SpellMaskRevenge
	// SpellMaskShatteringThrow
	// SpellMaskSlam
	// SpellMaskSunderArmor
	// SpellMaskThunderClap
	// SpellMaskWhirlwind
	// SpellMaskWhirlwindOh
	SpellMaskShieldSlam
	// SpellMaskConcussionBlow
	SpellMaskDevastate
	// SpellMaskShockwave
	// SpellMaskVictoryRush
	SpellMaskBloodthirst
	// SpellMaskRagingBlow
	SpellMaskMortalStrike
	// SpellMaskBladestorm
	// SpellMaskHeroicLeap

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
	Stance                 Stance
	EnrageEffectMultiplier float64
	CriticalBlockChance    []float64 // Can be gained as non-prot via certain talents and spells

	BattleShout     *core.Spell
	CommandingShout *core.Spell
	BattleStance    *core.Spell
	DefensiveStance *core.Spell
	BerserkerStance *core.Spell

	// BerserkerRage     *core.Spell
	// ColossusSmash     *core.Spell
	// DemoralizingShout *core.Spell
	// Execute           *core.Spell
	// Overpower         *core.Spell
	// Rend              *core.Spell
	// Revenge           *core.Spell
	// ShieldBlock       *core.Spell
	// Slam              *core.Spell
	// SunderArmor       *core.Spell
	// ThunderClap       *core.Spell
	// Whirlwind         *core.Spell
	// DeepWounds        *core.Spell
	// Charge            *core.Spell
	// ChargeAura        *core.Aura

	shoutsCD   *core.Timer
	hsCleaveCD *core.Timer
	// HeroicStrike             *core.Spell
	// Cleave                   *core.Spell

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	EnrageAura        *core.Aura
	BerserkerRageAura *core.Aura
	// BloodsurgeAura    *core.Aura
	// SuddenDeathAura   *core.Aura
	// ShieldBlockAura   *core.Aura
	// ThunderstruckAura *core.Aura
	// InnerRageAura     *core.Aura

	DemoralizingShoutAuras core.AuraArray
	SunderArmorAuras       core.AuraArray
	ThunderClapAuras       core.AuraArray
	ColossusSmashAuras     core.AuraArray
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	// if warrior.Talents.Rampage {
	// 	raidBuffs.Rampage = true
	// }
}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.EnrageEffectMultiplier = 1.0
	warrior.hsCleaveCD = warrior.NewTimer()
	warrior.shoutsCD = warrior.NewTimer()

	warrior.registerStances()
	warrior.registerShouts()
	warrior.registerPassives()
	// warrior.registerBerserkerRageSpell()
	// warrior.registerColossusSmash()
	// warrior.registerDemoralizingShoutSpell()
	// warrior.registerExecuteSpell()
	// warrior.registerHeroicStrikeSpell()
	// warrior.registerCleaveSpell()
	// warrior.registerHeroicLeap()
	// warrior.registerHeroicThrow()
	// warrior.registerInnerRage()
	// warrior.registerOverpowerSpell()
	// warrior.registerRecklessnessCD()
	// warrior.registerRendSpell()
	// warrior.registerRevengeSpell()
	// warrior.registerShatteringThrowCD()
	// warrior.registerShieldBlockCD()
	// warrior.registerShieldWallCD()
	// warrior.registerSlamSpell()
	// warrior.registerSunderArmor()
	// warrior.registerThunderClapSpell()
	// warrior.registerWhirlwindSpell()
	// warrior.registerCharge()
}

func (war *Warrior) registerPassives() {
	war.registerEnrage()
}

func (warrior *Warrior) Reset(_ *core.Simulation) {
	warrior.Stance = StanceNone
}

func NewWarrior(character *core.Character, options *proto.WarriorOptions, talents string, inputs WarriorInputs) *Warrior {
	warrior := &Warrior{
		Character:         *character,
		Talents:           &proto.WarriorTalents{},
		WarriorInputs:     inputs,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassWarrior),
	}
	core.FillTalentsProto(warrior.Talents.ProtoReflect(), talents)

	warrior.EnableRageBar(core.RageBarOptions{
		StartingRage:       options.StartingRage,
		BaseRageMultiplier: 1,
	})

	warrior.EnableAutoAttacks(warrior, core.AutoAttackOptions{
		MainHand:       warrior.WeaponFromMainHand(warrior.DefaultCritMultiplier()),
		OffHand:        warrior.WeaponFromOffHand(warrior.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	warrior.PseudoStats.CanParry = true

	warrior.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[character.Class])
	warrior.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	strengthToParryRating := (1 / 951.158596) * core.ParryRatingPerParryPercent
	warrior.AddStat(stats.ParryRating, -warrior.GetBaseStats()[stats.Strength]*strengthToParryRating) // Does not apply to base Strength
	warrior.AddStatDependency(stats.Strength, stats.ParryRating, strengthToParryRating)
	warrior.AddStatDependency(stats.Agility, stats.DodgeRating, 0.1/10000.0/100.0)
	warrior.AddStatDependency(stats.BonusArmor, stats.Armor, 1)

	// Base dodge unaffected by Diminishing Returns
	warrior.PseudoStats.BaseDodgeChance += 0.03
	warrior.PseudoStats.BaseParryChance += 0.03
	warrior.PseudoStats.BaseBlockChance += 0.03
	warrior.CriticalBlockChance = append(warrior.CriticalBlockChance, 0.0, 0.0)

	return warrior
}

func (warrior *Warrior) HasMajorGlyph(glyph proto.WarriorMajorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) HasMinorGlyph(glyph proto.WarriorMinorGlyph) bool {
	return warrior.HasGlyph(int32(glyph))
}

func (warrior *Warrior) GetCriticalBlockChance() float64 {
	return warrior.CriticalBlockChance[0] + warrior.CriticalBlockChance[1]
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
