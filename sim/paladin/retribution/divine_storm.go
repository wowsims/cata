package retribution

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/paladin"
	"time"
)

// Divine Storm is a non-ap normalised instant attack that has a weapon damage % modifier with a 1.0 coefficient.
// It does this damage to all targets in range.
// DS also heals up to 3 party or raid members for 25% of the total damage caused.
// The heal has threat implications, but given prot paladin cannot get enough talent
// points to take DS, we'll ignore it for now.

func (retPaladin *RetributionPaladin) RegisterDivineStorm() {
	if !retPaladin.Talents.DivineStorm {
		return
	}

	numTargets := retPaladin.Env.GetNumTargets()
	actionId := core.ActionID{SpellID: 53385}
	hpMetrics := retPaladin.NewHolyPowerMetrics(actionId)

	retPaladin.DivineStorm = retPaladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskDivineStorm | paladin.SpellMaskSpecialAttack,

		MaxRange: 8,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.05,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    retPaladin.NewTimer(),
				Duration: 4500 * time.Millisecond,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   retPaladin.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := 0
			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				result := spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				if result.Landed() {
					numHits += 1
				}
				results[idx] = result
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}

			if numHits >= 4 {
				retPaladin.GainHolyPower(sim, 1, hpMetrics)
			}
		},
	})
}
