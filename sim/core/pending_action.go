package core

import (
	"slices"
	"time"
)

type ActionPriority int32

const (
	ActionPriorityLow ActionPriority = -1
	ActionPriorityGCD ActionPriority = 0

	// Higher than GCD because regen can cause GCD actions (if we were waiting
	// for mana).
	ActionPriorityRegen ActionPriority = 1

	// Autos can cause regen (JoW, rage, energy procs, etc) so they should be
	// higher prio so that we never go backwards in the priority order.
	ActionPriorityAuto ActionPriority = 2

	// DOTs need to be higher than anything else so that dots can properly expire before we take other actions.
	ActionPriorityDOT ActionPriority = 3

	ActionPriorityPrePull ActionPriority = 10
)

type PendingAction struct {
	NextActionAt time.Duration
	Priority     ActionPriority

	// Action to perform (required).
	OnAction func(sim *Simulation)
	// Cleanup when the action is cancelled (optional).
	CleanUp func(sim *Simulation)

	cancelled bool
	consumed  bool
	canPool   bool // Flags the PA as safe to use in shared object pools.
}

func (pa *PendingAction) IsConsumed() bool {
	return pa == nil || pa.consumed
}

func (pa *PendingAction) Cancel(sim *Simulation) {
	if pa.cancelled {
		return
	}

	if pa.CleanUp != nil {
		pa.CleanUp(sim)
		pa.CleanUp = nil
	}

	pa.cancelled = true

	if i := slices.Index(sim.pendingActions, pa); i != -1 {
		sim.pendingActions = append(sim.pendingActions[:i], sim.pendingActions[i+1:]...)
	}
}
