package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerFrostBombSpell() {
	// Since Frost Bomb does double damage to all targets, these are the AOE values and the main target just gets double.
	var frostBombExplosionCoefficient = 1.725 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=113092 Field "EffetBonusCoefficient"
	var frostBombExplosionScaling = 2.21      // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=113092 Field "Coefficient"
	var frostBombVariance = 0.0
	var numTargets = mage.Env.GetNumTargets()

	frostBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 113092},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellFrostBomb,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultCritMultiplier(),
		BonusCoefficient:         frostBombExplosionCoefficient,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)
			baseDamage := mage.CalcAndRollDamageRange(sim, frostBombExplosionScaling, frostBombVariance)
			damage := 0.0
			for idx := int32(0); idx < numTargets; idx++ {
				if idx == 0 {
					damage = baseDamage * 2
				}
				results[idx] = spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit) //TODO How to tell if a spell is incapable of critting?
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	mage.FrostBomb = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44457},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostBombDot,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1.25,
		},
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
			NumberOfTicks:       1,
			TickLength:          time.Second * 4,
			AffectedByCastSpeed: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
