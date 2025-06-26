package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (war *Warrior) registerThunderClap() {
	war.ThunderClapAuras = war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.WeakenedBlowsAura(target)
	})

	results := make([]*core.SpellResult, war.Env.GetNumTargets())
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6343},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskThunderClap,

		RageCost: core.RageCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := war.CalcScalingSpellDmg(0.25) + spell.MeleeAttackPower()*0.44999998808

			for i, aoeTarget := range sim.Encounter.TargetUnits {
				results[i] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			}

			war.CastNormalizedSweepingStrikesAttack(results, sim, target)

			for _, result := range results {
				if result.Landed() {
					war.ThunderClapAuras.Get(result.Target).Activate(sim)
				}
				spell.DealDamage(sim, result)
			}
		},

		RelatedAuraArrays: war.ThunderClapAuras.ToMap(),
	})
}
