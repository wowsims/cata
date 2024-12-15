package core

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
)

type MovementAction struct {
	PendingAction
	srcPosition float64       // starting position
	startTime   time.Duration // starting time of the movement
	speed       float64       // theoretical movement speed, can be 0
}

func (action *MovementAction) GetCurrentPosition(sim *Simulation) float64 {
	return action.srcPosition + float64(sim.CurrentTime-action.startTime)*action.speed/float64(time.Second)
}

func (unit *Unit) initMovement() {
	unit.moveAura = unit.GetOrRegisterAura(Aura{
		Label:     "Movement",
		ActionID:  ActionID{OtherID: proto.OtherAction_OtherActionMove},
		Duration:  NeverExpires,
		MaxStacks: 30,

		OnGain: func(aura *Aura, sim *Simulation) {
			unit.Moving = true
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			unit.Moving = false
			unit.movementAction = nil
		},
	})

	unit.moveSpell = unit.GetOrRegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionMove},
		Flags:    SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			unit.moveAura.Activate(sim)
			unit.moveAura.SetStacks(sim, max(int32(unit.DistanceFromTarget), 1))
		},
	})
}

func (unit *Unit) MoveTo(moveRange float64, sim *Simulation) {
	if moveRange == unit.DistanceFromTarget {
		return
	}

	unit.UpdatePosition(sim)
	moveDistance := moveRange - unit.DistanceFromTarget
	timeToMove := time.Duration(math.Abs(moveDistance)/unit.GetMovementSpeed()*1000) * time.Millisecond
	registerMovementAction(unit, sim, unit.GetMovementSpeed()*TernaryFloat64(moveDistance < 0, -1., 1.), sim.CurrentTime+timeToMove)
}

func (unit *Unit) MoveDuration(duration time.Duration, sim *Simulation) {
	if duration == 0 {
		return
	}

	unit.UpdatePosition(sim)
	registerMovementAction(unit, sim, 0., sim.CurrentTime+duration)
}

func (unit *Unit) UpdatePosition(sim *Simulation) {
	if !unit.Moving {
		return
	}

	oldDist := unit.DistanceFromTarget
	unit.DistanceFromTarget = unit.movementAction.GetCurrentPosition(sim)
	if oldDist == unit.DistanceFromTarget {
		return
	}

	unit.OnMovement(sim, unit.DistanceFromTarget, MovementUpdate)

	// update auto attack state
	if unit.AutoAttacks.mh.enabled != unit.AutoAttacks.mh.IsInRange() {
		if unit.AutoAttacks.mh.IsInRange() {
			unit.AutoAttacks.EnableMeleeSwing(sim)
		} else {
			unit.AutoAttacks.CancelMeleeSwing(sim)
		}
	}

	if unit.AutoAttacks.ranged.enabled != unit.AutoAttacks.ranged.IsInRange() {
		if unit.AutoAttacks.ranged.IsInRange() {
			unit.AutoAttacks.EnableRangedSwing(sim)
		} else {
			unit.AutoAttacks.CancelRangedSwing(sim)
		}
	}

	yards := max(int32(unit.DistanceFromTarget), 1) // never set to 0 yards as we deactivate the aura
	if yards != unit.moveAura.GetStacks() {
		unit.moveAura.SetStacks(sim, yards)
	}
}

func (unit *Unit) FinalizeMovement(sim *Simulation) {
	if !unit.Moving {
		return
	}

	unit.UpdatePosition(sim)
	unit.moveAura.Deactivate(sim)

	unit.OnMovement(sim, unit.DistanceFromTarget, MovementEnd)
}

func registerMovementAction(unit *Unit, sim *Simulation, speed float64, endTime time.Duration) {
	if unit.movementAction != nil {
		unit.movementAction.Cancel(sim)
	} else {
		unit.moveSpell.Cast(sim, unit.CurrentTarget)
	}

	movementAction := MovementAction{
		startTime:   sim.CurrentTime,
		speed:       speed,
		srcPosition: unit.DistanceFromTarget,
	}

	movementAction.NextActionAt = endTime
	movementAction.OnAction = func(sim *Simulation) {
		unit.FinalizeMovement(sim)
	}

	unit.OnMovement(sim, unit.DistanceFromTarget, MovementStart)
	unit.movementAction = &movementAction
	sim.AddPendingAction(&movementAction.PendingAction)
}

type MovementUpdateType byte

const (
	MovementStart = iota
	MovementUpdate
	MovementEnd
)

type MovementCallback func(sim *Simulation, position float64, kind MovementUpdateType)

func (unit *Unit) RegisterMovementCallback(callback MovementCallback) {
	unit.movementCallbacks = append(unit.movementCallbacks, callback)
}

func (unit *Unit) OnMovement(sim *Simulation, position float64, kind MovementUpdateType) {
	for _, movementCallback := range unit.movementCallbacks {
		movementCallback(sim, position, kind)
	}
}

func (unit *Unit) MultiplyMovementSpeed(sim *Simulation, amount float64) {
	oldMultiplier := unit.PseudoStats.MovementSpeedMultiplier
	oldSpeed := unit.GetMovementSpeed()
	unit.PseudoStats.MovementSpeedMultiplier *= amount
	if sim.Log != nil {
		unit.Log(sim, "[DEBUG] Movement speed changed from %.2f (%.2f%%) to %.2f (%.2f%%)", oldSpeed, (oldMultiplier-1)*100.0, unit.GetMovementSpeed(), (unit.PseudoStats.MovementSpeedMultiplier-1)*100.0)
	}

	// we have a pending movement action that depends on our movement speed
	if unit.movementAction != nil && unit.movementAction.speed != 0 {
		dest := unit.movementAction.speed * float64(unit.movementAction.NextActionAt-unit.movementAction.startTime) / float64(time.Second)
		unit.MoveTo(dest, sim)
	}
}

// Returns the units current movement speed in yards / second
func (unit *Unit) GetMovementSpeed() float64 {
	if unit.Type == PlayerUnit {
		return 7. * unit.PseudoStats.MovementSpeedMultiplier
	}

	return 8. * unit.PseudoStats.MovementSpeedMultiplier
}

func (unit *Unit) NewMovementSpeedAura(label string, actionID ActionID, multiplier float64) *Aura {
	aura := MakePermanent(unit.GetOrRegisterAura(Aura{
		Label:    label,
		ActionID: actionID,
	}))

	aura.NewMovementSpeedEffect(multiplier)

	return aura
}

func (aura *Aura) NewMovementSpeedEffect(multiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("MovementSpeed", true, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyMovementSpeed(sim, 1+multiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyMovementSpeed(sim, 1.0/(1+multiplier))
		},
	})
}
