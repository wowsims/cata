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
		ItemID:     61462,
		AuraID:     92055,
		Bonus:      stats.Stats{stats.HasteRating: 1002},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Stonemother's Kiss",
		ItemID:     61411,
		AuraID:     90895,
		Bonus:      stats.Stats{stats.CritRating: 1164},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Witching Hourglass",
		ItemID:     56320,
		AuraID:     90887,
		Bonus:      stats.Stats{stats.HasteRating: 1710},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Witching Hourglass",
		ItemID:     55787,
		AuraID:     90885,
		Bonus:      stats.Stats{stats.HasteRating: 918},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Grace of the Herald",
		ItemID:     55266,
		AuraID:     92052,
		Bonus:      stats.Stats{stats.CritRating: 924},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Grace of the Herald (Heroic)",
		ItemID:     56295,
		AuraID:     92087,
		Bonus:      stats.Stats{stats.CritRating: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Porcelain Crab",
		ItemID:     55237,
		AuraID:     92166,
		Bonus:      stats.Stats{stats.MasteryRating: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Porcelain Crab (Heroic)",
		ItemID:     56280,
		AuraID:     92174,
		Bonus:      stats.Stats{stats.MasteryRating: 1710},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Key to the Endless Chamber",
		ItemID:     55795,
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
		ItemID:     56328,
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
		ItemID:     55810,
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
		ItemID:     56339,
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
		ItemID:     55819,
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
		ItemID:     56351,
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
		ItemID:     55854,
		AuraID:     91141,
		Bonus:      stats.Stats{stats.HasteRating: 1290},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Rainsong (Heroic)",
		ItemID:     56377,
		AuraID:     91143,
		Bonus:      stats.Stats{stats.HasteRating: 1710},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Tank-Commander Insignia",
		ItemID:     63841,
		AuraID:     91355,
		Bonus:      stats.Stats{stats.HasteRating: 1314},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Shrine-Cleansing Purifier",
		ItemID:     63838,
		AuraID:     91353,
		Bonus:      stats.Stats{stats.HasteRating: 1314},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.3,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Talisman of Sinister Order",
		ItemID:     65804,
		AuraID:     92166,
		Bonus:      stats.Stats{stats.MasteryRating: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Harrison's Insignia of Panache",
		ItemID:     65803,
		AuraID:     92164,
		Bonus:      stats.Stats{stats.MasteryRating: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Harrison's Insignia of Panache",
		ItemID:     65805,
		AuraID:     92164,
		Bonus:      stats.Stats{stats.MasteryRating: 918},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of the Vile",
		ItemID:     66969,
		AuraID:     92054,
		Bonus:      stats.Stats{stats.CritRating: 924},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Anhuur's Hymnal",
		ItemID:     55889,
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
		ItemID:     56407,
		AuraID:     90992,
		Bonus:      stats.Stats{stats.SpellPower: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Sorrowsong",
		ItemID:     55879,
		AuraID:     90990,
		Bonus:      stats.Stats{stats.SpellPower: 1512},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1.0,
		ICD:        time.Second * 20,

		CustomProcCondition: func(sim *core.Simulation, _ *core.Aura) bool {
			return sim.IsExecutePhase35()
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Sorrowsong (Heroic)",
		ItemID:     56400,
		AuraID:     91002,
		Bonus:      stats.Stats{stats.SpellPower: 1710},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1.0,
		ICD:        time.Second * 20,

		CustomProcCondition: func(sim *core.Simulation, _ *core.Aura) bool {
			return sim.IsExecutePhase35()
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Right Eye of Rajh",
		ItemID:     56100,
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
		ItemID:     56431,
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
		ItemID:     55995,
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
		ItemID:     56414,
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
		ItemID:     56121,
		AuraID:     92208,
		Bonus:      stats.Stats{stats.ParryRating: 1512},
		Duration:   time.Second * 12,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeParry,
		ProcChance: 1,
		ICD:        time.Minute,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Throngus's Finger (Heroic)",
		ItemID:     56449,
		AuraID:     92205,
		Bonus:      stats.Stats{stats.ParryRating: 1710},
		Duration:   time.Second * 12,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeParry,
		ProcChance: 1,
		ICD:        time.Minute,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of Solace",
		ItemID:     55868,
		AuraID:     91363,
		Bonus:      stats.Stats{stats.Strength: 1512},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of Solace (Heroic)",
		ItemID:     56393,
		AuraID:     91364,
		Bonus:      stats.Stats{stats.Strength: 1710},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Left Eye of Rajh",
		ItemID:     56102,
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
		ItemID:     56427,
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
		ItemID:     64762,
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
		ItemID:     64763,
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
		ItemID:     64761,
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
		ItemID:     62047,
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
		ProcMask:         core.ProcMaskSpellDamageProc,
		Outcome:          shared.OutcomeSpellNoMissCanCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Stump of Time (Horde)",
		ItemID:     62465,
		AuraID:     91047,
		Bonus:      stats.Stats{stats.SpellPower: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Stump of Time (Aliance)",
		ItemID:     62470,
		AuraID:     91047,
		Bonus:      stats.Stats{stats.SpellPower: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Unheeded Warning",
		ItemID:     59520,
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
		ItemID:     59224,
		AuraID:     91816,
		Bonus:      stats.Stats{stats.Strength: 1926},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heart of Rage (Heroic)",
		ItemID:     65072,
		AuraID:     92345,
		Bonus:      stats.Stats{stats.Strength: 2178},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Symbiotic Worm",
		ItemID:     59332,
		AuraID:     92235,
		Bonus:      stats.Stats{stats.MasteryRating: 963},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1,
		ICD:        time.Second * 30,

		CustomProcCondition: func(_ *core.Simulation, aura *core.Aura) bool {
			return aura.Unit.CurrentHealthPercent() < 0.35
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Symbiotic Worm (Heroic)",
		ItemID:     65048,
		AuraID:     92355,
		Bonus:      stats.Stats{stats.MasteryRating: 1089},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1,
		ICD:        time.Second * 30,

		CustomProcCondition: func(_ *core.Simulation, aura *core.Aura) bool {
			return aura.Unit.CurrentHealthPercent() < 0.35
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Mandala of Stirring Patterns (Horde)",
		ItemID:     62467,
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
		ItemID:     62472,
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
		ItemID:     59473,
		AuraID:     92126,
		Bonus:      stats.Stats{stats.CritRating: 1926},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Essence of the Cyclone (Heroic)",
		ItemID:     65140,
		AuraID:     92351,
		Bonus:      stats.Stats{stats.CritRating: 2178},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Theralion's Mirror",
		ItemID:     59519,
		AuraID:     91024,
		Bonus:      stats.Stats{stats.MasteryRating: 1926},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Theralion's Mirror (Heroic)",
		ItemID:     65105,
		AuraID:     92320,
		Bonus:      stats.Stats{stats.MasteryRating: 2178},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Crushing Weight",
		ItemID:     59506,
		AuraID:     91821,
		Bonus:      stats.Stats{stats.HasteRating: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Crushing Weight (Heroic)",
		ItemID:     65118,
		AuraID:     92342,
		Bonus:      stats.Stats{stats.HasteRating: 2178},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Bedrock Talisman",
		ItemID:     58182,
		AuraID:     92233,
		Bonus:      stats.Stats{stats.DodgeRating: 963},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1,
		ICD:        time.Second * 30,

		CustomProcCondition: func(_ *core.Simulation, aura *core.Aura) bool {
			return aura.Unit.CurrentHealthPercent() < 0.35
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Fall of Mortality",
		ItemID:     59500,
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
		ItemID:     65124,
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
		ItemID:     59326,
		AuraID:     91007,
		Bonus:      stats.Stats{stats.SpellPower: 1926},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellOrSpellProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Bell of Enraging Resonance (Heroic)",
		ItemID:     65053,
		AuraID:     92318,
		Bonus:      stats.Stats{stats.SpellPower: 2178},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellOrSpellProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 100,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Prestor's Talisman of Machination",
		ItemID:     59441,
		AuraID:     92124,
		Bonus:      stats.Stats{stats.HasteRating: 1926},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Prestor's Talisman of Machination (Heroic)",
		ItemID:     65026,
		AuraID:     92349,
		Bonus:      stats.Stats{stats.HasteRating: 2178},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 75,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Dominance - 365",
		ItemID:     61045,
		AuraID:     85027,
		Bonus:      stats.Stats{stats.SpellPower: 963},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Dominance - 371",
		ItemID:     70578,
		AuraID:     99719,
		Bonus:      stats.Stats{stats.SpellPower: 1077},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Ruthless Gladiator's Insignia of Dominance - 384",
		ItemID:     70402,
		AuraID:     99742,
		Bonus:      stats.Stats{stats.SpellPower: 1218},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Ruthless Gladiator's Insignia of Dominance - 390",
		ItemID:     72449,
		AuraID:     102435,
		Bonus:      stats.Stats{stats.SpellPower: 1287},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Cataclysmic Gladiator's Insignia of Dominance",
		ItemID:     73497,
		AuraID:     105137,
		Bonus:      stats.Stats{stats.SpellPower: 1452},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Victory - 365",
		ItemID:     61046,
		AuraID:     85032,
		Bonus:      stats.Stats{stats.Strength: 963},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Victory - 371",
		ItemID:     70579,
		AuraID:     99721,
		Bonus:      stats.Stats{stats.Strength: 1077},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Ruthless Gladiator's Insignia of Victory - 384",
		ItemID:     70403,
		AuraID:     99746,
		Bonus:      stats.Stats{stats.Strength: 1218},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Ruthless Gladiator's Insignia of Victory - 390",
		ItemID:     72455,
		AuraID:     102432,
		Bonus:      stats.Stats{stats.Strength: 1287},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Cataclysmic Gladiator's Insignia of Victory",
		ItemID:     73491,
		AuraID:     105139,
		Bonus:      stats.Stats{stats.Strength: 1452},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Conquest - 365",
		ItemID:     61047,
		AuraID:     85022,
		Bonus:      stats.Stats{stats.Agility: 963},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Vicious Gladiator's Insignia of Conquest - 371",
		ItemID:     70577,
		AuraID:     99717,
		Bonus:      stats.Stats{stats.Agility: 1077},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Ruthless Gladiator's Insignia of Conquest - 384",
		ItemID:     70404,
		AuraID:     99748,
		Bonus:      stats.Stats{stats.Agility: 1218},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Ruthless Gladiator's Insignia of Conquest - 390",
		ItemID:     72309,
		AuraID:     102439,
		Bonus:      stats.Stats{stats.Agility: 1287},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Cataclysmic Gladiator's Insignia of Conquest",
		ItemID:     73643,
		AuraID:     105135,
		Bonus:      stats.Stats{stats.Agility: 1452},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 55,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Dwyer's Caber",
		ItemID:     70141,
		AuraID:     100322,
		Bonus:      stats.Stats{stats.CritRating: 1020},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Coren's Chilled Chromium Coaster",
		ItemID:     232012,
		AuraID:     469349,
		Bonus:      stats.Stats{stats.AttackPower: 3576, stats.RangedAttackPower: 3576},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Petrified Pickled Egg",
		ItemID:     232014,
		AuraID:     469355,
		Bonus:      stats.Stats{stats.HasteRating: 1824},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Mithril Stopwatch",
		ItemID:     232013,
		AuraID:     469352,
		Bonus:      stats.Stats{stats.SpellPower: 1824},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "The Hungerer",
		ItemID:     68927,
		AuraID:     96911,
		Bonus:      stats.Stats{stats.HasteRating: 1532},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1, //TODO: verify proc chance, seems wrong?
		ICD:        time.Second * 60,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "The Hungerer Heroic",
		ItemID:     69112,
		AuraID:     97125,
		Bonus:      stats.Stats{stats.HasteRating: 1730},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 1, //TODO: verify proc chance, seems wrong?
		ICD:        time.Second * 60,
	})

	for _, version := range []ItemVersion{ItemVersionLFR, ItemVersionNormal, ItemVersionHeroic} {
		labelSuffix := []string{" (LFR)", "", " (Heroic)"}[version]

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:       "Creche of the Final Dragon" + labelSuffix,
			ItemID:     []int32{77972, 77205, 77992}[version],
			AuraID:     []int32{109742, 107988, 109744}[version],
			Bonus:      stats.Stats{stats.CritRating: []float64{2573, 2904, 3278}[version]},
			Duration:   time.Second * 20,
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 115,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:       "Insignia of the Corrupted Mind" + labelSuffix,
			ItemID:     []int32{77971, 77203, 77991}[version],
			AuraID:     []int32{109787, 107982, 109789}[version],
			Bonus:      stats.Stats{stats.HasteRating: []float64{2573, 2904, 3278}[version]},
			Duration:   time.Second * 20,
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskDirect | core.ProcMaskMeleeProc,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 115,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:       "Seal of the Seven Signs" + labelSuffix,
			ItemID:     []int32{77969, 77204, 77989}[version],
			AuraID:     []int32{109803, 107984, 109805}[version],
			Bonus:      stats.Stats{stats.HasteRating: []float64{2573, 2904, 3278}[version]},
			Duration:   time.Second * 20,
			Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			ProcMask:   core.ProcMaskSpellHealing,
			ProcChance: 0.15,
			ICD:        time.Second * 115,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:       "Soulshifter Vortex" + labelSuffix,
			ItemID:     []int32{77970, 77206, 77990}[version],
			AuraID:     []int32{109774, 107986, 109776}[version],
			Bonus:      stats.Stats{stats.MasteryRating: []float64{2573, 2904, 3278}[version]},
			Duration:   time.Second * 20,
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskDirect | core.ProcMaskMeleeProc,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 115,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:       "Starcatcher Compass" + labelSuffix,
			ItemID:     []int32{77973, 77202, 77993}[version],
			AuraID:     []int32{109709, 107982, 109711}[version],
			Bonus:      stats.Stats{stats.HasteRating: []float64{2573, 2904, 3278}[version]},
			Duration:   time.Second * 20,
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskDirect | core.ProcMaskMeleeProc,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 115,
		})

		// Ti'tahk, the Steps of Time
		// Equip: Your spells have a chance to grant you 1708/1928/2176 haste rating for 10 sec and 342/386/435 haste rating to up to 3 allies within 20 yards.
		// (Proc chance: 15%, 50s cooldown)
		// The buff has two effects, one for the caster and one shared.
		// * The first effect is 1366/1542/1741 haste rating on the caster.
		// * The second effect is 342/386/435 haste rating on the caster and up to 3 allies within 20 yards.
		// E.g. for the LFR version it's 1366 + 342 = 1708 haste rating for the caster with a shared ID so we just combine them.
		// TODO: Add the shared buff as an optional misc raid buff?
		//       Could be annoying with 3 different versions, uptime etc.
		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:       "Ti'tahk, the Steps of Time" + labelSuffix,
			ItemID:     []int32{78486, 77190, 78477}[version],
			AuraID:     []int32{109842, 107804, 109844}[version],
			Bonus:      stats.Stats{stats.HasteRating: []float64{1366 + 342, 1542 + 386, 1741 + 435}[version]},
			Duration:   time.Second * 10,
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellProc,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 50,
		})
	}

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Veil of Lies",
		ItemID:     72900,
		AuraID:     102667,
		Bonus:      stats.Stats{stats.DodgeRating: 1149},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Foul Gift of the Demon Lord",
		ItemID:     72898,
		AuraID:     102662,
		Bonus:      stats.Stats{stats.MasteryRating: 1149},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Arrow of Time",
		ItemID:     72897,
		AuraID:     102659,
		Bonus:      stats.Stats{stats.HasteRating: 1149},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Rosary of Light",
		ItemID:     72901,
		AuraID:     102660,
		Bonus:      stats.Stats{stats.CritRating: 1149},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 50,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Varo'then's Brooch",
		ItemID:     72899,
		AuraID:     102664,
		Bonus:      stats.Stats{stats.MasteryRating: 1149},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 50,
	})
}

var ItemSetAgonyAndTorment = core.NewItemSet(core.ItemSet{
	Name:  "Agony and Torment",
	Slots: core.MeleeWeaponSlots(),
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()

			procAura := character.NewTemporaryStatsAura(
				"Agony and Torment Proc",
				core.ActionID{SpellID: 95762},
				stats.Stats{stats.HasteRating: 1000},
				time.Second*10,
			)

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Agony and Torment Trigger",
				ActionID:   core.ActionID{SpellID: 95763},
				ProcMask:   core.ProcMaskMeleeOrRanged,
				Callback:   core.CallbackOnSpellHitDealt,
				ICD:        time.Second * 45,
				ProcChance: 0.1,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
	},
})
