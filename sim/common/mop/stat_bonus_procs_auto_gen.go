package mop

import (
	"github.com/wowsims/mop/sim/core"
 	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllProcs() {

	// Procs
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Agility-based
	// damage roles only.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102292, ItemName: "Assurance of Consequence (N)"},
	//	{ItemID: 104476, ItemName: "Assurance of Consequence (H)"},
	//	{ItemID: 104725, ItemName: "Assurance of Consequence (Flexible)"},
	//	{ItemID: 104974, ItemName: "Assurance of Consequence (LFR) (Celestial)"},
	//	{ItemID: 105223, ItemName: "Assurance of Consequence (Warforged)"},
	//	{ItemID: 105472, ItemName: "Assurance of Consequence (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your attacks have a chance to grant 1 Intellect for 20s. ( 15% chance, 115 sec cooldown)
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102293, ItemName: "Purified Bindings of Immerseus (N)"},
	//	{ItemID: 104426, ItemName: "Purified Bindings of Immerseus (H)"},
	//	{ItemID: 104675, ItemName: "Purified Bindings of Immerseus (Flexible)"},
	//	{ItemID: 104924, ItemName: "Purified Bindings of Immerseus (LFR) (Celestial)"},
	//	{ItemID: 105173, ItemName: "Purified Bindings of Immerseus (Warforged)"},
	//	{ItemID: 105422, ItemName: "Purified Bindings of Immerseus (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your heals have a 0.1% chance to trigger Multistrike, which causes instant additional healing to your
	// target equal to 33% of the original healing done.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
	//	ProcMask: core.ProcMaskSpellHealing,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102294, ItemName: "Nazgrim's Burnished Insignia (N)"},
	//	{ItemID: 104553, ItemName: "Nazgrim's Burnished Insignia (H)"},
	//	{ItemID: 104802, ItemName: "Nazgrim's Burnished Insignia (Flexible)"},
	//	{ItemID: 105051, ItemName: "Nazgrim's Burnished Insignia (LFR) (Celestial)"},
	//	{ItemID: 105300, ItemName: "Nazgrim's Burnished Insignia (Warforged)"},
	//	{ItemID: 105549, ItemName: "Nazgrim's Burnished Insignia (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102295, ItemName: "Fusion-Fire Core (N)"},
	//	{ItemID: 104463, ItemName: "Fusion-Fire Core (H)"},
	//	{ItemID: 104712, ItemName: "Fusion-Fire Core (Flexible)"},
	//	{ItemID: 104961, ItemName: "Fusion-Fire Core (LFR) (Celestial)"},
	//	{ItemID: 105210, ItemName: "Fusion-Fire Core (Warforged)"},
	//	{ItemID: 105459, ItemName: "Fusion-Fire Core (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based
	// damage roles only.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102298, ItemName: "Evil Eye of Galakras (N)"},
	//	{ItemID: 104495, ItemName: "Evil Eye of Galakras (H)"},
	//	{ItemID: 104744, ItemName: "Evil Eye of Galakras (Flexible)"},
	//	{ItemID: 104993, ItemName: "Evil Eye of Galakras (LFR) (Celestial)"},
	//	{ItemID: 105242, ItemName: "Evil Eye of Galakras (Warforged)"},
	//	{ItemID: 105491, ItemName: "Evil Eye of Galakras (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
	//	ProcMask: core.ProcMaskSpellHealing,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102299, ItemName: "Prismatic Prison of Pride (N)"},
	//	{ItemID: 104478, ItemName: "Prismatic Prison of Pride (H)"},
	//	{ItemID: 104727, ItemName: "Prismatic Prison of Pride (Flexible)"},
	//	{ItemID: 104976, ItemName: "Prismatic Prison of Pride (LFR) (Celestial)"},
	//	{ItemID: 105225, ItemName: "Prismatic Prison of Pride (Warforged)"},
	//	{ItemID: 105474, ItemName: "Prismatic Prison of Pride (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your
	// target equal to 33% of the original damage dealt.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102300, ItemName: "Kardris' Toxic Totem (N)"},
	//	{ItemID: 104544, ItemName: "Kardris' Toxic Totem (H)"},
	//	{ItemID: 104793, ItemName: "Kardris' Toxic Totem (Flexible)"},
	//	{ItemID: 105042, ItemName: "Kardris' Toxic Totem (LFR) (Celestial)"},
	//	{ItemID: 105291, ItemName: "Kardris' Toxic Totem (Warforged)"},
	//	{ItemID: 105540, ItemName: "Kardris' Toxic Totem (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your attacks have a 0.1% chance to trigger Multistrike, which deals instant additional damage to your
	// target equal to 33% of the original damage dealt.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102301, ItemName: "Haromm's Talisman (N)"},
	//	{ItemID: 104531, ItemName: "Haromm's Talisman (H)"},
	//	{ItemID: 104780, ItemName: "Haromm's Talisman (Flexible)"},
	//	{ItemID: 105029, ItemName: "Haromm's Talisman (LFR) (Celestial)"},
	//	{ItemID: 105278, ItemName: "Haromm's Talisman (Warforged)"},
	//	{ItemID: 105527, ItemName: "Haromm's Talisman (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102302, ItemName: "Sigil of Rampage (N)"},
	//	{ItemID: 104584, ItemName: "Sigil of Rampage (H)"},
	//	{ItemID: 104833, ItemName: "Sigil of Rampage (Flexible)"},
	//	{ItemID: 105082, ItemName: "Sigil of Rampage (LFR) (Celestial)"},
	//	{ItemID: 105331, ItemName: "Sigil of Rampage (Warforged)"},
	//	{ItemID: 105580, ItemName: "Sigil of Rampage (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your attacks have a 0.01% chance to Cleave, dealing the same damage to up to 5 other nearby targets.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102303, ItemName: "Frenzied Crystal of Rage (N)"},
	//	{ItemID: 104576, ItemName: "Frenzied Crystal of Rage (H)"},
	//	{ItemID: 104825, ItemName: "Frenzied Crystal of Rage (Flexible)"},
	//	{ItemID: 105074, ItemName: "Frenzied Crystal of Rage (LFR) (Celestial)"},
	//	{ItemID: 105323, ItemName: "Frenzied Crystal of Rage (Warforged)"},
	//	{ItemID: 105572, ItemName: "Frenzied Crystal of Rage (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your heals have a 0.01% chance to Cleave, dealing the same healing to up to 5 other nearby targets.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
	//	ProcMask: core.ProcMaskSpellHealing,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102304, ItemName: "Thok's Acid-Grooved Tooth (N)"},
	//	{ItemID: 104611, ItemName: "Thok's Acid-Grooved Tooth (H)"},
	//	{ItemID: 104860, ItemName: "Thok's Acid-Grooved Tooth (Flexible)"},
	//	{ItemID: 105109, ItemName: "Thok's Acid-Grooved Tooth (LFR) (Celestial)"},
	//	{ItemID: 105358, ItemName: "Thok's Acid-Grooved Tooth (Warforged)"},
	//	{ItemID: 105607, ItemName: "Thok's Acid-Grooved Tooth (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Amplifies your Critical Strike damage and healing, Haste, Mastery, and Spirit by 1%.
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102305, ItemName: "Thok's Tail Tip (N)"},
	//	{ItemID: 104613, ItemName: "Thok's Tail Tip (H)"},
	//	{ItemID: 104862, ItemName: "Thok's Tail Tip (Flexible)"},
	//	{ItemID: 105111, ItemName: "Thok's Tail Tip (LFR) (Celestial)"},
	//	{ItemID: 105360, ItemName: "Thok's Tail Tip (Warforged)"},
	//	{ItemID: 105609, ItemName: "Thok's Tail Tip (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your heals have a chance to grant you 19260 Spirit for 10s. Every 0.5 sec, this effect is reduced by 963
	// Spirit. (Approximately 0.92 procs per minute)
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
	//	ProcMask: core.ProcMaskSpellHealing,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102309, ItemName: "Dysmorphic Samophlange of Discontinuity (N)"},
	//	{ItemID: 104619, ItemName: "Dysmorphic Samophlange of Discontinuity (H)"},
	//	{ItemID: 104868, ItemName: "Dysmorphic Samophlange of Discontinuity (Flexible)"},
	//	{ItemID: 105117, ItemName: "Dysmorphic Samophlange of Discontinuity (LFR) (Celestial)"},
	//	{ItemID: 105366, ItemName: "Dysmorphic Samophlange of Discontinuity (Warforged)"},
	//	{ItemID: 105615, ItemName: "Dysmorphic Samophlange of Discontinuity (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Your melee and ranged attacks have a chance to grant you 19260 Agility for 10s. Every 0.5 sec this effect
	// decrements by 963 Agility. (Approximately 1.00 procs per minute)
	// shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true,
	// }, []shared.ItemVariant{
	//	{ItemID: 102311, ItemName: "Ticking Ebon Detonator (N)"},
	//	{ItemID: 104616, ItemName: "Ticking Ebon Detonator (H)"},
	//	{ItemID: 104865, ItemName: "Ticking Ebon Detonator (Flexible)"},
	//	{ItemID: 105114, ItemName: "Ticking Ebon Detonator (LFR) (Celestial)"},
	//	{ItemID: 105363, ItemName: "Ticking Ebon Detonator (Warforged)"},
	//	{ItemID: 105612, ItemName: "Ticking Ebon Detonator (Heroic Warforged)"},
	// })
	
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	// Teleport yourself to the Timeless Isle.
	// shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
	//	Name:     "Time-Lost Artifact",
	//	ItemID:   103678,
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
	//	Outcome:  core.OutcomeLanded,
	//	Harmful:  true
	// })
	
	// Your melee attacks have a chance to grant Blessing of the Celestials, increasing your Strength by 3027
	// for 15s. ( 20% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Xuen - Strength",
		ItemID:   79327,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deliver a melee or ranged critical strike, you have a chance to gain Blessing of the Celestials,
	// increasing your Agility by 3027 for 15s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Xuen - Agility",
		ItemID:   79328,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeCrit,
		Harmful:  true,
	})
	
	// When you cast healing spells, you have a chance to gain Blessing of the Celestials, increasing your Spirit
	// by 3027 for 20s. ( 20% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Chi-Ji",
		ItemID:   79330,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal spell damage, you have a chance to gain Blessing of the Celestials, increasing your Intellect
	// by 3027 for 15s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Yu'lon",
		ItemID:   79331,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Windswept Pages (H)",
		ItemID:   81125,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Empty Fruit Barrel (H)",
		ItemID:   81133,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 critical strike for 30s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Carbonic Carbuncle (H)",
		ItemID:   81138,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your spells deal damage you have a chance to gain 2573 critical strike for 30s. ( 15% chance, 115
	// sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vision of the Predator (H)",
		ItemID:   81192,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 81243, ItemName: "Iron Protector Talisman (H)"},
		{ItemID: 85181, ItemName: "Iron Protector Talisman (N)"},
	})
	
	// When your attacks critical strike your target you have a chance to gain 2573 Agility for 25s. ( 45% chance,
	// 85 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Searing Words (H)",
		ItemID:   81267,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
		Harmful:  true,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mark of the Catacombs",
		ItemID:   83731,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of the Catacombs",
		ItemID:   83732,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Emblem of the Catacombs",
		ItemID:   83733,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee attacks have a chance to grant 1851 parry for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Medallion of the Catacombs",
		ItemID:   83734,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Symbol of the Catacombs",
		ItemID:   83735,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 spirit for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Compassion",
		ItemID:   83736,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Fidelity",
		ItemID:   83737,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Grace",
		ItemID:   83738,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee attacks have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Patience",
		ItemID:   83739,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Devotion",
		ItemID:   83740,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fearwurm Relic",
		ItemID:   84070,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Charm of Ten Songs",
		ItemID:   84071,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Braid of Ten Songs",
		ItemID:   84072,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee attacks have a chance to grant 1851 parry for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Knot of Ten Songs",
		ItemID:   84073,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fearwurm Badge",
		ItemID:   84074,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing and damaging spells have a chance to grant 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Relic of Kypari Zar",
		ItemID:   84075,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1851 mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sigil of Kypari Zar",
		ItemID:   84076,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 haste for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Emblem of Kypari Zar",
		ItemID:   84077,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee attacks have a chance to grant 1851 dodge for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Insignia of Kypari Zar",
		ItemID:   84078,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your melee and ranged attacks have a chance to grant 1851 critical strike for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Badge of Kypari Zar",
		ItemID:   84079,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskMeleeProc | core.ProcMaskRangedProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dreadful Gladiator's Insignia of Conquest (Season 12)",
		ItemID:   84349,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dreadful Gladiator's Insignia of Dominance (Season 12)",
		ItemID:   84489,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dreadful Gladiator's Insignia of Victory (Season 12)",
		ItemID:   84495,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 84935, ItemName: "Malevolent Gladiator's Insignia of Conquest (Season 12)"},
		{ItemID: 91457, ItemName: "Malevolent Gladiator's Insignia of Conquest (Season 13)"},
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 84937, ItemName: "Malevolent Gladiator's Insignia of Victory (Season 12)"},
		{ItemID: 91768, ItemName: "Malevolent Gladiator's Insignia of Victory (Season 13)"},
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 84941, ItemName: "Malevolent Gladiator's Insignia of Dominance (Season 12)"},
		{ItemID: 91754, ItemName: "Malevolent Gladiator's Insignia of Dominance (Season 13)"},
	})
	
	// Your attacks have a chance to grant you 963 dodge for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86131, ItemName: "Vial of Dragon's Blood (N)"},
		{ItemID: 86790, ItemName: "Vial of Dragon's Blood (LFR) (Celestial)"},
		{ItemID: 87063, ItemName: "Vial of Dragon's Blood (H)"},
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86132, ItemName: "Bottle of Infinite Stars (N)"},
		{ItemID: 86791, ItemName: "Bottle of Infinite Stars (LFR) (Celestial)"},
		{ItemID: 87057, ItemName: "Bottle of Infinite Stars (H)"},
	})
	
	// Each time you deal periodic damage you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec
	// cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86133, ItemName: "Light of the Cosmos (N)"},
		{ItemID: 86792, ItemName: "Light of the Cosmos (LFR) (Celestial)"},
		{ItemID: 87065, ItemName: "Light of the Cosmos (H)"},
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86144, ItemName: "Lei Shen's Final Orders (N)"},
		{ItemID: 86802, ItemName: "Lei Shen's Final Orders (LFR) (Celestial)"},
		{ItemID: 87072, ItemName: "Lei Shen's Final Orders (H)"},
	})
	
	// Each time your spells heal you have a chance to gain 963 Intellect for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86147, ItemName: "Qin-xi's Polarizing Seal (N)"},
		{ItemID: 86805, ItemName: "Qin-xi's Polarizing Seal (LFR) (Celestial)"},
		{ItemID: 87075, ItemName: "Qin-xi's Polarizing Seal (H)"},
	})
	
	// Each time your attacks hit, you have a chance to gain 963 dodge for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86323, ItemName: "Stuff of Nightmares (N)"},
		{ItemID: 86881, ItemName: "Stuff of Nightmares (LFR) (Celestial)"},
		{ItemID: 87160, ItemName: "Stuff of Nightmares (H)"},
	})
	
	// Each time your spells heal you have a chance to gain 963 Spirit for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86327, ItemName: "Spirits of the Sun (N)"},
		{ItemID: 86885, ItemName: "Spirits of the Sun (LFR) (Celestial)"},
		{ItemID: 87163, ItemName: "Spirits of the Sun (H)"},
	})
	
	// Each time your attacks hit, you have a chance to gain 963 critical strike for 20s. ( 15% chance, 115 sec
	// cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86332, ItemName: "Terror in the Mists (N)"},
		{ItemID: 86890, ItemName: "Terror in the Mists (LFR) (Celestial)"},
		{ItemID: 87167, ItemName: "Terror in the Mists (H)"},
	})
	
	// Each time your attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86336, ItemName: "Darkmist Vortex (N)"},
		{ItemID: 86894, ItemName: "Darkmist Vortex (LFR) (Celestial)"},
		{ItemID: 87172, ItemName: "Darkmist Vortex (H)"},
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec
	// cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 86388, ItemName: "Essence of Terror (N)"},
		{ItemID: 86907, ItemName: "Essence of Terror (LFR) (Celestial)"},
		{ItemID: 87175, ItemName: "Essence of Terror (H)"},
	})
	
	// Your healing spells have a chance to grant 1926 spellpower for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Core of Decency",
		ItemID:   87497,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your harmful spells have a chance to increase your spell power by 2040 for 10s. ( 10% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mithril Wristwatch",
		ItemID:   87572,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your direct healing and heal over time spells have a chance to increase your haste by 2040 for 10s. (
	// 10% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Thousand-Year Pickled Egg",
		ItemID:   87573,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Chance on melee and ranged critical strike to increase your attack power by 4000 for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Coren's Cold Chromium Coaster",
		ItemID:   87574,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial,
		Outcome:  core.OutcomeCrit,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 91104, ItemName: "Tyrannical Gladiator's Insignia of Conquest (Season 13) (Alliance)"},
		{ItemID: 94356, ItemName: "Tyrannical Gladiator's Insignia of Conquest (Season 13) (Horde)"},
		{ItemID: 99777, ItemName: "Tyrannical Gladiator's Insignia of Conquest (Season 14) (Alliance)"},
		{ItemID: 100026, ItemName: "Tyrannical Gladiator's Insignia of Conquest (Season 14) (Horde)"},
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 91401, ItemName: "Tyrannical Gladiator's Insignia of Dominance (Season 13) (Alliance)"},
		{ItemID: 94482, ItemName: "Tyrannical Gladiator's Insignia of Dominance (Season 13) (Horde)"},
		{ItemID: 99938, ItemName: "Tyrannical Gladiator's Insignia of Dominance (Season 14) (Alliance)"},
		{ItemID: 100152, ItemName: "Tyrannical Gladiator's Insignia of Dominance (Season 14) (Horde)"},
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 91415, ItemName: "Tyrannical Gladiator's Insignia of Victory (Season 13) (Alliance)"},
		{ItemID: 94415, ItemName: "Tyrannical Gladiator's Insignia of Victory (Season 13) (Horde)"},
		{ItemID: 99948, ItemName: "Tyrannical Gladiator's Insignia of Victory (Season 14) (Alliance)"},
		{ItemID: 100085, ItemName: "Tyrannical Gladiator's Insignia of Victory (Season 14) (Horde)"},
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Dreadful Gladiator's Insignia of Conquest",
		ItemID:   93424,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Dreadful Gladiator's Insignia of Dominance",
		ItemID:   93601,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Dreadful Gladiator's Insignia of Victory",
		ItemID:   93611,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 15s. ( 15% chance, 85 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Brutal Talisman of the Shado-Pan Assault",
		ItemID:   94508,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 10s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Volatile Talisman of the Shado-Pan Assault",
		ItemID:   94510,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Talisman of the Shado-Pan Assault",
		ItemID:   94511,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your periodic damage spells have a chance to grant 1926 Intellect for 10s. (Approximately 1.10 procs per
	// minute)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 94521, ItemName: "Breath of the Hydra (N)"},
		{ItemID: 95711, ItemName: "Breath of the Hydra (LFR) (Celestial)"},
		{ItemID: 96083, ItemName: "Breath of the Hydra (Thunderforged)"},
		{ItemID: 96455, ItemName: "Breath of the Hydra (H)"},
		{ItemID: 96827, ItemName: "Breath of the Hydra (Heroic Thunderforged)"},
	})
	
	// When your attacks hit you have a chance to gain 2573 Agility and summon 3 Voodoo Gnomes for 10s. (Approximately
	// 1.10 procs per minute)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 94523, ItemName: "Bad Juju (N)"},
		{ItemID: 95665, ItemName: "Bad Juju (LFR) (Celestial)"},
		{ItemID: 96037, ItemName: "Bad Juju (Thunderforged)"},
		{ItemID: 96409, ItemName: "Bad Juju (H)"},
		{ItemID: 96781, ItemName: "Bad Juju (Heroic Thunderforged)"},
	})
	
	// When your spells deal critical damage, you have a chance to gain 1926 Intellect for 10s. (Approximately
	// 0.85 procs per minute)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeCrit,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 94531, ItemName: "Cha-Ye's Essence of Brilliance (N)"},
		{ItemID: 95772, ItemName: "Cha-Ye's Essence of Brilliance (LFR) (Celestial)"},
		{ItemID: 96144, ItemName: "Cha-Ye's Essence of Brilliance (Thunderforged)"},
		{ItemID: 96516, ItemName: "Cha-Ye's Essence of Brilliance (H)"},
		{ItemID: 96888, ItemName: "Cha-Ye's Essence of Brilliance (Heroic Thunderforged)"},
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Malevolent Gladiator's Insignia of Conquest",
		ItemID:   98760,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Malevolent Gladiator's Insignia of Dominance",
		ItemID:   98911,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crafted Malevolent Gladiator's Insignia of Victory",
		ItemID:   98917,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 100200, ItemName: "Grievous Gladiator's Insignia of Conquest (Season 14) (Alliance)"},
		{ItemID: 100586, ItemName: "Grievous Gladiator's Insignia of Conquest (Season 14) (Horde)"},
		{ItemID: 102840, ItemName: "Grievous Gladiator's Insignia of Conquest (Season 15) (Horde)"},
		{ItemID: 103150, ItemName: "Grievous Gladiator's Insignia of Conquest (Season 15) (Alliance)"},
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 100491, ItemName: "Grievous Gladiator's Insignia of Dominance (Season 14) (Alliance)"},
		{ItemID: 100712, ItemName: "Grievous Gladiator's Insignia of Dominance (Season 14) (Horde)"},
		{ItemID: 102963, ItemName: "Grievous Gladiator's Insignia of Dominance (Season 15) (Horde)"},
		{ItemID: 103309, ItemName: "Grievous Gladiator's Insignia of Dominance (Season 15) (Alliance)"},
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 100505, ItemName: "Grievous Gladiator's Insignia of Victory (Season 14) (Alliance)"},
		{ItemID: 100645, ItemName: "Grievous Gladiator's Insignia of Victory (Season 14) (Horde)"},
		{ItemID: 102896, ItemName: "Grievous Gladiator's Insignia of Victory (Season 15) (Horde)"},
		{ItemID: 103319, ItemName: "Grievous Gladiator's Insignia of Victory (Season 15) (Alliance)"},
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Stone of Battle",
		ItemID:   100990,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Idol of Battle",
		ItemID:   100991,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Defender Idol",
		ItemID:   100999,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart-Lesion Defender Stone",
		ItemID:   101002,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Idol of Rage",
		ItemID:   101009,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Stone of Rage",
		ItemID:   101012,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec
	// cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Idol of Destruction",
		ItemID:   101023,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Stone of Destruction",
		ItemID:   101026,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Springrain Stone of Wisdom",
		ItemID:   101041,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Trailseeker Idol of Rage",
		ItemID:   101054,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Trailseeker Stone of Rage",
		ItemID:   101057,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec
	// cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mountainsage Idol of Destruction",
		ItemID:   101069,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mountainsage Stone of Destruction",
		ItemID:   101072,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Defender Stone",
		ItemID:   101087,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Defender Idol",
		ItemID:   101089,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Stone of Wisdom",
		ItemID:   101107,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Idol of Rage",
		ItemID:   101113,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mistdancer Stone of Rage",
		ItemID:   101117,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Stone of Wisdom",
		ItemID:   101138,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Stone of Battle",
		ItemID:   101151,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Idol of Battle",
		ItemID:   101152,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Defender Idol",
		ItemID:   101160,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sunsoul Defender Stone",
		ItemID:   101163,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec
	// cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Communal Idol of Destruction",
		ItemID:   101168,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Communal Stone of Destruction",
		ItemID:   101171,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Communal Stone of Wisdom",
		ItemID:   101183,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightdrinker Idol of Rage",
		ItemID:   101200,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Lightdrinker Stone of Rage",
		ItemID:   101203,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Agility for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Idol of Rage",
		ItemID:   101217,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Stone of Rage",
		ItemID:   101220,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec
	// cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Idol of Destruction",
		ItemID:   101222,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Stone of Destruction",
		ItemID:   101225,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your healing spells have a chance to grant 1926 Intellect for 10s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Streamtalker Stone of Wisdom",
		ItemID:   101250,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Each time your harmful spells hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec
	// cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Felsoul Idol of Destruction",
		ItemID:   101263,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Felsoul Stone of Destruction",
		ItemID:   101266,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Stone of Battle",
		ItemID:   101294,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 Strength for 20s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Idol of Battle",
		ItemID:   101295,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// Your attacks have a chance to grant you 963 dodge for 15s. ( 15% chance, 55 sec cooldown)
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Defender Idol",
		ItemID:   101303,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s.
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Oathsworn Defender Stone",
		ItemID:   101306,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	})
	
	// When you deal damage you have a chance to gain 1287 Agility for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 102643, ItemName: "Prideful Gladiator's Insignia of Conquest (Season 15) (Alliance)"},
		{ItemID: 103347, ItemName: "Prideful Gladiator's Insignia of Conquest (Season 15) (Horde)"},
	})
	
	// When you deal damage you have a chance to gain 1287 Strength for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 102699, ItemName: "Prideful Gladiator's Insignia of Victory (Season 15) (Alliance)"},
		{ItemID: 103516, ItemName: "Prideful Gladiator's Insignia of Victory (Season 15) (Horde)"},
	})
	
	// When you deal damage or heal a target you have a chance to gain 1287 Intellect for 20s.
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 102766, ItemName: "Prideful Gladiator's Insignia of Dominance (Season 15) (Alliance)"},
		{ItemID: 103506, ItemName: "Prideful Gladiator's Insignia of Dominance (Season 15) (Horde)"},
	})
	
	// When your attacks hit you have a chance to gain 2573 Mastery for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskRangedAuto | core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage | core.ProcMaskMeleeProc | core.ProcMaskRangedProc | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 103686, ItemName: "Discipline of Xuen"},
		{ItemID: 103986, ItemName: "Discipline of Xuen (Timeless)"},
	})
	
	// When your spells deal damage you have a chance to gain 2573 critical strike for 20s. ( 15% chance, 115
	// sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellDamageProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 103687, ItemName: "Yu'lon's Bite"},
		{ItemID: 103987, ItemName: "Yu'lon's Bite (Timeless)"},
	})
	
	// Each time your melee attacks hit, you have a chance to gain 963 haste for 20s. ( 15% chance, 115 sec cooldown)
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeOHAuto | core.ProcMaskMeleeMHSpecial | core.ProcMaskMeleeOHSpecial | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
	}, []shared.ItemVariant{
		{ItemID: 103689, ItemName: "Alacrity of Xuen"},
		{ItemID: 103989, ItemName: "Alacrity of Xuen (Timeless)"},
	})
}