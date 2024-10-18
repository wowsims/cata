package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"time"
)

func (paladin *Paladin) registerExorcism() {
	exorcismMinDamage, exorcismMaxDamage :=
		core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 2.663, 0.11)

	paladin.Exorcism = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 879},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskExorcism,

		MaxRange: 30,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if paladin.CurrentMana() >= cast.Cost {
					castTime := paladin.ApplyCastSpeedForSpell(cast.CastTime, spell)
					if castTime > 0 {
						paladin.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
					}
				}
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.RollWithLabel(exorcismMinDamage, exorcismMaxDamage, "Exorcism") +
				0.344*max(spell.SpellPower(), spell.MeleeAttackPower())

			bonusCritPercent := core.TernaryFloat64(
				target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead,
				100,
				0)

			spell.BonusCritPercent += bonusCritPercent
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.BonusCritPercent -= bonusCritPercent

			spell.DealOutcome(sim, result)
		},
	})
}
