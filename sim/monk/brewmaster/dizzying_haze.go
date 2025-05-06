package brewmaster

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerDizzyingHaze() {
	spellActionID := core.ActionID{SpellID: 115180}
	debuffActionID := core.ActionID{SpellID: 116330}

	bm.DizzyingHazeAuras = bm.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Dizzying Haze",
			ActionID: debuffActionID,
			Duration: 15 * time.Second,
		})
	})

	projectile := bm.RegisterSpell(core.SpellConfig{
		ActionID:       spellActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: monk.MonkSpellDizzyingHazeProjectile,
		MaxRange:       8,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.ApplyAOEThreat(spell.MeleeAttackPower() * 1.1)
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)
				if result.Landed() {
					bm.DizzyingHazeAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},
		RelatedAuraArrays: bm.DizzyingHazeAuras.ToMap(),
	})

	bm.RegisterSpell(core.SpellConfig{
		ActionID:       spellActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellDizzyingHaze,
		MaxRange:       40,
		MissileSpeed:   15,

		EnergyCost: core.EnergyCostOptions{
			Cost: 20,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    bm.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplier: 0,
		ThreatMultiplier: 1,
		CritMultiplier:   bm.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bm.StanceMatches(monk.SturdyOx)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				projectile.Cast(sim, target)
			})
		},
	})
}
