package paladin

import (
	"github.com/wowsims/cata/sim/core"
)

// Divine Storm is a non-ap normalised instant attack that has a weapon damage % modifier with a 1.0 coefficient.
// It does this damage to all targets in range.
// DS also heals up to 3 party or raid members for 25% of the total damage caused.
// The heal has threat implications, but given prot paladin cannot get enough talent
// points to take DS, we'll ignore it for now.

func (paladin *Paladin) registerDivineStorm() {
	if !paladin.Talents.DivineStorm {
		return
	}
	results := make([]*core.SpellResult, paladin.Env.GetNumTargets())
	actionId := core.ActionID{SpellID: 53385}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskDivineStorm | SpellMaskSpecialAttack,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.05,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD:          *paladin.sharedBuilderCooldown,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := 0
			for idx := range results {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				if results[idx].Landed() {
					numHits += 1
				}
				target = sim.Environment.NextTargetUnit(target)
			}
			for _, result := range results {
				spell.DealDamage(sim, result)
			}
			if numHits >= 4 {
				paladin.GainHolyPower(sim, 1, hpMetrics)
			}
		},
	})
}
