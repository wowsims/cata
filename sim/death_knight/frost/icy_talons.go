package frost

import "github.com/wowsims/mop/sim/core"

// Your melee attack speed is increased by 45%.
func (fdk *FrostDeathKnight) registerIcyTalons() {
	core.MakePermanent(fdk.RegisterAura(core.Aura{
		Label:      "Icy Talons" + fdk.Label,
		ActionID:   core.ActionID{SpellID: 50887},
		BuildPhase: core.CharacterBuildPhaseTalents,
	})).AttachMultiplicativePseudoStatBuff(
		&fdk.PseudoStats.MeleeSpeedMultiplier, 1.45,
	)
}
