package cata

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {

	// Keep these in order by item ID

	// Fleet Shadowspirit Diamond
	core.NewItemEffect(52289, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.NewMovementSpeedAura("Minor Run Speed", core.ActionID{SpellID: 13889}, 0.08)
	})

	// Bracing Shadowspirit Diamond
	core.NewItemEffect(52292, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	// Eternal Shadowspirit Diamond
	core.NewItemEffect(52293, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BlockDamageReduction += 0.01
	})

	// Austere Shadowspirit Diamond
	core.NewItemEffect(52294, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.ApplyEquipScaling(stats.Armor, 1.02)
	})

	// Effulgent Shadowspirit Diamond
	core.NewItemEffect(52295, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 0.98
		character.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 0.98
	})

	// Ember Shadowspirit Diamond
	core.NewItemEffect(52296, func(agent core.Agent) {
		agent.GetCharacter().MultiplyStat(stats.Mana, 1.02)
	})

	// Chaotic Shadowspirit Diamond
	core.NewItemEffect(52291, core.ApplyMetaGemCriticalDamageEffect)
	// Revitalizing Shadowspirit Diamond
	core.NewItemEffect(52297, core.ApplyMetaGemCriticalDamageEffect)
	// Agile Shadowspirit Diamond
	core.NewItemEffect(68778, core.ApplyMetaGemCriticalDamageEffect)
	// Reverberating Shadowspirit Diamond
	core.NewItemEffect(68779, core.ApplyMetaGemCriticalDamageEffect)
	// Burning Shadowspirit Diamond
	core.NewItemEffect(68780, core.ApplyMetaGemCriticalDamageEffect)
}
