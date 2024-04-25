package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var TalentTreeSizes = [3]int{19, 19, 20}

// Start looking to refresh 5 minute totems at 4:55.
const TotemRefreshTime5M = time.Second * 295

const (
	SpellFlagShock     = core.SpellFlagAgentReserved1
	SpellFlagElectric  = core.SpellFlagAgentReserved2
	SpellFlagTotem     = core.SpellFlagAgentReserved3
	SpellFlagFocusable = core.SpellFlagAgentReserved4
)

func NewShaman(character *core.Character, talents string, totems *proto.ShamanTotems, selfBuffs SelfBuffs, thunderstormRange bool) *Shaman {
	shaman := &Shaman{
		Character:           *character,
		Talents:             &proto.ShamanTalents{},
		Totems:              totems,
		SelfBuffs:           selfBuffs,
		ThunderstormInRange: thunderstormRange,
	}
	// shaman.waterShieldManaMetrics = shaman.NewManaMetrics(core.ActionID{SpellID: 57960})

	core.FillTalentsProto(shaman.Talents.ProtoReflect(), talents, TalentTreeSizes)

	// Add Shaman stat dependencies
	shaman.AddStatDependency(stats.BonusArmor, stats.Armor, 1)
	shaman.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[shaman.Class]*core.CritRatingPerCritChance)
	shaman.EnableManaBarWithModifier(1.0)
	if shaman.Spec == proto.Spec_SpecEnhancementShaman {
		shaman.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*6)
		shaman.AddStatDependency(stats.AttackPower, stats.SpellPower, 0.55)
		shaman.AddStatDependency(stats.Agility, stats.AttackPower, 2.0)
		shaman.AddStatDependency(stats.Strength, stats.AttackPower, 1.0)
		shaman.PseudoStats.CanParry = true
	} else if shaman.Spec == proto.Spec_SpecElementalShaman {
		shaman.AddStatDependency(stats.Agility, stats.AttackPower, 2.0)
	}

	if selfBuffs.Shield == proto.ShamanShield_WaterShield {
		shaman.AddStat(stats.MP5, 354)
	}

	shaman.FireElemental = shaman.NewFireElemental(float64(totems.BonusSpellpower))
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

	LightningShield     *core.Spell
	LightningShieldAura *core.Aura
	Fulmination         *core.Spell

	Earthquake   *core.Spell
	Thunderstorm *core.Spell

	EarthShock    *core.Spell
	FlameShock    *core.Spell
	FlameShockDot *core.Spell
	FrostShock    *core.Spell

	FeralSpirit  *core.Spell
	SpiritWolves *SpiritWolves

	FireElemental      *FireElemental
	FireElementalTotem *core.Spell

	MagmaTotem           *core.Spell
	ManaSpringTotem      *core.Spell
	HealingStreamTotem   *core.Spell
	SearingTotem         *core.Spell
	StrengthOfEarthTotem *core.Spell
	TremorTotem          *core.Spell
	StoneskinTotem       *core.Spell
	WindfuryTotem        *core.Spell
	WrathOfAirTotem      *core.Spell
	FlametongueTotem     *core.Spell

	UnleashElements *core.Spell
	UnleashLife     *core.Spell
	UnleashFlame    *core.Spell
	UnleashFrost    *core.Spell
	UnleashWind     *core.Spell

	MaelstromWeaponAura *core.Aura
	SearingFlames       *core.Spell

	// Healing Spells
	tidalWaveProc          *core.Aura
	ancestralHealingAmount float64
	AncestralAwakening     *core.Spell
	HealingSurge           *core.Spell
	GreaterHealingWave     *core.Spell
	HealingWave            *core.Spell
	ChainHeal              *core.Spell
	Riptide                *core.Spell
	EarthShield            *core.Spell

	waterShieldManaMetrics *core.ResourceMetrics
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

func (shaman *Shaman) HasPrimeGlyph(glyph proto.ShamanPrimeGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}
func (shaman *Shaman) HasMajorGlyph(glyph proto.ShamanMajorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}
func (shaman *Shaman) HasMinorGlyph(glyph proto.ShamanMinorGlyph) bool {
	return shaman.HasGlyph(int32(glyph))
}

func (shaman *Shaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {

	if shaman.Totems.Fire != proto.FireTotem_NoFireTotem {
		raidBuffs.TotemicWrath = true
	}

	if shaman.Totems.Fire == proto.FireTotem_FlametongueTotem {
		raidBuffs.FlametongueTotem = true
	}

	if shaman.Totems.Water == proto.WaterTotem_ManaSpringTotem {
		raidBuffs.ManaSpringTotem = true
	}

	if shaman.Talents.ManaTideTotem {
		raidBuffs.ManaTideTotemCount++
	}

	switch shaman.Totems.Air {
	case proto.AirTotem_WrathOfAirTotem:
		raidBuffs.WrathOfAirTotem = true
	case proto.AirTotem_WindfuryTotem:
		raidBuffs.WindfuryTotem = true
	}

	switch shaman.Totems.Earth {
	case proto.EarthTotem_StrengthOfEarthTotem:
		raidBuffs.StrengthOfEarthTotem = true
	case proto.EarthTotem_StoneskinTotem:
		raidBuffs.StoneskinTotem = true
	}

	if shaman.Talents.UnleashedRage > 0 {
		raidBuffs.UnleashedRage = true
	}

	if shaman.Talents.ElementalOath > 0 {
		raidBuffs.ElementalOath = true
	}
}

func (shaman *Shaman) Initialize() {
	shaman.registerChainLightningSpell()
	shaman.registerFireElementalTotem()
	shaman.registerFireNovaSpell()
	shaman.registerLavaBurstSpell()
	shaman.registerLightningBoltSpell()
	shaman.registerLightningShieldSpell()
	shaman.registerMagmaTotemSpell()
	shaman.registerSearingTotemSpell()
	shaman.registerShocks()
	shaman.registerUnleashElements()

	shaman.registerStrengthOfEarthTotemSpell()
	shaman.registerFlametongueTotemSpell()
	shaman.registerTremorTotemSpell()
	shaman.registerStoneskinTotemSpell()
	shaman.registerWindfuryTotemSpell()
	shaman.registerWrathOfAirTotemSpell()
	shaman.registerManaSpringTotemSpell()
	shaman.registerHealingStreamTotemSpell()

	// // This registration must come after all the totems are registered
	shaman.registerCallOfTheElements()

	shaman.registerBloodlustCD()
	// shaman.NewTemporaryStatsAura("DC Pre-Pull SP Proc", core.ActionID{SpellID: 60494}, stats.Stats{stats.SpellPower: 765}, time.Second*10)
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

func (shaman *Shaman) GetOverloadChance() float64 {
	overloadChance := 0.0

	if shaman.Spec == proto.Spec_SpecElementalShaman {
		masteryPoints := shaman.GetMasteryPoints()
		overloadChance = 0.16 + masteryPoints*0.02
	}

	return overloadChance
}

func (shaman *Shaman) GetMentalQuicknessBonus() float64 {
	mentalQuicknessBonus := 0.0

	if shaman.Spec == proto.Spec_SpecEnhancementShaman {
		mentalQuicknessBonus += 0.55
	}

	return mentalQuicknessBonus
}

const (
	SpellMaskNone               int64 = 0
	SpellMaskFireElementalTotem int64 = 1 << iota
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
	SpellMaskStormstrike
	SpellMaskEarthShield
	SpellMaskFulmination
	SpellMaskFrostShock
	SpellMaskUnleashFrost
	SpellMaskUnleashFlame
	SpellMaskEarthquake

	SpellMaskFlameShock = SpellMaskFlameShockDirect | SpellMaskFlameShockDot
	SpellMaskFire       = SpellMaskFlameShock | SpellMaskLavaBurst | SpellMaskLavaBurstOverload | SpellMaskLavaLash | SpellMaskFireNova | SpellMaskUnleashFlame
	SpellMaskNature     = SpellMaskLightningBolt | SpellMaskLightningBoltOverload | SpellMaskChainLightning | SpellMaskChainLightningOverload | SpellMaskEarthShock | SpellMaskThunderstorm | SpellMaskFulmination
	SpellMaskFrost      = SpellMaskUnleashFrost | SpellMaskFrostShock
	SpellMaskOverload   = SpellMaskLavaBurstOverload | SpellMaskLightningBoltOverload | SpellMaskChainLightningOverload
)
