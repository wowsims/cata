package fire

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerCriticalMass() {

	getCritPercent := func() float64 {
		return fire.GetStat(stats.SpellCritPercent) * .5 // https://us.forums.blizzard.com/en/wow/t/feedback-mists-of-pandaria-class-changes/2117387/327
	}

	criticalMassCritBuffMod := fire.AddDynamicMod(core.SpellModConfig{
		FloatValue: getCritPercent(),
		ClassMask:  mage.MageSpellFireball | mage.MageSpellFrostfireBolt | mage.MageSpellScorch | mage.MageSpellPyroblast,
		Kind:       core.SpellMod_BonusCrit_Percent,
	})

	fire.AddOnTemporaryStatsChange(func(sim *core.Simulation, buffAura *core.Aura, statsChangeWithoutDeps stats.Stats) {
		critChance := getCritPercent()
		criticalMassCritBuffMod.UpdateFloatValue(critChance)
	})

	core.MakePermanent(fire.RegisterAura(core.Aura{
		Label: "Critical Mass",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			criticalMassCritBuffMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			criticalMassCritBuffMod.Deactivate()
		},
	}))

}
