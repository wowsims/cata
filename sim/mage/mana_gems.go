package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerManaGems() {

	var manaGain float64
	actionID := core.ActionID{ItemID: 36799}
	manaMetrics := mage.NewManaMetrics(actionID)
	hasMajorGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfManaGem)
	hasMinorGlyph := mage.HasMinorGlyph(proto.MageMinorGlyph_GlyphOfLooseMana)
	maxManaGems := core.Ternary(hasMajorGlyph, 10, 3)

	minManaGain := 42750.0
	maxManaGain := 47250.0

	var remainingManaGems int
	mage.RegisterResetEffect(func(sim *core.Simulation) {
		remainingManaGems = maxManaGems
	})

	minorGlyphAura := mage.RegisterAura(core.Aura{
		Label:    "Replenish Mana",
		ActionID: core.ActionID{SpellID: 5405},
		Duration: 6*time.Second + 1, // add 1 ns duration offset in order to guarantee that the final tick fires

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			manaPerTick := 45000.0 / 5
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second,
				NumTicks: 5,
				Priority: core.ActionPriorityRegen,

				OnAction: func(sim *core.Simulation) {
					if aura.IsActive() {
						mage.AddMana(sim, manaPerTick, manaMetrics)
					}
				},
			})
		},
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		// Don't use if we don't have any gems remaining!
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return remainingManaGems != 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if hasMinorGlyph {
				minorGlyphAura.Activate(sim)
			} else {
				manaGain = sim.Roll(minManaGain, maxManaGain)
				mage.AddMana(sim, manaGain, manaMetrics)
			}

			remainingManaGems--
			if remainingManaGems == 0 {
				mage.GetMajorCooldown(actionID).Disable()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Only pop if we have less than the max mana provided by the gem minus 1mp5 tick.
			totalRegen := character.ManaRegenPerSecondWhileCombat() * 5
			return character.MaxMana()-(character.CurrentMana()+totalRegen) >= maxManaGain
		},
	})
}
