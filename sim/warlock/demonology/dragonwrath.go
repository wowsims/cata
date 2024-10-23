package demonology

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	// https://docs.google.com/spreadsheets/d/12jnHZgMAYDTBmkeFjApaHL5yiiDlxXHYDbTXy2QCEBA/edit?gid=707775684#gid=707775684
	cata.CreateDTRClassConfig(proto.Spec_SpecDemonologyWarlock, 0.126).
		AddSpell(50589, cata.NewDragonwrathSpellConfig().SupressSpell()) // Immolation Aura TODO: Verify Spell Interaction

}
