package arms

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) ApplyTalents() {
	war.Warrior.ApplyCommonTalents()

	war.RegisterBladestorm()
	war.RegisterDeadlyCalm()
	war.RegisterSweepingStrikes()

	war.applyBloodFrenzy()
	war.applyImpale()
	war.applyImprovedSlam()
	war.applySlaughter()
	war.applySuddenDeath()
	war.applyTasteForBlood()
	war.applyWreckingCrew()
	war.applyJuggernaut()

	// Apply glyphs after talents so we can modify spells added from talents
	war.ApplyGlyphs()
}

func (war *ArmsWarrior) applyTasteForBlood() {
	if war.Talents.TasteForBlood == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskOverpower,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 20.0 * float64(war.Talents.TasteForBlood),
	})

	procChance := []float64{0, 0.33, 0.66, 1}[war.Talents.TasteForBlood]

	// Use a specific aura for TfB so we can track procs
	// Overpower will check for any aura with the EnableOverpowerTag when it tries to cast
	actionID := core.ActionID{SpellID: 60503}
	tfbAura := war.RegisterAura(core.Aura{
		Label:    "Taste for Blood",
		ActionID: actionID,
		Duration: time.Second * 9,
		Tag:      warrior.EnableOverpowerTag,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Taste for Blood Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnPeriodicDamageDealt,
		ClassSpellMask: warrior.SpellMaskRend,
		ICD:            5 * time.Second,
		ProcChance:     procChance,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			tfbAura.Activate(sim)
		},
	})
}

func (war *ArmsWarrior) applyImpale() {
	if war.Talents.Impale == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskMortalStrike | warrior.SpellMaskSlam | warrior.SpellMaskOverpower,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: 0.1 * float64(war.Talents.Impale),
	})
}

func (war *ArmsWarrior) applyImprovedSlam() {
	if war.Talents.ImprovedSlam == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskSlam,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * time.Duration(-500*war.Talents.ImprovedSlam),
	})

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskSlam,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(war.Talents.ImprovedSlam),
	})
}

func (war *ArmsWarrior) applySuddenDeath() {
	if war.Talents.SuddenDeath == 0 {
		return
	}

	procChance := 0.03 * float64(war.Talents.SuddenDeath)
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:       "Sudden Death Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.ColossusSmash.CD.Reset()
		},
	})
}
