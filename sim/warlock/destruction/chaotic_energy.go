package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warlock"
)

func (destruction DestructionWarlock) ApplyChaoticEnergy() {
	core.MakePermanent(destruction.RegisterAura(core.Aura{
		Label: "Chaotic Energy",
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		TimeValue: -time.Millisecond * 500,
		ClassMask: warlock.WarlockSpellAll,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: 3,
		ClassMask:  warlock.WarlockSpellsChaoticEnergyDestro,
	}))

	destruction.MultiplyStat(stats.MP5, 7.25)
	destruction.HasteEffectsRegen()
}
