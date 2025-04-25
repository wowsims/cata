package affliction

import (
	"github.com/wowsims/mop/sim/common/mop"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	// https://docs.google.com/spreadsheets/d/12jnHZgMAYDTBmkeFjApaHL5yiiDlxXHYDbTXy2QCEBA/edit?gid=552198772#gid=552198772
	mop.CreateDTRClassConfig(proto.Spec_SpecAfflictionWarlock, 0.105).
		AddSpell(47897, mop.NewDragonwrathSpellConfig().IsAoESpell()) // Shadowflame
}
