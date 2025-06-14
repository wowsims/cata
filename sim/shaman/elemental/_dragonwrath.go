package elemental

import (
	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/shaman"
)

func init() {
	//TODO (Chain Lightning) I think multiple overloads dup are possible, but couldn't
	//observe any during testing because of low chance of it happening ?
	//DTR proc also do not always start on the same target the base spell did
	//so DPS on a specific target might not be accurate.
	cata.CreateDTRClassConfig(proto.Spec_SpecElementalShaman, 0.108).
		AddSpell(88767, cata.NewDragonwrathSpellConfig().WithSpellHandler(customFulminationHandler)).       // Fullmination
		AddSpell(403, cata.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler)).               // Lightning Bolt
		AddSpell(421, cata.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler).ProcPerCast()). // Chain Lightning
		AddSpell(51505, cata.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler)).             // Lava Burst
		AddSpell(3599, cata.NewDragonwrathSpellConfig().SupressSpell()).                                    // Searing Totem
		AddSpell(8190, cata.NewDragonwrathSpellConfig().SupressSpell()).                                    // Magma Totem
		AddSpell(51490, cata.NewDragonwrathSpellConfig().ProcPerCast()).                                    // Thunderstorm
		AddSpell(77478, cata.NewDragonwrathSpellConfig().IsAoESpell()).                                     // Earthquake
		AddSpell(1535, cata.NewDragonwrathSpellConfig().IsAoESpell())                                       // Fire Nova
}

func overloadCopyHandler(unit *core.Unit, spell *core.Spell) {
	copySpell := cata.GetDRTSpellConfig(spell)
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
	totalDamage := (shaman.CalcScalingSpellDmg(0.38899999857) + 0.267*spell.SpellPower()) * (float64(shaman.LightningShieldAura.GetStacks()) - 3)
	if copySpell == nil {
		copyConfig := cata.GetDRTSpellConfig(spell)
		copyConfig.Cast.ModifyCast = nil
		copySpell = spell.Unit.RegisterSpell(copyConfig)
	}

	copySpell.ApplyEffects = damageFactory(totalDamage)
	cata.CastDTRSpell(sim, copySpell, result.Target)
}

func damageFactory(damage float64) func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	return func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
	}
}
