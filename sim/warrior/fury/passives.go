package fury

import (
	"github.com/wowsims/mop/sim/core"
)

func (war *FuryWarrior) registerCrazedBerserker() {

	war.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMeleeOH,
		FloatValue: 0.25,
	})
	war.AutoAttacks.MHConfig().DamageMultiplier *= 1.1
	war.AutoAttacks.OHConfig().DamageMultiplier *= 1.1
}
