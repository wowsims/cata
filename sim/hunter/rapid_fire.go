package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerRapidFireCD() {
	actionID := core.ActionID{SpellID: 3045}

	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 53232})
	var focusPA *core.PendingAction

	hasteMultiplier := 1.4 + core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfRapidFire), 0.1, 0)

	hunter.RapidFireAura = hunter.RegisterAura(core.Aura{
		Label:    "Rapid Fire",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, hasteMultiplier)

			focusPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 3,
				NumTicks: 5,
				OnAction: func(sim *core.Simulation) {
					if hunter.Talents.RapidRecuperation > 0 {
						hunter.AddFocus(sim, 6*float64(hunter.Talents.RapidRecuperation), focusMetrics)
					}
				},
			})

		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			focusPA.Cancel(sim)
			aura.Unit.MultiplyRangedSpeed(sim, 1/hasteMultiplier)
		},
	})

	hunter.RapidFire = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute*5 - time.Minute*time.Duration(hunter.Talents.Posthaste),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFireAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.RapidFire,
		Type:  core.CooldownTypeDPS,
	})
}
