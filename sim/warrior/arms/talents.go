package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ArmsWarrior) ApplyTalents() {
	war.Warrior.ApplyTalents()

	war.RegisterBladestorm()
	war.RegisterBloodFrenzy()
	war.RegisterDeadlyCalm()
	war.RegisterSlaughter()
	war.RegisterSuddenDeath()
	war.RegisterSweepingStrikes()
	war.RegisterTasteForBlood()
	war.RegisterWreckingCrew()
}

func (war *ArmsWarrior) RegisterTasteForBlood() {
	if war.Talents.TasteForBlood == 0 {
		return
	}

	procChance := []float64{0, 0.33, 0.66, 1}[war.Talents.TasteForBlood]

	icd := core.Cooldown{
		Timer:    war.NewTimer(),
		Duration: time.Second * 5,
	}

	// Use a specific aura for TfB so we can track procs
	// Overpower will check for any aura with the EnableOverpowerTag when it tries to cast
	tfbAura := war.RegisterAura(core.Aura{
		Label:    "Taste for Blood",
		ActionID: core.ActionID{SpellID: 60503},
		Duration: time.Second * 9,
		Tag:      warrior.EnableOverpowerTag,
	})

	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Taste for Blood Monitor",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != war.Rend {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Taste for Blood") < procChance {
				icd.Use(sim)
				tfbAura.Activate(sim)
			}
		},
	}))
}

func (war *ArmsWarrior) RegisterSuddenDeath() {
	if war.Talents.SuddenDeath == 0 {
		return
	}

	procChance := 0.03 * float64(war.Talents.SuddenDeath)
	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Sudden Death Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if sim.RandomFloat("Sudden Death") < procChance {
				war.ColossusSmash.CD.Reset()
			}
		},
	}))
}

func (war *ArmsWarrior) TriggerSlaughter(sim *core.Simulation, target *core.Unit) {
	if war.Talents.LambsToTheSlaughter == 0 {
		return
	}

	rend := war.Rend.Dot(target)
	if rend != nil && rend.IsActive() {
		rend.Refresh(sim)
	}

	if !war.slaughter.IsActive() {
		war.slaughter.Activate(sim)
	} else {
		war.slaughter.Refresh(sim)
		war.slaughter.AddStack(sim)
	}
}

func (war *ArmsWarrior) RegisterSlaughter() {
	if war.Talents.LambsToTheSlaughter == 0 {
		return
	}

	war.slaughter = war.RegisterAura(core.Aura{
		Label:     "Slaughter",
		ActionID:  core.ActionID{SpellID: 84586},
		Duration:  time.Second * 15,
		MaxStacks: war.Talents.LambsToTheSlaughter,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			// This being negative is a valid case as that means the aura is expiring (newStacks == 0)
			// so we should subtract whatever bonus had been applied
			diff := newStacks - oldStacks
			bonus := 0.1 * float64(diff)
			war.mortalStrike.DamageMultiplierAdditive += bonus
			war.Execute.DamageMultiplierAdditive += bonus
			war.Overpower.DamageMultiplierAdditive += bonus
			war.Slam.DamageMultiplierAdditive += bonus
		},
	})
}

func (war *ArmsWarrior) RegisterWreckingCrew() {
	if war.Talents.WreckingCrew == 0 {
		return
	}

	effect := 1.0 + (0.05 * float64(war.Talents.WreckingCrew))
	war.wreckingCrew = war.RegisterAura(core.Aura{
		Label:    "Wrecking Crew",
		ActionID: core.ActionID{SpellID: 56611},
		Duration: time.Second * 12,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] *= effect
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] /= effect
		},
	})

	core.RegisterPercentDamageModifierEffect(war.wreckingCrew, effect)
}

func (war *ArmsWarrior) TriggerWreckingCrew(sim *core.Simulation) {
	if war.Talents.WreckingCrew == 0 {
		return
	}

	procChance := 0.5 * float64(war.Talents.WreckingCrew)
	if sim.RandomFloat("Wrecking Crew") < procChance {
		war.wreckingCrew.Activate(sim)
	}
}
