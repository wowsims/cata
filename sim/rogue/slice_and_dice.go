package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerSliceAndDice() {
	actionID := core.ActionID{SpellID: 5171}

	rogue.SliceAndDiceBonusFlat = 0.4
	durationMultiplier := 1.0 + 0.25*float64(rogue.Talents.ImprovedSliceAndDice)
	durationBonus := time.Duration(0)
	if rogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfSliceAndDice) {
		durationBonus += time.Second * 6
	}
	rogue.sliceAndDiceDurations = [6]time.Duration{
		0,
		time.Duration(float64(time.Second*9+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*12+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*15+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*18+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*21+durationBonus) * durationMultiplier),
	}

	var slideAndDiceMod float64
	rogue.SliceAndDiceAura = rogue.RegisterAura(core.Aura{
		Label:    "Slice and Dice",
		ActionID: actionID,
		// This will be overridden on cast, but set a non-zero default so it doesn't crash when used in APL prepull
		Duration: rogue.sliceAndDiceDurations[5],
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryBonus := core.TernaryFloat64(rogue.Spec == proto.Spec_SpecSubtletyRogue, rogue.GetMasteryBonus(), 0)
			slideAndDiceMod = 1 + rogue.SliceAndDiceBonusFlat*(1+masteryBonus)
			rogue.MultiplyMeleeSpeed(sim, slideAndDiceMod)
			if sim.Log != nil {
				rogue.Log(sim, "[DEBUG]: Slice and Dice attack speed mod: %v", slideAndDiceMod)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, 1/slideAndDiceMod)
		},
	})

	rogue.SliceAndDice = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellSliceAndDice,

		EnergyCost: core.EnergyCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[rogue.ComboPoints()]
			rogue.ApplyFinisher(sim, spell)
			rogue.SliceAndDiceAura.Activate(sim)
		},
	})
}
