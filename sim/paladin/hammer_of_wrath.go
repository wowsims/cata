package paladin

import (
	"github.com/wowsims/cata/sim/core/proto"
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) registerHammerOfWrathSpell() {
	howMinDamage, howMaxDamage :=
		core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 3.9, 0.1)

	paladin.HammerOfWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 24275},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHammerOfWrath,

		MissileSpeed: 20,
		MaxRange:     30,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1,
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
				0.117*spell.SpellPower() +
				0.39*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			spell.WaitTravelTime(sim, func(simulation *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
