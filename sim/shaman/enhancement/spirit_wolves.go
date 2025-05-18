package enhancement

import (
	"strconv"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type SpiritWolf struct {
	core.Pet

	shamanOwner *EnhancementShaman
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
	stats.Stamina: 3137,
}

func (enh *EnhancementShaman) NewSpiritWolf(index int) *SpiritWolf {
	spiritWolf := &SpiritWolf{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Spirit Wolf " + strconv.Itoa(index),
			Owner:                           &enh.Character,
			BaseStats:                       spiritWolfBaseStats,
			StatInheritance:                 enh.makeStatInheritance(),
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
		shamanOwner: enh,
	}
	baseMeleeDamage := enh.CalcScalingSpellDmg(0.5)
	spiritWolf.EnableAutoAttacks(spiritWolf, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  baseMeleeDamage,
			BaseDamageMax:  baseMeleeDamage,
			SwingSpeed:     1.5,
			CritMultiplier: spiritWolf.DefaultCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	enh.AddPet(spiritWolf)

	return spiritWolf
}

func (enh *EnhancementShaman) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.Stamina:             ownerStats[stats.Stamina] * 0.3,
			stats.AttackPower:         ownerStats[stats.AttackPower] * 0.5,
			stats.PhysicalHitPercent:  ownerStats[stats.PhysicalHitPercent],
			stats.SpellHitPercent:     ownerStats[stats.SpellHitPercent],
			stats.ExpertiseRating:     ownerStats[stats.ExpertiseRating],
			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
			stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
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
