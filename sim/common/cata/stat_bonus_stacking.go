package cata

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	shared.NewStackingStatBonusCD(shared.StackingStatBonusCD{
		Name:                  "World-Queller Focus",
		ID:                    63842,
		AuraID:                90927,
		Bonus:                 stats.Stats{stats.SpellPower: 313},
		MaxStacks:             5,
		ProcMask:              core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Duration:              time.Second * 20,
		CD:                    time.Minute * 2,
		Callback:              core.CallbackOnCastComplete,
		Harmful:               false,
		TrinketLimitsDuration: true,
		ProcChance:            1,
		IsDefensive:           false,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Gale of Shadows",
		ItemID:     56138,
		AuraID:     90953,
		Bonus:      stats.Stats{stats.SpellPower: 15},
		MaxStacks:  20,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt,
		Harmful:    false,
		ProcChance: 1,
		Icd:        time.Millisecond * 500,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Gale of Shadows (Heroic)",
		ItemID:     56462,
		AuraID:     90985,
		Bonus:      stats.Stats{stats.SpellPower: 17},
		MaxStacks:  20,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt,
		Harmful:    false,
		ProcChance: 1,
		Icd:        time.Millisecond * 500,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Tia's Grace",
		ItemID:     55874,
		AuraID:     92085,
		Bonus:      stats.Stats{stats.Agility: 30},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Tia's Grace (Heroic)",
		ItemID:     56394,
		AuraID:     92089,
		Bonus:      stats.Stats{stats.Agility: 34},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Darkmoon Card: Tsunami",
		ItemID:     62050,
		AuraID:     92090,
		Bonus:      stats.Stats{stats.Spirit: 80},
		MaxStacks:  5,
		ProcMask:   core.ProcMaskSpellHealing,
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		Harmful:    false,
		ProcChance: 1,
		Icd:        time.Second * 2,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Fluid Death",
		ItemID:     58181,
		AuraID:     92104,
		Bonus:      stats.Stats{stats.Agility: 38},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "License to Slay",
		ItemID:     58180,
		AuraID:     91810,
		Bonus:      stats.Stats{stats.Strength: 38},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMeleeOrMeleeProc,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Necromantic Focus",
		ItemID:     68982,
		AuraID:     96962,
		Bonus:      stats.Stats{stats.MasteryRating: 39},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskSpellDamage,
		Duration:   time.Second * 10,
		Outcome:    core.OutcomeLanded,
		Callback:   core.CallbackOnPeriodicDamageDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Necromantic Focus Heroic",
		ItemID:     69139,
		AuraID:     97131,
		Bonus:      stats.Stats{stats.MasteryRating: 44},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskSpellDamage,
		Duration:   time.Second * 10,
		Outcome:    core.OutcomeLanded,
		Callback:   core.CallbackOnPeriodicDamageDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	for _, version := range []ItemVersion{ItemVersionLFR, ItemVersionNormal, ItemVersionHeroic} {
		labelSuffix := []string{" (LFR)", "", " (Heroic)"}[version]

		shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
			Name:       "Eye of Unmaking" + labelSuffix,
			ItemID:     []int32{77977, 77200, 77997}[version],
			AuraID:     []int32{109748, 107966, 109750}[version],
			Bonus:      stats.Stats{stats.Strength: []float64{78, 88, 99}[version]},
			MaxStacks:  10,
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskMelee,
			Duration:   time.Second * 10,
			ProcChance: 1,
		})

		shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
			Name:       "Resolve of Undying" + labelSuffix,
			ItemID:     []int32{77978, 77201, 77998}[version],
			AuraID:     []int32{109780, 107968, 109782}[version],
			Bonus:      stats.Stats{stats.DodgeRating: []float64{78, 88, 99}[version]},
			MaxStacks:  10,
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskMelee,
			Duration:   time.Second * 10,
			ProcChance: 1,
		})

		shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
			Name:       "Will of Unbinding" + labelSuffix,
			ItemID:     []int32{77975, 77198, 77995}[version],
			AuraID:     []int32{109793, 107970, 109795}[version],
			Bonus:      stats.Stats{stats.Intellect: []float64{78, 88, 99}[version]},
			MaxStacks:  10,
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			Duration:   time.Second * 10,
			ProcChance: 1,
		})

		shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
			Name:       "Wrath of Unchaining" + labelSuffix,
			ItemID:     []int32{77974, 77197, 77994}[version],
			AuraID:     []int32{109717, 107960, 109719}[version],
			Bonus:      stats.Stats{stats.Agility: []float64{78, 88, 99}[version]},
			MaxStacks:  10,
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Duration:   time.Second * 10,
			ProcChance: 1,
		})
	}
}
