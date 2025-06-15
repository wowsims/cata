package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Causes your Blood Boil to refresh your diseases on targets it damages, and your Blood Plague to also afflict enemies with Weakened Blows.

Weakened Blows
Demoralizes the target, reducing their physical damage dealt by 10% for 30 sec.
*/
func (bdk *BloodDeathKnight) registerScarletFever() {
	weakenedBlowsAuras := bdk.NewEnemyAuraArray(core.WeakenedBlowsAura)
	bdk.Env.RegisterPreFinalizeEffect(func() {
		bdk.BloodPlagueSpell.RelatedAuraArrays = bdk.BloodPlagueSpell.RelatedAuraArrays.Append(weakenedBlowsAuras)
	})

	var lastDiseaseTarget *core.Unit
	core.MakePermanent(bdk.GetOrRegisterAura(core.Aura{
		Label:    "Scarlet Fever" + bdk.Label,
		ActionID: core.ActionID{SpellID: 51160},
	})).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: death_knight.DeathKnightSpellBloodBoil,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			frostFever := bdk.FrostFeverSpell.Dot(result.Target)
			bloodPlague := bdk.BloodPlagueSpell.Dot(result.Target)

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
