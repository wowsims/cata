package windwalker

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (ww *WindwalkerMonk) registerTigereyeBrew() {
	buffActionID := core.ActionID{SpellID: 116740}
	stackActionID := core.ActionID{SpellID: 125195}

	ww.Monk.RegisterOnChiSpent(func(sim *core.Simulation, chiSpent int32) {
		accumulatedChi := ww.outstandingChi + chiSpent

		for accumulatedChi >= 4 {
			ww.AddBrewStacks(sim, 1)
			accumulatedChi -= 4
		}

		ww.outstandingChi = accumulatedChi
	})

	ww.TigereyeBrewStackAura = ww.RegisterAura(core.Aura{
		Label:     "Tigereye Brew Stacks" + ww.Label,
		ActionID:  stackActionID,
		Duration:  time.Minute * 2,
		MaxStacks: 20,
	})

	ww.Monk.RegisterOnNewBrewStacks(func(sim *core.Simulation, stacksToAdd int32) {
		ww.TigereyeBrewStackAura.Activate(sim)

		// Mastery: Bottled Fury
		// When you generate Tigereye Brew charges, you have a chance to generate an additional charge.
		// Can go above 100% and should then add the correct amount of guaranteed stacks.
		procChance := ww.getMasteryPercent()
		if sim.Proc(math.Mod(procChance, 1), "Mastery: Bottled Fury") {
			stacksToAdd += int32(math.Ceil(procChance))
		} else {
			stacksToAdd += int32(math.Floor(procChance))
		}

		ww.TigereyeBrewStackAura.SetStacks(sim, ww.TigereyeBrewStackAura.GetStacks()+stacksToAdd)
	})

	var damageMultiplier float64
	buffAura := ww.RegisterAura(core.Aura{
		Label:    "Tigereye Brew Buff" + ww.Label,
		ActionID: buffActionID,
		Duration: time.Second * 15,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			stacksToConsume := min(10, ww.TigereyeBrewStackAura.GetStacks())

			damageMultiplier = 1 + 0.06*float64(stacksToConsume)

			ww.PseudoStats.DamageDealtMultiplier *= damageMultiplier

			ww.TigereyeBrewStackAura.SetStacks(sim, ww.TigereyeBrewStackAura.GetStacks()-stacksToConsume)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ww.PseudoStats.DamageDealtMultiplier /= damageMultiplier
		},
	})

	ww.RegisterSpell(core.SpellConfig{
		ActionID:       buffActionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellTigereyeBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    ww.NewTimer(),
				Duration: time.Second * 5,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return ww.TigereyeBrewStackAura.GetStacks() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			buffAura.Activate(sim)
		},
	})
}
