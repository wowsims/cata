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
	Chance(sim *Simulation) float64
	Proc(sim *Simulation, label string) bool
}

func (dpm *DynamicProcManager) Reset() {
	for _, proc := range dpm.procChances {
		proc.Reset()
	}
}

// Returns whether the effect procced.
func (dpm *DynamicProcManager) Proc(sim *Simulation, procMask ProcMask, label string) bool {
	for i, m := range dpm.procMasks {
		if m.Matches(procMask) {
			return dpm.procChances[i].Proc(sim, label)
		}
	}

	return false
}

func (dpm *DynamicProcManager) Chance(procMask ProcMask, sim *Simulation) float64 {
	for i, m := range dpm.procMasks {
		if m.Matches(procMask) {
			return dpm.procChances[i].Chance(sim)
		}
	}

	return 0
}

// PPMManager for static ProcMasks
func (character *Character) NewLegacyPPMManager(ppm float64, procMask ProcMask) *DynamicProcManager {
	dpm := character.newDynamicWeaponProcManager(ppm, 0, procMask)

	character.RegisterItemSwapCallback(AllWeaponSlots(), func(sim *Simulation, slot proto.ItemSlot) {
		dpm = character.newDynamicWeaponProcManager(ppm, 0, procMask)
	})

	return &dpm
}

// PPMManager for static ProcMasks and no item swap callback
func (character *Character) NewStaticLegacyPPMManager(ppm float64, procMask ProcMask) *DynamicProcManager {
	dpm := character.newDynamicWeaponProcManager(ppm, 0, procMask)

	return &dpm
}

// Dynamic Proc Manager for static ProcMasks and no item swap callback
func (character *Character) NewFixedProcChanceManager(fixedProcChance float64, procMask ProcMask) *DynamicProcManager {
	dpm := character.newDynamicWeaponProcManager(0, fixedProcChance, procMask)

	return &dpm
}

// Dynamic Proc Manager for dynamic ProcMasks on weapon enchants
func (character *Character) NewDynamicLegacyProcForEnchant(effectID int32, ppm float64, fixedProcChance float64) *DynamicProcManager {
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

func (sp staticProc) Reset()                       {}
func (sp staticProc) Chance(_ *Simulation) float64 { return sp.chance }
func (sp staticProc) Proc(sim *Simulation, label string) bool {
	return sim.Proc(sp.chance, label)
}

type EffectType byte

const (
	ItemEffect EffectType = iota
	EnchantEffect
	WeaponEffect // Procs directly defined on weapons
)

// Creates a new RPPM proc manager for the given effectID.
// Will manage all equiped items that use the given effect ID and overwrite the given configuration's ilvl accordingly.
//
// # Example
//
//	character.NewRPPMProcManager(1.2, core.ProcMaskMelee | core.ProcMaskSpellDamage, func(r *core.RPPMProc) {
//		return r.WithSpecMod(-0.4, proto.Spec_SpecAfflictionWarlock).
//		WithHasteMod(1, core.HighestHaste)
//	})
func (character *Character) NewRPPMProcManager(effectID int32, effectType EffectType, procMask ProcMask, rppmConfig RPPMConfig) *DynamicProcManager {
	var slotList []proto.ItemSlot
	builder := func() DynamicProcManager {
		slotList = []proto.ItemSlot{}
		manager := DynamicProcManager{
			procMasks:   []ProcMask{},
			procChances: []DynamicProc{},
		}

		for slot, eq := range character.Equipment {
			if selectEffectId(eq, effectType) != effectID {
				continue
			}

			slotList = append(slotList, proto.ItemSlot(slot))
			rppmConfig.Ilvl = eq.ScalingOptions[int32(eq.UpgradeStep)].Ilvl
			proc := NewRPPMProc(character, rppmConfig)

			mask := procMask
			weaponMask := ProcMaskEmpty
			switch proto.ItemSlot(slot) {
			case proto.ItemSlot_ItemSlotMainHand:
				if eq.RangedWeaponType > 0 {
					weaponMask = ProcMaskRanged
					if mask.Matches(ProcMaskRangedProc) {
						weaponMask |= ProcMaskRangedProc
					}
				} else {
					weaponMask = ProcMaskMeleeMH
					if mask.Matches(ProcMaskMeleeProc) {
						weaponMask |= ProcMaskMeleeProc
					}
				}
			case proto.ItemSlot_ItemSlotOffHand:
				weaponMask = ProcMaskMeleeOH
				if mask.Matches(ProcMaskMeleeProc) {
					weaponMask |= ProcMaskMeleeProc
				}
			}

			// The current proc is attached to a weapon
			// In this case we want to make sure the proc can only proc off this weapon
			// Or spells, so we remove any melee proc masks and only add the ones determined
			if weaponMask != ProcMaskEmpty {
				mask &= ^(ProcMaskMeleeOrMeleeProc | ProcMaskRangedOrRangedProc)
				mask |= weaponMask
			}

			manager.procMasks = append(manager.procMasks, mask)
			manager.procChances = append(manager.procChances, proc)
		}

		return manager
	}

	dpm := builder()
	character.RegisterItemSwapCallback(slotList, func(_ *Simulation, _ proto.ItemSlot) {
		dpm = builder()
	})

	return &dpm
}

func selectEffectId(item Item, effectType EffectType) int32 {
	switch effectType {
	case ItemEffect:
		return item.ID

		// EnchantEffect and WeaponEffect sore some reason both use the enchant ID
	case EnchantEffect, WeaponEffect:
		return item.Enchant.EffectID
	default:
		return 0
	}
}
