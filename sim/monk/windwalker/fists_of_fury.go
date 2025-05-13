package windwalker

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

/*
Tooltip:
Pummel all targets in front of you with rapid hand strikes, stunning them and dealing ${7.5*0.89*$<low>} to ${7.5*0.89*$<high>} damage immediately and every 1 sec for 4 sec.
Damage is spread evenly over all targets.

-- Glyph of Fists of Fury --
Your parry chance is increased by 100% while channeling.
-- Glyph of Fists of Fury --
*/
func (ww *WindwalkerMonk) registerFistsOfFury() {
	actionID := core.ActionID{SpellID: 113656}
	debuffActionID := core.ActionID{SpellID: 117418}
	chiMetrics := ww.NewChiMetrics(actionID)
	numTargets := ww.Env.GetNumTargets()

	fistsOfFuryTickSpell := ww.RegisterSpell(core.SpellConfig{
		ActionID:       debuffActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: monk.MonkSpellFistsOfFury,
		MaxRange:       core.MaxMeleeRange,

		DamageMultiplier: 7.5 * 0.89,
		ThreatMultiplier: 1,
		CritMultiplier:   ww.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)
			baseDamage := ww.CalculateMonkStrikeDamage(sim, spell)

			// Damage is split between all mobs, each hit rolls for hit/crit separately
			baseDamage /= float64(numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				result := spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				results[idx] = result
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	fistsOfFuryBuff := ww.RegisterAura(core.Aura{
		Label:    "Fists of Fury" + ww.Label,
		ActionID: actionID,
		Duration: time.Second * 3,
	})

	ww.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagChanneled | monk.SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: monk.MonkSpellFistsOfFury,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    ww.NewTimer(),
				Duration: time.Second * 25,
			},
		},

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Fists of Fury" + ww.Label,
				ActionID: debuffActionID,
			},
			NumberOfTicks:        4,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				fistsOfFuryTickSpell.Cast(sim, target)
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !ww.Moving && ww.ComboPoints() >= 3
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ww.SpendChi(sim, 3, chiMetrics)

			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)

			expiresAt := dot.ExpiresAt()
			ww.AutoAttacks.DelayMeleeBy(sim, expiresAt-sim.CurrentTime)
			ww.ExtendGCDUntil(sim, expiresAt+ww.ReactionTime)

			remainingDuration := dot.RemainingDuration(sim)
			fistsOfFuryBuff.Duration = remainingDuration
			fistsOfFuryBuff.Activate(sim)
		},
	})
}
