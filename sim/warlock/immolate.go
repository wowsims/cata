package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerImmolateSpell() {
	fireAndBrimstoneBonus := 0.05 * float64(warlock.Talents.FireAndBrimstone)

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 348},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellImmolate,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.08,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.21999999881,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if warlock.Talents.ChaosBolt {
						warlock.ChaosBolt.DamageMultiplierAdditive += fireAndBrimstoneBonus
					}
					warlock.Incinerate.DamageMultiplierAdditive += fireAndBrimstoneBonus
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if warlock.Talents.ChaosBolt {
						warlock.ChaosBolt.DamageMultiplierAdditive -= fireAndBrimstoneBonus
					}
					warlock.Incinerate.DamageMultiplierAdditive -= fireAndBrimstoneBonus
				},
			},
			BonusCoefficient: 0.17599999905,
			NumberOfTicks:    5,
			TickLength:       time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 444
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 699, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}
