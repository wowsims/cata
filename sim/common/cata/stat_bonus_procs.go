package cata

import (
	"time"

	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Gear Detector",
		ID:         61462,
		AuraID:     92055,
		Bonus:      stats.Stats{stats.MeleeHaste: 1002, stats.SpellHaste: 1002},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Stonemother's Kiss",
		ID:         61411,
		AuraID:     90895,
		Bonus:      stats.Stats{stats.MeleeCrit: 1164, stats.SpellCrit: 1164},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Witching Hourglass",
		ID:         56320,
		AuraID:     90887,
		Bonus:      stats.Stats{stats.MeleeHaste: 1710, stats.SpellHaste: 1710},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Witching Hourglass",
		ID:         55787,
		AuraID:     90885,
		Bonus:      stats.Stats{stats.MeleeHaste: 918, stats.SpellHaste: 918},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Grace of the Herald",
		ID:         55266,
		AuraID:     92052,
		Bonus:      stats.Stats{stats.MeleeCrit: 924, stats.SpellCrit: 924},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Grace of the Herald (Heroic)",
		ID:         56295,
		AuraID:     92087,
		Bonus:      stats.Stats{stats.MeleeCrit: 1710, stats.SpellCrit: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Porcelain Crab",
		ID:         55237,
		AuraID:     92166,
		Bonus:      stats.Stats{stats.Mastery: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Porcelain Crab (Heroic)",
		ID:         56280,
		AuraID:     92174,
		Bonus:      stats.Stats{stats.Mastery: 1710},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Key to the Endless Chamber",
		ID:         55795,
		AuraID:     92069,
		Bonus:      stats.Stats{stats.Agility: 1290},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Key to the Endless Chamber (Heroic)",
		ID:         56328,
		AuraID:     92091,
		Bonus:      stats.Stats{stats.Agility: 1710},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Tendrils of Burrowing Dark",
		ID:         55810,
		AuraID:     90896,
		Bonus:      stats.Stats{stats.SpellPower: 1290},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Tendrils of Burrowing Dark (Heroic)",
		ID:         56339,
		AuraID:     90898,
		Bonus:      stats.Stats{stats.SpellPower: 1710},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Tear of Blood",
		ID:         55819,
		AuraID:     91138,
		Bonus:      stats.Stats{stats.Spirit: 1290},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Tear of Blood (Heroic)",
		ID:         56351,
		AuraID:     91139,
		Bonus:      stats.Stats{stats.Spirit: 1710},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Rainsong",
		ID:         55854,
		AuraID:     91141,
		Bonus:      stats.Stats{stats.MeleeHaste: 1290, stats.SpellHaste: 1290},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Rainsong (Heroic)",
		ID:         56377,
		AuraID:     91143,
		Bonus:      stats.Stats{stats.MeleeHaste: 1710, stats.SpellHaste: 1710},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Tank-Commander Insignia",
		ID:         63841,
		AuraID:     91355,
		Bonus:      stats.Stats{stats.MeleeHaste: 1314, stats.SpellHaste: 1314},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Shrine-Cleansing Purifier",
		ID:         63838,
		AuraID:     91353,
		Bonus:      stats.Stats{stats.MeleeHaste: 1314, stats.SpellHaste: 1314},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Talisman of Sinister Order",
		ID:         65804,
		AuraID:     92166,
		Bonus:      stats.Stats{stats.Mastery: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Harrison's Insignia of Panache",
		ID:         65803,
		AuraID:     92164,
		Bonus:      stats.Stats{stats.Mastery: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Harrison's Insignia of Panache",
		ID:         65805,
		AuraID:     92164,
		Bonus:      stats.Stats{stats.Mastery: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of the Vile",
		ID:         66969,
		AuraID:     92054,
		Bonus:      stats.Stats{stats.MeleeCrit: 924, stats.SpellCrit: 924},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Anhuur's Hymnal",
		ID:         55889,
		AuraID:     90989,
		Bonus:      stats.Stats{stats.SpellPower: 1512},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Anhuur's Hymnal (Heroic)",
		ID:         56407,
		AuraID:     90992,
		Bonus:      stats.Stats{stats.SpellPower: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffectWithCustomCondition(shared.ProcStatBonusEffect{
		Name:       "Sorrowsong",
		ID:         55879,
		AuraID:     90990,
		Bonus:      stats.Stats{stats.SpellPower: 1512},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 20,
	}, func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
		return sim.IsExecutePhase35()
	})

	shared.NewProcStatBonusEffectWithCustomCondition(shared.ProcStatBonusEffect{
		Name:       "Sorrowsong (Heroic)",
		ID:         56400,
		AuraID:     91002,
		Bonus:      stats.Stats{stats.SpellPower: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 20,
	}, func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
		return sim.IsExecutePhase35()
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Right Eye of Rajh",
		ID:         56100,
		AuraID:     91370,
		Bonus:      stats.Stats{stats.Strength: 1512},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.5,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Right Eye of Rajh (Heroic)",
		ID:         56431,
		AuraID:     91368,
		Bonus:      stats.Stats{stats.Strength: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.5,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Blood of Isiset",
		ID:         55995,
		AuraID:     91147,
		Bonus:      stats.Stats{stats.Spirit: 1512},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Blood of Isiset (Heroic)",
		ID:         56414,
		AuraID:     91149,
		Bonus:      stats.Stats{stats.Spirit: 1710},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Throngus's Finger",
		ID:         56121,
		AuraID:     92208,
		Bonus:      stats.Stats{stats.Parry: 1512},
		Duration:   time.Second * 12,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeParry,
		ProcChance: 1,
		ICD:        time.Minute,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Throngus's Finger (Heroic)",
		ID:         56449,
		AuraID:     92205,
		Bonus:      stats.Stats{stats.Parry: 1710},
		Duration:   time.Second * 12,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeParry,
		ProcChance: 1,
		ICD:        time.Minute,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of Solace",
		ID:         55868,
		AuraID:     91363,
		Bonus:      stats.Stats{stats.Strength: 1512},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeParry,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of Solace (Heroic)",
		ID:         56393,
		AuraID:     91364,
		Bonus:      stats.Stats{stats.Strength: 1710},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Left Eye of Rajh",
		ID:         56102,
		AuraID:     92096,
		Bonus:      stats.Stats{stats.Agility: 1512},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.5,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Left Eye of Rajh (Heroic)",
		ID:         56427,
		AuraID:     92094,
		Bonus:      stats.Stats{stats.Agility: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.5,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Bloodthirsty Gladiator's Insignia of Dominance",
		ID:         64762,
		AuraID:     92218,
		Bonus:      stats.Stats{stats.SpellPower: 912},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:   core.ProcMaskSpellHealing | core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeHit,
		ProcChance: 0.25,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Bloodthirsty Gladiator's Insignia of Victory",
		ID:         64763,
		AuraID:     92216,
		Bonus:      stats.Stats{stats.Strength: 912},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeHit,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Bloodthirsty Gladiator's Insignia of Conquest",
		ID:         64761,
		AuraID:     92220,
		Bonus:      stats.Stats{stats.Agility: 912},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeHit,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffectWithDamageProc(shared.ProcStatBonusEffect{
		Name:       "Darkmoon Card: Volcano",
		ID:         62047,
		AuraID:     89091,
		Bonus:      stats.Stats{stats.Intellect: 1600},
		Duration:   time.Second * 12,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.3,
		ICD:        time.Second * 45,
	}, shared.DamageEffect{
		SpellID:          89091,
		School:           core.SpellSchoolFire,
		MinDmg:           900,
		MaxDmg:           1500,
		BonusCoefficient: 0.1,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Stump of Time (Horde)",
		ID:         62465,
		AuraID:     91047,
		Bonus:      stats.Stats{stats.SpellPower: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Stump of Time (Aliance)",
		ID:         62470,
		AuraID:     91048,
		Bonus:      stats.Stats{stats.SpellPower: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Unheeded Warning",
		ID:         59520,
		AuraID:     92108,
		Bonus:      stats.Stats{stats.AttackPower: 1926, stats.RangedAttackPower: 1926},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of Rage",
		ID:         59224,
		AuraID:     91816,
		Bonus:      stats.Stats{stats.Strength: 1926},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of Rage (Heroic)",
		ID:         65072,
		AuraID:     92345,
		Bonus:      stats.Stats{stats.Strength: 2178},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffectWithCustomCondition(shared.ProcStatBonusEffect{
		Name:       "Symbiotic Worm",
		ID:         59332,
		AuraID:     92235,
		Bonus:      stats.Stats{stats.Mastery: 963},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1,
		ICD:        time.Second * 30,
	}, func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
		return result.Target.CurrentHealthPercent() < 0.35
	})

	shared.NewProcStatBonusEffectWithCustomCondition(shared.ProcStatBonusEffect{
		Name:       "Symbiotic Worm (Heroic)",
		ID:         65048,
		AuraID:     92355,
		Bonus:      stats.Stats{stats.Mastery: 1089},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1,
		ICD:        time.Second * 30,
	}, func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
		return result.Target.CurrentHealthPercent() < 0.35
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Mandala of Stirring Patterns (Horde)",
		ID:         62467,
		AuraID:     91192,
		Bonus:      stats.Stats{stats.Intellect: 1926},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Mandala of Stirring Patterns (Alliance)",
		ID:         62472,
		AuraID:     91192,
		Bonus:      stats.Stats{stats.Intellect: 1926},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Essence of the Cyclone",
		ID:         59473,
		AuraID:     92126,
		Bonus:      stats.Stats{stats.MeleeCrit: 1926, stats.SpellCrit: 1926},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Essence of the Cyclone (Heroic)",
		ID:         65140,
		AuraID:     92351,
		Bonus:      stats.Stats{stats.MeleeCrit: 2178, stats.SpellCrit: 2178},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Theralion's Mirror",
		ID:         59519,
		AuraID:     91024,
		Bonus:      stats.Stats{stats.Mastery: 1926},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Theralion's Mirror (Heroic)",
		ID:         65105,
		AuraID:     92320,
		Bonus:      stats.Stats{stats.Mastery: 2178},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Crushing Weight",
		ID:         59506,
		AuraID:     91821,
		Bonus:      stats.Stats{stats.MeleeHaste: 1926, stats.SpellHaste: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Crushing Weight (Heroic)",
		ID:         65118,
		AuraID:     92342,
		Bonus:      stats.Stats{stats.MeleeHaste: 2178, stats.SpellHaste: 2178},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffectWithCustomCondition(shared.ProcStatBonusEffect{
		Name:       "Bedrock Talisman",
		ID:         58182,
		AuraID:     92233,
		Bonus:      stats.Stats{stats.Dodge: 963},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1,
		ICD:        time.Second * 30,
	}, func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
		return result.Target.CurrentHealthPercent() < 0.35
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Fall of Mortality",
		ID:         59500,
		AuraID:     91821,
		Bonus:      stats.Stats{stats.Spirit: 1962},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Fall of Mortality (Heroic)",
		ID:         65124,
		AuraID:     92332,
		Bonus:      stats.Stats{stats.Spirit: 2178},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Bell of Enraging Resonance",
		ID:         59326,
		AuraID:     91007,
		Bonus:      stats.Stats{stats.SpellPower: 1926},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Bell of Enraging Resonance (Heroic)",
		ID:         65053,
		AuraID:     92318,
		Bonus:      stats.Stats{stats.SpellPower: 2178},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Prestor's Talisman of Machination",
		ID:         59441,
		AuraID:     92124,
		Bonus:      stats.Stats{stats.MeleeHaste: 1926, stats.SpellHaste: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Prestor's Talisman of Machination (Heroic)",
		ID:         65026,
		AuraID:     92349,
		Bonus:      stats.Stats{stats.MeleeHaste: 2178, stats.SpellHaste: 2178},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Dominance",
		ID:         61045,
		AuraID:     85027,
		Bonus:      stats.Stats{stats.SpellPower: 963},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Victory",
		ID:         61046,
		AuraID:     85032,
		Bonus:      stats.Stats{stats.Strength: 963},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Conquest",
		ID:         61047,
		AuraID:     85022,
		Bonus:      stats.Stats{stats.Agility: 963},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})
}
