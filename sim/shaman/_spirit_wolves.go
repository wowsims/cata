package shaman

import (
	"math"
	"strconv"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type SpiritWolf struct {
	core.Pet

	shamanOwner *Shaman
}

type SpiritWolves struct {
	SpiritWolf1 *SpiritWolf
	SpiritWolf2 *SpiritWolf
}

func (SpiritWolves *SpiritWolves) EnableWithTimeout(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.EnableWithTimeout(sim, SpiritWolves.SpiritWolf1, time.Second*30)
	SpiritWolves.SpiritWolf2.EnableWithTimeout(sim, SpiritWolves.SpiritWolf2, time.Second*30)
}

func (SpiritWolves *SpiritWolves) CancelGCDTimer(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.CancelGCDTimer(sim)
	SpiritWolves.SpiritWolf2.CancelGCDTimer(sim)
}

var spiritWolfBaseStats = stats.Stats{
	stats.Stamina:   389,
	stats.Spirit:    116,
	stats.Intellect: 69,
	stats.Armor:     12310,

	stats.Agility:     1218,
	stats.Strength:    476,
	stats.AttackPower: -20,

	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.PhysicalCritPercent: 1.1515 + 1.8,
}

func (shaman *Shaman) NewSpiritWolf(index int) *SpiritWolf {
	spiritWolf := &SpiritWolf{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Spirit Wolf " + strconv.Itoa(index),
			Owner:           &shaman.Character,
			BaseStats:       spiritWolfBaseStats,
			StatInheritance: shaman.makeStatInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      false,
		}),
		shamanOwner: shaman,
	}

	spiritWolf.EnableAutoAttacks(spiritWolf, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  583,
			BaseDamageMax:  876,
			SwingSpeed:     1.5,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	spiritWolf.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	spiritWolf.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[proto.Class_ClassWarrior])

	shaman.AddPet(spiritWolf)

	return spiritWolf
}

const PetExpertiseScale = 3.25

func (shaman *Shaman) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		flooredOwnerHitPercent := math.Floor(ownerStats[stats.PhysicalHitPercent])

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3189,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.AttackPower] * (core.TernaryFloat64(shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfFeralSpirit), 0.6296, 0.3296)),

			stats.HitRating:       flooredOwnerHitPercent * core.PhysicalHitRatingPerHitPercent,
			stats.ExpertiseRating: math.Floor(flooredOwnerHitPercent*PetExpertiseScale) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}

func (spiritWolf *SpiritWolf) Initialize() {
	// Nothing
}

func (spiritWolf *SpiritWolf) ExecuteCustomRotation(_ *core.Simulation) {
}

func (spiritWolf *SpiritWolf) Reset(sim *core.Simulation) {
	spiritWolf.Disable(sim)
	if sim.Log != nil {
		spiritWolf.Log(sim, "Base Stats: %s", spiritWolfBaseStats)
		inheritedStats := spiritWolf.shamanOwner.makeStatInheritance()(spiritWolf.shamanOwner.GetStats())
		spiritWolf.Log(sim, "Inherited Stats: %s", inheritedStats)
		spiritWolf.Log(sim, "Total Stats: %s", spiritWolf.GetStats())
	}
}

func (spiritWolf *SpiritWolf) GetPet() *core.Pet {
	return &spiritWolf.Pet
}
