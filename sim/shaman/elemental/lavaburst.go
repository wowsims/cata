package elemental

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/shaman"
)

func (ele *ElementalShaman) registerLavaBurstSpell() {
	ele.LavaBurst = ele.RegisterSpell(ele.newLavaBurstSpellConfig(false))
	ele.LavaBurstOverload = ele.RegisterSpell(ele.newLavaBurstSpellConfig(true))
}

func (ele *ElementalShaman) newLavaBurstSpellConfig(isElementalOverload bool) core.SpellConfig {
	actionID := core.ActionID{SpellID: 51505}

	mask := core.ProcMaskSpellDamage
	flags := shaman.SpellFlagFocusable
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
		ClassSpellMask: core.TernaryInt64(isElementalOverload, shaman.SpellMaskLavaBurstOverload, shaman.SpellMaskLavaBurst),

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(isElementalOverload, 0, 7.7),
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 2000,
				GCD:      core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   ele.DefaultCritMultiplier(),
		BonusCoefficient: 1,
	}

	if isElementalOverload {
		spellConfig.ActionID.Tag = shaman.CastTagLightningOverload
		spellConfig.Cast.DefaultCast.CastTime = 0
		spellConfig.Cast.DefaultCast.GCD = 0
		spellConfig.Cast.DefaultCast.Cost = 0
		spellConfig.Cast.ModifyCast = nil
		spellConfig.MetricSplits = 0
		spellConfig.DamageMultiplier *= 0.75
		spellConfig.ThreatMultiplier = 0
	} else {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    ele.NewTimer(),
			Duration: time.Second * 8,
		}
	}

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := ele.CalcAndRollDamageRange(sim, 1.41624999046, 0.10000000149)
		if ele.FlameShock.Dot(target).IsActive() {
			baseDamage *= 1.5
		}
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellProc) { //So that procs from DTR does not cast an overload
				if !isElementalOverload && result.Landed() && sim.Proc(ele.GetOverloadChance(), "Lava Burst Elemental Overload") {
					ele.LavaBurstOverload.Cast(sim, target)
				}
			}

			spell.DealDamage(sim, result)
		})
	}

	return spellConfig
}
