package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) registerSummonSuccubusSpell() {
	warlock.SummonImp = warlock.RegisterSpell(core.SpellConfig{
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
				CastTime: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//warlock.ChangeActivePet(sim, &warlock.Succubus.WarlockPet)
			warlock.ChangeActivePet(sim, PetSuccubus)
		},
	})
}

type SuccubusPet struct {
	*WarlockPet

	LashOfPain *core.Spell
}

func (warlock *Warlock) NewSuccubusPet() *SuccubusPet {
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

	succubus := &SuccubusPet{
		WarlockPet: NewWarlockPet(warlock, PetSuccubus, baseStats, autoAttackOptions),
	}

	succubus.EnableManaBarWithModifier(0.77)

	return succubus
}

func (succubus *SuccubusPet) Initialize() {

	succubus.registerLashOfPainSpell()
}

func (succubus *SuccubusPet) ExecuteCustomRotation(sim *core.Simulation) {
	if !succubus.LashOfPain.IsReady(sim) {
		succubus.WaitUntil(sim, succubus.LashOfPain.CD.ReadyAt())
		return
	}
}

func (succubus *SuccubusPet) registerLashOfPainSpell() {
	succubus.LashOfPain = succubus.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7814},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellSuccubusLashOfPain,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 187 + (0.612 * (0.5 * spell.SpellPower()))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
