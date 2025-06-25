package marksmanship

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (mmHunter *MarksmanshipHunter) registerAimedShotSpell() {
	mmHunter.RegisterSpell(mmHunter.getAimedShotSpell(82928, true))
	mmHunter.RegisterSpell(mmHunter.getAimedShotSpell(19434, false))
}
func (mmHunter *MarksmanshipHunter) getAimedShotSpell(spellID int32, isMasterMarksman bool) core.SpellConfig {
	normalCast := core.CastConfig{
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
			return time.Duration(float64(spell.DefaultCast.CastTime) / mmHunter.TotalRangedHasteMultiplier())
		},
	}
	freeCast := core.CastConfig{
		DefaultCast: core.Cast{
			GCD:      time.Second,
			CastTime: time.Second * 0,
		},
	}
	config := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: hunter.HunterSpellAimedShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagRanged | core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       0,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 50,
		},
		Cast:             normalCast,
		DamageMultiplier: 4.95,
		CritMultiplier:   mmHunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mmHunter.HasActiveAura("Ready, Set, Aim...")
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := spell.Unit.RangedNormalizedWeaponDamage(sim, spell.RangedAttackPower())
			baseDamage := wepDmg
			baseDamage += 2604.9 + sim.RandomFloat("Aimed Shot")*2742

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	}
	if isMasterMarksman {
		config.Cast = freeCast
		config.FocusCost.Cost = 0
		config.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
			return mmHunter.HasActiveAura("Ready, Set, Aim...")
		}
	}
	return config
}
