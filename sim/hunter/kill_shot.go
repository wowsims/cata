package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerKillShotSpell() {
	icd := core.Cooldown{
		Timer:    hunter.NewTimer(),
		Duration: time.Second * 6,
	}
	hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 53351},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskRangedSpecial,
		ClassSpellMask: HunterSpellKillShot,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagRanged,
		MissileSpeed:   40,
		MinRange:       0,
		MaxRange:       45,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20()
		},
		DamageMultiplier: 4.2,
		CritMultiplier:   hunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			normalizedWeaponDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower())

			result := spell.CalcDamage(sim, target, normalizedWeaponDamage, spell.OutcomeRangedHitAndCrit)
			if icd.IsReady(sim) {
				icd.Use(sim)
				spell.CD.Reset()
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
