package elemental

import (
	"github.com/wowsims/mop/sim/common/mop"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/shaman"
)

func init() {
	//TODO (Chain Lightning) I think multiple overloads dup are possible, but couldn't
	//observe any during testing because of low chance of it happening ?
	//DTR proc also do not always start on the same target the base spell did
	//so DPS on a specific target might not be accurate.
	mop.CreateDTRClassConfig(proto.Spec_SpecElementalShaman, 0.108).
		AddSpell(88767, mop.NewDragonwrathSpellConfig().WithSpellHandler(customFulminationHandler)).       // Fullmination
		AddSpell(403, mop.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler)).               // Lightning Bolt
		AddSpell(421, mop.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler).ProcPerCast()). // Chain Lightning
		AddSpell(51505, mop.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler)).             // Lava Burst
		AddSpell(3599, mop.NewDragonwrathSpellConfig().SupressSpell()).                                    // Searing Totem
		AddSpell(8190, mop.NewDragonwrathSpellConfig().SupressSpell()).                                    // Magma Totem
		AddSpell(51490, mop.NewDragonwrathSpellConfig().ProcPerCast()).                                    // Thunderstorm
		AddSpell(77478, mop.NewDragonwrathSpellConfig().IsAoESpell()).                                     // Earthquake
		AddSpell(1535, mop.NewDragonwrathSpellConfig().IsAoESpell())                                       // Fire Nova
}

func overloadCopyHandler(unit *core.Unit, spell *core.Spell) {
	copySpell := mop.GetDRTSpellConfig(spell)
	if spell.Tag == shaman.CastTagLightningOverload { // overload tag
		copySpell.DamageMultiplier = 0.75
	}

	copySpell.BonusCoefficient = spell.BonusCoefficient
	copySpell.Flags |= core.SpellFlagNoOnCastComplete
	unit.RegisterSpell(copySpell)
}

// need to calculate fulmination damage non-delayed to have correct stack amount
func customFulminationHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	copySpell := spell.Unit.GetSpell(spell.WithTag(71086))

	// some closure magic
	shaman := spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit).(shaman.ShamanAgent).GetShaman()
	totalDamage := (shaman.ClassSpellScaling*0.38899999857 + 0.267*spell.SpellPower()) * (float64(shaman.LightningShieldAura.GetStacks()) - 3)
	if copySpell == nil {
		copyConfig := mop.GetDRTSpellConfig(spell)
		copyConfig.Cast.ModifyCast = nil
		copySpell = spell.Unit.RegisterSpell(copyConfig)
	}

	copySpell.ApplyEffects = damageFactory(totalDamage)
	mop.CastDTRSpell(sim, copySpell, result.Target)
}

func damageFactory(damage float64) func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	return func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
	}
}
