package survival

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

// Todo: Should we support precasting freezing/ice trap?
func (hunter *SurvivalHunter) applyLNL() {
	if hunter.Spec != proto.Spec_SpecSurvivalHunter {
		return
	}

	actionID := core.ActionID{SpellID: 56343}
	procChance := 0.20

	icd := core.Cooldown{
		Timer:    hunter.NewTimer(),
		Duration: time.Second * 10,
	}

	hunter.LockAndLoadAura = hunter.RegisterAura(core.Aura{
		Icd:       &icd,
		Label:     "Lock and Load Proc",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.Cost.PercentModifier -= 100
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.Cost.PercentModifier += 100
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == hunter.ExplosiveShot {
				hunter.ExplosiveShot.CD.Reset()
				// Weird check but..
				if !aura.Unit.HasActiveAura("Burning Adrenaline") {
					aura.RemoveStack(sim)
				}
			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Lock and Load",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.BlackArrow {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Lock and Load") < procChance {
				icd.Use(sim)
				hunter.LockAndLoadAura.Activate(sim)
				hunter.LockAndLoadAura.SetStacks(sim, 2)
				if hunter.ExplosiveShot != nil {
					hunter.ExplosiveShot.CD.Reset()
				}
			}
		},
	})
}
