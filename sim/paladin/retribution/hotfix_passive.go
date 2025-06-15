package retribution

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

func (ret *RetributionPaladin) registerHotfixPassive() {
	core.MakePermanent(ret.RegisterAura(core.Aura{
		Label: "Hotfix Passive" + ret.Label,
	})).AttachSpellMod(core.SpellModConfig{
		// Beta changes 2025-06-13: https://www.wowhead.com/mop-classic/news/some-warlords-of-draenor-pre-patch-class-changes-coming-to-mists-of-pandaria-377239
		// - Divine Storm, Crusader Strike, Judgment, Hammer of the Righteous, Hammer of Wrath, and Exorcism have all had their damage raised by 10%. New
		// - Templar’s Verdict damage raised by 20%. 6.0.2 Change
		// EffectIndex 0 and 1 on the Retribution specific Hotfix Passive https://wago.tools/db2/SpellEffect?build=5.5.0.61411&filter%5BSpellID%5D=137027&page=1
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  paladin.SpellMaskBuilderRet | paladin.SpellMaskDivineStorm,
		FloatValue: 0.1,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  paladin.SpellMaskTemplarsVerdict,
		FloatValue: 0.2,
	})

}
