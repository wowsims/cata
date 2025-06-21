package mage

import (
	"sort"
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerLivingBomb() {
	// MOP version has a cap of 3 active dots at once
	// activeLivingBombs should only ever be 3 LBs long
	// When a dot is trying to be applied,
	// 1) it should remove the dot with the longest remaining duration
	// 2) to do this, when a dot is applied, it checks the length of the array
	//   2a) if the array is longer than 3, remove the last element
	// 3) append the dot to the array
	// 4) sort the array by remaining duration, such that the longest remaining duration is LAST, to fit step 1
	// When a dot expires, remove the 1st element.

	if !mage.Talents.LivingBomb {
		return
	}

	actionID := core.ActionID{SpellID: 44457}
	activeLivingBombs := make([]*core.Dot, 0)
	const maxLivingBombs int = 3

	mage.RegisterResetEffect(func(s *core.Simulation) {
		activeLivingBombs = make([]*core.Dot, 0)
	})

	livingBombExplosionCoefficient := 0.08 // Per https://wago.tools/db2/SpellEffect?build=4.4.2.60192&filter%5BSpellID%5D=44461 Field "EffetBonusCoefficient"
	livingBombExplosionScaling := 0.103    // Per https://wago.tools/db2/SpellEffect?build=4.4.2.60192&filter%5BSpellID%5D=44461 Field "Coefficient"
	livingBombDotCoefficient := 0.8        // Per https://wago.tools/db2/SpellEffect?build=4.4.2.60192&filter%5BSpellID%5D=44461 Field "EffetBonusCoefficient"
	livingBombDotScaling := 1.03           // Per https://wago.tools/db2/SpellEffect?build=4.4.2.60192&filter%5BSpellID%5D=44461 Field "Coefficient"

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(2), // Real Spell ID: 44461
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellLivingBombExplosion,
		Flags:          core.SpellFlagAoE | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: livingBombExplosionCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mage.CalcAndRollDamageRange(sim, livingBombExplosionScaling, 0)
			ticks := max(4, float64(mage.LivingBomb.RelatedDotSpell.Dot(target).Duration)/float64(mage.LivingBomb.RelatedDotSpell.Dot(target).TickPeriod()))
			spell.DamageMultiplier *= ticks
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			spell.DamageMultiplier /= ticks
		},
	})

	bombExplode := true

	mage.LivingBomb = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellLivingBombApply,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1.5,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		CritMultiplier:   mage.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				dot := spell.RelatedDotSpell.Dot(target)
				// If there is already an active dot on the target, just reapply
				if dot.IsActive() {
					spell.RelatedDotSpell.Cast(sim, target)
				} else {
					activeLbs := len(activeLivingBombs)

					if activeLbs >= maxLivingBombs {
						bombExplode = false
						activeLivingBombs[activeLbs-1].Deactivate(sim)
						if activeLbs != 0 {
							activeLivingBombs = activeLivingBombs[:1]
						}
						bombExplode = true
					}
					spell.RelatedDotSpell.Cast(sim, target)
					dot = spell.RelatedDotSpell.Dot(target)
					activeLivingBombs = append(activeLivingBombs, dot)
					sort.Slice(activeLivingBombs, func(i, j int) bool {
						return activeLivingBombs[i].Duration < activeLivingBombs[j].Duration
					})
				}
			}
			spell.DealOutcome(sim, result)
		},
	})

	mage.LivingBomb.RelatedDotSpell = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellLivingBombDot,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if bombExplode {
						livingBombExplosionSpell.Cast(sim, aura.Unit)
						mage.WaitUntil(sim, sim.CurrentTime+mage.ReactionTime)
						if len(activeLivingBombs) != 0 {
							activeLivingBombs = activeLivingBombs[1:]
						}
					}
				},
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,
			BonusCoefficient:    livingBombDotCoefficient,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, mage.CalcScalingSpellDmg(livingBombDotScaling))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			// The Bomb goes ot if the target has a Dot that has <= 1 tick remaining.
			if dot.IsActive() && dot.RemainingTicks() == 1 {
				livingBombExplosionSpell.Cast(sim, target)
			}
			dot.Apply(sim)
		},
	})
}
