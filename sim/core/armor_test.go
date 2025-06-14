package core

import (
	"testing"

	"github.com/wowsims/mop/sim/core/stats"
)

func TestWeakenedArmorStacks(t *testing.T) {
	sim := Simulation{}
	baseArmor := 24835.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        93,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 24835.0, target.Armor())
	}
	stacks := int32(1)
	sunderAura := WeakenedArmorAura(&target)
	sunderAura.Activate(&sim)
	sunderAura.SetStacks(&sim, stacks)
	tolerance := 0.001
	for stacks <= 3 {
		expectedArmor = baseArmor * (1.0 - float64(stacks)*0.04)
		if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
			t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
		}
		stacks++
		sunderAura.AddStack(&sim)
	}
}

func TestDamageReductionFromArmor(t *testing.T) {
	sim := Simulation{}
	baseArmor := 24835.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        93,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	attacker := Unit{
		Type:  PlayerUnit,
		Level: 90,
	}
	target.stats = target.initialStats
	expectedDamageReduction := 0.349334
	attackTable := NewAttackTable(&attacker, &target)
	tolerance := 0.0001
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.getArmorDamageModifier(), tolerance) {
		t.Fatalf("Expected no armor modifiers to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.getArmorDamageModifier())
	}

	// Major
	weakenedArmorAura := WeakenedArmorAura(&target)
	weakenedArmorAura.Activate(&sim)
	weakenedArmorAura.SetStacks(&sim, 3)
	expectedDamageReduction = 0.320864
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.getArmorDamageModifier(), tolerance) {
		t.Fatalf("Expected major armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.getArmorDamageModifier())
	}
	weakenedArmorAura.Deactivate(&sim)

	// Major Multi
	shatteringThrowAura := ShatteringThrowAura(&target, attacker.UnitIndex)
	shatteringThrowAura.Activate(&sim)
	expectedDamageReduction = 0.300459
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.getArmorDamageModifier(), tolerance) {
		t.Fatalf("Expected major & shattering modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.getArmorDamageModifier())
	}
}
