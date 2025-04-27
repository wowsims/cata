package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerSealOfRighteousness() {
	var numTargets int32
	if paladin.Talents.SealsOfCommand {
		numTargets = paladin.Env.GetNumTargets()
	} else {
		numTargets = 1
	}
	results := make([]*core.SpellResult, numTargets)

	// Judgement of Righteousness cast on Judgement
	paladin.JudgementOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20187},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,
		ClassSpellMask: SpellMaskJudgementOfRighteousness,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1 +
				0.32*spell.SpellPower() +
				0.2*spell.MeleeAttackPower()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	registerOnHitSpell := func(tag int32, applyEffects core.ApplySpellResults) *core.Spell {
		return paladin.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: 25742}.WithTag(tag),
			SpellSchool:    core.SpellSchoolHoly,
			ProcMask:       core.ProcMaskMeleeProc,
			Flags:          core.SpellFlagMeleeMetrics,
			ClassSpellMask: SpellMaskSealOfRighteousness,

			MaxRange: core.MaxMeleeRange,

			DamageMultiplier: 1,
			CritMultiplier:   paladin.DefaultCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: applyEffects,
		})
	}

	// Seal of Righteousness on-hit proc (single hit, for Divine Storm)
	// Divine Storm is special, SoR can only proc once per target of DS, not like with
	// e.g. CS or a white hit where one hit will proc SoR on all surrounding targets in range.
	// Example for 10 targets:
	// CS hits 1 target -> SoR procs 10 times
	// DS hits 10 targets -> SoR procs 10 times
	// otherwise it would be DS hits 10 targets -> SoR procs 100 times
	onHitSingleTarget := registerOnHitSpell(1, func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := paladin.GetMHWeapon().SwingSpeed *
			(0.022*spell.SpellPower() + 0.011*spell.MeleeAttackPower())

		// can't miss if melee swing landed, but can crit
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
	})

	// Seal of Righteousness on-hit proc (multi-target hit, for everything else)
	onHitMultiTarget := registerOnHitSpell(2, func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := paladin.GetMHWeapon().SwingSpeed *
			(0.022*spell.SpellPower() + 0.011*spell.MeleeAttackPower())

		for idx := int32(0); idx < numTargets; idx++ {
			// can't miss if melee swing landed, but can crit
			results[idx] = spell.CalcDamage(sim, sim.Environment.GetTargetUnit(idx), baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		}

		for idx := int32(0); idx < numTargets; idx++ {
			spell.DealDamage(sim, results[idx])
		}
	})

	paladin.SealOfRighteousnessAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Righteousness" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 20154},
		Duration: time.Minute * 30,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses
			if !result.Landed() {
				return
			}

			// SoR only procs on white hits, CS, DS, TV, HoW and the melee part of HotR
			if spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				spell.ClassSpellMask&SpellMaskCanTriggerSealOfRighteousness == 0 {
				return
			}

			if spell.ClassSpellMask&SpellMaskDivineStorm != 0 {
				onHitSingleTarget.Cast(sim, result.Target)
			} else {
				onHitMultiTarget.Cast(sim, result.Target)
			}
		},
	})

	// Seal of Righteousness self-buff.
	aura := paladin.SealOfRighteousnessAura
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20154},
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
			paladin.CurrentJudgement = paladin.JudgementOfRighteousness
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
