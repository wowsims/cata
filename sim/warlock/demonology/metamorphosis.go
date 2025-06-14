package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (demo *DemonologyWarlock) registerMetamorphosis() {
	metaActionId := core.ActionID{SpellID: 103958}
	var queueMetaCost func(sim *core.Simulation)
	var soulFireManaCost core.ResourceCostImpl
	var drainLifeManaCost core.ResourceCostImpl

	metaAura := demo.RegisterAura(core.Aura{
		Label:    "Metamorphosis",
		ActionID: metaActionId,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			queueMetaCost(sim)

			// update cast cost
			soulFireManaCost = demo.Soulfire.Cost.ResourceCostImpl
			drainLifeManaCost = demo.DrainLife.Cost
			demo.Soulfire.Cost.ResourceCostImpl = NewDemonicFuryCost(160)
			demo.DrainLife.Cost.ResourceCostImpl = NewDemonicFuryCost(0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			demo.Soulfire.Cost.ResourceCostImpl = soulFireManaCost
			demo.DrainLife.Cost.ResourceCostImpl = drainLifeManaCost
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		TimeValue: -time.Millisecond * 500,
		ClassMask: warlock.WarlockSpellSummonDoomguard | warlock.WarlockSpellSummonInfernal | warlock.WarlockSpellCarrionSwarm | warlock.WarlockSpellLifeTap,
	})

	queueMetaCost = func(sim *core.Simulation) {
		pa := core.PendingAction{
			NextActionAt: sim.CurrentTime + time.Second,
			Priority:     core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				if !metaAura.IsActive() {
					return
				}

				demo.DemonicFury.SpendUpTo(sim, core.TernaryInt32(demo.T15_2pc.IsActive(), 4, 6), metaActionId)
				if demo.DemonicFury.Value() < 50 {
					metaAura.Deactivate(sim)
					return
				}

				queueMetaCost(sim)
			},
		}

		sim.AddPendingAction(&pa)
	}

	demo.Metamorphosis = demo.RegisterSpell(core.SpellConfig{
		ActionID:    metaActionId,
		Flags:       core.SpellFlagAPL | core.SpellFlagNoOnCastComplete,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},

			CD: core.Cooldown{
				Timer:    demo.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !metaAura.IsActive() && demo.DemonicFury.Value() >= 50
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			metaAura.Activate(sim)
		},

		RelatedSelfBuff: metaAura,
	})
}
