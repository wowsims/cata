package paladin

import (
	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerCrusaderStrike() {
	actionId := core.ActionID{SpellID: 35395}
	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCrusaderStrike,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 10,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.BuilderCooldown(),
				Duration: paladin.sharedBuilderBaseCD,
			},
		},

		DamageMultiplier: 1.35,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				holyPowerGain := core.TernaryInt32(paladin.ZealotryAura.IsActive(), 3, 1)
				paladin.HolyPower.Gain(holyPowerGain, actionId, sim)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
