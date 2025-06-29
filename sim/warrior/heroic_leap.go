package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterHeroicLeap() {

	maxHits := warrior.Env.TotalTargetCount()
	results := make([]*core.SpellResult, maxHits)

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
			numHits := sim.Environment.ActiveTargetCount()

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				curTarget = sim.Environment.NextActiveTargetUnit(curTarget)
			}

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
			}

		},
	})
}
