package fire

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) GetCriticalMassCritPercentage() float64 {
	return fire.GetStat(stats.SpellCritPercent) * 1.3
}

func (fire *FireMage) registerCriticalMass() {

	criticalMassCritBuffMod := fire.AddDynamicMod(core.SpellModConfig{
		ClassMask: mage.MageSpellFireball | mage.MageSpellFrostfireBolt | mage.MageSpellScorch | mage.MageSpellPyroblast,
		Kind:      core.SpellMod_BonusCrit_Percent,
	})

	fire.AddOnTemporaryStatsChange(func(sim *core.Simulation, buffAura *core.Aura, statsChangeWithoutDeps stats.Stats) {
		critChance := fire.GetCriticalMassCritPercentage()
		criticalMassCritBuffMod.UpdateFloatValue(critChance)
	})

	core.MakePermanent(fire.RegisterAura(core.Aura{
		Label: "Mastery: Icicles - Water Elemental",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			criticalMassCritBuffMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			criticalMassCritBuffMod.Deactivate()
		},
	}))

}
