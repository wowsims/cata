package fire

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/mage"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecFireMage, 1.0/12.0).
		AddSpell(11129, cata.NewDragonwrathSpellConfig().SupressImpact()). // Combustion
		AddSpell(2120, cata.NewDragonwrathSpellConfig().SupressImpact()).  // Flamestrike
		AddSpell(88148, cata.NewDragonwrathSpellConfig().
			SupressImpact(). // Flamestrike (Blast Wave)
			WithCustomSpell(func(unit *core.Unit, spell *core.Spell) {
				config := unit.Env.GetAgentFromUnit(unit).(mage.MageAgent).GetMage().GetFlameStrikeConfig(88148, true)
				config.ActionID.Tag = 71086
				config.Dot.Aura.Label += " DTR"
				unit.RegisterSpell(config)
			})).
		AddSpell(11113, cata.NewDragonwrathSpellConfig().ProcPerCast())
}
