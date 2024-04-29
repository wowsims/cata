package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

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
	core.Pet

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

	autoAttackOptions := core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  88.8,
			BaseDamageMax:  133.3,
			SwingSpeed:     2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	}

	felhunter := &FelhunterPet{
		Pet: core.NewPet(PetFelhunter, &warlock.Character, baseStats, warlock.MakeStatInheritance(), false, false),
	}

	felhunter.owner = warlock

	felhunter.EnableManaBarWithModifier(0.77)
	felhunter.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	felhunter.AddStat(stats.AttackPower, -20)
	felhunter.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/52.0833)
	felhunter.EnableAutoAttacks(felhunter, autoAttackOptions)
	core.ApplyPetConsumeEffects(&warlock.Character, warlock.Consumes)
	warlock.AddPet(felhunter)

	return felhunter
}

func (felhunter *FelhunterPet) GetPet() *core.Pet {
	return &felhunter.Pet
}

func (felhunter *FelhunterPet) Reset(_ *core.Simulation) {
}

func (felhunter *FelhunterPet) Initialize() {
	felhunter.registerShadowBiteSpell()
}

func (felhunter *FelhunterPet) ExecuteCustomRotation(sim *core.Simulation) {
	if felhunter.ShadowBite.CanCast(sim, felhunter.CurrentTarget) {
		felhunter.ShadowBite.Cast(sim, felhunter.CurrentTarget)
	}
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

			activeDots := 0

			for _, spell := range felhunter.owner.Spellbook {
				if spell.ClassSpellMask&WarlockDoT > 0 && spell.Dot(target).IsActive() {
					activeDots++
				}
			}

			baseDamage *= 1 + 0.15*float64(activeDots)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}
