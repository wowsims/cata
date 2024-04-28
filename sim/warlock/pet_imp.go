package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) registerSummonImpSpell() {
	warlock.SummonImp = warlock.RegisterSpell(core.SpellConfig{
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
				CastTime: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//warlock.ChangeActivePet(sim, &warlock.Imp.WarlockPet)
			warlock.ChangeActivePet(sim, PetImp)
		},
	})
}

type ImpPet struct {
	*WarlockPet

	Firebolt *core.Spell
}

func (warlock *Warlock) NewImpPet() *ImpPet {
	baseStats := stats.Stats{
		stats.Strength:  297,
		stats.Agility:   79,
		stats.Stamina:   118,
		stats.Intellect: 369,
		stats.Spirit:    367,
		stats.Mana:      1174,
		stats.MP5:       270, // rough guess, unclear if it's affected by other stats
		stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
		stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
	}

	imp := &ImpPet{
		WarlockPet: NewWarlockPet(warlock, PetImp, baseStats, nil),
	}

	imp.EnableManaBarWithModifier(0.33)

	return imp
}

func (imp *ImpPet) Initialize() {
	imp.registerFireboltSpell()
}

func (imp *ImpPet) Reset(_ *core.Simulation) {
}

func (imp *ImpPet) ExecuteCustomRotation(sim *core.Simulation) {
	imp.Firebolt.Cast(sim, imp.CurrentTarget)
}

func (imp *ImpPet) registerFireboltSpell() {
	imp.Firebolt = imp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3110},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellImpFireBolt,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.02,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 146 + (0.657 * (0.5 * spell.SpellPower()))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
