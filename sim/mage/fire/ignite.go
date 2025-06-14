package fire

import (
	"github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerMastery() {

	fire.ignite = cata.RegisterIgniteEffect(&fire.Unit, cata.IgniteConfig{
		ActionID:       core.ActionID{SpellID: 12846},
		ClassSpellMask: mage.MageSpellIgnite,
		DotAuraLabel:   "Ignite",
		DotAuraTag:     "IgniteDot",

		ProcTrigger: core.ProcTrigger{
			Name:     "Ignite Talent",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellDamage,
			Outcome:  core.OutcomeLanded,

			ExtraCondition: func(_ *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
				return spell.SpellSchool.Matches(core.SpellSchoolFire)
			},
		},

		DamageCalculator: func(result *core.SpellResult) float64 {
			return result.Damage * fire.GetMasteryBonus()
		},
	})

	// This is needed because we want to listen for the spell "cast" event that refreshes the Dot
	fire.ignite.Flags ^= core.SpellFlagNoOnCastComplete

}

func (fire *FireMage) GetMasteryBonus() float64 {
	return (.12 + 0.015*fire.GetMasteryPoints())
}
