package shaman

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type FireElemental struct {
	core.Pet

	FireBlast *core.Spell
	FireNova  *core.Spell

	FireShieldAura *core.Aura

	shamanOwner *Shaman
}

var FireElementalSpellPowerScaling = 0.5883
var FireElementalIntellectScaling = 0.3198

func (shaman *Shaman) NewFireElemental() *FireElemental {
	fireElemental := &FireElemental{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Greater Fire Elemental",
			Owner:           &shaman.Character,
			BaseStats:       fireElementalPetBaseStats,
			StatInheritance: shaman.fireElementalStatInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
		shamanOwner: shaman,
	}
	fireElemental.EnableManaBar()
	fireElemental.AddStatDependency(stats.Intellect, stats.SpellPower, 1.0)
	fireElemental.AddStatDependency(stats.Intellect, stats.AttackPower, 7.0) // 1.0 * 7
	fireElemental.AddStat(stats.SpellPower, -10)
	fireElemental.AddStat(stats.AttackPower, -70) // -10 * 7
	fireElemental.EnableAutoAttacks(fireElemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  429, //Estimated from beta testing
			BaseDamageMax:  463, //Estimated from beta testing
			SwingSpeed:     2,
			CritMultiplier: 2.66, //Estimated from beta testing
			SpellSchool:    core.SpellSchoolFire,
		},
		AutoSwingMelee: true,
	})
	fireElemental.AutoAttacks.MHConfig().BonusCoefficient = 0

	if shaman.Race == proto.Race_RaceDraenei {
		fireElemental.AddStats(stats.Stats{
			stats.PhysicalHitPercent: -1,
			stats.SpellHitPercent:    -1,
			stats.ExpertiseRating:    math.Floor(-core.SpellHitRatingPerHitPercent * 0.79),
		})
	}

	fireElemental.OnPetEnable = fireElemental.enable
	fireElemental.OnPetDisable = fireElemental.disable

	shaman.AddPet(fireElemental)

	return fireElemental
}

func (fireElemental *FireElemental) enable(sim *core.Simulation) {
	fireElemental.ChangeStatInheritance(fireElemental.shamanOwner.fireElementalStatInheritance())
	fireElemental.FireShieldAura.Activate(sim)
}

func (fireElemental *FireElemental) disable(sim *core.Simulation) {
	fireElemental.FireShieldAura.Deactivate(sim)
}

func (fireElemental *FireElemental) GetPet() *core.Pet {
	return &fireElemental.Pet
}

func (fireElemental *FireElemental) Initialize() {

	fireElemental.registerFireBlast()
	fireElemental.registerFireNova()
	fireElemental.registerFireShieldAura()
}

func (fireElemental *FireElemental) Reset(_ *core.Simulation) {

}

func (fireElemental *FireElemental) ExecuteCustomRotation(sim *core.Simulation) {
	/*
		TODO this is a little dirty, can probably clean this up, the rotation might go through some more overhauls,
		the random AI is hard to emulate.
	*/
	target := fireElemental.CurrentTarget

	if fireElemental.FireNova.DefaultCast.Cost > fireElemental.CurrentMana() {
		return
	}

	random := sim.RandomFloat("Fire Elemental Pet Spell")

	//Melee the other 30%
	if random >= .75 {
		fireElemental.TryCast(sim, target, fireElemental.FireBlast)
	} else if random >= .40 && random < 0.75 {
		fireElemental.TryCast(sim, target, fireElemental.FireNova)
	}

	if !fireElemental.GCD.IsReady(sim) {
		return
	}

	minCd := min(fireElemental.FireBlast.CD.ReadyAt(), fireElemental.FireNova.CD.ReadyAt())
	fireElemental.ExtendGCDUntil(sim, max(minCd, sim.CurrentTime+time.Second))

	if !fireElemental.GCD.IsReady(sim) {
		return
	}

	fireElemental.ExtendGCDUntil(sim, sim.CurrentTime+time.Second)
}

func (fireElemental *FireElemental) TryCast(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	if !spell.Cast(sim, target) {
		return false
	}
	// all spell casts reset the elemental's swing timer
	fireElemental.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+spell.CurCast.CastTime, false)
	return true
}

var fireElementalPetBaseStats = stats.Stats{
	stats.Mana:        6803,
	stats.Health:      4903, //Estimated from beta testing
	stats.Intellect:   157,
	stats.Stamina:     0,
	stats.SpellPower:  0, //Estimated
	stats.AttackPower: 0, //Estimated

	// TODO : Log digging shows ~2% melee crit chance, and ~2% spell hit chance + 5% spell crit debuff
	stats.PhysicalCritPercent: 6.8,
	stats.SpellCritPercent:    6.8,
}

func (shaman *Shaman) fireElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerSpellHitPercent := ownerStats[stats.SpellHitPercent]

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.80,                              // Estimated from beta testing
			stats.Intellect:   ownerStats[stats.Intellect] * FireElementalIntellectScaling,   // Estimated from beta testing
			stats.SpellPower:  ownerStats[stats.SpellPower] * FireElementalSpellPowerScaling, // Estimated from beta testing
			stats.AttackPower: ownerStats[stats.SpellPower] * 4.9,                            // 0.7*7 Estimated from beta testing

			stats.PhysicalHitPercent: ownerSpellHitPercent / 17 * 8,
			stats.SpellHitPercent:    ownerSpellHitPercent,

			/*
				TODO working on figuring this out, getting close need more trials. will need to remove specific buffs,
				ie does not gain the benefit from draenei buff.
				Scaled linearly to reach Expertise Soft Cap (26) when Shaman is at 17% Spell Hit Cap
			*/
			stats.ExpertiseRating: math.Floor(ownerSpellHitPercent / 17 * 26 * core.ExpertisePerQuarterPercentReduction),
		}
	}
}
