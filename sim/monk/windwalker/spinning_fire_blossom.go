package windwalker

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

/*
Tooltip:
Deals ${2.1*$<low>} to ${2.1*$<high>} Fire damage to the first enemy target in front of you within 50 yards.

If Spinning Fire Blossom travels further than 10 yards, the damage is increased by 50% and you root the target for 2 sec.
*/
func (ww *WindwalkerMonk) registerSpinningFireBlossom() {
	actionID := core.ActionID{SpellID: 115073}
	chiMetrics := ww.NewChiMetrics(actionID)

	ww.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          monk.SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellSpinningFireBlossom,
		MissileSpeed:   20,
		MaxRange:       50,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 2.1,
		ThreatMultiplier: 1,
		CritMultiplier:   ww.DefaultCritMultiplier(), // TODO: Spell or melee?

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return ww.ComboPoints() >= 1
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := ww.CalculateMonkStrikeDamage(sim, spell)

			if ww.DistanceFromTarget >= 10 {
				baseDamage *= 1.5
			}

			ww.SpendChi(sim, 1, chiMetrics)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit) // TODO: Spell or melee?
		},
	})
}
