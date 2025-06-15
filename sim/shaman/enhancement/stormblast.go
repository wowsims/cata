package enhancement

import (
	"github.com/wowsims/mop/sim/core"
)

func (enh *EnhancementShaman) newStormblastHitSpell(isMh bool) *core.Spell {
	config := enh.newStormstrikeHitSpellConfig(115356, isMh)
	config.SpellSchool = core.SpellSchoolNature
	return enh.RegisterSpell(config)
}

func (enh *EnhancementShaman) registerStormblastSpell() {
	mhHit := enh.newStormblastHitSpell(true)
	ohHit := enh.newStormblastHitSpell(false)

	config := enh.newStormstrikeSpellConfig(115356, &enh.StormStrikeDebuffAuras, mhHit, ohHit)
	config.SpellSchool = core.SpellSchoolNature
	config.ManaCost.BaseCostPercent = 9.372
	config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return enh.AscendanceAura.IsActive()
	}

	enh.Stormblast = enh.RegisterSpell(config)
}
