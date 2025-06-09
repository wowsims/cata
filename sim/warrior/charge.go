package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (war *Warrior) registerCharge() {
	actionID := core.ActionID{SpellID: 100}
	metrics := war.NewRageMetrics(actionID)
	var chargeRageGenCD time.Duration

	hasRageGlyph := war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfBullRush)
	hasRangeGlyph := war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfLongCharge)

	chargeRageGain := 20 + core.TernaryFloat64(hasRageGlyph, 15, 0)
	chargeRange := 25 + core.TernaryFloat64(hasRangeGlyph, 5, 0)

	aura := war.RegisterAura(core.Aura{
		Label:    "Charge",
		ActionID: actionID,
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.MultiplyMovementSpeed(sim, 3.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.MultiplyMovementSpeed(sim, 1.0/3.0)
		},
	})

	war.RegisterMovementCallback(func(sim *core.Simulation, position float64, kind core.MovementUpdateType) {
		if kind == core.MovementEnd && aura.IsActive() {
			aura.Deactivate(sim)
		}
	})

	war.RegisterResetEffect(func(sim *core.Simulation) {
		chargeRageGenCD = 0
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCharge,
		MinRange:       8,
		MaxRange:       chargeRange,
		Charges:        core.TernaryInt(war.Talents.DoubleTime, 2, 0),
		RechargeTime:   core.TernaryDuration(war.Talents.DoubleTime, time.Second*20, 0),

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 20 * time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
			if !war.Talents.DoubleTime || chargeRageGenCD == 0 || sim.CurrentTime-chargeRageGenCD >= 12*time.Second {
				chargeRageGenCD = sim.CurrentTime
				war.AddRage(sim, chargeRageGain*war.GetRageMultiplier(target), metrics)
			}
			war.MoveTo(core.MaxMeleeRange-1, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
		},
	})
}
