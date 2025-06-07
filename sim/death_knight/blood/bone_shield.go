package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (bdk *BloodDeathKnight) registerBoneShield() {
	actionID := core.ActionID{SpellID: 49222}

	var boneShieldAura *core.Aura
	boneShieldAura = bdk.RegisterAura(core.Aura{
		Label:     "Bone Shield" + bdk.Label,
		ActionID:  actionID,
		Duration:  time.Minute * 5,
		MaxStacks: 6,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeLanded,
		ICD:      time.Second * 2,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			boneShieldAura.RemoveStack(sim)
		},
	}).AttachMultiplicativePseudoStatBuff(&bdk.PseudoStats.DamageTakenMultiplier, 0.8)

	spell := bdk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellBoneShield,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    bdk.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
			spell.RelatedSelfBuff.SetStacks(sim, spell.RelatedSelfBuff.MaxStacks)
		},

		RelatedSelfBuff: boneShieldAura,
	})

	if !bdk.Inputs.IsDps {
		bdk.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
			ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
				return bdk.CurrentHealthPercent() < 0.6
			},
		})
	}
}
