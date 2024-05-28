package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"time"
)

func (paladin *Paladin) RegisterExorcism() {
	paladin.Exorcism = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 48801},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskExorcism,

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

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(2483, 2771) +
				0.34*max(spell.SpellPower(), spell.MeleeAttackPower())

			bonusCrit := core.TernaryFloat64(
				target.MobType == proto.MobType_MobTypeDemon || target.MobType == proto.MobType_MobTypeUndead,
				100*core.CritRatingPerCritChance,
				0)

			spell.BonusCritRating += bonusCrit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.BonusCritRating -= bonusCrit
		},
	})
}
