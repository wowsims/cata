package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (rogue *Rogue) registerSliceAndDice() {
	actionID := core.ActionID{SpellID: 5171}

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

	rogue.SliceAndDiceAura = rogue.RegisterAura(core.Aura{
		Label:    "Slice and Dice",
		ActionID: actionID,
		// This will be overridden on cast, but set a non-zero default so it doesn't crash when used in APL prepull
		Duration: rogue.sliceAndDiceDurations[5],
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, 1+rogue.SliceAndDiceBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, 1/(1+rogue.SliceAndDiceBonus))
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
			rogue.SliceAndDiceAura.Activate(sim)
			rogue.ApplyFinisher(sim, spell)
		},
	})
}
