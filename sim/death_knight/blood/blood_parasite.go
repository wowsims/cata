package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Your melee attacks have a 10% chance to spawn a Bloodworm.
The Bloodworm attacks your enemies, gorging itself with blood until it bursts to heal nearby allies.
Lasts up to 20 sec.
(Proc chance: 10%, 5s cooldown)
*/
func (bdk *BloodDeathKnight) registerBloodParasite() {
	bloodParasiteSpell := bdk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50452},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// Summon Bloodworm
			for _, worm := range bdk.Bloodworm {
				if worm.IsActive() {
					continue
				}

				worm.EnableWithTimeout(sim, worm, time.Second*20)
				worm.CancelGCDTimer(sim)

				return
			}

			if sim.Log != nil {
				bdk.Log(sim, "No Bloodworm available for Blood Parasite to proc, this is unreasonable.")
			}
		},
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:       "Blood Parasite Trigger" + bdk.Label,
		ActionID:   core.ActionID{SpellID: 49542},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 5,
		ProcChance: 0.1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodParasiteSpell.Cast(sim, result.Target)
		},
	})
}
