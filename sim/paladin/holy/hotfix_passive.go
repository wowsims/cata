package holy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

func (holy *HolyPaladin) registerHotfixPassive() {
	core.MakePermanent(holy.RegisterAura(core.Aura{
		Label: "Hotfix Passive" + holy.Label,
	})).AttachSpellMod(core.SpellModConfig{
		// Beta changes 2025-06-13: https://www.wowhead.com/mop-classic/news/some-warlords-of-draenor-pre-patch-class-changes-coming-to-mists-of-pandaria-377239
		// - Eternal Flame’s periodic healing lowered by 28%. 5.4 Revert
		// EffectIndex 1 on the Holy specific Hotfix Passive https://wago.tools/db2/SpellEffect?build=5.5.0.61411&filter%5BSpellID%5D=137029&page=1
		Kind:       core.SpellMod_DotDamageDone_Pct,
		ClassMask:  paladin.SpellMaskWordOfGlory,
		FloatValue: 1.0 - 0.28,
	})
}
