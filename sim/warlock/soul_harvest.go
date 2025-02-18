package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerSoulHarvest() {

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 79268},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagChanneled | core.SpellFlagAPL,

		Cast: core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},
		ExtraCastCondition: func(sim *core.Simulation, _ *core.Unit) bool {
			return sim.CurrentTime <= 0 // only usable outside of combat
		},
		Hot: core.DotConfig{
			SelfOnly:            true,
			Aura:                core.Aura{Label: "Soul Harvest"},
			NumberOfTicks:       9,
			TickLength:          1 * time.Second,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealHealing(sim, target, warlock.MaxHealth()*0.05, dot.Spell.OutcomeMagicHit)
				remainingTicks := dot.RemainingTicks()
				if remainingTicks%3 == 0 {
					warlock.AddSoulShard()
					if sim.Log != nil {
						warlock.Log(sim, "Gained 1 soul shard (%v -> %v)", warlock.SoulShards-1, warlock.SoulShards)
					}
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})
}
