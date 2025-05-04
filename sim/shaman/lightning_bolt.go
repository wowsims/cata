package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) registerLightningBoltSpell() {
	shaman.LightningBolt = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(false))
	shaman.LightningBoltOverload = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(true))
}

func (shaman *Shaman) newLightningBoltSpellConfig(isElementalOverload bool) core.SpellConfig {
	spellConfig := shaman.newElectricSpellConfig(core.ActionID{SpellID: 403}, 7.1, time.Millisecond*2500, isElementalOverload, 0.73900002241)

	spellConfig.Flags |= core.SpellFlagCanCastWhileMoving

	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, SpellMaskLightningBoltOverload, SpellMaskLightningBolt)
	spellConfig.MissileSpeed = core.TernaryFloat64(isElementalOverload, 20, 35)

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 1.13999998569, 0.13300000131)
		result := shaman.calcDamageStormstrikeCritChance(sim, target, baseDamage, spell)

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellProc) { //So that procs from DTR does not cast an overload
				if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance(), "Lightning Bolt Elemental Overload") {
					shaman.LightningBoltOverload.Cast(sim, target)
				}
			}

			spell.DealDamage(sim, result)
		})
	}

	return spellConfig
}
