package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerEvocation() {

	if mage.Talents.RuneOfPower {
		return
	}

	actionID := core.ActionID{SpellID: 12051}
	manaMetrics := mage.NewManaMetrics(actionID)
	manaPerTick := 0.0

	evocation := mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: MageSpellEvocation,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		Hot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Evocation",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mage.invocationAura.Activate(sim)
				},
			},
			NumberOfTicks:        3,
			TickLength:           time.Second * 2,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.AddMana(sim, manaPerTick, manaMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			manaPerTick = mage.MaxMana() * 0.15
			spell.SelfHot().Apply(sim)
			spell.SelfHot().TickOnce(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: evocation,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if mage.invocationAura.TimeActive(sim) >= time.Duration(time.Second*55) {
				return true
			}
			return !mage.invocationAura.IsActive()
		},
	})
}
