package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (dk *DeathKnight) registerSoulReaper() {
	actionID := core.ActionID{SpellID: 114867}

	dotTickSpell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellSoulReaper,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Soul Reaper",
			},
			TickLength:    time.Second * 5,
			NumberOfTicks: 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if target.CurrentHealthPercent() >= 0.35 {
					return
				}

				baseDamage := dk.CalcAndRollDamageRange(sim, 48, 0.15000000596) +
					1.20000004768*dot.Spell.MeleeAttackPower()
				dot.Snapshot(target, baseDamage)
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickMagicCrit)
			},
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	runeCost := core.RuneCostOptions{
		RunicPowerGain: 10,
		Refundable:     true,
	}

	var tag int32
	switch dk.Spec {
	case proto.Spec_SpecBloodDeathKnight:
		tag = 1 // Actually 114866
		runeCost.BloodRuneCost = 1
	case proto.Spec_SpecFrostDeathKnight:
		tag = 2 // Actually 130735
		runeCost.FrostRuneCost = 1
	default:
		tag = 3 // Actually 130736
		runeCost.UnholyRuneCost = 1
	}
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(tag),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellSoulReaper,

		MaxRange: core.MaxMeleeRange,

		RuneCost: runeCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1.3,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.MHWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				spell.RelatedDotSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},

		RelatedDotSpell: dotTickSpell,
	})
}
