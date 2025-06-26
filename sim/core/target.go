package core

import (
	"strconv"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type Encounter struct {
	Duration          time.Duration
	DurationVariation time.Duration
	AllTargets        []*Target
	ActiveTargets     []*Target
	AllTargetUnits    []*Unit
	ActiveTargetUnits []*Unit

	ExecuteProportion_20 float64
	ExecuteProportion_25 float64
	ExecuteProportion_35 float64
	ExecuteProportion_90 float64

	EndFightAtHealth float64
	// DamageTaken is used to track health fights instead of duration fights.
	//  Once primary target has taken its health worth of damage, fight ends.
	DamageTaken float64
	// In health fight: set to true until we get something to base on
	DurationIsEstimate bool

	// Value to multiply by, for damage spells which are subject to the aoe cap.
	aoeCapMultiplier float64
}

func NewEncounter(options *proto.Encounter) Encounter {
	options.ExecuteProportion_25 = max(options.ExecuteProportion_25, options.ExecuteProportion_20)
	options.ExecuteProportion_35 = max(options.ExecuteProportion_35, options.ExecuteProportion_25)
	totalTargetCount := max(len(options.Targets), 1)

	encounter := Encounter{
		Duration:             DurationFromSeconds(options.Duration),
		DurationVariation:    DurationFromSeconds(options.DurationVariation),
		ExecuteProportion_20: max(options.ExecuteProportion_20, 0),
		ExecuteProportion_25: max(options.ExecuteProportion_25, 0),
		ExecuteProportion_35: max(options.ExecuteProportion_35, 0),
		ExecuteProportion_90: max(options.ExecuteProportion_90, 0),
		AllTargets:           make([]*Target, 0, totalTargetCount),
		ActiveTargets:        make([]*Target, 0, totalTargetCount),
		AllTargetUnits:       make([]*Unit, 0, totalTargetCount),
		ActiveTargetUnits:    make([]*Unit, 0, totalTargetCount),
	}

	for targetIndex, targetOptions := range options.Targets {
		target := NewTarget(targetOptions, int32(targetIndex))
		encounter.AllTargets = append(encounter.AllTargets, target)
		encounter.AllTargetUnits = append(encounter.AllTargetUnits, &target.Unit)

		if target.IsEnabled() {
			encounter.ActiveTargets = append(encounter.ActiveTargets, target)
			encounter.ActiveTargetUnits = append(encounter.ActiveTargetUnits, &target.Unit)
		}
	}

	if len(encounter.AllTargets) == 0 {
		// Add a dummy target. The only case where targets aren't specified is when
		// computing character stats, and targets won't matter there.
		target := NewTarget(&proto.Target{}, 0)
		encounter.AllTargets = append(encounter.AllTargets, target)
		encounter.ActiveTargets = append(encounter.ActiveTargets, target)
		encounter.AllTargetUnits = append(encounter.AllTargetUnits, &target.Unit)
		encounter.ActiveTargetUnits = append(encounter.ActiveTargetUnits, &target.Unit)
	}

	if len(encounter.ActiveTargets) == 0 {
		panic("At least one target must be active at the start of the simulation!")
	}

	// If UseHealth is set, we use the sum of targets health. After creating the targets to make sure stat modifications are done
	if options.UseHealth {
		for _, t := range options.Targets {
			encounter.EndFightAtHealth += t.Stats[stats.Health]
		}
		if encounter.EndFightAtHealth == 0 {
			encounter.EndFightAtHealth = 1 // default to something so we don't instantly end without anything.
		}
	}

	if encounter.EndFightAtHealth > 0 {
		// Until we pre-sim set duration to 10m
		encounter.Duration = time.Minute * 10
		encounter.DurationIsEstimate = true
	}

	encounter.updateAOECapMultiplier()

	return encounter
}

func (encounter *Encounter) AOECapMultiplier() float64 {
	return encounter.aoeCapMultiplier
}
func (encounter *Encounter) updateAOECapMultiplier() {
	encounter.aoeCapMultiplier = min(10/float64(len(encounter.ActiveTargets)), 1)
}

func (encounter *Encounter) doneIteration(sim *Simulation) {
	for _, target := range encounter.AllTargets {
		target.doneIteration(sim)
	}
}

func (encounter *Encounter) GetMetricsProto() *proto.EncounterMetrics {
	metrics := &proto.EncounterMetrics{
		Targets: make([]*proto.UnitMetrics, len(encounter.AllTargets)),
	}

	for idx, target := range encounter.AllTargets {
		metrics.Targets[idx] = target.GetMetricsProto()
	}

	return metrics
}

// Target is an enemy/boss that can be the target of player attacks/spells.
type Target struct {
	Unit

	AI TargetAI
}

func NewTarget(options *proto.Target, targetIndex int32) *Target {
	unitStats := stats.Stats{}
	if options.Stats != nil {
		unitStats = stats.FromProtoArray(options.Stats)
	}

	target := &Target{
		Unit: Unit{
			Type:        EnemyUnit,
			Index:       targetIndex,
			Label:       "Target " + strconv.Itoa(int(targetIndex)+1),
			Level:       options.Level,
			MobType:     options.MobType,
			auraTracker: newAuraTracker(),
			stats:       unitStats,
			PseudoStats: stats.NewPseudoStats(),
			Metrics:     NewUnitMetrics(),

			StatDependencyManager: stats.NewStatDependencyManager(),
			ReactionTime:          time.Millisecond * 1620,
			enabled:               !options.DisabledAtStart,
		},
	}
	defaultRaidBossLevel := int32(CharacterLevel + 3)
	target.GCD = target.NewTimer()
	target.RotationTimer = target.NewTimer()
	if target.Level == 0 {
		target.Level = defaultRaidBossLevel
	}

	// Default Crit chance for NPCs depends only on their level relative to the level of their
	// player target. If there is a need to model a custom Crit % for a Target, this can be
	// accomplished by specifying the extra (or reduced) Crit within the CritRating field of the
	// proto stats array. Any specified CritRating will be applied as an OFFSET to the default
	// level-based values, rather than as a replacement.
	target.stats[stats.PhysicalCritPercent] = UnitLevelFloat64(target.Level, 5.0, 5.2, 5.4, 5.6)
	target.addUniversalStatDependencies()

	target.PseudoStats.CanBlock = true
	target.PseudoStats.CanParry = true
	target.PseudoStats.ParryHaste = options.ParryHaste
	target.PseudoStats.InFrontOfTarget = true
	target.PseudoStats.DamageSpread = options.DamageSpread

	preset := GetPresetTargetWithID(options.Id)
	if preset != nil && preset.AI != nil {
		target.AI = preset.AI()
	}

	return target
}

func (target *Target) Reset(sim *Simulation) {
	target.Unit.reset(sim, nil)
	target.CurrentTarget = target.defaultTarget

	target.SetGCDTimer(sim, 0)
	if target.AI != nil {
		target.AI.Reset(sim)
	}
}

func (target *Target) NextActiveTarget() *Target {
	nextIndex := target.Index + 1

	if nextIndex >= target.Env.TotalTargetCount() {
		nextIndex = 0
	}

	nextTarget := target.Env.GetTargetByIndex(nextIndex)

	if nextTarget.IsEnabled() {
		return nextTarget
	} else {
		return nextTarget.NextActiveTarget()
	}
}

func (target *Target) GetMetricsProto() *proto.UnitMetrics {
	metrics := target.Metrics.ToProto()
	metrics.Name = target.Label
	metrics.UnitIndex = target.UnitIndex
	metrics.Auras = target.auraTracker.GetMetricsProto()
	return metrics
}

type DynamicDamageDoneByCaster func(sim *Simulation, spell *Spell, attackTable *AttackTable) float64

// Holds cached values for outcome/damage calculations, for a specific attacker+defender pair.
// These are updated dynamically when attacker or defender stats change.
type AttackTable struct {
	Attacker *Unit
	Defender *Unit

	BaseMissChance      float64
	BaseSpellMissChance float64
	BaseBlockChance     float64
	BaseDodgeChance     float64
	BaseParryChance     float64
	BaseGlanceChance    float64

	GlanceMultiplier     float64
	MeleeCritSuppression float64
	SpellCritSuppression float64

	DamageDealtMultiplier        float64 // attacker buff, applied in applyAttackerModifiers()
	DamageTakenMultiplier        float64 // defender debuff, applied in applyTargetModifiers()
	HauntSEDamageTakenMultiplier float64
	HealingDealtMultiplier       float64
	IgnoreArmor                  bool    // Ignore defender's armor for specifically this attacker's attacks
	ArmorIgnoreFactor            float64 // Percentage of armor to ignore for this attacker's attacks
	BonusSpellCritPercent        float64 // Analagous to Defender.PseudoStats.BonusSpellCritPercentTaken, but only for this attacker specifically

	// This is for "Apply Aura: Mod Damage Done By Caster" effects.
	// If set, the damage taken multiplier is multiplied by the callbacks result.
	DamageDoneByCasterMultiplier DynamicDamageDoneByCaster

	// When you need more then 1 active, default to using the above one
	// Used with EnableDamageDoneByCaster/DisableDamageDoneByCaster
	DamageDoneByCasterExtraMultiplier []DynamicDamageDoneByCaster
}

func NewAttackTable(attacker *Unit, defender *Unit) *AttackTable {
	table := &AttackTable{
		Attacker: attacker,
		Defender: defender,

		DamageDealtMultiplier:        1,
		DamageTakenMultiplier:        1,
		HauntSEDamageTakenMultiplier: 1,
		HealingDealtMultiplier:       1,
	}

	if defender.Type == EnemyUnit {
		// Assumes attacker (the Player) is level 80.
		table.BaseSpellMissChance = UnitLevelFloat64(defender.Level, 0.04, 0.05, 0.06, 0.17)
		table.BaseMissChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.08)
		table.BaseBlockChance = 0.05
		table.BaseDodgeChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.065)
		table.BaseParryChance = UnitLevelFloat64(defender.Level, 0.05, 0.055, 0.06, 0.14)
		table.BaseGlanceChance = UnitLevelFloat64(defender.Level, 0.06, 0.12, 0.18, 0.24)

		table.GlanceMultiplier = UnitLevelFloat64(defender.Level, 0.95, 0.95, 0.85, 0.75)
		table.MeleeCritSuppression = UnitLevelFloat64(defender.Level, 0, 0.01, 0.02, 0.048)
		table.SpellCritSuppression = UnitLevelFloat64(defender.Level, 0, 0, 0.003, 0.021)
	} else {
		// Assumes defender (the Player) is level 80.
		table.BaseSpellMissChance = 0.05
		table.BaseMissChance = UnitLevelFloat64(attacker.Level, 0.05, 0.048, 0.046, 0.044)
		table.BaseBlockChance = UnitLevelFloat64(attacker.Level, 0.05, 0.048, 0.046, 0.044)
		table.BaseDodgeChance = UnitLevelFloat64(attacker.Level, 0, -0.002, -0.004, -0.006)
		table.BaseParryChance = UnitLevelFloat64(attacker.Level, 0, -0.002, -0.004, -0.006)
	}

	return table
}

func EnableDamageDoneByCaster(index int, maxIndex int, attackTable *AttackTable, handler DynamicDamageDoneByCaster) {
	if attackTable.DamageDoneByCasterExtraMultiplier == nil {
		attackTable.DamageDoneByCasterExtraMultiplier = make([]DynamicDamageDoneByCaster, maxIndex)
	}
	attackTable.DamageDoneByCasterExtraMultiplier[index] = handler
}

func DisableDamageDoneByCaster(index int, attackTable *AttackTable) {
	attackTable.DamageDoneByCasterExtraMultiplier[index] = nil
}
