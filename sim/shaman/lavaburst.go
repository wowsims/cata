package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) registerLavaBurstSpell() {
	shaman.LavaBurst = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(false))
	shaman.LavaBurstOverload = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(true))
}

func (shaman *Shaman) newLavaBurstSpellConfig(isElementalOverload bool) core.SpellConfig {
	actionID := core.ActionID{SpellID: 51505}

	mask := core.ProcMaskSpellDamage
	flags := SpellFlagFocusable
	if isElementalOverload {
		mask = core.ProcMaskSpellProc
		flags |= core.SpellFlagPassiveSpell
	} else {
		flags |= core.SpellFlagAPL
	}

	spellConfig := core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       mask,
		Flags:          flags,
		MissileSpeed:   24,
		ClassSpellMask: core.TernaryInt64(isElementalOverload, SpellMaskLavaBurstOverload, SpellMaskLavaBurst),

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryInt32(isElementalOverload, 0, 10),
			PercentModifier: 100 - (5 * shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 2000,
				GCD:      core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				castTime := shaman.ApplyCastSpeedForSpell(cast.CastTime, spell)
				if castTime > 0 {
					shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
				}
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.628,
	}

	if isElementalOverload {
		spellConfig.ActionID.Tag = CastTagLightningOverload
		spellConfig.Cast.DefaultCast.CastTime = 0
		spellConfig.Cast.DefaultCast.GCD = 0
		spellConfig.Cast.DefaultCast.Cost = 0
		spellConfig.Cast.ModifyCast = nil
		spellConfig.MetricSplits = 0
		spellConfig.DamageMultiplier *= 0.75
		spellConfig.ThreatMultiplier = 0
	} else {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    shaman.NewTimer(),
			Duration: time.Second * 8,
		}
	}

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 1.57899999619, 0.24199999869)
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellProc) { //So that procs from DTR does not cast an overload
				if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance(), "Lava Burst Elemental Overload") {
					shaman.LavaBurstOverload.Cast(sim, target)
				}
			}

			spell.DealDamage(sim, result)
		})
	}

	return spellConfig
}
