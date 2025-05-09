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
			Priority:     ActionPriorityGCD,
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

	if !unit.rotationAction.consumed {
		unit.rotationAction.Cancel(sim)
	}

	unit.rotationAction.cancelled = false
	unit.rotationAction.NextActionAt = rotationReadyAt
	sim.AddPendingAction(unit.rotationAction)
}

// Call this when reacting to events that occur before the next scheduled rotation action
func (unit *Unit) ReactToEvent(sim *Simulation) {
	// If the next rotation action was already scheduled for this timestep then execute it now
	unit.Rotation.DoNextAction(sim)

	// Otherwise schedule an evaluation based on reaction time
	if unit.NextRotationActionAt() > sim.CurrentTime+unit.ReactionTime {
		unit.SetRotationTimer(sim, sim.CurrentTime+unit.ReactionTime)
	}
}

// Call this to stop the GCD loop for a unit.
// This is mostly used for pets that get summoned / expire.
func (unit *Unit) CancelGCDTimer(sim *Simulation) {
	unit.rotationAction.Cancel(sim)
}

func (unit *Unit) CancelHardcast(sim *Simulation) {
	unit.Hardcast.Expires = startingCDTime
	unit.SetGCDTimer(sim, sim.CurrentTime+unit.ReactionTime)
}

func (unit *Unit) WaitUntil(sim *Simulation, readyTime time.Duration) {
	if readyTime < sim.CurrentTime {
		panic(unit.Label + ": cannot wait negative time")
	}
	unit.SetRotationTimer(sim, readyTime)
	if sim.Log != nil && readyTime > sim.CurrentTime {
		unit.Log(sim, "Pausing rotation for %s due to resources / CDs.", readyTime-sim.CurrentTime)
	}
}

func (unit *Unit) ExtendGCDUntil(sim *Simulation, readyTime time.Duration) {
	if readyTime < sim.CurrentTime {
		panic(unit.Label + ": cannot wait negative time")
	}
	unit.SetGCDTimer(sim, readyTime)
	if sim.Log != nil && readyTime > sim.CurrentTime {
		unit.Log(sim, "Extending GCD for %s due to rotation / CDs.", readyTime-sim.CurrentTime)
	}
}
