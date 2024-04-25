package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (shaman *Shaman) registerLavaBurstSpell() {
	shaman.LavaBurst = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(false))
	shaman.LavaBurstOverload = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(true))
}

func (shaman *Shaman) newLavaBurstSpellConfig(isElementalOverload bool) core.SpellConfig {
	actionID := core.ActionID{SpellID: 51505}

	mask := core.ProcMaskSpellDamage
	if isElementalOverload {
		mask = core.ProcMaskProc
	}
	flags := SpellFlagFocusable
	if !isElementalOverload {
		flags |= core.SpellFlagAPL
	}

	spellConfig := core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       mask,
		Flags:          flags,
		ClassSpellMask: core.TernaryInt64(isElementalOverload, SpellMaskLavaBurstOverload, SpellMaskLavaBurst),

		ManaCost: core.ManaCostOptions{
			BaseCost:   core.TernaryFloat64(isElementalOverload, 0, 0.1),
			Multiplier: 1 - 0.05*float64(shaman.Talents.Convection),
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
		CritMultiplier:   shaman.SpellCritMultiplier(1.0, core.TernaryFloat64(shaman.Spec == proto.Spec_SpecElementalShaman, 1.0, 0)+float64(shaman.Talents.LavaFlows)*0.08),
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
		result := spell.CalcDamage(sim, target, 1586, spell.OutcomeMagicHitAndCrit)

		if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance(), "Lava Burst Elemental Overload") {
			shaman.LavaBurstOverload.Cast(sim, target)
		}

		spell.DealDamage(sim, result)
	}

	return spellConfig
}
