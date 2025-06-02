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

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           war.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

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
		TimeValue: -10,
	})
}

func (war *ArmsWarrior) registerSuddenDeath() {
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Sudden Death",
		ActionID: core.ActionID{SpellID: 29725},
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) || spell.ActionID.SpellID != StrikesOfOpportunityHitID {
				return
			}
			if sim.Proc(0.1, "Sudden Death") {
				war.ColossusSmash.CD.Reset()
			}
		},
	})

	war.SuddenExecute = war.RegisterAura(core.Aura{
		Label:    "Sudden Execute",
		ActionID: core.ActionID{SpellID: 139958},
		Duration: 10 * time.Second,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskOverpower,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Sudden Execute - Trigger",
		ClassSpellMask: warrior.SpellMaskExecute,
		Outcome:        core.OutcomeLanded,
		Callback:       core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.SuddenExecute.Activate(sim)
		},
	})
}
