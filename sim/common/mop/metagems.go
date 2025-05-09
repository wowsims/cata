package cata

import (
	"github.com/wowsims/mop/sim/core"
)

func init() {

	// Keep these in order by item ID
	// Agile Primal Diamond
	core.NewItemEffect(76884, core.ApplyCriticalDamageEffect)
	// Burning Primal Diamond
	core.NewItemEffect(76885, core.ApplyCriticalDamageEffect)
	// Reverberating Primal Diamond
	core.NewItemEffect(76886, core.ApplyCriticalDamageEffect)
	// Revitalizing Primal Diamond
	core.NewItemEffect(76888, core.ApplyCriticalDamageEffect)
	// Burning Primal Diamond
	core.NewItemEffect(97937, core.ApplyCriticalDamageEffect)
	// Burning Primal Diamond
	core.NewItemEffect(97534, core.ApplyCriticalDamageEffect)
	// Revitalizing Primal Diamond
	core.NewItemEffect(97306, core.ApplyCriticalDamageEffect)
}
