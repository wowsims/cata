package warlock

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

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
	core.Pet

	Firebolt *core.Spell
}

func (warlock *Warlock) NewImpPet() *ImpPet {
	baseStats := stats.Stats{
		stats.Health: 40962,
		stats.Mana:   32251,
		//EnableManaBarWithModifier is subtracting 10... we don't want that.
		stats.SpellPower:  10,
		stats.AttackPower: 0,
		stats.Agility:     0,
		stats.Stamina:     0,
		stats.Intellect:   0,
		stats.Strength:    0,
		stats.Spirit:      0,
		stats.MP5:         0, // rough guess, unclear if it's affected by other stats
		stats.MeleeCrit:   3.454 * core.CritRatingPerCritChance,
		stats.SpellCrit:   0.9075 * core.CritRatingPerCritChance,
	}

	imp := &ImpPet{
		Pet: core.NewPet(PetImp, &warlock.Character, baseStats, makeStatInheritance(), false, false),
	}

	imp.EnableManaBarWithModifier(0.33)
	imp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	imp.AddStat(stats.AttackPower, -20)
	imp.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/52.0833)

	core.ApplyPetConsumeEffects(&warlock.Character, warlock.Consumes)

	warlock.AddPet(imp)

	return imp
}

func makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		//starting with wotlk stats and adjusting from there
		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.75,
			stats.Armor:       ownerStats[stats.Armor] * 1.0,
			stats.AttackPower: ownerStats[stats.SpellPower] * 0.57,
			stats.SpellPower:  ownerStats[stats.SpellPower] * 0.15,
			stats.SpellHit:    math.Floor(ownerStats[stats.SpellHit] / 12.0 * 17.0),
		}
	}
}

func (imp *ImpPet) GetPet() *core.Pet {
	return &imp.Pet
}

func (imp *ImpPet) Reset(_ *core.Simulation) {
}

func (imp *ImpPet) Initialize() {
	imp.registerFireboltSpell()
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
			// The .5 seems to be based on the spellpower of the owner. So dividing this by .15 ratio of imp to owner spell power.
			baseDamage := 124 + (0.657 * (0.5 / .15 * spell.SpellPower()))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
