package fire

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerCriticalMass() {

	criticalMassCritBuff := fire.GetStat(stats.SpellCritPercent) * 1.3

	criticalMassCritBuffMod := fire.AddDynamicMod(core.SpellModConfig{
		FloatValue: criticalMassCritBuff,
		ClassMask:  mage.MageSpellFireball | mage.MageSpellFrostfireBolt | mage.MageSpellScorch | mage.MageSpellPyroblast,
		Kind:       core.SpellMod_BonusCrit_Percent,
	})

	fire.AddOnTemporaryStatsChange(func(sim *core.Simulation, buffAura *core.Aura, statsChangeWithoutDeps stats.Stats) {
		criticalMassCritBuffMod.UpdateFloatValue(criticalMassCritBuff)
	})

	criticalMassCritBuffMod.Activate()

}
