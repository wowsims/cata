package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) RegisterSealOfTruth() {

	// Censure DoT
	censureSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31803, Tag: 2},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskCensure,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Censure (DoT)",
				MaxStacks: 5,
			},
			NumberOfTicks:        5,
			HasteAffectsDuration: true,
			TickLength:           time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				tickValue := 0 +
					.014*dot.Spell.SpellPower() +
					.027*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage = tickValue * float64(dot.GetStacks())

				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},
	})

	// Judegment of Truth cast on Judgement
	judgementDmg := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31804},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,
		ClassSpellMask: SpellMaskJudgement | SpellMaskSpecialAttack,

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
	paladin.CurrentJudgement = judgementDmg

	// Seal of Truth on-hit proc
	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42463},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc, // does proc certain spell damage-based items, e.g. Black Magic, Pendulum of Telluric Currents
		Flags:       core.SpellFlagMeleeMetrics,

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
		/*
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfSealOfTruth) {
					expertise := core.ExpertisePerQuarterPercentReduction * 10
					paladin.AddStatDynamic(sim, stats.Expertise, expertise)
				}
			},

			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfSealOfTruth) {
					expertise := core.ExpertisePerQuarterPercentReduction * 10
					paladin.AddStatDynamic(sim, stats.Expertise, -expertise)
				}
			},
		*/
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses or our own procs.
			if !result.Landed() || spell == censureSpell || spell == judgementDmg || spell == onSpecialOrSwingProc {
				return
			}

			if spell.IsMelee() {
				if censureSpell.Dot(result.Target).GetStacks() == 5 {
					onSpecialOrSwingProc.Cast(sim, result.Target)
				}
			}

			if spell.ClassSpellMask&SpellMaskSingleTarget == 0 {
				return
			}

			dotResult := censureSpell.CalcOutcome(sim, result.Target, spell.OutcomeMeleeSpecialHit)

			if dotResult.Landed() {
				dot := censureSpell.Dot(result.Target)
				if !dot.IsActive() {
					dot.Apply(sim)
				}
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, false)
				dot.Activate(sim)
			}
		},
	})

	// Seal of Truth self-buff.
	aura := paladin.SealOfTruthAura
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 31801},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSealOfTruth,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentJudgement = judgementDmg
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
