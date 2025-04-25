package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) ChangeActivePet(sim *core.Simulation, newPet *WarlockPet) {
	if warlock.ActivePet != nil {
		warlock.ActivePet.Disable(sim)
	}
	newPet.Enable(sim, newPet)
	warlock.ActivePet = newPet
}

func (warlock *Warlock) GetSummonStunAura() core.Aura {
	return core.Aura{
		Label:    "Summoning Disorientation",
		ActionID: core.ActionID{SpellID: 32752},
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.WaitUntil(sim, sim.CurrentTime+5*time.Second)
			aura.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+5*time.Second, false)
		},
	}
}

func (warlock *Warlock) registerSummonDemon() {
	stunActionID := core.ActionID{SpellID: 32752}

	// Summon Felhunter
	warlock.Felhunter.RegisterAura(warlock.GetSummonStunAura())
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 691},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonFelhunter,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 80},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 6 * time.Second,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				warlock.ActivatePetSummonStun(sim, stunActionID)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Deactivate(sim)
			warlock.ChangeActivePet(sim, warlock.Felhunter)
		},
	})

	// Summon Imp
	warlock.Imp.RegisterAura(warlock.GetSummonStunAura())
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 688},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonImp,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 64},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 6 * time.Second,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				warlock.ActivatePetSummonStun(sim, stunActionID)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Deactivate(sim)
			warlock.ChangeActivePet(sim, warlock.Imp)
		},
	})

	warlock.Succubus.RegisterAura(warlock.GetSummonStunAura())
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 712},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonSuccubus,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 80},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 6 * time.Second,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				warlock.ActivatePetSummonStun(sim, stunActionID)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Deactivate(sim)
			warlock.ChangeActivePet(sim, warlock.Succubus)
		},
	})
}

func (warlock *Warlock) ActivatePetSummonStun(sim *core.Simulation, stunActionID core.ActionID) {
	if warlock.ActivePet != nil {
		warlock.ActivePet.GetAuraByID(stunActionID).Activate(sim)
	}
}
