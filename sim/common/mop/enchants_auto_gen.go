package mop

import (
	"github.com/wowsims/mop/sim/core"
 	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllEnchants() {

	// Enchants
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your critical strike, haste, or mastery by 1500
	// for 12s when dealing damage or healing with spells and melee attacks.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Windsong",
		EnchantID: 4441,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your Intellect by 0 when healing or dealing
	// damage with spells. If less than 25% of your mana remains when the effect is triggered, your Spirit will
	// also increase by 0.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Jade Spirit",
		EnchantID: 4442,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your Strength or Agility by 0 when dealing melee
	// damage. Your highest stat is always chosen.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Dancing Steel",
		EnchantID: 4444,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to make your damaging melee strikes sometimes activate a Mogu protection
	// spell, absorbing up to 0 damage.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Colossus",
		EnchantID: 4445,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your dodge by 1650 for 7s when dealing melee
	// damage.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - River's Song",
		EnchantID: 4446,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your Intellect by 165 when healing or dealing
	// damage with spells. If less than 25% of your mana remains when the effect is triggered, your Spirit will
	// also increase by 75. Requires a level 372 or higher item.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Jade Spirit",
		EnchantID: 5062,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your Intellect by 165 when healing or dealing
	// damage with spells. If less than 25% of your mana remains when the effect is triggered, your Spirit will
	// also increase by 75. Requires a level 372 or higher item.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Jade Spirit",
		EnchantID: 5098,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your Intellect by 0 when healing or dealing
	// damage with spells. If less than 25% of your mana remains when the effect is triggered, your Spirit will
	// also increase by 0.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Spirit of Conquest",
		EnchantID: 5124,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	//
	// Permanently enchants a melee weapon to sometimes increase your Strength or Agility by 0 when dealing melee
	// damage. Your highest stat is always chosen.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Enchant Weapon - Bloody Dancing Steel",
		EnchantID: 5125,
		Callback:  core.CallbackEmpty,
		ProcMask:  core.ProcMaskEmpty,
		Outcome:   core.OutcomeEmpty,
	})
	
	// Permanently attaches Lord Blastington's special scope to a ranged weapon, sometimes increasing Agility
	// by 1800 for 10s when dealing damage with ranged attacks.
	// 
	// Attaching this scope to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lord Blastington's Scope of Doom",
		EnchantID: 4699,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
	})
	
	// Permanently attaches a mirrored scope to a ranged weapon, sometimes increases critical strike by 900 for
	// 10s when dealing damage with ranged attacks.
	// 
	// Attaching this scope to a ranged weapon causes it to become soulbound.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mirror Scope",
		EnchantID: 4700,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
	})
	
	// Embroiders a subtle pattern of light into your cloak, giving you a chance to increase your Intellect by
	// 2000 for 15s when casting a spell.
	// 
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightweave Embroidery",
		EnchantID: 4892,
		Callback:  core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:   core.OutcomeLanded,
	})
	
	// Embroiders a magical pattern into your cloak, giving you a chance to increase your Spirit by 3000 for
	// 15s when you cast a spell.
	// 
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Darkglow Embroidery",
		EnchantID: 4893,
		Callback:  core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage,
		Outcome:   core.OutcomeLanded,
	})
	
	// Embroiders a magical pattern into your cloak, causing your damaging melee and ranged attacks to sometimes
	// increase your attack power by 4000 for 15s.
	// 
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Swordguard Embroidery",
		EnchantID: 4894,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:   core.OutcomeLanded,
	})
	
	// Embroiders a subtle pattern of light into your cloak, giving you a chance to increase your Intellect by
	// 2000 for 15s when casting a spell.
	// 
	// Embroidering your cloak will cause it to become soulbound and requires the Tailoring profession to remain
	// active.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightweave Embroidery",
		EnchantID: 5110,
		Callback:  core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:   core.OutcomeLanded,
	})
}