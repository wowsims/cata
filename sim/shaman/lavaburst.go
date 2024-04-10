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
	castTime := time.Millisecond * 2000
	spellCoeff := 0.628
	canOverload := false
	overloadChance := shaman.GetOverloadChance()
	if shaman.Spec == proto.Spec_SpecElementalShaman {
		//apply shamanism bonuses
		castTime -= 500
		spellCoeff += 0.36
		canOverload = true
	}

	actionID := core.ActionID{SpellID: 51505}
	spellCoeff += core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph(proto.ShamanPrimeGlyph_GlyphOfLavaBurst)), 0.1, 0)

	mask := core.ProcMaskSpellDamage
	if isElementalOverload {
		mask = core.ProcMaskProc
	}
	flags := SpellFlagElectric | SpellFlagFocusable
	if !isElementalOverload {
		flags |= core.SpellFlagAPL
	}

	spellConfig := core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    mask,
		Flags:       flags,

		ManaCost: core.ManaCostOptions{
			BaseCost:   core.TernaryFloat64(isElementalOverload, 0, 0.1),
			Multiplier: 1 - 0.05*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplier: 1 + 0.02*float64(shaman.Talents.Concussion) + 0.05*float64(shaman.Talents.CallOfFlame),
		CritMultiplier:   shaman.ElementalFuryCritMultiplier(0.08 * float64(shaman.Talents.LavaFlows)),
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
	}

	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 1586 + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		if canOverload && result.Landed() && sim.RandomFloat("Lava Burst Elemental Overload") < overloadChance {
			shaman.LavaBurstOverload.Cast(sim, target)
		}

		spell.DealDamage(sim, result)
	}

	return spellConfig
}
