package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// Increases your total Stamina by 9% and your chance to dodge by 2%.
func (bdk *BloodDeathKnight) registerVeteranOfTheThirdWar() {
	core.MakePermanent(bdk.GetOrRegisterAura(core.Aura{
		Label:    "Veteran of the Third War" + bdk.Label,
		ActionID: core.ActionID{SpellID: 50029},
	})).AttachMultiplicativePseudoStatBuff(
		&bdk.PseudoStats.BaseDodgeChance, 0.02,
	).AttachStatDependency(
		bdk.NewDynamicMultiplyStat(stats.Stamina, 1.09),
	)
}
