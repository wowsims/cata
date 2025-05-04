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
	fireBreath := FireBreathDebuff(&target)
	lightningBreath := MakePermanent(LightningBreathAura(&target))

	// Trauma in this case should *never* be overwritten
	// as its duration from 'MakePermanent' should make it non overwritable by 1 min duration mangles
	fireBreath.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	lightningBreath.Activate(sim)

	if !(lightningBreath.IsActive() && !fireBreath.IsActive()) {
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
	fireBreath := FireBreathDebuff(&target)
	lightningBreath := LightningBreathAura(&target)

	fireBreath.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	lightningBreath.Activate(sim)

	// In this case mangle should overwrite trauma as mangle will give a greater duration

	if !(lightningBreath.IsActive() && !fireBreath.IsActive()) {
		t.Fatalf("longer duration exclusive aura failed to overwrite")
	}
}
