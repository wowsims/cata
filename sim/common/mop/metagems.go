package cata

import (
	"github.com/wowsims/mop/sim/core"
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
	// Burning Primal Diamond
	core.NewItemEffect(97937, core.ApplyMetaGemCriticalDamageEffect)
	// Burning Primal Diamond
	core.NewItemEffect(97534, core.ApplyMetaGemCriticalDamageEffect)
	// Revitalizing Primal Diamond
	core.NewItemEffect(97306, core.ApplyMetaGemCriticalDamageEffect)
}
