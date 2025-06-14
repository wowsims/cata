package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Fills you with Holy Light, causing melee attacks to deal 12% additional weapon damage as Holy and apply Censure to the target.
Replaces Seal of Command.

Censure
Deals

-- Ardent Defender --
108 + 0.094 * <SP>
-- else --
108 * 5 + (0.094 * <SP>)
----------

additional Holy damage over 15 sec. Stacks up to 5 times.
*/
func (paladin *Paladin) registerSealOfTruth() {
	censureActionId := core.ActionID{SpellID: 31803}

	// Censure DoT application
	censureSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    censureActionId.WithTag(1),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc,
		Flags:       core.SpellFlagNoMetrics | core.SpellFlagNoLogs | core.SpellFlagPassiveSpell,

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
		CritMultiplier:   paladin.DefaultCritMultiplier(),
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
				tickValue := paladin.CalcScalingSpellDmg(0.09399999678) + 0.09399999678*dot.Spell.SpellPower()
				tickValue *= float64(dot.GetStacks())
				dot.SnapshotPhysical(target, tickValue)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Casts--
			dot := spell.Dot(target)

			dot.Apply(sim)
			dot.AddStack(sim)
			dot.TakeSnapshot(sim, false)
		},
	})

	// Seal of Truth on-hit proc
	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 42463},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskSealOfTruth,

		DamageMultiplier: 0.12,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHWeaponDamage(sim, spell.MeleeAttackPower())

			// can't miss if melee swing landed, but can crit
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	paladin.SealOfTruthAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Truth" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 31801},
		Duration: core.NeverExpires,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses.
			if !result.Landed() {
				return
			}

			// SoT only procs on white hits, CS, TV, Exo, Judge, HoW, HotR, ShoR
			if spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				!spell.Matches(SpellMaskCanTriggerSealOfTruth) {
				return
			}

			censureSpell.Cast(sim, result.Target)
			onSpecialOrSwingProc.Cast(sim, result.Target)
		},
	})

	// Seal of Truth self-buff.
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31801},
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
			paladin.CurrentSeal = paladin.SealOfTruthAura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
