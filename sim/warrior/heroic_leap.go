package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterHeroicLeap() {

	numHits := warrior.Env.GetNumTargets()
	warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6544},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHeroicLeap | SpellMaskSpecialAttack,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 1,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if warrior.AutoAttacks.MH().SwingSpeed == warrior.AutoAttacks.OH().SwingSpeed {
					warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, true)
				} else {
					warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, false)
				}
			},
			IgnoreHaste: true,
		},
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1 + 0.5*spell.MeleeAttackPower()
			curTarget := target

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				damageTarget := spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				spell.DealDamage(sim, damageTarget)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
