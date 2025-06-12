package fire

// import (
// 	"time"

// 	"github.com/wowsims/mop/sim/core"
// 	"github.com/wowsims/mop/sim/mage"
// )

// func (fire *FireMage) registerPyromaniac() {
// 	fire.pyromaniacAuras = fire.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
// 		return target.GetOrRegisterAura(core.Aura{
// 			Label:    "Pyromaniac",
// 			ActionID: core.ActionID{SpellID: 132209},
// 			Duration: time.Second * 15,
// 		})
// 	})

// 	core.MakeProcTriggerAura(&fire.Unit, core.ProcTrigger{
// 		Name:           "Pyromaniac - Trigger",
// 		ClassSpellMask: mage.MageSpellLivingBomb | mage.MageSpellFrostBomb | mage.MageSpellNetherTempest, // On application only
// 		Callback:       core.CallbackOnSpellHitDealt,
// 		Outcome:        core.OutcomeLanded,
// 		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			fire.pyromaniacAuras.Get(target).Activate(sim)
// 		},
// 	})
// }
