package core

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type UnitType int
type SpellRegisteredHandler func(spell *Spell)
type OnMasteryStatChanged func(sim *Simulation, oldMasteryRating float64, newMasteryRating float64)
type OnSpeedChanged func(oldSpeed float64, newSpeed float64)
type OnTemporaryStatsChange func(sim *Simulation, buffAura *Aura, statsChangeWithoutDeps stats.Stats)

const (
	PlayerUnit UnitType = iota
	EnemyUnit
	PetUnit
)

type PowerBarType int

const (
	ManaBar PowerBarType = iota
	EnergyBar
	RageBar
	RunicPower
	FocusBar
)

type DynamicDamageTakenModifier func(sim *Simulation, spell *Spell, result *SpellResult, isPeriodic bool)

type GetSpellPowerValue func(spell *Spell) float64
type GetAttackPowerValue func(spell *Spell) float64

// Unit is an abstraction of a Character/Boss/Pet/etc, containing functionality
// shared by all of them.
type Unit struct {
	Type UnitType

	// Index of this unit with its group.
	//  For Players, this is the 0-indexed raid index (0-24).
	//  For Enemies, this is its enemy index.
	//  For Pets, this is the same as the owner's index.
	Index int32

	// Unique index of this unit among all units in the environment.
	// This is used as the index for attack tables.
	UnitIndex int32

	// Unique label for logging.
	Label string

	Level int32 // Level of Unit, e.g. Bosses are 83.

	MobType proto.MobType

	// Amount of time it takes for the human agent to react to in-game events.
	// Used by certain APL values and actions.
	ReactionTime time.Duration

	// Amount of time following a post-GCD channel tick, to when the next action can be performed.
	ChannelClipDelay time.Duration

	// How far this unit is from its target(s). Measured in yards, this is used
	// for calculating spell travel time for certain spells.
	StartDistanceFromTarget float64
	DistanceFromTarget      float64
	Moving                  bool
	movementCallbacks       []MovementCallback
	moveAura                *Aura
	moveSpell               *Spell
	movementAction          *MovementAction

	// Environment in which this Unit exists. This will be nil until after the
	// construction phase.
	Env *Environment

	// Whether this unit is able to perform actions.
	enabled bool

	// Stats this Unit will have at the very start of each Sim iteration.
	// Includes all equipment / buffs / permanent effects but not temporary
	// effects from items / abilities.
	initialStats            stats.Stats
	initialStatsWithoutDeps stats.Stats

	initialPseudoStats stats.PseudoStats

	// Cast speed without any temporary effects.
	initialCastSpeed float64

	// Melee swing speed without any temporary effects.
	initialMeleeSwingSpeed float64

	// Ranged swing speed without any temporary effects.
	initialRangedSwingSpeed float64

	// Provides aura tracking behavior.
	auraTracker

	// Current stats, including temporary effects but not dependencies.
	statsWithoutDeps stats.Stats

	// Current stats, including temporary effects and dependencies.
	stats stats.Stats

	// Provides stat dependency management behavior.
	stats.StatDependencyManager

	PseudoStats stats.PseudoStats

	currentPowerBar PowerBarType
	healthBar
	manaBar
	rageBar
	energyBar
	focusBar
	runicPowerBar

	secondaryResourceBar SecondaryResourceBar

	// All spells that can be cast by this unit.
	Spellbook                 []*Spell
	spellRegistrationHandlers []SpellRegisteredHandler

	// Pets owned by this Unit.
	PetAgents []PetAgent

	DynamicStatsPets      []*Pet
	DynamicMeleeSpeedPets []*Pet
	DynamicCastSpeedPets  []*Pet

	// AutoAttacks is the manager for auto attack swings.
	// Must be enabled to use, with "EnableAutoAttacks()".
	AutoAttacks AutoAttacks

	Rotation *APLRotation

	// Statistics describing the results of the sim.
	Metrics UnitMetrics

	cdTimers []*Timer

	AttackTables                []*AttackTable
	DynamicDamageTakenModifiers []DynamicDamageTakenModifier
	Blockhandler                func(sim *Simulation, spell *Spell, result *SpellResult)
	avoidanceParams             DiminishingReturnsConstants

	GCD *Timer

	// Separate from GCD timer to support spell queueing and off-GCD actions
	RotationTimer *Timer

	// Used for applying the effect of a hardcast spell when casting finishes.
	//  For channeled spells, only Expires is set.
	// No more than one cast may be active at any given time.
	Hardcast Hardcast

	// Rotation-related PendingActions.
	rotationAction *PendingAction
	hardcastAction *PendingAction

	// Cached mana return values per tick.
	manaTickWhileCombat    float64
	manaTickWhileNotCombat float64

	CastSpeed         float64
	meleeAttackSpeed  float64
	rangedAttackSpeed float64

	CurrentTarget   *Unit
	defaultTarget   *Unit
	SecondaryTarget *Unit // Only used for NPCs in tank swap AIs currently.

	// The currently-channeled DOT spell, otherwise nil.
	ChanneledDot *Dot

	// Data about the most recently queued spell, otherwise nil.
	QueuedSpell *QueuedSpell

	// Used for reacting to mastery stat changes if a spec needs it
	OnMasteryStatChanged []OnMasteryStatChanged

	// Used for reacting to cast speed changes if a spec needs it (e.g. for cds reduced by haste)
	OnCastSpeedChanged []OnSpeedChanged

	// Used for reacting to melee attack speed changes if a spec needs it (e.g. for cds reduced by haste)
	OnMeleeAttackSpeedChanged []OnSpeedChanged

	// Used for reacting to ranged attack speed changes if a spec needs it (e.g. for cds reduced by haste)
	OnRangedAttackSpeedChanged []OnSpeedChanged

	// Used for reacting to transient stat changes if a spec needs if (e.g. for caching snapshotting calculations)
	OnTemporaryStatsChanges []OnTemporaryStatsChange

	GetSpellPowerValue GetSpellPowerValue

	GetAttackPowerValue GetAttackPowerValue
}

func (unit *Unit) getSpellPowerValueImpl(spell *Spell) float64 {
	return unit.stats[stats.SpellPower] + spell.BonusSpellPower
}
func (unit *Unit) getAttackPowerValueImpl(spell *Spell) float64 {
	return unit.stats[stats.AttackPower]
}

// Units can be disabled for several reasons:
//  1. Downtime for temporary pets (e.g. Water Elemental)
//  2. Enemy units in various phases (not yet implemented)
//  3. Dead units (not yet implemented)
func (unit *Unit) IsEnabled() bool {
	return unit.enabled
}

func (unit *Unit) IsActive() bool {
	return unit.IsEnabled() && unit.CurrentHealthPercent() > 0
}

func (unit *Unit) IsOpponent(other *Unit) bool {
	return (unit.Type == EnemyUnit) != (other.Type == EnemyUnit)
}

func (unit *Unit) GetOpponents() []*Unit {
	if unit.Type == EnemyUnit {
		return unit.Env.Raid.AllUnits
	} else {
		return unit.Env.Encounter.TargetUnits
	}
}

func (unit *Unit) LogLabel() string {
	return "[" + unit.Label + "]"
}

func (unit *Unit) Log(sim *Simulation, message string, vals ...interface{}) {
	sim.Log(unit.LogLabel()+" "+message, vals...)
}

func (unit *Unit) GetInitialStat(stat stats.Stat) float64 {
	return unit.initialStats[stat]
}
func (unit *Unit) GetStats() stats.Stats {
	return unit.stats
}

// Given an array of Stat types, return the Stat whose value is largest for this
// Unit.
func (unit *Unit) GetHighestStatType(statTypeOptions []stats.Stat) stats.Stat {
	return unit.stats.GetHighestStatType(statTypeOptions)
}
func (unit *Unit) GetStat(stat stats.Stat) float64 {
	return unit.stats[stat]
}
func (unit *Unit) GetMasteryPoints() float64 {
	return MasteryRatingToMasteryPoints(unit.GetStat(stats.MasteryRating))
}
func (unit *Unit) AddStats(stat stats.Stats) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, use AddStatsDynamic instead!")
	}
	unit.stats = unit.stats.Add(stat)
}
func (unit *Unit) AddStat(stat stats.Stat, amount float64) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, use AddStatDynamic instead!")
	}
	unit.stats[stat] += amount
}

func (unit *Unit) AddDynamicDamageTakenModifier(ddtm DynamicDamageTakenModifier) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, cannot add dynamic damage taken modifier!")
	}
	unit.DynamicDamageTakenModifiers = append(unit.DynamicDamageTakenModifiers, ddtm)
}

func (unit *Unit) AddOnMasteryStatChanged(omsc OnMasteryStatChanged) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, cannot add on mastery stat changed callback!")
	}
	unit.OnMasteryStatChanged = append(unit.OnMasteryStatChanged, omsc)
}

func (unit *Unit) AddOnCastSpeedChanged(ocsc OnSpeedChanged) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, cannot add on casting speed changed callback!")
	}
	unit.OnCastSpeedChanged = append(unit.OnCastSpeedChanged, ocsc)
}

func (unit *Unit) AddOnMeleeAttackSpeedChanged(ocsc OnSpeedChanged) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, cannot add on melee attack speed changed callback!")
	}
	unit.OnMeleeAttackSpeedChanged = append(unit.OnMeleeAttackSpeedChanged, ocsc)
}

func (unit *Unit) AddOnTemporaryStatsChange(otsc OnTemporaryStatsChange) {
	if unit.Env != nil && unit.Env.IsFinalized() {
		panic("Already finalized, cannot add on temporary stats change callback!")
	}
	unit.OnTemporaryStatsChanges = append(unit.OnTemporaryStatsChanges, otsc)
}

func (unit *Unit) AddStatsDynamic(sim *Simulation, bonus stats.Stats) {
	if unit.Env == nil {
		panic("Environment not constructed.")
	} else if !unit.Env.IsFinalized() && !unit.Env.MeasuringStats {
		panic("Not finalized, use AddStats instead!")
	}

	unit.statsWithoutDeps.AddInplace(&bonus)

	if !unit.Env.MeasuringStats || unit.Env.State == Finalized {
		bonus = unit.ApplyStatDependencies(bonus)
	}

	if sim.Log != nil {
		unit.Log(sim, "Dynamic stat change: %s", bonus.FlatString())
	}

	unit.stats.AddInplace(&bonus)
	unit.processDynamicBonus(sim, bonus)
}

func (unit *Unit) AddStatDynamic(sim *Simulation, stat stats.Stat, amount float64) {
	bonus := stats.Stats{}
	bonus[stat] = amount
	unit.AddStatsDynamic(sim, bonus)
}

func (unit *Unit) processDynamicBonus(sim *Simulation, bonus stats.Stats) {
	if bonus[stats.MP5] != 0 || bonus[stats.Intellect] != 0 || bonus[stats.Spirit] != 0 {
		unit.UpdateManaRegenRates()
	}
	if bonus[stats.Mana] != 0 && unit.HasManaBar() {
		if unit.CurrentMana() > unit.MaxMana() {
			unit.currentMana = unit.MaxMana()
		}
	}
	if bonus[stats.HasteRating] != 0 {
		unit.updateAttackSpeed()
		unit.AutoAttacks.UpdateSwingTimers(sim)
		unit.runicPowerBar.updateRegenTimes(sim)
		unit.energyBar.processDynamicHasteRatingChange(sim)
		unit.focusBar.processDynamicHasteRatingChange(sim)
		unit.updateCastSpeed()
	}
	if bonus[stats.MasteryRating] != 0 {
		newMasteryRating := unit.stats[stats.MasteryRating]
		oldMasteryRating := newMasteryRating - bonus[stats.MasteryRating]

		for i := range unit.OnMasteryStatChanged {
			unit.OnMasteryStatChanged[i](sim, oldMasteryRating, newMasteryRating)
		}
	}

	for _, pet := range unit.DynamicStatsPets {
		pet.addOwnerStats(sim, bonus)
	}
}

func (unit *Unit) EnableDynamicStatDep(sim *Simulation, dep *stats.StatDependency) {
	if unit.StatDependencyManager.EnableDynamicStatDep(dep) {
		oldStats := unit.stats
		unit.stats = unit.ApplyStatDependencies(unit.statsWithoutDeps)
		unit.processDynamicBonus(sim, unit.stats.Subtract(oldStats))

		if sim.Log != nil {
			unit.Log(sim, "Dynamic dep enabled (%s): %s", dep.String(), unit.stats.Subtract(oldStats).FlatString())
		}
	}
}
func (unit *Unit) DisableDynamicStatDep(sim *Simulation, dep *stats.StatDependency) {
	if unit.StatDependencyManager.DisableDynamicStatDep(dep) {
		oldStats := unit.stats
		unit.stats = unit.ApplyStatDependencies(unit.statsWithoutDeps)
		unit.processDynamicBonus(sim, unit.stats.Subtract(oldStats))

		if sim.Log != nil {
			unit.Log(sim, "Dynamic dep disabled (%s): %s", dep.String(), unit.stats.Subtract(oldStats).FlatString())
		}
	}
}

func (unit *Unit) UpdateDynamicStatDep(sim *Simulation, dep *stats.StatDependency, newAmount float64) {
	dep.UpdateValue(newAmount)

	if unit.Env.IsFinalized() {
		oldStats := unit.stats
		unit.stats = unit.ApplyStatDependencies(unit.statsWithoutDeps)
		statsChange := unit.stats.Subtract(oldStats)
		unit.processDynamicBonus(sim, statsChange)

		if sim.Log != nil {
			unit.Log(sim, "Dynamic dep updated (%s): %s", dep.String(), statsChange.FlatString())
		}
	}
}

func (unit *Unit) EnableBuildPhaseStatDep(sim *Simulation, dep *stats.StatDependency) {
	if unit.Env.MeasuringStats && unit.Env.State != Finalized {
		unit.StatDependencyManager.EnableDynamicStatDep(dep)
	} else {
		unit.EnableDynamicStatDep(sim, dep)
	}
}

func (unit *Unit) DisableBuildPhaseStatDep(sim *Simulation, dep *stats.StatDependency) {
	if unit.Env.MeasuringStats && unit.Env.State != Finalized {
		unit.StatDependencyManager.DisableDynamicStatDep(dep)
	} else {
		unit.DisableDynamicStatDep(sim, dep)
	}
}

// Returns whether the indicates stat is currently modified by a temporary bonus.
func (unit *Unit) HasTemporaryBonusForStat(stat stats.Stat) bool {
	return unit.initialStats[stat] != unit.stats[stat]
}

// Returns if spell casting has any temporary increases active.
func (unit *Unit) HasTemporarySpellCastSpeedIncrease() bool {
	return unit.CastSpeed != unit.initialCastSpeed
}

// Returns if melee swings have any temporary increases active.
func (unit *Unit) HasTemporaryMeleeSwingSpeedIncrease() bool {
	return unit.TotalMeleeHasteMultiplier() != unit.initialMeleeSwingSpeed
}

// Returns if ranged swings have any temporary increases active.
func (unit *Unit) HasTemporaryRangedSwingSpeedIncrease() bool {
	return unit.TotalRangedHasteMultiplier() != unit.initialRangedSwingSpeed
}

func (unit *Unit) InitialCastSpeed() float64 {
	return unit.initialCastSpeed
}

func (unit *Unit) SpellGCD() time.Duration {
	return max(GCDMin, unit.ApplyCastSpeed(GCDDefault))
}

func (unit *Unit) TotalSpellHasteMultiplier() float64 {
	return unit.PseudoStats.CastSpeedMultiplier * (1 + unit.stats[stats.HasteRating]/(HasteRatingPerHastePercent*100))
}

func (unit *Unit) updateCastSpeed() {
	oldCastSpeed := unit.CastSpeed
	unit.CastSpeed = 1 / unit.TotalSpellHasteMultiplier()
	newCastSpeed := unit.CastSpeed

	for i := range unit.OnCastSpeedChanged {
		unit.OnCastSpeedChanged[i](oldCastSpeed, newCastSpeed)
	}
}
func (unit *Unit) MultiplyCastSpeed(amount float64) {
	unit.PseudoStats.CastSpeedMultiplier *= amount

	for _, pet := range unit.DynamicCastSpeedPets {
		pet.dynamicCastSpeedInheritance(amount)
	}

	unit.updateCastSpeed()
}

func (unit *Unit) ApplyCastSpeed(dur time.Duration) time.Duration {
	return time.Duration(float64(dur) * unit.CastSpeed)
}
func (unit *Unit) ApplyCastSpeedForSpell(dur time.Duration, spell *Spell) time.Duration {
	return time.Duration(float64(dur) * unit.CastSpeed * max(0, spell.CastTimeMultiplier))
}

func (unit *Unit) TotalMeleeHasteMultiplier() float64 {
	return unit.PseudoStats.AttackSpeedMultiplier * unit.PseudoStats.MeleeSpeedMultiplier * (1 + (unit.stats[stats.HasteRating] / (HasteRatingPerHastePercent * 100)))
}

// Returns the melee haste multiplier only including equip haste and real haste modifiers like lust
// Same value for ranged and melee
func (unit *Unit) TotalRealHasteMultiplier() float64 {
	return unit.PseudoStats.AttackSpeedMultiplier * (1 + (unit.stats[stats.HasteRating] / (HasteRatingPerHastePercent * 100)))
}

func (unit *Unit) Armor() float64 {
	return unit.PseudoStats.ArmorMultiplier * unit.stats[stats.Armor]
}

func (unit *Unit) BlockDamageReduction() float64 {
	return unit.PseudoStats.BlockDamageReduction
}

func (unit *Unit) TotalRangedHasteMultiplier() float64 {
	return unit.PseudoStats.AttackSpeedMultiplier * unit.PseudoStats.RangedSpeedMultiplier * (1 + (unit.stats[stats.HasteRating] / (HasteRatingPerHastePercent * 100)))
}

func (unit *Unit) updateMeleeAttackSpeed() {
	oldMeleeAttackSpeed := unit.meleeAttackSpeed
	unit.meleeAttackSpeed = unit.TotalMeleeHasteMultiplier()

	for i := range unit.OnMeleeAttackSpeedChanged {
		unit.OnMeleeAttackSpeedChanged[i](oldMeleeAttackSpeed, unit.meleeAttackSpeed)
	}
}

// MultiplyMeleeSpeed will alter the attack speed multiplier and change swing speed of all autoattack swings in progress.
func (unit *Unit) MultiplyMeleeSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.MeleeSpeedMultiplier *= amount

	unit.updateMeleeAttackSpeed()

	for _, pet := range unit.DynamicMeleeSpeedPets {
		pet.dynamicMeleeSpeedInheritance(amount)
	}
	unit.AutoAttacks.UpdateSwingTimers(sim)
}

func (unit *Unit) updateRangedAttackSpeed() {
	oldRangedAttackSpeed := unit.rangedAttackSpeed
	unit.rangedAttackSpeed = unit.TotalRangedHasteMultiplier()

	for i := range unit.OnRangedAttackSpeedChanged {
		unit.OnRangedAttackSpeedChanged[i](oldRangedAttackSpeed, unit.rangedAttackSpeed)
	}
}

func (unit *Unit) MultiplyRangedSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.RangedSpeedMultiplier *= amount
	unit.updateRangedAttackSpeed()
	unit.AutoAttacks.UpdateSwingTimers(sim)
}

func (unit *Unit) updateAttackSpeed() {
	oldMeleeAttackSpeed := unit.meleeAttackSpeed
	oldRangedAttackSpeed := unit.rangedAttackSpeed

	unit.meleeAttackSpeed = unit.TotalMeleeHasteMultiplier()
	unit.rangedAttackSpeed = unit.TotalRangedHasteMultiplier()

	for i := range unit.OnMeleeAttackSpeedChanged {
		unit.OnMeleeAttackSpeedChanged[i](oldMeleeAttackSpeed, unit.meleeAttackSpeed)
	}

	for i := range unit.OnRangedAttackSpeedChanged {
		unit.OnRangedAttackSpeedChanged[i](oldRangedAttackSpeed, unit.rangedAttackSpeed)
	}
}

// Helper for when true haste effects are multiplied for i.E. Bloodlust
// Seems to also always impact the regen rate
func (unit *Unit) MultiplyAttackSpeed(sim *Simulation, amount float64) {
	unit.PseudoStats.AttackSpeedMultiplier *= amount
	unit.MultiplyResourceRegenSpeed(sim, amount)
	unit.updateAttackSpeed()

	for _, pet := range unit.DynamicMeleeSpeedPets {
		pet.dynamicMeleeSpeedInheritance(amount)
	}

	unit.AutoAttacks.UpdateSwingTimers(sim)
}

// Helper for multiplying resource generation speed
func (unit *Unit) MultiplyResourceRegenSpeed(sim *Simulation, amount float64) {
	if unit.HasRunicPowerBar() {
		unit.MultiplyRuneRegenSpeed(sim, amount)
	} else if unit.HasFocusBar() {
		unit.MultiplyFocusRegenSpeed(sim, amount)
	} else if unit.HasEnergyBar() {
		unit.MultiplyEnergyRegenSpeed(sim, amount)
	}
}

func (unit *Unit) AddBonusRangedHitPercent(percentage float64) {
	unit.OnSpellRegistered(func(spell *Spell) {
		if spell.ProcMask.Matches(ProcMaskRanged) {
			spell.BonusHitPercent += percentage
		}
	})
}
func (unit *Unit) AddBonusRangedCritPercent(percentage float64) {
	unit.OnSpellRegistered(func(spell *Spell) {
		if spell.ProcMask.Matches(ProcMaskRanged) {
			spell.BonusCritPercent += percentage
		}
	})
}

func (unit *Unit) SetCurrentPowerBar(bar PowerBarType) {
	unit.currentPowerBar = bar
}

func (unit *Unit) GetCurrentPowerBar() PowerBarType {
	return unit.currentPowerBar
}

// Stat dependencies that apply both to players/pets (represented as Character
// structs) and to NPCs (represented as Target structs).
func (unit *Unit) addUniversalStatDependencies() {
	unit.AddStatDependency(stats.HitRating, stats.PhysicalHitPercent, 1/PhysicalHitRatingPerHitPercent)
	unit.AddStatDependency(stats.HitRating, stats.SpellHitPercent, 1/SpellHitRatingPerHitPercent)
	unit.AddStatDependency(stats.ExpertiseRating, stats.SpellHitPercent, 1/SpellHitRatingPerHitPercent)
	unit.AddStatDependency(stats.CritRating, stats.PhysicalCritPercent, 1/CritRatingPerCritPercent)
	unit.AddStatDependency(stats.CritRating, stats.SpellCritPercent, 1/CritRatingPerCritPercent)
}

func (unit *Unit) finalize() {
	if unit.Env.IsFinalized() {
		panic("Unit already finalized!")
	}

	// Make sure we don't accidentally set initial stats instead of stats.
	if !unit.initialStats.Equals(stats.Stats{}) {
		panic("Initial stats may not be set before finalized: " + unit.initialStats.String())
	}

	if unit.ReactionTime == 0 {
		panic("Unset unit.ReactionTime")
	}

	unit.defaultTarget = unit.CurrentTarget
	unit.applyParryHaste()
	unit.updateCastSpeed()
	unit.updateAttackSpeed()
	unit.initMovement()

	// All stats added up to this point are part of the 'initial' stats.
	unit.initialStatsWithoutDeps = unit.stats
	unit.initialPseudoStats = unit.PseudoStats
	unit.initialCastSpeed = unit.CastSpeed
	unit.initialMeleeSwingSpeed = unit.TotalMeleeHasteMultiplier()
	unit.initialRangedSwingSpeed = unit.TotalRangedHasteMultiplier()

	unit.StatDependencyManager.FinalizeStatDeps()
	unit.initialStats = unit.ApplyStatDependencies(unit.initialStatsWithoutDeps)
	unit.statsWithoutDeps = unit.initialStatsWithoutDeps
	unit.stats = unit.initialStats

	unit.AutoAttacks.finalize()

	if unit.GetSpellPowerValue == nil {
		unit.GetSpellPowerValue = unit.getSpellPowerValueImpl
	}

	if unit.GetAttackPowerValue == nil {
		unit.GetAttackPowerValue = unit.getAttackPowerValueImpl
	}

	for _, spell := range unit.Spellbook {
		spell.finalize()
	}
}

func (unit *Unit) reset(sim *Simulation, _ Agent) {
	unit.enabled = true
	unit.resetCDs(sim)
	unit.Hardcast.Expires = startingCDTime
	unit.ChanneledDot = nil
	unit.QueuedSpell = nil
	unit.DistanceFromTarget = unit.StartDistanceFromTarget
	unit.Metrics.reset()
	unit.ResetStatDeps()
	unit.statsWithoutDeps = unit.initialStatsWithoutDeps
	unit.stats = unit.initialStats
	unit.PseudoStats = unit.initialPseudoStats
	unit.auraTracker.reset(sim)
	for _, spell := range unit.Spellbook {
		spell.reset(sim)
	}

	unit.manaBar.reset()
	unit.focusBar.reset(sim)
	unit.healthBar.reset(sim)
	unit.UpdateManaRegenRates()

	unit.energyBar.reset(sim)
	unit.rageBar.reset(sim)
	unit.runicPowerBar.reset(sim)

	if unit.secondaryResourceBar != nil {
		unit.secondaryResourceBar.Reset(sim)
	}

	unit.AutoAttacks.reset(sim)

	if unit.Rotation != nil {
		unit.Rotation.reset(sim)
	}

	unit.DynamicStatsPets = unit.DynamicStatsPets[:0]
	unit.DynamicMeleeSpeedPets = unit.DynamicMeleeSpeedPets[:0]

	if unit.Type != PetUnit {
		sim.addTracker(&unit.auraTracker)
	}
}

func (unit *Unit) startPull(sim *Simulation) {
	unit.AutoAttacks.startPull(sim)

	if unit.Type == PlayerUnit {
		unit.SetGCDTimer(sim, max(0, unit.GCD.ReadyAt()))
	}
}

func (unit *Unit) doneIteration(sim *Simulation) {
	unit.Hardcast = Hardcast{}

	unit.manaBar.doneIteration(sim)
	unit.rageBar.doneIteration()

	unit.auraTracker.doneIteration(sim)
	for _, spell := range unit.Spellbook {
		spell.doneIteration()
	}
}

func (unit *Unit) GetSpellsMatchingSchool(school SpellSchool) []*Spell {
	var spells []*Spell
	for _, spell := range unit.Spellbook {
		if spell.SpellSchool.Matches(school) {
			spells = append(spells, spell)
		}
	}
	return spells
}

func (unit *Unit) GetUnit(ref *proto.UnitReference) *Unit {
	return unit.Env.GetUnit(ref, unit)
}

func (unit *Unit) GetMetadata() *proto.UnitMetadata {
	metadata := &proto.UnitMetadata{
		Name: unit.Label,
	}

	metadata.Spells = MapSlice(unit.Spellbook, func(spell *Spell) *proto.SpellStats {
		return &proto.SpellStats{
			Id: spell.ActionID.ToProto(),

			IsCastable:      spell.Flags.Matches(SpellFlagAPL),
			IsChanneled:     spell.Flags.Matches(SpellFlagChanneled),
			IsMajorCooldown: spell.Flags.Matches(SpellFlagMCD),
			HasDot:          spell.dots != nil || spell.aoeDot != nil || (spell.RelatedDotSpell != nil && (spell.RelatedDotSpell.dots != nil || spell.RelatedDotSpell.aoeDot != nil)),
			HasShield:       spell.shields != nil || spell.selfShield != nil,
			PrepullOnly:     spell.Flags.Matches(SpellFlagPrepullOnly),
			EncounterOnly:   spell.Flags.Matches(SpellFlagEncounterOnly),
			HasCastTime:     spell.DefaultCast.CastTime > 0,
			IsFriendly:      spell.Flags.Matches(SpellFlagHelpful),
		}
	})

	aplAuras := FilterSlice(unit.auras, func(aura *Aura) bool {
		return !aura.ActionID.IsEmptyAction()
	})
	metadata.Auras = MapSlice(aplAuras, func(aura *Aura) *proto.AuraStats {
		return &proto.AuraStats{
			Id:                 aura.ActionID.ToProto(),
			MaxStacks:          aura.MaxStacks,
			HasIcd:             aura.Icd != nil,
			HasExclusiveEffect: len(aura.ExclusiveEffects) > 0,
		}
	})

	return metadata
}

func (unit *Unit) ExecuteCustomRotation(sim *Simulation) {
	panic("Unimplemented ExecuteCustomRotation")
}

func (unit *Unit) GetTotalDodgeChanceAsDefender(atkTable *AttackTable) float64 {
	chance := unit.PseudoStats.BaseDodgeChance +
		atkTable.BaseDodgeChance +
		unit.GetDiminishedDodgeChance()
	return math.Max(chance, 0.0)
}

func (unit *Unit) GetTotalParryChanceAsDefender(atkTable *AttackTable) float64 {
	chance := unit.PseudoStats.BaseParryChance +
		atkTable.BaseParryChance +
		unit.GetDiminishedParryChance()
	return math.Max(chance, 0.0)
}

func (unit *Unit) GetTotalChanceToBeMissedAsDefender(atkTable *AttackTable) float64 {
	chance := atkTable.BaseMissChance +
		unit.PseudoStats.ReducedPhysicalHitTakenChance
	return math.Max(chance, 0.0)
}

func (unit *Unit) GetTotalBlockChanceAsDefender(atkTable *AttackTable) float64 {
	chance := unit.PseudoStats.BaseBlockChance +
		atkTable.BaseBlockChance +
		unit.GetDiminishedBlockChance()
	return math.Max(chance, 0.0)
}

func (unit *Unit) GetTotalAvoidanceChance(atkTable *AttackTable) float64 {
	miss := unit.GetTotalChanceToBeMissedAsDefender(atkTable)
	dodge := unit.GetTotalDodgeChanceAsDefender(atkTable)
	parry := unit.GetTotalParryChanceAsDefender(atkTable)
	block := unit.GetTotalBlockChanceAsDefender(atkTable)
	return miss + dodge + parry + block
}
