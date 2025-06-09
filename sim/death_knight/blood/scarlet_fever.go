package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (dk *BloodDeathKnight) registerScarletFever() {
	weakenedBlowsAuras := dk.NewEnemyAuraArray(core.WeakenedBlowsAura)
	dk.Env.RegisterPreFinalizeEffect(func() {
		dk.BloodPlagueSpell.RelatedAuraArrays = dk.BloodPlagueSpell.RelatedAuraArrays.Append(weakenedBlowsAuras)
	})

	var lastDiseaseTarget *core.Unit
	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label:    "Scarlet Fever" + dk.Label,
		ActionID: core.ActionID{SpellID: 51160},
	})).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: death_knight.DeathKnightSpellBloodBoil,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			frostFever := dk.FrostFeverSpell.Dot(result.Target)
			bloodPlague := dk.BloodPlagueSpell.Dot(result.Target)

			if frostFever.IsActive() {
				frostFever.Refresh(sim)
				weakenedBlowsAuras.Get(result.Target).Activate(sim)
			}

			if bloodPlague.IsActive() {
				bloodPlague.Refresh(sim)
				weakenedBlowsAuras.Get(result.Target).Activate(sim)
			}
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: death_knight.DeathKnightSpellBloodPlague,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			lastDiseaseTarget = result.Target
			weakenedBlowsAuras.Get(result.Target).Activate(sim)
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: death_knight.DeathKnightSpellBloodPlague,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			weakenedBlowsAuras.Get(lastDiseaseTarget).UpdateExpires(spell.Dot(lastDiseaseTarget).ExpiresAt())
		},
	})
}
