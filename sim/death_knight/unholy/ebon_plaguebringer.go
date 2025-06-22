package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Increases the damage your diseases deal by 60%, causes your Plague Strike to also apply Frost Fever, and causes your Blood Plague to also apply the Physical Vulnerability effect.

Physical Vulnerability
Weakens the constitution of an enemy target, increasing their physical damage taken by 4% for 30 sec.
*/
func (uhdk *UnholyDeathKnight) registerEbonPlaguebringer() {
	physVulnAuras := uhdk.NewEnemyAuraArray(core.PhysVulnerabilityAura)

	uhdk.Env.RegisterPreFinalizeEffect(func() {
		uhdk.BloodPlagueSpell.RelatedAuraArrays = uhdk.BloodPlagueSpell.RelatedAuraArrays.Append(physVulnAuras)
	})

	var lastDiseaseTarget *core.Unit
	core.MakePermanent(uhdk.GetOrRegisterAura(core.Aura{
		Label:    "Ebon Plaguebringer" + uhdk.Label,
		ActionID: core.ActionID{SpellID: 51160},
	})).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: death_knight.DeathKnightSpellPlagueStrike,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			uhdk.FrostFeverSpell.Cast(sim, result.Target)
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: death_knight.DeathKnightSpellBloodPlague,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			lastDiseaseTarget = result.Target
			physVulnAuras.Get(result.Target).Activate(sim)
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: death_knight.DeathKnightSpellBloodPlague,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			physVulnAuras.Get(lastDiseaseTarget).UpdateExpires(spell.Dot(lastDiseaseTarget).ExpiresAt())
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.DeathKnightSpellDisease,
		FloatValue: 0.6,
	})
}
