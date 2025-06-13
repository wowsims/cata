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

	proc := NewRPPMProc(char, RPPMConfig{PPM: 1.2})
	procChance := proc.Chance(sim)
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

	proc := NewRPPMProc(char, RPPMConfig{PPM: 1.2})
	proc.Proc(sim, "UnitTest")
	procChance := proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Second proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestResetSetsCorrectState(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1}, false)

	const expectedChance = 0.74
	proc := NewRPPMProc(char, RPPMConfig{PPM: 1.2})
	proc.Proc(sim, "UnitTest")
	proc.Reset()
	procChance := proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Second proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestSpecModAppliesToCorrectSpec(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1}, false)

	const expectedChance = 0.74
	proc := NewRPPMProc(char, RPPMConfig{PPM: 1.2}.WithSpecMod(0.5, proto.Spec_SpecArmsWarrior))
	procChance := proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	const expectedChanceWithMod = 0.1
	proc = NewRPPMProc(char, RPPMConfig{PPM: 1.2}.WithSpecMod(-0.5, proto.Spec_SpecAfflictionWarlock))
	procChance = proc.Chance(sim)
	if math.Abs(procChance-expectedChanceWithMod) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChanceWithMod, procChance)
	}
}

func TestClassModAppliesToCorrectClass(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1}, false)

	const expectedChance = 0.74
	proc := NewRPPMProc(char, RPPMConfig{PPM: 1.2}.WithClassMod(0.5, 1))
	procChance := proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	const expectedChanceWithMod = 0.1
	proc = NewRPPMProc(char, RPPMConfig{PPM: 1.2}.WithClassMod(-0.5, 256)) // Class Mask Warlock
	procChance = proc.Chance(sim)
	if math.Abs(procChance-expectedChanceWithMod) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChanceWithMod, procChance)
	}
}

func TestHasteRatingMod(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1}, false)

	char.stats = stats.Stats{stats.HasteRating: HasteRatingPerHastePercent * 50}
	char.PseudoStats.AttackSpeedMultiplier = 1
	char.PseudoStats.MeleeSpeedMultiplier = 1.5
	char.PseudoStats.CastSpeedMultiplier = 1
	char.PseudoStats.RangedSpeedMultiplier = 1

	expectedChance := 2.19
	proc := NewRPPMProc(char, RPPMConfig{PPM: 1.2}.WithHasteMod())
	procChance := proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	char.PseudoStats.RangedSpeedMultiplier = 1.5
	char.updateAttackSpeed()
	procChance = proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	char.PseudoStats.AttackSpeedMultiplier = 1.5
	char.updateAttackSpeed()
	expectedChance = 5.715
	procChance = proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	char.PseudoStats.CastSpeedMultiplier = 1.5
	char.updateCastSpeed()
	procChance = proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	char.PseudoStats.CastSpeedMultiplier = 2
	char.updateCastSpeed()
	expectedChance = 10.86
	procChance = proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestCritRatingMod(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1}, false)
	char.stats = stats.Stats{stats.PhysicalCritPercent: 50}

	expectedChance := 2.19
	proc := NewRPPMProc(char, RPPMConfig{PPM: 1.2}.WithCritMod())
	procChance := proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}

	char.stats[stats.SpellCritPercent] = 100
	expectedChance = 4.36
	procChance = proc.Chance(sim)
	if math.Abs(procChance-expectedChance) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChance, procChance)
	}
}

func TestProcManagerShouldProcForItemEffect(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1}, false)

	procManager := char.NewRPPMProcManager(1, false, ProcMaskDirect, RPPMConfig{PPM: 2})
	if !procManager.Proc(sim, ProcMaskMeleeMHAuto, "Melee Swing") {
		t.Fatal("Did not proc for 100% proc chance")
	}
}

func TestProcManagerShouldBeConfigurable(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotTrinket1}, false)

	const expectedChanceWithMod = 0.1
	procManager := char.NewRPPMProcManager(1, false, ProcMaskDirect, RPPMConfig{PPM: 1.2}.WithClassMod(-0.5, 256))
	procChance := procManager.Chance(ProcMaskMelee, sim)
	if math.Abs(procChance-expectedChanceWithMod) > 0.001 {
		t.Fatalf("Proc chance wrong. Expected %f, got %f", expectedChanceWithMod, procChance)
	}
}

func TestProcManagerShouldProcOffCorrectWeaponForItemEffect(t *testing.T) {
	sim := SetupFakeSim()
	masks := []ProcMask{
		ProcMaskMeleeMH,
		ProcMaskMeleeOH,
		ProcMaskRanged,
	}

	for _, mask := range masks {
		itemSlot := proto.ItemSlot_ItemSlotMainHand
		if mask.Matches(ProcMaskMeleeOH) {
			itemSlot = proto.ItemSlot_ItemSlotOffHand
		}

		char := GetFakeCharacter([]proto.ItemSlot{itemSlot}, mask.Matches(ProcMaskRanged))

		procManager := char.NewRPPMProcManager(1, false, ProcMaskDirect, RPPMConfig{PPM: 2})
		if mask.Matches(ProcMaskMeleeOHAuto) != (procManager.Chance(ProcMaskMeleeOHAuto, sim) > 0) {
			t.Fatal("Wrong proc chance")
		}

		if mask.Matches(ProcMaskRanged) != (procManager.Chance(ProcMaskRanged, sim) > 0) {
			t.Fatal("Wrong proc chance")
		}

		if mask.Matches(ProcMaskMeleeMHAuto) != (procManager.Chance(ProcMaskMeleeMHAuto, sim) > 0) {
			t.Fatal("Wrong proc chance")
		}
	}
}

func TestProcManagerShouldProcOffCorrectWeaponForEnchantEffect(t *testing.T) {
	sim := SetupFakeSim()
	masks := []ProcMask{
		ProcMaskMeleeMH,
		ProcMaskMeleeOH,
		ProcMaskRanged,
	}

	for _, mask := range masks {
		itemSlot := proto.ItemSlot_ItemSlotMainHand
		if mask.Matches(ProcMaskMeleeOH) {
			itemSlot = proto.ItemSlot_ItemSlotOffHand
		}

		char := GetFakeCharacter([]proto.ItemSlot{itemSlot}, mask.Matches(ProcMaskRanged))
		procManager := char.NewRPPMProcManager(1234, true, ProcMaskDirect, RPPMConfig{PPM: 2})
		if mask.Matches(ProcMaskMeleeOHAuto) != (procManager.Chance(ProcMaskMeleeOHAuto, sim) > 0) {
			t.Fatal("Wrong proc chance")
		}

		if mask.Matches(ProcMaskRanged) != (procManager.Chance(ProcMaskRanged, sim) > 0) {
			t.Fatal("Wrong proc chance")
		}

		if mask.Matches(ProcMaskMeleeMHAuto) != (procManager.Chance(ProcMaskMeleeMHAuto, sim) > 0) {
			t.Fatal("Wrong proc chance")
		}
	}
}

func TestProcManagerShouldProcIndependentForSameEffect(t *testing.T) {
	sim := SetupFakeSim()
	char := GetFakeCharacter([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}, false)
	procManager := char.NewRPPMProcManager(1234, true, ProcMaskDirect, RPPMConfig{PPM: 2})
	if !procManager.Proc(sim, ProcMaskMeleeMHAuto, "MH Auto") {
		t.Fatal("Wrong proc")
	}

	if procManager.Chance(ProcMaskMeleeMH, sim) == procManager.Chance(ProcMaskMeleeOH, sim) {
		t.Fatal("Proc chances are not independent")
	}

	if !procManager.Proc(sim, ProcMaskMeleeOHAuto, "OH Auto") {
		t.Fatal("Wrong proc")
	}
}

func GetFakeCharacter(slots []proto.ItemSlot, withRangedWeapon bool) *Character {
	character := &Character{
		Unit:  Unit{},
		Spec:  proto.Spec_SpecAfflictionWarlock,
		Class: proto.Class_ClassWarlock,
	}
	character.ItemSwap.character = character
	for idx := range slots {
		item := Item{
			ID: int32(idx + 1),
			Enchant: Enchant{
				EffectID: 1234,
			},
			ScalingOptions: map[int32]*proto.ScalingItemProperties{
				int32(proto.ItemLevelState_Base): {
					Ilvl: 528,
				},
			},
		}

		if slots[idx] == proto.ItemSlot_ItemSlotMainHand && withRangedWeapon {
			item.RangedWeaponType = proto.RangedWeaponType_RangedWeaponTypeBow
			item.Type = proto.ItemType_ItemTypeRanged
		} else if slots[idx] == proto.ItemSlot_ItemSlotMainHand || slots[idx] == proto.ItemSlot_ItemSlotOffHand {
			item.Type = proto.ItemType_ItemTypeWeapon
			item.HandType = proto.HandType_HandTypeOneHand
		} else {
			item.Type = proto.ItemType_ItemTypeTrinket
		}

		character.Equipment[slots[idx]] = item
		ItemsByID[int32(idx+1)] = item
	}

	return character
}
