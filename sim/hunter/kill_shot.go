package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerKillShotSpell() {
	hunter.KillShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 53351},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskRangedSpecial,
		ClassSpellMask: HunterSpellKillShot,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       5,
		MaxRange:       40,
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
		DamageMultiplier: 1.5,
		CritMultiplier:   hunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			normalizedWeaponDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower())
			rapBonusDamage := spell.RangedAttackPower() * 0.45
			flatBonus := 543.0

			baseDamage := normalizedWeaponDamage + rapBonusDamage + flatBonus
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
