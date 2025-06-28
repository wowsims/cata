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
)

const (
	SpellMaskNone int64 = 0
	// Abilities that don't cost rage and aren't attacks
	SpellMaskBattleShout int64 = 1 << iota
	SpellMaskCommandingShout
	SpellMaskBerserkerRage
	SpellMaskRallyingCry
	SpellMaskRecklessness
	SpellMaskShieldWall
	SpellMaskLastStand
	SpellMaskCharge
	SpellMaskSkullBanner
	SpellMaskDemoralizingBanner
	SpellMaskAvatar
	SpellMaskDemoralizingShout

	// Special attacks
	SpellMaskSweepingStrikes
	SpellMaskSweepingStrikesHit
	SpellMaskSweepingStrikesNormalizedHit
	SpellMaskCleave
	SpellMaskColossusSmash
	SpellMaskExecute
	SpellMaskHeroicStrike
	SpellMaskHeroicThrow
	SpellMaskOverpower
	SpellMaskRevenge
	SpellMaskShatteringThrow
	SpellMaskSlam
	SpellMaskSweepingSlam
	SpellMaskSunderArmor
	SpellMaskThunderClap
	SpellMaskWhirlwind
	SpellMaskWhirlwindOh
	SpellMaskShieldBarrier
	SpellMaskShieldSlam
	SpellMaskDevastate
	SpellMaskBloodthirst
	SpellMaskRagingBlow
	SpellMaskRagingBlowMH
	SpellMaskRagingBlowOH
	SpellMaskMortalStrike
	SpellMaskHeroicLeap
	SpellMaskWildStrike
	SpellMaskShieldBlock

	// Talents
	SpellMaskImpendingVictory
	SpellMaskBladestorm
	SpellMaskBladestormMH
	SpellMaskBladestormOH
	SpellMaskDragonRoar
	SpellMaskBloodbath
	SpellMaskBloodbathDot
	SpellMaskStormBolt
	SpellMaskStormBoltOH
	SpellMaskShockwave

	SpellMaskShouts = SpellMaskCommandingShout | SpellMaskBattleShout
)

const EnrageTag = "EnrageEffect"

type Warrior struct {
	core.Character

	ClassSpellScaling float64

	Talents *proto.WarriorTalents

	WarriorInputs

	// Current state
	Stance              Stance
	CriticalBlockChance []float64 // Can be gained as non-prot via certain talents and spells

	HeroicStrikeCleaveCostMod *core.SpellMod

	BattleShout     *core.Spell
	CommandingShout *core.Spell
	BattleStance    *core.Spell
	DefensiveStance *core.Spell
	BerserkerStance *core.Spell

	ColossusSmash                   *core.Spell
	MortalStrike                    *core.Spell
	DeepWounds                      *core.Spell
	ShieldSlam                      *core.Spell
	SweepingStrikesNormalizedAttack *core.Spell

	sharedShoutsCD   *core.Timer
	sharedHSCleaveCD *core.Timer

	BattleStanceAura    *core.Aura
	DefensiveStanceAura *core.Aura
	BerserkerStanceAura *core.Aura

	InciteAura          *core.Aura
	UltimatumAura       *core.Aura
	SweepingStrikesAura *core.Aura
	EnrageAura          *core.Aura
	BerserkerRageAura   *core.Aura
	ShieldBlockAura     *core.Aura
	LastStandAura       *core.Aura
	RallyingCryAura     *core.Aura
	VictoryRushAura     *core.Aura
	SwordAndBoardAura   *core.Aura
	ShieldBarrierAura   *core.DamageAbsorptionAura

	SkullBannerAura         *core.Aura
	DemoralizingBannerAuras core.AuraArray

	DemoralizingShoutAuras core.AuraArray
	SunderArmorAuras       core.AuraArray
	ThunderClapAuras       core.AuraArray
	ColossusSmashAuras     core.AuraArray
	WeakenedArmorAuras     core.AuraArray

	// Set Bonuses
	T14Tank2P *core.Aura
	T15Tank2P *core.Aura
	T15Tank4P *core.Aura
	T16Dps4P  *core.Aura
}

func (warrior *Warrior) GetCharacter() *core.Character {
	return &warrior.Character
}

func (warrior *Warrior) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {

}

func (warrior *Warrior) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (warrior *Warrior) Initialize() {
	warrior.sharedHSCleaveCD = warrior.NewTimer()
	warrior.sharedShoutsCD = warrior.NewTimer()

	warrior.WeakenedArmorAuras = warrior.NewEnemyAuraArray(core.WeakenedArmorAura)

	warrior.registerStances()
	warrior.registerShouts()
	warrior.registerPassives()
	warrior.registerBanners()
	warrior.ApplyGlyphs()

	warrior.registerBerserkerRage()
	warrior.registerRallyingCry()
	warrior.registerColossusSmash()
	warrior.registerExecuteSpell()
	warrior.registerHeroicStrikeSpell()
	warrior.registerCleaveSpell()
	warrior.registerHeroicLeap()
	warrior.registerHeroicThrow()
	warrior.registerRecklessness()
	warrior.registerVictoryRush()
	warrior.registerShatteringThrow()
	warrior.registerShieldWall()
	warrior.registerSunderArmor()
	warrior.registerThunderClap()
	warrior.registerWhirlwind()
	warrior.registerCharge()
}

func (warrior *Warrior) registerPassives() {
	warrior.registerEnrage()
	warrior.registerDeepWounds()
	warrior.registerBloodAndThunder()
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
		MaxRage:            core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfUnendingRage), 120, 100),
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

	warrior.AddStat(stats.ParryRating, -warrior.GetBaseStats()[stats.Strength]*core.StrengthToParryRating) // Does not apply to base Strength
	warrior.AddStatDependency(stats.Strength, stats.ParryRating, core.StrengthToParryRating)
	warrior.AddStatDependency(stats.Agility, stats.DodgeRating, 0.1/10000.0/100.0)
	warrior.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	warrior.MultiplyStat(stats.HasteRating, 1.5)

	// Base dodge unaffected by Diminishing Returns
	warrior.PseudoStats.BaseDodgeChance += 0.03
	warrior.PseudoStats.BaseParryChance += 0.03
	warrior.PseudoStats.BaseBlockChance += 0.03
	warrior.CriticalBlockChance = append(warrior.CriticalBlockChance, 0.0, 0.0)

	warrior.HeroicStrikeCleaveCostMod = warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskHeroicStrike | SpellMaskCleave,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -2,
	})

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

// Used for T15 Protection 4P bonus.
func (warrior *Warrior) GetRageMultiplier(target *core.Unit) float64 {
	// At the moment only protection warriors use this bonus.
	if warrior.Spec != proto.Spec_SpecProtectionWarrior {
		return 1.0
	}

	if warrior.T15Tank4P != nil && warrior.T15Tank4P.IsActive() && warrior.DemoralizingShoutAuras.Get(target).IsActive() {
		return 1.5
	}

	return 1.0
}

func (warrior *Warrior) CastNormalizedSweepingStrikesAttack(results []*core.SpellResult, sim *core.Simulation) {
	if warrior.SweepingStrikesAura != nil && warrior.SweepingStrikesAura.IsActive() {
		for _, result := range results {
			if result.Landed() {
				warrior.SweepingStrikesNormalizedAttack.Cast(sim, warrior.Env.NextTargetUnit(result.Target))
				break
			}
		}
	}
}

// Agent is a generic way to access underlying warrior on any of the agents.
type WarriorAgent interface {
	GetWarrior() *Warrior
}
