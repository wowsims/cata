package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/rogue"
)

func (subRogue *SubtletyRogue) registerShadowstepCD() {
	actionID := core.ActionID{SpellID: 36554}

	affectedClassSpellMasks := rogue.RogueSpellAmbush | rogue.RogueSpellGarrote
	damageMultiMod := subRogue.AddDynamicMod(core.SpellModConfig{
		ClassMask:  rogue.RogueSpellAmbush | rogue.RogueSpellGarrote,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1.2,
	})

	subRogue.ShadowstepAura = subRogue.RegisterAura(core.Aura{
		Label:    "Shadowstep",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Damage of your next ability is increased by 20% and the threat caused is reduced by 50%.
			damageMultiMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMultiMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(affectedClassSpellMasks) {
				aura.Deactivate(sim)
			}
		},
	})

	subRogue.Shadowstep = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellShadowstep,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    subRogue.NewTimer(),
				Duration: time.Second * 24,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: subRogue.ShadowstepAura,
	})
}
