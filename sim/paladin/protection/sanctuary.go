package protection

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// Decreases damage taken by 15%, increases armor value from items by 10%, and increases your chance to dodge by 2%.
func (prot *ProtectionPaladin) registerSanctuary() {
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Sanctuary" + prot.Label,
		ActionID:   core.ActionID{SpellID: 105805},
		BuildPhase: core.CharacterBuildPhaseTalents,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.ApplyDynamicEquipScaling(sim, stats.Armor, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.RemoveDynamicEquipScaling(sim, stats.Armor, 1.1)
		},
	})).AttachAdditivePseudoStatBuff(
		&prot.PseudoStats.BaseDodgeChance, 0.02,
	).AttachMultiplicativePseudoStatBuff(
		// Beta changes 2025-06-13: https://www.wowhead.com/mop-classic/news/some-warlords-of-draenor-pre-patch-class-changes-coming-to-mists-of-pandaria-377239
		// - The damage reduction from Sanctuary has been raised to 20% (was 15%). New
		// EffectIndex 1 on the Protection specific Hotfix Passive https://wago.tools/db2/SpellEffect?build=5.5.0.61411&filter%5BSpellID%5D=137028&page=1
		&prot.PseudoStats.DamageTakenMultiplier, 0.8,
	)
}
