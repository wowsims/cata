package mop

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	// Keep these in order by item ID
	// Agile Primal Diamond
	core.NewItemEffect(76884, core.ApplyMetaGemCriticalDamageEffect)
	// Burning Primal Diamond
	core.NewItemEffect(76885, core.ApplyMetaGemCriticalDamageEffect)
	// Reverberating Primal Diamond
	core.NewItemEffect(76886, core.ApplyMetaGemCriticalDamageEffect)
	// Revitalizing Primal Diamond
	core.NewItemEffect(76888, core.ApplyMetaGemCriticalDamageEffect)

	// Austere Primal Diamond
	core.NewItemEffect(76895, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		character.ApplyEquipScaling(stats.Armor, 1.02)
	})
}
