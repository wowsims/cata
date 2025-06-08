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
		Flags:       core.SpellFlagPassiveSpell,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// Summon Bloodworm
			i := 0
			max := len(bdk.Bloodworm)
			for ; i < max; i++ {
				if !bdk.Bloodworm[i].IsActive() {
					break
				}
			}
			if i == max {
				// No free worms - increase cap
				return
			}
			bdk.Bloodworm[i].EnableWithTimeout(sim, bdk.Bloodworm[i], time.Second*20)
			bdk.Bloodworm[i].CancelGCDTimer(sim)
		},
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:       "Blood Parasite Trigger" + bdk.Label,
		ActionID:   core.ActionID{SpellID: 49542},
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 5,
		ProcChance: 0.1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodParasiteSpell.Cast(sim, result.Target)
		},
	})
}
