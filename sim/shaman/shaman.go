package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// Start looking to refresh 5 minute totems at 4:55.
const TotemRefreshTime5M = time.Second * 295

// Damage Done By Caster setup
const (
	DDBC_4pcT12 int = iota
	DDBC_FrostbrandWeapon

	DDBC_Total
)

const (
	SpellFlagShock     = core.SpellFlagAgentReserved1
	SpellFlagElectric  = core.SpellFlagAgentReserved2
	SpellFlagTotem     = core.SpellFlagAgentReserved3
	SpellFlagFocusable = core.SpellFlagAgentReserved4
)

func NewShaman(character *core.Character, talents string, selfBuffs SelfBuffs, thunderstormRange bool) *Shaman {
	shaman := &Shaman{
		Character:           *character,
		Talents:             &proto.ShamanTalents{},
		Totems:              &proto.ShamanTotems{},
		SelfBuffs:           selfBuffs,
		ThunderstormInRange: thunderstormRange,
		ClassSpellScaling:   core.GetClassSpellScalingCoefficient(proto.Class_ClassShaman),
	}
	// shaman.waterShieldManaMetrics = shaman.NewManaMetrics(core.ActionID{SpellID: 57960})

	core.FillTalentsProto(shaman.Talents.ProtoReflect(), talents)

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	shaman.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[shaman.Class])
	shaman.EnableManaBarWithModifier(1.0)
	if shaman.Spec == proto.Spec_SpecEnhancementShaman {
		shaman.AddStat(stats.PhysicalHitPercent, 6)
		shaman.AddStatDependency(stats.AttackPower, stats.SpellPower, 0.55)
		shaman.AddStatDependency(stats.Agility, stats.AttackPower, 2.0)
		shaman.AddStatDependency(stats.Strength, stats.AttackPower, 1.0)
		shaman.PseudoStats.CanParry = true
	} else if shaman.Spec == proto.Spec_SpecElementalShaman {
		shaman.AddStatDependency(stats.Agility, stats.AttackPower, 2.0)
	}

	if selfBuffs.Shield == proto.ShamanShield_WaterShield {
		shaman.AddStat(stats.MP5, 2138)
	}

	shaman.FireElemental = shaman.NewFireElemental()
	shaman.EarthElemental = shaman.NewEarthElemental()

	return shaman
}

// Which buffs this shaman is using.
type SelfBuffs struct {
	Shield  proto.ShamanShield
	ImbueMH proto.ShamanImbue
	ImbueOH proto.ShamanImbue
}

// Indexes into NextTotemDrops for self buffs
const (
	AirTotem int = iota
	EarthTotem
	FireTotem
	WaterTotem
)

// Shaman represents a shaman character.
type Shaman struct {
	core.Character

	ClassSpellScaling float64

	ThunderstormInRange bool // flag if thunderstorm will be in range.

	Talents   *proto.ShamanTalents
	SelfBuffs SelfBuffs

	Totems *proto.ShamanTotems

	// The expiration time of each totem (earth, air, fire, water).
	TotemExpirations [4]time.Duration

	LightningBolt         *core.Spell
	LightningBoltOverload *core.Spell

	ChainLightning          *core.Spell
	ChainLightningHits      []*core.Spell
	ChainLightningOverloads []*core.Spell

	LavaBurst         *core.Spell
	LavaBurstOverload *core.Spell
	FireNova          *core.Spell
	LavaLash          *core.Spell
	Stormstrike       *core.Spell
	PrimalStrike      *core.Spell

	LightningShield       *core.Spell
	LightningShieldDamage *core.Spell
	LightningShieldAura   *core.Aura
	Fulmination           *core.Spell

	ElementalBlast *core.Spell

	Earthquake   *core.Spell
	Thunderstorm *core.Spell

	EarthShock *core.Spell
	FlameShock *core.Spell
	FrostShock *core.Spell

	FeralSpirit *core.Spell

	FireElemental      *FireElemental
	FireElementalTotem *core.Spell

	EarthElemental      *EarthElemental
	EarthElementalTotem *core.Spell

	ElementalSharedCDTimer *core.Timer

	MagmaTotem         *core.Spell
	HealingStreamTotem *core.Spell
	SearingTotem       *core.Spell
	TremorTotem        *core.Spell

	UnleashElements *core.Spell
	UnleashLife     *core.Spell
	UnleashFlame    *core.Spell
	UnleashFrost    *core.Spell
	UnleashWind     *core.Spell

	MaelstromWeaponAura *core.Aura
	SearingFlames       *core.Spell

	SearingFlamesMultiplier float64

	// Healing Spells
	tidalWaveProc          *core.Aura
	ancestralHealingAmount float64
	AncestralAwakening     *core.Spell
	HealingSurge           *core.Spell

	GreaterHealingWave *core.Spell
	HealingWave        *core.Spell
	ChainHeal          *core.Spell
	Riptide            *core.Spell
	EarthShield        *core.Spell

	waterShieldManaMetrics   *core.ResourceMetrics
	VolcanicRegalia4PT12Aura *core.Aura

	// Item sets
	DungeonSet3 *core.Aura
	T12Enh2pc   *core.Aura
	T12Ele4pc   *core.Aura
}

// Implemented by each Shaman spec.
type ShamanAgent interface {
	core.Agent

	// The Shaman controlled by this Agent.
	GetShaman() *Shaman
}

func (shaman *Shaman) GetCharacter() *core.Character {
	return &shaman.Character
}

func (shaman *Shaman) HasMajorGlyph(glyph proto.ShamanMajorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}
func (shaman *Shaman) HasMinorGlyph(glyph proto.ShamanMinorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}

// func (shaman *Shaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {

// 	if shaman.Talents.UnleashedRage > 0 {
// 		raidBuffs.UnleashedRage = true
// 	}

// 	if shaman.Talents.ElementalOath > 0 {
// 		raidBuffs.ElementalOath = true
// 	}
// }

func (shaman *Shaman) Initialize() {
	// shaman.registerChainLightningSpell()
	// shaman.registerFireElementalTotem()
	// shaman.registerEarthElementalTotem()
	// shaman.registerFireNovaSpell()
	// shaman.registerLavaBurstSpell()
	// shaman.registerLightningBoltSpell()
	// shaman.registerLightningShieldSpell()
	shaman.registerSpiritwalkersGraceSpell()
	// shaman.registerMagmaTotemSpell()
	// shaman.registerSearingTotemSpell()
	// shaman.registerShocks()
	// shaman.registerUnleashElements()

	// shaman.registerStrengthOfEarthTotemSpell()
	// shaman.registerFlametongueTotemSpell()
	// shaman.registerTremorTotemSpell()
	// shaman.registerStoneskinTotemSpell()
	// shaman.registerWindfuryTotemSpell()
	// shaman.registerWrathOfAirTotemSpell()
	// shaman.registerManaSpringTotemSpell()
	// shaman.registerHealingStreamTotemSpell()

	// // This registration must come after all the totems are registered
	// shaman.registerCallOfTheElements()
	// shaman.registerCallOfTheAncestors()
	// shaman.registerCallOfTheSpirits()

	shaman.registerBloodlustCD()
}

func (shaman *Shaman) RegisterHealingSpells() {
	// shaman.registerAncestralHealingSpell()
	// shaman.registerHealingSurgeSpell()
	// shaman.registerHealingWaveSpell()
	// shaman.registerRiptideSpell()
	// shaman.registerEarthShieldSpell()
	// shaman.registerChainHealSpell()

	// if shaman.Talents.TidalWaves > 0 {
	// 	shaman.tidalWaveProc = shaman.GetOrRegisterAura(core.Aura{
	// 		Label:    "Tidal Wave Proc",
	// 		ActionID: core.ActionID{SpellID: 53390},
	// 		Duration: core.NeverExpires,
	// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
	// 			aura.Deactivate(sim)
	// 		},
	// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
	// 			shaman.HealingWave.CastTimeMultiplier *= 0.7
	// 			shaman.HealingSurge.BonusCritRating += core.CritRatingPerCritChance * 25
	// 		},
	// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
	// 			shaman.HealingWave.CastTimeMultiplier /= 0.7
	// 			shaman.HealingSurge.BonusCritRating -= core.CritRatingPerCritChance * 25
	// 		},
	// 		MaxStacks: 2,
	// 	})
	// }
}

func (shaman *Shaman) Reset(sim *core.Simulation) {

}

func (shaman *Shaman) calcDamageStormstrikeCritChance(sim *core.Simulation, target *core.Unit, baseDamage float64, spell *core.Spell) *core.SpellResult {
	var result *core.SpellResult
	if target.HasActiveAura("Stormstrike-" + shaman.Label) {
		critPercentBonus := 25.0
		spell.BonusCritPercent += critPercentBonus
		result = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		spell.BonusCritPercent -= critPercentBonus
	} else {
		result = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}
	return result
}

func (shaman *Shaman) GetOverloadChance() float64 {
	overloadChance := 0.0

	if shaman.Spec == proto.Spec_SpecElementalShaman {
		masteryPoints := shaman.GetMasteryPoints()
		overloadChance = 0.16 + masteryPoints*0.02
	}

	return overloadChance
}

const (
	SpellMaskNone               int64 = 0
	SpellMaskFireElementalTotem int64 = 1 << iota
	SpellMaskEarthElementalTotem
	SpellMaskFlameShockDirect
	SpellMaskFlameShockDot
	SpellMaskLavaBurst
	SpellMaskLavaBurstOverload
	SpellMaskLavaLash
	SpellMaskLightningBolt
	SpellMaskLightningBoltOverload
	SpellMaskChainLightning
	SpellMaskChainLightningOverload
	SpellMaskEarthShock
	SpellMaskLightningShield
	SpellMaskThunderstorm
	SpellMaskFireNova
	SpellMaskMagmaTotem
	SpellMaskSearingTotem
	SpellMaskPrimalStrike
	SpellMaskStormstrikeCast
	SpellMaskStormstrikeDamage
	SpellMaskEarthShield
	SpellMaskFulmination
	SpellMaskFrostShock
	SpellMaskUnleashFrost
	SpellMaskUnleashFlame
	SpellMaskEarthquake
	SpellMaskFlametongueWeapon
	SpellMaskFeralSpirit
	SpellMaskElementalMastery
	SpellMaskSpiritwalkersGrace
	SpellMaskShamanisticRage
	SpellMaskElementalBlast

	SpellMaskStormstrike = SpellMaskStormstrikeCast | SpellMaskStormstrikeDamage
	SpellMaskFlameShock  = SpellMaskFlameShockDirect | SpellMaskFlameShockDot
	SpellMaskFire        = SpellMaskFlameShock | SpellMaskLavaBurst | SpellMaskLavaBurstOverload | SpellMaskLavaLash | SpellMaskFireNova | SpellMaskUnleashFlame
	SpellMaskNature      = SpellMaskLightningBolt | SpellMaskLightningBoltOverload | SpellMaskChainLightning | SpellMaskChainLightningOverload | SpellMaskEarthShock | SpellMaskThunderstorm | SpellMaskFulmination
	SpellMaskFrost       = SpellMaskUnleashFrost | SpellMaskFrostShock
	SpellMaskOverload    = SpellMaskLavaBurstOverload | SpellMaskLightningBoltOverload | SpellMaskChainLightningOverload
	SpellMaskShock       = SpellMaskFlameShock | SpellMaskEarthShock | SpellMaskFrostShock
	SpellMaskTotem       = SpellMaskMagmaTotem | SpellMaskSearingTotem | SpellMaskFireElementalTotem | SpellMaskEarthElementalTotem
)
