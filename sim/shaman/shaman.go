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
	DDBC_UnleashedFury
	DDBC_2PT16

	DDBC_Total
)

const (
	SpellFlagShamanSpell = core.SpellFlagAgentReserved1
	SpellFlagShock       = core.SpellFlagAgentReserved2
	SpellFlagIsEcho      = core.SpellFlagAgentReserved3
	SpellFlagFocusable   = core.SpellFlagAgentReserved4
)

func NewShaman(character *core.Character, talents string, selfBuffs SelfBuffs, thunderstormRange bool, feleAutocastOptions *proto.FeleAutocastSettings) *Shaman {
	if feleAutocastOptions == nil {
		feleAutocastOptions = &proto.FeleAutocastSettings{
			AutocastFireblast: true,
			AutocastFirenova:  true,
			AutocastImmolate:  true,
			AutocastEmpower:   false,
		}
	}
	shaman := &Shaman{
		Character:           *character,
		Talents:             &proto.ShamanTalents{},
		Totems:              &proto.ShamanTotems{},
		FeleAutocast:        feleAutocastOptions,
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

	shaman.AddStatDependency(stats.Agility, stats.AttackPower, 2.0)
	shaman.AddStat(stats.AttackPower, -20)

	shaman.AddStatDependency(stats.Strength, stats.AttackPower, 1.0)
	shaman.AddStat(stats.AttackPower, -10)

	if selfBuffs.Shield == proto.ShamanShield_WaterShield {
		shaman.AddStat(stats.MP5, 2138)
	}

	shaman.FireElemental = shaman.NewFireElemental(!shaman.Talents.PrimalElementalist)
	shaman.EarthElemental = shaman.NewEarthElemental(!shaman.Talents.PrimalElementalist)

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

	FeleAutocast *proto.FeleAutocastSettings

	// The expiration time of each totem (earth, air, fire, water).
	TotemExpirations [4]time.Duration

	LightningBolt         *core.Spell
	LightningBoltOverload [2]*core.Spell

	ChainLightning          *core.Spell
	ChainLightningOverloads [2][]*core.Spell

	LavaBeam          *core.Spell
	LavaBeamOverloads [2][]*core.Spell

	LavaBurst         *core.Spell
	LavaBurstOverload [2]*core.Spell
	FireNova          *core.Spell
	FireNovas         []*core.Spell
	LavaLash          *core.Spell
	Stormstrike       *core.Spell
	PrimalStrike      *core.Spell
	Stormblast        *core.Spell

	LightningShield       *core.Spell
	LightningShieldDamage *core.Spell
	LightningShieldAura   *core.Aura
	Fulmination           *core.Spell

	ElementalBlast         *core.Spell
	ElementalBlastOverload [2]*core.Spell

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

	Ascendance     *core.Spell
	AscendanceAura *core.Aura

	MagmaTotem         *core.Spell
	HealingStreamTotem *core.Spell
	SearingTotem       *core.Spell
	TremorTotem        *core.Spell

	UnleashElements *core.Spell
	UnleashLife     *core.Spell
	UnleashFlame    *core.Spell
	UnleashFrost    *core.Spell
	UnleashWind     *core.Spell

	MaelstromWeaponAura           *core.Aura
	AncestralSwiftnessInstantAura *core.Aura
	SearingFlames                 *core.Spell

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
	T14Ele4pc   *core.Aura
	T14Enh4pc   *core.Aura
	T15Enh2pc   *core.Aura
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

func (shaman *Shaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.GraceOfAir = true
	raidBuffs.BurningWrath = true
}

func (shaman *Shaman) Initialize() {
	shaman.registerChainLightningSpell()
	shaman.registerFireElementalTotem(!shaman.Talents.PrimalElementalist)
	shaman.registerEarthElementalTotem(!shaman.Talents.PrimalElementalist)
	shaman.registerLightningBoltSpell()
	shaman.registerLightningShieldSpell()
	shaman.registerSpiritwalkersGraceSpell()
	shaman.registerMagmaTotemSpell()
	shaman.registerSearingTotemSpell()
	shaman.registerShocks()
	shaman.registerUnleashElements()
	shaman.registerAscendanceSpell()

	shaman.registerBloodlustCD()
	shaman.registerStormlashCD()
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
		critPercentBonus := core.TernaryFloat64(shaman.T14Enh4pc.IsActive(), 40.0, 25.0)
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
	SpellMaskFireElementalMelee
	SpellMaskFlameShockDirect
	SpellMaskFlameShockDot
	SpellMaskLavaBurst
	SpellMaskLavaBurstOverload
	SpellMaskLavaLash
	SpellMaskLightningBolt
	SpellMaskLightningBoltOverload
	SpellMaskChainLightning
	SpellMaskChainLightningOverload
	SpellMaskLavaBeam
	SpellMaskLavaBeamOverload
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
	SpellMaskUnleashElements
	SpellMaskUnleashFrost
	SpellMaskUnleashFlame
	SpellMaskEarthquake
	SpellMaskFlametongueWeapon
	SpellMaskWindfuryWeapon
	SpellMaskFrostbrandWeapon
	SpellMaskFeralSpirit
	SpellMaskElementalMastery
	SpellMaskAscendance
	SpellMaskSpiritwalkersGrace
	SpellMaskShamanisticRage
	SpellMaskElementalBlast
	SpellMaskElementalBlastOverload
	SpellMaskStormlashTotem
	SpellMaskBloodlust

	SpellMaskStormstrike  = SpellMaskStormstrikeCast | SpellMaskStormstrikeDamage
	SpellMaskFlameShock   = SpellMaskFlameShockDirect | SpellMaskFlameShockDot
	SpellMaskFire         = SpellMaskFlameShock | SpellMaskLavaBurst | SpellMaskLavaBurstOverload | SpellMaskLavaLash | SpellMaskFireNova | SpellMaskUnleashFlame | SpellMaskLavaBeam | SpellMaskLavaBeamOverload | SpellMaskElementalBlast | SpellMaskElementalBlastOverload
	SpellMaskNature       = SpellMaskLightningBolt | SpellMaskLightningBoltOverload | SpellMaskChainLightning | SpellMaskChainLightningOverload | SpellMaskEarthShock | SpellMaskThunderstorm | SpellMaskFulmination | SpellMaskElementalBlast | SpellMaskElementalBlastOverload
	SpellMaskFrost        = SpellMaskUnleashFrost | SpellMaskFrostShock | SpellMaskElementalBlast | SpellMaskElementalBlastOverload
	SpellMaskOverload     = SpellMaskLavaBurstOverload | SpellMaskLightningBoltOverload | SpellMaskChainLightningOverload | SpellMaskElementalBlastOverload | SpellMaskLavaBeamOverload
	SpellMaskShock        = SpellMaskFlameShock | SpellMaskEarthShock | SpellMaskFrostShock
	SpellMaskTotem        = SpellMaskMagmaTotem | SpellMaskSearingTotem | SpellMaskFireElementalTotem | SpellMaskEarthElementalTotem | SpellMaskStormlashTotem
	SpellMaskInstantSpell = SpellMaskAscendance | SpellMaskFeralSpirit | SpellMaskUnleashElements | SpellMaskBloodlust
)
