package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// TODO: No patch notes for this ability, need to validate the damage and threat coefficients haven't changed
func (war *Warrior) registerHeroicThrow() {
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 57755},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHeroicThrow,
		MaxRange:       30,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 30,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				war.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime)
			},
			IgnoreHaste: true,
		},
		DamageMultiplier: 0.5,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := war.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
}
