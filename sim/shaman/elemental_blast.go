package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (shaman *Shaman) registerElementalBlastSpell() {
	shaman.ElementalBlast = shaman.RegisterSpell(shaman.newElementalBlastSpellConfig(false))
	shaman.ElementalBlastOverload = shaman.RegisterSpell(shaman.newElementalBlastSpellConfig(true))
}

func (shaman *Shaman) newElementalBlastSpellConfig(isElementalOverload bool) core.SpellConfig {

	masteryAura := shaman.NewTemporaryStatsAura("Elemental Blast Mastery", core.ActionID{SpellID: 118522, Tag: 1}, stats.Stats{stats.MasteryRating: 3500}, time.Second*8)
	hasteAura := shaman.NewTemporaryStatsAura("Elemental Blast Haste", core.ActionID{SpellID: 118522, Tag: 2}, stats.Stats{stats.HasteRating: 3500}, time.Second*8)
	critAura := shaman.NewTemporaryStatsAura("Elemental Blast Crit", core.ActionID{SpellID: 118522, Tag: 3}, stats.Stats{stats.CritRating: 3500}, time.Second*8)

	mask := core.ProcMaskSpellDamage
	flags := SpellFlagShamanSpell | SpellFlagFocusable
	if isElementalOverload {
		mask = core.ProcMaskSpellProc
		flags |= core.SpellFlagPassiveSpell
	} else {
		flags |= core.SpellFlagAPL
	}

	spellConfig := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 117014},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature,
		ProcMask:       mask,
		Flags:          flags,
		MissileSpeed:   40,
		ClassSpellMask: SpellMaskElementalBlast,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Second * 2,
				GCD:      core.GCDDefault,
			},
		},
		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: 2.11199998856,
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
			Duration: time.Second * 12,
		}
	}

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := shaman.CalcAndRollDamageRange(sim, 4.23999977112, 0.15000000596)
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		if !isElementalOverload {
			rand := sim.RandomFloat("Elemental Blast buff")
			if rand < 1.0/3.0 {
				hasteAura.Deactivate(sim)
				critAura.Deactivate(sim)
				masteryAura.Activate(sim)
			} else {
				if rand < 2.0/3.0 {
					masteryAura.Deactivate(sim)
					critAura.Deactivate(sim)
					hasteAura.Activate(sim)
				} else {
					masteryAura.Deactivate(sim)
					hasteAura.Deactivate(sim)
					critAura.Activate(sim)
				}
			}
		}

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellProc) { //So that procs from DTR does not cast an overload
				if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance(), "Elemental Blast Elemental Overload") {
					shaman.ElementalBlastOverload.Cast(sim, target)
				}
			}
			spell.DealDamage(sim, result)
		})
	}

	return spellConfig
}
