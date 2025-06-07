package core

import (
	"math"
	"testing"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func TestFirstCheckOnPull(t *testing.T) {
	sim := &Simulation{}
	char := &Character{
		Unit: Unit{},
	}

	const expectedChance = 0.74

	proc := NewRPPMProc(1.20)
	procChance := proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("First proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestTwoChecksInOneStep(t *testing.T) {
	sim := SetupFakeSim()
	char := &Character{
		Unit: Unit{},
	}

	const expectedChance = 0.0

	proc := NewRPPMProc(1.20)
	proc.Proc(char, sim)
	procChance := proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Second proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestResetSetsCorrectState(t *testing.T) {
	sim := SetupFakeSim()
	char := &Character{
		Unit: Unit{},
	}

	const expectedChance = 0.74
	proc := NewRPPMProc(1.20)
	proc.Proc(char, sim)
	proc.Reset()
	procChance := proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Second proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestSpecModAppliesToCorrectSpec(t *testing.T) {
	sim := SetupFakeSim()
	char := &Character{
		Unit: Unit{},
		Spec: proto.Spec_SpecAfflictionWarlock,
	}

	const expectedChance = 0.74
	proc := NewRPPMProc(1.20).WithSpecMod(0.5, proto.Spec_SpecArmsWarrior)
	procChance := proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	const expectedChanceWithMod = 0.37
	proc.WithSpecMod(-0.5, proto.Spec_SpecAfflictionWarlock)
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChanceWithMod) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChanceWithMod, procChance)
	}
}

func TestClassModAppliesToCorrectClass(t *testing.T) {
	sim := SetupFakeSim()
	char := &Character{
		Unit:  Unit{},
		Class: proto.Class_ClassWarlock,
	}

	const expectedChance = 0.74
	proc := NewRPPMProc(1.20).WithClassMod(0.5, 1) // Class Mask Warrior
	procChance := proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	const expectedChanceWithMod = 0.37
	proc.WithClassMod(-0.5, 256) // Class Mask Warlock
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChanceWithMod) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChanceWithMod, procChance)
	}
}

func TestHasteRatingMod(t *testing.T) {
	sim := SetupFakeSim()

	stats := stats.Stats{stats.HasteRating: HasteRatingPerHastePercent * 50}
	char := &Character{
		Unit: Unit{
			stats: stats,
		},
		Class: proto.Class_ClassWarlock,
	}

	char.PseudoStats.AttackSpeedMultiplier = 1
	char.PseudoStats.MeleeSpeedMultiplier = 1
	char.PseudoStats.CastSpeedMultiplier = 1
	char.PseudoStats.RangedSpeedMultiplier = 1

	expectedChance := 0.74 * 1.5
	proc := NewRPPMProc(1.20).WithHasteMod(1, MeleeHaste)
	procChance := proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	char.PseudoStats.CastSpeedMultiplier = 1.5
	char.PseudoStats.RangedSpeedMultiplier = 1.5
	char.updateCastSpeed()

	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	expectedChance = 0.74 * 1.5 * 1.5
	proc = NewRPPMProc(1.2).WithHasteMod(1, RangedHaste)
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	proc = NewRPPMProc(1.2).WithHasteMod(1, SpellHaste)
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	char.PseudoStats.AttackSpeedMultiplier = 1.5
	char.updateCastSpeed()
	expectedChance = 0.74 * 1.5 * 1.5 * 1.5
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	expectedChance = 0.74 * 1.5 * 1.5
	proc = NewRPPMProc(1.2).WithHasteMod(1, HighestHaste)
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestCritRatingMod(t *testing.T) {
	sim := SetupFakeSim()

	stats := stats.Stats{stats.PhysicalCritPercent: 50, stats.SpellCritPercent: 100}
	char := &Character{
		Unit: Unit{
			stats: stats,
		},
		Class: proto.Class_ClassWarlock,
	}

	expectedChance := 0.74 * 1.5
	proc := NewRPPMProc(1.20).WithCritMod(1, MeleeCrit)
	procChance := proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	proc = NewRPPMProc(1.20).WithCritMod(1, RangedCrit)
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	proc = NewRPPMProc(1.20).WithCritMod(1, LowestCrit)
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	expectedChance = 0.74 * 2
	proc = NewRPPMProc(1.20).WithCritMod(1, SpellCrit)
	procChance = proc.getProcChance(char, sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}
