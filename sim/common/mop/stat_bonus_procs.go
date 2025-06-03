package mop

import (
	"github.com/wowsims/mop/sim/core"
 	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllOnUseProcs() {

	// Procs
	
	// Your melee attacks have a chance to grant Blessing of the Celestials, increasing your Strength by 3027 for 15s. ( 20% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Xuen",
		ItemID:   79327,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deliver a melee or ranged critical strike, you have a chance to gain Blessing of the Celestials, increasing your Agility by 3027 for 15s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Xuen",
		ItemID:   79328,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When you cast healing spells, you have a chance to gain Blessing of the Celestials, increasing your Spirit by 3027 for 20s. ( 20% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Chi-Ji",
		ItemID:   79330,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal spell damage, you have a chance to gain Blessing of the Celestials, increasing your Intellect by 3027 for 15s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Yu'lon",
		ItemID:   79331,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Windswept Pages",
		ItemID:   81125,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Empty Fruit Barrel",
		ItemID:   81133,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 critical strike for 30s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Carbonic Carbuncle",
		ItemID:   81138,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal damage you have a chance to gain 2573 critical strike for 30s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vision of the Predator",
		ItemID:   81192,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Iron Protector Talisman",
		ItemID:   81243,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks critical strike your target you have a chance to gain 2573 Agility for 25s. ( 45% chance, 85 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Searing Words",
		ItemID:   81267,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mark of the Catacombs",
		ItemID:   83731,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of the Catacombs",
		ItemID:   83732,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Emblem of the Catacombs",
		ItemID:   83733,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee attacks have a chance to grant 1851 parry for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Medallion of the Catacombs",
		ItemID:   83734,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Symbol of the Catacombs",
		ItemID:   83735,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 spirit for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Compassion",
		ItemID:   83736,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Fidelity",
		ItemID:   83737,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Grace",
		ItemID:   83738,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee attacks have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Patience",
		ItemID:   83739,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Devotion",
		ItemID:   83740,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fearwurm Relic",
		ItemID:   84070,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Charm of Ten Songs",
		ItemID:   84071,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Braid of Ten Songs",
		ItemID:   84072,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee attacks have a chance to grant 1851 parry for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Knot of Ten Songs",
		ItemID:   84073,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fearwurm Badge",
		ItemID:   84074,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Kypari Zar",
		ItemID:   84075,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Kypari Zar",
		ItemID:   84076,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Emblem of Kypari Zar",
		ItemID:   84077,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee attacks have a chance to grant 1851 dodge for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Insignia of Kypari Zar",
		ItemID:   84078,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Badge of Kypari Zar",
		ItemID:   84079,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dreadful Gladiator's Insignia of Conquest",
		ItemID:   84349,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dreadful Gladiator's Insignia of Dominance",
		ItemID:   84489,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dreadful Gladiator's Insignia of Victory",
		ItemID:   84495,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Malevolent Gladiator's Insignia of Conquest",
		ItemID:   84935,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Malevolent Gladiator's Insignia of Victory",
		ItemID:   84937,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Malevolent Gladiator's Insignia of Dominance",
		ItemID:   84941,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Iron Protector Talisman",
		ItemID:   85181,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vial of Dragon's Blood",
		ItemID:   86131,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bottle of Infinite Stars",
		ItemID:   86132,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time you deal periodic damage you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Light of the Cosmos",
		ItemID:   86133,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lei Shen's Final Orders",
		ItemID:   86144,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your spells heal you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Qin-xi's Polarizing Seal",
		ItemID:   86147,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 dodge for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Stuff of Nightmares",
		ItemID:   86323,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your spells heal you have a chance to gain 963 Spirit for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Spirits of the Sun",
		ItemID:   86327,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 critical strike for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Terror in the Mists",
		ItemID:   86332,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Darkmist Vortex",
		ItemID:   86336,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Essence of Terror",
		ItemID:   86388,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vial of Dragon's Blood",
		ItemID:   86790,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bottle of Infinite Stars",
		ItemID:   86791,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time you deal periodic damage you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Light of the Cosmos",
		ItemID:   86792,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lei Shen's Final Orders",
		ItemID:   86802,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your spells heal you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Qin-xi's Polarizing Seal",
		ItemID:   86805,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 dodge for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Stuff of Nightmares",
		ItemID:   86881,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your spells heal you have a chance to gain 963 Spirit for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Spirits of the Sun",
		ItemID:   86885,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 critical strike for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Terror in the Mists",
		ItemID:   86890,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Darkmist Vortex",
		ItemID:   86894,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Essence of Terror",
		ItemID:   86907,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bottle of Infinite Stars",
		ItemID:   87057,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vial of Dragon's Blood",
		ItemID:   87063,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time you deal periodic damage you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Light of the Cosmos",
		ItemID:   87065,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lei Shen's Final Orders",
		ItemID:   87072,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your spells heal you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Qin-xi's Polarizing Seal",
		ItemID:   87075,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 dodge for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Stuff of Nightmares",
		ItemID:   87160,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your spells heal you have a chance to gain 963 Spirit for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Spirits of the Sun",
		ItemID:   87163,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 critical strike for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Terror in the Mists",
		ItemID:   87167,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Darkmist Vortex",
		ItemID:   87172,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Essence of Terror",
		ItemID:   87175,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing spells have a chance to grant 1926 spellpower for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Core of Decency",
		ItemID:   87497,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your harmful spells have a chance to increase your spell power by 2040 for 10s. ( 10% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mithril Wristwatch",
		ItemID:   87572,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your direct healing and heal over time spells have a chance to increase your haste by 2040 for 10s. ( 10% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thousand-Year Pickled Egg",
		ItemID:   87573,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Chance on melee and ranged critical strike to increase your attack power by 4000 for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Coren's Cold Chromium Coaster",
		ItemID:   87574,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:  core.OutcomeCrit,
	})
	
	// You gain an additional 375 critical strike for 10s. This effect stacks up to 3 times.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "The Gloaming Blade",
		ItemID:   88149,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Conquest",
		ItemID:   91104,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Dominance",
		ItemID:   91401,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Victory",
		ItemID:   91415,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Malevolent Gladiator's Insignia of Conquest",
		ItemID:   91457,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Malevolent Gladiator's Insignia of Dominance",
		ItemID:   91754,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Malevolent Gladiator's Insignia of Victory",
		ItemID:   91768,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Dreadful Gladiator's Insignia of Conquest",
		ItemID:   93424,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Dreadful Gladiator's Insignia of Dominance",
		ItemID:   93601,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Dreadful Gladiator's Insignia of Victory",
		ItemID:   93611,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Conquest",
		ItemID:   94356,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Victory",
		ItemID:   94415,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Dominance",
		ItemID:   94482,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 15s. ( 15% chance, 85 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Brutal Talisman of the Shado-Pan Assault",
		ItemID:   94508,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 10s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Volatile Talisman of the Shado-Pan Assault",
		ItemID:   94510,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Talisman of the Shado-Pan Assault",
		ItemID:   94511,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you dodge, you have a 4% chance to gain 963 mastery for 20s. This effect can stack up to 3 times.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Delicate Vial of the Sanguinaire",
		ItemID:   94518,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Primordius' Talisman of Rage",
		ItemID:   94519,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your periodic damage spells have a chance to grant 1926 Intellect for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Breath of the Hydra",
		ItemID:   94521,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 haste for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Talisman of Bloodlust",
		ItemID:   94522,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Agility and summon 3 Voodoo Gnomes for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bad Juju",
		ItemID:   94523,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your critical attacks have a chance to grant you 963 Critical Strike for 20s. This effect can stack up to 3 times. (Approximately 0.72 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Gaze of the Twins",
		ItemID:   94529,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When your spells deal critical damage, you have a chance to gain 1926 Intellect for 10s. (Approximately 0.85 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cha-Ye's Essence of Brilliance",
		ItemID:   94531,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When your attacks hit you have a chance to gain 2573 Agility and summon 3 Voodoo Gnomes for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bad Juju",
		ItemID:   95665,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your periodic damage spells have a chance to grant 1926 Intellect for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Breath of the Hydra",
		ItemID:   95711,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 haste for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Talisman of Bloodlust",
		ItemID:   95748,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Primordius' Talisman of Rage",
		ItemID:   95757,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal critical damage, you have a chance to gain 1926 Intellect for 10s. (Approximately 0.85 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cha-Ye's Essence of Brilliance",
		ItemID:   95772,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When you dodge, you have a 4% chance to gain 963 mastery for 20s. This effect can stack up to 3 times.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Delicate Vial of the Sanguinaire",
		ItemID:   95779,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your critical attacks have a chance to grant you 963 Critical Strike for 20s. This effect can stack up to 3 times. (Approximately 0.72 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Gaze of the Twins",
		ItemID:   95799,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When your attacks hit you have a chance to gain 2573 Agility and summon 3 Voodoo Gnomes for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bad Juju",
		ItemID:   96037,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your periodic damage spells have a chance to grant 1926 Intellect for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Breath of the Hydra",
		ItemID:   96083,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 haste for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Talisman of Bloodlust",
		ItemID:   96120,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Primordius' Talisman of Rage",
		ItemID:   96129,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal critical damage, you have a chance to gain 1926 Intellect for 10s. (Approximately 0.85 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cha-Ye's Essence of Brilliance",
		ItemID:   96144,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When you dodge, you have a 4% chance to gain 963 mastery for 20s. This effect can stack up to 3 times.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Delicate Vial of the Sanguinaire",
		ItemID:   96151,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your critical attacks have a chance to grant you 963 Critical Strike for 20s. This effect can stack up to 3 times. (Approximately 0.72 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Gaze of the Twins",
		ItemID:   96171,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When your attacks hit you have a chance to gain 2573 Agility and summon 3 Voodoo Gnomes for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bad Juju",
		ItemID:   96409,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your periodic damage spells have a chance to grant 1926 Intellect for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Breath of the Hydra",
		ItemID:   96455,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 haste for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Talisman of Bloodlust",
		ItemID:   96492,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Primordius' Talisman of Rage",
		ItemID:   96501,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal critical damage, you have a chance to gain 1926 Intellect for 10s. (Approximately 0.85 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cha-Ye's Essence of Brilliance",
		ItemID:   96516,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When you dodge, you have a 4% chance to gain 963 mastery for 20s. This effect can stack up to 3 times.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Delicate Vial of the Sanguinaire",
		ItemID:   96523,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your critical attacks have a chance to grant you 963 Critical Strike for 20s. This effect can stack up to 3 times. (Approximately 0.72 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Gaze of the Twins",
		ItemID:   96543,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When your attacks hit you have a chance to gain 2573 Agility and summon 3 Voodoo Gnomes for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bad Juju",
		ItemID:   96781,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your periodic damage spells have a chance to grant 1926 Intellect for 10s. (Approximately 1.10 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Breath of the Hydra",
		ItemID:   96827,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 haste for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Talisman of Bloodlust",
		ItemID:   96864,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 10s. This effect can stack up to 5 times. (Approximately 3.50 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Primordius' Talisman of Rage",
		ItemID:   96873,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal critical damage, you have a chance to gain 1926 Intellect for 10s. (Approximately 0.85 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cha-Ye's Essence of Brilliance",
		ItemID:   96888,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// When you dodge, you have a 4% chance to gain 963 mastery for 20s. This effect can stack up to 3 times.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Delicate Vial of the Sanguinaire",
		ItemID:   96895,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your critical attacks have a chance to grant you 963 Critical Strike for 20s. This effect can stack up to 3 times. (Approximately 0.72 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Gaze of the Twins",
		ItemID:   96915,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Empty Fruit Barrel",
		ItemID:   97304,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time you deal periodic damage you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Light of the Cosmos",
		ItemID:   98019,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Essence of Terror",
		ItemID:   98020,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal spell damage, you have a chance to gain Blessing of the Celestials, increasing your Intellect by 3027 for 15s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Yu'lon",
		ItemID:   98049,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time you deal periodic damage you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Light of the Cosmos",
		ItemID:   98050,
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal spell damage, you have a chance to gain Blessing of the Celestials, increasing your Intellect by 3027 for 15s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Yu'lon",
		ItemID:   98075,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Essence of Terror",
		ItemID:   98076,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Malevolent Gladiator's Insignia of Conquest",
		ItemID:   98760,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Malevolent Gladiator's Insignia of Dominance",
		ItemID:   98911,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Malevolent Gladiator's Insignia of Victory",
		ItemID:   98917,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Conquest",
		ItemID:   99777,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Dominance",
		ItemID:   99938,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Victory",
		ItemID:   99948,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Conquest",
		ItemID:   100026,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Victory",
		ItemID:   100085,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tyrannical Gladiator's Insignia of Dominance",
		ItemID:   100152,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Conquest",
		ItemID:   100200,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Dominance",
		ItemID:   100491,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Victory",
		ItemID:   100505,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Conquest",
		ItemID:   100586,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Victory",
		ItemID:   100645,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Dominance",
		ItemID:   100712,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Stone of Battle",
		ItemID:   100990,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Idol of Battle",
		ItemID:   100991,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Defender Idol",
		ItemID:   100999,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Defender Stone",
		ItemID:   101002,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Idol of Rage",
		ItemID:   101009,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Stone of Rage",
		ItemID:   101012,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Idol of Destruction",
		ItemID:   101023,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Stone of Destruction",
		ItemID:   101026,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Stone of Wisdom",
		ItemID:   101041,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Trailseeker Idol of Rage",
		ItemID:   101054,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Trailseeker Stone of Rage",
		ItemID:   101057,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mountainsage Idol of Destruction",
		ItemID:   101069,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mountainsage Stone of Destruction",
		ItemID:   101072,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Defender Stone",
		ItemID:   101087,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Defender Idol",
		ItemID:   101089,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Stone of Wisdom",
		ItemID:   101107,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Idol of Rage",
		ItemID:   101113,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Stone of Rage",
		ItemID:   101117,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Stone of Wisdom",
		ItemID:   101138,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Stone of Battle",
		ItemID:   101151,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Idol of Battle",
		ItemID:   101152,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Defender Idol",
		ItemID:   101160,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Defender Stone",
		ItemID:   101163,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Communal Idol of Destruction",
		ItemID:   101168,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Communal Stone of Destruction",
		ItemID:   101171,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Communal Stone of Wisdom",
		ItemID:   101183,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightdrinker Idol of Rage",
		ItemID:   101200,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightdrinker Stone of Rage",
		ItemID:   101203,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Idol of Rage",
		ItemID:   101217,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Stone of Rage",
		ItemID:   101220,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Idol of Destruction",
		ItemID:   101222,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Stone of Destruction",
		ItemID:   101225,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Stone of Wisdom",
		ItemID:   101250,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Felsoul Idol of Destruction",
		ItemID:   101263,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Felsoul Stone of Destruction",
		ItemID:   101266,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Stone of Battle",
		ItemID:   101294,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Idol of Battle",
		ItemID:   101295,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Defender Idol",
		ItemID:   101303,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Defender Stone",
		ItemID:   101306,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Agility-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Assurance of Consequence",
		ItemID:   102292,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a chance to grant 1 Intellect for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Purified Bindings of Immerseus",
		ItemID:   102293,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.1% chance to trigger Multistrike, which causes instant additional healing to your target equal to 33% of the original healing done.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Nazgrim's Burnished Insignia",
		ItemID:   102294,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fusion-Fire Core",
		ItemID:   102295,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Evil Eye of Galakras",
		ItemID:   102298,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prismatic Prison of Pride",
		ItemID:   102299,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Kardris' Toxic Totem",
		ItemID:   102300,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Haromm's Talisman",
		ItemID:   102301,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Rampage",
		ItemID:   102302,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Frenzied Crystal of Rage",
		ItemID:   102303,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.01% chance to Cleave, dealing the same healing to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Acid-Grooved Tooth",
		ItemID:   102304,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Tail Tip",
		ItemID:   102305,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// Your heals have a chance to grant you 19260 Spirit for 10s. Every 0.5 sec, this effect is reduced by 963 Spirit. (Approximately 0.92 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dysmorphic Samophlange of Discontinuity",
		ItemID:   102309,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your melee and ranged attacks have a chance to grant you 19260 Agility for 10s. Every 0.5 sec this effect decrements by 963 Agility. (Approximately 1.00 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ticking Ebon Detonator",
		ItemID:   102311,
		Callback: core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "5.4 Raid - Normal - Siege of Orgrimmar - Boss X Loot X - Agi DPS Trinket (5)",
		ItemID:   102312,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal damage you have a chance to gain 2573 critical strike for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "5.4 Raid - Normal - Siege of Orgrimmar - Boss X Loot X - Int Hit Trinket (5)",
		ItemID:   102313,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your melee attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "5.4 Raid - Normal - Siege of Orgrimmar - Boss X Loot X - Str DPS Trinket (5)",
		ItemID:   102315,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prideful Gladiator's Insignia of Conquest",
		ItemID:   102643,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prideful Gladiator's Insignia of Victory",
		ItemID:   102699,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prideful Gladiator's Insignia of Dominance",
		ItemID:   102766,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Conquest",
		ItemID:   102840,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Victory",
		ItemID:   102896,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Dominance",
		ItemID:   102963,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Conquest",
		ItemID:   103150,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Dominance",
		ItemID:   103309,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grievous Gladiator's Insignia of Victory",
		ItemID:   103319,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prideful Gladiator's Insignia of Conquest",
		ItemID:   103347,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prideful Gladiator's Insignia of Dominance",
		ItemID:   103506,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prideful Gladiator's Insignia of Victory",
		ItemID:   103516,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// TODO: Overwrite me
	// Teleport yourself to the Timeless Isle.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Time-Lost Artifact",
		ItemID:   103678,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Discipline of Xuen",
		ItemID:   103686,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal damage you have a chance to gain 2573 critical strike for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Yu'lon's Bite",
		ItemID:   103687,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your melee attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Alacrity of Xuen",
		ItemID:   103689,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Discipline of Xuen",
		ItemID:   103986,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// When your spells deal damage you have a chance to gain 2573 critical strike for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Yu'lon's Bite",
		ItemID:   103987,
		Callback: core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// Each time your melee attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Alacrity of Xuen",
		ItemID:   103989,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// TODO: Overwrite me
	// Your attacks have a chance to grant 1 Intellect for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Purified Bindings of Immerseus",
		ItemID:   104426,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fusion-Fire Core",
		ItemID:   104463,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Agility-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Assurance of Consequence",
		ItemID:   104476,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prismatic Prison of Pride",
		ItemID:   104478,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Evil Eye of Galakras",
		ItemID:   104495,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Haromm's Talisman",
		ItemID:   104531,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Kardris' Toxic Totem",
		ItemID:   104544,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.1% chance to trigger Multistrike, which causes instant additional healing to your target equal to 33% of the original healing done.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Nazgrim's Burnished Insignia",
		ItemID:   104553,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Frenzied Crystal of Rage",
		ItemID:   104576,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Rampage",
		ItemID:   104584,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.01% chance to Cleave, dealing the same healing to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Acid-Grooved Tooth",
		ItemID:   104611,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Tail Tip",
		ItemID:   104613,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// Your melee and ranged attacks have a chance to grant you 19260 Agility for 10s. Every 0.5 sec this effect decrements by 963 Agility. (Approximately 1.00 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ticking Ebon Detonator",
		ItemID:   104616,
		Callback: core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your heals have a chance to grant you 19260 Spirit for 10s. Every 0.5 sec, this effect is reduced by 963 Spirit. (Approximately 0.92 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dysmorphic Samophlange of Discontinuity",
		ItemID:   104619,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// TODO: Overwrite me
	// Your attacks have a chance to grant 1 Intellect for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Purified Bindings of Immerseus",
		ItemID:   104675,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fusion-Fire Core",
		ItemID:   104712,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Agility-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Assurance of Consequence",
		ItemID:   104725,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prismatic Prison of Pride",
		ItemID:   104727,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Evil Eye of Galakras",
		ItemID:   104744,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Haromm's Talisman",
		ItemID:   104780,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Kardris' Toxic Totem",
		ItemID:   104793,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.1% chance to trigger Multistrike, which causes instant additional healing to your target equal to 33% of the original healing done.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Nazgrim's Burnished Insignia",
		ItemID:   104802,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Frenzied Crystal of Rage",
		ItemID:   104825,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Rampage",
		ItemID:   104833,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.01% chance to Cleave, dealing the same healing to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Acid-Grooved Tooth",
		ItemID:   104860,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Tail Tip",
		ItemID:   104862,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// Your melee and ranged attacks have a chance to grant you 19260 Agility for 10s. Every 0.5 sec this effect decrements by 963 Agility. (Approximately 1.00 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ticking Ebon Detonator",
		ItemID:   104865,
		Callback: core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your heals have a chance to grant you 19260 Spirit for 10s. Every 0.5 sec, this effect is reduced by 963 Spirit. (Approximately 0.92 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dysmorphic Samophlange of Discontinuity",
		ItemID:   104868,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// TODO: Overwrite me
	// Your attacks have a chance to grant 1 Intellect for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Purified Bindings of Immerseus",
		ItemID:   104924,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fusion-Fire Core",
		ItemID:   104961,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Agility-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Assurance of Consequence",
		ItemID:   104974,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prismatic Prison of Pride",
		ItemID:   104976,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Evil Eye of Galakras",
		ItemID:   104993,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Haromm's Talisman",
		ItemID:   105029,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Kardris' Toxic Totem",
		ItemID:   105042,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.1% chance to trigger Multistrike, which causes instant additional healing to your target equal to 33% of the original healing done.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Nazgrim's Burnished Insignia",
		ItemID:   105051,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Frenzied Crystal of Rage",
		ItemID:   105074,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Rampage",
		ItemID:   105082,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.01% chance to Cleave, dealing the same healing to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Acid-Grooved Tooth",
		ItemID:   105109,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Tail Tip",
		ItemID:   105111,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// Your melee and ranged attacks have a chance to grant you 19260 Agility for 10s. Every 0.5 sec this effect decrements by 963 Agility. (Approximately 1.00 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ticking Ebon Detonator",
		ItemID:   105114,
		Callback: core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your heals have a chance to grant you 19260 Spirit for 10s. Every 0.5 sec, this effect is reduced by 963 Spirit. (Approximately 0.92 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dysmorphic Samophlange of Discontinuity",
		ItemID:   105117,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// TODO: Overwrite me
	// Your attacks have a chance to grant 1 Intellect for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Purified Bindings of Immerseus",
		ItemID:   105173,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fusion-Fire Core",
		ItemID:   105210,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Agility-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Assurance of Consequence",
		ItemID:   105223,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prismatic Prison of Pride",
		ItemID:   105225,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Evil Eye of Galakras",
		ItemID:   105242,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Haromm's Talisman",
		ItemID:   105278,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Kardris' Toxic Totem",
		ItemID:   105291,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.1% chance to trigger Multistrike, which causes instant additional healing to your target equal to 33% of the original healing done.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Nazgrim's Burnished Insignia",
		ItemID:   105300,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Frenzied Crystal of Rage",
		ItemID:   105323,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Rampage",
		ItemID:   105331,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.01% chance to Cleave, dealing the same healing to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Acid-Grooved Tooth",
		ItemID:   105358,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Tail Tip",
		ItemID:   105360,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// Your melee and ranged attacks have a chance to grant you 19260 Agility for 10s. Every 0.5 sec this effect decrements by 963 Agility. (Approximately 1.00 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ticking Ebon Detonator",
		ItemID:   105363,
		Callback: core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your heals have a chance to grant you 19260 Spirit for 10s. Every 0.5 sec, this effect is reduced by 963 Spirit. (Approximately 0.92 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dysmorphic Samophlange of Discontinuity",
		ItemID:   105366,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
	
	// TODO: Overwrite me
	// Your attacks have a chance to grant 1 Intellect for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Purified Bindings of Immerseus",
		ItemID:   105422,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fusion-Fire Core",
		ItemID:   105459,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Agility-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Assurance of Consequence",
		ItemID:   105472,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prismatic Prison of Pride",
		ItemID:   105474,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based damage roles only.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Evil Eye of Galakras",
		ItemID:   105491,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Haromm's Talisman",
		ItemID:   105527,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your target equal to 33% of the original damage dealt.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Kardris' Toxic Totem",
		ItemID:   105540,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.1% chance to trigger Multistrike, which causes instant additional healing to your target equal to 33% of the original healing done.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Nazgrim's Burnished Insignia",
		ItemID:   105549,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Frenzied Crystal of Rage",
		ItemID:   105572,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Rampage",
		ItemID:   105580,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Your heals have a 0.01% chance to Cleave, dealing the same healing to up to 5 other nearby targets.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Acid-Grooved Tooth",
		ItemID:   105607,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// TODO: Overwrite me
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thok's Tail Tip",
		ItemID:   105609,
		Callback: core.CallbackEmpty,
		ProcMask: core.ProcMaskEmpty,
		Outcome:  core.OutcomeEmpty,
	})
	
	// Your melee and ranged attacks have a chance to grant you 19260 Agility for 10s. Every 0.5 sec this effect decrements by 963 Agility. (Approximately 1.00 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ticking Ebon Detonator",
		ItemID:   105612,
		Callback: core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})
	
	// Your heals have a chance to grant you 19260 Spirit for 10s. Every 0.5 sec, this effect is reduced by 963 Spirit. (Approximately 0.92 procs per minute)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dysmorphic Samophlange of Discontinuity",
		ItemID:   105615,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
	})
}