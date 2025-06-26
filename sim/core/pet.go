package core

import (
	"fmt"
	"slices"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// Extension of Agent interface, for Pets.
type PetAgent interface {
	Agent

	// The Pet controlled by this PetAgent.
	GetPet() *Pet
}

type OnPetEnable func(sim *Simulation)
type OnPetDisable func(sim *Simulation)

type PetStatInheritance func(ownerStats stats.Stats) stats.Stats
type PetSpeedInheritance func(amount float64)

type PetConfig struct {
	Name      string
	Owner     *Character
	BaseStats stats.Stats
	// Hit and Expertise are always inherited by combining the owners physical hit and expertise, then halving it
	// For casters this will automatically give spell hit cap at 7.5% physical hit and exp
	NonHitExpStatInheritance        PetStatInheritance
	EnabledOnStart                  bool
	IsGuardian                      bool
	HasDynamicMeleeSpeedInheritance bool
	HasDynamicCastSpeedInheritance  bool
	HasResourceRegenInheritance     bool
}

// Pet is an extension of Character, for any entity created by a player that can
// take actions on its own.
type Pet struct {
	Character

	Owner *Character

	isGuardian     bool
	enabledOnStart bool

	OnPetEnable  OnPetEnable
	OnPetDisable OnPetDisable

	// Calculates inherited stats based on owner stats or stat changes.
	statInheritance        PetStatInheritance
	dynamicStatInheritance PetStatInheritance
	inheritedStats         stats.Stats

	// In MoP pets inherit their owners melee speed and cast speed
	// rather than having auras such as Heroism being applied to them.
	dynamicMeleeSpeedInheritance PetSpeedInheritance
	dynamicCastSpeedInheritance  PetSpeedInheritance

	// If true the pet will automatically inherit the owner's melee speed
	hasDynamicMeleeSpeedInheritance bool
	// If true the pet will automatically inherit the owner's cast speed
	hasDynamicCastSpeedInheritance bool
	// If true the pet will automatically inherit the owner's regen speed multiplier
	hasResourceRegenInheritance bool

	isReset bool

	// Some pets expire after a certain duration. This is the pending action that disables
	// the pet on expiration.
	timeoutAction *PendingAction

	// Examples:
	// DK Raise Dead is doing its whole RP thing by climbing out of the ground before attacking.
	// Monk clones Rush towards targets before attacking.
	startAttackDelay time.Duration
}

func NewPet(config PetConfig) Pet {
	pet := Pet{
		Character: Character{
			Unit: Unit{
				Type:        PetUnit,
				Index:       config.Owner.Party.Raid.getNextPetIndex(),
				Label:       fmt.Sprintf("%s - %s", config.Owner.Label, config.Name),
				Level:       CharacterLevel,
				PseudoStats: stats.NewPseudoStats(),
				auraTracker: newAuraTracker(),
				Metrics:     NewUnitMetrics(),

				StatDependencyManager: stats.NewStatDependencyManager(),

				ReactionTime: config.Owner.ReactionTime,

				StartDistanceFromTarget: MaxMeleeRange, // TODO: Match to owner and add movement logic to pet rotations
			},
			Name:       config.Name,
			Party:      config.Owner.Party,
			PartyIndex: config.Owner.PartyIndex,
			baseStats:  config.BaseStats,
		},
		Owner:                           config.Owner,
		statInheritance:                 makeStatInheritanceFunc(config.NonHitExpStatInheritance),
		hasDynamicMeleeSpeedInheritance: config.HasDynamicMeleeSpeedInheritance,
		hasDynamicCastSpeedInheritance:  config.HasDynamicCastSpeedInheritance,
		hasResourceRegenInheritance:     config.HasResourceRegenInheritance,
		enabledOnStart:                  config.EnabledOnStart,
		isGuardian:                      config.IsGuardian,
	}

	pet.GCD = pet.NewTimer()
	pet.RotationTimer = pet.NewTimer()

	pet.AddStats(config.BaseStats)
	pet.addUniversalStatDependencies()
	pet.PseudoStats.InFrontOfTarget = config.Owner.PseudoStats.InFrontOfTarget
	return pet
}

func (pet *Pet) Initialize() {
	if pet.hasResourceRegenInheritance {
		pet.enableResourceRegenInheritance()
	}
}

func makeStatInheritanceFunc(nonHitExpStatInheritance PetStatInheritance) PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		inheritedStats := nonHitExpStatInheritance(ownerStats)

		hitRating := ownerStats[stats.HitRating]
		expertiseRating := ownerStats[stats.ExpertiseRating]
		combined := (hitRating + expertiseRating) * 0.5

		inheritedStats[stats.HitRating] = combined
		inheritedStats[stats.ExpertiseRating] = combined

		return inheritedStats
	}
}

// Updates the stats for this pet in response to a stat change on the owner.
// addedStats is the amount of stats added to the owner (will be negative if the
// owner lost stats).
func (pet *Pet) addOwnerStats(sim *Simulation, addedStats stats.Stats) {
	inheritedChange := pet.dynamicStatInheritance(addedStats)

	pet.inheritedStats.AddInplace(&inheritedChange)
	pet.AddStatsDynamic(sim, inheritedChange)
}

func (pet *Pet) reset(sim *Simulation, agent PetAgent) {
	if pet.isReset {
		return
	}
	pet.isReset = true

	pet.Character.reset(sim, agent)

	pet.CancelGCDTimer(sim)
	pet.AutoAttacks.CancelAutoSwing(sim)

	pet.enabled = false
	if pet.enabledOnStart {
		pet.Enable(sim, agent)
	}
}
func (pet *Pet) doneIteration(sim *Simulation) {
	pet.Character.doneIteration(sim)
	pet.Disable(sim)
	pet.isReset = false
}

func (pet *Pet) IsGuardian() bool {
	return pet.isGuardian
}

// petAgent should be the PetAgent which embeds this Pet.
func (pet *Pet) Enable(sim *Simulation, petAgent PetAgent) {
	if pet.enabled {
		if sim.Log != nil {
			pet.Log(sim, "Pet already summoned")
		}
		return
	}

	// In case of Pre-pull guardian summoning we need to reset
	// TODO: Check if this has side effects
	if !pet.isReset {
		pet.reset(sim, petAgent)
	}

	pet.inheritedStats = pet.statInheritance(pet.Owner.GetStats())
	pet.AddStatsDynamic(sim, pet.inheritedStats)
	pet.Owner.DynamicStatsPets = append(pet.Owner.DynamicStatsPets, pet)
	pet.dynamicStatInheritance = pet.statInheritance

	//reset current mana after applying stats
	pet.manaBar.reset()

	//reset current health after applying stats
	pet.healthBar.reset(sim)

	// Call onEnable callbacks before enabling auto swing
	// to not have to reorder PAs multiple times
	pet.enabled = true

	if pet.OnPetEnable != nil {
		pet.OnPetEnable(sim)
	}

	if pet.hasDynamicMeleeSpeedInheritance {
		pet.enableDynamicMeleeSpeed(func(amount float64) {
			pet.MultiplyMeleeSpeed(sim, amount)
		})
	}

	if pet.hasDynamicCastSpeedInheritance {
		pet.enableDynamicCastSpeed(func(amount float64) {
			pet.MultiplyCastSpeed(sim, amount)
		})
	}

	pet.SetGCDTimer(sim, max(0, sim.CurrentTime+pet.startAttackDelay, sim.CurrentTime))
	pet.AutoAttacks.EnableAutoSwing(sim)

	if (sim.CurrentTime < 0) || (pet.startAttackDelay > 0) {
		pet.AutoAttacks.StopMeleeUntil(sim, max(0, sim.CurrentTime+pet.startAttackDelay)-pet.AutoAttacks.MainhandSwingSpeed())
	}

	if sim.Log != nil {
		pet.Log(sim, "Pet stats: %s", pet.GetStats().FlatString())
		pet.Log(sim, "Pet inherited stats: %s", pet.ApplyStatDependencies(pet.inheritedStats).FlatString())
		pet.Log(sim, "Pet summoned")
	}

	sim.addTracker(&pet.auraTracker)

	if pet.HasFocusBar() {
		// make sure to reset it to refresh focus
		pet.focusBar.reset(sim)
		pet.focusBar.enable(sim, sim.CurrentTime)
		pet.focusBar.focusRegenMultiplier *= pet.Owner.PseudoStats.AttackSpeedMultiplier
	}

	if pet.HasEnergyBar() {
		// make sure to reset it to refresh energy
		pet.energyBar.reset(sim)
		pet.energyBar.enable(sim, sim.CurrentTime)
		pet.energyBar.energyRegenMultiplier *= pet.Owner.PseudoStats.AttackSpeedMultiplier
	}
}

func (pet *Pet) EnableWithStartAttackDelay(sim *Simulation, petAgent PetAgent, startAttackDelay time.Duration) {
	pet.startAttackDelay = startAttackDelay
	pet.Enable(sim, petAgent)
}

// Helper for enabling a pet that will expire after a certain duration.
func (pet *Pet) EnableWithTimeout(sim *Simulation, petAgent PetAgent, petDuration time.Duration) {
	pet.Enable(sim, petAgent)
	pet.SetTimeoutAction(sim, petDuration)
}

func (pet *Pet) SetTimeoutAction(sim *Simulation, duration time.Duration) {
	pet.timeoutAction = sim.GetConsumedPendingActionFromPool()
	pet.timeoutAction.NextActionAt = sim.CurrentTime + duration
	pet.timeoutAction.OnAction = pet.Disable
	sim.AddPendingAction(pet.timeoutAction)
}

func (pet *Pet) SetStartAttackDelay(startAttackDelay time.Duration) {
	pet.startAttackDelay = startAttackDelay
}

// Enables and possibly updates how the pet inherits its owner's stats.
func (pet *Pet) EnableDynamicStats(inheritance PetStatInheritance) {
	if !slices.Contains(pet.Owner.DynamicStatsPets, pet) {
		pet.Owner.DynamicStatsPets = append(pet.Owner.DynamicStatsPets, pet)
	}
	pet.dynamicStatInheritance = inheritance
}

// Enables and possibly updates how the pet inherits its owner's melee speed.
func (pet *Pet) EnableDynamicMeleeSpeed(inheritance PetSpeedInheritance) {
	if pet.hasDynamicMeleeSpeedInheritance {
		panic("To use custom EnableDynamicMeleeSpeed remove hasDynamicMeleeSpeedInheritance from the Pet constructor")
	}
	pet.enableDynamicMeleeSpeed(inheritance)
}

func (pet *Pet) enableDynamicMeleeSpeed(inheritance PetSpeedInheritance) {
	if !slices.Contains(pet.Owner.DynamicMeleeSpeedPets, pet) {
		pet.Owner.DynamicMeleeSpeedPets = append(pet.Owner.DynamicMeleeSpeedPets, pet)
		inheritance(pet.Owner.PseudoStats.MeleeSpeedMultiplier)
		inheritance(pet.Owner.PseudoStats.AttackSpeedMultiplier)
	}
	pet.dynamicMeleeSpeedInheritance = inheritance
}

// Enables and possibly updates how the pet inherits its owner's cast speed.
func (pet *Pet) EnableDynamicCastSpeed(inheritance PetSpeedInheritance) {
	if pet.hasDynamicCastSpeedInheritance {
		panic("To use custom EnableDynamicCastSpeed remove hasDynamicCastSpeedInheritance from the Pet constructor")
	}
	pet.enableDynamicCastSpeed(inheritance)
}

func (pet *Pet) enableDynamicCastSpeed(inheritance PetSpeedInheritance) {
	if !slices.Contains(pet.Owner.DynamicCastSpeedPets, pet) {
		pet.Owner.DynamicCastSpeedPets = append(pet.Owner.DynamicCastSpeedPets, pet)
		inheritance(pet.Owner.PseudoStats.CastSpeedMultiplier)
	}
	pet.dynamicCastSpeedInheritance = inheritance
}

func (pet *Pet) enableResourceRegenInheritance() {
	if !slices.Contains(pet.Owner.RegenInheritancePets, pet) {
		pet.Owner.RegenInheritancePets = append(pet.Owner.RegenInheritancePets, pet)
	}
}

func (pet *Pet) Disable(sim *Simulation) {
	if !pet.enabled {
		if sim.Log != nil {
			pet.Log(sim, "No pet summoned")
		}
		return
	}

	pet.AddStatsDynamic(sim, pet.inheritedStats.Invert())
	pet.inheritedStats = stats.Stats{}

	if pet.dynamicStatInheritance != nil {
		if idx := slices.Index(pet.Owner.DynamicStatsPets, pet); idx != -1 {
			pet.Owner.DynamicStatsPets = removeBySwappingToBack(pet.Owner.DynamicStatsPets, idx)
		}
		pet.dynamicStatInheritance = nil
	}

	if pet.dynamicMeleeSpeedInheritance != nil {
		if idx := slices.Index(pet.Owner.DynamicMeleeSpeedPets, pet); idx != -1 {
			pet.Owner.DynamicMeleeSpeedPets = removeBySwappingToBack(pet.Owner.DynamicMeleeSpeedPets, idx)
		}

		// reset melee speed inheritance here so pets that get enabled later to not keep it
		pet.dynamicMeleeSpeedInheritance(1 / pet.Owner.PseudoStats.MeleeSpeedMultiplier)
		pet.dynamicMeleeSpeedInheritance(1 / pet.Owner.PseudoStats.AttackSpeedMultiplier)
		pet.dynamicMeleeSpeedInheritance = nil
	}

	if pet.dynamicCastSpeedInheritance != nil {
		if idx := slices.Index(pet.Owner.DynamicCastSpeedPets, pet); idx != -1 {
			pet.Owner.DynamicCastSpeedPets = removeBySwappingToBack(pet.Owner.DynamicCastSpeedPets, idx)
		}

		// reset cast speed inheritance here so pets that get enabled later to not keep it
		pet.dynamicCastSpeedInheritance(1 / pet.Owner.PseudoStats.CastSpeedMultiplier)
		pet.dynamicCastSpeedInheritance = nil
	}

	pet.CancelGCDTimer(sim)
	pet.focusBar.disable(sim)
	pet.energyBar.disable(sim)
	pet.AutoAttacks.CancelAutoSwing(sim)
	pet.enabled = false

	// If a pet is immediately re-summoned it might try to use GCD, so we need to clear it.
	pet.Hardcast = Hardcast{}

	if (pet.timeoutAction != nil) && !pet.timeoutAction.consumed {
		pet.timeoutAction.Cancel(sim)
	}

	if pet.OnPetDisable != nil {
		pet.OnPetDisable(sim)
	}

	pet.auraTracker.expireAll(sim)

	sim.removeTracker(&pet.auraTracker)

	if sim.Log != nil {
		pet.Log(sim, "Pet dismissed")
		pet.Log(sim, pet.GetStats().FlatString())
	}
}

func (pet *Pet) ChangeStatInheritance(nonHitExpStatInheritance PetStatInheritance) {
	pet.statInheritance = makeStatInheritanceFunc(nonHitExpStatInheritance)
}

func (pet *Pet) GetInheritedStats() stats.Stats {
	return pet.inheritedStats
}

func (pet *Pet) DisableOnStart() {
	pet.enabledOnStart = false
}

// Default implementations for some Agent functions which most Pets don't need.
func (pet *Pet) GetCharacter() *Character {
	return &pet.Character
}
func (pet *Pet) AddRaidBuffs(_ *proto.RaidBuffs)   {}
func (pet *Pet) AddPartyBuffs(_ *proto.PartyBuffs) {}
func (pet *Pet) ApplyTalents()                     {}
func (pet *Pet) OnGCDReady(_ *Simulation)          {}

func (env *Environment) triggerDelayedPetInheritance(sim *Simulation, dynamicPets []*Pet, inheritanceFunc func(*Simulation, *Pet)) {
	for _, pet := range dynamicPets {
		if !pet.IsEnabled() {
			continue
		}

		numHeartbeats := (sim.CurrentTime - env.heartbeatOffset) / PetUpdateInterval
		nextHeartbeat := PetUpdateInterval*(numHeartbeats+1) + env.heartbeatOffset

		pa := sim.GetConsumedPendingActionFromPool()
		pa.NextActionAt = nextHeartbeat
		pa.Priority = ActionPriorityDOT

		pa.OnAction = func(sim *Simulation) {
			if pet.enabled {
				inheritanceFunc(sim, pet)
			}
		}

		sim.AddPendingAction(pa)
	}
}
