package rogue

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeLeather, 87504)
	rogue.PseudoStats.MeleeSpeedMultiplier *= []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.LightningReflexes]
	rogue.AddStat(stats.PhysicalHitPercent, 2*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.SpellHitPercent, 2*float64(rogue.Talents.Precision))

	if rogue.Talents.SavageCombat > 0 {
		rogue.MultiplyStat(stats.AttackPower, 1.0+0.03*float64(rogue.Talents.SavageCombat))
	}

	if rogue.Talents.Ruthlessness > 0 {
		rogue.ruthlessnessMetrics = rogue.NewComboPointMetrics(core.ActionID{SpellID: 14161})
	}

	if rogue.Talents.RelentlessStrikes > 0 {
		rogue.relentlessStrikesMetrics = rogue.NewEnergyMetrics(core.ActionID{SpellID: 14179})
	}

	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade) {
		rogue.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			IntValue:  -15,
			ClassMask: RogueSpellTricksOfTheTrade,
		})
	}

	if rogue.Talents.SlaughterFromTheShadows > 0 {
		rogue.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			IntValue:  []int32{-0, -7, -14, -20}[rogue.Talents.SlaughterFromTheShadows],
			ClassMask: RogueSpellBackstab | RogueSpellAmbush,
		})

		rogue.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			IntValue:  -2 * rogue.Talents.SlaughterFromTheShadows,
			ClassMask: RogueSpellHemorrhage | RogueSpellFanOfKnives,
		})
	}

	if rogue.Talents.Quickening > 0 {
		spellID := []int32{0, 24090, 31209}[rogue.Talents.Quickening]
		multiplier := []float64{0, 0.08, 0.15}[rogue.Talents.Quickening]
		rogue.NewMovementSpeedAura("Quickening", core.ActionID{SpellID: spellID}, multiplier)
	}
}

// DWSMultiplier returns the offhand damage multiplier
func (rogue *Rogue) DWSMultiplier() float64 {
	// DWS (Now named Ambidexterity) is now a Combat rogue passive
	return core.TernaryFloat64(rogue.Spec == proto.Spec_SpecCombatRogue, 1.75, 1)
}
