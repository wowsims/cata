package warlock

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// TODO: Pet Mana
func (warlock *Warlock) registerLifeTap() {
	actionID := core.ActionID{SpellID: 1454}
	impLifetap := 1.0 + 0.1*float64(warlock.Talents.ImprovedLifeTap)
	manaMetrics := warlock.NewManaMetrics(actionID)
	//petManaGain := 0.3 * float64(warlock.Talents.ManaFeed)

	//var petManaMetrics *core.ResourceMetrics
	//if warlock.Talents.ManaFeed > 0 && warlock.ActivePet != nil {
	//	petManaMetrics = warlock.ActivePet.NewManaMetrics(actionID)
	//}

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellLifeTap,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			restore := 0.15 * warlock.GetStat(stats.Health) * 1.2 * impLifetap
			warlock.AddMana(sim, restore, manaMetrics)

			//if warlock.Talents.ManaFeed > 0 && warlock.ActivePet != nil {
			//warlock.ActivePet.AddMana(sim, restore*petManaGain, petManaMetrics)
			//}
		},
	})
}
