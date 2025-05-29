package hunter

import (
	"github.com/wowsims/mop/sim/core"
)

func (hunt *Hunter) ApplyHotfixes() {
	hunt.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterSpellExplosiveShot,
		FloatValue: 0.1,
	})
	hunt.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterSpellChimeraShot,
		FloatValue: 0.5,
	})
}
