package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerPowerShotSpell() {
	if !hunter.Talents.Powershot {
		return
	}
	hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 109259},
		SpellSchool: core.SpellSchoolPhysical,
		//ClassSpellMask: hunter.HunterSpellAimedShot,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagRanged,
		MissileSpeed: 40,
		MinRange:     0,
		MaxRange:     40,
		FocusCost: core.FocusCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Millisecond * 2250,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()

				hunter.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime+spell.CastTime())
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.TotalRangedHasteMultiplier())
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: 45 * time.Second,
			},
		},
		DamageMultiplier: 6,
		CritMultiplier:   hunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := spell.Unit.RangedNormalizedWeaponDamage(sim, spell.RangedAttackPower())

			result := spell.CalcDamage(sim, target, wepDmg, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
