package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerSealOfTruth() {
	hasteMultiplier := 1 + 0.01*3*float64(paladin.Talents.JudgementsOfThePure)

	censureActionId := core.ActionID{SpellID: 31803}

	// Censure DoT application
	censureSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    censureActionId.WithTag(1),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc,
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dotResult := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)

			if dotResult.Landed() {
				spell.RelatedDotSpell.Cast(sim, target)
			}
		},
	})

	// Censure DoT
	censureSpell.RelatedDotSpell = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       censureActionId.WithTag(2),
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskCensure,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Censure (DoT)" + paladin.Label,
				MaxStacks: 5,
			},

			NumberOfTicks:       5,
			AffectedByCastSpeed: true,
			TickLength:          time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				tickValue := float64(dot.GetStacks()) * (0 +
					0.01400000043*dot.Spell.SpellPower() +
					0.0270000007*dot.Spell.MeleeAttackPower())

				dot.SnapshotBaseDamage = tickValue
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Casts--
			dot := spell.Dot(target)

			undoJotpForInitialTick := !dot.IsActive() &&
				paladin.JudgementsOfThePureAura != nil &&
				paladin.JudgementsOfThePureAura.IsActive()

			if undoJotpForInitialTick {
				paladin.MultiplyCastSpeed(1 / hasteMultiplier)
			}

			dot.Apply(sim)
			dot.AddStack(sim)

			if undoJotpForInitialTick {
				paladin.MultiplyCastSpeed(hasteMultiplier)
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
				0.22300000489*spell.SpellPower() +
				0.14200000465*spell.MeleeAttackPower()

			baseDamage *= 1 + .2*float64(censureSpell.Dot(target).GetStacks())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	// Seal of Truth on-hit proc
	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 42463},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskSealOfTruth,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.15 * paladin.MHWeaponDamage(sim, spell.MeleeAttackPower()) *
				(0.2 * float64(censureSpell.Dot(target).GetStacks()))

			// can't miss if melee swing landed, but can crit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	paladin.SealOfTruthAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Truth" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 31801},
		Duration: time.Minute * 30,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses.
			if !result.Landed() {
				return
			}

			// SoT only procs on white hits, CS, TV, Exo, Judge, HoW, HotR, ShoR
			if spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				spell.ClassSpellMask&SpellMaskCanTriggerSealOfTruth == 0 {
				return
			}

			censureSpell.Cast(sim, result.Target)
			onSpecialOrSwingProc.Cast(sim, result.Target)
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
			BaseCostPercent: 14,
			PercentModifier: 100,
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
