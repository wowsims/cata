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
	shaman.ElementalBlastOverload[0] = shaman.RegisterSpell(shaman.newElementalBlastSpellConfig(true))
	shaman.ElementalBlastOverload[1] = shaman.RegisterSpell(shaman.newElementalBlastSpellConfig(true))
}

func (shaman *Shaman) newElementalBlastSpellConfig(isElementalOverload bool) core.SpellConfig {

	actionID := core.ActionID{SpellID: 118522}

	masteryAura := shaman.NewTemporaryStatsAura("Elemental Blast Mastery", actionID.WithTag(9), stats.Stats{stats.MasteryRating: 3500}, time.Second*8)
	hasteAura := shaman.NewTemporaryStatsAura("Elemental Blast Haste", actionID.WithTag(10), stats.Stats{stats.HasteRating: 3500}, time.Second*8)
	critAura := shaman.NewTemporaryStatsAura("Elemental Blast Crit", actionID.WithTag(11), stats.Stats{stats.CritRating: 3500}, time.Second*8)
	agiAura := shaman.NewTemporaryStatsAura("Elemental Blast Agi", actionID.WithTag(12), stats.Stats{stats.Agility: 3500}, time.Second*8)
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
		MetricSplits:   6,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Second * 2,
				GCD:      core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(shaman.MaelstromWeaponAura.GetStacks())
				castTime := shaman.ApplyCastSpeedForSpell(cast.CastTime, spell)
				if sim.CurrentTime+castTime > shaman.AutoAttacks.NextAttackAt() {
					shaman.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+castTime, false)
				}
			},
		},
		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: 2.11199998856,
		ThreatMultiplier: 1,
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
		result := shaman.calcDamageStormstrikeCritChance(sim, target, baseDamage, spell)

		idx := core.TernaryInt32(spell.Flags.Matches(SpellFlagIsEcho), 1, 0)
		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			if !isElementalOverload && result.Landed() && sim.Proc(shaman.GetOverloadChance(), "Elemental Blast Elemental Overload") {
				shaman.ElementalBlastOverload[idx].Cast(sim, target)
			}
			spell.DealDamage(sim, result)
		})
	}

	return spellConfig
}
