package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) registerSummonFelHunterSpell() {
	warlock.SummonImp = warlock.RegisterSpell(core.SpellConfig{
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
				CastTime: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//warlock.ChangeActivePet(sim, &warlock.Imp.WarlockPet)
			warlock.ChangeActivePet(sim, PetImp)
		},
	})
}

type FelhunterPet struct {
	*WarlockPet

	owner *Warlock

	ShadowBite *core.Spell
}

func (warlock *Warlock) NewFelhunterPet() *FelhunterPet {
	baseStats := stats.Stats{
		stats.Strength:  314,
		stats.Agility:   90,
		stats.Stamina:   328,
		stats.Intellect: 150,
		stats.Spirit:    209,
		stats.Mana:      1559,
		stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
		stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
	}

	autoAttackOptions := &core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  88.8,
			BaseDamageMax:  133.3,
			SwingSpeed:     2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	}

	felhunter := &FelhunterPet{
		WarlockPet: NewWarlockPet(warlock, PetFelhunter, baseStats, autoAttackOptions),
	}

	felhunter.owner = warlock

	felhunter.EnableManaBarWithModifier(0.77)

	return felhunter
}

func (felhunter *FelhunterPet) Initialize() {
	felhunter.registerShadowBiteSpell()
}

func (felhunter *FelhunterPet) ExecuteCustomRotation(sim *core.Simulation) {
	felhunter.ShadowBite.Cast(sim, felhunter.CurrentTarget)
}

func (felhunter *FelhunterPet) registerShadowBiteSpell() {
	felhunter.ShadowBite = felhunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 54049},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellFelHunterShadowBite,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    felhunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 166 + (1.228 * (0.5 * spell.SpellPower()))

			owner := felhunter.owner
			spells := []*core.Spell{
				owner.UnstableAffliction,
				owner.Immolate,
				owner.BaneOfAgony,
				owner.BaneOfDoom,
				owner.Corruption,
				owner.Seed,
				owner.DrainSoul,
				// missing: drain life, shadowflame
			}
			counter := 0
			for _, spell := range spells {
				if spell != nil && spell.Dot(target).IsActive() {
					counter++
				}
			}

			baseDamage *= 1 + 0.15*float64(counter)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}
