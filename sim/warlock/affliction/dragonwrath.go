package affliction

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	// https://docs.google.com/spreadsheets/d/12jnHZgMAYDTBmkeFjApaHL5yiiDlxXHYDbTXy2QCEBA/edit?gid=552198772#gid=552198772
	cata.CreateDTRClassConfig(proto.Spec_SpecAfflictionWarlock, 0.105).
		AddSpell(47897, cata.NewDragonwrathSpellConfig().IsAoESpell()) // Shadowflame
}
