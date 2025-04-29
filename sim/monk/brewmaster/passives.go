package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerPassives() {
	bm.registerBrewmasterTraining()
	bm.registerElusiveBrew()
	bm.registerGiftOfTheOx()
}

func (bm *BrewmasterMonk) registerBrewmasterTraining() {
	// Fortifying Brew
	// Also increases your Stagger amount by 20% while active.

	// Tiger Palm
	// Tiger Palm no longer costs Chi, and when you deal damage with Tiger Palm the amount of your next Guard is increased by 15%. Lasts 30 sec.
	// Tiger Palm Chi mod is implemented in tiger_palm.go
	bm.PowerGuardAura = bm.RegisterAura(core.Aura{
		Label:    "Power Guard",
		ActionID: core.ActionID{SpellID: 118636},
		Duration: 30 * time.Second,
	})

	core.MakeProcTriggerAura(&bm.Unit, core.ProcTrigger{
		Name:           "Power Guard Trigger",
		ClassSpellMask: monk.MonkSpellTigerPalm,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bm.PowerGuardAura.Activate(sim)
		},
	})

	// Blackout Kick
	// After you Blackout Kick, you gain Shuffle, increasing your parry chance by 20%
	// and your Stagger amount by an additional 20% for 6 sec.
	// TODO
}
