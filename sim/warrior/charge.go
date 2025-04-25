package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterCharge() {
	actionID := core.ActionID{SpellID: 100}
	metrics := warrior.NewRageMetrics(actionID)
	rage := float64(15 + 5*warrior.Talents.Blitz)

	warrior.ChargeAura = warrior.RegisterAura(core.Aura{
		Label:    "Charge",
		ActionID: actionID,
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMovementSpeed(sim, 3.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.MultiplyMovementSpeed(sim, 1.0/3.0)
		},
	})

	warrior.RegisterMovementCallback(func(sim *core.Simulation, position float64, kind core.MovementUpdateType) {
		if kind == core.MovementEnd && warrior.ChargeAura.IsActive() {
			warrior.ChargeAura.Deactivate(sim)
		}
	})

	warrior.Charge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCharge,
		MinRange:       8,
		MaxRange:       25,

		Cast: core.CastConfig{
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
			warrior.ChargeAura.Activate(sim)
			warrior.AddRage(sim, rage, metrics)
			warrior.MoveTo(core.MaxMeleeRange-1, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
		},
	})
}
