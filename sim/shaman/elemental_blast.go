package shaman

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
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
	agiAura := shaman.NewTemporaryStatsAura("Elemental Blast Agi", core.ActionID{SpellID: 118522, Tag: 4}, stats.Stats{stats.Agility: 3500}, time.Second*8)
	eleBlastAuras := []*core.StatBuffAura{masteryAura, hasteAura, critAura, agiAura}

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
		ClassSpellMask: core.TernaryInt64(isElementalOverload, SpellMaskElementalBlastOverload, SpellMaskElementalBlast),
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
		if !isElementalOverload {
			var rand int
			if shaman.Spec == proto.Spec_SpecEnhancementShaman {
				rand = int(math.Floor(sim.RollWithLabel(0, 4, "Elemental Blast buff")))
			} else {
				rand = int(math.Floor(sim.RollWithLabel(0, 3, "Elemental Blast buff")))
			}
			for i, aura := range eleBlastAuras {
				if i == rand {
					aura.Activate(sim)
				} else {
					aura.Deactivate(sim)
				}
			}
		}

		baseDamage := shaman.CalcAndRollDamageRange(sim, 4.23999977112, 0.15000000596)
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

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
