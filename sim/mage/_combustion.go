package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (mage *Mage) registerCombustionSpell() {
	if !mage.Talents.Combustion {
		return
	}

	actionID := core.ActionID{SpellID: 11129}

	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage, // need to check proc mask for impact damage
		ClassSpellMask: MageSpellCombustionApplication,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultCritMultiplier(),
		BonusCoefficient:         1.113,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429 * mage.ClassSpellScaling
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.DealDamage(sim, result)
				spell.RelatedDotSpell.Cast(sim, target)
			}
			if mage.T13_4pc.IsActive() && spell.ProcMask&core.ProcMaskSpellProc == 0 {
				spell.CD.Reduce(time.Second * time.Duration(5*mage.t13ProcAura.GetStacks()))
			}
		},
	})

	dotBase := map[int64]float64{
		MageSpellLivingBombDot: 0.25 * mage.ClassSpellScaling,
		MageSpellPyroblastDot:  0.175 * mage.ClassSpellScaling,
	}

	calculatedDotTick := func(target *core.Unit) float64 {
		tickDamage := 0.0
		dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite, mage.Pyroblast.RelatedDotSpell}
		for _, spell := range dotSpells {
			dot := spell.Dot(target)
			if dot.IsActive() {
				if spell.ClassSpellMask&(MageSpellLivingBombDot|MageSpellPyroblastDot) != 0 {
					dps := dotBase[spell.ClassSpellMask] + dot.BonusCoefficient*dot.Spell.SpellPower()
					dps *= spell.DamageMultiplier * spell.DamageMultiplierAdditive
					tickDamage += dps / dot.BaseTickLength.Seconds()
				} else {
					tickDamage += dot.SnapshotBaseDamage / 2
				}
			}
		}
		return tickDamage
	}

	mage.Combustion.RelatedDotSpell = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: MageSpellCombustion,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Combustion Dot",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if mage.t13ProcAura != nil {
						mage.t13ProcAura.Deactivate(sim)
					}
				},
			},
			NumberOfTicks:       10,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				tickBase := calculatedDotTick(target)
				dot.Snapshot(target, tickBase)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			tickBase := calculatedDotTick(target)
			result := spell.CalcPeriodicDamage(sim, target, tickBase, spell.OutcomeExpectedMagicAlwaysHit)

			critChance := spell.SpellCritChance(target)
			critMod := (critChance * (spell.CritMultiplier - 1))
			result.Damage *= 1 + critMod

			return result
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	combustionTickCount := 0
	combustionTickDamage := 0.0
	updateCombustionTickCountEstimate := func() {
		combustionTickCount = int(mage.Combustion.RelatedDotSpell.Dot(mage.CurrentTarget).ExpectedTickCount())
	}
	updateCombustionTickDamageEstimate := func(sim *core.Simulation) {
		combustionTickDamage = mage.Combustion.RelatedDotSpell.ExpectedTickDamage(sim, mage.CurrentTarget)
	}

	updateCombustionTotalDamageEstimate := func() {
		combustionDotDamage := int32(float64(combustionTickCount) * combustionTickDamage)
		mage.combustionDotEstimate = combustionDotDamage
	}

	mage.AddOnCastSpeedChanged(func(old float64, new float64) {
		updateCombustionTickCountEstimate()
		updateCombustionTotalDamageEstimate()
	})

	mage.AddOnTemporaryStatsChange(func(sim *core.Simulation, _ *core.Aura, stats stats.Stats) {
		updateCombustionTickDamageEstimate(sim)
		updateCombustionTotalDamageEstimate()
	})

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Combustion Dot Calculation Listener",
		ClassSpellMask: MageSpellIgnite | MageSpellPyroblastDot | MageSpellLivingBombDot,
		Callback:       core.CallbackOnCastComplete | core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			updateCombustionTickDamageEstimate(sim)
			updateCombustionTotalDamageEstimate()
		},
	})
}
