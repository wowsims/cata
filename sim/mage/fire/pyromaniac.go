package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerPyromaniac() {
	fire.pyromaniacAuras = fire.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Pyromaniac",
			ActionID: core.ActionID{SpellID: 132209},
			Duration: time.Second * 15,
		}).AttachDDBC(DDBC_Pyromaniac, DDBC_Total, &fire.AttackTables, fire.pyromaniacDDBCHandler)
	})

	core.MakeProcTriggerAura(&fire.Unit, core.ProcTrigger{
		Name:           "Pyromaniac - Trigger",
		ClassSpellMask: mage.MageSpellLivingBombApply | mage.MageSpellFrostBomb | mage.MageSpellNetherTempest,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			fire.pyromaniacAuras.Get(fire.CurrentTarget).Activate(sim)
		},
	})
}

func (fire *FireMage) pyromaniacDDBCHandler(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
	if spell.Matches(mage.FireSpellIgnitable ^ mage.MageSpellScorch) {
		return 1.1
	}
	return 1.0
}
