package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) ChangeActivePet(sim *core.Simulation, newPet *WarlockPet) {
	for _, pet := range warlock.Pets {
		if !pet.IsGuardian() && pet.IsEnabled() {
			pet.Disable(sim)
		}
	}

	newPet.Enable(sim, newPet)
}

func (warlock *Warlock) registerSummonFelHunterSpell() {
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 691},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonFelhunter,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.80,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 6 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Deactivate(sim)
			warlock.ChangeActivePet(sim, warlock.Felhunter)
		},
	})
}

func (warlock *Warlock) registerSummonImpSpell() {
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 688},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonImp,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.64,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 6 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Deactivate(sim)
			warlock.ChangeActivePet(sim, warlock.Imp)
		},
	})
}

func (warlock *Warlock) registerSummonSuccubusSpell() {
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 712},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonSuccubus,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.80,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 6 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Deactivate(sim)
			warlock.ChangeActivePet(sim, warlock.Succubus)
		},
	})
}
