package marksmanship

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (mmHunter *MarksmanshipHunter) registerAimedShotSpell() {

	mmHunter.AimedShot = mmHunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 19434},
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: hunter.HunterSpellAimedShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   40,
		FocusCost: core.FocusCostOptions{
			Cost: 50 - mmHunter.Talents.Efficiency*2,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Second * 3,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
				// Aimed Shot on Beta currently is a full reset
				mmHunter.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime+spell.CastTime())
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / mmHunter.RangedSwingSpeed())
			},
		},
		DamageMultiplier: 1.6,
		CritMultiplier:   mmHunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := spell.Unit.RangedNormalizedWeaponDamage(sim, spell.RangedAttackPower())
			rap := spell.RangedAttackPower() * 0.723
			baseDamage := (wepDmg + rap) + sim.Roll(776, 866)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
