package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func performCharge(sim *core.Simulation, warrior *Warrior, metrics *core.ResourceMetrics, rage float64) {
	warrior.AddRage(sim, rage, metrics)
	if warrior.DistanceFromTarget > core.MaxMeleeRange {
		warrior.PseudoStats.MovementSpeedMultiplier *= 2
		warrior.MoveTo(core.MaxMeleeRange-1, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
	} else {
		warrior.ChargeAura.Deactivate(sim)
	}
}

func (warrior *Warrior) RegisterCharge() {
	metrics := warrior.NewRageMetrics(core.ActionID{SpellID: 100})
	rage := float64(15 + 5*warrior.Talents.Blitz)
	chargeMinRange := 8.0

	warrior.ChargeAura = warrior.RegisterAura(core.Aura{
		Label:    "Charge",
		ActionID: core.ActionID{SpellID: 100},
		Duration: 5 * time.Second,
	})

	warrior.RegisterMovementCallback(func(sim *core.Simulation, position float64, kind core.MovementUpdateType) {
		if warrior.ChargeAura.IsActive() && kind == core.MovementEnd {
			if warrior.DistanceFromTarget >= chargeMinRange {
				// Has moved out, charge back in
				performCharge(sim, warrior, metrics, rage)
			} else {
				// Has charged back in, reset movement speed
				warrior.PseudoStats.MovementSpeedMultiplier /= 2
				warrior.ChargeAura.Deactivate(sim)
			}
		}
	})

	warrior.Charge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 100},
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCharge,
		MaxRange:       25,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 15,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: core.Ternary(warrior.Talents.Juggernaut || warrior.Talents.Warbringer, nil, func(sim *core.Simulation, target *core.Unit) bool {
			return sim.CurrentTime < 0 && warrior.StanceMatches(BattleStance)
		}),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Don't want to deal with adding tons of "not moving" conditions to the APL atm
			if sim.CurrentTime < 0 {
				performCharge(sim, warrior, metrics, rage)
				return
			}

			warrior.ChargeAura.Activate(sim)

			if warrior.StartDistanceFromTarget <= core.MaxMeleeRange && sim.CurrentTime < 0 && warrior.DistanceFromTarget < chargeMinRange {
				// Always charge from min range in prepull
				warrior.DistanceFromTarget = chargeMinRange
			}

			if warrior.DistanceFromTarget >= chargeMinRange {
				performCharge(sim, warrior, metrics, rage)
				return
			}

			warrior.MoveTo(chargeMinRange, sim)
		},
	})
}
