package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerHammerOfWrathSpell() {
	howMinDamage, howMaxDamage :=
		core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 3.90000009537, 0.10000000149)

	paladin.HammerOfWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 24275},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHammerOfWrath,

		MissileSpeed: 20,
		MaxRange:     30,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 12,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 6 * time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20() || (paladin.Talents.SanctifiedWrath > 0 && paladin.AvengingWrathAura.IsActive())
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.RollWithLabel(howMinDamage, howMaxDamage, "Hammer of Wrath"+paladin.Label) +
				0.11699999869*spell.SpellPower() +
				0.38999998569*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
