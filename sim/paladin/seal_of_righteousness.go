package paladin

import (
	"github.com/wowsims/mop/sim/core"
)

// Fills you with Holy Light, causing melee attacks to deal 9% weapon damage to all targets within 8 yards.
func (paladin *Paladin) registerSealOfRighteousness() {
	numTargets := paladin.Env.GetNumTargets()

	registerOnHitSpell := func(tag int32, applyEffects core.ApplySpellResults) *core.Spell {
		return paladin.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: 101423}.WithTag(tag),
			SpellSchool:    core.SpellSchoolHoly,
			ProcMask:       core.ProcMaskMeleeProc,
			Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagAoE,
			ClassSpellMask: SpellMaskSealOfRighteousness,

			MaxRange: 8,

			DamageMultiplier: 0.09,
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
		baseDamage := paladin.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

		// can't miss if melee swing landed, but can crit
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
	})

	// Seal of Righteousness on-hit proc (multi-target hit, for everything else)
	onHitMultiTarget := registerOnHitSpell(2, func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		results := make([]*core.SpellResult, numTargets)

		for idx := range numTargets {
			currentTarget := sim.Environment.GetTargetUnit(idx)
			baseDamage := paladin.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			// can't miss if melee swing landed, but can crit
			results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		}

		for idx := range numTargets {
			spell.DealDamage(sim, results[idx])
		}
	})

	paladin.SealOfRighteousnessAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Righteousness" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 20154},
		Duration: core.NeverExpires,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			divineStorm := spell.Matches(SpellMaskDivineStorm)

			// Don't proc on misses, **except for Divine Storm**
			if !result.Landed() && !divineStorm {
				return
			}

			// SoR only procs on white hits, CS, DS, TV, ShotR and the melee part of HotR
			if spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				!spell.Matches(SpellMaskCanTriggerSealOfRighteousness) {
				return
			}

			if divineStorm {
				onHitSingleTarget.Cast(sim, result.Target)
			} else {
				onHitMultiTarget.Cast(sim, result.Target)
			}
		},
	})

	// Seal of Righteousness self-buff.
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20154},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 16.4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = paladin.SealOfRighteousnessAura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
