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
	lightningBreath := MakePermanent(LightningBreathDebuff(&target))
	fireBreath := FireBreathDebuff(&target)

	// Lightning Breath in this case should *never* be overwritten
	// as its duration from 'MakePermanent' should make it non overwritable by Fire Breath
	lightningBreath.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	fireBreath.Activate(sim)

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
	lightningBreath := LightningBreathDebuff(&target)

	fireBreath.Activate(sim)

	sim.CurrentTime = 1 * time.Second

	lightningBreath.Activate(sim)

	// In this case Lightning Breath should overwrite Fire Breath as Lightning Breath will give a greater duration

	if !(lightningBreath.IsActive() && !fireBreath.IsActive()) {
		t.Fatalf("longer duration exclusive aura failed to overwrite")
	}
}
