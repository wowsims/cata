package core

import (
	"slices"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
)

type EnvironmentState int

const (
	Created EnvironmentState = iota
	Constructed
	Initialized
	Finalized
)

// Callback for doing something after finalization.
type PostFinalizeEffect func()

// Callback for doing something on prepull.
type PrepullAction struct {
	DoAt   time.Duration
	Action func(*Simulation)
}

type Environment struct {
	State EnvironmentState

	// Whether stats are currently being measured. Used to disable some validation
	// checks which are otherwise helpful.
	MeasuringStats bool

	Raid      *Raid
	Encounter Encounter
	AllUnits  []*Unit

	BaseDuration      time.Duration // base duration
	DurationVariation time.Duration // variation per duration

	// Effects to invoke when the Env is finalized.
	preFinalizeEffects  []PostFinalizeEffect
	postFinalizeEffects []PostFinalizeEffect

	prepullActions []PrepullAction
}

func NewEnvironment(raidProto *proto.Raid, encounterProto *proto.Encounter, runFakePrepull bool) (*Environment, *proto.RaidStats, *proto.EncounterStats) {
	env := &Environment{
		State: Created,
	}

	env.construct(raidProto, encounterProto)
	raidStats := env.initialize(raidProto, encounterProto)
	env.finalize(raidProto, encounterProto, raidStats, runFakePrepull)

	encounterStats := &proto.EncounterStats{}
	for _, target := range env.Encounter.AllTargets {
		encounterStats.Targets = append(encounterStats.Targets, &proto.TargetStats{
			Metadata: target.GetMetadata(),
		})
	}

	return env, raidStats, encounterStats
}

// The construction phase.
func (env *Environment) construct(raidProto *proto.Raid, encounterProto *proto.Encounter) {
	env.Encounter = NewEncounter(encounterProto)
	env.BaseDuration = env.Encounter.Duration
	env.DurationVariation = env.Encounter.DurationVariation
	env.Raid = NewRaid(raidProto)

	env.Raid.updatePlayersAndPets()

	env.AllUnits = append(env.Encounter.AllTargetUnits, env.Raid.AllUnits...)

	for unitIndex, unit := range env.AllUnits {
		unit.Env = env
		unit.UnitIndex = int32(unitIndex)
	}

	for _, unit := range env.Raid.AllUnits {
		unit.CurrentTarget = env.Encounter.ActiveTargetUnits[0]
	}

	// Apply extra debuffs from raid.
	if raidProto.Debuffs != nil && len(env.Encounter.AllTargetUnits) > 0 {
		for targetIdx, targetUnit := range env.Encounter.AllTargetUnits {
			applyDebuffEffects(targetUnit, targetIdx, raidProto.Debuffs, raidProto)
		}
	}

	tankTargetSet := map[*Unit]bool{}
	// Assign target-of-target using Tanks field.
	for _, target := range env.Encounter.AllTargets {
		if target.Index < int32(len(encounterProto.Targets)) {
			targetProto := encounterProto.Targets[target.Index]
			env.setupTankTarget(target, targetProto.TankIndex, raidProto.Tanks, true, tankTargetSet)

			if targetProto.SecondTankIndex != targetProto.TankIndex {
				env.setupTankTarget(target, targetProto.SecondTankIndex, raidProto.Tanks, false, tankTargetSet)
			}
		}
	}

	env.State = Constructed
}

// The initialization phase.
func (env *Environment) initialize(raidProto *proto.Raid, encounterProto *proto.Encounter) *proto.RaidStats {
	for _, target := range env.Encounter.AllTargets {
		if target.Index < int32(len(encounterProto.Targets)) {
			target.initialize(encounterProto.Targets[target.Index])
		} else {
			target.initialize(nil)
		}
	}

	for _, party := range env.Raid.Parties {
		for _, playerOrPet := range party.PlayersAndPets {
			playerOrPet.GetCharacter().initialize(playerOrPet)
		}
	}

	raidStats := env.Raid.applyCharacterEffects(raidProto)

	for _, party := range env.Raid.Parties {
		for _, playerOrPet := range party.PlayersAndPets {
			playerOrPet.Initialize()
		}
	}

	env.State = Initialized
	return raidStats
}

// The finalization phase.
func (env *Environment) finalize(raidProto *proto.Raid, _ *proto.Encounter, raidStats *proto.RaidStats, runFakePrepull bool) {
	for _, finalizeEffect := range env.preFinalizeEffects {
		finalizeEffect()
	}
	env.preFinalizeEffects = nil

	for _, target := range env.Encounter.AllTargets {
		target.finalize()
		if target.AI != nil {
			target.Rotation = target.newCustomRotation()
		}
	}

	for _, party := range env.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			character.Finalize()
			for _, pet := range character.Pets {
				pet.Finalize()
				pet.Rotation = pet.newCustomRotation()
			}
		}
	}

	for partyIdx, party := range env.Raid.Parties {
		partyProto := raidProto.Parties[partyIdx]
		for playerIdx, player := range party.Players {
			if playerIdx >= len(partyProto.Players) {
				// This happens for target dummies.
				continue
			}
			playerProto := partyProto.Players[playerIdx]
			char := player.GetCharacter()
			char.Rotation = char.newAPLRotation(playerProto.Rotation)
		}
	}

	// Use a traditional for loop here to accomodate callback chains that
	// queue up additional delayed evaluations.
	for i := 0; i < len(env.postFinalizeEffects); i++ {
		env.postFinalizeEffects[i]()
	}
	env.postFinalizeEffects = nil

	slices.SortStableFunc(env.prepullActions, func(a1, a2 PrepullAction) int {
		return int(a1.DoAt - a2.DoAt)
	})

	env.setupAttackTables()

	env.State = Finalized

	for partyIdx, party := range env.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			character.FillPlayerStats(raidStats.Parties[partyIdx].Players[character.PartyIndex])
		}
	}

	if runFakePrepull {
		// Runs prepull only, for a single iteration. This lets us detect misconfigured
		// prepull spells (e.g. GCD not available) in APL.
		sim := newSimWithEnv(env, &proto.SimOptions{
			Iterations: 1,
		}, simsignals.CreateSignals())
		sim.reset()
		sim.PrePull()
		sim.Cleanup()
	}
}

func (env *Environment) setupTankTarget(npc *Target, tankIndex int32, tankList []*proto.UnitReference, isFirstTank bool, tankHasTarget map[*Unit]bool) {
	if (tankIndex < 0) || (tankIndex >= int32(len(tankList))) {
		return
	}

	tankRef := tankList[tankIndex]

	if tankRef == nil {
		return
	}

	tankUnit := env.GetUnit(tankRef, nil)

	if tankUnit == nil {
		return
	}

	if isFirstTank {
		npc.CurrentTarget = tankUnit
	} else {
		npc.SecondaryTarget = tankUnit
	}

	// Set the tank's target to the first unit tanked.
	if !tankHasTarget[tankUnit] {
		tankUnit.CurrentTarget = &npc.Unit
		tankHasTarget[tankUnit] = true
	}
}

func (env *Environment) setupAttackTables() {
	raidUnits := env.Raid.AllUnits
	if len(raidUnits) == 0 {
		return
	}

	for _, attacker := range env.AllUnits {
		attacker.AttackTables = make([]*AttackTable, len(env.AllUnits))
		for idx, defender := range env.AllUnits {
			attacker.AttackTables[idx] = NewAttackTable(attacker, defender)
		}
	}
}

func (env *Environment) IsFinalized() bool {
	return env.State >= Finalized
}

func (env *Environment) reset(sim *Simulation) {
	// Reset primary targets damage taken for tracking health fights.
	env.Encounter.DamageTaken = 0

	// Targets need to be reset before the raid, so that players can check for
	// the presence of permanent target auras in their Reset handlers.
	for _, target := range env.Encounter.AllTargets {
		target.Reset(sim)
	}

	env.Raid.reset(sim)
}

// The maximum possible duration for any iteration.
func (env *Environment) GetMaxDuration() time.Duration {
	return env.BaseDuration + env.DurationVariation
}

func (env *Environment) ActiveTargetCount() int32 {
	return int32(len(env.Encounter.ActiveTargets))
}
func (env *Environment) TotalTargetCount() int32 {
	return int32(len(env.Encounter.AllTargets))
}

func (env *Environment) GetActiveTargets() []*Target {
	return env.Encounter.ActiveTargets
}

func (env *Environment) GetTargetByIndex(index int32) *Target {
	return env.Encounter.AllTargets[index]
}
func (env *Environment) GetTargetUnitByIndex(index int32) *Unit {
	return env.Encounter.AllTargetUnits[index]
}
func (env *Environment) NextActiveTarget(target *Unit) *Target {
	return env.Encounter.AllTargets[target.Index].NextActiveTarget()
}
func (env *Environment) NextActiveTargetUnit(target *Unit) *Unit {
	return &env.NextActiveTarget(target).Unit
}
func (env *Environment) GetAgentFromUnit(unit *Unit) Agent {
	raidAgent := env.Raid.GetPlayerFromUnit(unit)
	if raidAgent != nil {
		return raidAgent
	}

	for _, target := range env.Encounter.AllTargets {
		if unit == &target.Unit {
			return target
		}
	}

	return nil
}

func (env *Environment) GetUnit(ref *proto.UnitReference, contextUnit *Unit) *Unit {
	if ref == nil {
		return nil
	}

	switch ref.Type {
	case proto.UnitReference_Player:
		raidIndex := ref.Index
		partyIndex := int(raidIndex / 5)
		if partyIndex < 0 || partyIndex >= len(env.Raid.Parties) {
			return nil
		}

		party := env.Raid.Parties[partyIndex]
		for _, player := range party.Players {
			if player.GetCharacter().Index == raidIndex {
				return &player.GetCharacter().Unit
			}
		}
	case proto.UnitReference_Pet:
		ownerAgent := env.Raid.GetPlayerFromUnit(env.GetUnit(ref.Owner, contextUnit))
		if ownerAgent == nil {
			return nil
		}
		pets := ownerAgent.GetCharacter().PetAgents
		if int(ref.Index) < len(pets) {
			return &pets[ref.Index].GetCharacter().Unit
		} else {
			return nil
		}
	case proto.UnitReference_Target:
		if int(ref.Index) < len(env.Encounter.AllTargetUnits) {
			return env.Encounter.AllTargetUnits[ref.Index]
		} else {
			return nil
		}
	case proto.UnitReference_Self:
		return contextUnit
	case proto.UnitReference_CurrentTarget:
		if contextUnit == nil {
			return nil
		}
		return contextUnit.CurrentTarget
	}

	return nil
}

// Registers a callback to this Character which will be invoked BEFORE all Units
// are finalized, but after they are all initialized and have other effects applied.
func (env *Environment) RegisterPreFinalizeEffect(preFinalizeEffect PostFinalizeEffect) {
	if env.IsFinalized() {
		panic("Pre-Finalize effects may not be added once finalized!")
	}

	env.preFinalizeEffects = append(env.preFinalizeEffects, preFinalizeEffect)
}

// Registers a callback to this Character which will be invoked AFTER all Units
// are finalized.
func (env *Environment) RegisterPostFinalizeEffect(postFinalizeEffect PostFinalizeEffect) {
	if env.IsFinalized() {
		panic("Post-Finalize effects may not be added once finalized!")
	}

	env.postFinalizeEffects = append(env.postFinalizeEffects, postFinalizeEffect)
}

// Registers a callback to this Unit which will be invoked on the prepull at the specified
// negative time.
func (unit *Unit) RegisterPrepullAction(doAt time.Duration, action func(*Simulation)) {
	env := unit.Env
	if env.IsFinalized() {
		panic("Prepull actions may not be added once finalized!")
	}
	if doAt > 0 {
		panic("Prepull DoAt must not be positive!")
	}

	env.prepullActions = append(env.prepullActions, PrepullAction{
		DoAt:   doAt,
		Action: action,
	})
}

func (env *Environment) PrepullStartTime() time.Duration {
	if !env.IsFinalized() {
		panic("Env not yet finalized")
	}

	if len(env.prepullActions) == 0 {
		return 0
	} else {
		return env.prepullActions[0].DoAt
	}
}
