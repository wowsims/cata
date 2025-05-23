package warlock

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (warlock *Warlock) registerSummonInfernal(timer *core.Timer) {
	summonInfernalAura := warlock.RegisterAura(core.Aura{
		Label:    "Summon Infernal",
		ActionID: core.ActionID{SpellID: 1122},
		Duration: time.Duration(45+10*warlock.Talents.AncientGrimoire) * time.Second,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1122},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonInfernal,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 80},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: 1500 * time.Millisecond,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: 10 * time.Minute,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   warlock.DefaultCritMultiplier(),
		BonusCoefficient: 0.76499998569,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Encounter.AOECapMultiplier() *
					warlock.CalcAndRollDamageRange(sim, 0.48500001431, 0.11999999732)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			warlock.Infernal.EnableWithTimeout(sim, warlock.Infernal, spell.RelatedSelfBuff.Duration)
			// fake aura to show duration
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: summonInfernalAura,
	})
}

type InfernalPet struct {
	core.Pet
	owner          *Warlock
	immolationAura *core.Spell
}

func (warlock *Warlock) NewInfernalPet() *InfernalPet {
	statInheritance := func(ownerStats stats.Stats) stats.Stats {
		ownerHitPercent := math.Floor(ownerStats[stats.SpellHitPercent])

		return stats.Stats{
			stats.Stamina:            ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:          ownerStats[stats.Intellect] * 0.3,
			stats.Armor:              ownerStats[stats.Armor] * 1.0,
			stats.AttackPower:        ownerStats[stats.SpellPower] * 0.57,
			stats.SpellPower:         ownerStats[stats.SpellPower] * 0.15,
			stats.SpellPenetration:   ownerStats[stats.SpellPenetration],
			stats.PhysicalHitPercent: ownerHitPercent,
			stats.SpellHitPercent:    ownerHitPercent,
			stats.ExpertiseRating:    ownerStats[stats.SpellHitPercent] * PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,
		}
	}

	infernal := &InfernalPet{
		Pet: core.NewPet(core.PetConfig{
			Name:  "Infernal",
			Owner: &warlock.Character,
			BaseStats: stats.Stats{
				stats.Strength:            331,
				stats.Agility:             113,
				stats.Stamina:             361,
				stats.Intellect:           65,
				stats.Spirit:              109,
				stats.Mana:                0,
				stats.PhysicalCritPercent: 3.192,
			},
			StatInheritance: statInheritance,
			EnabledOnStart:  false,
			IsGuardian:      false,
		}),
		owner: warlock,
	}

	infernal.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	infernal.AddStat(stats.AttackPower, -20)

	// infernal is classified as a warrior class, so we assume it gets the
	// same agi crit coefficient
	infernal.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, 1/62.5)

	// command doesn't apply to infernal
	if warlock.Race == proto.Race_RaceOrc {
		infernal.PseudoStats.DamageDealtMultiplier /= 1.05
	}

	infernal.EnableAutoAttacks(infernal, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  330,
			BaseDamageMax:  494.9,
			SwingSpeed:     2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})
	infernal.AutoAttacks.MHConfig().DamageMultiplier *= 3.2

	warlock.AddPet(infernal)

	return infernal
}

func (infernal *InfernalPet) GetPet() *core.Pet {
	return &infernal.Pet
}

func (infernal *InfernalPet) Initialize() {
	infernal.immolationAura = infernal.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20153},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Immolation",
				ActionID: core.ActionID{SpellID: 19483},
			},
			NumberOfTicks:       31,
			TickLength:          2 * time.Second,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// base formula is 25 + (lvl-50)*0.5 * Warlock_SP*0.2
				// note this scales with the warlocks SP, NOT with the pets
				warlockSP := infernal.owner.Unit.GetStat(stats.SpellPower)
				baseDmg := (40 + warlockSP*0.2) * sim.Encounter.AOECapMultiplier()

				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}

func (infernal *InfernalPet) Reset(_ *core.Simulation) {
}

func (infernal *InfernalPet) ExecuteCustomRotation(sim *core.Simulation) {
	infernal.immolationAura.Cast(sim, nil)
}
