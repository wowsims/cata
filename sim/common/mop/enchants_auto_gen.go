package mop

import (
	"github.com/wowsims/mop/sim/core"
 	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllEnchants() {

	// Enchants
	
	// Permanently attaches Lord Blastington's special scope to a ranged weapon, sometimes increasing Agility
	// by 1800 for 10s when dealing damage with ranged attacks.
	// 
	// Attaching this scope to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Lord Blastington's Scope of Doom",
		EnchantID: 4699,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	
	// Permanently attaches a mirrored scope to a ranged weapon, sometimes increases critical strike by 900 for
	// 10s when dealing damage with ranged attacks.
	// 
	// Attaching this scope to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Mirror Scope",
		EnchantID: 4700,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	
	// Embroiders a subtle pattern of light into your cloak, giving you a chance to increase your Intellect by
	// 2000 for 15s when casting a spell.
	// 
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Lightweave Embroidery (Rank 3)",
		EnchantID: 4892,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellDamageProc,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	
	// Embroiders a magical pattern into your cloak, giving you a chance to increase your Spirit by 3000 for
	// 15s when you cast a spell.
	// 
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Darkglow Embroidery (Rank 3)",
		EnchantID: 4893,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
	
	// Embroiders a magical pattern into your cloak, causing your damaging melee and ranged attacks to sometimes
	// increase your attack power by 4000 for 15s.
	// 
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Swordguard Embroidery (Rank 3)",
		EnchantID: 4894,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
		Harmful:   true,
	})
}