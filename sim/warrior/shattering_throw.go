package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (war *Warrior) registerShatteringThrow() {
	shattDebuffs := war.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return core.ShatteringThrowAura(unit, war.UnitIndex)
	})

	ShatteringThrowSpell := war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 64382},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskShatteringThrow,
		MaxRange:       30,
		MissileSpeed:   50,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 5,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 12 + spell.MeleeAttackPower()*0.5
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if result.Landed() {
				shattDebuffs.Get(target).Activate(sim)
			}
		},

		RelatedAuraArrays: shattDebuffs.ToMap(),
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: ShatteringThrowSpell,
		Type:  core.CooldownTypeDPS,
	})
}
