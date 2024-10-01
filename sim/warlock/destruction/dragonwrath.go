package destruction

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecDestructionWarlock, 0.11).
		AddSpell(348, cata.NewDragonwrathSpellConfig().SupressImpact()).                          // Immolate
		AddSpell(47897, cata.NewDragonwrathSpellConfig().SupressImpact()).                        // Shadowflame
		AddSpell(17962, cata.NewDragonwrathSpellConfig().WithSpellHandler(customImmolateHandler)) // Conflagrate
}

// TODO: Verify this is how it's supposed to work
func customImmolateHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	copySpell := spell.Unit.GetSpell(spell.WithTag(71086))

	// some closure magic
	warlock := spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit).(warlock.WarlockAgent).GetWarlock()
	baseDamage := warlock.CalcScalingSpellDmg(0.43900001049)
	mulitplier := spell.DamageMultiplier
	if copySpell == nil {
		copyConfig := cata.GetDRTSpellConfig(spell)
		copyConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.DamageMultiplier *= mulitplier
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier /= mulitplier
		}
		copySpell = spell.Unit.RegisterSpell(copyConfig)
	}

	cata.CastDTRSpell(sim, copySpell, result.Target)
}
