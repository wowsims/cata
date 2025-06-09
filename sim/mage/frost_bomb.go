package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerFrostBombSpell() {

	if !mage.Talents.FrostBomb {
		return
	}

	// Since Frost Bomb does double damage to all targets, these are the AOE values and the main target just gets double.
	frostBombExplosionCoefficient := 1.725 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=113092 Field "EffetBonusCoefficient"
	frostBombExplosionScaling := 2.21      // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=113092 Field "Coefficient"
	frostBombVariance := 0.0

	frostBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 113092},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellFrostBombExplosion,
		Flags:          core.SpellFlagAoE,

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultCritMultiplier(),
		BonusCoefficient:         frostBombExplosionCoefficient,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx, aoeTarget := range sim.Encounter.TargetUnits {
				if idx == 0 {
					spell.DamageMultiplier *= 2
				}
				baseDamage := mage.CalcAndRollDamageRange(sim, frostBombExplosionScaling, frostBombVariance)
				result := spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				spell.DealDamage(sim, result)
				if idx == 0 {
					spell.DamageMultiplier /= 2
				}
			}
		},
	})

	mage.FrostBombAuras = mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Frost Bomb",
			ActionID: core.ActionID{SpellID: 113092},
			Duration: time.Second * 4,
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				frostBombExplosionSpell.Cast(sim, aura.Unit)
				mage.WaitUntil(sim, sim.CurrentTime+mage.ReactionTime)
			},
		})
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 113092},
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
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				mage.FrostBombAuras.Get(target).Activate(sim)
			}
			spell.DealOutcome(sim, result)
		},
		RelatedAuraArrays: mage.FrostBombAuras.ToMap(),
	})

}
