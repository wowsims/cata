package protection

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (prot *ProtectionPaladin) registerSanctuary() {
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Sanctuary" + prot.Label,
		ActionID:   core.ActionID{SpellID: 105805},
		BuildPhase: core.CharacterBuildPhaseBuffs,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.ApplyDynamicEquipScaling(sim, stats.Armor, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.RemoveDynamicEquipScaling(sim, stats.Armor, 1.1)
		},
	})).AttachAdditivePseudoStatBuff(
		&prot.PseudoStats.BaseDodgeChance,
		0.02,
	).AttachMultiplicativePseudoStatBuff(
		&prot.PseudoStats.DamageTakenMultiplier,
		0.85,
	)
}
