package elemental

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/shaman"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecElementalShaman, 1.0/12.0). // TODO: Proc chance unknown for now
										AddSpell(8050, cata.NewDragonwrathSpellConfig().SupressImpact()).                             // Flame Shock
										AddSpell(88767, cata.NewDragonwrathSpellConfig().WithSpellHandler(customFulminationHandler)). // Fullmination
										AddSpell(403, cata.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler)).         // Lightning Bold
										AddSpell(1, cata.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler)).           // Chain Lightning
										AddSpell(1, cata.NewDragonwrathSpellConfig().WithCustomSpell(overloadCopyHandler))            // Lava Burst
}

func overloadCopyHandler(unit *core.Unit, spell *core.Spell) {
	copySpell := cata.GetDRTSpellConfig(spell)
	if spell.Tag == shaman.CastTagLightningOverload { // overload tag
		copySpell.DamageMultiplier = 0.75
	}

	copySpell.BonusCoefficient = spell.BonusCoefficient
	unit.RegisterSpell(copySpell)
}

// need to calculate fulmination damage non-delayed to have correct stack amount
func customFulminationHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	copySpell := spell.Unit.GetSpell(spell.WithTag(71086))

	// some closure magic
	shaman := spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit).(shaman.ShamanAgent).GetShaman()
	totalDamage := (shaman.ClassSpellScaling*0.38899999857 + 0.267*spell.SpellPower()) * (float64(shaman.LightningShieldAura.GetStacks()) - 3)
	if copySpell == nil {
		copyConfig := cata.GetDRTSpellConfig(spell)
		copyConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			copySpell.CalcAndDealDamage(sim, result.Target, totalDamage, spell.OutcomeMagicHitAndCrit)
		}
		copyConfig.Cast.ModifyCast = nil
		copySpell = spell.Unit.RegisterSpell(copyConfig)
	}

	cata.CastDTRSpell(sim, copySpell, result.Target)
}
