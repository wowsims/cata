package destruction

import (
	"github.com/wowsims/mop/sim/common/mop"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
)

func init() {
	// https://docs.google.com/spreadsheets/d/12jnHZgMAYDTBmkeFjApaHL5yiiDlxXHYDbTXy2QCEBA/edit?gid=1393367300#gid=1393367300
	mop.CreateDTRClassConfig(proto.Spec_SpecDestructionWarlock, 0.116).
		AddSpell(17962, mop.NewDragonwrathSpellConfig().WithSpellHandler(customImmolateHandler)). // Conflagrate
		AddSpell(47897, mop.NewDragonwrathSpellConfig().IsAoESpell())                             // Shadowflame
}

// TODO: Verify this is how it's supposed to work
func customImmolateHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	copySpell := spell.Unit.GetSpell(spell.WithTag(71086))
	warlock := spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit).(warlock.WarlockAgent).GetWarlock()
	baseDamage := warlock.CalcScalingSpellDmg(0.43900001049)

	if copySpell == nil {
		copyConfig := mop.GetDRTSpellConfig(spell)
		copySpell = spell.Unit.RegisterSpell(copyConfig)
	}

	copySpell.ApplyEffects = immolationFactory(baseDamage, float64(warlock.Immolate.Dot(result.Target).HastedTickCount())*0.6)
	mop.CastDTRSpell(sim, copySpell, result.Target)
}

func immolationFactory(damage float64, multiplier float64) func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	return func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		spell.DamageMultiplier *= multiplier
		spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		spell.DamageMultiplier /= multiplier
	}
}
