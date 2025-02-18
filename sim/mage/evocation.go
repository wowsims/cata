package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerEvocation() {
	actionID := core.ActionID{SpellID: 12051}
	manaMetrics := mage.NewManaMetrics(actionID)
	hasGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfEvocation)

	evocation := mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL | core.SpellFlagIgnoreAttackerModifiers,
		ClassSpellMask: MageSpellEvocation,
		ProcMask:       core.Ternary(hasGlyph, core.ProcMaskSpellHealing, core.ProcMaskUnknown),

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 4,
			},
		},

		Hot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Evocation",
			},
			NumberOfTicks:        3,
			TickLength:           time.Second * 2,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.AddMana(sim, mage.MaxMana()*0.15, manaMetrics)
				if hasGlyph {
					// Evo counts as "Tick as Cast"
					dot.Spell.CalcAndDealHealing(sim, target, mage.MaxHealth()*0.10, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
			spell.SelfHot().TickOnce(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: evocation,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if character.HasActiveAuraWithTag(core.InnervateAuraTag) || character.HasActiveAuraWithTag(core.ManaTideTotemAuraTag) {
				return false
			}

			if sim.GetRemainingDuration() < 12*time.Second {
				return false
			}

			return character.CurrentManaPercent() < 0.1
		},
	})
}
