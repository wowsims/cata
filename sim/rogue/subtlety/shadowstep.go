package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/rogue"
)

func (subRogue *SubtletyRogue) registerShadowstepCD() {
	actionID := core.ActionID{SpellID: 36554}
	var affectedSpells []*core.Spell

	subRogue.ShadowstepAura = subRogue.RegisterAura(core.Aura{
		Label:    "Shadowstep",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = append(affectedSpells, subRogue.Ambush)
			affectedSpells = append(affectedSpells, subRogue.Garrote)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Damage of your next ability is increased by 20% and the threat caused is reduced by 50%.
			for _, spell := range affectedSpells {
				spell.DamageMultiplier *= 1.2
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplier *= 1 / 1.2
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			for _, affectedSpell := range affectedSpells {
				if spell == affectedSpell {
					aura.Deactivate(sim)
				}
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
