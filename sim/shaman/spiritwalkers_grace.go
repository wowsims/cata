package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) spiritwalkersGraceBaseDuration() time.Duration {
	return 15 * time.Second
}

func (shaman *Shaman) registerSpiritwalkersGraceSpell() {
	actionID := core.ActionID{SpellID: 79206}

	castWhileMovingMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_AllowCastWhileMoving,
		ClassMask: SpellMaskNone,
	})

	shaman.SpiritwalkersGraceAura = shaman.RegisterAura(core.Aura{
		Label:    "Spiritwalker's Grace" + shaman.Label,
		ActionID: actionID,
		Duration: 15 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castWhileMovingMod.Activate()
			if shaman.hasT13Resto4pc() {
				shaman.SpiritwalkersVestments4PT13Aura.Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castWhileMovingMod.Deactivate()
		},
	})

	shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 120,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.SpiritwalkersGraceAura.Activate(sim)
		},
	})
}
