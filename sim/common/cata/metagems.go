package cata

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {

	// Keep these in order by item ID

	// Fleet Shadowspirit Diamond
	// PLACEHODLER - Until we implement movement speed
	// core.NewItemEffect(52289, func(agent core.Agent) {
	// 	agent.GetCharacter().MultiplyStat(Stats.Speed, 1.09)
	// })

	// Bracing Shadowspirit Diamond
	core.NewItemEffect(52292, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	// Eternal Shadowspirit Diamond
	core.NewItemEffect(52293, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.BlockValueMultiplier += 0.01
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

	// These are handled in character.go, but create empty effects so they are included in tests.
	core.NewItemEffect(41285, func(_ core.Agent) {}) // Chaotic Skyflare Diamond
	core.NewItemEffect(41376, func(_ core.Agent) {}) // Revitalizing Skyflare Diamond
	core.NewItemEffect(41398, func(_ core.Agent) {}) // Relentless Earthsiege Diamond
	core.NewItemEffect(52291, func(_ core.Agent) {}) // Chaotic Shadowspirit Diamond
	core.NewItemEffect(52297, func(_ core.Agent) {}) // Revitalizing Shadowspirit Diamond
	core.NewItemEffect(68778, func(_ core.Agent) {}) // Agile Shadowspirit Diamond
	core.NewItemEffect(68779, func(_ core.Agent) {}) // Reverberating Shadowspirit Diamond
	core.NewItemEffect(68780, func(_ core.Agent) {}) // Burning Shadowspirit Diamond
}
