package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (uhdk *UnholyDeathKnight) registerUnholyMight() {
	core.MakePermanent(uhdk.RegisterAura(core.Aura{
		Label:      "Unholy Might" + uhdk.Label,
		ActionID:   core.ActionID{SpellID: 91107},
		BuildPhase: core.CharacterBuildPhaseTalents,
	})).AttachStatDependency(uhdk.NewDynamicMultiplyStat(stats.Strength, 1.35))
}
