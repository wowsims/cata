package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

// Rather than update a variable somewhere for one effect (Fury's Unshackled Fury) just take a callback
// to fetch its multiplier when needed
type RageMultiplierCB func() float64

func (warrior *Warrior) RegisterBerserkerRageSpell() {

	actionID := core.ActionID{SpellID: 18499}
	rageBonus := float64(core.TernaryInt(warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfBerserkerRage), 5, 0))
	rageMetrics := warrior.NewRageMetrics(actionID)

	warrior.BerserkerRageAura = warrior.RegisterAura(core.Aura{
		Label:    "Berserker Rage",
		Tag:      EnrageTag,
		ActionID: actionID,
		Duration: time.Second * 10,
		// TODO: rage from damage taken multiplier
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {

		},
	})

	warrior.BerserkerRage = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskBerserkerRage,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// Fetch multiplier in here so we capture any mastery buffs that might affect it
			warrior.AddRage(sim, rageBonus*warrior.EnrageEffectMultiplier, rageMetrics)
			warrior.BerserkerRageAura.Activate(sim)
		},
	})
}
