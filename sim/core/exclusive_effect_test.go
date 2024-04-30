package core

import (
	"testing"
	"time"
)

func TestSingleAuraExclusiveDurationNoOverwrite(t *testing.T) {
	sim := &Simulation{}

	target := Unit{
		Type:        EnemyUnit,
		Index:       0,
		Level:       83,
		auraTracker: newAuraTracker(),
	}
	mangle := MangleAura(&target)
	hemorrhage := MakePermanent(HemorrhageAura(&target))

	// Trauma in this case should *never* be overwritten
	// as its duration from 'MakePermanent' should make it non overwritable by 1 min duration mangles
	hemorrhage.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	mangle.Activate(sim)

	if !(hemorrhage.IsActive() && !mangle.IsActive()) {
		t.Fatalf("lower duration exclusive aura overwrote previous!")
	}
}

func TestSingleAuraExclusiveDurationOverwrite(t *testing.T) {
	sim := &Simulation{}

	target := Unit{
		Type:        EnemyUnit,
		Index:       0,
		Level:       83,
		auraTracker: newAuraTracker(),
	}
	mangle := MangleAura(&target)
	hemorrhage := HemorrhageAura(&target)

	hemorrhage.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	mangle.Activate(sim)

	// In this case mangle should overwrite trauma as mangle will give a greater duration

	if !(mangle.IsActive() && !hemorrhage.IsActive()) {
		t.Fatalf("longer duration exclusive aura failed to overwrite")
	}
}
