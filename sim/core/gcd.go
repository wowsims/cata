package core

import (
	"time"
)

// Note that this is only used when the hardcast and GCD actions happen at different times.
func (unit *Unit) newHardcastAction(sim *Simulation) {
	if unit.hardcastAction != nil && !unit.hardcastAction.consumed {
		unit.hardcastAction.Cancel(sim)
		unit.hardcastAction = nil
	}

	if unit.hardcastAction == nil {
		pa := &PendingAction{
			NextActionAt: unit.Hardcast.Expires,
			OnAction: func(sim *Simulation) {
				if hc := &unit.Hardcast; hc.Expires != startingCDTime && hc.Expires <= sim.CurrentTime {
					hc.Expires = startingCDTime
					if hc.OnComplete != nil {
						hc.OnComplete(sim, hc.Target)
					}
				}
			},
		}
		unit.hardcastAction = pa
	} else {
		unit.hardcastAction.cancelled = false
		unit.hardcastAction.NextActionAt = unit.Hardcast.Expires
	}

	sim.AddPendingAction(unit.hardcastAction)
}

func (unit *Unit) NextGCDAt() time.Duration {
	return unit.GCD.ReadyAt()
}

func (unit *Unit) NextRotationActionAt() time.Duration {
	return unit.RotationTimer.ReadyAt()
}

func (unit *Unit) SetGCDTimer(sim *Simulation, gcdReadyAt time.Duration) {
	if unit.rotationAction == nil {
		return
	}

	unit.GCD.Set(gcdReadyAt)
	unit.SetRotationTimer(sim, gcdReadyAt)
}

func (unit *Unit) SetRotationTimer(sim *Simulation, rotationReadyAt time.Duration) {
	if unit.rotationAction == nil {
		return
	}

	unit.RotationTimer.Set(rotationReadyAt)

	if unit.rotationAction.consumed {
		unit.rotationAction.cancelled = false
		unit.rotationAction.NextActionAt = rotationReadyAt
	} else {
		unit.rotationAction.Cancel(sim)
		oldAction := unit.rotationAction.OnAction
		unit.rotationAction = &PendingAction{
			NextActionAt: rotationReadyAt,
			Priority:     ActionPriorityGCD,
			OnAction:     oldAction,
		}
	}
	sim.AddPendingAction(unit.rotationAction)
}

// Call this when reacting to events that occur before the next scheduled rotation action
func (unit *Unit) ReactToEvent(sim *Simulation) {
	unit.RotationTimer.Reset()
	unit.Rotation.DoNextAction(sim)
}

// Call this to stop the GCD loop for a unit.
// This is mostly used for pets that get summoned / expire.
func (unit *Unit) CancelGCDTimer(sim *Simulation) {
	unit.rotationAction.Cancel(sim)
}

func (unit *Unit) WaitUntil(sim *Simulation, readyTime time.Duration) {
	if readyTime < sim.CurrentTime {
		panic(unit.Label + ": cannot wait negative time")
	}
	unit.SetGCDTimer(sim, readyTime)
	if sim.Log != nil && readyTime > sim.CurrentTime {
		unit.Log(sim, "Pausing GCD for %s due to rotation / CDs.", readyTime-sim.CurrentTime)
	}
}
