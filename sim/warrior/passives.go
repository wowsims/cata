package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (war *Warrior) registerEnrage() {
	actionID := core.ActionID{SpellID: 12880}
	var bonusSnapshot float64
	war.EnrageAura = war.RegisterAura(core.Aura{
		Label:    "Enrage",
		Tag:      EnrageTag,
		ActionID: actionID,
		Duration: 6 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusSnapshot = 1.0 + (0.1 * war.EnrageEffectMultiplier)
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= bonusSnapshot
			core.RegisterPercentDamageModifierEffect(aura, bonusSnapshot)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= bonusSnapshot
		},
	})

	core.RegisterPercentDamageModifierEffect(war.EnrageAura, 1+0.1)

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Enrage Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskColossusSmash | SpellMaskShieldSlam | SpellMaskDevastate | SpellMaskBloodthirst | SpellMaskMortalStrike,
		Outcome:        core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.EnrageAura.Deactivate(sim)
			war.EnrageAura.Activate(sim)
		},
	})
}
