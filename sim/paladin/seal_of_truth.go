package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) registerSealOfTruth() {

	// Censure DoT
	censureSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31803},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: SpellMaskCensure,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				ActionID:  core.ActionID{SpellID: 31803},
				Label:     "Censure",
				MaxStacks: 5,
			},

			NumberOfTicks:       5,
			AffectedByCastSpeed: true,
			TickLength:          time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				tickValue := float64(dot.GetStacks()) * (0 +
					.014*dot.Spell.SpellPower() +
					.027*dot.Spell.MeleeAttackPower())

				dot.Snapshot(target, tickValue)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dotResult := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			spell.SpellMetrics[target.UnitIndex].Hits--

			if dotResult.Landed() {
				dot := spell.Dot(target)
				if dot.IsActive() {
					dot.AddStack(sim)
					dot.TakeSnapshot(sim, false)
					dot.Refresh(sim)
				} else {
					dot.Apply(sim)
					dot.SetStacks(sim, 1)
					dot.TakeSnapshot(sim, false)
				}
			}
		},
	})

	// Judgement of Truth cast on Judgement
	paladin.JudgementOfTruth = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31804},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,
		ClassSpellMask: SpellMaskJudgementOfTruth,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1 +
				.223*spell.SpellPower() +
				.142*spell.MeleeAttackPower()

			baseDamage *= 1 + .2*float64(censureSpell.Dot(target).GetStacks())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	// Seal of Truth on-hit proc
	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 42463},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskProc, // does proc certain spell damage-based items, e.g. Black Magic, Pendulum of Telluric Currents
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskSealOfTruth,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := .15 * paladin.MHWeaponDamage(sim, spell.MeleeAttackPower())

			// can't miss if melee swing landed, but can crit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	paladin.SealOfTruthAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Truth",
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 31801},
		Duration: time.Minute * 30,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses.
			if !result.Landed() {
				return
			}

			// SoT only procs on white hits, CS, TV, Exo, Judge and HoW
			if spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				spell.ClassSpellMask&SpellMaskCanTriggerSealOfTruth == 0 {
				return
			}

			if censureSpell.Dot(result.Target).GetStacks() == 5 {
				onSpecialOrSwingProc.Cast(sim, result.Target)
			}

			censureSpell.Cast(sim, result.Target)
		},
	})

	// Seal of Truth self-buff.
	aura := paladin.SealOfTruthAura
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31801},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentJudgement = paladin.JudgementOfTruth
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
