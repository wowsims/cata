package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) registerCurseOfElements() {
	warlock.CurseOfElementsAuras = warlock.NewEnemyAuraArray(core.CurseOfElementsAura)

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1490},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfElements,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 10},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		ThreatMultiplier: 1,
		FlatThreatBonus:  104,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfElementsAuras.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: warlock.CurseOfElementsAuras.ToMap(),
	})
}

func (warlock *Warlock) registerCurseOfWeakness() {
	warlock.CurseOfWeaknessAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.CurseOfWeaknessAura(target)
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 702},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfWeakness,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 10},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		ThreatMultiplier: 1,
		FlatThreatBonus:  32,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfWeaknessAuras.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: warlock.CurseOfWeaknessAuras.ToMap(),
	})
}

func (warlock *Warlock) registerCurseOfTongues() {
	actionID := core.ActionID{SpellID: 1714}

	// Empty aura so we can simulate cost/time to keep tongues up
	warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Curse of Tongues",
			ActionID: actionID,
			Duration: 30 * time.Second,
		})
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfTongues,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 4},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		ThreatMultiplier: 1,
		FlatThreatBonus:  52,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfTonguesAuras.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: warlock.CurseOfTonguesAuras.ToMap(),
	})
}

func (warlock *Warlock) registerBaneOfAgony() {
	baseTickDmg := warlock.CalcScalingSpellDmg(0.13300000131)

	warlock.BaneOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 980},
		SpellSchool:    core.SpellSchoolShadow,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellBaneOfAgony,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 10},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Bane of Agony",
			},
			NumberOfTicks:       12,
			TickLength:          2 * time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.08799999952,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 0.5*baseTickDmg)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				if dot.TickCount()%4 == 0 { // CoA ramp up
					dot.SnapshotBaseDamage += 0.5 * baseTickDmg
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				warlock.BaneOfDoom.Dot(target).Deactivate(sim)
				//TODO: Cancel BaneOfHavoc
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

func (warlock *Warlock) registerBaneOfDoom() {
	ebonImpBonusSummon := 0.1 * float64(warlock.Talents.ImpendingDoom)

	warlock.BaneOfDoom = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 603},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellBaneOfDoom,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 15},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		FlatThreatBonus:          40,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Bane of Doom",
			},
			NumberOfTicks:    4,
			TickLength:       15 * time.Second,
			BonusCoefficient: 0.87999999523,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, warlock.CalcScalingSpellDmg(2.02399992943))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				if sim.Proc(0.2+ebonImpBonusSummon, "Ebon Imp") {
					warlock.EbonImp.EnableWithTimeout(sim, warlock.EbonImp, 15*time.Second)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.BaneOfAgony.Dot(target).Deactivate(sim)
				//TODO: Cancel BaneOfHavoc
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
