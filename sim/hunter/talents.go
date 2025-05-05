package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) applyBlinkStrike() {
	if !hunter.Talents.BlinkStrikes {
		return
	}
	hunter.Pet.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterPetFocusDump,
		FloatValue: 0.5,
	})
}
func (hunter *Hunter) applyThrillOfTheHunt() {
	if !hunter.Talents.ThrillOfTheHunt {
		return
	}

	actionID := core.ActionID{SpellID: 109306}
	procChance := 0.30

	tothAura := hunter.RegisterAura(core.Aura{
		Label:     "Thrill of the Hunt",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.MultiShot.Cost.PercentModifier -= 50
			hunter.ArcaneShot.Cost.PercentModifier -= 50

		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.MultiShot.Cost.PercentModifier += 50
			hunter.ArcaneShot.Cost.PercentModifier += 50

		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == hunter.MultiShot || spell == hunter.ArcaneShot {
				aura.RemoveStack(sim)

			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Thrill of the Hunt Proccer",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Either currcast or base cast idk
			if spell.CurCast.Cost <= 0 {
				return
			}

			if sim.RandomFloat("Thrill of the Hunt") < procChance {
				tothAura.Activate(sim)
				tothAura.SetStacks(sim, 2)
			}
		},
	})
}
