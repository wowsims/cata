package marksmanship

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (mm *MarksmanshipHunter) registerSteadyShotSpell() {

	ssMetrics := mm.NewFocusMetrics(core.ActionID{SpellID: 56641})

	mm.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 56641},
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: hunter.HunterSpellSteadyShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagRanged,
		MissileSpeed:   40,
		MinRange:       0,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Millisecond * 2000,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / mm.TotalRangedHasteMultiplier())
			},
		},
		BonusCritPercent:         0,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         0.726,
		CritMultiplier:           mm.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mm.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower())
			baseDamage += mm.GetBaseDamageFromCoeff(2.112)

			intFocus := 14.0
			mm.AddFocus(sim, intFocus, ssMetrics)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
