package arms

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) registerMastery() {
	procAttackConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: StrikesOfOpportunityHitID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 0.55,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	}

	procAttack := war.RegisterSpell(procAttackConfig)

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Strikes of Opportunity",
		ActionID: procAttackConfig.ActionID,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		ProcMask: core.ProcMaskMelee,
		ICD:      100 * time.Millisecond,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			// Implement the proc in here so we can get the most up to date proc chance from mastery
			return sim.Proc(war.GetMasteryProcChance(), "Strikes of Opportunity")
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAttack.Cast(sim, result.Target)
		},
	})
}

func (war *ArmsWarrior) registerSeasonedSoldier() {
	actionID := core.ActionID{SpellID: 12712}
	core.MakePermanent(war.RegisterAura(core.Aura{
		Label:    "Seasoned Soldier",
		ActionID: actionID,
		Duration: core.NeverExpires,
	}).AttachMultiplicativePseudoStatBuff(
		&war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.25,
	))

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskThunderClap | warrior.SpellMaskWhirlwind,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -10,
	})
}

func (war *ArmsWarrior) registerSuddenDeath() {

	suddenDeathAura := war.RegisterAura(core.Aura{
		Label:    "Sudden Death",
		ActionID: core.ActionID{SpellID: 52437},
		Duration: 2 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.ColossusSmash.CD.Reset()
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Sudden Death - Trigger",
		ActionID: core.ActionID{SpellID: 29725},
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && spell.ActionID.SpellID != StrikesOfOpportunityHitID {
				return
			}

			if sim.Proc(0.1, "Sudden Death") {
				suddenDeathAura.Activate(sim)
			}
		},
	})

	executeAura := war.RegisterAura(core.Aura{
		Label:    "Sudden Execute",
		ActionID: core.ActionID{SpellID: 139958},
		Duration: 10 * time.Second,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskOverpower,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -2,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Sudden Execute - Trigger",
		ClassSpellMask: warrior.SpellMaskExecute,
		Outcome:        core.OutcomeLanded,
		Callback:       core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			executeAura.Activate(sim)
		},
	})
}

func (war *ArmsWarrior) registerTasteForBlood() {
	actionID := core.ActionID{SpellID: 60503}

	war.TasteForBloodAura = war.RegisterAura(core.Aura{
		Label:     "Taste For Blood",
		ActionID:  actionID,
		Duration:  12 * time.Second,
		MaxStacks: 5,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Taste For Blood: Mortal Strike - Trigger",
		ClassSpellMask: warrior.SpellMaskMortalStrike,
		Outcome:        core.OutcomeLanded,
		Callback:       core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.TasteForBloodAura.Activate(sim)
			war.TasteForBloodAura.SetStacks(sim, war.TasteForBloodAura.GetStacks()+2)
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Taste For Blood: Dodge - Trigger",
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeDodge,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.TasteForBloodAura.Activate(sim)
			war.TasteForBloodAura.AddStack(sim)
		},
	})

}
