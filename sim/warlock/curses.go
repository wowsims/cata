package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerCurseOfElementsSpell() {
	warlock.CurseOfElementsAuras = warlock.NewEnemyAuraArray(core.CurseOfElementsAura)

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1490},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfElements,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  156,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfElementsAuras.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{warlock.CurseOfElementsAuras},
	})
}

func (warlock *Warlock) registerCurseOfWeaknessSpell() {
	warlock.CurseOfWeaknessAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.CurseOfWeaknessAura(target)
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 702},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfWeakness,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  142,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfWeaknessAuras.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{warlock.CurseOfWeaknessAuras},
	})
}

func (warlock *Warlock) registerCurseOfTonguesSpell() {
	actionID := core.ActionID{SpellID: 1714}

	// Empty aura so we can simulate cost/time to keep tongues up
	warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Curse of Tongues",
			ActionID: actionID,
			Duration: time.Second * 30,
		})
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCurseOfTongues,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.04,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.CurseOfTonguesAuras.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{warlock.CurseOfTonguesAuras},
	})
}

func (warlock *Warlock) registerBaneOfAgonySpell() {
	baseTickDmg := warlock.ScalingBaseDamage * Coefficient_BaneOfAgony / 12.0

	warlock.BaneOfAgony = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 980},
		SpellSchool:    core.SpellSchoolShadow,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellBaneOfAgony,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Bane of Agony",
			},
			NumberOfTicks:    12,
			TickLength:       time.Second * 2,
			BonusCoefficient: 0.088,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 0.5*baseTickDmg)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				if dot.TickCount%4 == 0 { // CoA ramp up
					dot.SnapshotBaseDamage += 0.5 * baseTickDmg
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				warlock.BaneOfDoom.Dot(target).Cancel(sim)
				//TODO: Cancel BaneOfHavoc
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}

// TODO: Does this benefit from haunt?
func (warlock *Warlock) registerBaneOfDoomSpell() {

	ebonImpBonusSummon := 0.1 * float64(warlock.Talents.ImpendingDoom)

	warlock.BaneOfDoom = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 603},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellBaneOfDoom,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.15,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		FlatThreatBonus:          160,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Bane of Doom",
			},
			NumberOfTicks:    4,
			TickLength:       time.Second * 15,
			BonusCoefficient: 0.88,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, warlock.ScalingBaseDamage*Coefficient_BaneOfDoom)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				//TODO: Can this crit?
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				if sim.Proc(0.2+ebonImpBonusSummon, "Ebon Imp") {
					warlock.EbonImp.EnableWithTimeout(sim, warlock.EbonImp, time.Second*15)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.BaneOfAgony.Dot(target).Cancel(sim)
				//TODO: Cancel BaneOfHavoc
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
