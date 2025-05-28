package cata

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Gear Detector",
		ItemID:   61462,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Stonemother's Kiss",
		ItemID:   61411,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Witching Hourglass",
		ItemID:   56320,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Witching Hourglass",
		ItemID:   55787,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grace of the Herald",
		ItemID:   55266,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Grace of the Herald (Heroic)",
		ItemID:   56295,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Porcelain Crab",
		ItemID:   55237,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Porcelain Crab (Heroic)",
		ItemID:   56280,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Key to the Endless Chamber",
		ItemID:   55795,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Key to the Endless Chamber (Heroic)",
		ItemID:   56328,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tendrils of Burrowing Dark",
		ItemID:   55810,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tendrils of Burrowing Dark (Heroic)",
		ItemID:   56339,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tear of Blood",
		ItemID:   55819,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tear of Blood (Heroic)",
		ItemID:   56351,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Rainsong",
		ItemID:   55854,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Rainsong (Heroic)",
		ItemID:   56377,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Tank-Commander Insignia",
		ItemID:   63841,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Shrine-Cleansing Purifier",
		ItemID:   63838,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Talisman of Sinister Order",
		ItemID:   65804,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Harrison's Insignia of Panache",
		ItemID:   65803,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Harrison's Insignia of Panache",
		ItemID:   65805,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart of the Vile",
		ItemID:   66969,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Anhuur's Hymnal",
		ItemID:   55889,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Anhuur's Hymnal (Heroic)",
		ItemID:   56407,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sorrowsong",
		ItemID:   55879,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,

		CustomProcCondition: func(sim *core.Simulation, _ *core.Aura) bool {
			return sim.IsExecutePhase35()
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Sorrowsong (Heroic)",
		ItemID:   56400,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,

		CustomProcCondition: func(sim *core.Simulation, _ *core.Aura) bool {
			return sim.IsExecutePhase35()
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Right Eye of Rajh",
		ItemID:   56100,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Right Eye of Rajh (Heroic)",
		ItemID:   56431,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Blood of Isiset",
		ItemID:   55995,
		Callback: core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Blood of Isiset (Heroic)",
		ItemID:   56414,
		Callback: core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Throngus's Finger",
		ItemID:   56121,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeParry,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Throngus's Finger (Heroic)",
		ItemID:   56449,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeParry,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart of Solace",
		ItemID:   55868,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart of Solace (Heroic)",
		ItemID:   56393,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Left Eye of Rajh",
		ItemID:   56102,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Left Eye of Rajh (Heroic)",
		ItemID:   56427,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bloodthirsty Gladiator's Insignia of Dominance",
		ItemID:   64762,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeHit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bloodthirsty Gladiator's Insignia of Victory",
		ItemID:   64763,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskDirect,
		Outcome:  core.OutcomeHit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bloodthirsty Gladiator's Insignia of Conquest",
		ItemID:   64761,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskDirect,
		Outcome:  core.OutcomeHit,
	})

	shared.NewProcStatBonusEffectWithDamageProc(shared.ProcStatBonusEffect{
		Name:     "Darkmoon Card: Volcano",
		ItemID:   62047,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
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
		Name:     "Stump of Time (Horde)",
		ItemID:   62465,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Stump of Time (Aliance)",
		ItemID:   62470,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Unheeded Warning",
		ItemID:   59520,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart of Rage",
		ItemID:   59224,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Heart of Rage (Heroic)",
		ItemID:   65072,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Symbiotic Worm",
		ItemID:   59332,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskDirect,
		Outcome:  core.OutcomeLanded,

		CustomProcCondition: func(_ *core.Simulation, aura *core.Aura) bool {
			return aura.Unit.CurrentHealthPercent() < 0.35
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Symbiotic Worm (Heroic)",
		ItemID:   65048,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskDirect,
		Outcome:  core.OutcomeLanded,

		CustomProcCondition: func(_ *core.Simulation, aura *core.Aura) bool {
			return aura.Unit.CurrentHealthPercent() < 0.35
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mandala of Stirring Patterns (Horde)",
		ItemID:   62467,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mandala of Stirring Patterns (Alliance)",
		ItemID:   62472,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Essence of the Cyclone",
		ItemID:   59473,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Essence of the Cyclone (Heroic)",
		ItemID:   65140,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Theralion's Mirror",
		ItemID:   59519,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Theralion's Mirror (Heroic)",
		ItemID:   65105,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crushing Weight",
		ItemID:   59506,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Crushing Weight (Heroic)",
		ItemID:   65118,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bedrock Talisman",
		ItemID:   58182,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskDirect,
		Outcome:  core.OutcomeLanded,

		CustomProcCondition: func(_ *core.Simulation, aura *core.Aura) bool {
			return aura.Unit.CurrentHealthPercent() < 0.35
		},
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fall of Mortality",
		ItemID:   59500,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Fall of Mortality (Heroic)",
		ItemID:   65124,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bell of Enraging Resonance",
		ItemID:   59326,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellOrSpellProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Bell of Enraging Resonance (Heroic)",
		ItemID:   65053,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellOrSpellProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prestor's Talisman of Machination",
		ItemID:   59441,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Prestor's Talisman of Machination (Heroic)",
		ItemID:   65026,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Gladiator's Insignia of Dominance - 365",
		ItemID:   61045,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Gladiator's Insignia of Dominance - 371",
		ItemID:   70578,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ruthless Gladiator's Insignia of Dominance - 384",
		ItemID:   70402,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ruthless Gladiator's Insignia of Dominance - 390",
		ItemID:   72449,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cataclysmic Gladiator's Insignia of Dominance",
		ItemID:   73497,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Gladiator's Insignia of Victory - 365",
		ItemID:   61046,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Gladiator's Insignia of Victory - 371",
		ItemID:   70579,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ruthless Gladiator's Insignia of Victory - 384",
		ItemID:   70403,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ruthless Gladiator's Insignia of Victory - 390",
		ItemID:   72455,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cataclysmic Gladiator's Insignia of Victory",
		ItemID:   73491,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Gladiator's Insignia of Conquest - 365",
		ItemID:   61047,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Vicious Gladiator's Insignia of Conquest - 371",
		ItemID:   70577,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ruthless Gladiator's Insignia of Conquest - 384",
		ItemID:   70404,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Ruthless Gladiator's Insignia of Conquest - 390",
		ItemID:   72309,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Cataclysmic Gladiator's Insignia of Conquest",
		ItemID:   73643,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Dwyer's Caber",
		ItemID:   70141,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskDirect,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Coren's Chilled Chromium Coaster",
		ItemID:   232012,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeCrit,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Petrified Pickled Egg",
		ItemID:   232014,
		Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellHealing,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Mithril Stopwatch",
		ItemID:   232013,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "The Hungerer",
		ItemID:   68927,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "The Hungerer Heroic",
		ItemID:   69112,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
	})

	for _, version := range []ItemVersion{ItemVersionLFR, ItemVersionNormal, ItemVersionHeroic} {
		labelSuffix := []string{" (LFR)", "", " (Heroic)"}[version]

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     "Creche of the Final Dragon" + labelSuffix,
			ItemID:   []int32{77972, 77205, 77992}[version],
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
			Outcome:  core.OutcomeLanded,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     "Insignia of the Corrupted Mind" + labelSuffix,
			ItemID:   []int32{77971, 77203, 77991}[version],
			Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask: core.ProcMaskDirect | core.ProcMaskSpellProc,
			Outcome:  core.OutcomeLanded,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     "Seal of the Seven Signs" + labelSuffix,
			ItemID:   []int32{77969, 77204, 77989}[version],
			Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			ProcMask: core.ProcMaskSpellHealing,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     "Soulshifter Vortex" + labelSuffix,
			ItemID:   []int32{77970, 77206, 77990}[version],
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect | core.ProcMaskMeleeProc,
			Outcome:  core.OutcomeLanded,
		})

		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     "Starcatcher Compass" + labelSuffix,
			ItemID:   []int32{77973, 77202, 77993}[version],
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect | core.ProcMaskMeleeProc,
			Outcome:  core.OutcomeLanded,
		})
	}

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Veil of Lies",
		ItemID:   72900,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Foul Gift of the Demon Lord",
		ItemID:   72898,
		Callback: core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing | core.ProcMaskSpellProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Arrow of Time",
		ItemID:   72897,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Rosary of Light",
		ItemID:   72901,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Varo'then's Brooch",
		ItemID:   72899,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
		Outcome:  core.OutcomeLanded,
	})

	boostedProcAgility := map[int32]string{
		92116: "Thundercaller Idol of Rage",
		92118: "Darkwalker Idol of Rage",
		92133: "Naturalist Idol of Rage",
		92142: "Forestwalker Idol of Rage",
		92401: "Waterdancer Idol of Rage",
	}
	for id, name := range boostedProcAgility {
		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     name,
			ItemID:   id,
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
			Outcome:  core.OutcomeLanded,
		})
	}

	boostedProcStrength := map[int32]string{
		92128: "Martial Idol of Battle",
		92148: "Partisan Idol of Battle",
		92167: "Scourgeheart Idol of Battle",
	}
	for id, name := range boostedProcStrength {
		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     name,
			ItemID:   id,
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
			Outcome:  core.OutcomeLanded,
		})
	}

	boostedProcDodge := map[int32]string{
		92127: "Martial Defender Idol",
		92135: "Scourgeheart Defender Idol",
		92147: "Partisan Defender Idol",
		92399: "Waterdancer Defender Idol",
	}
	for id, name := range boostedProcDodge {
		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     name,
			ItemID:   id,
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
			Outcome:  core.OutcomeLanded,
		})
	}

	boostedHealProcIntellect := map[int32]string{
		92115: "Deliverer Stone of Wisdom",
		92122: "Thundercaller Stone of Wisdom",
		92139: "Naturalist Stone of Wisdom",
		92145: "Partisan Stone of Wisdom",
		92402: "Waterdancer Stone of Wisdom",
	}
	for id, name := range boostedHealProcIntellect {
		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     name,
			ItemID:   id,
			Callback: core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			ProcMask: core.ProcMaskSpellHealing,
		})
	}

	boostedProcMastery := map[int32]string{
		92114: "Partisan Defender Stone",
		92117: "Darkwalker Stone of Rage",
		92119: "Thundercaller Stone of Destruction",
		92121: "Thundercaller Stone of Rage",
		92124: "Soulseizer Stone of Destruction",
		92126: "Martial Defender Stone",
		92129: "Martial Stone of Battle",
		92134: "Scourgeheart Defender Stone",
		92136: "Naturalist Stone of Destruction",
		92138: "Naturalist Stone of Rage",
		92141: "Forestwalker Stone of Rage",
		92143: "Enlightened Stone of Destruction",
		92149: "Partisan Stone of Battle",
		92151: "Deliverer Stone of Destruction",
		92168: "Scourgeheart Stone of Battle",
		92398: "Waterdancer Defender Stone",
		92400: "Waterdancer Stone of Rage",
	}
	for id, name := range boostedProcMastery {
		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     name,
			ItemID:   id,
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect | core.ProcMaskProc,
			Outcome:  core.OutcomeLanded,
		})
	}

	boostedSpellHitProcHaste := map[int32]string{
		92113: "Deliverer Idol of Destruction",
		92120: "Thundercaller Idol of Destruction",
		92125: "Soulseizer Idol of Destruction",
		92137: "Naturalist Idol of Destruction",
		92144: "Enlightened Idol of Destruction",
	}
	for id, name := range boostedSpellHitProcHaste {
		shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
			Name:     name,
			ItemID:   id,
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellOrSpellProc,
			Outcome:  core.OutcomeLanded,
		})
	}
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
