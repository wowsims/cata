package warlock

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) registerInfernoSpell() {
	duration := time.Second * time.Duration(45+10*warlock.Talents.AncientGrimoire)

	summonInfernalAura := warlock.RegisterAura(core.Aura{
		Label:    "Summon Infernal",
		ActionID: core.ActionID{SpellID: 1122},
		Duration: duration,
	})

	warlock.Inferno = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1122},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 1500,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.765,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := warlock.CalcBaseDamageWithVariance(sim, 0.485, 0.119) * sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			warlock.Infernal.EnableWithTimeout(sim, warlock.Infernal, duration)

			// fake aura to show duration
			summonInfernalAura.Activate(sim)
		},
	})
}

type InfernalPet struct {
	core.Pet
	owner          *Warlock
	immolationAura *core.Spell
}

func (warlock *Warlock) NewInfernal() *InfernalPet {
	statInheritance := func(ownerStats stats.Stats) stats.Stats {
		ownerHitChance := math.Floor(ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance)

		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 1.0,
			stats.AttackPower:      ownerStats[stats.SpellPower] * 0.57,
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.MeleeHit:         ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit:         ownerHitChance * core.SpellHitRatingPerHitChance,
			stats.Expertise: (ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance) *
				PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,
		}
	}

	infernal := &InfernalPet{
		Pet: core.NewPet("Infernal", &warlock.Character, stats.Stats{
			stats.Strength:  331,
			stats.Agility:   113,
			stats.Stamina:   361,
			stats.Intellect: 65,
			stats.Spirit:    109,
			stats.Mana:      0,
			stats.MeleeCrit: 3.192 * core.CritRatingPerCritChance,
		}, statInheritance, false, false),
		owner: warlock,
	}

	infernal.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	infernal.AddStat(stats.AttackPower, -20)

	// infernal is classified as a warrior class, so we assume it gets the
	// same agi crit coefficient
	infernal.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/62.5)

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

	core.ApplyPetConsumeEffects(&infernal.Character, warlock.Consumes)

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
			TickLength:          time.Second * 2,
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
