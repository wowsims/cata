package core

import (
	"github.com/wowsims/mop/sim/core/proto"
)

type DynamicProcManager struct {
	procMasks   []ProcMask
	procChances []DynamicProc
}

type DynamicProc interface {
	Reset()
	Chance(unit *Unit, sim *Simulation) float64
	Proc(unit *Unit, sim *Simulation, label string) bool
}

func (dpm *DynamicProcManager) Reset() {
	for _, proc := range dpm.procChances {
		proc.Reset()
	}
}

// Returns whether the effect procced.
func (dpm *DynamicProcManager) Proc(unit *Unit, sim *Simulation, procMask ProcMask, label string) bool {
	for i, m := range dpm.procMasks {
		if m.Matches(procMask) {
			return dpm.procChances[i].Proc(unit, sim, label)
		}
	}

	return false
}

func (dpm *DynamicProcManager) Chance(procMask ProcMask, unit *Unit, sim *Simulation) float64 {
	for i, m := range dpm.procMasks {
		if m.Matches(procMask) {
			return dpm.procChances[i].Chance(unit, sim)
		}
	}

	return 0
}

// PPMManager for static ProcMasks
func (character *Character) NewPPMManager(ppm float64, procMask ProcMask) *DynamicProcManager {
	dpm := character.newDynamicWeaponProcManager(ppm, 0, procMask)

	character.RegisterItemSwapCallback(AllWeaponSlots(), func(sim *Simulation, slot proto.ItemSlot) {
		dpm = character.newDynamicWeaponProcManager(ppm, 0, procMask)
	})

	return &dpm
}

// PPMManager for static ProcMasks and no item swap callback
func (character *Character) NewStaticPPMManager(ppm float64, procMask ProcMask) *DynamicProcManager {
	dpm := character.newDynamicWeaponProcManager(ppm, 0, procMask)

	return &dpm
}

// Dynamic Proc Manager for static ProcMasks and no item swap callback
func (character *Character) NewStaticDynamicProcManager(fixedProcChance float64, procMask ProcMask) *DynamicProcManager {
	dpm := character.newDynamicWeaponProcManager(0, fixedProcChance, procMask)

	return &dpm
}

// Dynamic Proc Manager for dynamic ProcMasks on weapon enchants
func (character *Character) NewDynamicProcForEnchant(effectID int32, ppm float64, fixedProcChance float64) *DynamicProcManager {
	return character.newDynamicProcManagerWithDynamicProcMask(ppm, fixedProcChance, func() ProcMask {
		return character.getCurrentProcMaskForWeaponEnchant(effectID)
	})
}

// Dynamic Proc Manager for dynamic ProcMasks on weapon effects
func (character *Character) NewDynamicProcForWeapon(itemID int32, ppm float64, fixedProcChance float64) *DynamicProcManager {
	return character.newDynamicProcManagerWithDynamicProcMask(ppm, fixedProcChance, func() ProcMask {
		return character.getCurrentProcMaskForWeaponEffect(itemID)
	})
}

func (character *Character) newDynamicProcManagerWithDynamicProcMask(ppm float64, fixedProcChance float64, procMaskFn func() ProcMask) *DynamicProcManager {
	dpm := character.newDynamicWeaponProcManager(ppm, fixedProcChance, procMaskFn())
	character.RegisterItemSwapCallback(AllWeaponSlots(), func(sim *Simulation, slot proto.ItemSlot) {
		dpm = character.newDynamicWeaponProcManager(ppm, fixedProcChance, procMaskFn())
	})

	return &dpm
}

func (character *Character) newDynamicWeaponProcManager(ppm float64, fixedProcChance float64, procMask ProcMask) DynamicProcManager {
	if (ppm != 0) && (fixedProcChance != 0) {
		panic("Cannot simultaneously specify both a ppm and a fixed proc chance!")
	}

	aa := character.AutoAttacks
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return DynamicProcManager{}
	}

	dpm := DynamicProcManager{procMasks: make([]ProcMask, 0, 2), procChances: []DynamicProc{}}

	chances := make([]staticProc, 0, 2)
	mergeOrAppend := func(speed float64, mask ProcMask) {
		if speed == 0 || mask == 0 {
			return
		}

		for i, proc := range chances {
			if proc.chance == speed {
				dpm.procMasks[i] |= mask
				return
			}
		}

		dpm.procMasks = append(dpm.procMasks, mask)
		chances = append(chances, staticProc{chance: speed})
	}

	mergeOrAppend(aa.mh.SwingSpeed, procMask&^ProcMaskRanged&^ProcMaskMeleeOH) // "everything else", even if not explicitly flagged MH
	mergeOrAppend(aa.oh.SwingSpeed, procMask&ProcMaskMeleeOH)
	mergeOrAppend(aa.ranged.SwingSpeed, procMask&ProcMaskRanged)

	for i := range chances {
		if fixedProcChance != 0 {
			chances[i].chance = fixedProcChance
		} else {
			chances[i].chance *= ppm / 60
		}

		dpm.procChances = append(dpm.procChances, chances[i])
	}

	return dpm
}

type staticProc struct {
	chance float64
}

func (sp staticProc) Reset()                                {}
func (sp staticProc) Chance(_ *Unit, _ *Simulation) float64 { return sp.chance }
func (sp staticProc) Proc(_ *Unit, sim *Simulation, label string) bool {
	return sim.Proc(sp.chance, label)
}

func (character *Character) NewRPPMProcManagerForEnchant(ppm float64, procMask ProcMask, enchantId int32, configure func(*RPPMProc)) *DynamicProcManager {
	builder := func() DynamicProcManager {
		mh := character.MainHand()
		oh := character.OffHand()

		realPPM := ppm
		if mh != nil && oh != nil && mh.Enchant.EffectID == enchantId && oh.Enchant.EffectID == enchantId {
			realPPM *= 2
		}
		return *character.NewRPPMProcManager(realPPM, procMask, configure)
	}

	dpm := builder()
	character.RegisterItemSwapCallback(MeleeWeaponSlots(), func(s *Simulation, is proto.ItemSlot) {
		dpm = builder()
	})

	return &dpm
}

// Creates a new RPPM proc manager
//
//	configure func(*RPPMProc)
//
// argument can be used to modify the properties of the RPPM proc
//
// Enchants should be configured through
//
//	NewRPPMProcManagerForEnchant()
//
// # Example
//
//	character.NewRPPMProcManager(1.2, core.ProcMaskMelee | core.ProcMaskSpellDamage, func(r *core.RPPMProc) {
//		r.WithSpecMod(-0.4, proto.Spec_SpecAfflictionWarlock).
//		WithHasteMod(1, core.HighestHaste)
//	})
func (character *Character) NewRPPMProcManager(ppm float64, procMask ProcMask, configure func(*RPPMProc)) *DynamicProcManager {
	proc := NewRPPMProc(ppm)
	if configure != nil {
		configure(proc)
	}

	proc.ForCharacter(character)
	return &DynamicProcManager{
		procMasks:   []ProcMask{procMask},
		procChances: []DynamicProc{proc},
	}
}
