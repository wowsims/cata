package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerLightningBoltSpell() {
	shaman.LightningBolt = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(false))
	shaman.LightningBoltOverload = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(true))
}

func (shaman *Shaman) newLightningBoltSpellConfig(isElementalOverload bool) core.SpellConfig {
	spellConfig := shaman.newElectricSpellConfig(core.ActionID{SpellID: 403}, 0.06, time.Millisecond*2500, isElementalOverload, 0.714)

	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, SpellMaskLightningBoltOverload, SpellMaskLightningBolt)

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := shaman.calcDamageStormstrikeCritChance(sim, target, 770, spell)

		if result.Landed() && sim.RandomFloat("Lightning Bolt Elemental Overload") < shaman.GetOverloadChance() {
			shaman.LightningBoltOverload.Cast(sim, target)
		}

		spell.DealDamage(sim, result)
	}

	return spellConfig
}
