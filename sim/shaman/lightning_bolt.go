package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) registerLightningBoltSpell() {
	shaman.LightningBolt = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(false))
	shaman.LightningBoltOverload[0] = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(true))
	shaman.LightningBoltOverload[1] = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(true))
}

func (shaman *Shaman) newLightningBoltSpellConfig(isElementalOverload bool) core.SpellConfig {
	shamConfig := ShamSpellConfig{
		ActionID:            core.ActionID{SpellID: 403},
		IsElementalOverload: isElementalOverload,
		BaseCostPercent:     7.1,
		BonusCoefficient:    0.73900002241,
		BaseCastTime:        time.Millisecond * 2500,
	}
	spellConfig := shaman.newElectricSpellConfig(shamConfig)

	spellConfig.Flags |= core.SpellFlagCanCastWhileMoving

	spellConfig.ClassSpellMask = core.TernaryInt64(isElementalOverload, SpellMaskLightningBoltOverload, SpellMaskLightningBolt)
	spellConfig.MissileSpeed = core.TernaryFloat64(isElementalOverload, 30, 35)

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 1.13999998569, 0.13300000131)
		result := shaman.calcDamageStormstrikeCritChance(sim, target, baseDamage, spell)

		idx := core.TernaryInt32(spell.Flags.Matches(SpellFlagIsEcho), 1, 0)
		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance(), "Lightning Bolt Elemental Overload") {
				shaman.LightningBoltOverload[idx].Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		})
	}

	return spellConfig
}
