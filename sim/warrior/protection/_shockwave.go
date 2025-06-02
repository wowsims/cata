package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerShockwave() {
	if !war.Talents.Shockwave {
		return
	}

	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 46968},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagAoE | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskShockwave,

		RageCost: core.RageCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 20 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.75 * spell.MeleeAttackPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}
		},
	})
}
