package core

import (
	"testing"

	"github.com/wowsims/mop/sim/core/stats"
)

func TestSunderArmorStacks(t *testing.T) {
	sim := Simulation{}
	baseArmor := 11977.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        88,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 11977.0, target.Armor())
	}
	stacks := int32(1)
	sunderAura := SunderArmorAura(&target)
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

func TestAcidSpitStacks(t *testing.T) {
	sim := Simulation{}
	baseArmor := 11977.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        88,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 11977.0, target.Armor())
	}
	stacks := int32(1)
	corrosiveSpitAura := CorrosiveSpitAura(&target)
	corrosiveSpitAura.Activate(&sim)
	corrosiveSpitAura.SetStacks(&sim, stacks)
	tolerance := 0.001
	for stacks <= 3 {
		expectedArmor = baseArmor * (1.0 - float64(stacks)*0.04)
		if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
			t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
		}
		stacks++
		corrosiveSpitAura.AddStack(&sim)
	}
}

func TestExposeArmor(t *testing.T) {
	sim := Simulation{}
	baseArmor := 11977.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        88,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 11977.0, target.Armor())
	}
	exposeAura := ExposeArmorAura(&target, false)
	exposeAura.Activate(&sim)
	tolerance := 0.001
	expectedArmor = baseArmor * (1.0 - 0.12)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
}

func TestMajorArmorReductionAurasDoNotStack(t *testing.T) {
	sim := Simulation{}
	baseArmor := 11977.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        88,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	target.stats = target.initialStats
	expectedArmor := baseArmor
	if target.Armor() != expectedArmor {
		t.Fatalf("Armor value for target should be %f but found %f", 11977.0, target.Armor())
	}
	corrosiveSpitAura := CorrosiveSpitAura(&target)
	corrosiveSpitAura.Activate(&sim)
	corrosiveSpitAura.SetStacks(&sim, 3)
	tolerance := 0.001
	expectedArmor = baseArmor * (1.0 - 0.12)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
	exposeArmorAura := ExposeArmorAura(&target, false)
	exposeArmorAura.Activate(&sim)
	if !WithinToleranceFloat64(expectedArmor, target.Armor(), tolerance) {
		t.Fatalf("Armor value for target should be %f but found %f", expectedArmor, target.Armor())
	}
}

func TestDamageReductionFromArmor(t *testing.T) {
	sim := Simulation{}
	baseArmor := 11977.0
	target := Unit{
		Type:         EnemyUnit,
		Index:        0,
		Level:        88,
		auraTracker:  newAuraTracker(),
		initialStats: stats.Stats{stats.Armor: baseArmor},
		PseudoStats:  stats.NewPseudoStats(),
		Metrics:      NewUnitMetrics(),
	}
	attacker := Unit{
		Type:  PlayerUnit,
		Level: 85,
	}
	spell := &Spell{}
	target.stats = target.initialStats
	expectedDamageReduction := 0.314795
	attackTable := NewAttackTable(&attacker, &target)
	tolerance := 0.0001
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected no armor modifiers to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Major
	corrosiveSpitAura := CorrosiveSpitAura(&target)
	corrosiveSpitAura.Activate(&sim)
	corrosiveSpitAura.SetStacks(&sim, 3)
	expectedDamageReduction = 0.287895
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}
	corrosiveSpitAura.Deactivate(&sim)

	// Major
	faerieFireAura := FaerieFireAura(&target)
	faerieFireAura.Activate(&sim)
	faerieFireAura.SetStacks(&sim, 3)
	expectedDamageReduction = 0.287895
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & minor armor modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}

	// Major Multi
	shatteringThrowAura := ShatteringThrowAura(&target, attacker.UnitIndex)
	shatteringThrowAura.Activate(&sim)
	expectedDamageReduction = 0.244387
	if !WithinToleranceFloat64(1-expectedDamageReduction, attackTable.GetArmorDamageModifier(spell), tolerance) {
		t.Fatalf("Expected major & shattering modifier to result in %f damage reduction got %f", expectedDamageReduction, 1-attackTable.GetArmorDamageModifier(spell))
	}
}
