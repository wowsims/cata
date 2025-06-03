package mop

import (
	"github.com/wowsims/mop/sim/core"
 	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllEnchants() {

	// Procs
	
	// Increases agility by 1800.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lord Blastington's Scope of Doom",
		EnchantID:   4699,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:  core.OutcomeLanded,
	})
	
	// Increases critical strike by 900.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mirror Scope",
		EnchantID:   4700,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:  core.OutcomeLanded,
	})
	
	// 
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightweave Embroidery",
		EnchantID:   4892,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// 
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Darkglow Embroidery",
		EnchantID:   4893,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})
	
	// 
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Swordguard Embroidery",
		EnchantID:   4894,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:  core.OutcomeLanded,
	})
	
	// 
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightweave Embroidery",
		EnchantID:   5110,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
}