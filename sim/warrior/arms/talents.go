package arms

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) ApplyTalents() {
	war.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86526)
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

func (war *ArmsWarrior) TriggerSlaughter(sim *core.Simulation, target *core.Unit) {
	if war.Talents.LambsToTheSlaughter == 0 {
		return
	}

	rend := war.Rend.Dot(target)
	if rend != nil && rend.IsActive() {
		rend.Apply(sim)
		rend.TickOnce(sim)
	}

	if !war.slaughter.IsActive() {
		war.slaughter.Activate(sim)
		war.slaughter.AddStack(sim)
	} else {
		war.slaughter.Refresh(sim)
		war.slaughter.AddStack(sim)
	}
}

func (war *ArmsWarrior) applySlaughter() {
	if war.Talents.LambsToTheSlaughter == 0 {
		return
	}

	damageMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskMortalStrike | warrior.SpellMaskExecute | warrior.SpellMaskOverpower | warrior.SpellMaskSlam,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.0,
	})

	war.slaughter = war.RegisterAura(core.Aura{
		Label:     "Slaughter",
		ActionID:  core.ActionID{SpellID: 84586},
		Duration:  time.Second * 15,
		MaxStacks: war.Talents.LambsToTheSlaughter,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks != 0 {
				bonus := 0.1 * float64(newStacks)
				damageMod.UpdateFloatValue(bonus)
				damageMod.Activate()
			} else {
				damageMod.Deactivate()
			}
		},
	})
}

func (war *ArmsWarrior) applyWreckingCrew() {
	if war.Talents.WreckingCrew == 0 {
		return
	}

	effect := 1.0 + (0.05 * float64(war.Talents.WreckingCrew))
	war.wreckingCrew = war.RegisterAura(core.Aura{
		Label:    "Wrecking Crew",
		ActionID: core.ActionID{SpellID: 57519},
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= effect
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= effect
		},
	})

	core.RegisterPercentDamageModifierEffect(war.wreckingCrew, effect)
}

func (war *ArmsWarrior) TriggerWreckingCrew(sim *core.Simulation) {
	if war.Talents.WreckingCrew == 0 {
		return
	}

	procChance := 0.5 * float64(war.Talents.WreckingCrew)
	if sim.Proc(procChance, "Wrecking Crew") {
		war.wreckingCrew.Activate(sim)
	}
}

func (war *ArmsWarrior) applyJuggernaut() {
	if !war.Talents.Juggernaut {
		return
	}
	// Charge cooldown is sharded with intercept, but as intercept is not implemented will ignore that for now
	war.AddStaticMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskCharge,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -2 * time.Second,
	})

	modJuggernaut := war.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 25,
		ClassMask:  warrior.SpellMaskMortalStrike | warrior.SpellMaskSlam,
	})
	actionId := core.ActionID{SpellID: 65156}
	auraJugg := war.RegisterAura(core.Aura{
		Label:    "Juggernaut",
		ActionID: actionId,
		Duration: 10 * time.Second,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			modJuggernaut.Activate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask&warrior.SpellMaskSlam) != 0 || (spell.ClassSpellMask&warrior.SpellMaskMortalStrike) != 0 {
				aura.Deactivate(sim)
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			modJuggernaut.Deactivate()
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Juggernaut Trigger",
		ActionID:       core.ActionID{SpellID: 65156},
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: warrior.SpellMaskCharge,

		ProcChance: 1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			auraJugg.Activate(sim)
		},
	})
}
