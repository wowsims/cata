package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerFrostBomb() {

	if !mage.Talents.FrostBomb {
		return
	}

	// Since Frost Bomb does double damage to all targets, these are the AOE values and the main target just gets double.
	frostBombExplosionCoefficient := 1.725 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=113092 Field "EffetBonusCoefficient"
	frostBombExplosionScaling := 2.21      // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=113092 Field "Coefficient"
	frostBombVariance := 0.0
	actionID := core.ActionID{SpellID: 112948}

	frostBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(2), // Real Spell ID: 113092
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellFrostBombExplosion,
		Flags:          core.SpellFlagAoE,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: frostBombExplosionCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx, aoeTarget := range sim.Encounter.TargetUnits {
				if idx == 0 {
					spell.DamageMultiplier *= 2
				}
				baseDamage := mage.CalcAndRollDamageRange(sim, frostBombExplosionScaling, frostBombVariance)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				if idx == 0 {
					spell.DamageMultiplier /= 2
				}
			}
		},
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostBomb,
		ManaCost:       core.ManaCostOptions{BaseCostPercent: 1.25},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FrostBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					frostBombExplosionSpell.Cast(sim, aura.Unit)
					mage.WaitUntil(sim, sim.CurrentTime+mage.ReactionTime)
				},
			},
			NumberOfTicks:        4,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Empty onTick, we don't want to deal damage over time.
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			dot := spell.Dot(target)
			if result.Landed() {
				dot.Apply(sim)
				spell.CD.Set(sim.CurrentTime + mage.ApplyCastSpeedForSpell(spell.CD.Duration, spell))
			}
		},
	})
}
