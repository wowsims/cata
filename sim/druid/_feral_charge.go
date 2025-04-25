package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerCatCharge() {
	druid.CatCharge = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49376},
		Flags:    core.SpellFlagAPL,
		MinRange: 8,
		MaxRange: 25,

		EnergyCost: core.EnergyCostOptions{
			Cost: 10,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFeralCharge), time.Second*28, time.Second*30),
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return !druid.PseudoStats.InFrontOfTarget && !druid.CannotShredTarget
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// Leap speed is around 80 yards/second according to measurements
			// from boЯsch. This is too fast to be modeled accurately using
			// movement aura stacks, so do it directly here by setting the
			// position to 0 instantaneously but introducing a GCD delay based
			// on the distance traveled.
			travelTime := core.DurationFromSeconds(druid.DistanceFromTarget / 80)
			druid.ExtendGCDUntil(sim, max(druid.NextGCDAt(), sim.CurrentTime+travelTime))
			druid.DistanceFromTarget = 0
			druid.MoveDuration(travelTime, sim)

			// Measurements from boЯsch indicate that while travel speed (and
			// therefore special ability delays) is fairly consistent, there
			// is an additional variable delay on auto-attacks after landing,
			// likely due to the server needing to perform positional checks.
			minAutoDelaySeconds := 0.150
			autoDelaySpreadSeconds := 0.6
			randomDelayTime := core.DurationFromSeconds(minAutoDelaySeconds + sim.RandomFloat("Cat Charge")*autoDelaySpreadSeconds)

			druid.AutoAttacks.CancelMeleeSwing(sim)
			pa := &core.PendingAction{
				NextActionAt: sim.CurrentTime + travelTime + randomDelayTime,
				Priority:     core.ActionPriorityDOT,

				OnAction: func(sim *core.Simulation) {
					druid.AutoAttacks.EnableMeleeSwing(sim)
				},
			}
			sim.AddPendingAction(pa)

			if druid.StampedeCatAura != nil {
				druid.StampedeCatAura.Activate(sim)
			}
		},
	})
}
