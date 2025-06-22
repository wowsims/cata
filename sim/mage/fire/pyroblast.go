package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerPyroblastSpell() {
	actionID := core.ActionID{SpellID: 11366}
	pyroblastVariance := 0.24    // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Variance"
	pyroblastScaling := 1.98     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Coefficient"
	pyroblastCoefficient := 1.98 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "BonusCoefficient"
	pyroblastDotScaling := .36
	pyroblastDotCoefficient := .36

	fire.Pyroblast = fire.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellPyroblast,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 3500,
			},
		},

		DamageMultiplier: 1 * 1.12, //https://us.forums.blizzard.com/en/wow/t/feedback-mists-of-pandaria-class-changes/2117387/327
		CritMultiplier:   fire.DefaultCritMultiplier(),
		BonusCoefficient: pyroblastCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !fire.InstantPyroblastAura.IsActive() && fire.PresenceOfMindAura != nil {
				fire.PresenceOfMindAura.Deactivate(sim)
			}
			fire.InstantPyroblastAura.Deactivate(sim)
			baseDamage := fire.CalcAndRollDamageRange(sim, pyroblastScaling, pyroblastVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			fire.HeatingUpSpellHandler(sim, spell, result, func() {
				spell.RelatedDotSpell.Cast(sim, target)
				spell.DealDamage(sim, result)
			})
		},
	})

	fire.Pyroblast.RelatedDotSpell = fire.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: mage.MageSpellPyroblastDot,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1 * 1.12, //https://us.forums.blizzard.com/en/wow/t/feedback-mists-of-pandaria-class-changes/2117387/327
		CritMultiplier:   fire.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "PyroblastDoT",
			},
			NumberOfTicks:       6,
			TickLength:          time.Second * 3,
			BonusCoefficient:    pyroblastDotCoefficient,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, fire.CalcScalingSpellDmg(pyroblastDotScaling))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
