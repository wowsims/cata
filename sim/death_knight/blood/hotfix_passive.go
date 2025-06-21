package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (bdk *BloodDeathKnight) registerHotfixPassive() {
	core.MakePermanent(bdk.RegisterAura(core.Aura{
		Label: "Hotfix Passive" + bdk.Label,
	})).AttachSpellMod(core.SpellModConfig{
		// Beta changes 2025-06-16: https://www.wowhead.com/mop-classic/news/blood-death-knights-buffed-and-even-more-class-balance-adjustments-mists-of-377292
		// - Outbreak’s base cooldown for Blood Death Knights has been decreased to 30 seconds (was 60 seconds). [New]
		// EffectIndex 0 on the Blood specific Hotfix Passive https://wago.tools/db2/SpellEffect?build=5.5.0.61496&filter%5BSpellID%5D=137008&page=1&sort%5BEffectIndex%5D=asc
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: death_knight.DeathKnightSpellOutbreak,
		TimeValue: time.Second * -30,
	})
}
