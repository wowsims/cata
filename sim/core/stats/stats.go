package stats

import (
	"fmt"
	"math"
	"strings"

	"github.com/wowsims/mop/sim/core/proto"
)

type Stats [SimStatsLen]float64

type Stat byte

// Use internal representation instead of proto.Stat so we can add functions
// and use 'byte' as the data type.
//
// This needs to stay synced with proto.Stat: it is okay for SimStatsLen to
// exceed ProtoStatsLen, but the shared indices between the two must match 1:1 .
const (
	Strength Stat = iota
	Agility
	Stamina
	Intellect
	Spirit
	HitRating
	CritRating
	HasteRating
	ExpertiseRating
	DodgeRating
	ParryRating
	MasteryRating
	AttackPower
	RangedAttackPower
	SpellPower
	PvpResilienceRating
	PvpPowerRating
	Armor
	BonusArmor
	Health
	Mana
	MP5
	// end of Stat enum in proto/common.proto

	// The remaining stats below are stored as PseudoStats rather than as
	// Stats in UnitStats proto messages, since they are not required in the
	// database files. However, it is valuable to keep these as proper Stats
	// in the back-end, since they are used in various stat dependencies.
	// The units for all 5 of these are percentages (between 0 and 100).
	PhysicalHitPercent
	SpellHitPercent
	PhysicalCritPercent
	SpellCritPercent
	BlockPercent
	// DO NOT add new stats here without discussing it first; new stats come
	// with a performance penalty.

	SimStatsLen
)

var ProtoStatsLen = len(proto.Stat_name)
var PseudoStatsLen = len(proto.PseudoStat_name)
var UnitStatsLen = ProtoStatsLen + PseudoStatsLen

type SchoolIndex byte

const (
	SchoolIndexNone     SchoolIndex = 0
	SchoolIndexPhysical SchoolIndex = iota
	SchoolIndexArcane
	SchoolIndexFire
	SchoolIndexFrost
	SchoolIndexHoly
	SchoolIndexNature
	SchoolIndexShadow

	SchoolLen
)

func NewSchoolFloatArray() [SchoolLen]float64 {
	return [SchoolLen]float64{
		1, 1, 1, 1, 1, 1, 1, 1,
	}
}

func ProtoArrayToStatsList(protoStats []proto.Stat) []Stat {
	stats := make([]Stat, len(protoStats))
	for i, v := range protoStats {
		stats[i] = Stat(v)
	}
	return stats
}

func IntTupleToStatsList(statType1 int32, statType2 int32, statType3 int32) []Stat {
	statTypes := make([]Stat, 0, 3)

	for _, statIdx := range []int32{statType1, statType2, statType3} {
		if statIdx >= 0 {
			statTypes = append(statTypes, Stat(statIdx))
		}
	}

	return statTypes
}

func (s Stat) StatName() string {
	switch s {
	case Strength:
		return "Strength"
	case Agility:
		return "Agility"
	case Stamina:
		return "Stamina"
	case Intellect:
		return "Intellect"
	case Spirit:
		return "Spirit"
	case HitRating:
		return "HitRating"
	case CritRating:
		return "CritRating"
	case HasteRating:
		return "HasteRating"
	case ExpertiseRating:
		return "ExpertiseRating"
	case DodgeRating:
		return "DodgeRating"
	case ParryRating:
		return "ParryRating"
	case MasteryRating:
		return "MasteryRating"
	case AttackPower:
		return "AttackPower"
	case RangedAttackPower:
		return "RangedAttackPower"
	case SpellPower:
		return "SpellPower"
	case PvpResilienceRating:
		return "PvpResilienceRating"
	case PvpPowerRating:
		return "PvpPowerRating"
	case Armor:
		return "Armor"
	case BonusArmor:
		return "BonusArmor"
	case Health:
		return "Health"
	case Mana:
		return "Mana"
	case MP5:
		return "MP5"
	case PhysicalHitPercent:
		return "PhysicalHitPercent"
	case SpellHitPercent:
		return "SpellHitPercent"
	case PhysicalCritPercent:
		return "PhysicalCritPercent"
	case SpellCritPercent:
		return "SpellCritPercent"
	case BlockPercent:
		return "BlockPercent"
	}

	return "none"
}

func FromProtoArray(values []float64) Stats {
	// SimStatsLen can be larger than ProtoStatsLen, but the built-in copy
	// function will only import the shared indices between the two.
	var stats Stats
	copy(stats[:], values)
	return stats
}

// Runs FromProtoArray() on the stats array embedded in the UnitStats message, but additionally imports any
// PseudoStats that we want to model as proper Stats in the back-end. This allows us to include only essential
// basic properties in database stats arrays, while still letting the back-end Stat enum include derived
// properties when it is computationally convenient to do so (such as for automatically applying stat
// dependencies). Make sure to update this function if you add any back-end Stat entries that are modeled as
// PseudoStats in the front-end.
func FromUnitStatsProto(unitStatsMessage *proto.UnitStats) Stats {
	simStats := FromProtoArray(unitStatsMessage.Stats)

	if unitStatsMessage.PseudoStats != nil {
		pseudoStatsMessage := unitStatsMessage.PseudoStats
		simStats[PhysicalHitPercent] = pseudoStatsMessage[proto.PseudoStat_PseudoStatPhysicalHitPercent]
		simStats[SpellHitPercent] = pseudoStatsMessage[proto.PseudoStat_PseudoStatSpellHitPercent]
		simStats[PhysicalCritPercent] = pseudoStatsMessage[proto.PseudoStat_PseudoStatPhysicalCritPercent]
		simStats[SpellCritPercent] = pseudoStatsMessage[proto.PseudoStat_PseudoStatSpellCritPercent]
		simStats[BlockPercent] = pseudoStatsMessage[proto.PseudoStat_PseudoStatBlockPercent]
	}

	return simStats
}

// Adds two Stats together, returning the new Stats.
func (stats Stats) Add(other Stats) Stats {
	for k := range stats {
		stats[k] += other[k]
	}
	return stats
}

// Adds another to Stats to this, in-place. For performance, only.
func (stats *Stats) AddInplace(other *Stats) {
	for k := range stats {
		stats[k] += other[k]
	}
}

// Subtracts another Stats from this one, returning the new Stats.
func (stats Stats) Subtract(other Stats) Stats {
	for k := range stats {
		stats[k] -= other[k]
	}
	return stats
}

func (stats Stats) Invert() Stats {
	for k, v := range stats {
		stats[k] = -v
	}
	return stats
}

// Rounds all stat values down to the nearest integer, returning the new Stats.
// Used for random suffix stats currently.
func (stats Stats) Floor() Stats {
	for k, v := range stats {
		stats[k] = math.Floor(v)
	}
	return stats
}

func (stats Stats) Multiply(multiplier float64) Stats {
	for k := range stats {
		stats[k] *= multiplier
	}
	return stats
}

// Multiplies two Stats together by multiplying the values of corresponding
// stats, like a dot product operation.
func (stats Stats) DotProduct(other Stats) Stats {
	for k := range stats {
		stats[k] *= other[k]
	}
	return stats
}

// Higher performance variant of the above.
func (stats Stats) ApplyMultipliers(multipliers map[Stat]float64) Stats {
	for k, v := range multipliers {
		stats[k] *= v
	}
	return stats
}

func (stats Stats) Equals(other Stats) bool {
	return stats == other
}

func (stats Stats) EqualsWithTolerance(other Stats, tolerance float64) bool {
	for k, v := range stats {
		if v < other[k]-tolerance || v > other[k]+tolerance {
			return false
		}
	}
	return true
}

// Given an array of Stat types, return the Stat whose value is largest within
// this Stats array.
func (stats Stats) GetHighestStatType(statTypeOptions []Stat) Stat {
	if len(statTypeOptions) < 1 {
		panic("Must supply at least one Stat type option!")
	}

	var highestStatType Stat
	var highestStatValue float64

	for idx, statType := range statTypeOptions {
		if (idx == 0) || (stats[statType] > highestStatValue) {
			highestStatType = statType
			highestStatValue = stats[statType]
		}
	}

	return highestStatType
}

// Returns all Stat types with positive representation in this Stats array.
func (stats Stats) GetBuffedStatTypes() []Stat {
	buffedStatTypes := make([]Stat, 0, SimStatsLen)

	for statIdx, statValue := range stats {
		if statValue > 0 {
			buffedStatTypes = append(buffedStatTypes, Stat(statIdx))
		}
	}

	return buffedStatTypes
}

func (stats Stats) String() string {
	var sb strings.Builder
	sb.WriteString("\n{\n")

	for statIdx, statValue := range stats {
		if statValue == 0 {
			continue
		}
		if name := Stat(statIdx).StatName(); name != "none" {
			_, _ = fmt.Fprintf(&sb, "\t%s: %0.3f,\n", name, statValue)
		}
	}

	sb.WriteString("\n}")
	return sb.String()
}

// Like String() but without the newlines.
func (stats Stats) FlatString() string {
	var sb strings.Builder
	sb.WriteString("{")

	for statIdx, statValue := range stats {
		if statValue == 0 {
			continue
		}
		if name := Stat(statIdx).StatName(); name != "none" {
			_, _ = fmt.Fprintf(&sb, "\"%s\": %0.3f,", name, statValue)
		}
	}

	sb.WriteString("}")
	return sb.String()
}

func (stats Stats) ToProtoArray() []float64 {
	// SimStatsLen can be larger than ProtoStatsLen, so export only the
	// shared indices between the two.
	return stats[:ProtoStatsLen]
}
func (stats Stats) ToProtoMap() map[int32]float64 {
	m := make(map[int32]float64, ProtoStatsLen)
	for i := 0; i < int(ProtoStatsLen); i++ {
		if stats[i] != 0 {
			m[int32(i)] = stats[i]
		}
	}
	return m
}

func FromProtoMap(m map[int32]float64) Stats {
	var stats Stats
	for k, v := range m {
		stats[k] = v

	}
	return stats
}

type PseudoStats struct {
	///////////////////////////////////////////////////
	// Effects that apply when this unit is the attacker.
	///////////////////////////////////////////////////

	SpellCostPercentModifier int32 // Multiplies spell cost.

	CastSpeedMultiplier   float64
	MeleeSpeedMultiplier  float64
	RangedSpeedMultiplier float64
	AttackSpeedMultiplier float64 // Used for real haste effects like Bloodlust that modify resoruce regen and are used for RPPM effects

	SpiritRegenRateCombat float64 // percentage of spirit regen allowed during combat

	// Both of these are currently only used for innervate.
	ForceFullSpiritRegen  bool    // If set, automatically uses full spirit regen regardless of FSR refresh time.
	SpiritRegenMultiplier float64 // Multiplier on spirit portion of mana regen.

	// If true, allows block/parry.
	InFrontOfTarget bool

	// "Apply Aura: Mod Damage Done (Physical)", applies to abilities with EffectSpellCoefficient > 0.
	//  This includes almost all "(Normalized) Weapon Damage", but also some "School Damage (Physical)" abilities.
	BonusDamage float64 // Comes from '+X Weapon Damage' effects

	BonusMHDps     float64
	BonusOHDps     float64
	BonusRangedDps float64

	DisableDWMissPenalty bool // Used by Heroic Strike and Cleave

	ThreatMultiplier float64 // Modulates the threat generated. Affected by things like salv.

	DamageDealtMultiplier          float64            // All damage
	SchoolDamageDealtMultiplier    [SchoolLen]float64 // For specific spell schools (arcane, fire, shadow, etc).
	DotDamageMultiplierAdditive    float64            // All periodic damage
	HealingDealtMultiplier         float64            // All non-shield healing
	PeriodicHealingDealtMultiplier float64            // All periodic healing (on top of HealingDealtMultiplier)
	CritDamageMultiplier           float64            // All multiplicative crit damage

	// Important when unit is attacker or target
	BlockDamageReduction float64

	// Only used for NPCs, governs variance in enemy auto-attack damage
	DamageSpread float64

	///////////////////////////////////////////////////
	// Effects that apply when this unit is the target.
	///////////////////////////////////////////////////

	CanBlock bool
	CanParry bool
	Stunned  bool // prevents blocks, dodges, and parries

	ParryHaste bool

	// Avoidance % not affected by Diminishing Returns, represented as
	// probabilities (between 0 and 1).
	BaseDodgeChance float64
	BaseParryChance float64
	BaseBlockChance float64

	ReducedCritTakenChance float64 // Reduces chance to be crit.

	BonusHealingTaken float64 // Talisman of Troll Divinity

	DamageTakenMultiplier       float64            // All damage
	SchoolDamageTakenMultiplier [SchoolLen]float64 // For specific spell schools (arcane, fire, shadow, etc.)

	DiseaseDamageTakenMultiplier          float64
	PeriodicPhysicalDamageTakenMultiplier float64

	ArmorMultiplier float64 // Major/minor/special multiplicative armor modifiers

	ReducedPhysicalHitTakenChance float64
	ReducedArcaneHitTakenChance   float64
	ReducedFireHitTakenChance     float64
	ReducedFrostHitTakenChance    float64
	ReducedNatureHitTakenChance   float64
	ReducedShadowHitTakenChance   float64

	HealingTakenMultiplier         float64 // All healing sources including self-healing
	ExternalHealingTakenMultiplier float64 // Modulates the output of the individual tank sim healing model
	MovementSpeedMultiplier        float64 // Multiplier for movement speed, default to 1. Player base movement 7 yards/s. All effects affecting movements are multipliers.
}

func NewPseudoStats() PseudoStats {
	return PseudoStats{
		SpellCostPercentModifier: 100,

		CastSpeedMultiplier:   1,
		MeleeSpeedMultiplier:  1,
		RangedSpeedMultiplier: 1,
		AttackSpeedMultiplier: 1,
		SpiritRegenMultiplier: 1,

		ThreatMultiplier: 1,

		DamageDealtMultiplier:          1,
		SchoolDamageDealtMultiplier:    NewSchoolFloatArray(),
		DotDamageMultiplierAdditive:    1,
		HealingDealtMultiplier:         1,
		PeriodicHealingDealtMultiplier: 1,
		CritDamageMultiplier:           1,

		BlockDamageReduction: 0.3,

		DamageSpread: 0.3333,

		// Target effects.
		DamageTakenMultiplier:       1,
		SchoolDamageTakenMultiplier: NewSchoolFloatArray(),

		DiseaseDamageTakenMultiplier:          1,
		PeriodicPhysicalDamageTakenMultiplier: 1,

		ArmorMultiplier: 1,

		HealingTakenMultiplier:         1,
		ExternalHealingTakenMultiplier: 1,
		MovementSpeedMultiplier:        1,
	}
}

type UnitStat int

func (s UnitStat) IsStat() bool               { return int(s) < int(ProtoStatsLen) }
func (s UnitStat) IsPseudoStat() bool         { return !s.IsStat() }
func (s UnitStat) EqualsStat(other Stat) bool { return s.IsStat() && (s.StatIdx() == int(other)) }
func (s UnitStat) EqualsPseudoStat(other proto.PseudoStat) bool {
	return s.IsPseudoStat() && (s.PseudoStatIdx() == int(other))
}
func (s UnitStat) StatIdx() int {
	if !s.IsStat() {
		panic("Is a pseudo stat")
	}
	return int(s)
}
func (s UnitStat) PseudoStatIdx() int {
	if s.IsStat() {
		panic("Is a regular stat")
	}
	return int(s) - int(ProtoStatsLen)
}
func (s UnitStat) AddToStatsProto(p *proto.UnitStats, value float64) {
	if s.IsStat() {
		p.Stats[s.StatIdx()] += value
	} else {
		p.PseudoStats[s.PseudoStatIdx()] += value
	}
}

func UnitStatFromIdx(s int) UnitStat   { return UnitStat(s) }
func UnitStatFromStat(s Stat) UnitStat { return UnitStat(s) }
func UnitStatFromPseudoStat(s proto.PseudoStat) UnitStat {
	return UnitStat(int(s) + int(ProtoStatsLen))
}
