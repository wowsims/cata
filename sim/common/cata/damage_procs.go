package cata

import (
	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
)

func init() {
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:  62049,
		SpellID: 89087,
		School:  core.SpellSchoolNature,
		MinDmg:  5250,
		MaxDmg:  8750,
		Flags:   core.SpellFlagNoSpellMods | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,
		Outcome: shared.OutcomeMeleeNoBlockDodgeParryCrit,
		Trigger: core.ProcTrigger{
			Name:     "Darkmoon Card: Hurricane",
			ProcMask: core.ProcMaskMeleeOrRanged,
			PPM:      1,
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
		},
	})

	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:  62051,
		SpellID: 89087,
		School:  core.SpellSchoolNature,
		MinDmg:  5250,
		MaxDmg:  8750,
		Outcome: shared.OutcomeMeleeNoBlockDodgeParryCrit,
		Flags:   core.SpellFlagNoSpellMods | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,
		Trigger: core.ProcTrigger{
			Name:     "Darkmoon Card: Hurricane",
			ProcMask: core.ProcMaskProc,
			PPM:      1,
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
		},
	})
}
