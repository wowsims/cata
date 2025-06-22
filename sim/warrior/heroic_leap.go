package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) registerHeroicLeap() {
	results := make([]*core.SpellResult, warrior.Env.GetNumTargets())

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6544},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAoE | core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: SpellMaskHeroicLeap,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 45,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, warrior.AutoAttacks.MH().SwingSpeed == warrior.AutoAttacks.OH().SwingSpeed)
			},
			IgnoreHaste: true,
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   warrior.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1 + 0.5*spell.MeleeAttackPower()

			for i, enemyTarget := range sim.Encounter.TargetUnits {
				results[i] = spell.CalcDamage(sim, enemyTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
