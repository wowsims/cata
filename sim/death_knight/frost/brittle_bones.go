package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Your Frost Fever also applies the Physical Vulnerability effect.

Physical Vulnerability
Weakens the constitution of an enemy target, increasing their physical damage taken by 4% for 30 sec.
*/
func (fdk *FrostDeathKnight) registerBrittleBones() {
	physVulnAuras := fdk.NewEnemyAuraArray(core.PhysVulnerabilityAura)

	fdk.Env.RegisterPreFinalizeEffect(func() {
		fdk.FrostFeverSpell.RelatedAuraArrays = fdk.FrostFeverSpell.RelatedAuraArrays.Append(physVulnAuras)
	})

	var lastDiseaseTarget *core.Unit
	core.MakePermanent(fdk.GetOrRegisterAura(core.Aura{
		Label:    "Brittle Bones" + fdk.Label,
		ActionID: core.ActionID{SpellID: 81328},
	})).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: death_knight.DeathKnightSpellFrostFever,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			lastDiseaseTarget = result.Target
			physVulnAuras.Get(result.Target).Activate(sim)
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: death_knight.DeathKnightSpellFrostFever,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			physVulnAuras.Get(lastDiseaseTarget).UpdateExpires(spell.Dot(lastDiseaseTarget).ExpiresAt())
		},
	})
}
