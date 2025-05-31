package fury

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (war *FuryWarrior) registerCrazedBerserker() {

	war.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMeleeOH,
		FloatValue: 0.25,
	})
	war.AutoAttacks.MHConfig().DamageMultiplier *= 1.1
	war.AutoAttacks.OHConfig().DamageMultiplier *= 1.1
}

func (war *FuryWarrior) registerFlurry() {

	atkSpeedBonus := 1.0 + 0.25
	flurryAura := war.RegisterAura(core.Aura{
		Label:     "Flurry",
		ActionID:  core.ActionID{SpellID: 12968},
		Duration:  15 * time.Second,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.MultiplyMeleeSpeed(sim, atkSpeedBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.MultiplyAttackSpeed(sim, 1.0/atkSpeedBonus)
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Flurry Trigger",
		ActionID: core.ActionID{SpellID: 12972},
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(0.09, "Flurry") {
				flurryAura.Activate(sim)
				flurryAura.SetStacks(sim, flurryAura.MaxStacks)
				return
			}
			if flurryAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				flurryAura.RemoveStack(sim)
			}
		},
	})
}
