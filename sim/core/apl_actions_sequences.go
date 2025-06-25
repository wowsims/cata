package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
)

type APLActionSequence struct {
	defaultAPLActionImpl
	unit       *Unit
	name       string
	subactions []*APLAction
	curIdx     int
}

func (rot *APLRotation) newActionSequence(config *proto.APLActionSequence) APLActionImpl {
	subactions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return rot.newAPLAction(action)
	})
	subactions = FilterSlice(subactions, func(action *APLAction) bool { return action != nil })
	if len(subactions) == 0 {
		return nil
	}

	return &APLActionSequence{
		unit:       rot.unit,
		name:       config.Name,
		subactions: subactions,
	}
}
func (action *APLActionSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.subactions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionSequence) Finalize(rot *APLRotation) {
	for _, subaction := range action.subactions {
		subaction.impl.Finalize(rot)
	}
}
func (action *APLActionSequence) PostFinalize(rot *APLRotation) {
	for _, subaction := range action.subactions {
		subaction.impl.PostFinalize(rot)
	}
}
func (action *APLActionSequence) Reset(*Simulation) {
	action.curIdx = 0
}
func (action *APLActionSequence) IsReady(sim *Simulation) bool {
	action.unit.Rotation.inSequence = true
	isReady := (action.curIdx < len(action.subactions)) && action.subactions[action.curIdx].IsReady(sim)
	action.unit.Rotation.inSequence = false
	return isReady
}
func (action *APLActionSequence) Execute(sim *Simulation) {
	action.unit.Rotation.inSequence = true
	action.subactions[action.curIdx].Execute(sim)

	if action.unit.CanQueueSpell(sim) {
		// Only advance to the next step in the sequence if we actually cast a spell rather than simply queueing one up.
		action.curIdx++
	} else {
		// If we did queue up a spell, then modify the queue action to advance the sequence when it fires.
		queueAction := action.unit.QueuedSpell.queueAction
		oldFunc := queueAction.OnAction
		queueAction.OnAction = func(sim *Simulation) {
			oldFunc(sim)
			action.curIdx++
			queueAction.OnAction = oldFunc
		}
		action.unit.SetRotationTimer(sim, queueAction.NextActionAt+time.Duration(1))
	}

	action.unit.Rotation.inSequence = false
}
func (action *APLActionSequence) String() string {
	return "Sequence(" + strings.Join(MapSlice(action.subactions, func(subaction *APLAction) string { return fmt.Sprintf("(%s)", subaction) }), "+") + ")"
}

type APLActionResetSequence struct {
	defaultAPLActionImpl
	name     string
	sequence *APLActionSequence
}

func (rot *APLRotation) newActionResetSequence(config *proto.APLActionResetSequence) APLActionImpl {
	if config.SequenceName == "" {
		rot.ValidationMessage(proto.LogLevel_Warning, "Reset Sequence must provide a sequence name")
		return nil
	}
	return &APLActionResetSequence{
		name: config.SequenceName,
	}
}
func (action *APLActionResetSequence) Finalize(rot *APLRotation) {
	for _, otherAction := range rot.allAPLActions() {
		if sequence, ok := otherAction.impl.(*APLActionSequence); ok && sequence.name == action.name {
			action.sequence = sequence
			return
		}
	}
	rot.ValidationMessage(proto.LogLevel_Warning, "No sequence with name: '%s'", action.name)
}
func (action *APLActionResetSequence) IsReady(sim *Simulation) bool {
	return action.sequence != nil && action.sequence.curIdx != 0
}
func (action *APLActionResetSequence) Execute(sim *Simulation) {
	action.sequence.curIdx = 0
}
func (action *APLActionResetSequence) String() string {
	return fmt.Sprintf("Reset Sequence(name = '%s')", action.name)
}

type APLActionStrictSequence struct {
	defaultAPLActionImpl
	unit       *Unit
	subactions []*APLAction
	curIdx     int

	subactionSpells []*Spell
}

func (rot *APLRotation) newActionStrictSequence(config *proto.APLActionStrictSequence) APLActionImpl {
	subactions := MapSlice(config.Actions, func(action *proto.APLAction) *APLAction {
		return rot.newAPLAction(action)
	})
	subactions = FilterSlice(subactions, func(action *APLAction) bool { return action != nil })
	if len(subactions) == 0 {
		return nil
	}

	return &APLActionStrictSequence{
		unit:       rot.unit,
		subactions: subactions,
	}
}
func (action *APLActionStrictSequence) GetInnerActions() []*APLAction {
	return Flatten(MapSlice(action.subactions, func(action *APLAction) []*APLAction { return action.GetAllActions() }))
}
func (action *APLActionStrictSequence) Finalize(rot *APLRotation) {
	for _, subaction := range action.subactions {
		subaction.impl.Finalize(rot)
		action.subactionSpells = append(action.subactionSpells, subaction.GetAllSpells()...)
	}
}
func (action *APLActionStrictSequence) PostFinalize(rot *APLRotation) {
	for _, subaction := range action.subactions {
		subaction.impl.PostFinalize(rot)
	}
}
func (action *APLActionStrictSequence) Reset(*Simulation) {
	action.curIdx = 0
	action.unit.Rotation.inSequence = false
}
func (action *APLActionStrictSequence) IsReady(sim *Simulation) bool {
	action.unit.Rotation.inSequence = true

	if action.unit.GCD.TimeToReady(sim) > MaxSpellQueueWindow {
		action.unit.Rotation.inSequence = false
		return false
	}
	if !action.subactions[0].IsReady(sim) {
		action.unit.Rotation.inSequence = false
		return false
	}
	for _, spell := range action.subactionSpells {
		if !spell.IsReady(sim) {
			action.unit.Rotation.inSequence = false
			return false
		}
	}

	return true
}
func (action *APLActionStrictSequence) Execute(sim *Simulation) {
	action.unit.Rotation.pushControllingAction(action)
}
func (action *APLActionStrictSequence) relinquishControl() {
	action.curIdx = 0
	action.unit.Rotation.inSequence = false
	action.unit.Rotation.popControllingAction(action)
}
func (action *APLActionStrictSequence) advanceSequence() {
	action.curIdx++
	if action.curIdx == len(action.subactions) {
		action.relinquishControl()
	}
}
func (action *APLActionStrictSequence) GetNextAction(sim *Simulation) *APLAction {
	if action.subactions[action.curIdx].IsReady(sim) {
		nextAction := action.subactions[action.curIdx]

		if action.unit.GCD.IsReady(sim) {
			action.advanceSequence()
		} else {
			pa := sim.GetConsumedPendingActionFromPool()
			pa.NextActionAt = action.unit.NextGCDAt()
			pa.Priority = ActionPriorityLow

			pa.OnAction = func(_ *Simulation) {
				if action.unit.Rotation.inSequence {
					action.advanceSequence()
				}
			}

			sim.AddPendingAction(pa)
			action.unit.SetRotationTimer(sim, pa.NextActionAt+time.Duration(1))
		}

		return nextAction
	} else if action.unit.GCD.TimeToReady(sim) <= MaxSpellQueueWindow {
		// If the GCD is ready when the next subaction isn't, it means the sequence is bad
		// so reset and exit the sequence.
		action.relinquishControl()
		return action.unit.Rotation.getNextAction(sim)
	} else {
		// Return nil to wait for the GCD to become ready.
		return nil
	}
}
func (action *APLActionStrictSequence) String() string {
	return "Strict Sequence(" + strings.Join(MapSlice(action.subactions, func(subaction *APLAction) string { return fmt.Sprintf("(%s)", subaction) }), "+") + ")"
}
