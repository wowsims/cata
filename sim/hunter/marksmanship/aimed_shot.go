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
		MinRange:       0,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Second * 3,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()

				mmHunter.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime+spell.CastTime())
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / mmHunter.RangedSwingSpeed())
			},
		},
		DamageMultiplier: 4.5,
		CritMultiplier:   mmHunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := spell.Unit.RangedNormalizedWeaponDamage(sim, spell.RangedAttackPower())
			baseDamage := wepDmg
			baseDamage += 2604.9 + sim.RandomFloat("Aimed Shot")*2742

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
