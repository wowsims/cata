package demonology

import (
	"github.com/wowsims/mop/sim/common/mop"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	// https://docs.google.com/spreadsheets/d/12jnHZgMAYDTBmkeFjApaHL5yiiDlxXHYDbTXy2QCEBA/edit?gid=707775684#gid=707775684
	mop.CreateDTRClassConfig(proto.Spec_SpecDemonologyWarlock, 0.126).
		AddSpell(50589, mop.NewDragonwrathSpellConfig().TreatCastAsTick()). // Immolation Aura
		AddSpell(47897, mop.NewDragonwrathSpellConfig().IsAoESpell())       // Shadowflame
}
