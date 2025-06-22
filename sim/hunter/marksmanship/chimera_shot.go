package marksmanship

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (mmHunter *MarksmanshipHunter) registerChimeraShotSpell() {
	mmHunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53209},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,

		ClassSpellMask: hunter.HunterSpellChimeraShot,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagRanged,
		MissileSpeed:   40,
		MinRange:       0,
		MaxRange:       40,

		FocusCost: core.FocusCostOptions{
			Cost: 45,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    mmHunter.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 4.57,
		CritMultiplier:   mmHunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := mmHunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower())
			baseDamage := mmHunter.GetBaseDamageFromCoeff(1.25)

			result := spell.CalcDamage(sim, target, wepDmg+baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					if mmHunter.SerpentSting.Dot(target).IsActive() {
						mmHunter.SerpentSting.Dot(target).Apply(sim)
					}
				}
				spell.DealDamage(sim, result)
			})
		},
	})
}
